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
	query := fmt.Sprintf("INSERT INTO %s (name,description,category_id,tags,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING id;", itemsTable)
	fmt.Println(item.Tags)
	row := r.db.QueryRow(query, item.Name, item.Description, item.CategoryId, pq.Array(item.Tags), item.CreatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
