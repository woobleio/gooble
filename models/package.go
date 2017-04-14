package model

import (
	"fmt"
	"wooble/lib"
	enum "wooble/models/enums"
)

// Package is a bundle of creations
type Package struct {
	ID lib.ID `json:"id" db:"pkg.id"`

	Title string `json:"title" validate:"required" db:"pkg.title"`

	Referer   *lib.NullString `json:"referer,omitempty" db:"referer"`
	UserID    uint64          `json:"-" db:"pkg.user_id"`
	User      User            `json:"-" db:""`
	Creations []Creation      `json:"creations,omitempty" db:""`
	Source    *lib.NullString `json:"source,omitempty" db:"source"`
	NbCrea    *lib.NullInt64  `json:"nbCreations" db:"nb_creations"`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"pkg.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"pkg.updated_at"`
}

// PopulateCreations populates creations in the package
func (p *Package) PopulateCreations() error {
	q := `
	SELECT
		c.id "crea.id",
		pc.version,
		c.title,
    c.creator_id,
		CASE WHEN c.state = 'draft' THEN c.versions[0:array_length(c.versions, 1) - 1] ELSE c.versions END AS versions,
		c.price,
		CASE WHEN pc.alias != '' THEN pc.alias ELSE c.alias END AS alias,
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
  		pkg.referer,
			pkg.source,
  		pkg.created_at "pkg.created_at",
  		pkg.updated_at "pkg.updated_at",
			(SELECT COUNT(pc.package_id) FROM package_creation pc WHERE pc.package_id = pkg.id GROUP BY pc.package_id) AS nb_creations
  	FROM package pkg
		WHERE pkg.user_id = $1
    `,
		Opt: opt,
	}

	if err := lib.DB.Select(&packages, q.String(), userID); err != nil {
		return nil, err
	}

	if q.Opt.HasPopulate("creations") {
		for i := range packages {
			packages[i].PopulateCreations()
		}
	}

	return &packages, nil
}

// PackageByID returns package with id "id"
func PackageByID(uID uint64, id lib.ID) (*Package, error) {
	pkg := new(Package)
	q := `
	SELECT
		pkg.id "pkg.id",
		pkg.user_id "pkg.user_id",
		pkg.title "pkg.title",
		pkg.referer,
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
	q := `INSERT INTO package(title, user_id, referer) VALUES ($1, $2, $3) RETURNING id`
	return pkg, lib.DB.QueryRow(q, pkg.Title, pkg.UserID, pkg.Referer).Scan(&pkg.ID)
}

// DeletePackage delete a package
func DeletePackage(uID uint64, id lib.ID) error {
	q := `DELETE FROM package WHERE id = $2 AND user_id = $1`
	_, err := lib.DB.Exec(q, uID, id)
	return err
}

// UpdatePackageCreation udpates package creation information
func UpdatePackageCreation(pkg *Package) error {
	crea := pkg.Creations[0]
	q := `
	UPDATE package_creation
	SET alias = $3,
	version =	CASE WHEN (
		SELECT ` + fmt.Sprintf("%d", crea.Version) + ` FROM creation
		WHERE ((
			$4 = ANY (versions[1:array_length(versions, 1)-1])
			AND state = $5
		)
		OR (
			$4 = ANY (versions)
			AND state != $5
		))
		AND id = $6)
	IS NOT NULL THEN $4 ELSE version END
	WHERE package_id = (
		SELECT id FROM package WHERE user_id = $1 AND id = $2
	) AND creation_id = $6
	`
	_, err := lib.DB.Exec(q, pkg.UserID, pkg.ID, crea.Alias, crea.Version, enum.Draft, crea.ID)
	return err
}

// NewPackageCreation create a new relationship with package and creation
func NewPackageCreation(pkgID lib.ID, creaID lib.ID, version uint64, alias string) error {
	q := `INSERT INTO package_creation(package_id, creation_id, version, alias) VALUES ($1, $2, $3, $4)`
	_, err := lib.DB.Exec(q, pkgID, creaID, version, alias)
	return err
}

// UpdatePackage updates package form
func UpdatePackage(pkg *Package) error {
	q := `UPDATE package SET title = $3, referer = $4 WHERE id = $1 AND user_id = $2`
	_, err := lib.DB.Exec(q, pkg.ID, pkg.UserID, pkg.Title, pkg.Referer)
	return err
}

// UpdatePackagePatch updates user informations
func UpdatePackagePatch(uID uint64, pkgID lib.ID, patch lib.SQLPatch) error {
	q := patch.GetUpdateQuery("package") +
		` WHERE user_id = $` + fmt.Sprintf("%d", patch.Index+1) +
		` AND id = $` + fmt.Sprintf("%d", patch.Index+2)

	patch.Args = append(patch.Args, uID)
	patch.Args = append(patch.Args, pkgID)
	_, err := lib.DB.Exec(q, patch.Args...)

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
