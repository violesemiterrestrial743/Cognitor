package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type AccessCheckRule struct{}

func (AccessCheckRule) ID() string { return "added-access-check" }

func (AccessCheckRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(change.AddedCalls, change.AddedOps...), "AccessCheck", "SeAccessCheck", "NtAccessCheck", "access mask", "object type validation")
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "access-control", "Added access validation before privileged behavior", hits, 0.9*change.Similarity)}
}
