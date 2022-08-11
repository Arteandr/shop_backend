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

func (r *DeliveryRepo) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
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

func (r *DeliveryRepo) GetDB(ctx context.Context) SqlxDB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx
	}
	return r.db
}

func (r *DeliveryRepo) Create(ctx context.Context, delivery models.Delivery) (int, error) {
	db := r.GetDB(ctx)
	var id int
	subquery := fmt.Sprintf("SELECT id FROM %s WHERE name=$2", deliveryCompanyTable)
	query := fmt.Sprintf("INSERT INTO %s (name,company_id,price) VALUES($1,(%s),$3) RETURNING id;", deliveryTable, subquery)
	if err := db.GetContext(ctx, &id, query, delivery.Name, delivery.CompanyName, delivery.Price); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *DeliveryRepo) CreateCompany(ctx context.Context, name string) error {
	db := r.GetDB(ctx)
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES($1) RETURNING id;", deliveryCompanyTable)
	_, err := db.ExecContext(ctx, query, name)

	return err
}

func (r *DeliveryRepo) ExistCompany(ctx context.Context, name string) (bool, error) {
	db := r.GetDB(ctx)
	var exist bool
	subquery := fmt.Sprintf("SELECT * FROM %s WHERE name=$1", deliveryCompanyTable)
	query := fmt.Sprintf("SELECT exists (%s)", subquery)
	if err := db.GetContext(ctx, &exist, query, name); err != nil {
		return false, err
	}

	return exist, nil
}

func (r *DeliveryRepo) GetById(ctx context.Context, deliveryId int) (models.Delivery, error) {
	var delivery models.Delivery
	query := fmt.Sprintf("SELECT d.id, d.name, dc.name company_name, d.price FROM %s d JOIN %s dc ON dc.id=d.company_id WHERE d.id=$1 LIMIT 1;", deliveryTable, deliveryCompanyTable)
	rows, err := r.db.QueryxContext(ctx, query, deliveryId)
	if err != nil {
		return models.Delivery{}, err
	}

	for rows.Next() {
		if err := rows.StructScan(&delivery); err != nil {
			return models.Delivery{}, err
		}
	}

	if delivery == (models.Delivery{}) {
		return models.Delivery{}, models.ErrDeliveryNotFound
	}

	return delivery, nil
}

func (r *DeliveryRepo) Update(ctx context.Context, delivery models.Delivery) error {
	db := r.GetDB(ctx)
	subquery := fmt.Sprintf("SELECT id FROM %s WHERE name=$2", deliveryCompanyTable)
	query := fmt.Sprintf("UPDATE %s SET name=$1,company_id=(%s),price=$3 WHERE id=$4;", deliveryTable, subquery)
	rows, err := db.ExecContext(ctx, query, delivery.Name, delivery.CompanyName, delivery.Price, delivery.Id)
	rowsAffected, _ := rows.RowsAffected()
	fmt.Println(rowsAffected)
	if rowsAffected < 1 {
		return models.ErrDeliveryNotFound
	}

	return err
}

func (r *DeliveryRepo) Delete(ctx context.Context, deliveryId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", deliveryTable)
	_, err := r.db.ExecContext(ctx, query, deliveryId)

	return err
}
