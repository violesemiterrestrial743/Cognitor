package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type RPCRule struct{}

func (RPCRule) ID() string { return "added-rpc-auth-validation" }

func (RPCRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(append(change.AddedCalls, change.AddedOps...), change.AddedStrings...),
		"RpcBindingInqAuthClient", "RpcBindingInqAuthClientEx", "RpcImpersonateClient", "RpcRevertToSelf",
		"RpcServerRegisterIf3", "RpcServerRegisterIfEx", "RPC_IF_ALLOW_SECURE_ONLY", "RPC_C_AUTHN_LEVEL",
		"RPC_C_AUTHZ", "rpc auth", "authentication level", "authorization service", "interface security callback",
	)
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "rpc-hardening", "Added RPC authentication validation", hits, 0.8*change.Similarity)}
}
