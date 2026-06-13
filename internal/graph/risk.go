package graph

import "github.com/kernelstub/cognitor/pkg/model"

func SiblingPotential(findings []model.Finding) []model.Finding {
	var out []model.Finding
	for _, finding := range findings {
		if finding.Confidence >= 0.7 && finding.RiskScore >= 5 {
			out = append(out, finding)
		}
	}
	return out
}
