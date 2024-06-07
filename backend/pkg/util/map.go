package browserkubeutil

import "github.com/pkg/errors"

// Map converts an array of one type to another
func Map[T, V any](items []T, mapF func(T) V) []V {
	if mapF == nil {
		empty := make([]V, 0)
		return empty
	}
	transformed := make([]V, len(items))
	for i, item := range items {
		transformed[i] = mapF(item)
	}
	return transformed
}

// MapErr converts an array of one type to another with respect to transformation errors
func MapErr[T, V any](items []T, mapF func(T) (V, error)) ([]V, error) {
	if mapF == nil {
		empty := make([]V, 0)
		return empty, nil
	}

	var err error
	transformed := make([]V, len(items))
	for i, item := range items {
		transformed[i], err = mapF(item)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return transformed, nil
}

// MapPageErr converts a page of one type to another with respect to transformation errors
func MapPageErr[T, V any](p *Page[T], mapF func(T) (V, error)) (*Page[V], error) {
	transformed, err := MapErr[T, V](p.Items, mapF)
	if err != nil {
		return nil, err
	}

	return &Page[V]{Items: transformed, Remaining: p.Remaining, ContinueToken: p.ContinueToken}, nil
}
