package model

import "wooble/lib"

// Package is a bundle of creations
type Package struct {
	ID lib.ID `json:"id" db:"pkg.id"`

	Title string `json:"title" binding:"required" db:"pkg.title"`

	Domains   lib.StringSlice `json:"domains" db:"domains"`
	UserID    uint64          `json:"-" db:"pkg.user_id"`
	User      User            `json:"-" db:""`
	Creations []Creation      `json:"creations,omitempty" db:""`
	Source    *lib.NullString `json:"source,omitempty" db:"source"`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"pkg.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"pkg.updated_at"`
}

// PopulateCreations populates creations in the package
func (p *Package) PopulateCreations() error {
	q := `
	SELECT
		c.id "crea.id",
    pc.alias,
		pc.version,
		c.title,
    c.creator_id,
		c.versions,
		c.has_document,
		c.has_script,
		c.has_style,
		c.price,
		u.id "user.id",
		u.name,
    CASE WHEN c.price = 0  THEN 'false'
         WHEN u.id = $2 THEN 'false'
         WHEN cp.purchased_at IS NULL THEN 'true'
         ELSE 'false'
    END AS is_to_buy
	FROM package_creation pc
	INNER JOIN creation c ON (pc.creation_id = c.id)
	INNER JOIN app_user u ON (c.creator_id = u.id)
  LEFT OUTER JOIN creation_purchase cp ON (cp.user_id = $2 AND cp.creation_id = c.id)
	WHERE pc.package_id = $1
	`

	return lib.DB.Select(&p.Creations, q, p.ID.ValueDecoded, p.UserID)
}

// AllPackages returns all packages
func AllPackages(opt lib.Option, userID uint64) (*[]Package, error) {
	var packages []Package
	q := lib.Query{
		Q: `SELECT
  		pkg.id "pkg.id",
  		pkg.user_id "pkg.user_id",
  		pkg.title "pkg.title",
  		pkg.domains,
			pkg.source,
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
func PackageByID(pkg *Package) (*Package, error) {
	q := `
	SELECT
		pkg.id "pkg.id",
		pkg.user_id "pkg.user_id",
		pkg.title "pkg.title",
		pkg.domains,
		pkg.source,
		pkg.created_at "pkg.created_at",
		pkg.updated_at "pkg.updated_at"
	FROM package pkg
	WHERE pkg.id = $1
  AND pkg.user_id = $2
	`

	if err := lib.DB.Get(pkg, q, pkg.ID.ValueDecoded, pkg.UserID); err != nil {
		return nil, err
	}

	return pkg, pkg.PopulateCreations()
}

// PackageNbCrea returns the number of creations in the package id "id"
func PackageNbCrea(pkg *Package) uint64 {
	var nbCrea struct {
		Value uint64 `db:"nb_crea"`
	}

	q := `
	SELECT
		COUNT(p.id) nb_crea
	FROM package p
	JOIN package_creation pc ON (pc.package_id = p.id)
	WHERE p.id = $1
	`

	lib.DB.Get(&nbCrea, q, pkg.ID.ValueDecoded)

	return nbCrea.Value
}

// NewPackage creates a new package
func NewPackage(pkg *Package) (string, error) {
	var pkgID int64
	q := `INSERT INTO package(title, user_id, domains) VALUES ($1, $2, $3) RETURNING id`
	if err := lib.DB.QueryRow(q, pkg.Title, pkg.UserID, pkg.Domains).Scan(&pkgID); err != nil {
		return "", err
	}
	return lib.HashID(pkgID)
}

// PushCreation pushes a creation in the package
func PushCreation(pkg *Package, crea *Creation) error {
	q := `INSERT INTO package_creation(package_id, creation_id) VALUES ($1, $2)`
	_, err := lib.DB.Exec(q, pkg.ID.ValueDecoded, crea.ID.ValueDecoded)
	return err
}

// UpdatePackage updates package form
func UpdatePackage(pkg *Package) error {
	q := `UPDATE package SET title=$2, domains=$3 WHERE id=$1`
	_, err := lib.DB.Exec(q, pkg.ID, pkg.Title, pkg.Domains)
	return err
}

// UpdatePackageSource updates package source
func UpdatePackageSource(pkg *Package) error {
	q := `UPDATE package SET source=$2 WHERE id=$1`
	_, err := lib.DB.Exec(q, pkg.ID.ValueDecoded, pkg.Source.String)
	return err
}
