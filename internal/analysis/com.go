package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type COMRule struct{}

func (COMRule) ID() string { return "changed-com-permissions" }

func (COMRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(append(change.AddedCalls, change.AddedStrings...), change.AddedOps...),
		"COM elevation", "LaunchPermission", "AccessPermission", "AppID", "CLSID", "CoInitializeSecurity",
		"CoImpersonateClient", "CoRevertToSelf", "EOAC_DISABLE_AAA", "RPC_C_IMP_LEVEL", "DCOM hardening",
	)
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "com-hardening", "Changed COM elevation or launch permission validation", hits, 0.76*change.Similarity)}
}
