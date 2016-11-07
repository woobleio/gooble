package model

type Source struct {
	ID uint64 `json:"id"   db:"src.id"`

	Host string `json:"host" db:"host"`
}
