// Package wbzr provides tool to create a wooble in a given language and to package
// (wrap) some woobles.
package wbzr

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/woobleio/wooblizer/wbzr/engine"
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
  	{{range $i, $o := .Scripts}}
			"{{$o.GetName}}":{{$o.GetSource}},
			"__{{$o.GetName}}":{
			{{$lenParams := len $o.Params}}
			{{range $i, $p := $o.Params}}
				"{{$p.Field}}":{{$p.Value}}{{if ne (plus1 $i) $lenParams}},{{end}}
			{{end}}
			}{{if ne (plus1 $i) $lenScripts}},{{end}}
		{{end}}
  }

  var c = cs[id];
  if(typeof c == 'undefined') {
  	console.log("Wooble error : creation", id, "not found");
    return undefined;
  }

  this.init = function (tar, p) {
    if(document.querySelector(tar) == null) {
    	console.log("Wooble error : Element", tar, "not found in the document");
      return;
    }

		if (p) {
			var _ = cs['__'+id];
			for (prop in p) {
				if (_.hasOwnProperty(prop)) _[prop] = p[prop];
			}
			p = _;
		} else p = cs['__'+id];

		var t = this;
    return new Promise(function(r, e) {
      if (!document.head.attachShadow) {
        // Browsers shadow dom support with polyfill
        var s = document.createElement('script');
        s.type = 'text/javascript';
        s.src = 'https://cdnjs.cloudflare.com/ajax/libs/webcomponentsjs/1.0.0-rc.11/webcomponents-lite.js';
        document.getElementsByTagName('head')[0].appendChild(s);
        s.onload = function() {
          r(new c(tar,p));
        }
      } else {
        r(new c(tar,p));
      }
    });
  }

  return this;
}
`
