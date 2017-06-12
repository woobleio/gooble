package form

// UserForm is the form for users
type UserForm struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,ascii,min=4,max=18,alphanum"`
	Fullname string `json:"fullname" validate:"required,ascii,min=4,max=50"`
	Secret   string `json:"secret" validate:"required,min=8"`
	Plan     struct {
		Label string `json:"label" validate:"required,ascii"`
	} `json:"plan" validate:"required"`

	CardToken string `json:"cardToken"`

	IsCreator bool `json:"isCreator"`
}

// UserPatchForm is the form for patching users
type UserPatchForm struct {
	Email     *string `json:"email" validate:"omitempty,email" db:"email"`
	Name      *string `json:"name" validate:"omitempty,ascii,min=4,max=18,alphanum" db:"name"`
	Fullname  *string `json:"fullname" validate:"omitempty,ascii,min=4,max=50" db:"fullname"`
	NewSecret *string `json:"newSecret" validate:"omitempty,min=8" db:"passwd"`
	IsCreator *bool   `json:"isCreator" db:"is_creator"`

	Plan *struct {
		Label string `json:"label" validate:"required,ascii"`
	} `json:"plan" validate:"omitempty"`
	CardToken *string `json:"cardToken"`

	PicPath      *string `json:"profilePath" validate:"omitempty" db:"pic_path"`
	Website      *string `json:"website" validate:"omitempty" db:"website"`
	CodepenName  *string `json:"codepenName" validate:"omitempty,ascii,max=30" db:"codepen_name"`
	DribbbleName *string `json:"dribbbleName" validate:"omitempty,ascii,max=30" db:"dribbble_name"`
	GithubName   *string `json:"githubName" validate:"omitempty,ascii,max=30" db:"github_name"`
	TwitterName  *string `json:"twitterName" validate:"omitempty,ascii,max=30" db:"twitter_name"`

	OldSecret *string `json:"secret" validate:"omitempty,min=8"`
	BankToken *string `json:"bankToken"`

	Salt *string `json:"-" db:"salt_key"`

	// Operations
	Withdraw *bool `json:"withdraw"`
}
