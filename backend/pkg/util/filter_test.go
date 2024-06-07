package browserkubeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	type args struct {
		items  []string
		filter func(string) bool
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "has items",
			args: args{items: []string{"A", "B", "C"}, filter: func(str string) bool { return str == "A" }},
			want: []string{"A"},
		},
		{
			name: "all items",
			args: args{items: []string{"A", "B", "C"}, filter: func(str string) bool { return true }},
			want: []string{"A", "B", "C"},
		},
		{
			name: "no items",
			args: args{items: []string{"A", "B", "C"}, filter: func(str string) bool { return false }},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, Filter(tt.args.items, tt.args.filter), "Filter(%v, %v)", tt.args.items, tt.args.filter)
		})
	}
}

func TestFindFirst(t *testing.T) {
	type args struct {
		items  []string
		filter func(string) bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "found",
			args: args{items: []string{"A", "B", "C"}, filter: func(str string) bool { return str == "B" }},
			want: "B",
		},
		{
			name: "all items",
			args: args{items: []string{"A", "B", "C"}, filter: func(str string) bool { return true }},
			want: "A",
		},
		{
			name: "no items",
			args: args{items: []string{"A", "B", "C"}, filter: func(str string) bool { return false }},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FindFirst(tt.args.items, tt.args.filter), "FindFirst(%v, %v)", tt.args.items, tt.args.filter)
		})
	}
}

func TestReverseDeduplicateBY(t *testing.T) {
	type args struct {
		items []string
		uniqF func(string) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "happy",
			args: args{items: []string{"one", "two", "one"}, uniqF: func(s string) string { return s }},
			want: []string{"one", "two"},
		},
		{
			name: "happy-start-with",
			args: args{items: []string{"one[1]", "two", "one[2]"}, uniqF: func(s string) string { return s[0:1] }},
			want: []string{"one[2]", "two"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ReverseDeduplicateBY(tt.args.items, tt.args.uniqF), "ReverseDeduplicateBY(%v, %v)", tt.args.items, tt.args.uniqF)
		})
	}
}
