package model

import (
	"wooble/lib"
)

// Package is a bundle of creations
type Package struct {
	ID uint64 `json:"id" db:"pkg.id"`

	Title string `json:"title" binding:"required" db:"pkg.title"`

	Domains   lib.StringSlice `json:"domains" db:"domains"`
	Key       string          `json:"-" db:"key"`
	UserID    uint64          `json:"-" db:"user_id"`
	User      User            `json:"-" db:""`
	Creations []Creation      `json:"creations,omitempty" db:""`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"pkg.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"pkg.updated_at"`
}

// PackageForm is a form standard for package
type PackageForm struct {
	UserID uint64
	Title  string `json:"title" binding:"required"`

	Domains lib.StringSlice `json:"domains"`

	Key string
}

// PopulateCreations populates creations in the package
func (p *Package) PopulateCreations() error {
	q := `
	SELECT
		c.id "crea.id",
		c.title,
		c.version,
		c.has_document,
		c.has_script,
		c.has_style,
		u.id "user.id",
		u.name
	FROM package_creation pc
	INNER JOIN creation c ON (pc.creation_id = c.id)
	INNER JOIN app_user u ON (c.creator_id = u.id)
	WHERE pc.package_id = $1
	`

	return lib.DB.Select(&p.Creations, q, p.ID)
}

// AllPackages returns all packages
func AllPackages(opt lib.Option, userID uint64) (*[]Package, error) {
	var packages []Package
	q := lib.Query{
		Q: `SELECT
  		pkg.id "pkg.id",
  		pkg.user_id,
  		pkg.title "pkg.title",
  		pkg.domains,
  		pkg.key,
  		pkg.created_at "pkg.created_at",
  		pkg.updated_at "pkg.updated_at"
  	FROM package pkg
		WHERE pkg.user_id = $1
    `,
		Opt: &opt,
	}

	query := q.String()

	return &packages, lib.DB.Select(&packages, query, userID)
}

// PackageByID returns package with id "id"
func PackageByID(id string, userID uint64) (*Package, error) {
	var pkg Package
	q := `
	SELECT
		pkg.id "pkg.id",
		pkg.user_id,
		pkg.title "pkg.title",
		pkg.domains,
		pkg.key,
		pkg.created_at "pkg.created_at",
		pkg.updated_at "pkg.updated_at"
	FROM package pkg
	WHERE pkg.id = $1
  AND pkg.user_id = $2
	`

	if err := lib.DB.Get(&pkg, q, id, userID); err != nil {
		return nil, err
	}

	return &pkg, pkg.PopulateCreations()
}

// NewPackage created a new package
func NewPackage(data *PackageForm) (pkgID uint64, err error) {
	q := `INSERT INTO package(title, user_id, domains, key) VALUES ($1, $2, $3, $4) RETURNING id`
	err = lib.DB.QueryRow(q, data.Title, data.UserID, data.Domains, data.Key).Scan(&pkgID)
	return pkgID, err
}

// PushCreation pushes a creation in the package
func PushCreation(pkgID uint64, creaID uint64) error {
	q := `INSERT INTO package_creation(package_id, creation_id) VALUES ($1, $2)`
	_, err := lib.DB.Exec(q, pkgID, creaID)
	return err
}
