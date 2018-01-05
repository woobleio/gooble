package model

import (
	"strings"
	"wooble/lib"
)

// TagSeparator specifies which character delimites a tag
const TagSeparator string = ","

// Tag is for arranging crations
type Tag struct {
	ID    uint64 `json:"id,omitempty" db:"tag.id"`
	Title string `json:"title" db:"tag.title"`
}

// NewOrGetTag creates a tag or get one if exists
func NewOrGetTag(tag *Tag) {
	q := `
	INSERT INTO tag (title) VALUES ($1) RETURNING id
	`

	// Makes sure we split to build one tag
	oneTag := strings.Split(tag.Title, TagSeparator)[0]

	tagLength := 18

	// A tag is not longer than tagLength characters
	if len(oneTag) > tagLength {
		oneTag = oneTag[:tagLength]
	}

	oneTag = strings.ToLower(oneTag)

	if err := lib.DB.QueryRow(q, oneTag).Scan(&tag.ID); err != nil {
		q = `SELECT id "tag.id", title "tag.title" FROM tag WHERE title = $1`
		lib.DB.Get(tag, q, oneTag)
	}
}
