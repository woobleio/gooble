package model

import (
	"wooble/lib"
	enum "wooble/models/enums"
)

// Creation is a Wooble creation
type Creation struct {
	ID lib.ID `json:"id"      db:"crea.id"`

	Title       string          `json:"title"  db:"title"`
	Description *lib.NullString `json:"description,omitempty" db:"description"`
	Creator     User            `json:"creator" db:""`
	Versions    lib.StringSlice `json:"versions" db:"versions"`
	Version     string          `json:"version,omitempty" db:"version"`
	Alias       *lib.NullString `json:"alias,omitempty" db:"alias"`
	State       string          `json:"state" db:"state"`

	CreatorID uint64 `json:"-"       db:"creator_id"`
	HasDoc    bool   `json:"-"       db:"has_document"`
	HasStyle  bool   `json:"-"       db:"has_style"`
	Engine    Engine `json:"-" db:""`
	Price     uint64 `json:"price" db:"price"` // in cents euro
	IsToBuy   bool   `json:"isToBuy" db:"is_to_buy"`

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
const BaseVersion string = "1.0"

// AllCreations returns all creations
func AllCreations(opt lib.Option) (*[]Creation, error) {
	var creations []Creation
	q := lib.Query{
		Q: `SELECT
	    c.id "crea.id",
	    c.title,
			c.description,
	    c.created_at "crea.created_at",
	    c.updated_at "crea.updated_at",
	    c.versions,
			c.has_document,
			c.has_style,
			c.price,
			c.state,
			e.name "eng.name",
			e.extension,
			e.content_type,
	    u.id "user.id",
	    u.name
	  FROM creation c
	  INNER JOIN app_user u ON (c.creator_id = u.id)
		INNER JOIN engine e ON (c.engine=e.name)
		WHERE c.state = 'public' OR c.state = 'delete'
		`,
		Opt: &opt,
	}

	query := q.String()

	return &creations, lib.DB.Select(&creations, query)
}

// CreationByID returns a creation with the id "id"
func CreationByID(id lib.ID) (*Creation, error) {
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
		c.state,
		c.has_document,
		c.has_style,
		e.name "eng.name",
		e.extension,
		e.content_type,
    u.id "user.id",
    u.name
  FROM creation c
  INNER JOIN app_user u ON (c.creator_id = u.id)
	INNER JOIN engine e ON (c.engine=e.name)
  WHERE c.id = $1 AND (c.state = 'public' OR c.state = 'delete')
	`

	return &crea, lib.DB.Get(&crea, q, id)
}

// UpdateCreation update creation's information
func UpdateCreation(crea *Creation) error {
	q := `
  UPDATE creation
  SET title = $3, description = $4, price = $5
  WHERE id = $1
  AND creator_id = $2 
  `
	_, err := lib.DB.Exec(q, crea.ID, crea.CreatorID, crea.Title, crea.Description, crea.Price)
	return err
}

// CreationByIDAndVersion returns the creation "creaID" and check if the version "version" exists
func CreationByIDAndVersion(id lib.ID, version string) (*Creation, error) {
	crea := new(Creation)
	if version == "" {
		version = BaseVersion
	}
	q := `
  SELECT id "crea.id", versions, state 
  FROM creation WHERE id = $1 
  AND $2 = ANY (versions)
  `
	return crea, lib.DB.Get(crea, q, id, version)
}

// NewCreation creates a creation
func NewCreation(crea *Creation) (*Creation, error) {
	q := `
  INSERT INTO creation(
    title, 
    description, 
    creator_id, 
    versions, 
    price, 
    engine,
		state
  ) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id
  `

	stringSliceVersions := make(lib.StringSlice, 0, 1)

	return crea, lib.DB.QueryRow(q, crea.Title, crea.Description, crea.CreatorID, append(stringSliceVersions, BaseVersion), crea.Price, crea.Engine.Name, crea.State).Scan(&crea.ID)
}

// NewCreationVersion create a new version
func NewCreationVersion(crea *Creation) (int64, error) {
	q := `UPDATE creation SET versions=$3 WHERE id = $1 AND creator_id = $2 AND state = $3`
	res, err := lib.DB.Exec(q, crea.ID, crea.CreatorID, crea.Versions, enum.Draft)
	rowAff, _ := res.RowsAffected()
	return rowAff, err
}

// NewCreationPurchases creates a creation purchase
func NewCreationPurchases(buyerID uint64, chargeID string, creations *[]Creation) error {
	qPurchase := `INSERT INTO creation_purchase(user_id, creation_id,	price, charge_id) VALUES ($1, $2, $3, $4)`

	qSellerDue := `UPDATE app_user SET total_due=total_due + $2 WHERE id = $1`

	tx := lib.DB.MustBegin()
	for _, crea := range *creations {
		tx.Exec(qPurchase, buyerID, crea.ID, crea.Price, chargeID)

		// Credit the seller
		tx.Exec(qSellerDue, crea.Creator.ID, crea.Price)
	}

	return tx.Commit()
}

// UpdateCreationCode updates creation information
func UpdateCreationCode(crea *Creation) (int64, error) {
	q := `
  UPDATE creation SET has_document = $2, has_style = $3
  WHERE id = $1
  AND state = 'draft'
  AND versions[array_length(versions, 1)] = $5
  `
	res, err := lib.DB.Exec(q, crea.ID, crea.HasDoc, crea.HasStyle, crea.Version)
	rowAff, _ := res.RowsAffected()
	return rowAff, err
}

// PublishCreation switches creation "creaID" state to "public"
func PublishCreation(crea *Creation) error {
	q := `
  UPDATE creation SET state = 'public'
  WHERE creator_id = $1
  AND state = $3
  AND id = $2
  `
	_, err := lib.DB.Exec(q, crea.CreatorID, crea.ID, enum.Draft)
	return err
}
