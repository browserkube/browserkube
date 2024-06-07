package broadcast

import (
	"context"
	"time"
)

// Batcher creates time-based batches
type Batcher[T any] struct {
	tFrame     time.Duration
	emptyBatch bool
}

// NewBatcher creates new batcher
func NewBatcher[T any](tFrame time.Duration, emptyBatch bool) *Batcher[T] {
	return &Batcher[T]{tFrame: tFrame, emptyBatch: emptyBatch}
}

// Batch starts batcher
func (b *Batcher[T]) Batch(ctx context.Context, inc <-chan T, callb func([]T) error) error {
	ticker := time.NewTicker(b.tFrame)
	defer ticker.Stop()

	var tSlice []T
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if len(tSlice) == 0 && !b.emptyBatch {
				continue
			}
			if err := callb(tSlice); err != nil {
				return err
			}
			tSlice = nil
		case t := <-inc:
			tSlice = append(tSlice, t)
		}
	}
}
