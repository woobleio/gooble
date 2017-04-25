# Wooblizer - A library wrapper for front-end languages

Wooblizer wrap and build front-end codes in one library. Wooblizer can do :
- Packaging of HTML, CSS and JavaScript in one JavaScript file. The ouput is called a **wooble**.
- Packaging some woobles in one library.

# Usage for JSES5

```go

js1 :=
`
obj = {
  foo: "bar",
  _init: function() {
    console.log(this.foo);
  }
}
`

js2 :=
`
obj = {
  fn: function(arg) {
    document.appendChild(document.createElement("div"));
  }
}
`

wb := wbzr.New(wbzr.JSES5)
wb.Inject(js1, "firstObj")
wb.Inject(js2, "secObj")

// Create a Wooble
wb.BuildScriptFile("firstObj", "/tmp", "first_script")

sc, _ := wb.Get("firstObj")
sc.IncludeHtml("<div id='test'>hello</div>")
sc.IncludeCss("#test{ background-color: red; }")

// js1 and js2 are injected, it will package both in one Wooble library
wb.WrapAndBuildFile("/tmp", "wrap")
```

\* The field `_init` is not mandatory.

# Supported script languages

## JS ECMAScript Edition 5

A source should be an object in the following form

```js
obj = {
  foo: 1,
  "bar": function(arg1, arg2) {
    // stuff
  },
  _init: function() {
    // Used by Wooble library to execute the code when the obj is initialized
  }
}
```

The following field names are reserved :
- `_buildDoc`
- `_buildStyle`

# Supported markup languages

## HTML5

A source should only contain the content, no header, no body, no doctype, etc.

# Supported style sheet languages

## CSS

There is no restriction.
