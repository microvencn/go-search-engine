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

// Intersection 计算交集大小，要求 source 和 target 均有序
func Intersection(source []string, target []string) int {
	i, j := 0, 0
	var count int = 0
	for {
		for i < len(target) && target[i] < source[j] {
			i++
		}
		if i == len(target) {
			break
		}

		for j < len(source) && source[j] < target[i] {
			j++
		}
		if j == len(source) {
			break
		}

		if target[i] == source[j] {
			count++
			i++
			j++
			if i == len(target) || j == len(source) {
				break
			}
		}
	}
	return count
}

// HasIntersection 判断是否有交集，要求 source 和 target 均有序
func HasIntersection(source []string, target []string) bool {
	i, j := 0, 0
	for {
		for i < len(target) && target[i] < source[j] {
			i++
		}
		if i == len(target) {
			break
		}

		for j < len(source) && source[j] < target[i] {
			j++
		}
		if j == len(source) {
			break
		}

		if target[i] == source[j] {
			return true
		}
	}
	return false
}
