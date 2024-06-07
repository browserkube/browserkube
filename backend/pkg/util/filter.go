package browserkubeutil

func Filter[T any](items []T, filter func(T) bool) []T {
	var filtered []T
	for _, item := range items {
		if filter(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func FindFirst[T any](items []T, filter func(T) bool) T {
	for _, item := range items {
		if filter(item) {
			return item
		}
	}
	return *new(T)
}

// ReverseDeduplicateBY creates an array of deduplicated items
func ReverseDeduplicateBY[T any, V comparable](items []T, uniqF func(T) V) []T {
	result := make([]T, 0, len(items))
	seen := make(map[V]struct{}, len(items))

	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]
		itemHash := uniqF(item)
		if _, ok := seen[itemHash]; ok {
			continue
		}

		seen[itemHash] = struct{}{}
		result = append(result, item)
	}

	return result
}
