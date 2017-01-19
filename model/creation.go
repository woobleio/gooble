package model

import "wooble/lib"

// Creation is a Wooble creation
type Creation struct {
	ID uint64 `json:"id"      db:"crea.id"`

	Title   string `json:"title"  db:"title"`
	Creator User   `json:"creator" db:""`
	Version string `json:"version" db:"version"`

	CreatorID uint64 `json:"-"       db:"creator_id"`
	SourceID  uint64 `json:"-"       db:"source_id"`
	HasDoc    bool   `json:"-"       db:"has_document"`
	HasScript bool   `json:"-"       db:"has_script"`
	HasStyle  bool   `json:"-"       db:"has_style"`
	Engine    Engine `json:"-" db:""`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"crea.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"crea.updated_at"`
}

// CreationForm is a form for creation
type CreationForm struct {
	CreatorID uint64

	Engine string `json:"engine" binding:"required"`
	Title  string `json:"title" binding:"required"`

	Document string `json:"document"`
	Script   string `json:"script"`
	Style    string `json:"style"`

	Version string
}

// BaseVersion is creation default version
const BaseVersion string = "1.0"

// AllCreations returns all creations
func AllCreations(opt lib.Option) (*[]Creation, error) {
	var creations []Creation
	q := lib.Query{
		Q: `SELECT
	    c.id "crea.id",
	    c.title,
	    c.created_at "crea.created_at",
	    c.updated_at "crea.updated_at",
	    c.version,
			c.has_document,
			c.has_script,
			c.has_style,
			e.name "eng.name",
			e.extension,
			e.content_type,
	    u.id "user.id",
	    u.name
	  FROM creation c
	  INNER JOIN app_user u ON (c.creator_id = u.id)
		INNER JOIN engine e ON (c.engine=e.name)
		`,
		Opt: &opt,
	}

	query := q.String()

	return &creations, lib.DB.Select(&creations, query)
}

// CreationByTitle returns a creation with the title "title"
func CreationByTitle(title string) (*Creation, error) {
	var crea Creation
	q := `
  SELECT
    c.id "crea.id",
    c.title,
    c.created_at "crea.created_at",
    c.updated_at "crea.updated_at",
    c.version,
		c.has_document,
		c.has_script,
		c.has_style,
		e.name "eng.name",
		e.extension,
		e.content_type,
    u.id "user.id",
    u.name
  FROM creation c
  INNER JOIN app_user u ON (c.creator_id = u.id)
	INNER JOIN engine e ON (c.engine=e.name)
  WHERE c.title = $1
	`

	return &crea, lib.DB.Get(&crea, q, title)
}

// CreationByID returns a creation with the id "id"
func CreationByID(id string) (*Creation, error) {
	var crea Creation
	q := `
  SELECT
    c.id "crea.id",
    c.title,
    c.created_at "crea.created_at",
    c.updated_at "crea.updated_at",
    c.version,
		c.has_document,
		c.has_script,
		c.has_style,
		e.name "eng.name",
		e.extension,
		e.content_type,
    u.id "user.id",
    u.name
  FROM creation c
  INNER JOIN app_user u ON (c.creator_id = u.id)
	INNER JOIN engine e ON (c.engine=e.name)
  WHERE c.id = $1
	`

	return &crea, lib.DB.Get(&crea, q, id)
}

// DeleteCreation deletes creation id "id"
func DeleteCreation(id uint64) error {
	_, err := lib.DB.Exec(`DELETE FROM creation WHERE id = $1`, id)
	return err
}

// NewCreation creates a creation
func NewCreation(data *CreationForm) (creaID uint64, err error) {
	q := `INSERT INTO creation(title, creator_id, version, has_document, has_script, has_style, engine) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = lib.DB.QueryRow(q, data.Title, data.CreatorID, data.Version, data.Document != "", data.Script != "", data.Style != "", data.Engine).Scan(&creaID)
	return creaID, err
}
