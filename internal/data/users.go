package data

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	FullName       string    `json:"full_name"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	RoleID         int64     `json:"role_id"`
	Image          string    `json:"image_url"`
}

type UserModel struct {
	DB *sql.DB
}

func (um UserModel) GetTrainers() ([]*User, error) {
	query := `SELECT t.user_id, full_name, u.image FROM users u JOIN trainers t ON t.user_id = u.id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := um.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	trainers := []*User{}

	for rows.Next() {
		var trainer User

		err := rows.Scan(&trainer.ID, &trainer.FullName, &trainer.Image)
		if err != nil {
			return nil, err
		}

		trainers = append(trainers, &trainer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return trainers, nil
}
