package lib

import (
  _ "github.com/lib/pq"
  "github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func InitDB(url string) {
  DB = sqlx.MustOpen("postgres", url)
}
