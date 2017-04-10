package model

import (
	"fmt"
	"wooble/lib"
	enum "wooble/models/enums"
)

// Creation is a Wooble creation
type Creation struct {
	ID lib.ID `json:"id"      db:"crea.id"`

	Title       string          `json:"title"  db:"title"`
	Description *lib.NullString `json:"description,omitempty" db:"description"`
	Creator     User            `json:"creator,omitempty" db:""`
	Versions    lib.UintSlice   `json:"versions,omitempty" db:"versions"`
	Alias       string          `json:"alias,omitempty" db:"alias"`
	State       string          `json:"state,omitempty" db:"state"`
	IsOwner     bool            `json:"isOwner,omitempty" db:"is_owner"`

	PreviewURL string `json:"previewUrl,omitempty"`
	Version    uint64 `json:"version,omitempty"`

	Script   string `json:"script,omitempty"`
	Document string `json:"document,omitempty"`
	Style    string `json:"style,omitempty"`

	CreatorID uint64 `json:"-"       db:"creator_id"`
	Engine    Engine `json:"-" db:""`
	Price     uint64 `json:"price" db:"price"` // in cents euro
	IsToBuy   bool   `json:"isToBuy,omitempty" db:"is_to_buy"`

	CreatedAt *lib.NullTime `json:"createdAt,omitempty" db:"crea.created_at"`
	UpdatedAt *lib.NullTime `json:"updatedAt,omitempty" db:"crea.updated_at"`
}

// CreationPurchase is an history of all creations purchase
type CreationPurchase struct {
	UserID uint64 `db:"crea_pur.user_id"`
	CreaID uint64 `db:"crea_pur.creation_id"`

	Price    uint64 `db:"crea_pur.price"`
	ChargeID string `db:"charge_id"`

	PurchasedAt *lib.NullTime `db:"purchased_at"`
}

// BaseVersion is creation default version
const BaseVersion uint64 = 1

