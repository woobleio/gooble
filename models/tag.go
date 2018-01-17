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
	NbUse uint64 `json:"nbUse" db:"tag_nb_use"`
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

	// if unique constrain error (on tag title) then get the tag
	if err := lib.DB.QueryRow(q, oneTag).Scan(&tag.ID); err != nil {
		q = `SELECT id "tag.id", title "tag.title" FROM tag WHERE title = $1`
		lib.DB.Get(tag, q, oneTag)
	}
}

// AllTags returns all tags
// Filters : search
// Orders : nb_use
func AllTags(opt lib.Option) ([]Tag, error) {
	var tags []Tag
	q := lib.NewQuery(`SELECT
		t.id "tag.id",
		t.title "tag.title",
		COUNT(ct.tag_id) AS tag_nb_use
	FROM tag t
	LEFT OUTER JOIN creation_tag ct ON (ct.tag_id=t.id)
	`, &opt)

	q.SetFilters(false, lib.SEARCH, "t.title")

	q.Q += `GROUP BY t.id`

	q.SetOrder(lib.NB_USE, "tag_nb_use")

	return tags, lib.DB.Select(&tags, q.String(), q.Values...)
}
