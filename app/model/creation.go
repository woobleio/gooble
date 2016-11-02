package model

import (
  "wooblapp/lib"
  "wooblapp/app/util"
)

type Creation struct {
  ID        uint64  `json:"id"      db:"crea.id"`

  CreatorID uint64  `json:"-"       db:"creator_id"`
  Creator   User    `json:"creator" db:""`
  SourceID  uint64  `json:"-"       db:"source_id"`
  Source    Source  `json:"source"  db:""`
  Title     string  `json:"title"   db:"title"`
  Version   string  `json:"version" db:"version"`

  CreatedAt *util.NullTime `json:"createdAt,omitempty" db:"crea.created_at"`
  UpdatedAt *util.NullTime `json:"updatedAt,omitempty" db:"crea.updated_at"`
}

func AllCreations(opt util.Option) (*[]Creation, error) {
  var creations []Creation
  q := util.Query{`
    SELECT
      c.id "crea.id",
      c.title,
      c.created_at "crea.created_at",
      c.updated_at "crea.updated_at",
      c.version,
      s.id "src.id",
      s.host,
      u.id "user.id",
      u.name
    FROM creation c
    INNER JOIN source s ON (c.source_id = s.id)
    INNER JOIN app_user u ON (c.creator_id = u.id)`,
    &opt,
  }

  query := q.String()

  if err := lib.DB.Select(&creations, query); err != nil {
    return nil, err
  }

  return &creations, nil
}

func CreationByTitle(title string, opt util.Option) (*Creation, error) {
  var crea Creation
  q := util.Query{`
    SELECT
      c.id "crea.id",
      c.title,
      c.created_at "crea.created_at",
      c.updated_at "crea.updated_at",
      c.version,
      s.id "src.id",
      s.host,
      u.id "user.id",
      u.name
    FROM creation c
    INNER JOIN source s ON (c.source_id = s.id)
    INNER JOIN app_user u ON (c.creator_id = u.id)
    WHERE c.title = $1`,
    &opt,
  }

  query := q.String()

  if err := lib.DB.Get(&crea, query, title); err != nil {
    return nil, err
  }

  return &crea, nil
}
