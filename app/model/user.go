package model

import (
  "database/sql"
  "time"
)

type User struct {
  ID        int64  `db:"id"`

  Email     string `db:"email"`
  Name      string `db:"name"`
  IsCreator bool   `db:"is_creator"`

  CreatedAt time.Time `db:"created_at"`
  UpdatedAt sql.NullString `db:"updated_at"`
}
