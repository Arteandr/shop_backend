package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

type CategoriesRepo struct {
	db *sqlx.DB
}

func NewCategoriesRepo(db *sqlx.DB) *CategoriesRepo {
	return &CategoriesRepo{db: db}
}

func (r *CategoriesRepo) Create(category models.Category) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id;", categoriesTable)
	row := r.db.QueryRow(query, category.Name)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
