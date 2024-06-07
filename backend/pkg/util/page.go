package browserkubeutil

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Page[T any] struct {
	Items         []T
	ContinueToken string
	Remaining     int64
}

func Empty[T any]() *Page[T] {
	return &Page[T]{}
}

func AsPage[T any](l metav1.ListInterface, items []T) *Page[T] {
	var remaining int64
	if remCount := l.GetRemainingItemCount(); remCount != nil {
		remaining = *remCount
	}
	return &Page[T]{
		Items:         items,
		ContinueToken: l.GetContinue(),
		Remaining:     remaining,
	}
}
