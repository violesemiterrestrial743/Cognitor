package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type BoundsCheckRule struct{}

func (BoundsCheckRule) ID() string { return "added-bounds-check" }

func (BoundsCheckRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(change.AddedCalls, change.AddedOps...), "ProbeForRead", "ProbeForWrite", "length check", "bounds check", "integer overflow", "size validation", "null check")
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "memory-safety", "Added input or memory safety validation", hits, 0.86*change.Similarity)}
}
