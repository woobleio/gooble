package model

import (
	"fmt"
	"strings"
	"wooble/lib"
	enum "wooble/models/enums"
)

// Creation is a Wooble object
type Creation struct {
	ID lib.ID `json:"id"      db:"crea.id"`

	Title       string          `json:"title"  db:"title"`
	Description *lib.NullString `json:"description,omitempty" db:"description"`
	ThumbPath   *lib.NullString `json:"thumbPath,omitempty" db:"thumb_path"`
	Creator     User            `json:"creator,omitempty" db:""`

	// When a creation is in draft, the very last version is ignored by most queries,
	// thus the -1 from versions array within most queries
	Versions lib.UintSlice `json:"versions,omitempty" db:"versions"`

	Alias      string `json:"alias,omitempty" db:"alias"`
	State      string `json:"state,omitempty" db:"state"`
	IsOwner    bool   `json:"isOwner,omitempty" db:"is_owner"`
	IsFeatured bool   `json:"-" db:"is_featured"`

	NbUse uint64 `json:"nbUse" db:"nb_use"`

	Params     []CreationParam `json:"params,omitempty" db:""`
	PreviewURL string          `json:"previewUrl,omitempty"`
	Version    uint64          `json:"version,omitempty"`

	Script       string `json:"script,omitempty"`
	ParsedScript string `json:"parsedScript,omitempty"`
	Document     string `json:"document,omitempty"`
	Style        string `json:"style,omitempty"`

	CreatorID uint64 `json:"-" db:"creator_id"`
	Engine    Engine `json:"-" db:""`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"crea.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"crea.updated_at"`
}

// CreationParam is a creation parameter
type CreationParam struct {
	Field string `json:"field" db:"field"`
	Value string `json:"value" db:"value"`
}

// BaseVersion is creation default version
const BaseVersion uint64 = 1

const lastVersionQuery = `SELECT
CASE WHEN state = 'draft' THEN versions[array_length(versions, 1)-1]
ELSE versions[array_length(versions, 1)]
END AS version
FROM creation WHERE id = $1`

// PopulateParams populates creation's parameters (for previewing and building into package)
func (c *Creation) PopulateParams() {
	q := `SELECT field, value FROM creation_param WHERE creation_id = $1 AND version`
	if c.Version == 0 {
		q += `(` + lastVersionQuery + `)`
		lib.DB.Select(&c.Params, q, c.ID)
	} else {
		q += `= $2`
		lib.DB.Select(&c.Params, q, c.ID, c.Version)
	}
}

// RetrieveSourceCode request the source files in the cloud and set the content to the Creation
func (c *Creation) RetrieveSourceCode(version string, files ...string) error {
	uIDStr := fmt.Sprintf("%d", c.CreatorID)
	creaIDStr := fmt.Sprintf("%d", c.ID.ValueDecoded)

	storage := lib.NewStorage(lib.SrcCreations)
	for _, f := range files {
		source := storage.GetFileContent(uIDStr, creaIDStr, version, f)
		switch f {
		case enum.Script:
			c.Script = source
		case enum.Document:
			c.Document = source
		case enum.Style:
			c.Style = source
		case enum.ParsedScript:
			c.ParsedScript = source
		}
	}

	return storage.Error()
}

