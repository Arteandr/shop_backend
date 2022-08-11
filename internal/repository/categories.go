package repository

import (
	"context"
	"database/sql"
	"errors"
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

func (r *CategoriesRepo) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
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

func (r *CategoriesRepo) GetInstance(ctx context.Context) SqlxDB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *CategoriesRepo) Create(ctx context.Context, category models.Category) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id;", categoriesTable)
	row := r.db.QueryRow(query, category.Name)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *CategoriesRepo) Exist(ctx context.Context, categoryId int) (bool, error) {
	var exist bool
	queryMain := fmt.Sprintf("SELECT name FROM %s WHERE id=$1", categoriesTable)
	query := fmt.Sprintf("SELECT exists (%s)", queryMain)
	if err := r.db.QueryRow(query, categoryId).Scan(&exist); err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exist, nil
}

func (r *CategoriesRepo) Delete(ctx context.Context, categoryId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", categoriesTable)
	_, err := r.db.Exec(query, categoryId)

	return err
}

func (r *CategoriesRepo) GetAllC(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	query := fmt.Sprintf("SELECT * FROM %s;", categoriesTable)
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoriesRepo) GetById(ctx context.Context, categoryId int) (models.Category, error) {
	var category models.Category
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", categoriesTable)
	if err := r.db.QueryRow(query, categoryId).Scan(&category.Id, &category.Name); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.Category{}, err
	}

	return category, nil
}

// $1 = category.Name
// $2 = category.Id
func (r *CategoriesRepo) Update(ctx context.Context, category models.Category) error {
	query := fmt.Sprintf("UPDATE %s SET name=$1 WHERE id=$2;", categoriesTable)
	_, err := r.db.Exec(query, category.Name, category.Id)

	return err
}
