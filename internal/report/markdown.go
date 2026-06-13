package report

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/kernelstub/cognitor/pkg/model"
)

func Markdown(report model.Report) ([]byte, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "# Cognitor Defensive Semantic Diff Report\n\n")
	if report.Metadata.GeneratedAt != "" || report.Metadata.ToolVersion != "" {
		fmt.Fprintf(&buf, "- Generated: `%s`\n", valueOrUnknown(report.Metadata.GeneratedAt))
		fmt.Fprintf(&buf, "- Tool version: `%s`\n", valueOrUnknown(report.Metadata.ToolVersion))
		if report.Metadata.OldSnapshot.Path != "" || report.Metadata.NewSnapshot.Path != "" {
			fmt.Fprintf(&buf, "- Snapshots: `%s` -> `%s`\n", valueOrUnknown(report.Metadata.OldSnapshot.Path), valueOrUnknown(report.Metadata.NewSnapshot.Path))
		}
		fmt.Fprintf(&buf, "\n")
	}
	fmt.Fprintf(&buf, "## Executive Summary\n\n")
	fmt.Fprintf(&buf, "Cognitor identified %d security-relevant semantic change candidates. These findings are intended for defensive validation, patch comprehension, sibling-bug hunting at a high level, and coordinated disclosure workflows.\n\n", report.Summary.TotalFindings)
	fmt.Fprintf(&buf, "Risk posture: `%s`; priority: `%s`; changed items: `%d`; high-confidence findings: `%d`.\n\n", report.Executive.RiskLevel, report.Executive.Priority, report.Executive.TotalChangedItems, report.Executive.HighConfidenceFindings)
	fmt.Fprintf(&buf, "Change inventory: %d added binaries, %d removed binaries, %d modified binaries, %d added artifacts, %d removed artifacts, %d changed artifacts, %d service changes, and %d registry changes.\n\n",
		len(report.Changes.AddedBinaries),
		len(report.Changes.RemovedBinaries),
		len(report.Changes.ModifiedBinaries),
		len(report.Changes.AddedArtifacts),
		len(report.Changes.RemovedArtifacts),
		len(report.Changes.ChangedArtifacts),
		len(report.Changes.AddedServices)+len(report.Changes.RemovedServices)+len(report.Changes.ChangedServices),
		len(report.Changes.AddedRegistry)+len(report.Changes.RemovedRegistry)+len(report.Changes.ChangedRegistry),
	)
	writeCounts(&buf, "Severity", report.Summary.BySeverity)
	writeCounts(&buf, "Likely Vulnerability Classes", report.Summary.ByCategory)
	writeGuidance(&buf, report.Executive)
	writeReviewQueue(&buf, report.Executive)
	fmt.Fprintf(&buf, "## Automatic Change Inventory\n\n")
	writeBinaryChanges(&buf, "Modified Binaries", report.Changes.ModifiedBinaries, 20)
	writeBinaryChanges(&buf, "Added Binaries", report.Changes.AddedBinaries, 15)
	writeBinaryChanges(&buf, "Removed Binaries", report.Changes.RemovedBinaries, 15)
	writeArtifactChanges(&buf, "Changed Evidence Artifacts", report.Changes.ChangedArtifacts, 20)
	writeArtifactChanges(&buf, "Added Evidence Artifacts", report.Changes.AddedArtifacts, 15)
	writeArtifactChanges(&buf, "Removed Evidence Artifacts", report.Changes.RemovedArtifacts, 15)
	writeServiceRegistryChanges(&buf, report.Changes)
	fmt.Fprintf(&buf, "## Top Changed Components\n\n")
	if len(report.Summary.TopChangedComponents) == 0 {
		fmt.Fprintf(&buf, "No changed components with findings were recorded.\n\n")
	} else {
		for _, component := range report.Summary.TopChangedComponents {
			fmt.Fprintf(&buf, "- `%s`\n", component)
		}
		fmt.Fprintf(&buf, "\n")
	}
	fmt.Fprintf(&buf, "## Top Findings\n\n")
	if len(report.Findings) == 0 {
		fmt.Fprintf(&buf, "No findings were generated.\n\n")
	} else {
		limit := len(report.Findings)
		if limit > 20 {
			limit = 20
		}
		for i := 0; i < limit; i++ {
			f := report.Findings[i]
			fmt.Fprintf(&buf, "### %s\n\n", f.Title)
			fmt.Fprintf(&buf, "- ID: `%s`\n", f.ID)
			fmt.Fprintf(&buf, "- Component: `%s`\n", f.AffectedBinary)
			fmt.Fprintf(&buf, "- Function: `%s` -> `%s`\n", f.OldFunction, f.NewFunction)
			fmt.Fprintf(&buf, "- Category: `%s`\n", f.Category)
			fmt.Fprintf(&buf, "- Severity: `%s`\n", f.Severity)
			fmt.Fprintf(&buf, "- Confidence: `%.2f`\n", f.Confidence)
			fmt.Fprintf(&buf, "- Risk score: `%.1f`\n", f.RiskScore)
			fmt.Fprintf(&buf, "- Evidence: %s\n", inlineList(f.Evidence))
			fmt.Fprintf(&buf, "- Reasoning: %s\n\n", f.Reasoning)
		}
	}
	fmt.Fprintf(&buf, "## Semantic Change Clusters\n\n")
	writeCounts(&buf, "Categories", report.Summary.ByCategory)
	fmt.Fprintf(&buf, "## Sibling Bug Hypotheses\n\n")
	for _, f := range report.Findings {
		if len(f.SiblingBugSearchHints) == 0 {
			continue
		}
		fmt.Fprintf(&buf, "- `%s`: %s\n", f.AffectedBinary, strings.Join(f.SiblingBugSearchHints, " "))
	}
	if len(report.Findings) == 0 {
		fmt.Fprintf(&buf, "No sibling-bug hypotheses were generated.\n")
	}
	fmt.Fprintf(&buf, "\n## Recommended Manual Review Plan\n\n")
	fmt.Fprintf(&buf, "1. Validate whether each newly added check dominates the sensitive operation in real control flow.\n")
	fmt.Fprintf(&buf, "2. Inspect sibling functions and callers that still reach the same privileged API, object type, IOCTL, registry key, service, RPC interface, COM class, or ALPC port.\n")
	fmt.Fprintf(&buf, "3. Review modified EDB, event, registry, service, manifest, and configuration artifacts for newly introduced policy gates, defaults, ACLs, identifiers, or telemetry strings.\n")
	fmt.Fprintf(&buf, "4. Confirm reachability from relevant trust boundaries such as user-to-kernel, service, network, AppContainer, or low-privilege local callers.\n")
	fmt.Fprintf(&buf, "5. Record conclusions in a responsible disclosure workflow without exploit payloads or bypass instructions.\n")
	return buf.Bytes(), nil
}

