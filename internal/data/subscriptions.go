package data

import "database/sql"

type Subscription struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	VisitsPerWeek uint8   `json:"visits_per_week"`
	Price         float64 `json:"price"`
}

type SubscriptionModel struct {
	DB *sql.DB
}
