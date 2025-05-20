package data

import (
	"context"
	"database/sql"
	"time"
)

type Subscription struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	VisitsPerWeek uint8   `json:"visits_per_week"`
	Price         float64 `json:"price"`
}

type SubscriptionModel struct {
	DB *sql.DB
}

func (sm SubscriptionModel) GetAll() ([]*Subscription, error) {
	query := `SELECT id, name, visits_per_week, price FROM subscriptions`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := sm.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	subscriptions := []*Subscription{}

	for rows.Next() {
		var subscription Subscription

		err := rows.Scan(&subscription.ID, &subscription.Name, &subscription.VisitsPerWeek, &subscription.Price)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, &subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}
