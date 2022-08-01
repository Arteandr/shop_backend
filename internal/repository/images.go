package repository

import (
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

func (r *ImagesRepo) Upload(filename string) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (filename) VALUES($1) RETURNING id;", imagesTable)
	if err := r.db.QueryRow(query, filename).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ImagesRepo) GetById(imageId int) (models.Image, error) {
	var image models.Image
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", imagesTable)
	if err := r.db.QueryRow(query, imageId).Scan(&image.Id, &image.Filename, &image.CreatedAt); err != nil {
		return models.Image{}, err
	}

	return image, nil
}

func (r *ImagesRepo) GetAll() ([]models.Image, error) {
	var images []models.Image
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC;", imagesTable)
	if err := r.db.Select(&images, query); err != nil {
		return nil, err
	}

	return images, nil
}

func (r *ImagesRepo) Exist(imageId int) (bool, error) {
	var exist bool
	queryMain := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", imagesTable)
	query := fmt.Sprintf("SELECT exists (%s)", queryMain)
	if err := r.db.QueryRow(query, imageId).Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

func (r *ImagesRepo) Delete(imageId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", imagesTable)
	_, err := r.db.Exec(query, imageId)

	return err
}

func (r *ImagesRepo) DeleteFromItems(imageId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE image_id=$1;", itemsImagesTable)
	_, err := r.db.Exec(query, imageId)

	return err
}
