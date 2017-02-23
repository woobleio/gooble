package form

// UserForm is the form for users
type UserForm struct {
	Email  string `json:"email" binding:"required"`
	Name   string `json:"name" binding:"required"`
	Secret string `json:"secret" binding:"required"`
	Plan   string `json:"plan" binding:"required"`

	CardToken string `json:"cardToken"`

	IsCreator bool `json:"isCreator"`
}
