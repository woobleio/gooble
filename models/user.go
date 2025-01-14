package model

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"io"
	"strconv"
	"strings"
	"wooble/lib"

	gomail "gopkg.in/gomail.v2"
)

// User is a Wooble user
type User struct {
	ID         uint64         `json:"-" db:"user.id"`
	CustomerID string         `json:"-" db:"customer_id"`
	AccountID  lib.NullString `json:"-" db:"account_id"`

	Email    string          `json:"email,omitempty" db:"email"`
	Name     string          `json:"name,omitempty" db:"name"`
	PicPath  *lib.NullString `json:"profilePath,omitempty" db:"pic_path"`
	Fullname *lib.NullString `json:"fullname,omitempty" db:"fullname"`
	IsVIP    bool            `json:"-" db:"is_vip"`
	IsActive bool            `json:"-" db:"is_active"`

	Website      *lib.NullString `json:"website,omitempty" db:"website"`
	CodepenName  *lib.NullString `json:"codepenName,omitempty" db:"codepen_name"`
	DribbbleName *lib.NullString `json:"dribbbleName,omitempty" db:"dribbble_name"`
	GithubName   *lib.NullString `json:"githubName,omitempty" db:"github_name"`
	TwitterName  *lib.NullString `json:"twitterName,omitempty" db:"twitter_name"`

	PlanID   *lib.NullInt64 `json:"-" db:"plan_user.id"`
	Plan     *Plan          `json:"plan,omitempty" db:""`
	Packages *[]Package     `json:"packages,omitempty" db:""`

	IsCreator bool   `json:"isCreator,omitempty" db:"is_creator"`
	Fund      uint64 `json:"fund,omitempty" db:"fund"`

	Secret string `json:"-" db:"passwd"`
	Salt   string `json:"-" db:"salt_key"`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"user.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"user.updated_at"`
	DeletedAt *lib.NullTime `json:"deletedAt,omitempty" db:"user.deleted_at"`
}

// AllActiveUsers returns all public and activated users
func AllActiveUsers() ([]User, error) {
	var users []User
	q := `
	SELECT
		u.name,
		u.fullname,
		u.is_creator,
		u.pic_path,
		u.website,
		u.codepen_name,
		u.dribbble_name,
		u.github_name,
		u.twitter_name,
		u.created_at "user.created_at"
	FROM app_user u
	WHERE u.deleted_at IS NULL AND u.is_active = TRUE
	`
	return users, lib.DB.Select(&users, q)
}

// UserPublicByName returns user public profile with the name "username"
func UserPublicByName(username string) (*User, error) {
	var user User
	q := `
		SELECT
			u.name,
			u.fullname,
			u.is_creator,
			u.pic_path,
			u.website,
			u.codepen_name,
			u.dribbble_name,
			u.github_name,
			u.twitter_name,
			u.created_at "user.created_at"
    FROM app_user u
		WHERE u.name = $1
		AND u.deleted_at IS NULL
	`
	return &user, lib.DB.Get(&user, q, username)
}

// UserPrivateByID returns user with id "id"
// It'll select the most recent plan but ignore it if the end_date expired
func UserPrivateByID(userID uint64) (*User, error) {
	var user User
	q := `SELECT DISTINCT ON (u.id)
			u.id "user.id",
			u.email,
			u.name,
			u.fullname,
			u.is_active,
			u.pic_path,
			u.website,
			u.codepen_name,
			u.dribbble_name,
			u.github_name,
			u.twitter_name,
			u.is_creator,
			u.created_at "user.created_at",
			u.updated_at "user.updated_at",
			u.salt_key,
      u.customer_id,
			u.account_id,
			pu.id "plan_user.id",
      pu.start_date,
      pu.end_date,
			pu.unsub_date,
      pl.label "plan.label",
			pl.level,
      pl.nb_pkg,
      pl.nb_crea
    FROM app_user u
    LEFT OUTER JOIN plan_user pu ON (pu.user_id = u.id AND pu.end_date >= now())
    LEFT OUTER JOIN plan pl ON (pl.label = pu.plan_label)
		WHERE u.id = $1
		AND u.deleted_at IS NULL
    ORDER BY u.id, pu.id DESC`

	if err := lib.DB.Get(&user, q, userID); err != nil {
		return nil, err
	}

	if user.Plan.Label == nil {
		user.Plan, _ = DefaultPlan(userID)
	}

	return &user, nil
}

