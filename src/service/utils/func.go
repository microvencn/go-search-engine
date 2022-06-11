package utils

func MinInt(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func MaxInt(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
