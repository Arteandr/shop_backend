package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	query := fmt.Sprintf("INSERT INTO %s (name,description,category_id,tags,sku) VALUES ($1,$2,$3,$4,$5) RETURNING id;", itemsTable)
	row := r.db.QueryRow(query, item.Name, item.Description, item.CategoryId, pq.Array(item.Tags), item.Sku)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *ItemsRepo) LinkColor(itemId int, colorId int) error {
	query := fmt.Sprintf("INSERT INTO %s (item_id,color_id) VALUES ($1,$2);", itemColorsTable)
	_, err := r.db.Exec(query, itemId, colorId)

	return err
}

func (r *ItemsRepo) GetById(itemId int) (models.Item, error) {
	var item models.Item
	query := fmt.Sprintf("SELECT * from %s WHERE id=$1;", itemsTable)
	if err := r.db.QueryRow(query, itemId).Scan(&item); err != nil {
		return models.Item{}, err
	}

	return item, nil
}

func (r *ItemsRepo) GetColors(itemId int) ([]models.Color, error) {
	var colors []models.Color
	query := fmt.Sprintf("SELECT colors.id, colors.name, colors.hex, colors.price FROM %s, %s WHERE colors.id = item_colors.color_id AND item_colors.item_id = $1;", colorsTable, itemColorsTable)
	if err := r.db.Select(&colors, query, itemId); err != nil {
		return []models.Color{}, err
	}

	return colors, nil
}
