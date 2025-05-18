package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Pools PoolModel
}

func NewModels(db *sql.DB) Models {
	return Models{Pools: PoolModel{DB: db}}
}