func writeCounts(buf *bytes.Buffer, title string, counts map[string]int) {
	fmt.Fprintf(buf, "## %s\n\n", title)
	if len(counts) == 0 {
		fmt.Fprintf(buf, "No data.\n\n")
		return
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(buf, "- `%s`: %d\n", key, counts[key])
	}
	fmt.Fprintf(buf, "\n")
}

func inlineList(values []string) string {
	if len(values) == 0 {
		return "none"
	}
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, "`"+value+"`")
	}
	return strings.Join(quoted, ", ")
}

func writeBinaryChanges(buf *bytes.Buffer, title string, changes []model.BinaryChange, limit int) {
	fmt.Fprintf(buf, "### %s\n\n", title)
	if len(changes) == 0 {
		fmt.Fprintf(buf, "No entries.\n\n")
		return
	}
	for i, change := range changes {
		if i >= limit {
			fmt.Fprintf(buf, "- ... %d more entries omitted from this section.\n\n", len(changes)-limit)
			return
		}
		fmt.Fprintf(buf, "- `%s`", change.Path)
		if change.Kind != "" {
			fmt.Fprintf(buf, " (%s)", change.Kind)
		}
		if change.OldSize > 0 || change.NewSize > 0 {
			fmt.Fprintf(buf, " size `%d` -> `%d`", change.OldSize, change.NewSize)
			if change.SizeDelta != 0 {
				fmt.Fprintf(buf, " (delta `%+d`)", change.SizeDelta)
			}
		}
		if change.ChangeClass != "" {
			fmt.Fprintf(buf, " class `%s`", change.ChangeClass)
		}
		if change.OldVersion != "" || change.NewVersion != "" {
			fmt.Fprintf(buf, " version `%s` -> `%s`", change.OldVersion, change.NewVersion)
		}
		fmt.Fprintf(buf, "\n")
		if len(change.RiskSignals) > 0 {
			fmt.Fprintf(buf, "  - Signals: %s\n", inlineList(change.RiskSignals))
		}
		if len(change.AddedImports) > 0 {
			fmt.Fprintf(buf, "  - Added imports: %s\n", inlineList(change.AddedImports))
		}
		if len(change.AddedExports) > 0 {
			fmt.Fprintf(buf, "  - Added exports: %s\n", inlineList(change.AddedExports))
		}
		if len(change.AddedStrings) > 0 {
			fmt.Fprintf(buf, "  - Added strings: %s\n", inlineList(change.AddedStrings))
		}
	}
	fmt.Fprintf(buf, "\n")
}

