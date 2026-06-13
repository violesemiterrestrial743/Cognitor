package model

type Finding struct {
	ID                        string   `json:"id"`
	Title                     string   `json:"title"`
	AffectedBinary            string   `json:"affected_binary"`
	OldFunction               string   `json:"old_function"`
	NewFunction               string   `json:"new_function"`
	Category                  string   `json:"category"`
	Confidence                float64  `json:"confidence"`
	Severity                  string   `json:"severity"`
	RiskScore                 float64  `json:"risk_score"`
	Evidence                  []string `json:"evidence"`
	OldEvidence               []string `json:"old_evidence"`
	NewEvidence               []string `json:"new_evidence"`
	Reasoning                 string   `json:"reasoning"`
	SiblingBugSearchHints     []string `json:"sibling_bug_search_hints"`
	RecommendedAuditTargets   []string `json:"recommended_audit_targets"`
	ResponsibleDisclosureNote string   `json:"responsible_disclosure_note"`
}
