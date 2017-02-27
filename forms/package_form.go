package form

import "wooble/lib"

// PackageForm is a form standard for package
type PackageForm struct {
	Title string `json:"title" validate:"required,ascii,max=15"`

	Domains lib.StringSlice `json:"domains"`
}

// PackageCreationForm if a form standard for pushing creation in a package
type PackageCreationForm struct {
	CreationID string `json:"creation" validate:"required"`
	Version    string `json:"version"`
	Alias      string `json:"alias"`
}
