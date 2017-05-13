// Package wbzr provides tool to create a wooble in a given language and to package
// (wrap) some woobles.
package wbzr

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	"text/template"

	"github.com/woobleio/wooblizer/wbzr/engine"
)

// ScriptLang are constants for implemented script languages.
type ScriptLang int

// Supported engines
const (
	JS ScriptLang = iota
)

type Wbzr struct {
	DomainsSec []string
	Scripts    []engine.Script

	lang     ScriptLang
	skeleton string
}

// New takes a script language which is used to inject and output a file.
func New(sl ScriptLang) *Wbzr {
	var skeleton string
	switch sl {
	case JS:
		skeleton = wbJS
	default:
		panic("Language not supported")
	}

	return &Wbzr{
		nil,
		make([]engine.Script, 0),
		sl,
		skeleton,
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
func (wb *Wbzr) Inject(src string, name string) (engine.Script, error) {
	if _, err := wb.Get(name); err == nil {
		return nil, ErrUniqueName
	}
	var sc engine.Script
	var err error

	switch wb.lang {
	case JS:
		rmVar := regexp.MustCompile(`^var Woobly[ ]?=`)
		src = rmVar.ReplaceAllString(src, "")
		src = strings.TrimRight(src, ";")
		sc = &engine.JS{
			Src:  src,
			Name: name,
		}
	}

	if err != nil {
		return nil, err
	}

	wb.Scripts = append(wb.Scripts, sc)
	return sc, nil
}

// InjectFile injects a source from a file.
func (wb *Wbzr) InjectFile(path string, name string) (engine.Script, error) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return wb.Inject(string(c[:]), name)
}

// Secure set some domains to protect the script and make it works only for specific domains
func (wb *Wbzr) Secure(domains ...string) {
	wb.DomainsSec = domains
}

func (wb *Wbzr) SecureAndWrap(domains ...string) (*bytes.Buffer, error) {
	wb.Secure(domains...)
	return wb.Wrap()
}

// Wrap packages some woobles (all the woobles injected in the Wbzr)
// and build a file which contains the wooble library.
func (wb *Wbzr) Wrap() (*bytes.Buffer, error) {
	fns := template.FuncMap{
		"plus1": func(x int) int {
			return x + 1
		},
	}

	tmpl := template.Must(template.New("wbJS").Funcs(fns).Parse(wbJS))

	var out bytes.Buffer
	if err := tmpl.Execute(&out, wb); err != nil {
		return nil, err
	}

	return &out, nil
}

// WooblyJS is a Wooble creation template for JS
var WooblyJS = `class Woobly {

	constructor() {
		// This is mandatory. Use this.document to query elements in your creation
		// (this.document.queryAll('div')), use document for manipulating the
		// document parent (document.createElement('div'))
		this.document = document;

		/*
		 * Your creation start-up code
		 */
	}

	/*
	 * You can create all methods you need
	 */
}`

var wbJS = `
var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

function Wb(id) {
	{{if .DomainsSec}}
	{{$lenDoms := len .DomainsSec}}
	var ah = [{{range $i, $o := .DomainsSec}}"{{$o}}"{{if ne (plus1 $i) $lenDoms}},{{end}}{{end}}];
  var xx = ah.indexOf(window.location.hostname);
  if(ah.indexOf(window.location.hostname) == -1) {
  	console.log("Wooble error : domain restricted");
    return;
  }
	{{end}}

	if(window === this) {
  	return new Wb(id);
  }

  var cs = {
		{{$lenScripts := len .Scripts}}
  	{{range $i, $o := .Scripts}}"{{$o.GetName}}":{{$o.GetSource}}{{if ne (plus1 $i) $lenScripts}},{{end}}{{end}}
  }

  var c = cs[id];
  if(typeof c == 'undefined') {
  	console.log("Wooble error : creation", id, "not found");
    return undefined;
  }

  this.init = function (tar) {
    if(document.querySelector(tar) == null) {
    	console.log("Wooble error : Element", target, "not found in the document");
      return;
    }

		var t = this;
    return new Promise(function(r, e) {
      if (!document.head.attachShadow) {
        // Browsers shadow dom support with polyfill
        var s = document.createElement('script');
        s.type = 'text/javascript';
        s.src = 'https://cdnjs.cloudflare.com/ajax/libs/webcomponentsjs/1.0.0-rc.11/webcomponents-lite.js';
        document.getElementsByTagName('head')[0].appendChild(s);
        s.onload = function() {
          r(new c(tar));
        }
      } else {
        r(new c(tar));
      }
    });
  }

  return this;
}
`
