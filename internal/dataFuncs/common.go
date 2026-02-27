package datafuncs

func Common[T comparable](sl1, sl2 []T) []T {
	mp1 := make(map[T]int, len(sl1))
	for _, val := range sl1 {
		mp1[val] = 0
	}

	commonMp := make(map[T]int, min(len(sl1), len(sl2)))
	var common []T
	for _, val := range sl2 {
		if _, exists := mp1[val]; exists {
			if _, added := commonMp[val]; !added {
				commonMp[val] = 0
				common = append(common, val)
			}
		}
	}
	return common
}
