package report

import (
	"encoding/json"

	"github.com/kernelstub/cognitor/pkg/model"
)

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name  string      `json:"name"`
	Rules []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	ShortDescription sarifText `json:"shortDescription"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifText       `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifText struct {
	Text string `json:"text"`
}

func SARIF(report model.Report) ([]byte, error) {
	rulesByID := map[string]sarifRule{}
	var results []sarifResult
	for _, finding := range report.Findings {
		rulesByID[finding.Category] = sarifRule{ID: finding.Category, Name: finding.Category, ShortDescription: sarifText{Text: finding.Title}}
		results = append(results, sarifResult{
			RuleID:    finding.Category,
			Level:     sarifLevel(finding.Severity),
			Message:   sarifText{Text: finding.Title + ": " + finding.Reasoning},
			Locations: []sarifLocation{{PhysicalLocation: sarifPhysicalLocation{ArtifactLocation: sarifArtifactLocation{URI: finding.AffectedBinary}}}},
		})
	}
	var rules []sarifRule
	for _, rule := range rulesByID {
		rules = append(rules, rule)
	}
	log := sarifLog{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Runs: []sarifRun{{
			Tool:    sarifTool{Driver: sarifDriver{Name: "cognitor", Rules: rules}},
			Results: results,
		}},
	}
	return json.MarshalIndent(log, "", "  ")
}

func sarifLevel(severity string) string {
	switch severity {
	case "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "note"
	}
}
