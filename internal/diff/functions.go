package diff

import (
	"math"
	"sort"

	"github.com/kernelstub/cognitor/pkg/model"
)

func MatchFunctions(oldBinary model.Binary, newBinary model.Binary) []model.FunctionPair {
	used := map[string]struct{}{}
	var pairs []model.FunctionPair
	for _, oldFn := range oldBinary.Functions {
		bestScore := 0.0
		var best model.Function
		reason := ""
		for _, newFn := range newBinary.Functions {
			if _, ok := used[newFn.ID]; ok {
				continue
			}
			score, why := similarity(oldFn, newFn)
			if score > bestScore {
				bestScore = score
				best = newFn
				reason = why
			}
		}
		if bestScore >= 0.45 {
			used[best.ID] = struct{}{}
			pairs = append(pairs, model.FunctionPair{Old: oldFn, New: best, Similarity: bestScore, Reason: reason})
		}
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Old.Name < pairs[j].Old.Name
	})
	return pairs
}

func similarity(a model.Function, b model.Function) (float64, string) {
	if a.Name != "" && a.Name == b.Name {
		return 1, "exact symbol match"
	}
	if a.NormalizedName != "" && a.NormalizedName == b.NormalizedName {
		return 0.95, "normalized symbol match"
	}
	score := 0.0
	reason := "semantic neighborhood match"
	score += 0.35 * jaccard(a.Imports, b.Imports)
	score += 0.35 * jaccard(a.Strings, b.Strings)
	score += 0.2 * jaccard(a.Calls, b.Calls)
	if a.BasicBlockCount > 0 && b.BasicBlockCount > 0 {
		delta := math.Abs(float64(a.BasicBlockCount - b.BasicBlockCount))
		base := math.Max(float64(a.BasicBlockCount), float64(b.BasicBlockCount))
		score += 0.1 * (1 - delta/base)
	}
	return score, reason
}

func jaccard(a []string, b []string) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 0
	}
	set := map[string]struct{}{}
	for _, value := range a {
		set[value] = struct{}{}
	}
	intersection := 0
	for _, value := range b {
		if _, ok := set[value]; ok {
			intersection++
		}
		set[value] = struct{}{}
	}
	if len(set) == 0 {
		return 0
	}
	return float64(intersection) / float64(len(set))
}