// AllCreations returns all creations
func AllCreations(opt lib.Option, uID uint64) (*[]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
	    c.id "crea.id",
	    c.title,
			c.description,
	    c.created_at "crea.created_at",
	    c.updated_at "crea.updated_at",
	    c.versions,
			c.price,
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
		WHERE (c.state = 'public'
		OR (
			c.state = 'draft' AND array_length(c.versions, 1) > 1
		))
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(lib.SEARCH, "c.title|u.name", lib.CREATOR, "u.name")
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
	    c.created_at "crea.created_at",
	    c.updated_at "crea.updated_at",
	    c.versions,
			c.price,
			c.alias,
			CASE WHEN c.creator_id = $1 THEN true ELSE false END "is_owner",
			c.state,
			COUNT(c.id) AS nb_crea,
	    u.id "user.id",
	    u.name
	  FROM creation c
	  INNER JOIN app_user u ON (c.creator_id = u.id)
		LEFT JOIN package_creation pc ON (pc.creation_id = c.id)
		WHERE (c.state = 'public'
		OR (
			c.state = 'draft' AND array_length(c.versions, 1) > 1
		))
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(lib.SEARCH, "c.title|u.name", lib.CREATOR, "u.name")

	q.Q += "GROUP BY c.id, u.id ORDER BY nb_crea DESC"

	return &creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// AllUsedCreations return creations used in some packages
func AllUsedCreations(opt lib.Option, uID uint64) (*[]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
			DISTINCT c.id "crea.id",
			c.title,
			c.created_at "crea.created_at",
			c.versions,
			c.price,
			u.id "user.id",
			u.name
		FROM creation c
		INNER JOIN app_user u ON (c.creator_id = u.id)
		INNER JOIN package_creation pc ON (pc.creation_id = c.id)
    INNER JOIN package p ON (p.id = pc.package_id)
		WHERE p.user_id = $1
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(lib.SEARCH, "c.title")

	return &creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// AllPurchasedCreations returns all purchased creation by the user 'uID'
func AllPurchasedCreations(opt lib.Option, uID uint64) (*[]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
			DISTINCT c.id "crea.id",
			c.title,
			c.created_at "crea.created_at",
			c.versions,
			c.price,
			u.id "user.id",
			u.name
		FROM creation c
		INNER JOIN app_user u ON (c.creator_id = u.id)
		INNER JOIN creation_purchase cp ON (cp.creation_id = c.id)
		WHERE cp.user_id = $1
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(lib.SEARCH, "c.title")

	return &creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// AllDraftCreations returns all creation in draft of authenticated user
func AllDraftCreations(opt lib.Option, uID uint64) (*[]Creation, error) {
	var creations []Creation
	q := lib.NewQuery(`SELECT
			c.id "crea.id",
			c.title,
			c.created_at "crea.created_at",
			c.versions,
			c.price,
			u.id "user.id",
			u.name
		FROM creation c
		INNER JOIN app_user u ON (c.creator_id = u.id)
		WHERE u.id = $1 AND c.state = 'draft'
		`, &opt)

	q.AddValues(uID)
	q.SetFilters(lib.SEARCH, "c.title")

	return &creations, lib.DB.Select(&creations, q.String(), q.Values...)
}

// CreationByID returns a creation with the id "id"
func CreationByID(id lib.ID, uID uint64) (*Creation, error) {
	var crea Creation
	q := `
  SELECT
    c.id "crea.id",
    c.title,
		c.description,
    c.created_at "crea.created_at",
    c.updated_at "crea.updated_at",
    c.versions,
		c.price,
		c.alias,
		CASE WHEN c.creator_id = $2 THEN true ELSE false END "is_owner",
		c.state,
		e.name "eng.name",
		e.extension,
		e.content_type,
    u.id "user.id",
    u.name
  FROM creation c
  INNER JOIN app_user u ON (c.creator_id = u.id)
	INNER JOIN engine e ON (c.engine=e.name)
  WHERE c.id = $1 AND (c.state = 'public' OR c.state = 'delete' OR array_length(c.versions, 1) > 1)
	`

	if err := lib.DB.Get(&crea, q, id, uID); err != nil {
		return nil, err
	}

	if len(crea.Versions) > 1 {
		crea.Versions = crea.Versions[:len(crea.Versions)-1]
	}

	return &crea, nil
}

// CreationPrivateByID returns a creation as private
func CreationPrivateByID(uID uint64, creaID lib.ID) (*Creation, error) {
	var crea Creation
	q := `
	SELECT
		c.id "crea.id",
		c.title,
		c.description,
		c.alias,
		c.creator_id,
		c.created_at "crea.created_at",
		c.updated_at "crea.updated_at",
		c.versions,
		c.price,
		c.state,
		e.name "eng.name",
		e.extension,
		e.content_type
	FROM creation c
	INNER JOIN engine e ON (c.engine=e.name)
	WHERE creator_id = $1 AND c.id = $2 AND c.state != 'delete'
	`

	return &crea, lib.DB.Get(&crea, q, uID, creaID)
}

// UpdateCreation update creation's information
func UpdateCreation(crea *Creation) error {
	q := `
  UPDATE creation
  SET title = $3, description = $4, price = $5, state = $6, alias = $7
  WHERE id = $1
  AND creator_id = $2
  `
	_, err := lib.DB.Exec(q, crea.ID, crea.CreatorID, crea.Title, crea.Description, crea.Price, crea.State, crea.Alias)
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
func NewCreationVersion(uID uint64, creaID lib.ID, version lib.UintSlice) error {
	q := `UPDATE creation SET versions = $4, state = $5 WHERE id = $2 AND creator_id = $1 AND state = $3`
	_, err := lib.DB.Exec(q, uID, creaID, enum.Public, version, enum.Draft)
	return err
}

// NewCreationPurchases creates a creation purchase
func NewCreationPurchases(buyerID uint64, chargeID string, creations *[]Creation) error {
	qPurchase := `INSERT INTO creation_purchase(user_id, creation_id,	price, charge_id) VALUES ($1, $2, $3, $4)`

	qSellerDue := `UPDATE app_user SET fund=fund + $2 WHERE id = $1`

	tx := lib.DB.MustBegin()
	for _, crea := range *creations {
		tx.Exec(qPurchase, buyerID, crea.ID, crea.Price, chargeID)

		// Credit the seller
		tx.Exec(qSellerDue, crea.Creator.ID, crea.Price)
	}

	return tx.Commit()
}

// CreationInUse returns true if the creation is used by anyone
func CreationInUse(creaID lib.ID) bool {
	var inUse struct {
		InUse bool `db:"creation_in_use"`
	}
	q := `
	SELECT EXISTS (SELECT creation_id FROM creation_purchase WHERE creation_id = $1)
	OR EXISTS (SELECT creation_id FROM package_creation WHERE creation_id = $1) AS creation_in_use
	`
	if err := lib.DB.Get(&inUse, q, creaID); err != nil {
		return false
	}
	return inUse.InUse
}
