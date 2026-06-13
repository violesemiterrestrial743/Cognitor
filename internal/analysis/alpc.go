package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type ALPCRule struct{}

func (ALPCRule) ID() string { return "added-alpc-validation" }

func (ALPCRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(change.AddedCalls, change.AddedOps...), "Alpc", "port security", "message attribute validation")
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "alpc-hardening", "Added ALPC validation", hits, 0.76*change.Similarity)}
}
