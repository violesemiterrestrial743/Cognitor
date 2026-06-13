package model

type SemanticChange struct {
	Binary       Binary
	OldFunction  Function
	NewFunction  Function
	AddedCalls   []string
	AddedStrings []string
	AddedOps     []string
	Similarity   float64
	MatchReason  string
}
