package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

	query := `SELECT p.id, p.name. p.address, p.type FROM pools p WHERE p.id = $1`

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
	query := fmt.Sprintf("SELECT p.id, p.name, p.address, p.type FROM pools p ORDER BY name ASC")

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
