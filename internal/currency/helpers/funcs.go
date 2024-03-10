package helpers

func GetUnique[T comparable](sl []T) []T {
	m := make(map[T]struct{})

	var res []T
	for i := range sl {
		if _, ok := m[sl[i]]; !ok {
			m[sl[i]] = struct{}{}
			res = append(res, sl[i])
		}
	}

	return res
}