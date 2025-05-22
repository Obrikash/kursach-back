package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrTrainerOnlyOnePool = errors.New("Trainer could only work in one pool, not at both")
)

type Group struct {
	ID       int64 `json:"id"`
	Pool     int64 `json:"pool"`
	Category int64 `json:"category"`
	Trainer  User  `json:"trainer"`
}

type Groups struct {
	ID          int64  `json:"id"`
	Category    string `json:"category"`
	PoolName    string `json:"pool_name"`
	TrainerName string `json:"trainer_name"`
	UserID      int64  `json:"user_id"`
	Image       string `json:"image"`
}

type GroupModel struct {
	DB *sql.DB
}

func (gm GroupModel) GetGroups() ([]*Groups, error) {
	query := `SELECT g.id, c.name as "category", p.name as "pool", tr.full_name, t.user_id, tr.image
	FROM training_groups g JOIN group_category c ON g.category_id = c.id 
	JOIN pools p ON g.pool_id = p.id JOIN trainers t ON g.trainer_id = t.id JOIN users tr ON t.user_id = tr.id;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := gm.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	groups := []*Groups{}

	for rows.Next() {
		var group Groups

		err := rows.Scan(&group.ID, &group.Category, &group.PoolName, &group.TrainerName, &group.UserID, &group.Image)
		if err != nil {
			return nil, err
		}

		groups = append(groups, &group)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (gm GroupModel) AddToPool(group *Group) error {
	// we make this check so one trainer would not work in 2 pools, only at 1
	query := `INSERT INTO training_groups (pool_id, category_id, trainer_id) SELECT $1, $2, id FROM trainers WHERE id = $3 and pool_id = $1 RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{group.Pool, group.Category, group.Trainer.ID}

	err := gm.DB.QueryRowContext(ctx, query, args...).Scan(&group.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrTrainerOnlyOnePool
		default:
			return err
		}
	}

	return nil
}
