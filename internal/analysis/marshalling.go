package analysis

import (
	"context"

	"github.com/kernelstub/cognitor/pkg/model"
)

type MarshallingRule struct{}

func (MarshallingRule) ID() string { return "added-marshalling-validation" }

func (MarshallingRule) Evaluate(ctx context.Context, change model.SemanticChange) []model.Finding {
	hits := hasAny(append(append(change.AddedCalls, change.AddedOps...), change.AddedStrings...),
		"Ndr", "MIDL", "marshalling", "unmarshal", "deserialize", "wire length", "conformant array",
		"range check", "max count", "string binding validation",
	)
	if len(hits) == 0 {
		return nil
	}
	return []model.Finding{finding(change, "marshalling-validation", "Added RPC or structured input marshalling validation", hits, 0.79*change.Similarity)}
}
