package engine

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
	"text/template"

	h "golang.org/x/net/html"

	"github.com/robertkrimen/otto"
	"github.com/woobleio/wooblizer/wbzr/engine/doc"
)

type jses5 struct {
	hasHtml bool
	name    string
	obj     *otto.Object
	src     string
}

const docVar string = "_doc()"

func NewJSES5(src string, name string) (*jses5, error) {
	if src == "" {
		src = `obj = {}`
	}
	vm := otto.New()
	obj, err := vm.Object(src)
	if err != nil {
		return nil, err
	}
	return &jses5{
		false,
		name,
		obj,
		src,
	}, nil
}

func (js *jses5) AddAttr(name string, val interface{}) error {
	if !isAcceptedFieldName(name) {
		return errors.New("The attribute name should be an alphanumerical word")
	}
	if err := js.obj.Set(name, val); err != nil {
		return err
	}
	return nil
}

func (js *jses5) AddMethod(name string, src string) error {
	if !isAcceptedFieldName(name) {
		return errors.New("This method name should be an alphanumerical word")
	}

	vm := otto.New()

	// TODO this is a workaround to build a fn with Otto
	tmpObj, err := vm.Object("({tmp:" + src + "})")
	if err != nil {
		return err
	}
	fn, err := tmpObj.Get("tmp")
	if err != nil {
		return err
	}

	if !fn.IsFunction() {
		return errors.New("The source should be a JS function(...params){}")
	}
	if err := js.obj.Set(name, fn); err != nil {
		return err
	}

	return nil
}

func (js *jses5) Build() (*template.Template, error) {
	jsw := newJsWriter(js.name, "document")

	//jsw.affectObj("", "")
	if err := buildInnerObject(js.obj, jsw); err != nil {
		return nil, err
	}

	toStr := jsw.bf.String()

	js.src = toStr

	tmpl := template.Must(template.New("jsObject").Parse(templateStr))

	return tmpl, nil
}

func (js *jses5) GetName() string { return js.name }

func (js *jses5) GetSource() string { return js.src }

func (js *jses5) IncludeHtml(src string) error {
	// Fixes net/html new line reading as text node... It breaks the generated script
	doc, err := doc.NewHTML(sanitize(src))
	if err != nil {
		return err
	}
	/*
	 * buildDoc: function(target){
	 *   var _sr = document.querySelector(target).attachShadow({mode:'open'});
	 *   // create nodes, add theirs attrs and append children to parents
	 *   this._doc = _sr; // To be query in place of the document
	 * }
	 */

	sRootVar := "_sr" // Shadow root element
	jsw := newJsWriter(sRootVar, "_d")
	jsw.makeFunction("target")
	jsw.affectVar("_d", "document")
	jsw.affectVar(sRootVar, "_d.querySelector(target).attachShadow({mode:'open'})")
	doc.ReadAndExecute(jsw.buildNode, 0)
	jsw.affectAttr("this", "_doc", "function() { return "+sRootVar+"}")
	jsw.closeExpr()

	if err := js.AddMethod("_buildDoc", jsw.bf.String()); err != nil {
		return err
	}

	js.hasHtml = true

	return nil
}

func (js *jses5) IncludeCss(css string) error {
	jsw := newJsWriter("a", docVar)

	jsw.makeFunction()
	jsw.affectVar("a", "")
	jsw.createElement("", "style")
	jsw.affectAttr("a", "innerHTML", "\""+sanitize(css)+"\"")

	jsw.appendChild("this."+jsw.doc, "a")
	jsw.closeExpr()

	err := js.AddMethod("_buildStyle", jsw.bf.String())

	return err
}

func isAcceptedFieldName(str string) bool {
	res, _ := regexp.MatchString("\\w", str)
	return res
}

func buildField(field otto.Value, jsw *jsWriter) error {
	var str string
	var err error

	switch {
	case field.IsObject() && !field.IsFunction() && field.Class() != "Array":
		jsw.makeObj()
		if err := buildInnerObject(field.Object(), jsw); err != nil {
			return err
		}
		jsw.closeExpr()
	case field.Class() == "Array":
		if str, err = formatArray(field); err != nil {
			return err
		}
	default:
		if str, err = formatVar(field); err != nil {
			return err
		}
	}
	jsw.bf.WriteString(str)

	return nil
}

func buildInnerObject(obj *otto.Object, jsw *jsWriter) error {
	keys := obj.Keys()

	for i, fieldName := range keys {
		val, err := obj.Get(fieldName)
		if err != nil {
			return err
		}

		// foo: field
		jsw.affectField(fieldName, "")
		if err := buildField(val, jsw); err != nil {
			return err
		}
		if i < len(keys)-1 {
			jsw.endField()
		}
	}

	return nil
}

func formatArray(arr otto.Value) (string, error) {
	val, errExp := arr.Export()
	if errExp != nil {
		return "", errExp
	}
	formatArr, errStr := arr.ToString()
	if errStr != nil {
		return "", errExp
	}
	switch val.(type) {
	case []string:
		var arrBf bytes.Buffer
		split := strings.Split(formatArr, ",")
		for i, val := range split {
			// ["foo","bar"]
			arrBf.WriteRune('"')
			arrBf.WriteString(val)
			arrBf.WriteRune('"')
			if i < len(split)-1 {
				arrBf.WriteRune(',')
			}
		}
		formatArr = arrBf.String()
	}
	return ("[" + formatArr + "]"), nil
}

