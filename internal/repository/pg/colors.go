package pg

import (
	"context"
	"fmt"

	"shop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type ColorsRepo struct {
	db *sqlx.DB
}

func NewColorsRepo(db *sqlx.DB) *ColorsRepo {
	return &ColorsRepo{db: db}
}

func (r *ColorsRepo) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	var (
		tx  *sqlx.Tx
		err error
	)

	// Check if transaction is existed in ctx
	existingTx := extractTx(ctx)
	if existingTx != nil {
		tx = existingTx
	} else {
		tx, err = r.db.Beginx()
		if err != nil {
			return fmt.Errorf("begin transcation: %w", err)
		}
	}

	if err := tFunc(injectTx(ctx, tx)); err != nil {
		if existingTx == nil {
			tx.Rollback()
		}

		return err
	}

	if existingTx == nil {
		tx.Commit()
	}

	return nil
}

func (r *ColorsRepo) GetInstance(ctx context.Context) SqlxDB {
	if tx := extractTx(ctx); tx != nil {
		return tx
	}

	return r.db
}

func (r *ColorsRepo) Create(ctx context.Context, color models.Color) (int, error) {
	var (
		db = r.GetInstance(ctx)
		id int
	)

	query := fmt.Sprintf("INSERT INTO %s (name,hex,price) VALUES ($1,$2,$3) RETURNING id;", colorsTable)
	if err := db.GetContext(ctx, &id, query, color.Name, color.Hex, color.Price); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ColorsRepo) Exist(ctx context.Context, colorId int) (bool, error) {
	var (
		db    = r.GetInstance(ctx)
		exist bool
	)

	queryMain := fmt.Sprintf("SELECT name FROM %s WHERE id=$1", colorsTable)

	query := fmt.Sprintf("SELECT exists (%s)", queryMain)
	if err := db.GetContext(ctx, &exist, query, colorId); err != nil {
		return false, err
	}

	return exist, nil
}

func (r *ColorsRepo) Delete(ctx context.Context, colorId int) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", colorsTable)
	_, err := db.ExecContext(ctx, query, colorId)

	return err
}

func (r *ColorsRepo) DeleteFromItems(ctx context.Context, colorId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE color_id=$1;", itemsColorsTable)
	_, err := r.db.Exec(query, colorId)

	return err
}

func (r *ColorsRepo) AddToItems(ctx context.Context, colorId int) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("INSERT INTO %s (item_id, color_id) SELECT id, $1 from %s ON CONFLICT (item_id, color_id) DO NOTHING;", itemsColorsTable, itemsTable)
	_, err := db.ExecContext(ctx, query, colorId)

	return err
}

func (r *ColorsRepo) Update(ctx context.Context, color models.Color) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("UPDATE %s SET name=$1,hex=$2,price=$3 WHERE id=$4", colorsTable)
	_, err := db.ExecContext(ctx, query, color.Name, color.Hex, color.Price, color.Id)

	return err
}

func (r *ColorsRepo) GetById(ctx context.Context, colorId int) (models.Color, error) {
	var (
		db    = r.GetInstance(ctx)
		color models.Color
	)

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", colorsTable)
	if err := db.QueryRowContext(ctx, query, colorId).Scan(&color.Id, &color.Name, &color.Hex, &color.Price); err != nil {
		return models.Color{}, err
	}

	return color, nil
}

func (r *ColorsRepo) GetAll(ctx context.Context) ([]models.Color, error) {
	var (
		db     = r.GetInstance(ctx)
		colors []models.Color
	)

	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id;", colorsTable)

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var color models.Color
		if err := rows.StructScan(&color); err != nil {
			return nil, err
		}

		colors = append(colors, color)
	}

	return colors, nil
}
