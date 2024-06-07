package broadcast

import (
	"sync"
	"testing"
)

func TestBroadcast(t *testing.T) {
	wg := sync.WaitGroup{}

	b := NewBroadcaster[int](100)
	defer b.Close()

	for i := 0; i < 5; i++ {
		wg.Add(1)

		cch := make(chan int)

		b.Register(cch)

		go func() {
			defer wg.Done()
			defer b.Deregister(cch)
			<-cch
		}()

	}

	b.Submit(1)

	wg.Wait()
}

func TestBroadcastTrySubmit(t *testing.T) {
	b := NewBroadcaster[int](1)
	defer b.Close()

	if ok := b.TrySubmit(0); !ok {
		t.Fatalf("1st TrySubmit assert error expect=true actual=%v", ok)
	}

	if ok := b.TrySubmit(1); ok {
		t.Fatalf("2nd TrySubmit assert error expect=false actual=%v", ok)
	}

	cch := make(chan int)
	b.Register(cch)

	if ok := b.TrySubmit(1); !ok {
		t.Fatalf("3rd TrySubmit assert error expect=true actual=%v", ok)
	}
}

func TestBroadcastCleanup(t *testing.T) {
	b := NewBroadcaster[int](100)
	b.Register(make(chan int))
	b.Close()
}
