package wbzr_test

import (
	"testing"

	"github.com/woobleio/wooblizer"
	"github.com/woobleio/wooblizer/engine"
)

func TestInject(t *testing.T) {
	wb := wbzr.New(wbzr.JS)
	if _, errs := wb.Inject("var Woobly=function Woobly(){};", "foo", make([]interface{}, 0)); len(errs) == 0 || (len(errs) > 0 && errs[0] != engine.ErrNoDocInit) {
		t.Error("Inject 1 : Should trigger an error => No document initializer")
	}

	if _, errs := wb.Inject("var Foobar = function(){function Foobar(){this.document=document}}", "bar", nil); len(errs) == 0 || (len(errs) > 0 && errs[0] != engine.ErrNoClassFound) {
		t.Error("Inject 2 : Should tigger an error => No class found")
	}

	if _, errs := wb.Inject("var Woobly = function(){}", "foobar", nil); len(errs) == 0 || (len(errs) > 0 && errs[0] != engine.ErrNoConstructor) {
		t.Error("Inject 3 : Should tigger an error => No constructor found")
	}

	// foo already exists
	if _, err := wb.Inject("otherObj={}", "foo", nil); err == nil {
		t.Error("Inject 4 : Should trigger an error => Unique alias only")
	}
}

func TestSecureAndWrap(t *testing.T) {
	wb := wbzr.New(wbzr.JS)

	var params = make([]interface{}, 2)
	params[0] = engine.JSParam{"par1", "'value1'"}
	params[1] = engine.JSParam{"par2", "'value2'"}

	script1, errs := wb.Inject(`var Woobly = function () {
	  function Woobly(params) {
	    _classCallCheck(this, Woobly);
			this.document = document.body.shadowRoot;
	  }

	  _createClass(Woobly, [{
	    key: "getDocument",
	    value: function getDocument() {
	      return document;
	    }
	  }]);

	  return Woobly;
	}();`, "obj1", params)

	if len(errs) > 0 {
		t.Error("Failed to inject the first script, error : %s", errs)
	}

	script2, _ := wb.Inject(`var Woobly = function(){function Woobly(params){_classCallCheck(this,Woobly);this.document=document.body.shadowRoot}_createClass(Woobly,[{key:"toto",value:function toto(lol){}}]);return Woobly}();`, "obj2", params)

	src := `
	<div id='divid'>
		yoyo
	</div>`

	if err := script1.IncludeHTMLCSS(src, "div { color: red; }"); err != nil {
		t.Errorf("Failed to include HTML in script1, error : %s", err)
	}

	if err := script2.IncludeHTMLCSS("", "div { color: red; }"); err != nil {
		t.Errorf("Failed to include HTML in script2, error : %s", err)
	}

	bf, errWrap := wb.SecureAndWrap("toto.com", "tata.com")
	if errWrap != nil {
		t.Error("Failed to wrap, error %s", errWrap)
	}

	t.Log(bf.String())
}
