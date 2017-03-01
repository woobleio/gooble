package model

import (
	"encoding/hex"

	"wooble/lib"

	"golang.org/x/crypto/scrypt"
)

// User is a Wooble user
type User struct {
	ID         uint64         `json:"-" db:"user.id"`
	CustomerID string         `json:"-" db:"customer_id"`
	AccountID  lib.NullString `json:"-" db:"account_id"`

	Email string `json:"email,omitempty" db:"email"`
	Name  string `json:"name" db:"name"`

	Plan     *Plan      `json:"plan,omitempty" db:""`
	Packages *[]Package `json:"packages,omitempty" db:""`

	IsCreator bool   `json:"isCreator,omitempty" db:"is_creator"`
	Fund      uint64 `json:"totalDue,omitempty" db:"fund"`

	Secret string `json:"-" db:"passwd"`
	Salt   string `json:"-" db:"salt_key"`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"user.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"user.updated_at"`
	DeletedAt *lib.NullTime `json:"deletedAt,omitempty" db:"user.deleted_at"`
}

// UserPublicByName returns user public profile with the name "username"
func UserPublicByName(username string) (*User, error) {
	var user User
	q := `
		SELECT
			u.name,
			u.is_creator,
			u.created_at "user.created_at"
    FROM app_user u
		WHERE u.name = $1
		AND u.deleted_at IS NULL
	`
	return &user, lib.DB.Get(&user, q, username)
}

// UserPrivateByID returns user with id "id"
func UserPrivateByID(id uint64) (*User, error) {
	var user User
	q := `
		SELECT DISTINCT ON (u.id)
			u.id "user.id",
			u.email,
			u.name,
			u.is_creator,
			u.created_at "user.created_at",
			u.updated_at "user.updated_at",
			u.salt_key,
      u.customer_id,
			u.account_id,
      pu.start_date,
      pu.end_date,
      pl.label "plan.label",
      pl.nb_pkg,
      pl.nb_crea,
      pl.nb_domains
    FROM app_user u
    LEFT OUTER JOIN plan_user pu ON (pu.user_id = u.id)
    LEFT OUTER JOIN plan pl ON (pl.label = pu.plan_label)
		WHERE u.id = $1
		AND u.deleted_at IS NULL
    ORDER BY u.id, pu.start_date DESC
	`

	if err := lib.DB.Get(&user, q, id); err != nil {
		return nil, err
	}

	if user.Plan == nil {
		user.Plan, _ = DefaultPlan()
	}

	return &user, nil
}

// NewUser creates a new user
func NewUser(user *User) (uID uint64, err error) {
	salt := lib.GenKey()
	cp, errPasswd := getPassword(user.Secret, []byte(salt))
	if errPasswd != nil {
		return 0, errPasswd
	}
	q := `INSERT INTO app_user(name, email, is_creator, passwd, salt_key, customer_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = lib.DB.QueryRow(q, user.Name, user.Email, user.IsCreator, cp, salt, user.CustomerID).Scan(&uID)

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
      u.passwd,
      u.is_creator,
      u.created_at "user.created_at",
      u.updated_at "user.updated_at",
      u.salt_key,
      u.customer_id,
      pu.start_date,
      pu.end_date,
      pl.label "plan.label",
      pl.nb_pkg,
      pl.nb_crea,
      pl.nb_domains
    FROM app_user u
    LEFT OUTER JOIN plan_user pu ON (pu.user_id = u.id)
    LEFT OUTER JOIN plan pl ON (pl.label = pu.plan_label)
    WHERE u.email = $1
		AND u.deleted_at IS NULL
    ORDER BY u.id, pu.start_date DESC
	`

	if err := lib.DB.Get(&user, q, email); err != nil {
		return nil, err
	}

	if user.Plan == nil {
		user.Plan, _ = DefaultPlan()
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
	_, err := lib.DB.Exec(q, uID)
	return err
}

// UpdateUser updates user form (password not included)
func UpdateUser(uID uint64, user *User) error {
	q := `UPDATE app_user SET name=$2, email=$3, is_creator=$4 WHERE id=$1 AND deleted_at IS NULL`
	_, err := lib.DB.Exec(q, uID, user.Name, user.Email, user.IsCreator)
	return err
}

// UpdateUserPassword updates user "uID" passwd and its salt key
func UpdateUserPassword(uID uint64, newSecret string) error {
	salt := lib.GenKey()
	cp, errPasswd := getPassword(newSecret, []byte(salt))
	if errPasswd != nil {
		return errPasswd
	}
	q := `UPDATE app_user SET passwd=$2, salt_key=$3 WHERE id=$1 AND deleted_at IS NULL`
	_, err := lib.DB.Exec(q, uID, cp, salt)
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
	cp, err := getPassword(passwd, []byte(u.Salt))
	if err != nil || cp == "" {
		return false
	}
	return u.Secret == cp
}

func getPassword(passwd string, salt []byte) (string, error) {
	cp, err := scrypt.Key([]byte(passwd), []byte(salt), 16384, 8, 1, 32)
	return hex.EncodeToString(cp), err
}
