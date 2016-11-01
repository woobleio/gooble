package model

type Source struct {
  ID   int64 `db:"id"`

  Host string `db:"host"`
}
