package browserkubeutil

import (
	"fmt"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	type args struct {
		items []string
		mapF  func(string) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{items: []string{}, mapF: func(s string) string { return strings.ToUpper(s) }},
			want: []string{},
		},
		{
			name: "nil F",
			args: args{items: []string{"one"}, mapF: nil},
			want: []string{},
		},
		{
			name: "happy",
			args: args{items: []string{"one", "two"}, mapF: func(s string) string { return strings.ToUpper(s) }},
			want: []string{"ONE", "TWO"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Map(tt.args.items, tt.args.mapF), "Map(%v, %v)", tt.args.items, tt.args.mapF)
		})
	}
}

func TestMapErr(t *testing.T) {
	type args[T any, V any] struct {
		items []T
		mapF  func(T) (V, error)
	}
	type testCase[T any, V any] struct {
		name    string
		args    args[T, V]
		want    []V
		wantErr assert.ErrorAssertionFunc
	}
	tests := []testCase[string, string]{
		{
			name: "empty",
			args: args[string, string]{
				items: []string{},
				mapF:  nil,
			},
			want:    []string{},
			wantErr: assert.NoError,
		},
		{
			name: "happy",
			args: args[string, string]{
				items: []string{"hello"},
				mapF: func(s string) (string, error) {
					return strings.ToUpper(s), nil
				},
			},
			want:    []string{"HELLO"},
			wantErr: assert.NoError,
		},
		{
			name: "err",
			args: args[string, string]{
				items: []string{"err"},
				mapF: func(s string) (string, error) {
					return s, errors.New("some horrible error")
				},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapErr(tt.args.items, tt.args.mapF)
			if !tt.wantErr(t, err, fmt.Sprintf("MapErr(%v, %v)", tt.args.items, err)) {
				return
			}
			assert.Equalf(t, tt.want, got, "MapErr(%v, %v)", tt.args.items, err)
		})
	}
}
