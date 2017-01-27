package model

// Plan is a Wooble plan with some restrictions
// NbPkg is the number of packages allowed (0 means infinite)
// NbCrea is the number of creations per package (0 means infinite)
// NbDomains is the number of domains per package (0 means infinite)
type Plan struct {
	ID uint64 `json:"id" db:"plan.id"`

	Label      string  `json:"label" db:"plan.label"`
	PriceMonth float64 `json:"pricePerMonth" db:"price_per_month"`
	PriceYear  float64 `json:"pricePerYear" db:"price_per_year"`

	NbPkg     int `json:"nbPackages" db:"nb_pkg"`
	NbCrea    int `json:"nbCreations" db:"nb_crea"`
	NbDomains int `json:"nbDomains" db:"nb_domains"`
}
