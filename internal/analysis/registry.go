package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type RegistryRule struct{}

func (RegistryRule) ID() string { return "added-registry-hardening" }

func (RegistryRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(change.AddedCalls, change.AddedOps...), "RegSetKeySecurity", "registry ACL", "CmCheckRegistryAccess", "KEY_SET_VALUE")
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "registry-hardening", "Added registry permission hardening", hits, 0.78*change.Similarity)}
}
