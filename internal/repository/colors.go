package repository

import (
	"database/sql"
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

func (r *ColorsRepo) Exist(colorId int) (bool, error) {
	var exist bool
	queryMain := fmt.Sprintf("SELECT name FROM %s WHERE id=$1", colorsTable)
	query := fmt.Sprintf("SELECT exists (%s)", queryMain)
	if err := r.db.QueryRow(query, colorId).Scan(&exist); err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exist, nil
}

func (r *ColorsRepo) Delete(colorId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", colorsTable)
	_, err := r.db.Exec(query, colorId)

	return err
}
