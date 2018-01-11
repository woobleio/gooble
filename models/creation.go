package model

import (
	"bytes"
	"fmt"
	"image/png"
	"strings"
	"wooble/lib"
	enum "wooble/models/enums"

	"github.com/nfnt/resize"
)

// Creation is a Wooble object
type Creation struct {
	ID lib.ID `json:"id"      db:"crea.id"`

	Title       string          `json:"title"  db:"title"`
	Description *lib.NullString `json:"description,omitempty" db:"description"`
	Tags        []Tag           `json:"tags" db:""`
	ThumbPath   *lib.NullString `json:"thumbPath,omitempty" db:"thumb_path"`
	Creator     User            `json:"creator,omitempty" db:""`

	// When a creation is in draft, the very last version is ignored by most queries
	Versions lib.UintSlice `json:"versions,omitempty" db:"versions"`

	Alias      string `json:"alias,omitempty" db:"alias"`
	State      string `json:"state,omitempty" db:"state"`
	IsOwner    bool   `json:"isOwner,omitempty" db:"is_owner"`
	IsFeatured bool   `json:"-" db:"is_featured"`

	NbUse uint64 `json:"nbUse" db:"nb_use"`

	Params     []CreationParam    `json:"params,omitempty" db:""`
	Functions  []CreationFunction `json:"functions,omitempty" db:""`
	PreviewURL string             `json:"previewUrl,omitempty"`
	Version    uint64             `json:"version,omitempty"`

	PreviewPos       PreviewPosition   `json:"previewPosition" db:""`
	PreviewPositions []PreviewPosition `json:"previewPositions,omitempty" db:""`
	IsThumbPreview   bool              `json:"isThumbPreview" db:"is_thumb_preview"`

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
	CreationID lib.ID `json:"-" db:"creation_id"`
	Field      string `json:"field" db:"field"`
	Value      string `json:"value" db:"value"`
}

// CreationFunction is a creation function
type CreationFunction struct {
	Call   string `json:"call" db:"call"`
	Detail string `json:"detail" db:"detail"`
}

// PreviewPosition is the creation position in the preview
type PreviewPosition struct {
	Position    string `json:"position" db:"position_id"`
	StyleSource string `json:"styleSource" db:"style_source"`
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
	q := `SELECT field, value FROM creation_param WHERE creation_id = $1 AND version ` + c.getLastVersionQuery()
	lib.DB.Select(&c.Params, q, c.ID, c.Version)
}

// PopulateFunctions populates creation's functions for documentation
func (c *Creation) PopulateFunctions() {
	q := `SELECT call, detail FROM creation_function WHERE creation_id = $1 AND version ` + c.getLastVersionQuery()
	lib.DB.Select(&c.Functions, q, c.ID, c.Version)
}

// PopulateTags populates creation's tags
func (c *Creation) PopulateTags() error {
	q := `
	SELECT
		tag.id "tag.id",
		tag.title "tag.title"
	FROM creation_tag ct
	INNER JOIN creation ON (creation.id=ct.creation_id)
	INNER JOIN tag ON (tag.id=ct.tag_id)
	WHERE ct.creation_id = $1
	`

	return lib.DB.Select(&c.Tags, q, c.ID)
}

// PopulatePreviewPositions populates available creation's preview positions
func (c *Creation) PopulatePreviewPositions() {
	q := `SELECT position_id, style_source FROM preview_position`
	lib.DB.Select(&c.PreviewPositions, q)
}

// RetrievePreviewURL construct the creation's preview URL and sets the `PreviewUrl` field
// It uses the current `Version` of the creation c
func (c *Creation) RetrievePreviewURL() {
	s := lib.NewStorage(lib.SrcPreview)

	if c.Version == 0 {
		c.Version = CreationLastVersion(c.ID)
	}

	creaLastVersion := fmt.Sprintf("%d", c.Version)

	creatorID := fmt.Sprintf("%d", c.Creator.ID)
	creaID := fmt.Sprintf("%d", c.ID.ValueDecoded)

	previewURL := s.GetPathFor(creatorID, creaID, creaLastVersion, "index.html")

	spltPath := strings.Split(previewURL, "/")

	c.PreviewURL = strings.Join(spltPath[1:], "/")
}

