package repository

import (
	"context"
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

func (r *ColorsRepo) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("begin transcation: %w", err)
	}

	if err := tFunc(injectTx(ctx, tx)); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *ColorsRepo) GetInstance(ctx context.Context) SqlxDB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *ColorsRepo) Create(ctx context.Context, color models.Color) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name,hex,price) VALUES ($1,$2,$3) RETURNING id;", colorsTable)
	row := r.db.QueryRow(query, color.Name, color.Hex, color.Price)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ColorsRepo) Exist(ctx context.Context, colorId int) (bool, error) {
	var exist bool
	queryMain := fmt.Sprintf("SELECT name FROM %s WHERE id=$1", colorsTable)
	query := fmt.Sprintf("SELECT exists (%s)", queryMain)
	if err := r.db.QueryRow(query, colorId).Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

func (r *ColorsRepo) Delete(ctx context.Context, colorId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", colorsTable)
	_, err := r.db.Exec(query, colorId)

	return err
}

func (r *ColorsRepo) DeleteFromItems(ctx context.Context, colorId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE color_id=$1;", itemsColorsTable)
	_, err := r.db.Exec(query, colorId)

	return err
}

func (r *ColorsRepo) AddToItems(ctx context.Context, colorId int) error {
	query := fmt.Sprintf("INSERT INTO %s (item_id,color_id) SELECT id, %d from %s;", itemsColorsTable, colorId, itemsTable)
	_, err := r.db.Exec(query)

	return err
}

func (r *ColorsRepo) Update(ctx context.Context, color models.Color) error {
	query := fmt.Sprintf("UPDATE %s SET name=$1,hex=$2,price=$3 WHERE id=$4", colorsTable)
	_, err := r.db.Exec(query, color.Name, color.Hex, color.Price, color.Id)

	return err
}

func (r *ColorsRepo) GetById(ctx context.Context, colorId int) (models.Color, error) {
	var color models.Color
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", colorsTable)
	if err := r.db.QueryRow(query, colorId).Scan(&color.Id, &color.Name, &color.Hex, &color.Price); err != nil {
		return models.Color{}, err
	}

	return color, nil
}

func (r *ColorsRepo) GetAll(ctx context.Context) ([]models.Color, error) {
	var colors []models.Color
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id;", colorsTable)
	if err := r.db.Select(&colors, query); err != nil {
		return []models.Color{}, err
	}

	return colors, nil
}
