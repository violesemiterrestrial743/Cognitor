package diff

import (
	"context"
	"sort"

	"github.com/kernelstub/cognitor/internal/analysis"
	"github.com/kernelstub/cognitor/internal/score"
	"github.com/kernelstub/cognitor/pkg/model"
)

func Analyze(ctx context.Context, oldSnapshot model.Snapshot, newSnapshot model.Snapshot) []model.Finding {
	var changes []SemanticChange
	for _, pair := range MatchBinaries(oldSnapshot, newSnapshot) {
		oldBinary, newBinary := pair[0], pair[1]
		for _, fnPair := range MatchFunctions(oldBinary, newBinary) {
			change := SemanticChange{
				Binary:       newBinary,
				OldFunction:  fnPair.Old,
				NewFunction:  fnPair.New,
				AddedCalls:   AddedStrings(fnPair.Old.Calls, fnPair.New.Calls),
				AddedStrings: AddedStrings(fnPair.Old.Strings, fnPair.New.Strings),
				AddedOps:     AddedStrings(fnPair.Old.Operations, fnPair.New.Operations),
				Similarity:   fnPair.Similarity,
				MatchReason:  fnPair.Reason,
			}
			if len(change.AddedCalls)+len(change.AddedStrings)+len(change.AddedOps) > 0 {
				changes = append(changes, change)
			}
		}
	}
	findings := analysis.DefaultEngine().Evaluate(ctx, changes)
	scorer := score.DefaultScorer()
	for i := range findings {
		findings[i] = scorer.Score(findings[i])
	}
	sort.Slice(findings, func(i, j int) bool {
		if findings[i].RiskScore == findings[j].RiskScore {
			return findings[i].ID < findings[j].ID
		}
		return findings[i].RiskScore > findings[j].RiskScore
	})
	return findings
}