// RetrieveSourceCode request the source files in the cloud and set the content to the Creation
func (c *Creation) RetrieveSourceCode(version string, files ...string) error {
	uIDStr := fmt.Sprintf("%d", c.Creator.ID)
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
func AllCreations(opt lib.Option, uID uint64) ([]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
	    c.id "crea.id",
			COUNT(pc.creation_id) AS nb_use,
	    c.title,
			c.description,
			c.thumb_path,
			c.is_thumb_preview,
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
		LEFT JOIN creation_tag ct ON (ct.creation_id = c.id)
		LEFT JOIN tag t ON (t.id = ct.tag_id)
		LEFT JOIN package_creation pc ON (pc.creation_id = c.id)
		WHERE (c.state = 'public' OR array_length(versions, 1) > 1) AND c.state != 'delete'
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(true, lib.SEARCH, "c.title|u.name|t.title", lib.CREATOR, "u.name")

	q.Q += " GROUP BY c.id, u.id, e.name"

	q.SetOrder(lib.CREATED_AT, "c.created_at")

	return creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// AllPopularCreations returns all popular creations
func AllPopularCreations(opt lib.Option, uID uint64) ([]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
	    c.id "crea.id",
	    c.title,
			c.description,
			c.thumb_path,
			c.is_thumb_preview,
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
		LEFT JOIN creation_tag ct ON (ct.creation_id = c.id)
		LEFT JOIN tag t ON (t.id = ct.tag_id)
		LEFT JOIN package_creation pc ON (pc.creation_id = c.id)
		WHERE (c.state = 'public' OR array_length(versions, 1) > 1) AND c.state != 'delete'
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(true, lib.SEARCH, "c.title|u.name|t.title", lib.CREATOR, "u.name")

	q.Q += " GROUP BY c.id, u.id ORDER BY c.is_featured DESC, nb_use DESC"

	return creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// AllUsedCreations return creations used in some packages
func AllUsedCreations(opt lib.Option, uID uint64) ([]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
			DISTINCT c.id "crea.id",
			COUNT(pcc.creation_id) AS nb_use,
			c.title,
			c.thumb_path,
			c.is_thumb_preview,
			c.created_at "crea.created_at",
			CASE WHEN c.state = 'public' THEN c.versions[0:array_length(c.versions, 1)] ELSE c.versions[0:array_length(c.versions, 1)-1] END AS versions,
			c.thumb_path,
			u.id "user.id",
			u.name
		FROM creation c
		INNER JOIN app_user u ON (c.creator_id = u.id)
		INNER JOIN package_creation pc ON (pc.creation_id = c.id)
		LEFT JOIN creation_tag ct ON (ct.creation_id = c.id)
		LEFT JOIN tag t ON (t.id = ct.tag_id)
		LEFT JOIN package_creation pcc ON (pcc.creation_id = c.id)
    INNER JOIN package p ON (p.id = pc.package_id)
		WHERE p.user_id = $1
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(true, lib.SEARCH, "c.title|t.title")

	q.Q += " GROUP BY c.id, u.id"

	return creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// AllDraftCreations returns all creation in draft of authenticated user
func AllDraftCreations(opt lib.Option, uID uint64) ([]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
			c.id "crea.id",
			COUNT(pc.creation_id) AS nb_use,
			c.title,
			c.thumb_path,
			c.is_thumb_preview,
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
	q.SetFilters(true, lib.SEARCH, "c.title")

	q.Q += " GROUP BY c.id, u.id"

	return creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// CreationByID returns a creation with the id "id"
func CreationByID(id lib.ID, uID uint64, latestVersion bool) (*Creation, error) {
	var crea Creation
	q := `
  SELECT
    c.id "crea.id",
    c.title,
		c.thumb_path,
		c.is_thumb_preview,
		c.description,
    c.created_at "crea.created_at",
    c.updated_at "crea.updated_at",
		CASE WHEN c.state = 'public' THEN versions[0:array_length(c.versions, 1)] ELSE versions[0:array_length(c.versions, 1)-$3] END AS versions,
		c.alias,
		CASE WHEN c.creator_id = $2 THEN true ELSE false END "is_owner",
		c.state,
		e.name "eng.name",
		e.extension,
		e.content_type,
		pp.position_id,
		pp.style_source,
    u.id "user.id",
    u.name,
		u.pic_path
  FROM creation c
  INNER JOIN app_user u ON (c.creator_id = u.id)
	INNER JOIN engine e ON (c.engine=e.name)
	INNER JOIN preview_position pp ON (c.preview_position_id=pp.position_id)
  WHERE c.id = $1
	`

	decVersion := 1
	if latestVersion {
		decVersion = 0
	}

	if err := lib.DB.Get(&crea, q, id, uID, decVersion); err != nil {
		return nil, err
	}

	nbVersions := len(crea.Versions)
	if nbVersions > 0 {
		crea.Version = crea.Versions[nbVersions-1] // Latest version
	} else {
		crea.Version = 1 // If no versions, the default is one
	}

	crea.PopulateParams()
	crea.PopulateFunctions()
	crea.PopulatePreviewPositions()
	crea.PopulateTags()

	return &crea, nil
}

// UpdateCreation update creation's information
func UpdateCreation(crea *Creation) error {
	q := `
  UPDATE creation
  SET title = $3, description = $4, state = $5, alias = $6, thumb_path = $7, is_thumb_preview = $8, preview_position_id = $9
  WHERE id = $1
  AND creator_id = $2
  `

	if _, err := lib.DB.Exec(q, crea.ID, crea.Creator.ID, crea.Title, crea.Description, crea.State, crea.Alias, crea.ThumbPath, crea.IsThumbPreview, crea.PreviewPos.Position); err != nil {
		return err
	}

	if len(crea.Params) > 0 {
		if err := UpdateCreationParams(crea); err != nil {
			return err
		}
	}

	if len(crea.Functions) > 0 {
		if err := UpdateCreationFunctions(crea); err != nil {
			return err
		}
	}

	return UpdateCreationTags(crea)
}

// UpdateCreationPatch updates a creation
func UpdateCreationPatch(uID uint64, creaID lib.ID, patch lib.SQLPatch) error {
	updateQuery := patch.GetUpdateQuery("creation")
	if len(updateQuery) == 0 {
		return nil
	}
	q := updateQuery +
		` WHERE creator_id = $` + fmt.Sprintf("%d", patch.Index+1) +
		` AND id = $` + fmt.Sprintf("%d", patch.Index+2)

	patch.Args = append(patch.Args, uID)
	patch.Args = append(patch.Args, creaID)

	_, err := lib.DB.Exec(q, patch.Args...)

	return err
}

// UpdateCreationTags updates creation tags
func UpdateCreationTags(crea *Creation) (err error) {
	// Reset all tags
	lib.DB.Exec(`DELETE FROM creation_tag WHERE creation_id = $1`, crea.ID)

	if len(crea.Tags) > 0 {
		bulk := lib.NewQuery(`INSERT INTO creation_tag(creation_id, tag_id) VALUES`, nil)

		values := make([]interface{}, 0)
		for _, t := range crea.Tags {
			values = append(values, t)
		}

		bulk.SetBulkInsert([]interface{}{crea.ID}, []string{"ID"}, values...)

		_, err = lib.DB.Exec(bulk.String(), bulk.Values...)
	}

	return err
}

// UpdateCreationFunctions updates creation functions
func UpdateCreationFunctions(crea *Creation) error {
	lastVersion := crea.Version
	if lastVersion == 0 {
		if err := lib.DB.Get(&lastVersion, lastVersionQuery, crea.ID); err != nil {
			return err
		}
	}

	q := `
	DELETE FROM creation_function
	WHERE creation_id = $1
	AND version = $2
	`
	lib.DB.Exec(q, crea.ID, lastVersion)

	bulk := lib.NewQuery(`INSERT INTO creation_function(creation_id, version, call, detail) VALUES`, nil)

	values := make([]interface{}, 0)
	for _, fn := range crea.Functions {
		values = append(values, fn)
	}

	bulk.SetBulkInsert([]interface{}{crea.ID, lastVersion}, []string{"Call", "Detail"}, values...)

	_, err := lib.DB.Exec(bulk.String(), bulk.Values...)
	return err
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

	bulk := lib.NewQuery(`INSERT INTO creation_param(creation_id, version, field, value) VALUES`, nil)

	values := make([]interface{}, 0)
	for _, p := range crea.Params {
		values = append(values, p)
	}

	bulk.SetBulkInsert([]interface{}{crea.ID, lastVersion}, []string{"Field", "Value"}, values...)

	_, err := lib.DB.Exec(bulk.String(), bulk.Values...)
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
	if err := lib.DB.Get(crea, q, id, version); err != nil {
		return nil, err
	}

	crea.PopulateParams()

	return crea, nil
}

// CreationLastVersion gets creation's last version
func CreationLastVersion(id lib.ID) uint64 {
	var version uint64
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

	if err := lib.DB.QueryRow(q, crea.Title, crea.CreatorID, append(stringSliceVersions, fmt.Sprintf("%d", BaseVersion)), crea.Engine.Name, crea.State, crea.Alias).Scan(&crea.ID); err != nil {
		return nil, err
	}

	_, path := crea.GenerateDefaultThumb()

	q = `UPDATE creation SET thumb_path = $2 WHERE id = $1`
	lib.DB.Exec(q, crea.ID, path)

	return crea, nil
}

// CopyCreationParamsAndFunctions copie creation parameters and functions from current version to a new version
func CopyCreationParamsAndFunctions(id lib.ID, curVersion uint64, newVersion uint64) error {
	q := `
	INSERT INTO creation_param(creation_id, version, field, value)
	SELECT creation_id, $3 AS version, field, value
	FROM creation_param
	WHERE creation_id = $1 AND version = $2
	`
	if _, err := lib.DB.Exec(q, id, curVersion, newVersion); err != nil {
		return err
	}

	q = `
	INSERT INTO creation_function(creation_id, version, call, detail)
	SELECT creation_id, $3 AS version, call, detail
	FROM creation_function
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
	if _, err := lib.DB.Exec(q, crea.Creator.ID, crea.ID, enum.Public, crea.Versions, enum.Draft); err != nil {
		return err
	}

	return CopyCreationParamsAndFunctions(crea.ID, crea.Version, newVersion)
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

// GenerateDefaultThumb generates unique thumbnail for the creation
// it returns the image buffer and the image path
func (c *Creation) GenerateDefaultThumb() ([]byte, string) {
	buff := new(bytes.Buffer)

	img := resize.Resize(100, 0, lib.GenImage(c.ID.ValueDecoded), resize.Lanczos3)

	png.Encode(buff, img)

	pngToBytes := buff.Bytes()

	source := lib.SrcCreaThumb
	storage := lib.NewStorage(source)
	path := storage.StoreFile(pngToBytes, "image/png", fmt.Sprintf("%d", c.CreatorID), source+c.ID.ValueEncoded, "", "gen.png")

	splPath := strings.Split(path, "/")
	path = strings.Join(splPath[1:len(splPath)], "/")

	return pngToBytes, path
}

func (c *Creation) getLastVersionQuery() string {
	if c.Version == 0 {
		return `(` + lastVersionQuery + `)`
	}
	return `= $2`
}
