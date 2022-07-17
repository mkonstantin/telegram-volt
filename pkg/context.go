package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TX struct {
	Tx       *sqlx.Tx
	IsActive bool
}

type transactionCtxKey struct{}

func FromContext(ctx context.Context) *TX {
	v := ctx.Value(&transactionCtxKey{})
	if v == nil {
		return nil
	}
	return v.(*TX)
}

func ToContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, &transactionCtxKey{}, &TX{Tx: tx, IsActive: true})
}
