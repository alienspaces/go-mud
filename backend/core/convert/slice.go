package convert

func GenericSlice[T any](s []T) []any {
	var g []any

	for _, t := range s {
		g = append(g, t)
	}

	return g
}
