package graph

import (
	"sort"

	"github.com/kernelstub/cognitor/pkg/model"
)

func FunctionsNewlyProtected(findings []model.Finding) []string {
	seen := map[string]struct{}{}
	for _, finding := range findings {
		if finding.Category == "access-control" || finding.Category == "privilege-boundary" {
			seen[finding.AffectedBinary+"!"+finding.NewFunction] = struct{}{}
		}
	}
	return sortedKeys(seen)
}

func BinariesWithValidationAdditions(findings []model.Finding) []string {
	seen := map[string]struct{}{}
	for _, finding := range findings {
		seen[finding.AffectedBinary] = struct{}{}
	}
	return sortedKeys(seen)
}

func sortedKeys(seen map[string]struct{}) []string {
	values := make([]string, 0, len(seen))
	for value := range seen {
		values = append(values, value)
	}
	sort.Strings(values)
	return values
}
