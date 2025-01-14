package form

// PackageForm is a form standard for package
type PackageForm struct {
	Title string `json:"title" validate:"required,ascii,max=15"`

	Referer string `json:"referer"`
}

// PackageCreationForm if a form standard for pushing creation in a package
type PackageCreationForm struct {
	CreationID string `json:"creationId"`
	Version    uint64 `json:"version"`
	Alias      string `json:"alias" validate:"omitempty,ascii,min=2,max=10,alpha"`
}

// PackagePatchForm is the form for patching package
type PackagePatchForm struct {
	Title *string `json:"title" validate:"omitempty,ascii,max=25" db:"title"`

	Source  *string `json:"-" db:"source"`
	NbBuild *uint64 `json:"-" db:"nb_build"`
	BuiltAt *string `json:"-" db:"built_at"`

	// Operation build
	Operation *string `json:"operation" validate:"omitempty"`
}
