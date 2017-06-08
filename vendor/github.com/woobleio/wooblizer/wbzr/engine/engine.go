package engine

// StdName is the name standard for a wooble object
const StdName string = "Woobly"

// Script is an interface for script langage
type Script interface {
	GetName() string
	GetSource() string
	GetParams() []interface{}

	// IncludeHTMLCSS includes HTML and CSS code into the script object
	IncludeHTMLCSS(srcHTML string, srcCSS string) error

	// Control controles wether the object is valid or not
	Control() []error
}
