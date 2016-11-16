package model

import (
	"wooble/lib"
)

type Package struct {
	ID uint64 `json:"id" db:"pkg.id"`

	UserID uint64 `json:"userId" db:"user_id"`
	Title  string `json:"title" db:"pkg.title"`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"pkg.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"pkg.updated_at"`
}

type PackageForm struct {
	Title string `json:"title" binding:"required"`

	UserID uint64
}

func NewPackage(data *PackageForm) (pkgId uint64, err error) {
	q := `INSERT INTO package(title, user_id) VALUES ($1, $2) RETURNING id`
	err = lib.DB.QueryRow(q, data.Title, data.UserID).Scan(&pkgId)
	return pkgId, err
}

func PushCreation(pkgID uint64, creaID uint64) error {
	q := `INSERT INTO package_creation(package_id, creation_id) VALUES ($1, $2)`
	_, err := lib.DB.Exec(q, pkgID, creaID)
	return err
}
