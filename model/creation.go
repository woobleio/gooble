package model

import (
	"wooble/lib"
)

type Creation struct {
	ID uint64 `json:"id"      db:"crea.id"`

	CreatorID uint64 `json:"-"       db:"creator_id"`
	Creator   User   `json:"creator" db:""`
	SourceID  uint64 `json:"-"       db:"source_id"`
	Title     string `json:"title"   db:"title"`
	Version   string `json:"version" db:"version"`
	HasDoc    bool   `json:"-"       db:"has_document"`
	HasScript bool   `json:"-"       db:"has_script"`
	HasStyle  bool   `json:"-"       db:"has_style"`
	Engine    string `json:"engine"  db:"engine"`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"crea.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"crea.updated_at"`
}

type CreationForm struct {
	Engine string `json:"engine" binding:"required"`
	Title  string `json:"title" binding:"required"`

	CreatorID uint64

	Document string `json:"document"`
	Script   string `json:"script"`
	Style    string `json:"style"`

	Version string `json:"version"`
}

const BASE_VERSION string = "1.0"

func AllCreations(opt lib.Option) (*[]Creation, error) {
	var creations []Creation
	q := lib.Query{`
    SELECT
      c.id "crea.id",
      c.title,
      c.created_at "crea.created_at",
      c.updated_at "crea.updated_at",
      c.version,
			c.has_document,
			c.has_script,
			c.has_style,
			c.engine,
      u.id "user.id",
      u.name
    FROM creation c
    INNER JOIN app_user u ON (c.creator_id = u.id)`,
		&opt,
	}

	query := q.String()

	if err := lib.DB.Select(&creations, query); err != nil {
		return nil, err
	}

	return &creations, nil
}

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
			c.engine,
      u.id "user.id",
      u.name
    FROM creation c
    INNER JOIN app_user u ON (c.creator_id = u.id)
    WHERE c.title = $1
	`

	if err := lib.DB.Get(&crea, q, title); err != nil {
		return nil, err
	}

	return &crea, nil
}

func NewCreation(data *CreationForm) (creaId int64, err error) {
	q := `INSERT INTO creation(title, creator_id, version, has_document, has_script, has_style, engine) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	crea, err := lib.DB.Exec(q, data.Title, data.CreatorID, data.Version, data.Document != "", data.Script != "", data.Style != "", data.Engine)
	if err != nil {
		return 0, err
	}
	creaId, _ = crea.LastInsertId()
	return creaId, nil
}
