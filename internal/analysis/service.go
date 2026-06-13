package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type ServiceRule struct{}

func (ServiceRule) ID() string { return "added-service-hardening" }

func (ServiceRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(change.AddedCalls, change.AddedOps...), "ChangeServiceConfig2", "service permission", "service ACL", "SERVICE_CHANGE_CONFIG")
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "service-hardening", "Added service permission hardening", hits, 0.78*change.Similarity)}
}
