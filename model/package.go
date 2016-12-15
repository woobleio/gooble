package model

import (
	"wooble/lib"
)

type Package struct {
	ID uint64 `json:"id" db:"pkg.id"`

	Domains lib.StringSlice `json:"domains" binding:"required" db:"domains"`
	FEngine string          `json:"engine" binding:"required" db:"eng.name"`
	Title   string          `json:"title" binding:"required" db:"pkg.title"`

	Key       string     `json:"-" db:"key"`
	Engine    Engine     `json:"-" db:""`
	UserID    uint64     `json:"-" db:"user_id"`
	User      User       `json:"user" db:""`
	Creations []Creation `json:"creations" db:""`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"pkg.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"pkg.updated_at"`
}

func (p *Package) PopulateCreations() error {
	q := `
	SELECT
		c.id "crea.id",
		c.title,
		c.version,
		c.has_document,
		c.has_script,
		c.has_style,
		e.name "eng.name",
		e.extension,
		e.content_type,
		u.id "user.id",
		u.name
	FROM package_creation pc
	INNER JOIN creation c ON (pc.creation_id = c.id)
	INNER JOIN app_user u ON (c.creator_id = u.id)
	INNER JOIN engine e ON (c.engine = e.name)
	WHERE pc.package_id = $1
	`

	return lib.DB.Select(&p.Creations, q, p.ID)
}

func PackageByID(id uint64) (*Package, error) {
	var pkg Package
	q := `
	SELECT
		pkg.id "pkg.id",
		pkg.user_id,
		pkg.title "pkg.title",
		pkg.domains,
		pkg.key,
		pkg.created_at "pkg.created_at",
		pkg.updated_at "pkg.updated_at",
		u.id "user.id",
    u.name,
		e.name "eng.name",
		e.content_type,
		e.extension
	FROM package pkg
	INNER JOIN app_user u ON (pkg.user_id = u.id)
	INNER JOIN engine e ON (pkg.engine = e.name)
	WHERE pkg.id = $1
	`

	if err := lib.DB.Get(&pkg, q, id); err != nil {
		return nil, err
	}

	return &pkg, pkg.PopulateCreations()
}

func NewPackage(data *Package) (pkgId uint64, err error) {
	q := `INSERT INTO package(title, engine, user_id, domains, key) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err = lib.DB.QueryRow(q, data.Title, data.FEngine, data.UserID, data.Domains, data.Key).Scan(&pkgId)
	return pkgId, err
}

func PushCreation(pkgID uint64, creaID uint64) error {
	q := `INSERT INTO package_creation(package_id, creation_id) VALUES ($1, $2)`
	_, err := lib.DB.Exec(q, pkgID, creaID)
	return err
}