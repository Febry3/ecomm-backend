package pg

import (
	"context"

	"github.com/febry3/gamingin/internal/repository"
	"gorm.io/gorm"
)

// TxManagerPg implements repository.TxManager using GORM
type TxManagerPg struct {
	db *gorm.DB
}

// NewTxManager creates a new TxManagerPg
func NewTxManager(db *gorm.DB) repository.TxManager {
	return &TxManagerPg{db: db}
}

// WithTransaction executes fn within a database transaction.
// The transaction is stored in context so repositories can access it via TxFromContext.
func (t *TxManagerPg) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, repository.TxKey{}, tx)
		return fn(txCtx)
	})
}

// TxFromContext retrieves the transaction from context, or returns the default DB if none exists.
// Repositories should use this function to get the appropriate database connection.
func TxFromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(repository.TxKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
