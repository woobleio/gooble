package model

import (
	"wooble/lib"
	enum "wooble/models/enums"
)

// Package is a bundle of creations
type Package struct {
	ID lib.ID `json:"id" db:"pkg.id"`

	Title string `json:"title" validate:"required" db:"pkg.title"`

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
func AllPackages(opt *lib.Option, userID uint64) (*[]Package, error) {
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
		Opt: opt,
	}

	query := q.String()

	return &packages, lib.DB.Select(&packages, query, userID)
}

// PackageByID returns package with id "id"
func PackageByID(uID uint64, id lib.ID) (*Package, error) {
	pkg := new(Package)
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
	WHERE pkg.id = $2
  AND pkg.user_id = $1
	`

	if err := lib.DB.Get(pkg, q, uID, id); err != nil {
		return nil, err
	}

	return pkg, pkg.PopulateCreations()
}

// PackageNbCrea returns the number of creations in the package id "id"
func PackageNbCrea(id lib.ID) uint64 {
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

	lib.DB.Get(&nbCrea, q, id)

	return nbCrea.Value
}

// NewPackage creates a new package
func NewPackage(pkg *Package) (*Package, error) {
	q := `INSERT INTO package(title, user_id, domains) VALUES ($1, $2, $3) RETURNING id`
	return pkg, lib.DB.QueryRow(q, pkg.Title, pkg.UserID, pkg.Domains).Scan(&pkg.ID)
}

// DeletePackage delete a package
func DeletePackage(uID uint64, id lib.ID) error {
	q := `DELETE FROM package WHERE id = $2 AND user_id = $1`
	_, err := lib.DB.Exec(q, uID, id)
	return err
}

// UpdatePackageCreation udpates package creation information
func UpdatePackageCreation(pkg *Package) error {
	q := `
	UPDATE package_creation
	SET alias = $3,
	version = (
		SELECT $4 FROM creation
		WHERE ((
			$4 = ANY (versions[1:array_length(versions, 1)-1])
			AND state = $5
		)
		OR (
			$4 = ANY (versions)
			AND state != $5
		))
		AND id = $6
	)
	WHERE package_id = (
		SELECT id FROM package WHERE user_id = $1 AND id = $2
	)
	`
	crea := pkg.Creations[0]
	_, err := lib.DB.Exec(q, pkg.UserID, pkg.ID, crea.Alias, crea.Version, enum.Draft, crea.ID)
	return err
}

// NewPackageCreation create a new relationship with package and creation
func NewPackageCreation(pkgID lib.ID, creaID lib.ID, version string) error {
	q := `INSERT INTO package_creation(package_id, creation_id, version) VALUES ($1, $2, $3)`
	_, err := lib.DB.Exec(q, pkgID, creaID, version)
	return err
}

// UpdatePackage updates package form
func UpdatePackage(pkg *Package) error {
	q := `UPDATE package SET title=$3, domains=$4 WHERE id=$1 AND user_id=$2`
	_, err := lib.DB.Exec(q, pkg.ID, pkg.UserID, pkg.Title, pkg.Domains)
	return err
}

// UpdatePackageSource updates package source
func UpdatePackageSource(uID uint64, pkgID lib.ID, source string) error {
	q := `UPDATE package SET source = $3 WHERE id = $2 AND user_id = $1`
	_, err := lib.DB.Exec(q, uID, pkgID, source)
	return err
}

// BulkUpdatePackageSource updates somes packages "ids" source
func BulkUpdatePackageSource(ids lib.StringSlice, source string) error {
	q := `UPDATE package SET source = $2 WHERE id = ANY($1)`
	_, err := lib.DB.Exec(q, ids, source)
	return err
}

// DeletePackageCreation delete a creation from a package
func DeletePackageCreation(uID uint64, pkgID lib.ID, creaID lib.ID) error {
	q := `
	DELETE FROM package_creation 
	USING package 
	WHERE package_id = $2 AND package.id = package_id
	AND creation_id = $3
	AND package.user_id = $1
	`

	_, err := lib.DB.Exec(q, uID, pkgID, creaID)
	return err
}
