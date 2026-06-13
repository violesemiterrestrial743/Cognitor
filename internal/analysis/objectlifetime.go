package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type ObjectLifetimeRule struct{}

func (ObjectLifetimeRule) ID() string { return "added-lifetime-reference" }

func (ObjectLifetimeRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(change.AddedCalls, change.AddedOps...),
		"ObReferenceObject", "ObReferenceObjectByHandle", "ObDereferenceObject", "ExAcquireRundownProtection",
		"ExReleaseRundownProtection", "reference count", "lifetime reference", "rundown protection", "use-after-free guard",
	)
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "object-lifetime", "Added object lifetime handling", hits, 0.82*change.Similarity)}
}
