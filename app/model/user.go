package model

import (
  "wooblapp/app/util"
)

type User struct {
  ID        uint64  `json:"id"                 db:"user.id"`

  Email     string `json:"email,omitempty"     db:"email"`
  Name      string `json:"name"                db:"name"`
  IsCreator bool   `json:"isCreator,omitempty" db:"is_creator"`

  CreatedAt *util.NullTime `json:"createdAt,omitempty" db:"user.created_at"`
  UpdatedAt *util.NullTime `json:"updatedAt,omitempty" db:"user.updated_at"`
}
