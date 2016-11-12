package model

import (
	"wooble/lib"
)

type User struct {
	ID uint64 `json:"id" db:"user.id"`

	Email     string `json:"email,omitempty"     db:"email"`
	Name      string `json:"name"                db:"name"`
	IsCreator bool   `json:"isCreator,omitempty" db:"is_creator"`

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
			u.updated_at "user.updated_at"
		FROM app_user u
		WHERE u.id = $1
	`

	if err := lib.DB.Get(&user, q, id); err != nil {
		return nil, err
	}

	return &user, nil
}
