package utils

import (
	"sort"
	"strings"
	"time"
)

func SliceContains[T comparable](s []T, e T, comp func(o1 T, o2 T) bool) bool {
	for _, v := range s {
		if comp(v, e) {
			return true
		}
	}
	return false
}

func StringSliceContains(s []string, e string) bool {
	for _, v := range s {
		if strings.Contains(e, v) {
			return true
		}
	}
	return false
}

func SortStringTime(m map[string]time.Time) []string {
	type kv struct {
		Key   string
		Value time.Time
	}

	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{Key: k, Value: v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value.After(ss[j].Value)
	})
	ranked := make([]string, len(m))
	for i, kv := range ss {
		ranked[i] = kv.Key
	}
	return ranked
}
