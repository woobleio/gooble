package form

// PackageForm is a form standard for package
type PackageForm struct {
	Title string `json:"title" validate:"required,ascii,max=15"`

	Referer string `json:"referer"`
}

// PackageCreationForm if a form standard for pushing creation in a package
type PackageCreationForm struct {
	CreationID string `json:"creationId" validate:"required"`
	Version    uint64 `json:"version"`
	Alias      string `json:"alias"`
}

// PackagePatchForm is the form for patching users
type PackagePatchForm struct {
	Title  *string `json:"title" validate:"omitempty,ascii,max=15" db:"title"`
	Source *string `json:"-" db:"source"`

	// Operation build
	Build *bool `json:"build" validate:"omitempty"`
}
