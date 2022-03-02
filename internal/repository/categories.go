package repository

import "github.com/jmoiron/sqlx"

type CategoriesRepo struct {
	db *sqlx.DB
}

func NewCategoriesRepo(db *sqlx.DB) *CategoriesRepo {
	return &CategoriesRepo{db: db}
}
