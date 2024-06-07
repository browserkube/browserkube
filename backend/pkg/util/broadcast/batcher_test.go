package broadcast

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestBatcher_Batch(t *testing.T) {
	batcher := NewBatcher[int](2*time.Second, false)

	ctx, cancelF := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelF()

	ch := make(chan int)
	go func() {
		idx := 0
		for {
			select {
			case <-time.After(1 * time.Second):
				idx++
				ch <- idx
			case <-ctx.Done():
				return
			}
		}
	}()

	var prev []int
	err := batcher.Batch(ctx, ch, func(batch []int) error {
		assert.Greater(t, len(batch), 0)

		if len(prev) > 0 {
			lastInt := prev[len(prev)-1]
			assert.Equal(t, lastInt, batch[0])
		}
		return nil
	})
	assert.NoError(t, err)
}

func TestBatcher_Empty(t *testing.T) {
	t.Run("Empty Disabled", func(t *testing.T) {
		batcher := NewBatcher[int](1*time.Second, false)

		ctx, cancelF := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelF()

		ch := make(chan int)
		err := batcher.Batch(ctx, ch, func(batch []int) error {
			assert.Fail(t, "Nothing is sent")
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Empty Enabled", func(t *testing.T) {
		batcher := NewBatcher[int](1*time.Second, true)

		ctx, cancelF := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancelF()

		ch := make(chan int)

		sent := atomic.NewBool(false)
		err := batcher.Batch(ctx, ch, func(batch []int) error {
			sent.Store(true)
			return nil
		})
		assert.NoError(t, err)
		assert.True(t, sent.Load(), "Empty Batch isn't sent")
	})
}
