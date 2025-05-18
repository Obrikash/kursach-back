package data

import (
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
}

type UserModel struct {
	DB *sql.DB
}
