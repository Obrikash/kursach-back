package data

import (
	"context"
	"database/sql"
	"time"
)

type Group struct {
	ID       int64  `json:"id"`
	Pool     string `json:"pool"`
	Category string `json:"category"`
	Trainer  User   `json:"trainer"`
}

type GroupModel struct {
	DB *sql.DB
}

func (gm GroupModel) GetGroups() ([]*Group, error) {
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

	groups := []*Group{}

	for rows.Next() {
		var group Group

		err := rows.Scan(&group.ID, &group.Category, &group.Pool, &group.Trainer.FullName, &group.Trainer.ID, &group.Trainer.Image)
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
