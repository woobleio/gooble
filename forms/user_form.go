package form

// UserForm is the form for users
type UserForm struct {
	Email  string `json:"email" validate:"required,email"`
	Name   string `json:"name" validate:"required,min=4,max=18,alpha"`
	Secret string `json:"secret" validate:"required,min=8"`
	Plan   string `json:"plan" validate:"required"`

	CardToken string `json:"cardToken"`

	IsCreator bool `json:"isCreator"`
}

// UserPatchForm is the form for patching users
type UserPatchForm struct {
	Email     *string `json:"email" validate:"omitempty,email" db:"email"`
	Name      *string `json:"name" validate:"omitempty,min=2,max=18,alpha" db:"name"`
	NewSecret *string `json:"newSecret" validate:"omitempty,min=8" db:"passwd"`
	IsCreator *bool   `json:"isCreator" db:"is_creator"`

	PicPath       *string `json:"profilePath" validate:"omitempty" db:"pic_path"`
	Website       *string `json:"website" validate:"omitempty" db:"website"`
	CodepenioName *string `json:"codepenioName" validate:"omitempty,ascii,max=30" db:"codepenio_name"`
	DribbbleName  *string `json:"dribbbleName" validate:"omitempty,ascii,max=30" db:"dribbble_name"`
	GithubName    *string `json:"githubName" validate:"omitempty,ascii,max=30" db:"github_name"`
	TwitterName   *string `json:"twitterName" validate:"omitempty,ascii,max=30" db:"twitter_name"`

	OldSecret *string `json:"secret" validate:"omitempty,min=8"`
	BankToken *string `json:"bankToken"`
	CardToken *string `json:"cardToken"`

	Salt *string `json:"-" db:"salt_key"`

	// Operation withdraw
	Withdraw *bool `json:"withdraw"`
}
