package repository

import "github.com/jmoiron/sqlx"

type ItemsRepo struct {
	db *sqlx.DB
}

func NewItemsRepo(db *sqlx.DB) *ItemsRepo {
	return &ItemsRepo{db: db}
}
