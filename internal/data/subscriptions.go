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

type Subscriptions struct {
	UserID           int64     `json:"user_id"`
	FullName         string    `json:"full_name"`
	SubscriptionID   int64     `json:"sub_id"`
	SubscriptionName string    `json:"sub_name"`
	VisitsPerWeek    uint8     `json:"visits_per_week"`
	Price            float64   `json:"price"`
	DateStart        time.Time `json:"date_start"`
	DateEnd          time.Time `json:"date_end"`
}

func (sm SubscriptionModel) UserSubscriptions(id int64) ([]*Subscriptions, error) {
	query := `SELECT 
    u.id AS user_id,
    u.full_name AS user_name,
    sub.id AS subscription_id,
    sub.name AS subscription_name,
    sub.visits_per_week,
    sub.price,
    us.date_start,
    us.date_end
FROM user_subscriptions us
JOIN users u ON us.user_id = u.id
JOIN subscriptions sub ON us.subscription_id = sub.id
WHERE u.id = $1
ORDER BY us.date_start DESC;`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := sm.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	subscriptions := []*Subscriptions{}

	for rows.Next() {
		var subscription Subscriptions

		err := rows.Scan(&subscription.UserID, &subscription.FullName, &subscription.SubscriptionID,
			&subscription.SubscriptionName, &subscription.VisitsPerWeek,
			&subscription.Price, &subscription.DateStart, &subscription.DateEnd)
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
