package model

import (
	"wooble/lib"
)

type Engine struct {
	Name string `json:"name"   db:"name"`

	ContentType string `json:"-" db:"content_type"`
	Extension   string `json:"-" db:"extension"`
}

func EngineByName(name string) (*Engine, error) {
	var eng Engine
	q := `
		SELECT
			*
		FROM engine
		WHERE name = $1
	`

	if err := lib.DB.Get(&eng, q, name); err != nil {
		return nil, err
	}

	return &eng, nil
}
