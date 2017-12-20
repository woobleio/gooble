// Package engine defines everything related to script code for running a creation
package engine

// StdName is the name standard for a wooble object
const StdName string = "Woobly"

// Script is an interface for script langage
// ex : JS for JavaScript ES2015
type Script interface {
	// GetName returns obj name
	GetName() string

	// GetSource returns obj source code
	GetSource() string

	// GetParams returns obj parameters
	GetParams() []interface{}

	// IncludeHTMLCSS includes HTML and CSS code into the script object
	IncludeHTMLCSS(srcHTML string, srcCSS string) error

	// Control controles wether the object is valid or not
	Control() []error
}
