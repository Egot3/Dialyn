package datafuncs

func SliceFromKeys[K comparable, T any](mp map[K]T) []K {
	var sl []K
	for val := range mp {
		sl = append(sl, val)
	}

	return sl
}
