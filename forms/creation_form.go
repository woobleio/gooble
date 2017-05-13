package form

// CreationForm is a form for creation
type CreationForm struct {
	Title string `json:"title" validate:"required,min=3,max=30"`

	Alias       string `json:"alias" validate:"omitempty,alpha,max=10"`
	Engine      string `json:"engine" validate:"omitempty,alphanum"`
	State       string `json:"state" validate:"omitempty,alpha"`
	Description string `json:"description" validate:"omitempty"`
	ThumbPath   string `json:"thumbPath" validate:"omitempty,max=300"`
}

// CreationCodeForm is a form for creation code
type CreationCodeForm struct {
	Script string `json:"script" validate:"required"`

	Title        string `json:"title"`
	Style        string `json:"style"`
	Document     string `json:"document"`
	ParsedScript string `json:"parsedScript"`
}
