package form

import (
	model "wooble/models"
)

// CreationForm is a form for creation
type CreationForm struct {
	Title string `json:"title" validate:"required,min=3,max=30"`

	Alias       string                   `json:"alias" validate:"omitempty,alpha,max=10"`
	Engine      string                   `json:"engine" validate:"omitempty,alphanum"`
	State       string                   `json:"state" validate:"omitempty,alpha"`
	Description string                   `json:"description" validate:"omitempty,max=10000"`
	ThumbPath   string                   `json:"thumbPath" validate:"omitempty,max=300"`
	PreviewPos  model.PreviewPosition    `json:"previewPosition" validate:"omitempty"`
	Params      []model.CreationParam    `json:"params" validate:"omitempty"`
	Functions   []model.CreationFunction `json:"functions" validate:"omitempty"`
	Version     int                      `json:"version" validate:"omitempty"`
}

// CreationPatchForm is the form for patching creation
type CreationPatchForm struct {
	PreviewPositionID *string `json:"position" validate:"omitempty" db:"preview_position_id"`

	// Operation generateDefaultThumb
	Operation *string `json:"operation" validate:"omitempty"`
}

// CreationCodeForm is a form for creation code
type CreationCodeForm struct {
	Script string                `json:"script" validate:"required"`
	Params []model.CreationParam `json:"params" validate:"omitempty"`

	Title        string `json:"title"`
	Style        string `json:"style"`
	Document     string `json:"document"`
	ParsedScript string `json:"parsedScript"`
}
