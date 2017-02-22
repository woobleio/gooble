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

// PackageForm is a form standard for package
type PackageForm struct {
	UserID uint64
	Title  string `json:"title" binding:"required"`

	Domains lib.StringSlice `json:"domains"`
}

// PackageCreationForm if a form standard for pushing creation in a package
type PackageCreationForm struct {
	PackageID  uint64
	CreationID string `json:"creation" binding:"required"`
	Version    string `json:"version"`
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
func PackageByID(id string, userID uint64) (*Package, error) {
	var pkg Package
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

	decodeID, _ := lib.DecodeHash(id)

	if err := lib.DB.Get(&pkg, q, decodeID, userID); err != nil {
		return nil, err
	}

	return &pkg, pkg.PopulateCreations()
}

// PackageNbCrea returns the number of creations in the package id "id"
func PackageNbCrea(id string) int64 {
	var nbCrea struct {
		Value int64 `db:"nb_crea"`
	}

	q := `
	SELECT
		COUNT(p.id) nb_crea
	FROM package p
	JOIN package_creation pc ON (pc.package_id = p.id)
	WHERE p.id = $1
	`

	decodeID, _ := lib.DecodeHash(id)

	lib.DB.Get(&nbCrea, q, decodeID)

	return nbCrea.Value
}

// NewPackage creates a new package
func NewPackage(data *PackageForm) (string, error) {
	var pkgID int64
	q := `INSERT INTO package(title, user_id, domains) VALUES ($1, $2, $3) RETURNING id`
	if err := lib.DB.QueryRow(q, data.Title, data.UserID, data.Domains).Scan(&pkgID); err != nil {
		return "", err
	}
	return lib.HashID(pkgID)
}

// PushCreation pushes a creation in the package
func PushCreation(pkgID uint64, form *PackageCreationForm) error {
	decodeCreaID, _ := lib.DecodeHash(form.CreationID)
	q := `INSERT INTO package_creation(package_id, creation_id) VALUES ($1, $2)`
	_, err := lib.DB.Exec(q, pkgID, decodeCreaID)
	return err
}

// UpdatePackage updates package form
func UpdatePackage(pkg *Package) error {
	q := `UPDATE package SET title=$2, domains=$3 WHERE id=$1`
	_, err := lib.DB.Exec(q, pkg.ID, pkg.Title, pkg.Domains)
	return err
}

// UpdatePackageSource updates package source
func UpdatePackageSource(source string, pkgID lib.ID) error {
	q := `UPDATE package SET source=$2 WHERE id=$1`
	_, err := lib.DB.Exec(q, pkgID.ValueDecoded, source)
	return err
}
