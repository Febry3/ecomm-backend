package repository

import "context"

type TxKey struct{}

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
