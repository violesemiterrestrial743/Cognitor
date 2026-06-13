package ingest

import (
	"os"
	"sort"
	"unicode"
)

func ExtractStrings(path string, minLen int) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	var current []rune
	flush := func() {
		if len(current) >= minLen {
			seen[string(current)] = struct{}{}
		}
		current = current[:0]
	}
	for _, r := range string(data) {
		if r < 128 && (unicode.IsPrint(r) || r == '\t') {
			current = append(current, r)
			continue
		}
		flush()
	}
	flush()
	values := make([]string, 0, len(seen))
	for value := range seen {
		values = append(values, value)
	}
	sort.Strings(values)
	if len(values) > 256 {
		values = values[:256]
	}
	return values, nil
}
