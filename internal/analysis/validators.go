package analysis

import (
	"strings"

	"github.com/kernelstub/cognitor/internal/util"
	"github.com/kernelstub/cognitor/pkg/model"
)

func finding(change model.SemanticChange, category string, title string, evidence []string, confidence float64) model.Finding {
	return model.Finding{
		ID:             util.StableID(change.Binary.Path, change.NewFunction.Name, category, strings.Join(evidence, ",")),
		Title:          title,
		AffectedBinary: change.Binary.Path,
		OldFunction:    change.OldFunction.Name,
		NewFunction:    change.NewFunction.Name,
		Category:       category,
		Confidence:     confidence,
		Evidence:       evidence,
		OldEvidence:    append(change.OldFunction.Calls, change.OldFunction.Operations...),
		NewEvidence:    append(change.NewFunction.Calls, change.NewFunction.Operations...),
		Reasoning:      "New defensive validation or authorization behavior appears before or near a sensitive operation in the matched function.",
		SiblingBugSearchHints: []string{
			"Review sibling functions with the same operation pattern but without the new validation.",
			"Compare callers that reach the same privileged API or object type.",
		},
		RecommendedAuditTargets: []string{
			change.Binary.Path,
			change.NewFunction.Name,
		},
		ResponsibleDisclosureNote: "Use findings for defensive review and coordinated disclosure. Do not publish exploit details or weaponized proof of concept material.",
	}
}

func hasAny(values []string, needles ...string) []string {
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
