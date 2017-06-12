package model

import (
	"time"
	"wooble/lib"
)

// Plans label
const (
	Free    string = "Visitor"
	Woobler string = "Woobler"
)

// Plan is a Wooble plan with some restrictions
// NbPkg is the number of packages allowed (0 means infinite)
// NbCrea is the number of creations per package (0 means infinite)
type Plan struct {
	Label      *lib.NullString `json:"label,omitempty" db:"plan.label"`
	PriceMonth uint64          `json:"pricePerMonth,omitempty" db:"price_per_month"`
	PriceYear  uint64          `json:"pricePerYear,omitempty" db:"price_per_year"`

	Level *lib.NullInt64 `json:"level,omitempty" db:"level"`

	NbPkg  *lib.NullInt64 `json:"nbPkg,omitempty" db:"nb_pkg"`
	NbCrea *lib.NullInt64 `json:"nbCrea,omitempty" db:"nb_crea"`

	StartDate *lib.NullTime `json:"startDate,omitempty" db:"start_date"`
	EndDate   *lib.NullTime `json:"endDate,omitempty" db:"end_date"`
	UnsubDate *lib.NullTime `json:"unsubDate,omitempty" db:"unsub_date"`
}

// NewPlanUser logs user subscription
func NewPlanUser(uID uint64, planLabel string, periodEnd int64) (id uint64, err error) {
	q := `INSERT INTO plan_user(user_id, nb_renew, end_date, plan_label) VALUES ($1, 0, $2, $3) RETURNING id`

	if periodEnd == 0 {
		lib.DB.QueryRow(q, uID, nil, planLabel).Scan(&uID)
	} else {
		lib.DB.QueryRow(q, uID, time.Unix(periodEnd, 0), planLabel).Scan(&uID)
	}

	return id, err
}

// UnsubUserPlan unsubscribe user from his current plan (will set to free by default)
func UnsubUserPlan(planID uint64) error {
	q := `UPDATE plan_user SET unsub_date = now() WHERE id = $1`
	_, err := lib.DB.Exec(q, planID)
	return err
}

// AllPlans returns all plans
func AllPlans() (*[]Plan, error) {
	var plans []Plan
	q := `SELECT
		p.label "plan.label",
		p.price_per_month,
		p.price_per_year,
		p.nb_pkg,
		p.level,
		p.nb_crea
	FROM plan p
	ORDER BY p.label ASC
	`
	return &plans, lib.DB.Select(&plans, q)
}

// PlanByLabel gets plan with label "label"
func PlanByLabel(label string) (*Plan, error) {
	var plan Plan
	q := `SELECT
		p.label "plan.label",
		p.price_per_month,
		p.price_per_year,
		p.level,
		p.nb_pkg,
		p.nb_crea
	FROM plan p
	WHERE p.label = $1
	`
	return &plan, lib.DB.Get(&plan, q, label)
}

// DefaultPlan gets free plan or second plan if user is VIP
func DefaultPlan(userID uint64) (*Plan, error) {
	var plan Plan
	q := `
		SELECT
      pl.label "plan.label",
      pl.nb_pkg,
			pl.level,
      pl.nb_crea
		FROM plan pl, app_user u
		WHERE CASE WHEN u.is_vip THEN pl.label = $2 ELSE pl.label = $3 END
		AND u.id = $1
	`

	return &plan, lib.DB.Get(&plan, q, userID, Woobler, Free)
}

// HasExpired returns true if the plan has expired
func (plan *Plan) HasExpired() bool {
	if plan.EndDate == nil {
		return true
	}
	return plan.EndDate.Time.Unix() < time.Now().Unix()
}
