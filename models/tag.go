package model

import (
	"strings"
	"wooble/lib"
)

// TagSeparator specifies which character delimites a tag
const TagSeparator string = ","

// Tag is for arranging crations
type Tag struct {
	ID    uint64 `json:"id" db:"tag.id"`
	Title string `json:"title" db:"tag.title"`
}

// NewTag creates a tag
// The only error possible here is unique title constrain
// No need to warn if this error is raised
func NewTag(tag *Tag) *Tag {
	q := `
	INSERT INTO tag (title) VALUES ($1) RETURNING id
	`

	// Makes sure we split to build one tag
	oneTag := strings.Split(tag.Title, TagSeparator)

	tagLength := 18

	// A tag is not longer than tagLength characters
	lib.DB.QueryRow(q, strings.ToLower(oneTag[0])[:tagLength]).Scan(&tag.ID)

	return tag
}
