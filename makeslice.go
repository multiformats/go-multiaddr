package multiaddr

// makeSlice is like make([]T, len) but it perform a class size capacity
// extension. In other words, if the allocation gets rounded to a bigger
// allocation class, instead of wasting the unused space it is gonna return it
// as extra capacity.
func makeSlice[T any](len int) []T {
	return append([]T(nil), make([]T, len)...)
}
