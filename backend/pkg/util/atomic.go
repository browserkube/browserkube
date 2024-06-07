package browserkubeutil

import (
	"sync/atomic"
)

type TypedAtomic[T any] struct {
	val *atomic.Value
}

func NewTypedAtomic[T any]() *TypedAtomic[T] {
	at := &atomic.Value{}
	return &TypedAtomic[T]{val: at}
}

func (ta *TypedAtomic[T]) Load() T {
	val, ok := ta.val.Load().(T)
	if !ok {
		var r T
		return r
	}
	return val
}

func (ta *TypedAtomic[T]) Set(val T) {
	ta.val.Store(val)
}
