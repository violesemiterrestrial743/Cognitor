package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type HandleValidationRule struct{}

func (HandleValidationRule) ID() string { return "added-handle-object-validation" }

func (HandleValidationRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(change.AddedCalls, change.AddedOps...),
		"ObReferenceObjectByHandle", "ObOpenObjectByPointer", "ObGetObjectType", "handle type validation", "object type validation",
		"GrantedAccess", "DesiredAccess", "access mask validation", "OBJ_KERNEL_HANDLE",
	)
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "handle-validation", "Added handle or object type validation", hits, 0.84*change.Similarity)}
}
