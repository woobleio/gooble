package api

// JS2015 is native JS API
var JS2015 = `
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
		var _cs = [];
    return new Promise(function(r, e) {
      if (!document.head.attachShadow) {
        // Browsers shadow dom support with polyfill
        var s = document.createElement('script');
        s.type = 'text/javascript';
        s.src = 'https://cdnjs.cloudflare.com/ajax/libs/webcomponentsjs/1.0.14/webcomponents-sd-ce.js';
        document.getElementsByTagName('head')[0].appendChild(s);
        s.onload = function() {
					for (var d of document.querySelectorAll(tar)) _cs.push(new c(d,p));
          r(_cs);
        }
      } else {
				for (var d of document.querySelectorAll(tar)) _cs.push(new c(d,p));
        r(_cs);
      }
    });
  }

  return this;
}
`
