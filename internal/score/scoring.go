package score

import (
	"strings"

	"github.com/kernelstub/cognitor/pkg/model"
)

type Scorer struct{}

func DefaultScorer() Scorer {
	return Scorer{}
}

func (Scorer) Score(finding model.Finding) model.Finding {
	score := 2.0 + confidenceBand(finding.Confidence)*3
	text := strings.ToLower(finding.AffectedBinary + " " + finding.Category + " " + strings.Join(finding.Evidence, " "))
	if strings.Contains(text, ".sys") || strings.Contains(text, "kernel") || strings.Contains(text, "driver") {
		score += 1.4
	}
	if strings.Contains(text, "rpc") || strings.Contains(text, "alpc") || strings.Contains(text, "service") {
		score += 1.1
	}
	if strings.Contains(text, "ioctl") || strings.Contains(text, "probefor") || strings.Contains(text, "user-to-kernel") {
		score += 1.3
	}
	if strings.Contains(text, "access") || strings.Contains(text, "privilege") || strings.Contains(text, "token") {
		score += 1.2
	}
	if strings.Contains(text, "network") {
		score += 0.8
	}
	if score > 10 {
		score = 10
	}
	finding.RiskScore = score
	finding.Severity = severity(score)
	return finding
}
