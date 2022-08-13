package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

type ImagesRepo struct {
	db *sqlx.DB
}

func NewImagesRepo(db *sqlx.DB) *ImagesRepo {
	return &ImagesRepo{
		db: db,
	}
}

func (r *ImagesRepo) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	var tx *sqlx.Tx
	var err error
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

func (r *ImagesRepo) GetInstance(ctx context.Context) SqlxDB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *ImagesRepo) Upload(ctx context.Context, filename string) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("INSERT INTO %s (filename) VALUES($1) RETURNING id;", imagesTable)
	_, err := db.ExecContext(ctx, query, filename)

	return err
}

func (r *ImagesRepo) GetById(ctx context.Context, imageId int) (models.Image, error) {
	var image models.Image
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", imagesTable)
	if err := r.db.QueryRow(query, imageId).Scan(&image.Id, &image.Filename, &image.CreatedAt); err != nil {
		return models.Image{}, err
	}

	return image, nil
}

func (r *ImagesRepo) GetAll(context.Context) ([]models.Image, error) {
	var images []models.Image
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC;", imagesTable)
	if err := r.db.Select(&images, query); err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImagesRepo) Exist(ctx context.Context, imageId int) (bool, error) {
	var exist bool
	queryMain := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", imagesTable)
	query := fmt.Sprintf("SELECT exists (%s)", queryMain)
	if err := r.db.QueryRow(query, imageId).Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

func (r *ImagesRepo) Delete(ctx context.Context, imageId int) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", imagesTable)
	_, err := db.ExecContext(ctx, query, imageId)

	return err
}

func (r *ImagesRepo) DeleteFromItems(ctx context.Context, imageId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE image_id=$1;", itemsImagesTable)
	_, err := r.db.Exec(query, imageId)

	return err
}
