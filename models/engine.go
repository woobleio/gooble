package model

import (
	"wooble/lib"
)

// Engine is the engine used by a creation (JSES5 is the only one available for now)
type Engine struct {
	Name string `json:"name"   db:"eng.name"`

	ContentType string `json:"-" db:"content_type"`
	Extension   string `json:"-" db:"extension"`
}

// EngineByName returns engine with name "name"
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
