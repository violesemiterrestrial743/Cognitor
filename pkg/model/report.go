package model

type Report struct {
	Metadata  ReportMetadata  `json:"metadata"`
	Summary   ReportSummary   `json:"summary"`
	Executive ReportExecutive `json:"executive"`
	Changes   ChangeSummary   `json:"changes"`
	Findings  []Finding       `json:"findings"`
	Graph     Graph           `json:"graph"`
}

type ReportMetadata struct {
	GeneratedAt string        `json:"generated_at,omitempty"`
	ToolVersion string        `json:"tool_version,omitempty"`
	OldSnapshot SnapshotBrief `json:"old_snapshot,omitempty"`
	NewSnapshot SnapshotBrief `json:"new_snapshot,omitempty"`
}

type SnapshotBrief struct {
	Name          string `json:"name,omitempty"`
	Path          string `json:"path,omitempty"`
	BinaryCount   int    `json:"binary_count"`
	ArtifactCount int    `json:"artifact_count"`
	ServiceCount  int    `json:"service_count"`
	RegistryCount int    `json:"registry_count"`
}

type ReportSummary struct {
	TotalFindings        int            `json:"total_findings"`
	BySeverity           map[string]int `json:"by_severity"`
	ByCategory           map[string]int `json:"by_category"`
	TopChangedComponents []string       `json:"top_changed_components"`
}

type ReportExecutive struct {
	RiskLevel              string         `json:"risk_level"`
	Priority               string         `json:"priority"`
	TotalChangedItems      int            `json:"total_changed_items"`
	HighConfidenceFindings int            `json:"high_confidence_findings"`
	ChangedArtifactKinds   []string       `json:"changed_artifact_kinds"`
	TopReviewTargets       []ReviewTarget `json:"top_review_targets"`
	BeginnerNotes          []string       `json:"beginner_notes"`
	ResearcherChecklist    []string       `json:"researcher_checklist"`
	NextActions            []string       `json:"next_actions"`
}

type ReviewTarget struct {
	Target   string   `json:"target"`
	Kind     string   `json:"kind"`
	Reason   string   `json:"reason"`
	Priority string   `json:"priority"`
	Signals  []string `json:"signals,omitempty"`
}

type ChangeSummary struct {
	AddedBinaries    []BinaryChange   `json:"added_binaries"`
	RemovedBinaries  []BinaryChange   `json:"removed_binaries"`
	ModifiedBinaries []BinaryChange   `json:"modified_binaries"`
	AddedArtifacts   []ArtifactChange `json:"added_artifacts"`
	RemovedArtifacts []ArtifactChange `json:"removed_artifacts"`
	ChangedArtifacts []ArtifactChange `json:"changed_artifacts"`
	AddedServices    []Service        `json:"added_services"`
	RemovedServices  []Service        `json:"removed_services"`
	ChangedServices  []ServiceChange  `json:"changed_services"`
	AddedRegistry    []RegistryKey    `json:"added_registry"`
	RemovedRegistry  []RegistryKey    `json:"removed_registry"`
	ChangedRegistry  []RegistryChange `json:"changed_registry"`
}

type BinaryChange struct {
	Path         string   `json:"path"`
	Name         string   `json:"name"`
	Kind         string   `json:"kind"`
	ChangeClass  string   `json:"change_class,omitempty"`
	RiskSignals  []string `json:"risk_signals,omitempty"`
	OldSHA256    string   `json:"old_sha256,omitempty"`
	NewSHA256    string   `json:"new_sha256,omitempty"`
	OldSize      int64    `json:"old_size,omitempty"`
	NewSize      int64    `json:"new_size,omitempty"`
	SizeDelta    int64    `json:"size_delta,omitempty"`
	OldVersion   string   `json:"old_version,omitempty"`
	NewVersion   string   `json:"new_version,omitempty"`
	AddedImports []string `json:"added_imports,omitempty"`
	AddedExports []string `json:"added_exports,omitempty"`
	AddedStrings []string `json:"added_strings,omitempty"`
}

type ArtifactChange struct {
	Path         string   `json:"path"`
	Name         string   `json:"name"`
	Kind         string   `json:"kind"`
	ChangeClass  string   `json:"change_class,omitempty"`
	RiskSignals  []string `json:"risk_signals,omitempty"`
	OldSHA256    string   `json:"old_sha256,omitempty"`
	NewSHA256    string   `json:"new_sha256,omitempty"`
	OldSize      int64    `json:"old_size,omitempty"`
	NewSize      int64    `json:"new_size,omitempty"`
	SizeDelta    int64    `json:"size_delta,omitempty"`
	AddedStrings []string `json:"added_strings,omitempty"`
}

type ServiceChange struct {
	Name string  `json:"name"`
	Old  Service `json:"old"`
	New  Service `json:"new"`
}

type RegistryChange struct {
	Path string      `json:"path"`
	Old  RegistryKey `json:"old"`
	New  RegistryKey `json:"new"`
}
