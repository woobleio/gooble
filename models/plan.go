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
	Label      string  `json:"label" db:"plan.label"`
	PriceMonth float64 `json:"pricePerMonth" db:"price_per_month"`
	PriceYear  float64 `json:"pricePerYear" db:"price_per_year"`

	NbPkg     int `json:"nbPackages" db:"nb_pkg"`
	NbCrea    int `json:"nbCreations" db:"nb_crea"`
	NbDomains int `json:"nbDomains" db:"nb_domains"`
}

// NewPlanUser logs user subscription
func NewPlanUser(uID uint64, planLabel string, periodEnd int64) (id uint64, err error) {
	q := `INSERT INTO plan_user(user_id, nb_renew, end_date, plan_label) VALUES ($1, $2, $3, $4) RETURNING id`
	err = lib.DB.QueryRow(q, uID, 0, time.Unix(periodEnd, 0), planLabel).Scan(&uID)
	return id, err
}
