package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

type TxManagerPg struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) repository.TxManager {
	return &TxManagerPg{db: db}
}

func (t *TxManagerPg) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, repository.TxKey{}, tx)
		return fn(txCtx)
	})
}

func TxFromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(repository.TxKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
