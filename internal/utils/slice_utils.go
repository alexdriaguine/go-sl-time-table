package utils

func Map[T any, V any](collection []T, mapFn func(T) V) []V {
	mapped := make([]V, len(collection))

	for i, item := range collection {
		mapped[i] = mapFn(item)
	}

	return mapped
}

func Filter[T any](collection []T, filterFn func(T) bool) []T {
	filtered := []T{}

	for _, item := range collection {
		if filterFn(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
