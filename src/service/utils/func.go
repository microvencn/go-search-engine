package utils

import (
	"strings"
	"unicode"
)

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

func RemoveSpace(str string) string {
	builder := strings.Builder{}

	for _, s := range strings.FieldsFunc(str, func(r rune) bool {
		if r <= 32 || unicode.IsSpace(r) {
			return true
		}
		return false
	}) {
		builder.WriteString(s)
	}

	return builder.String()
}
