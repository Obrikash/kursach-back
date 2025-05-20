package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Pools         PoolModel
	Users         UserModel
	Groups        GroupModel
	Subscriptions SubscriptionModel
}

func NewModels(db *sql.DB) Models {
	return Models{Pools: PoolModel{DB: db},
		Users:         UserModel{DB: db},
		Groups:        GroupModel{DB: db},
		Subscriptions: SubscriptionModel{DB: db}}
}
