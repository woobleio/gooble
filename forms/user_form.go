package form

// UserForm is the form for users
type UserForm struct {
	Email  string `json:"email" validate:"required,email"`
	Name   string `json:"name" validate:"required,min=4,max=12,alpha"`
	Secret string `json:"secret" validate:"required,min=8"`
	Plan   string `json:"plan" validate:"required"`

	CardToken string `json:"cardToken"`

	IsCreator bool `json:"isCreator"`
}