func writeArtifactChanges(buf *bytes.Buffer, title string, changes []model.ArtifactChange, limit int) {
	fmt.Fprintf(buf, "### %s\n\n", title)
	if len(changes) == 0 {
		fmt.Fprintf(buf, "No entries.\n\n")
		return
	}
	for i, change := range changes {
		if i >= limit {
			fmt.Fprintf(buf, "- ... %d more entries omitted from this section.\n\n", len(changes)-limit)
			return
		}
		fmt.Fprintf(buf, "- `%s`", change.Path)
		if change.Kind != "" {
			fmt.Fprintf(buf, " (%s)", change.Kind)
		}
		if change.OldSize > 0 || change.NewSize > 0 {
			fmt.Fprintf(buf, " size `%d` -> `%d`", change.OldSize, change.NewSize)
			if change.SizeDelta != 0 {
				fmt.Fprintf(buf, " (delta `%+d`)", change.SizeDelta)
			}
		}
		if change.ChangeClass != "" {
			fmt.Fprintf(buf, " class `%s`", change.ChangeClass)
		}
		fmt.Fprintf(buf, "\n")
		if len(change.RiskSignals) > 0 {
			fmt.Fprintf(buf, "  - Signals: %s\n", inlineList(change.RiskSignals))
		}
		if len(change.AddedStrings) > 0 {
			fmt.Fprintf(buf, "  - Added strings: %s\n", inlineList(change.AddedStrings))
		}
	}
	fmt.Fprintf(buf, "\n")
}

func valueOrUnknown(value string) string {
	if value == "" {
		return "unknown"
	}
	return value
}

func writeReviewQueue(buf *bytes.Buffer, executive model.ReportExecutive) {
	fmt.Fprintf(buf, "## Priority Review Queue\n\n")
	if len(executive.TopReviewTargets) == 0 {
		fmt.Fprintf(buf, "No priority targets were generated.\n\n")
		return
	}
	for _, target := range executive.TopReviewTargets {
		fmt.Fprintf(buf, "- `%s` (%s, %s): %s", target.Target, target.Kind, target.Priority, target.Reason)
		if len(target.Signals) > 0 {
			fmt.Fprintf(buf, " Signals: %s", inlineList(target.Signals))
		}
		fmt.Fprintf(buf, "\n")
	}
	if len(executive.NextActions) > 0 {
		fmt.Fprintf(buf, "\n### Next Actions\n\n")
		for _, action := range executive.NextActions {
			fmt.Fprintf(buf, "- %s\n", action)
		}
	}
	fmt.Fprintf(buf, "\n")
}

func writeGuidance(buf *bytes.Buffer, executive model.ReportExecutive) {
	fmt.Fprintf(buf, "## Analyst Guidance\n\n")
	if len(executive.BeginnerNotes) > 0 {
		fmt.Fprintf(buf, "### Beginner Read\n\n")
		for _, note := range executive.BeginnerNotes {
			fmt.Fprintf(buf, "- %s\n", note)
		}
		fmt.Fprintf(buf, "\n")
	}
	if len(executive.ResearcherChecklist) > 0 {
		fmt.Fprintf(buf, "### Researcher Checklist\n\n")
		for _, item := range executive.ResearcherChecklist {
			fmt.Fprintf(buf, "- %s\n", item)
		}
		fmt.Fprintf(buf, "\n")
	}
}

func writeServiceRegistryChanges(buf *bytes.Buffer, changes model.ChangeSummary) {
	fmt.Fprintf(buf, "### Service And Registry Changes\n\n")
	total := len(changes.AddedServices) + len(changes.RemovedServices) + len(changes.ChangedServices) + len(changes.AddedRegistry) + len(changes.RemovedRegistry) + len(changes.ChangedRegistry)
	if total == 0 {
		fmt.Fprintf(buf, "No entries.\n\n")
		return
	}
	for _, service := range changes.AddedServices {
		fmt.Fprintf(buf, "- Added service `%s` -> `%s`\n", service.Name, service.BinaryPath)
	}
	for _, service := range changes.RemovedServices {
		fmt.Fprintf(buf, "- Removed service `%s` -> `%s`\n", service.Name, service.BinaryPath)
	}
	for _, service := range changes.ChangedServices {
		fmt.Fprintf(buf, "- Changed service `%s`: binary `%s` -> `%s`, permissions `%s` -> `%s`, start `%s` -> `%s`\n", service.Name, service.Old.BinaryPath, service.New.BinaryPath, service.Old.Permissions, service.New.Permissions, service.Old.StartType, service.New.StartType)
	}
	for _, key := range changes.AddedRegistry {
		fmt.Fprintf(buf, "- Added registry key `%s`\n", key.Path)
	}
	for _, key := range changes.RemovedRegistry {
		fmt.Fprintf(buf, "- Removed registry key `%s`\n", key.Path)
	}
	for _, key := range changes.ChangedRegistry {
		fmt.Fprintf(buf, "- Changed registry key `%s`: ACL `%s` -> `%s`\n", key.Path, key.Old.ACL, key.New.ACL)
	}
	fmt.Fprintf(buf, "\n")
}
