package common

import "context"

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) (err error)
}
