package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"shop_backend/internal/models"
)

type DeliveryRepo struct {
	db *sqlx.DB
}

func NewDeliveryRepo(db *sqlx.DB) *DeliveryRepo {
	return &DeliveryRepo{
		db: db,
	}
}

func (r *DeliveryRepo) Create(ctx context.Context, delivery models.Delivery) (int, error) {
	var id int
	subquery := fmt.Sprintf("SELECT id FROM %s WHERE name=$2", deliveryCompanyTable)
	query := fmt.Sprintf("INSERT INTO %s (name,company_id,price) VALUES($1,(%s),$3) RETURNING id;", deliveryTable, subquery)
	if err := r.db.QueryRowContext(ctx, query, delivery.Name, delivery.CompanyName, delivery.Price).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DeliveryRepo) CreateCompany(ctx context.Context, name string) error {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES($1) RETURNING id;", deliveryCompanyTable)
	if err := r.db.QueryRowContext(ctx, query, name).Scan(&id); err != nil {
		return err
	}

	return nil
}

func (r *DeliveryRepo) ExistCompany(ctx context.Context, name string) (bool, error) {
	var exist bool
	queryMain := fmt.Sprintf("SELECT * FROM %s WHERE name=$1", deliveryCompanyTable)
	query := fmt.Sprintf("SELECT exists (%s)", queryMain)
	if err := r.db.QueryRowContext(ctx, query, name).Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}