// NewUser creates a new user
func NewUser(user *User) (uID uint64, err error) {
	salt := lib.GenKey()
	cp, errPasswd := lib.Encrypt(user.Secret, []byte(salt))
	if errPasswd != nil {
		return 0, errPasswd
	}
	q := `INSERT INTO app_user(name, fullname, email, is_creator, passwd, salt_key, customer_id, is_vip) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err = lib.DB.QueryRow(q, user.Name, user.Fullname, user.Email, user.IsCreator, cp, salt, user.CustomerID, user.IsVIP).Scan(&uID)

	return uID, err
}

// UserByEmail returns user with a specified email or name
func UserByEmail(email string) (*User, error) {
	var user User
	q := `
    SELECT DISTINCT ON (u.id)
      u.id "user.id",
      u.email,
      u.name,
			u.fullname,
			u.is_active,
      u.passwd,
      u.is_creator,
      u.salt_key,
      u.customer_id,
      pu.start_date,
      pu.end_date,
      pl.label "plan.label",
      pl.nb_pkg,
      pl.nb_crea
    FROM app_user u
    LEFT OUTER JOIN plan_user pu ON (pu.user_id = u.id)
    LEFT OUTER JOIN plan pl ON (pl.label = pu.plan_label AND pu.end_date >= now())
    WHERE u.email = $1
		AND u.deleted_at IS NULL
    ORDER BY u.id, pu.id DESC
	`

	if err := lib.DB.Get(&user, q, email); err != nil {
		return nil, err
	}

	if user.Plan == nil {
		user.Plan, _ = DefaultPlan(user.ID)
	}

	return &user, lib.DB.Get(&user, q, email)
}

// UserCustomerID returns customerID of user "uID"
func UserCustomerID(uID uint64) (string, error) {
	var user User
	q := `SELECT customer_id FROM app_user WHERE app_user.id = $1 AND deleted_at IS NULL`
	return user.CustomerID, lib.DB.Get(&user, q, uID)
}

// DeleteUser deletes the user "uID" from the DB
func DeleteUser(uID uint64) error {
	q := `DELETE FROM app_user WHERE id = $1`
	_, err := lib.DB.Exec(q, uID)
	return err
}

// SafeDeleteUser sets deleted at to current date, meaning this user is disactivated
func SafeDeleteUser(uID uint64) error {
	q := `UPDATE app_user SET deleted_at = CURRENT_DATE WHERE id = $1`
	if _, err := lib.DB.Exec(q, uID); err != nil {
		return err
	}

	q = `DELETE FROM package WHERE user_id = $1`
	if _, err := lib.DB.Exec(q, uID); err != nil {
		return err
	}

	q = `UPDATE creation SET state = 'delete' WHERE creator_id = $1`
	_, err := lib.DB.Exec(q, uID)
	return err
}

// UpdateUserPatch updates user informations
func UpdateUserPatch(uID uint64, patch lib.SQLPatch) error {
	builtQ := patch.GetUpdateQuery("app_user")
	if builtQ == "" {
		return nil
	}
	q := builtQ +
		`WHERE id = $` + fmt.Sprintf("%d", patch.Index+1) +
		` AND deleted_at IS NULL`

	patch.Args = append(patch.Args, uID)
	_, err := lib.DB.Exec(q, patch.Args...)
	return err
}

// UpdateUserAccountID updates user "uID" account ID
func UpdateUserAccountID(uID uint64, accID string) error {
	q := `UPDATE app_user SET account_id = $2 WHERE id = $1 AND account_id IS NULL`
	_, err := lib.DB.Exec(q, uID, accID)
	return err
}

// UpdateCustomerID updates user's customer ID
func UpdateCustomerID(uID uint64, customerID string) error {
	q := `UPDATE app_user SET customer_id=$2 WHERE id=$1`
	_, err := lib.DB.Exec(q, uID, customerID)
	return err
}

// UserNbPackages returns numbers of package the user uID has
func UserNbPackages(uID uint64) int64 {
	var nbPackages struct {
		Value int64 `db:"nb_pkg"`
	}

	q := `
		SELECT
			COUNT(p.user_id) nb_pkg
		FROM package p
		WHERE p.user_id = $1
	`

	lib.DB.Get(&nbPackages, q, uID)

	return nbPackages.Value
}

// UserSubFund substracts user "uID" fund by "amount"
func UserSubFund(uID uint64, amount uint64) error {
	q := `UPDATE appe_user SET fund = fund - $2 WHERE id = $1`
	_, err := lib.DB.Exec(q, uID, amount)
	return err
}

// SendActivationEmail prepares and sends an activation email
func SendActivationEmail(name string, to string) error {
	addToBytes := []byte(to)

	block, err := aes.NewCipher([]byte(lib.GetEncKey()))
	if err != nil {
		return err
	}

	ciphertext := make([]byte, aes.BlockSize+len(addToBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], addToBytes)

	// convert to base64
	validationToken := base64.URLEncoding.EncodeToString(ciphertext)

	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@wooble.io")
	m.SetHeader("To", to)
	m.SetHeader("Subject", fmt.Sprintf("Hi %s, Validate your Wooble account!", name))

	var out bytes.Buffer
	tmpl := template.Must(template.New("email.html").ParseFiles("./templates/email.html"))
	data := Email{
		Title: "Activate your Wooble account!",
		Name:  name,
		Elements: []Element{
			{"p", "I am happy to see you here! To use Wooble you have to activate your email address. Please click the link below so Wooble can activate your account :", ""},
			{"a", "Activate my account", fmt.Sprintf(`%s/settings/account?t=%s`, lib.GetOrigins()[0], validationToken)},
			{"p", "", ""},
			{"p", "Thank you!", ""},
		},
	}

	if err := tmpl.Execute(&out, data); err != nil {
		return err
	}

	m.SetBody("text/html", string(out.Bytes()))

	emailHost := strings.Split(lib.GetEmailHost(), ":")
	port, _ := strconv.Atoi(emailHost[1])

	d := gomail.NewDialer(emailHost[0], port, "no-reply@wooble.io", lib.GetEmailPasswd())

	return d.DialAndSend(m)
}

// ActivateUser activate user's email if the token is valid
func ActivateUser(email string, token string) error {
	ciphertext, _ := base64.URLEncoding.DecodeString(token)

	block, err := aes.NewCipher([]byte(lib.GetEncKey()))
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	if string(ciphertext) != email {
		return errors.New("Invalid token")
	}

	q := `UPDATE app_user SET is_active = TRUE WHERE email = $1`
	if _, err = lib.DB.Exec(q, email); err != nil {
		return err
	}

	q = `INSERT INTO package(title, user_id, referer, nb_build) SELECT 'Super Package' AS title, id, '' AS referer, 0 AS nb_build FROM app_user WHERE email = $1`
	_, err = lib.DB.Exec(q, email)

	return err
}

// IsPasswordValid checks if a password is valid
func (u *User) IsPasswordValid(passwd string) bool {
	if u.Secret == "" && u.Salt == "" {
		var userCreds struct {
			Salt   string `db:"salt_key"`
			Secret string `db:"passwd"`
		}
		q := `SELECT salt_key, passwd FROM app_user WHERE id=$1`
		lib.DB.Get(&userCreds, q, u.ID)
		u.Salt = userCreds.Salt
		u.Secret = userCreds.Secret
	}
	cp, err := lib.Encrypt(passwd, []byte(u.Salt))
	if err != nil || cp == "" {
		return false
	}
	return u.Secret == cp
}
