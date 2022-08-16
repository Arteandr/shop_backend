package pg

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type Transactor interface {
	GetInstance(ctx context.Context) SqlxDB
	WithinTransaction(context.Context, func(ctx context.Context) error) error
}

type txKey struct{}

func injectTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}
