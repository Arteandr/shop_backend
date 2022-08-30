package pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"shop_backend/internal/models"
	apperrors "shop_backend/pkg/errors"
	"strings"
)

type OrdersRepo struct {
	db *sqlx.DB
}

func NewOrdersRepo(db *sqlx.DB) *OrdersRepo {
	return &OrdersRepo{
		db: db,
	}
}

func (r *OrdersRepo) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
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

func (r *OrdersRepo) GetInstance(ctx context.Context) SqlxDB {
	tx := extractTx(ctx)
	if tx != nil {
		return tx
	}
	return r.db
}

// $1 = userId
// $2 = deliveryId
func (r *OrdersRepo) Create(ctx context.Context, userId int, deliveryId int) (int, error) {
	db := r.GetInstance(ctx)
	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_id,delivery_id) VALUES ($1,$2) RETURNING id;", ordersTable)
	err := db.GetContext(ctx, &id, query, userId, deliveryId)
	if err != nil {
		//pqError, ok := err.(*pq.Error)
		//if ok {
		//	return 0, err
		//}
		return 0, err
	}

	return id, nil
}

// $1 = orderId
// $2 = itemId
// $3 = colorId
// $4 = quantity
func (r *OrdersRepo) LinkItem(ctx context.Context, orderId, itemId, colorId, quantity int) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("INSERT INTO %s (order_id,item_id,color_id,quantity) VALUES ($1,$2,$3,$4);", orderItemsTable)
	_, err := db.ExecContext(ctx, query, orderId, itemId, colorId, quantity)
	if err != nil {
		pqError, ok := err.(*pq.Error)
		if ok {
			field := strings.Split(pqError.Constraint, "_")[2]
			var id int
			switch field {
			case "item":
				id = itemId
			case "color":
				id = colorId
			}
			return apperrors.ErrIdNotFound(field, id)
		}
	}
	return err
}

// $1 = orderId
func (r *OrdersRepo) Exist(ctx context.Context, orderId int) (bool, error) {
	db := r.GetInstance(ctx)
	var exist bool
	subquery := fmt.Sprintf("SELECT * FROM %s WHERE id=$1", ordersTable)
	query := fmt.Sprintf("SELECT exists (%s)", subquery)
	if err := db.GetContext(ctx, &exist, query, orderId); err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exist, nil
}

// $1 = orderId
func (r *OrdersRepo) Delete(ctx context.Context, orderId int) error {
	db := r.GetInstance(ctx)
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1;", ordersTable)
	_, err := db.ExecContext(ctx, query, orderId)

	return err
}

// $1 = userId
func (r *OrdersRepo) GetAllByUserId(ctx context.Context, userId int) ([]models.Order, error) {
	db := r.GetInstance(ctx)
	var orders []models.Order
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1;", ordersTable)
	rows, err := db.QueryxContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.Order
		if err := rows.StructScan(&order); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrdersRepo) GetItems(ctx context.Context, orderId int) ([]models.ServiceOrderItem, error) {
	db := r.GetInstance(ctx)
	var items []models.ServiceOrderItem
	query := fmt.Sprintf("SELECT o.item_id, o.color_id, o.quantity FROM %s o WHERE o.order_id=$1;", orderItemsTable)
	rows, err := db.QueryxContext(ctx, query, orderId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var item models.ServiceOrderItem
		if err := rows.StructScan(&item); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *OrdersRepo) GetStatus(ctx context.Context, statusId int) (string, error) {
	db := r.GetInstance(ctx)
	var name string
	query := fmt.Sprintf("SELECT name FROM %s WHERE id=$1;", statusTable)
	if err := db.GetContext(ctx, &name, query, statusId); err != nil {
		return "", err
	}

	return name, nil
}
