package common

import "context"

// TxManager is an abstraction usecases depend on.
type TxManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type txMarkerKey struct{}

// MarkTx returns a ctx that is known to be transactional - for testing purposes
func MarkTx(ctx context.Context) context.Context {
	return context.WithValue(ctx, txMarkerKey{}, true)
}

// IsTransactional tells whether ctx is inside TxManager.Do
func IsTransactional(ctx context.Context) bool {
	v, ok := ctx.Value(txMarkerKey{}).(bool)
	return ok && v
}