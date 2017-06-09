package engine

import (
	"bytes"
	"errors"
	"regexp"
	"strings"

	h "golang.org/x/net/html"

	"github.com/woobleio/wooblizer/wbzr/engine/doc"
)

// JS Object
type JS struct {
	Name   string
	Src    string
	Params []JSParam
}

// JSParam is a object parameter
type JSParam struct {
	Field string
	Value string
}

const docVar string = "this.document"

const (
	docRegex         string = `this.document[ ]?=[ ]?document[;]?`
	constructorRegex string = `.*function Woobly\(`
	classRegex       string = `var Woobly[ ]?=`
)

// NewJS initializes a JS and returns errors if the class isn't in the standard
func NewJS(name string, src string, params []JSParam) (*JS, []error) {
	js := &JS{
		Name:   name,
		Src:    src,
		Params: params,
	}

	errs := js.Control()

	rmVar := regexp.MustCompile(classRegex)
	src = rmVar.ReplaceAllString(src, "")
	src = strings.TrimRight(src, ";")

	js.Src = src

	return js, errs
}

// GetName returns obj name
func (js *JS) GetName() string { return js.Name }

// GetSource returns obj code source
func (js *JS) GetSource() string { return js.Src }

// GetParams returns obj parameters
func (js *JS) GetParams() []interface{} {
	var interf = make([]interface{}, len(js.Params))
	for i, p := range js.Params {
		interf[i] = p
	}
	return interf
}

// IncludeHTMLCSS includes HTML and CSS in the object
func (js *JS) IncludeHTMLCSS(srcHTML string, srcCSS string) error {
	// Fixes net/html new line reading as text node... It breaks the generated script
	doc, err := doc.NewHTML(sanitize(srcHTML))
	if err != nil {
		return errors.New("DOM error : " + err.Error())
	}

	initDocRegex := regexp.MustCompile(docRegex)
	if !initDocRegex.MatchString(js.Src) {
		return errors.New("No document initilization found. this.document = document is required in the object consctructor")
	}

	sRootVar := "_sr_" // Shadow root element
	jsw := newJsWriter(sRootVar)
	if srcHTML != "" {
		constructorRegex := regexp.MustCompile(constructorRegex)
		constructorIdx := constructorRegex.FindIndex([]byte(js.Src))
		srcToBytes := []byte(js.Src)
		index := constructorIdx[1]

		coma := ","
		if string(srcToBytes[index:index+1]) == ")" {
			coma = ""
		}
		// Insert target parameter in the object constructor
		js.Src = string(append(srcToBytes[:index], append([]byte("_t_"+coma), srcToBytes[index:]...)...))

		jsw.affectVar(sRootVar, "document.querySelector(_t_).attachShadow({mode:'open'})")
		doc.ReadAndExecute(jsw.buildNode, 0)
		jsw.affectAttr("this", "document", sRootVar)
	}

	if srcCSS != "" {
		styleVar := "__s"
		jsw.affectVar(styleVar, "")
		jsw.createElement("style")
		jsw.affectAttr(styleVar, "innerHTML", "'"+sanitize(sanitizeString(srcCSS))+"'")

		if srcHTML != "" {
			jsw.appendChild(docVar, styleVar)
		} else {
			// Rewrite document initialization to override the replaceAll
			jsw.bf.WriteString("this.document=document;")
			jsw.appendChild(docVar+".head", styleVar)
		}
	}

	js.Src = string(initDocRegex.ReplaceAll([]byte(js.Src), jsw.bf.Bytes()))

	return nil
}

// Control checks if the class is valid
func (js *JS) Control() []error {
	docR := regexp.MustCompile(docRegex)
	constR := regexp.MustCompile(constructorRegex)
	classR := regexp.MustCompile(classRegex + "[ ]?function.*")

	errs := make([]error, 0)

	if !classR.MatchString(js.GetSource()) {
		errs = append(errs, ErrNoClassFound)
	}
	if !constR.MatchString(js.GetSource()) {
		errs = append(errs, ErrNoConstructor)
	}
	if !docR.MatchString(js.GetSource()) {
		errs = append(errs, ErrNoDocInit)
	}

	return errs
}

func sanitize(src string) string {
	rpcer := strings.NewReplacer("\n", "", "\t", "", "\r", "")
	return rpcer.Replace(src)
}

func sanitizeString(src string) string {
	rpcer := strings.NewReplacer("'", "\\'")
	return rpcer.Replace(src)
}

type jsWriter struct {
	bf      bytes.Buffer
	baseVar string
	cVar    string
	vars    []string
}

func newJsWriter(baseVar string) *jsWriter {
	vars := make([]string, 0)
	vars = append(vars, baseVar)
	return &jsWriter{
		bytes.Buffer{},
		baseVar,
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
		jsw.createElement(node.Data)
	case h.TextNode:
		jsw.createTextNode(node.Data)
	}
	jsw.setAttributes(node.Attr)

	jsVar := jsw.vars[pIndex]
	if jsVar != jsw.baseVar {
		jsVar = "__" + jsVar
	}
	jsw.appendChild(jsVar, "")

	return len(jsw.vars) - 1
}

func (jsw *jsWriter) createElement(el string) {
	jsw.bf.WriteString("document.createElement('")
	jsw.bf.WriteString(el)
	jsw.bf.WriteString("')")
	jsw.endExpr()
}

func (jsw *jsWriter) createTextNode(text string) {
	jsw.bf.WriteString("document.createTextNode('")
	jsw.bf.WriteString(text)
	jsw.bf.WriteString("')")
	jsw.endExpr()
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
	jsw.cVar = "__" + jsw.vars[len(jsw.vars)-1]
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
		jsw.bf.WriteString(".setAttribute('")
		jsw.bf.WriteString(attrKey)
		jsw.bf.WriteString("', '")
		jsw.bf.WriteString(attr.Val)
		jsw.bf.WriteString("')")
		jsw.endExpr()
	}
}
