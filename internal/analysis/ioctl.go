package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type IOCTLRule struct{}

func (IOCTLRule) ID() string { return "changed-ioctl-validation" }

func (IOCTLRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(change.AddedStrings, change.AddedOps...), "ioctl", "METHOD_NEITHER", "input buffer validation", "FILE_READ_DATA", "FILE_WRITE_DATA")
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "ioctl-hardening", "Changed IOCTL validation requirements", hits, 0.78*change.Similarity)}
}
