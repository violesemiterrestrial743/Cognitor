package diff

import "strings"

func ContainsAny(values []string, needles ...string) []string {
	var hits []string
	for _, value := range values {
		lower := strings.ToLower(value)
		for _, needle := range needles {
			if strings.Contains(lower, strings.ToLower(needle)) {
				hits = append(hits, value)
				break
			}
		}
	}
	return hits
}
