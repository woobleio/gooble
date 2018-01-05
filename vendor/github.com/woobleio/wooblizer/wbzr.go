// Package wbzr provides tool to create a wooble in a given language and to package
// (wrap) some woobles.
package wbzr

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/woobleio/wooblizer/api"
	"github.com/woobleio/wooblizer/engine"
)

// ScriptLang are constants for implemented script languages.
type ScriptLang int

// Supported engines
const (
	JS ScriptLang = iota
)

// Wbzr is the wooblizer system
type Wbzr struct {
	DomainsSec []string
	Scripts    []engine.Script

	lang    ScriptLang
	api     string
	apiName string
}

// New takes a script language which is used to inject and output a file.
func New(sl ScriptLang) *Wbzr {
	var apiLib string
	var apiName string
	switch sl {
	case JS:
		apiLib = api.JS2015
		apiName = "js2015"
	default:
		panic("Language not supported")
	}

	return &Wbzr{
		nil,
		make([]engine.Script, 0),
		sl,
		apiLib,
		apiName,
	}
}

// Get returns an injected source.
func (wb *Wbzr) Get(name string) (engine.Script, error) {
	for _, sc := range wb.Scripts {
		if sc.GetName() == name {
			return sc, nil
		}
	}
	return nil, ErrUniqueName
}

// Inject injects a source code to be wooblized. It takes a name which must be
// unique. Src can be empty, it'll create a default object
func (wb *Wbzr) Inject(src string, name string, params []interface{}) (engine.Script, []error) {
	errs := make([]error, 0)
	if _, err := wb.Get(name); err == nil {
		errs = append(errs, err)
		return nil, errs
	}
	var sc engine.Script

	switch wb.lang {
	case JS:
		var jsParams = make([]engine.JSParam, len(params))
		for i, p := range params {
			jsParams[i] = p.(engine.JSParam)
		}
		sc, errs = engine.NewJS(name, src, jsParams)
	}

	if len(errs) > 0 {
		return sc, errs
	}

	wb.Scripts = append(wb.Scripts, sc)

	return sc, errs
}

// InjectFile injects a source from a file.
func (wb *Wbzr) InjectFile(path string, name string, params []interface{}) (engine.Script, []error) {
	errs := make([]error, 0)
	c, err := ioutil.ReadFile(path)
	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	return wb.Inject(string(c[:]), name, params)
}

// Secure set some domains to protect the script and make it works only for specific domains
func (wb *Wbzr) Secure(domains ...string) {
	wb.DomainsSec = domains
}

// SecureAndWrap wrap all scripts in the wooblizer and secure it with domains
func (wb *Wbzr) SecureAndWrap(domains ...string) (*bytes.Buffer, error) {
	wb.Secure(domains...)
	return wb.Wrap()
}

// Wrap packages some creations (all the creations injected in the Wbzr)
// and build a file which contains the wooble library.
func (wb *Wbzr) Wrap() (*bytes.Buffer, error) {
	fns := template.FuncMap{
		"plus1": func(x int) int {
			return x + 1
		},
	}

	tmpl := template.Must(template.New(wb.apiName).Funcs(fns).Parse(wb.api))

	var out bytes.Buffer
	if err := tmpl.Execute(&out, wb); err != nil {
		return nil, err
	}

	return &out, nil
}

// WooblyJS is a Wooble base code for creation
var WooblyJS = `class Woobly {

	constructor(params) {
		// This is mandatory.
		// Use this.document.querySelector to query elements in your creation
		// Use document to call document prototypes such as document.createElement
		this.document = document.body.shadowRoot;

		/*
		 * Your creation start-up code
		 */
	}

	/*
	 * You can create all methods you need
	 */
}`
