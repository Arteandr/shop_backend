package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

type ColorsRepo struct {
	db *sqlx.DB
}

func NewColorsRepo(db *sqlx.DB) *ColorsRepo {
	return &ColorsRepo{db: db}
}

func (r *ColorsRepo) Create(color models.Color) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name,hex,price) VALUES ($1,$2,$3) RETURNING id;", colorsTable)
	row := r.db.QueryRow(query, color.Name, color.Hex, color.Price)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
