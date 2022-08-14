package repository

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
		pqError, ok := err.(*pq.Error)
		if ok {
			fmt.Println(pqError.Code)
		}
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
			fmt.Println(field)
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
