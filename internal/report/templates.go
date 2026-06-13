package report

import (
	"sort"

	"github.com/kernelstub/cognitor/internal/util"
	"github.com/kernelstub/cognitor/pkg/model"
)

func Build(findings []model.Finding, graph model.Graph, changes ...model.ChangeSummary) model.Report {
	summary := model.ReportSummary{
		TotalFindings: len(findings),
		BySeverity:    map[string]int{},
		ByCategory:    map[string]int{},
	}
	componentCounts := map[string]int{}
	for _, finding := range findings {
		summary.BySeverity[finding.Severity]++
		summary.ByCategory[finding.Category]++
		componentCounts[finding.AffectedBinary]++
	}
	var components []string
	for component := range componentCounts {
		components = append(components, component)
	}
	sort.Slice(components, func(i, j int) bool {
		if componentCounts[components[i]] == componentCounts[components[j]] {
			return components[i] < components[j]
		}
		return componentCounts[components[i]] > componentCounts[components[j]]
	})
	if len(components) > 10 {
		components = components[:10]
	}
	summary.TopChangedComponents = components
	var changeSummary model.ChangeSummary
	if len(changes) > 0 {
		changeSummary = changes[0]
	}
	metadata := model.ReportMetadata{GeneratedAt: util.NowUTC().Format("2006-01-02T15:04:05Z07:00")}
	return model.Report{Metadata: metadata, Summary: summary, Executive: buildExecutive(findings, changeSummary), Changes: changeSummary, Findings: findings, Graph: graph}
}

func BuildWithMetadata(findings []model.Finding, graph model.Graph, metadata model.ReportMetadata, changes model.ChangeSummary) model.Report {
	report := Build(findings, graph, changes)
	if metadata.GeneratedAt == "" {
		metadata.GeneratedAt = report.Metadata.GeneratedAt
	}
	report.Metadata = metadata
	return report
}

func buildExecutive(findings []model.Finding, changes model.ChangeSummary) model.ReportExecutive {
	totalChanged := len(changes.AddedBinaries) + len(changes.RemovedBinaries) + len(changes.ModifiedBinaries) +
		len(changes.AddedArtifacts) + len(changes.RemovedArtifacts) + len(changes.ChangedArtifacts) +
		len(changes.AddedServices) + len(changes.RemovedServices) + len(changes.ChangedServices) +
		len(changes.AddedRegistry) + len(changes.RemovedRegistry) + len(changes.ChangedRegistry)
	highConfidence := 0
	maxScore := 0.0
	bySeverity := map[string]int{}
	for _, finding := range findings {
		if finding.Confidence >= 0.8 {
			highConfidence++
		}
		if finding.RiskScore > maxScore {
			maxScore = finding.RiskScore
		}
		bySeverity[finding.Severity]++
	}
	targets := topReviewTargets(findings, changes)
	riskLevel := riskLevel(maxScore, bySeverity, totalChanged, targets)
	return model.ReportExecutive{
		RiskLevel:              riskLevel,
		Priority:               priorityForRisk(riskLevel),
		TotalChangedItems:      totalChanged,
		HighConfidenceFindings: highConfidence,
		ChangedArtifactKinds:   changedArtifactKinds(changes),
		TopReviewTargets:       targets,
		BeginnerNotes:          beginnerNotes(riskLevel, len(findings), totalChanged),
		ResearcherChecklist:    researcherChecklist(findings, changes),
		NextActions:            nextActions(riskLevel, len(findings), totalChanged),
	}
}

func riskLevel(maxScore float64, bySeverity map[string]int, totalChanged int, targets []model.ReviewTarget) string {
	switch {
	case bySeverity["high"] > 0 || maxScore >= 8:
		return "high"
	case bySeverity["medium"] >= 3 || len(targets) >= 5:
		return "elevated"
	case totalChanged > 0 || bySeverity["medium"] > 0:
		return "moderate"
	default:
		return "informational"
	}
}

func priorityForRisk(risk string) string {
	switch risk {
	case "high":
		return "review immediately"
	case "elevated":
		return "same-day review"
	case "moderate":
		return "scheduled review"
	default:
		return "archive or spot-check"
	}
}

