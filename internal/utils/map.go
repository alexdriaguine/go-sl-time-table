package utils

func Map[T any, V any](collection []T, mapFn func(T) V) []V {
	mapped := make([]V, len(collection))

	for i, m := range collection {
		mapped[i] = mapFn(m)
	}

	return mapped
}
