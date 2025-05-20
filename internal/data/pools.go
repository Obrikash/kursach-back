package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Pool struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	PoolType string `json:"type"`
}

type PoolModel struct {
	DB *sql.DB
}

func (pm PoolModel) Get(id int64) (*Pool, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT p.id, p.name, p.address, p.type FROM pools p WHERE p.id = $1`

	pool := &Pool{}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pm.DB.QueryRowContext(ctx, query, id).Scan(&pool.ID, &pool.Name, &pool.Address, &pool.PoolType)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return pool, nil
}

func (pm PoolModel) GetAll() ([]*Pool, error) {
	query := "SELECT p.id, p.name, p.address, p.type FROM pools p ORDER BY name ASC"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := pm.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pools := []*Pool{}

	for rows.Next() {
		var pool Pool

		err := rows.Scan(&pool.ID, &pool.Name, &pool.Address, &pool.PoolType)
		if err != nil {
			return nil, err
		}
		pools = append(pools, &pool)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pools, nil
}

func (pm PoolModel) MaxProfit() (*Pool, float64, error) {
	query := `SELECT p.id AS pool_id, p.name AS pool_name, p.address, p.type, SUM(sub.price) AS total_revenue 
	FROM user_subscriptions us JOIN user_groups ug ON us.user_id = ug.user_id JOIN training_groups tg ON ug.group_id = tg.id 
	JOIN trainers tr ON tg.trainer_id = tr.id JOIN pools p ON tr.pool_id = p.id 
	JOIN subscriptions sub ON us.subscription_id = sub.id GROUP BY p.id, p.name, p.address, p.type ORDER BY total_revenue DESC LIMIT 1;`

	pool := &Pool{}
	var profit float64
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := pm.DB.QueryRowContext(ctx, query).Scan(&pool.ID, &pool.Name, &pool.Address, &pool.PoolType, &profit)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, 0, ErrRecordNotFound
		default:
			return nil, 0, err
		}
	}
	return pool, profit, nil
}
