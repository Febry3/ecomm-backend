package repository

import "context"

// TxKey is the context key for transactions
type TxKey struct{}

// TxManager defines an abstract transaction manager interface.
// This allows usecases to manage transactions without depending on specific database implementations.
type TxManager interface {
	// WithTransaction executes fn within a database transaction.
	// If fn returns an error, the transaction is rolled back.
	// If fn returns nil, the transaction is committed.
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
