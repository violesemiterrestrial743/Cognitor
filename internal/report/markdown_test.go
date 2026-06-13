package report

import (
	"strings"
	"testing"

	"github.com/kernelstub/cognitor/pkg/model"
)

func TestMarkdownIncludesRequiredSections(t *testing.T) {
	r := Build([]model.Finding{{
		ID:                    "f1",
		Title:                 "Added access validation before privileged behavior",
		AffectedBinary:        "driver.sys",
		OldFunction:           "Old",
		NewFunction:           "New",
		Category:              "access-control",
		Severity:              "high",
		Confidence:            0.9,
		RiskScore:             8.1,
		Evidence:              []string{"SeAccessCheck"},
		Reasoning:             "defensive hardening",
		SiblingBugSearchHints: []string{"Review sibling functions."},
	}}, model.Graph{})
	data, err := Markdown(r)
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, section := range []string{"Executive Summary", "Analyst Guidance", "Priority Review Queue", "Top Changed Components", "Top Findings", "Semantic Change Clusters", "Sibling Bug Hypotheses", "Recommended Manual Review Plan"} {
		if !strings.Contains(text, section) {
			t.Fatalf("missing section %q in %s", section, text)
		}
	}
}
