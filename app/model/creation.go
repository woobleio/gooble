package model

import (
  "database/sql"
  "time"

  "wooblapp/lib"
  "wooblapp/app/util"
)

type Creation struct {
  ID        int64   `db:"id"`

  CreatorID int64   `db:"creator_id"`
  Creator   User    `db:""`
  SourceID  int64   `db:"source_id"`
  Source    Source  `db:""`
  Title     string  `db:"title"`
  Version   string  `db:"version"`

  CreatedAt time.Time `db:"created_at"`
  UpdatedAt sql.NullString `db:"updated_at"`
}

func AllCreations(opt util.Option) (*[]Creation, error) {
  var creations []Creation
  q := util.Query{`
    SELECT
      c.*,
      s.*,
      u.id,
      u.email,
      u.name,
      u.is_creator,
      u.created_at,
      u.updated_at
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
      c.*,
      s.*,
      u.id,
      u.email,
      u.name,
      u.is_creator,
      u.created_at,
      u.updated_at
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
