package engine

import (
	"text/template"
)

type Script interface {
	AddAttr(name string, val interface{}) error
	AddMethod(name string, src string) error
	GetName() string
	GetSource() string
	Build() (*template.Template, error)
	IncludeCss(css string) error
	IncludeHtml(src string) error
}
