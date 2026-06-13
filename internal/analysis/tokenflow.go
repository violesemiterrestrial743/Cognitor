package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type TokenFlowRule struct{}

func (TokenFlowRule) ID() string { return "added-token-check" }

func (TokenFlowRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(change.AddedCalls, change.AddedOps...), "SeSinglePrivilegeCheck", "privilege check", "token privilege", "impersonation")
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "privilege-boundary", "Added token or impersonation validation", hits, 0.88*change.Similarity)}
}
