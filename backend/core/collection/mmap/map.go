package mmap

// FromSlice converts a slice to a map where the slice elements are the map values,
// and the corresponding keys are computed by the keyFn.
func FromSlice[K comparable, V any](keyFn func(V) K, slices ...[]V) map[K]V {
	m := map[K]V{}

	for _, items := range slices {
		for _, i := range items {
			m[keyFn(i)] = i
		}
	}

	return m
}
