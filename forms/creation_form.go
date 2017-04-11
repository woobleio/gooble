package form

// CreationForm is a form for creation
type CreationForm struct {
	Title string `json:"title" validate:"required,min=3,max=20"`

	Alias       string `json:"alias" validate:"omitempty,alpha,max=10"`
	Engine      string `json:"engine" validate:"omitempty,alphanum"`
	State       string `json:"state" validate:"omitempty,alpha"`
	Description string `json:"description" validate:"ascii"`
	ThumbPath   string `json:"thumbPath" validate:"omitempty,max=300"`

	Price uint64 `json:"price,omitempty"`
}

// CreationCodeForm is a form for creation code
type CreationCodeForm struct {
	Script string `json:"script" validate:"required"`

	Title    string `json:"title"`
	Style    string `json:"style"`
	Document string `json:"document"`
}
