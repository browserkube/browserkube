package broadcast

import "go.uber.org/zap"

//go:generate mockery --name Broadcaster --output mocks
type Broadcaster[T any] interface {
	Register(chan<- T)
	Deregister(chan<- T)
	Close() error
	Submit(T)
	TrySubmit(T) bool
}

type broadcaster[T any] struct {
	input  chan T
	reg    chan chan<- T
	unreg  chan chan<- T
	closed chan struct{}

	receivers map[chan<- T]bool
	logger    *zap.SugaredLogger
}

// NewBroadcaster creates a new broadcast with the given input
// channel buffer length.
func NewBroadcaster[T any](buflen int) Broadcaster[T] {
	b := &broadcaster[T]{
		input:     make(chan T, buflen),
		reg:       make(chan chan<- T),
		unreg:     make(chan chan<- T),
		receivers: make(map[chan<- T]bool),
		closed:    make(chan struct{}),
		logger:    zap.S(),
	}

	go b.run()

	return b
}

func (b *broadcaster[T]) broadcast(m T) {
	for ch := range b.receivers {
		ch <- m
	}
}

func (b *broadcaster[T]) run() {
	for {
		select {
		case m := <-b.input:
			b.broadcast(m)
		case ch, ok := <-b.reg:
			if ok {
				b.receivers[ch] = true
			} else {
				return
			}
			b.logger.Infof("Adding receiver. Total receivers: %d", len(b.receivers))

		case ch := <-b.unreg:
			delete(b.receivers, ch)
			b.logger.Infof("Deleting receiver. Total receivers: %d", len(b.receivers))
		case <-b.closed:
			return
		}
	}
}

func (b *broadcaster[T]) Register(newch chan<- T) {
	b.reg <- newch
}

func (b *broadcaster[T]) Deregister(ch chan<- T) {
	b.unreg <- ch
}

func (b *broadcaster[T]) Close() error {
	close(b.reg)
	close(b.unreg)
	close(b.closed)
	return nil
}

// Submit an item to be broadcast to all listeners.
func (b *broadcaster[T]) Submit(m T) {
	if b != nil {
		b.input <- m
	}
}

// TrySubmit attempts to submit an item to be broadcast, returning
// true iff it the item was broadcast, else false.
func (b *broadcaster[T]) TrySubmit(m T) bool {
	if b == nil {
		return false
	}
	select {
	case b.input <- m:
		return true
	default:
		return false
	}
}
