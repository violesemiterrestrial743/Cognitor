package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type NativeAPIRule struct{}

func (NativeAPIRule) ID() string { return "changed-native-api-boundary" }

func (NativeAPIRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(append(change.AddedCalls, change.AddedOps...), change.AddedStrings...),
		"Nt", "Zw", "syscall", "system call", "previous mode", "probe user", "capture user", "user mode buffer", "kernel mode caller",
	)
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "native-api-boundary", "Changed native API or syscall boundary validation", hits, 0.74*change.Similarity)}
}