// AllCreations returns all creations
func AllCreations(opt lib.Option, uID uint64) (*[]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
	    c.id "crea.id",
			COUNT(pc.creation_id) AS nb_use,
	    c.title,
			c.description,
			c.thumb_path,
	    c.created_at "crea.created_at",
	    c.updated_at "crea.updated_at",
	    CASE WHEN c.state = 'public' THEN versions[0:array_length(c.versions, 1)] ELSE versions[0:array_length(c.versions, 1)-1] END AS versions,
			c.alias,
			CASE WHEN c.creator_id = $1 THEN true ELSE false END "is_owner",
			c.state,
			e.name "eng.name",
			e.extension,
			e.content_type,
	    u.id "user.id",
	    u.name
	  FROM creation c
	  INNER JOIN app_user u ON (c.creator_id = u.id)
		INNER JOIN engine e ON (c.engine=e.name)
		LEFT JOIN package_creation pc ON (pc.creation_id = c.id)
		WHERE c.state = 'public' OR array_length(versions, 1) > 1
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(lib.SEARCH, "c.title|u.name", lib.CREATOR, "u.name")

	q.Q += "GROUP BY c.id, u.id, e.name"

	q.SetOrder(lib.CREATED_AT, "c.created_at")

	return &creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// AllPopularCreations returns all popular creations
func AllPopularCreations(opt lib.Option, uID uint64) (*[]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
	    c.id "crea.id",
	    c.title,
			c.description,
			c.thumb_path,
	    c.created_at "crea.created_at",
	    c.updated_at "crea.updated_at",
	    CASE WHEN c.state = 'public' THEN c.versions[0:array_length(c.versions, 1)] ELSE c.versions[0:array_length(c.versions, 1)-1] END AS versions,
			c.alias,
			CASE WHEN c.creator_id = $1 THEN true ELSE false END "is_owner",
			c.state,
			COUNT(pc.creation_id) AS nb_use,
	    u.id "user.id",
	    u.name
	  FROM creation c
	  INNER JOIN app_user u ON (c.creator_id = u.id)
		LEFT JOIN package_creation pc ON (pc.creation_id = c.id)
		WHERE c.state = 'public' OR array_length(versions, 1) > 1
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(lib.SEARCH, "c.title|u.name", lib.CREATOR, "u.name")

	q.Q += "GROUP BY c.id, u.id ORDER BY c.is_featured DESC, nb_use DESC"

	return &creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// AllUsedCreations return creations used in some packages
func AllUsedCreations(opt lib.Option, uID uint64) (*[]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
			DISTINCT c.id "crea.id",
			COUNT(pcc.creation_id) AS nb_use,
			c.title,
			c.thumb_path,
			c.created_at "crea.created_at",
			c.versions,
			c.thumb_path,
			u.id "user.id",
			u.name
		FROM creation c
		INNER JOIN app_user u ON (c.creator_id = u.id)
		INNER JOIN package_creation pc ON (pc.creation_id = c.id)
		LEFT JOIN package_creation pcc ON (pcc.creation_id = c.id)
    INNER JOIN package p ON (p.id = pc.package_id)
		WHERE p.user_id = $1
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(lib.SEARCH, "c.title")

	q.Q += "GROUP BY c.id, u.id"

	return &creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// AllDraftCreations returns all creation in draft of authenticated user
func AllDraftCreations(opt lib.Option, uID uint64) (*[]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
			c.id "crea.id",
			COUNT(pc.creation_id) AS nb_use,
			c.title,
			c.thumb_path,
			c.created_at "crea.created_at",
			c.versions[array_length(c.versions, 1)] AS version,
			c.versions[array_length(c.versions, 1)-1] AS versions,
			c.state,
			u.id "user.id",
			u.name
		FROM creation c
		INNER JOIN app_user u ON (c.creator_id = u.id)
		LEFT JOIN package_creation pc ON (pc.creation_id = c.id)
		WHERE u.id = $1 AND c.state = 'draft'
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(lib.SEARCH, "c.title")

	q.Q += "GROUP BY c.id, u.id"

	return &creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// CreationByID returns a creation with the id "id"
func CreationByID(id lib.ID, uID uint64) (*Creation, error) {
	var crea Creation
	q := `
  SELECT
    c.id "crea.id",
    c.title,
		c.thumb_path,
		c.description,
    c.created_at "crea.created_at",
    c.updated_at "crea.updated_at",
		CASE WHEN c.state = 'public' THEN versions[0:array_length(c.versions, 1)] ELSE versions[0:array_length(c.versions, 1)-1] END AS versions,
		c.alias,
		CASE WHEN c.creator_id = $2 THEN true ELSE false END "is_owner",
		c.state,
		e.name "eng.name",
		e.extension,
		e.content_type,
    u.id "user.id",
    u.name,
		u.pic_path
  FROM creation c
  INNER JOIN app_user u ON (c.creator_id = u.id)
	INNER JOIN engine e ON (c.engine=e.name)
  WHERE c.id = $1
	`

	if err := lib.DB.Get(&crea, q, id, uID); err != nil {
		return nil, err
	}

	crea.Version = crea.Versions[len(crea.Versions)-1] // Latest version

	crea.PopulateParams()
	return &crea, nil
}

// CreationPrivateByID returns a creation as private
func CreationPrivateByID(uID uint64, creaID lib.ID, version string) (*Creation, error) {
	var crea Creation
	q := `
	SELECT
		c.id "crea.id",
		c.title,
		c.description,
		c.thumb_path,
		c.alias,
		c.creator_id,
		c.created_at "crea.created_at",
		c.updated_at "crea.updated_at",
		c.versions,
		CASE WHEN c.state = 'draft' THEN versions[array_length(c.versions, 1)-$3] ELSE versions[array_length(c.versions, 1)] END AS version,
		c.state,
		e.name "eng.name",
		e.extension,
		e.content_type
	FROM creation c
	INNER JOIN engine e ON (c.engine=e.name)
	WHERE creator_id = $1 AND c.id = $2 AND c.state != 'delete'
	`

	// decVersion 0 means it'll get the latest version,
	// decVersion 1 means it'll get the latest version - 1,
	decVersion := 1
	if version == "latest" {
		decVersion = 0
	}

	if err := lib.DB.Get(&crea, q, uID, creaID, decVersion); err != nil {
		return nil, err
	}

	crea.PopulateParams()

	return &crea, nil
}

// UpdateCreation update creation's information
func UpdateCreation(crea *Creation) error {
	q := `
  UPDATE creation
  SET title = $3, description = $4, state = $5, alias = $6, thumb_path = $7
  WHERE id = $1
  AND creator_id = $2
  `

	if _, err := lib.DB.Exec(q, crea.ID, crea.CreatorID, crea.Title, crea.Description, crea.State, crea.Alias, crea.ThumbPath); err != nil {
		return err
	}

	if len(crea.Params) > 0 {
		return UpdateCreationParams(crea)
	}
	return nil
}

// UpdateCreationParams update all creation params for a given version (crea.Version)
func UpdateCreationParams(crea *Creation) error {
	lastVersion := crea.Version
	if lastVersion == 0 {
		if err := lib.DB.Get(&lastVersion, lastVersionQuery, crea.ID); err != nil {
			return err
		}
	}

	q := `
	DELETE FROM creation_param
	WHERE creation_id = $1
	AND version = $2
	`
	lib.DB.Exec(q, crea.ID, lastVersion)

	// Bulk insert
	index := 3
	vals := []interface{}{crea.ID, lastVersion}
	q = `INSERT INTO creation_param(creation_id, version, field, value) VALUES`
	for _, p := range crea.Params {
		q += `($1, $2, $` + fmt.Sprintf("%d", index) + `, $` + fmt.Sprintf("%d", index+1) + `),`
		vals = append(vals, p.Field, p.Value)
		index += 2
	}
	q = strings.TrimRight(q, ",")

	_, err := lib.DB.Exec(q, vals...)
	return err
}

// CreationByIDAndVersion returns the creation "creaID" and check if the version "version" exists
func CreationByIDAndVersion(id lib.ID, version uint64) (*Creation, error) {
	crea := new(Creation)
	if version == 0 {
		version = BaseVersion
	}
	q := `
  SELECT
		id "crea.id",
		CASE WHEN state = 'draft' THEN versions[0:array_length(versions, 1)-1] ELSE versions END AS versions,
		state
  FROM creation
	WHERE id = $1
  AND $2 = ANY (versions)
  `
	return crea, lib.DB.Get(crea, q, id, version)
}

// CreationLastVersion gets creation's last version
func CreationLastVersion(id lib.ID) int64 {
	var version int64
	q := `SELECT versions[array_length(versions,1)] AS version FROM creation WHERE id = $1`
	lib.DB.Get(&version, q, id)
	return version
}

// NewCreation creates a creation
func NewCreation(crea *Creation) (*Creation, error) {
	q := `
  INSERT INTO creation(
    title,
    creator_id,
    versions,
    engine,
		state,
		alias
  ) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
  `

	stringSliceVersions := make(lib.StringSlice, 0, 1)

	if crea.Alias == "" {
		crea.Alias = "woobly"
	}

	return crea, lib.DB.QueryRow(q, crea.Title, crea.CreatorID, append(stringSliceVersions, fmt.Sprintf("%d", BaseVersion)), crea.Engine.Name, crea.State, crea.Alias).Scan(&crea.ID)
}

// CopyCreationParams copie creation parameters from current version to a new version
func CopyCreationParams(id lib.ID, curVersion uint64, newVersion uint64) error {
	q := `
	INSERT INTO creation_param(creation_id, version, field, value)
	SELECT creation_id, $3 AS version, field, value
	FROM creation_param
	WHERE creation_id = $1 AND version = $2
	`
	if _, err := lib.DB.Exec(q, id, curVersion, newVersion); err != nil {
		return err
	}

	return nil
}

// DeleteCreation deletes a creation
func DeleteCreation(uID uint64, creaID lib.ID) error {
	q := `
	DELETE FROM creation
	WHERE creator_id = $1
	AND id = $2
	`
	_, err := lib.DB.Exec(q, uID, creaID)
	return err
}

// SafeDeleteCreation sets creation's state to 'Deleted'
func SafeDeleteCreation(uID uint64, creaID lib.ID) error {
	q := `
	UPDATE creation
	SET state = $3
	WHERE creator_id = $1
	AND id = $2
	`
	_, err := lib.DB.Exec(q, uID, creaID, enum.Deleted)
	return err
}

// NewCreationVersion create a new version
func NewCreationVersion(crea *Creation, newVersion uint64) error {
	crea.Versions = append(crea.Versions, newVersion)
	q := `UPDATE creation SET versions = $4, state = $5 WHERE id = $2 AND creator_id = $1 AND state = $3`
	if _, err := lib.DB.Exec(q, crea.CreatorID, crea.ID, enum.Public, crea.Versions, enum.Draft); err != nil {
		return err
	}

	return CopyCreationParams(crea.ID, crea.Version, newVersion)
}

// CreationInUse returns true if the creation is used by anyone
func CreationInUse(creaID lib.ID) bool {
	var inUse struct {
		InUse bool `db:"creation_in_use"`
	}
	q := `SELECT EXISTS (SELECT creation_id FROM package_creation WHERE creation_id = $1) AS creation_in_use`
	if err := lib.DB.Get(&inUse, q, creaID); err != nil {
		return false
	}
	return inUse.InUse
}
