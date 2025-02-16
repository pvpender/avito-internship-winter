package mock

import "context"

type MockTransactionManager struct {
}

func (m MockTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) (err error) {
	if err = fn(ctx); err != nil {
		return err
	}

	return nil
}
