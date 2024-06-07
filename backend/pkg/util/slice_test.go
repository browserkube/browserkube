package browserkubeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlice(t *testing.T) {
	type testCase[T any] struct {
		name string
		arg1 []T
		arg2 T
		want bool
	}

	tests := []testCase[string]{
		{
			name: "empty",
			arg1: []string{},
			arg2: "",
			want: false,
		},
		{
			name: "happycase",
			arg1: []string{"test1", "test3", "test2"},
			arg2: "test2",
			want: true,
		},
		{
			name: "badcase",
			arg1: []string{"test5"},
			arg2: "test8",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := SliceContains(tt.arg1, tt.arg2)
			assert.Equal(t, tt.want, res)
		})
	}
}
