package pg

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"shop_backend/internal/models"
	apperrors "shop_backend/pkg/errors"
)

type CategoriesRepo struct {
	db *sqlx.DB
}

func NewCategoriesRepo(db *sqlx.DB) *CategoriesRepo {
	return &CategoriesRepo{db: db}
}

func (r *CategoriesRepo) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
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

func (r *CategoriesRepo) GetImage(ctx context.Context, categoryId int) (models.Image, error) {
	var (
		db    = r.GetInstance(ctx)
		image models.Image
	)

	query := fmt.Sprintf("SELECT i.id,i.filename,i.created_at FROM %s AS i, %s as c_i WHERE i.id=c_i.image_id AND c_i.category_id=$1 LIMIT 1;", imagesTable, categoriesImagesTable)
	if err := db.QueryRowContext(ctx, query, categoryId).Scan(&image.Id, &image.Filename, &image.CreatedAt); err != nil {
		return models.Image{}, err
	}

	return image, nil
}

func (r *CategoriesRepo) LinkImage(ctx context.Context, categoryId, imageId int) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("INSERT INTO %s (category_id,image_id) VALUES ($1,$2);", categoriesImagesTable)
	_, err := db.ExecContext(ctx, query, categoryId, imageId)

	return err
}

func (r *CategoriesRepo) GetInstance(ctx context.Context) SqlxDB {
	if tx := extractTx(ctx); tx != nil {
		return tx
	}

	return r.db
}

func (r *CategoriesRepo) Create(ctx context.Context, category models.Category) (int, error) {
	var (
		db = r.GetInstance(ctx)
		id int
	)

	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id;", categoriesTable)
	if err := db.GetContext(ctx, &id, query, category.Name); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *CategoriesRepo) Exist(ctx context.Context, categoryId int) (bool, error) {
	var (
		db    = r.GetInstance(ctx)
		exist bool
	)

	queryMain := fmt.Sprintf("SELECT name FROM %s WHERE id=$1", categoriesTable)

	query := fmt.Sprintf("SELECT exists (%s)", queryMain)
	if err := db.GetContext(ctx, &exist, query, categoryId); err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exist, nil
}

func (r *CategoriesRepo) Delete(ctx context.Context, categoryId int) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", categoriesTable)
	_, err := db.ExecContext(ctx, query, categoryId)

	pqError, ok := err.(*pq.Error)
	if ok {
		if pqError.Code == "23503" {
			return apperrors.ErrViolatesKey
		} else {
			return err
		}
	}

	return err
}

func (r *CategoriesRepo) GetAll(ctx context.Context) ([]models.Category, error) {
	var (
		db         = r.GetInstance(ctx)
		categories []models.Category
	)

	query := fmt.Sprintf("SELECT * FROM %s;", categoriesTable)

	rows, err := db.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.Id, &c.Name); err != nil {
			return nil, err
		}

		categories = append(categories, c)
	}

	return categories, nil
}

func (r *CategoriesRepo) GetById(ctx context.Context, categoryId int) (models.Category, error) {
	var (
		db       = r.GetInstance(ctx)
		category models.Category
	)

	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", categoriesTable)

	rows, err := db.QueryxContext(ctx, query, categoryId)
	if err != nil {
		return models.Category{}, err
	}

	for rows.Next() {
		if err := rows.StructScan(&category); err != nil {
			return models.Category{}, err
		}
	}

	return category, nil
}

// $1 = category.Name
// $2 = category.Id
func (r *CategoriesRepo) Update(ctx context.Context, category models.Category) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("UPDATE %s SET name=$1 WHERE id=$2;", categoriesTable)
	_, err := db.ExecContext(ctx, query, category.Name, category.Id)

	return err
}

func (r *CategoriesRepo) UpdateImage(ctx context.Context, categoryId, imageId int) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("UPDATE %s SET image_id=$1 WHERE category_id=$2;", categoriesImagesTable)
	_, err := db.ExecContext(ctx, query, imageId, categoryId)

	return err
}
