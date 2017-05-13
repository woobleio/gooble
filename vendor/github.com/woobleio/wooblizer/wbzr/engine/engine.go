package engine

// StdName is the name standard for a wooble object
const StdName string = "Woobly"

// Script is an interface for script langage
type Script interface {
	GetName() string
	GetSource() string
	IncludeHTMLCSS(srcHTML string, srcCSS string) error
}
