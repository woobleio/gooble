package form

// CreationForm is a form for creation
type CreationForm struct {
	Engine string `json:"engine" validate:"required,alphanum"`
	Title  string `json:"title" validate:"required,alpha,max=15,min=3"`

	State       string `json:"state" validate:"alpha"`
	Description string `json:"description" validate:"ascii"`

	Price uint64 `json:"price,omitempty"`
}

// CreationCodeForm is a form for creation code
type CreationCodeForm struct {
	Script string `json:"script" validate:"required"`

	Style    string `json:"style"`
	Document string `json:"document"`
}
