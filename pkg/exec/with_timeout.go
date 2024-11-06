package exec

import (
	"context"
)

func WithCtx(ctx context.Context, fn func() error) error {
	end := make(chan struct{})

	var err error
	go func() {
		err = fn()

		close(end)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-end:
		return err
	}
}
