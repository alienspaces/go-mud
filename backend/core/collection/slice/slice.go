package slice

// FromMap converts a map to a slice where the map values are the slice elements.
func FromMap[K comparable, V any](m map[K]V) []V {
	var s []V

	for _, v := range m {
		s = append(s, v)
	}

	return s
}

// Map maps the slice values using the mapFn.
func Map[T any, R any](s []T, mapFn func(T) R) []R {
	var r []R

	for _, t := range s {
		r = append(r, mapFn(t))
	}

	return r
}

func ToSliceOfPointers[T any](s []T) []*T {
	var ptrs []*T

	for i := range s {
		ptrs = append(ptrs, &s[i])
	}

	return ptrs
}
