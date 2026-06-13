package report

import (
	"bytes"
	"encoding/csv"
	"fmt"

	"github.com/kernelstub/cognitor/pkg/model"
)

func CSV(report model.Report) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"type", "target", "category", "severity", "confidence", "risk_score", "reason", "signals"}); err != nil {
		return nil, err
	}
	for _, finding := range report.Findings {
		if err := writer.Write([]string{
			"finding",
			finding.AffectedBinary,
			finding.Category,
			finding.Severity,
			fmt.Sprintf("%.2f", finding.Confidence),
			fmt.Sprintf("%.1f", finding.RiskScore),
			finding.Title,
			joinCSVSignals(finding.Evidence),
		}); err != nil {
			return nil, err
		}
	}
	for _, change := range report.Changes.ModifiedBinaries {
		if err := writer.Write([]string{
			"binary-change",
			change.Path,
			change.ChangeClass,
			priorityFromSignals(change.RiskSignals),
			"",
			"",
			fmt.Sprintf("size %d -> %d", change.OldSize, change.NewSize),
			joinCSVSignals(change.RiskSignals),
		}); err != nil {
			return nil, err
		}
	}
	for _, change := range report.Changes.ChangedArtifacts {
		if err := writer.Write([]string{
			"artifact-change",
			change.Path,
			change.ChangeClass,
			priorityFromSignals(change.RiskSignals),
			"",
			"",
			fmt.Sprintf("size %d -> %d", change.OldSize, change.NewSize),
			joinCSVSignals(change.RiskSignals),
		}); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return buf.Bytes(), writer.Error()
}

func joinCSVSignals(values []string) string {
	var out string
	for i, value := range values {
		if i > 0 {
			out += "; "
		}
		out += value
	}
	return out
}
