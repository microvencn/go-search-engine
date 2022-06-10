package utils

func RemoveRepeatedElement[T int | string](s []T) []T {
	result := make([]T, 0)
	m := make(map[T]bool)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

type HasValue interface {
	Value() string
}

func RemoveRepeated[T HasValue](s []T) []T {
	result := make([]T, len(s))
	m := make(map[string]bool)
	for _, v := range s {
		if _, ok := m[v.Value()]; !ok {
			result = append(result, v)
			m[v.Value()] = true
		}
	}
	return result
}
