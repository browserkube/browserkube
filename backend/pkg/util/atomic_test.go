package browserkubeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomic(t *testing.T) {
	ta := NewTypedAtomic[string]()
	ta.Set("hello world")
	val := ta.Load()
	assert.Equal(t, "hello world", val)
}

func TestAtomicEmptyLoad(t *testing.T) {
	ta := NewTypedAtomic[string]()
	val := ta.Load()
	assert.Equal(t, "", val)
}

func TestAtomicEmptyStruct(t *testing.T) {
	ta := NewTypedAtomic[*testStruct]()
	val := ta.Load()
	assert.True(t, val == nil)
}

type testStruct struct{}
