package model

// Email is for sending emails (used in templates)
type Email struct {
	Title    string
	Name     string
	Elements []Element
}

// Element represents an HTML element (DOM node)
type Element struct {
	Type    string
	Content string
	Attr    string
}
