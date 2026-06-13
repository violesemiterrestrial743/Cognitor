package ingest

import (
	"strings"
	"unicode"
)

func NormalizeSymbol(value string) string {
	value = strings.ToLower(value)
	value = strings.TrimPrefix(value, "_")
	value = strings.TrimSuffix(value, "w")
	value = strings.TrimSuffix(value, "a")
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return r
		}
		return -1
	}, value)
}
