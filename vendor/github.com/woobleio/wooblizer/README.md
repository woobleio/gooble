# Wooblizer - A wrapped library for front-end languages

Wooblizer wrap and build front-end codes and outputs a Wooble library.

[What is Wooble?](https://github.com/woobleio/wooble/blob/master/doc/whitepaper.md)

Wooblizer what allows to packed creations in one single library for using it in a website or an application.

# Usage for JS ES2015

```go

// Babelified creation to es2015
js1 :=
`
var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var Woobly = function () {
  function Woobly(params) {
    _classCallCheck(this, Woobly);

    this.document = document.body.shadowRoot;
  }

  _createClass(Woobly, [{
    key: "doSomething",
    value: function doSomething(something) {
      console.log(something);
    }
  }]);

  return Woobly;
}();
`

// Babelified creation to es2015
js2 :=
`
// Another creation
function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var Woobly = function Woobly(params) {
	_classCallCheck(this, Woobly);

	this.document = document.body.shadowRoot;

	console.log(params);
};
`

wb := wbzr.New(wbzr.JSES5) // JSES5 is a constant to specify in which standard or langage you want to build your Wooble

sc1 := wb.Inject(js1, "firstObj")
wb.Inject(js2, "secObj")

// Adds HTML and CSS code in creation js1
sc1.IncludeHTMLCSS("<div id='test'>hello</div>", "#test{ background-color: red; }")

// js1 and js2 are injected, it will package both in the Wooble API. It returns a buffer
// containing the source code of the API
bf, err := wb.Wrap()
```

# Supported script languages and frameworks

Wooble consider two types of engines, as everything if very different, I choose
to make a structure that separates them through the Engine interface.

For instance the Engine JS is for JavaScript ES2015 only. It cant handle ES6, it
has to be [babelified](https://babeljs.io/repl/) to ES2015 to be compatible with this engine.

An Engine can be a framework such as Angular4 even though we might think that it
is a ES6 code that could fit in an Engine that we might call "JS6". Angular4 has some
features to the point it's a better design to define an engine specifically for
this framework. So let's apply this rule for everything that embed some useful features
and allow Wooble to offer the best experience for every frameworks / libraries / languages.

## JS ECMAScript Edition 5 (JS ES2015)

A source should be an object in the following form

```js
var Woobly = function () {
  function Woobly(params) {
    _classCallCheck(this, Woobly);

    this.document = document.body.shadowRoot; // this is mandatory
  }

  _createClass(Woobly, [{
    key: "doSomething",
    value: function doSomething(something) {
      console.log(something);
    }
  }]);

  return Woobly;
}();
```

Or, it is possible to make a class with JavaScript ES6 and to "babelify" it to ES2015 in order to be processed by Wooble (wooblelized).

```js
class Woobly {

  constructor(params) {
    this.document = document.body.shadowRoot;

    console.log(params);
  }
}
```

# Supported markup languages

## HTML5

A source should only contain the content, no header, no body, no doctype, etc.

# Supported style sheet languages

## CSS

There is no restriction.

[Contributing](https://github.com/woobleio/wooblizer/blob/master/CONTRIBUTING.md)
