package model

import "wooble/lib"

// Creation is a Wooble creation
type Creation struct {
	ID lib.ID `json:"id"      db:"crea.id"`

	Title       string          `json:"title"  db:"title"`
	Description *lib.NullString `json:"description,omitempty" db:"description"`
	Creator     User            `json:"creator" db:""`
	Version     string          `json:"version" db:"version"`
	Alias       *lib.NullString `json:"alias,omitempty" db:"alias"`

	CreatorID uint64 `json:"-"       db:"creator_id"`
	HasDoc    bool   `json:"-"       db:"has_document"`
	HasScript bool   `json:"-"       db:"has_script"`
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

	Total    uint64 `db:"crea_pur.total"`
	ChargeID string `db:"charge_id"`

	PurchasedAt *lib.NullTime `db:"purchased_at"`
}

// CreationForm is a form for creation
type CreationForm struct {
	CreatorID uint64

	Engine string `json:"engine" binding:"required"`
	Title  string `json:"title" binding:"required"`

	Description string `json:"description"`
	Document    string `json:"document"`
	Script      string `json:"script"`
	Style       string `json:"style"`

	Price uint64 `json:"price,omitempty"`

	Version string
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
	    c.version,
			c.has_document,
			c.has_script,
			c.has_style,
			c.price,
			e.name "eng.name",
			e.extension,
			e.content_type,
	    u.id "user.id",
	    u.name
	  FROM creation c
	  INNER JOIN app_user u ON (c.creator_id = u.id)
		INNER JOIN engine e ON (c.engine=e.name)
		`,
		Opt: &opt,
	}

	query := q.String()

	return &creations, lib.DB.Select(&creations, query)
}

// CreationByID returns a creation with the id "id"
func CreationByID(id string) (*Creation, error) {
	var crea Creation
	q := `
  SELECT
    c.id "crea.id",
    c.title,
		c.description,
    c.created_at "crea.created_at",
    c.updated_at "crea.updated_at",
    c.version,
		c.price,
		c.has_document,
		c.has_script,
		c.has_style,
		e.name "eng.name",
		e.extension,
		e.content_type,
    u.id "user.id",
    u.name
  FROM creation c
  INNER JOIN app_user u ON (c.creator_id = u.id)
	INNER JOIN engine e ON (c.engine=e.name)
  WHERE c.id = $1
	`

	encodedID, _ := lib.DecodeHash(id)

	return &crea, lib.DB.Get(&crea, q, encodedID)
}

// DeleteCreation deletes creation id "id"
func DeleteCreation(id string) error {
	encodedID, _ := lib.DecodeHash(id)
	_, err := lib.DB.Exec(`DELETE FROM creation WHERE id = $1`, encodedID)
	return err
}

// NewCreation creates a creation
func NewCreation(data *CreationForm) (string, error) {
	var creaID int64
	q := `
  INSERT INTO creation(
    title, 
    description, 
    creator_id, 
    version, 
    price, 
    has_document, 
    has_script, 
    has_style, 
    engine
  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id
  `

	err := lib.DB.QueryRow(q, data.Title, data.Description, data.CreatorID, data.Version, data.Price, data.Document != "", data.Script != "", data.Style != "", data.Engine).Scan(&creaID)
	if err != nil {
		return "", err
	}
	return lib.HashID(creaID)
}

// NewCreationPurchase creates a creation purchase
func NewCreationPurchase(data *CreationPurchase) error {
	q := `INSERT INTO creation_purchase(
		user_id, 
		creation_id,
		total,
		charge_id
	) VALUES ($1, $2, $3, $4)
	`
	_, err := lib.DB.Exec(q, data.UserID, data.CreaID, data.Total, data.ChargeID)
	return err
}
