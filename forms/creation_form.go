package form

// CreationForm is a form for creation
type CreationForm struct {
	Engine string `json:"engine" binding:"required"`
	Title  string `json:"title" binding:"required"`

	State       string `json:"state"`
	Description string `json:"description"`

	Price uint64 `json:"price,omitempty"`
}

// CreationCodeForm is a form for creation code
type CreationCodeForm struct {
	Script string `json:"script" binding:"required"`

	Style    string `json:"style"`
	Document string `json:"document"`
}
