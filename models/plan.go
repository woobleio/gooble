package model

import (
	"time"
	"wooble/lib"
)

const (
	Free string = "free"
)

// Plan is a Wooble plan with some restrictions
// NbPkg is the number of packages allowed (0 means infinite)
// NbCrea is the number of creations per package (0 means infinite)
// NbDomains is the number of domains per package (0 means infinite)
type Plan struct {
	Label      lib.NullString `json:"label" db:"plan.label"`
	PriceMonth float64        `json:"pricePerMonth" db:"price_per_month"`
	PriceYear  float64        `json:"pricePerYear" db:"price_per_year"`

	NbPkg     lib.NullInt64 `json:"nbPkg" db:"nb_pkg"`
	NbCrea    lib.NullInt64 `json:"nbCrea" db:"nb_crea"`
	NbDomains lib.NullInt64 `json:"nbDomains" db:"nb_domains"`

	StartDate *lib.NullTime `json:"startDate,omitempty" db:"start_date"`
	EndDate   *lib.NullTime `json:"endDate,omitempty" db:"end_date"`
}

// NewPlanUser logs user subscription
func NewPlanUser(uID uint64, planLabel string, periodEnd int64) (id uint64, err error) {
	q := `INSERT INTO plan_user(user_id, nb_renew, end_date, plan_label) VALUES ($1, $2, $3, $4) RETURNING id`
	err = lib.DB.QueryRow(q, uID, 0, time.Unix(periodEnd, 0), planLabel).Scan(&uID)
	return id, err
}

// DefaultPlan get free plan from the db
func DefaultPlan() (*Plan, error) {
	var plan Plan
	q := `
		SELECT
      pl.label "plan.label",
      pl.nb_pkg,
      pl.nb_crea,
      pl.nb_domains
		FROM plan pl
		WHERE pl.label = $1
	`

	return &plan, lib.DB.Get(&plan, q, Free)
}
