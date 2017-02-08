package model

import (
	"encoding/hex"

	"wooble/lib"

	"golang.org/x/crypto/scrypt"
)

// User is a Wooble user
type User struct {
	ID         uint64 `json:"-" db:"user.id"`
	CustomerID string `json:"-" db:"customer_id"`

	Email string `json:"email,omitempty" db:"email"`
	Name  string `json:"name" db:"name"`

	Plan Plan `json:"plan" db:""`

	IsCreator bool `json:"isCreator" db:"is_creator"`

	Secret string `json:"-" db:"passwd"`
	Salt   string `json:"-" db:"salt_key"`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"user.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"user.updated_at"`
}

// UserForm is the form for users
type UserForm struct {
	Email  string `json:"email" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Secret string `json:"secret" binding:"required"`
	Plan   string `json:"plan" binding:"required"`

	CardToken string `json:"cardToken"`

	IsCreator bool `json:"isCreator"`
}

// UserByID returns user with id "id"
func UserByID(id uint64) (*User, error) {
	var user User
	q := `
		SELECT
			u.id "user.id",
			u.email,
			u.name,
			u.is_creator,
			u.created_at "user.created_at",
			u.updated_at "user.updated_at",
			u.salt_key
		FROM app_user u
		WHERE u.id = $1
	`

	return &user, lib.DB.Get(&user, q, id)
}

// NewUser creates a new user
func NewUser(user *UserForm) (uID uint64, err error) {
	customer, err := lib.NewCustomer(user.Email, user.Plan, user.CardToken)
	if err != nil {
		return 0, err
	}

	salt := lib.GenKey()
	cp, err := getPassword(user.Secret, []byte(salt))
	if err != nil {
		return 0, err
	}
	q := `INSERT INTO app_user(name, email, is_creator, passwd, salt_key, customer_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err = lib.DB.QueryRow(q, user.Name, user.Email, user.IsCreator, cp, salt, customer.ID).Scan(&uID)
	return uID, err
}

// UserByEmail returns user with a specified email or name
func UserByEmail(email string) (*User, error) {
	var user User
	q := `
		SELECT
			u.id "user.id",
			u.name,
			u.passwd,
			u.salt_key
		FROM app_user u
		WHERE u.email = $1
	`
	return &user, lib.DB.Get(&user, q, email)
}

// IsPasswordValid checks if a password is valid
func (u *User) IsPasswordValid(passwd string) bool {
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
