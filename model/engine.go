package model

import (
	"wooble/lib"
)

type Engine struct {
	Name string `json:"name"   db:"eng.name"`

	ContentType string `json:"-" db:"content_type"`
	Extension   string `json:"-" db:"extension"`
}

func EngineByName(name string) (*Engine, error) {
	var eng Engine
	q := `
		SELECT
			e.name "eng.name",
			e.extension,
			e.content_type
		FROM engine e
		WHERE name = $1
	`

	if err := lib.DB.Get(&eng, q, name); err != nil {
		return nil, err
	}

	return &eng, nil
}
