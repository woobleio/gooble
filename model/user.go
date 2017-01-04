package model

import (
	"encoding/hex"

	"wooble/lib"

	"golang.org/x/crypto/scrypt"
)

type User struct {
	ID uint64 `json:"-" db:"user.id"`

	Email  string `json:"email,omitempty" binding:"required" db:"email"`
	Name   string `json:"name" binding:"required" db:"name"`
	Passwd string `json:"secret,omitempty" binding:"required" db:"passwd"`

	IsCreator bool `json:"isCreator" db:"is_creator"`

	IsAuth bool   `json:"-" db:"is_auth"`
	Salt   string `json:"-" db:"salt_key"`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"user.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"user.updated_at"`
}

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

func NewUser(user *User) (uId uint64, err error) {
	salt := lib.GenKey()
	cp, err := getPassword(user.Passwd, []byte(salt))
	if err != nil {
		return 0, err
	}
	q := `INSERT INTO app_user(name, email, is_creator, passwd, salt_key) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = lib.DB.QueryRow(q, user.Name, user.Email, user.IsCreator, cp, salt).Scan(&uId)
	return uId, err
}

func UserByLogin(login string) (*User, error) {
	var user User
	q := `
		SELECT
			u.id "user.id",
			u.name,
			u.passwd,
			u.salt_key
		FROM app_user u
		WHERE u.name = $1
		OR u.email = $1
	`
	return &user, lib.DB.Get(&user, q, login)
}

func (u *User) IsPasswordValid(passwd string) bool {
	cp, err := getPassword(passwd, []byte(u.Salt))
	if err != nil || cp == "" {
		return false
	}
	return u.Passwd == cp
}

func getPassword(passwd string, salt []byte) (string, error) {
	cp, err := scrypt.Key([]byte(passwd), []byte(salt), 16384, 8, 1, 32)
	return hex.EncodeToString(cp), err
}