func topReviewTargets(findings []model.Finding, changes model.ChangeSummary) []model.ReviewTarget {
	var targets []model.ReviewTarget
	for _, finding := range findings {
		targets = append(targets, model.ReviewTarget{
			Target:   finding.AffectedBinary,
			Kind:     "finding",
			Reason:   finding.Title,
			Priority: finding.Severity,
			Signals:  finding.Evidence,
		})
	}
	for _, change := range changes.ModifiedBinaries {
		if len(change.RiskSignals) == 0 {
			continue
		}
		targets = append(targets, model.ReviewTarget{
			Target:   change.Path,
			Kind:     "binary",
			Reason:   change.ChangeClass,
			Priority: priorityFromSignals(change.RiskSignals),
			Signals:  change.RiskSignals,
		})
	}
	for _, change := range changes.ChangedArtifacts {
		if len(change.RiskSignals) == 0 {
			continue
		}
		targets = append(targets, model.ReviewTarget{
			Target:   change.Path,
			Kind:     "artifact",
			Reason:   change.ChangeClass,
			Priority: priorityFromSignals(change.RiskSignals),
			Signals:  change.RiskSignals,
		})
	}
	sort.Slice(targets, func(i, j int) bool {
		pi, pj := priorityRank(targets[i].Priority), priorityRank(targets[j].Priority)
		if pi == pj {
			return targets[i].Target < targets[j].Target
		}
		return pi < pj
	})
	if len(targets) > 10 {
		targets = targets[:10]
	}
	return targets
}

func priorityFromSignals(signals []string) string {
	for _, signal := range signals {
		switch signal {
		case "authorization", "token-or-identity", "kernel-boundary", "memory-safety":
			return "high"
		}
	}
	if len(signals) > 0 {
		return "medium"
	}
	return "low"
}

func priorityRank(priority string) int {
	switch priority {
	case "high":
		return 0
	case "medium":
		return 1
	case "low":
		return 2
	default:
		return 3
	}
}

func changedArtifactKinds(changes model.ChangeSummary) []string {
	seen := map[string]struct{}{}
	for _, change := range append(append([]model.ArtifactChange{}, changes.AddedArtifacts...), changes.ChangedArtifacts...) {
		if change.Kind != "" {
			seen[change.Kind] = struct{}{}
		}
	}
	var kinds []string
	for kind := range seen {
		kinds = append(kinds, kind)
	}
	sort.Strings(kinds)
	return kinds
}

func nextActions(risk string, findingCount int, changeCount int) []string {
	if findingCount == 0 && changeCount == 0 {
		return []string{"Archive the run as a clean comparison baseline.", "Spot-check snapshot completeness if this result was unexpected."}
	}
	actions := []string{
		"Review top targets before broad inventory triage.",
		"Confirm whether added validation dominates sensitive operations in real control flow.",
		"Check sibling components that share the same API, object, service, registry, or artifact signals.",
	}
	if risk == "high" || risk == "elevated" {
		actions = append(actions, "Prepare a concise validation note for coordinated defensive follow-up.")
	}
	return actions
}

func beginnerNotes(risk string, findingCount int, changeCount int) []string {
	if findingCount == 0 && changeCount == 0 {
		return []string{
			"No security-relevant rule matched, and no tracked inventory drift was found.",
			"This does not prove two snapshots are identical; it means Cognitor did not find supported signals in the scanned inputs.",
		}
	}
	return []string{
		"Start with the Priority Review Queue instead of reading every changed string.",
		"A finding means Cognitor saw newly added defensive logic, not that an exploitable vulnerability is proven.",
		"Risk posture combines rule findings, confidence, and changed inventory signals; use it to choose review order.",
		"Compare old and new function sidecars when available to confirm the check is actually on the sensitive path.",
		"Current risk posture: " + risk + ". Changed tracked items: " + intString(changeCount) + ".",
	}
}

func researcherChecklist(findings []model.Finding, changes model.ChangeSummary) []string {
	checklist := []string{
		"Identify the exact caller-controlled inputs that reach each changed function.",
		"Check whether new validation dominates the privileged operation, allocation, copy, object lookup, RPC dispatch, or COM activation path.",
		"Look for sibling APIs or protocol methods with the same operation but without the new guard.",
		"Compare public symbols, exports, imports, and newly added strings for hints about feature flags or policy rollout.",
	}
	if hasCategory(findings, "native-api-boundary") {
		checklist = append(checklist, "For native API changes, verify PreviousMode, probing/capture, handle access, and user/kernel buffer assumptions.")
	}
	if hasCategory(findings, "rpc-hardening") || hasCategory(findings, "marshalling-validation") {
		checklist = append(checklist, "For RPC changes, review interface security callbacks, authn/authz levels, impersonation windows, and NDR range checks.")
	}
	if hasCategory(findings, "com-hardening") {
		checklist = append(checklist, "For COM changes, compare LaunchPermission, AccessPermission, AppID/CLSID mapping, and impersonation level changes.")
	}
	if len(changes.ChangedArtifacts)+len(changes.AddedArtifacts) > 0 {
		checklist = append(checklist, "For changed artifacts, correlate policy/configuration strings with the binaries that consume them.")
	}
	return checklist
}

func hasCategory(findings []model.Finding, category string) bool {
	for _, finding := range findings {
		if finding.Category == category {
			return true
		}
	}
	return false
}

func intString(value int) string {
	if value == 0 {
		return "0"
	}
	var digits [20]byte
	i := len(digits)
	for value > 0 {
		i--
		digits[i] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[i:])
}