func formatVar(pVar otto.Value) (string, error) {
	formatVar, errStr := pVar.ToString()
	if errStr != nil {
		return "", errStr
	}

	switch {
	case pVar.IsString():
		formatVar = "\"" + formatVar + "\""
	case pVar.IsFunction():
		formatVar = sanitize(formatVar)
	}

	return formatVar, nil
}

func sanitize(src string) string {
	rpcer := strings.NewReplacer("\n", "", "\t", "", "\r", "")
	return rpcer.Replace(src)
}

type jsWriter struct {
	bf   bytes.Buffer
	doc  string
	cVar string
	vars []string
}

func newJsWriter(baseVar string, docTarget string) *jsWriter {
	vars := make([]string, 0)
	vars = append(vars, baseVar)
	return &jsWriter{
		bytes.Buffer{},
		docTarget,
		baseVar,
		vars,
	}
}

func (jsw *jsWriter) affectAttr(context string, attrName string, expr string) {
	jsw.bf.WriteString(context)
	jsw.bf.WriteRune('.')
	jsw.bf.WriteString(attrName)
	jsw.bf.WriteString(" = ")
	jsw.bf.WriteString(expr)
	jsw.endExpr()
}

func (jsw *jsWriter) affectField(name string, expr string) {
	jsw.bf.WriteString(name)
	jsw.bf.WriteRune(':')
	jsw.bf.WriteString(expr)
}

func (jsw *jsWriter) affectObj(varName string, obj string) {
	if len(varName) == 0 {
		varName = jsw.cVar
	}
	jsw.bf.WriteString(varName)
	jsw.bf.WriteString("={")
	if len(obj) > 0 {
		jsw.bf.WriteString(obj)
		jsw.bf.WriteRune('}')
		jsw.endExpr()
	}
}

func (jsw *jsWriter) affectVar(varName string, expr string) {
	if len(varName) == 0 {
		varName = jsw.cVar
	}
	jsw.bf.WriteString("var ")
	jsw.bf.WriteString(varName)
	jsw.bf.WriteString(" = ")
	if len(expr) > 0 {
		jsw.bf.WriteString(expr)
		jsw.endExpr()
	}
}

func (jsw *jsWriter) appendChild(to string, toAppend string) {
	if len(toAppend) == 0 {
		toAppend = jsw.cVar
	}
	jsw.bf.WriteString(to)
	jsw.bf.WriteString(".appendChild(")
	jsw.bf.WriteString(toAppend)
	jsw.bf.WriteRune(')')
	jsw.endExpr()
}

func (jsw *jsWriter) buildNode(node *h.Node, pIndex int) int {
	jsw.genUniqueVar()
	jsw.affectVar("", "")
	switch node.Type {
	case h.ElementNode:
		jsw.createElement("_d", node.Data)
	case h.TextNode:
		jsw.createTextNode("_d", node.Data)
	}
	jsw.setAttributes(node.Attr)

	jsw.appendChild(jsw.vars[pIndex], "")

	return len(jsw.vars) - 1
}

func (jsw *jsWriter) closeExpr() {
	jsw.bf.WriteRune('}')
}

func (jsw *jsWriter) createElement(docVar string, el string) {
	if len(docVar) == 0 {
		docVar = "document"
	}
	jsw.bf.WriteString(docVar)
	jsw.bf.WriteString(".createElement(\"")
	jsw.bf.WriteString(el)
	jsw.bf.WriteString("\")")
	jsw.endExpr()
}

func (jsw *jsWriter) createTextNode(docVar string, text string) {
	if len(docVar) == 0 {
		docVar = "document"
	}
	jsw.bf.WriteString(docVar)
	jsw.bf.WriteString(".createTextNode(\"")
	jsw.bf.WriteString(text)
	jsw.bf.WriteString("\")")
	jsw.endExpr()
}

func (jsw *jsWriter) endField() {
	jsw.bf.WriteRune(',')
}

func (jsw *jsWriter) endExpr() {
	jsw.bf.WriteRune(';')
}

func (jsw *jsWriter) genUniqueVar() {
	baseNames := [26]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	tLength := len(jsw.vars)
	bLength := len(baseNames)
	if tLength >= bLength {
		mod := tLength % bLength
		time := tLength / bLength

		jsw.vars = append(jsw.vars, jsw.vars[(time-1)*bLength]+baseNames[mod])
	} else {
		jsw.vars = append(jsw.vars, baseNames[tLength])
	}
	jsw.cVar = jsw.vars[len(jsw.vars)-1]
}

func (jsw *jsWriter) setAttributes(attrs []h.Attribute) {
	var attrKey string
	for _, attr := range attrs {
		if len(attr.Namespace) > 0 {
			attrKey = attr.Namespace + ":" + attr.Key
		} else {
			attrKey = attr.Key
		}
		jsw.bf.WriteString(jsw.cVar)
		jsw.bf.WriteString(".setAttribute(\"")
		jsw.bf.WriteString(attrKey)
		jsw.bf.WriteString("\", \"")
		jsw.bf.WriteString(attr.Val)
		jsw.bf.WriteString("\")")
		jsw.endExpr()
	}
}

func (jsw *jsWriter) makeFunction(args ...string) {
	jsw.bf.WriteString("function(")
	for i, arg := range args {
		jsw.bf.WriteString(arg)
		if i < len(args)-1 {
			jsw.bf.WriteRune(',')
		}
	}
	jsw.bf.WriteString("){")
}

func (jsw *jsWriter) makeObj() {
	jsw.bf.WriteRune('{')
}

const templateStr = `{{.GetName}}={{"{"}}{{.GetSource}}{{"}"}}`
