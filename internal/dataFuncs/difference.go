package datafuncs

// first - second
func Difference[T comparable](sl1, sl2 []T) []T {
	mp1 := make(map[T]int, len(sl1))
	for _, val := range sl1 {
		mp1[val] = 0
	}
	mp2 := make(map[T]int, len(sl2))
	for _, val := range sl2 {
		mp2[val] = 0
	}

	for val := range mp2 {
		delete(mp1, val)
	}

	var diff []T
	for val := range mp1 {
		diff = append(diff, val)
	}

	return diff
}
