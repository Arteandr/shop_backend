package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

type ItemsRepo struct {
	db *sqlx.DB
}

func NewItemsRepo(db *sqlx.DB) *ItemsRepo {
	return &ItemsRepo{db: db}
}

func (r *ItemsRepo) Create(item models.Item) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name,description,category_id,sku,price) VALUES ($1,$2,$3,$4,$5) RETURNING id;", itemsTable)
	row := r.db.QueryRow(query, item.Name, item.Description, item.CategoryId, item.Sku, item.Price)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ItemsRepo) LinkColor(itemId int, colorId int) error {
	query := fmt.Sprintf("INSERT INTO %s (item_id,color_id) VALUES ($1,$2);", itemsColorsTable)
	_, err := r.db.Exec(query, itemId, colorId)

	return err
}

func (r *ItemsRepo) LinkTag(itemId int, tag string) error {
	query := fmt.Sprintf("INSERT INTO %s (item_id, name) VALUES($1,$2);", tagsTable)
	_, err := r.db.Exec(query, itemId, tag)

	return err
}

func (r *ItemsRepo) LinkImage(itemId, imageId int) error {
	query := fmt.Sprintf("INSERT INTO %s (item_id, image_id) VALUES ($1, $2);", itemsImagesTable)
	_, err := r.db.Exec(query, itemId, imageId)

	return err
}

func (r *ItemsRepo) GetById(itemId int) (models.Item, error) {
	var item models.Item
	query := fmt.Sprintf("SELECT * FROM %s WHERE id=$1;", itemsTable)
	if err := r.db.QueryRow(query, itemId).Scan(&item.Id, &item.Name, &item.Description, &item.CategoryId, &item.Price, &item.Sku); err != nil {
		return models.Item{}, err
	}

	return item, nil
}

func (r *ItemsRepo) GetBySku(sku string) (models.Item, error) {
	var item models.Item
	query := fmt.Sprintf("SELECT * FROM %s where sku=$1;", itemsTable)
	if err := r.db.QueryRow(query, sku).Scan(&item.Id, &item.Name, &item.Description, &item.CategoryId, &item.Price, &item.Sku); err != nil {
		return models.Item{}, err
	}

	return item, nil
}

func (r *ItemsRepo) GetByCategory(categoryId int) ([]models.Item, error) {
	var items []models.Item
	query := fmt.Sprintf("SELECT * FROM %s WHERE category_id=$1;", itemsTable)
	if err := r.db.Select(&items, query, categoryId); err != nil {
		return []models.Item{}, err
	}

	return items, nil
}

func (r *ItemsRepo) GetByTag(tag string) ([]models.Item, error) {
	var items []models.Item
	query := fmt.Sprintf("SELECT items.id,items.name,items.description,items.category_id,items.price,items.sku FROM %s,%s WHERE %s.name = $1 AND %s.id = %s.item_id;", itemsTable, tagsTable, tagsTable, itemsTable, tagsTable)
	if err := r.db.Select(&items, query, tag); err != nil {
		return []models.Item{}, err
	}

	return items, nil
}

func (r *ItemsRepo) GetColors(itemId int) ([]models.Color, error) {
	var colors []models.Color
	query := fmt.Sprintf("SELECT colors.id, colors.name, colors.hex, colors.price FROM %s, %s WHERE colors.id = %s.color_id AND %s.item_id = $1;", colorsTable, itemsColorsTable, itemsColorsTable, itemsColorsTable)
	if err := r.db.Select(&colors, query, itemId); err != nil {
		return []models.Color{}, err
	}

	return colors, nil
}

func (r *ItemsRepo) GetTags(itemId int) ([]models.Tag, error) {
	var tags []models.Tag
	query := fmt.Sprintf("SELECT * FROM %s WHERE tags.item_id = $1;", tagsTable)
	if err := r.db.Select(&tags, query, itemId); err != nil && err != sql.ErrNoRows {
		return []models.Tag{}, err
	}

	return tags, nil
}

func (r *ItemsRepo) Delete(itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", itemsTable)
	_, err := r.db.Exec(query, itemId)

	return err
}

func (r *ItemsRepo) Exist(itemId int) (bool, error) {
	var exist bool
	queryMain := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", itemsTable)
	query := fmt.Sprintf("SELECT exists (%s)", queryMain)
	if err := r.db.QueryRow(query, itemId).Scan(&exist); err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exist, nil
}
