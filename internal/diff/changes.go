package diff

import (
	"sort"
	"strings"

	"github.com/kernelstub/cognitor/pkg/model"
)

func SummarizeChanges(oldSnapshot model.Snapshot, newSnapshot model.Snapshot) model.ChangeSummary {
	return normalizeChangeSummary(model.ChangeSummary{
		AddedBinaries:    addedBinaries(oldSnapshot.Binaries, newSnapshot.Binaries),
		RemovedBinaries:  removedBinaries(oldSnapshot.Binaries, newSnapshot.Binaries),
		ModifiedBinaries: modifiedBinaries(oldSnapshot.Binaries, newSnapshot.Binaries),
		AddedArtifacts:   addedArtifacts(oldSnapshot.Artifacts, newSnapshot.Artifacts),
		RemovedArtifacts: removedArtifacts(oldSnapshot.Artifacts, newSnapshot.Artifacts),
		ChangedArtifacts: changedArtifacts(oldSnapshot.Artifacts, newSnapshot.Artifacts),
		AddedServices:    addedServices(oldSnapshot.Services, newSnapshot.Services),
		RemovedServices:  removedServices(oldSnapshot.Services, newSnapshot.Services),
		ChangedServices:  changedServices(oldSnapshot.Services, newSnapshot.Services),
		AddedRegistry:    addedRegistry(oldSnapshot.Registry, newSnapshot.Registry),
		RemovedRegistry:  removedRegistry(oldSnapshot.Registry, newSnapshot.Registry),
		ChangedRegistry:  changedRegistry(oldSnapshot.Registry, newSnapshot.Registry),
	})
}

func normalizeChangeSummary(summary model.ChangeSummary) model.ChangeSummary {
	if summary.AddedBinaries == nil {
		summary.AddedBinaries = []model.BinaryChange{}
	}
	if summary.RemovedBinaries == nil {
		summary.RemovedBinaries = []model.BinaryChange{}
	}
	if summary.ModifiedBinaries == nil {
		summary.ModifiedBinaries = []model.BinaryChange{}
	}
	if summary.AddedArtifacts == nil {
		summary.AddedArtifacts = []model.ArtifactChange{}
	}
	if summary.RemovedArtifacts == nil {
		summary.RemovedArtifacts = []model.ArtifactChange{}
	}
	if summary.ChangedArtifacts == nil {
		summary.ChangedArtifacts = []model.ArtifactChange{}
	}
	if summary.AddedServices == nil {
		summary.AddedServices = []model.Service{}
	}
	if summary.RemovedServices == nil {
		summary.RemovedServices = []model.Service{}
	}
	if summary.ChangedServices == nil {
		summary.ChangedServices = []model.ServiceChange{}
	}
	if summary.AddedRegistry == nil {
		summary.AddedRegistry = []model.RegistryKey{}
	}
	if summary.RemovedRegistry == nil {
		summary.RemovedRegistry = []model.RegistryKey{}
	}
	if summary.ChangedRegistry == nil {
		summary.ChangedRegistry = []model.RegistryChange{}
	}
	return summary
}

func addedBinaries(oldValues []model.Binary, newValues []model.Binary) []model.BinaryChange {
	oldByPath := binariesByPath(oldValues)
	var out []model.BinaryChange
	for _, binary := range newValues {
		if _, ok := oldByPath[binary.Path]; !ok {
			out = append(out, binaryChange(model.Binary{}, binary))
		}
	}
	sortBinaryChanges(out)
	return out
}

func removedBinaries(oldValues []model.Binary, newValues []model.Binary) []model.BinaryChange {
	newByPath := binariesByPath(newValues)
	var out []model.BinaryChange
	for _, binary := range oldValues {
		if _, ok := newByPath[binary.Path]; !ok {
			out = append(out, binaryChange(binary, model.Binary{}))
		}
	}
	sortBinaryChanges(out)
	return out
}

func modifiedBinaries(oldValues []model.Binary, newValues []model.Binary) []model.BinaryChange {
	oldByPath := binariesByPath(oldValues)
	var out []model.BinaryChange
	for _, newBinary := range newValues {
		oldBinary, ok := oldByPath[newBinary.Path]
		if !ok || oldBinary.SHA256 == newBinary.SHA256 {
			continue
		}
		change := binaryChange(oldBinary, newBinary)
		change.AddedImports = limitStrings(AddedStrings(oldBinary.Imports, newBinary.Imports), 20)
		change.AddedExports = limitStrings(AddedStrings(oldBinary.Exports, newBinary.Exports), 20)
		change.AddedStrings = limitStrings(AddedStrings(oldBinary.Strings, newBinary.Strings), 30)
		change.RiskSignals = riskSignals(append(append([]string{}, change.AddedImports...), change.AddedStrings...))
		change.ChangeClass = classifyBinaryChange(change)
		out = append(out, change)
	}
	sortBinaryChanges(out)
	return out
}

func binaryChange(oldBinary model.Binary, newBinary model.Binary) model.BinaryChange {
	path := newBinary.Path
	name := newBinary.Name
	kind := newBinary.Kind
	if path == "" {
		path = oldBinary.Path
		name = oldBinary.Name
		kind = oldBinary.Kind
	}
	return model.BinaryChange{
		Path:       path,
		Name:       name,
		Kind:       kind,
		OldSHA256:  oldBinary.SHA256,
		NewSHA256:  newBinary.SHA256,
		OldSize:    oldBinary.Size,
		NewSize:    newBinary.Size,
		SizeDelta:  newBinary.Size - oldBinary.Size,
		OldVersion: oldBinary.Version,
		NewVersion: newBinary.Version,
	}
}

func addedArtifacts(oldValues []model.Artifact, newValues []model.Artifact) []model.ArtifactChange {
	oldByPath := artifactsByPath(oldValues)
	var out []model.ArtifactChange
	for _, artifact := range newValues {
		if _, ok := oldByPath[artifact.Path]; !ok {
			out = append(out, artifactChange(model.Artifact{}, artifact))
		}
	}
	sortArtifactChanges(out)
	return out
}

func removedArtifacts(oldValues []model.Artifact, newValues []model.Artifact) []model.ArtifactChange {
	newByPath := artifactsByPath(newValues)
	var out []model.ArtifactChange
	for _, artifact := range oldValues {
		if _, ok := newByPath[artifact.Path]; !ok {
			out = append(out, artifactChange(artifact, model.Artifact{}))
		}
	}
	sortArtifactChanges(out)
	return out
}

func changedArtifacts(oldValues []model.Artifact, newValues []model.Artifact) []model.ArtifactChange {
	oldByPath := artifactsByPath(oldValues)
	var out []model.ArtifactChange
	for _, newArtifact := range newValues {
		oldArtifact, ok := oldByPath[newArtifact.Path]
		if !ok || oldArtifact.SHA256 == newArtifact.SHA256 {
			continue
		}
		change := artifactChange(oldArtifact, newArtifact)
		change.AddedStrings = limitStrings(AddedStrings(oldArtifact.Strings, newArtifact.Strings), 30)
		change.RiskSignals = riskSignals(change.AddedStrings)
		change.ChangeClass = classifyArtifactChange(change)
		out = append(out, change)
	}
	sortArtifactChanges(out)
	return out
}

func artifactChange(oldArtifact model.Artifact, newArtifact model.Artifact) model.ArtifactChange {
	path := newArtifact.Path
	name := newArtifact.Name
	kind := newArtifact.Kind
	if path == "" {
		path = oldArtifact.Path
		name = oldArtifact.Name
		kind = oldArtifact.Kind
	}
	return model.ArtifactChange{
		Path:      path,
		Name:      name,
		Kind:      kind,
		OldSHA256: oldArtifact.SHA256,
		NewSHA256: newArtifact.SHA256,
		OldSize:   oldArtifact.Size,
		NewSize:   newArtifact.Size,
		SizeDelta: newArtifact.Size - oldArtifact.Size,
	}
}

func riskSignals(values []string) []string {
	type signal struct {
		label    string
		needles  []string
		priority int
	}
	signals := []signal{
		{label: "authorization", needles: []string{"accesscheck", "access check", "acl", "securitydescriptor", "privilege", "admin", "auth", "allow", "deny"}, priority: 1},
		{label: "token-or-identity", needles: []string{"token", "sid", "principal", "impersonat", "identity", "logon"}, priority: 2},
		{label: "kernel-boundary", needles: []string{"ioctl", "irp", "probeforread", "probeforwrite", "devicecontrol", "user buffer"}, priority: 3},
		{label: "memory-safety", needles: []string{"memcpy", "memmove", "copy", "length", "bounds", "overflow", "buffer", "size"}, priority: 4},
		{label: "service-control", needles: []string{"service", "scm", "starttype", "launch", "daemon"}, priority: 5},
		{label: "registry-policy", needles: []string{"registry", "hklm", "hkcu", "policy", "regkey"}, priority: 6},
		{label: "crypto-or-secret", needles: []string{"crypto", "cert", "secret", "credential", "password", "dpapi", "key"}, priority: 7},
		{label: "network-surface", needles: []string{"rpc", "alpc", "socket", "http", "named pipe", "endpoint", "port"}, priority: 8},
		{label: "sandbox-boundary", needles: []string{"appcontainer", "sandbox", "integrity", "lowbox", "capability"}, priority: 9},
		{label: "telemetry-or-audit", needles: []string{"event", "audit", "telemetry", "trace", "etw", "log"}, priority: 10},
	}
	found := map[string]int{}
	for _, value := range values {
		lower := strings.ToLower(value)
		for _, candidate := range signals {
			for _, needle := range candidate.needles {
				if strings.Contains(lower, needle) {
					if previous, ok := found[candidate.label]; !ok || candidate.priority < previous {
						found[candidate.label] = candidate.priority
					}
				}
			}
		}
	}
	out := make([]string, 0, len(found))
	for label := range found {
		out = append(out, label)
	}
	sort.Slice(out, func(i, j int) bool {
		return found[out[i]] < found[out[j]]
	})
	return out
}

func classifyBinaryChange(change model.BinaryChange) string {
	switch {
	case hasAnySignal(change.RiskSignals, "authorization", "token-or-identity", "kernel-boundary", "memory-safety"):
		return "security-hardening"
	case len(change.AddedImports) > 0 || len(change.AddedExports) > 0:
		return "interface-surface"
	case len(change.RiskSignals) > 0:
		return "behavioral-signal"
	default:
		return "binary-drift"
	}
}

func classifyArtifactChange(change model.ArtifactChange) string {
	switch {
	case hasAnySignal(change.RiskSignals, "authorization", "registry-policy", "service-control", "sandbox-boundary"):
		return "policy-or-configuration"
	case hasAnySignal(change.RiskSignals, "telemetry-or-audit"):
		return "telemetry"
	case len(change.RiskSignals) > 0:
		return "evidence-signal"
	default:
		return "artifact-drift"
	}
}

func hasAnySignal(signals []string, wanted ...string) bool {
	set := map[string]struct{}{}
	for _, signal := range signals {
		set[signal] = struct{}{}
	}
	for _, value := range wanted {
		if _, ok := set[value]; ok {
			return true
		}
	}
	return false
}

func addedServices(oldValues []model.Service, newValues []model.Service) []model.Service {
	oldByName := servicesByName(oldValues)
	var out []model.Service
	for _, service := range newValues {
		if _, ok := oldByName[service.Name]; !ok {
			out = append(out, service)
		}
	}
	sortServices(out)
	return out
}

func removedServices(oldValues []model.Service, newValues []model.Service) []model.Service {
	newByName := servicesByName(newValues)
	var out []model.Service
	for _, service := range oldValues {
		if _, ok := newByName[service.Name]; !ok {
			out = append(out, service)
		}
	}
	sortServices(out)
	return out
}

func changedServices(oldValues []model.Service, newValues []model.Service) []model.ServiceChange {
	oldByName := servicesByName(oldValues)
	var out []model.ServiceChange
	for _, service := range newValues {
		oldService, ok := oldByName[service.Name]
		if ok && oldService != service {
			out = append(out, model.ServiceChange{Name: service.Name, Old: oldService, New: service})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func addedRegistry(oldValues []model.RegistryKey, newValues []model.RegistryKey) []model.RegistryKey {
	oldByPath := registryByPath(oldValues)
	var out []model.RegistryKey
	for _, key := range newValues {
		if _, ok := oldByPath[key.Path]; !ok {
			out = append(out, key)
		}
	}
	sortRegistry(out)
	return out
}

func removedRegistry(oldValues []model.RegistryKey, newValues []model.RegistryKey) []model.RegistryKey {
	newByPath := registryByPath(newValues)
	var out []model.RegistryKey
	for _, key := range oldValues {
		if _, ok := newByPath[key.Path]; !ok {
			out = append(out, key)
		}
	}
	sortRegistry(out)
	return out
}

func changedRegistry(oldValues []model.RegistryKey, newValues []model.RegistryKey) []model.RegistryChange {
	oldByPath := registryByPath(oldValues)
	var out []model.RegistryChange
	for _, key := range newValues {
		oldKey, ok := oldByPath[key.Path]
		if ok && oldKey != key {
			out = append(out, model.RegistryChange{Path: key.Path, Old: oldKey, New: key})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Path < out[j].Path })
	return out
}

func binariesByPath(values []model.Binary) map[string]model.Binary {
	out := map[string]model.Binary{}
	for _, value := range values {
		out[windowsPathKey(value.Path)] = value
	}
	return out
}

func artifactsByPath(values []model.Artifact) map[string]model.Artifact {
	out := map[string]model.Artifact{}
	for _, value := range values {
		out[windowsPathKey(value.Path)] = value
	}
	return out
}

func servicesByName(values []model.Service) map[string]model.Service {
	out := map[string]model.Service{}
	for _, value := range values {
		out[value.Name] = value
	}
	return out
}

func registryByPath(values []model.RegistryKey) map[string]model.RegistryKey {
	out := map[string]model.RegistryKey{}
	for _, value := range values {
		out[value.Path] = value
	}
	return out
}

func limitStrings(values []string, limit int) []string {
	if len(values) <= limit {
		return values
	}
	return values[:limit]
}

func sortBinaryChanges(values []model.BinaryChange) {
	sort.Slice(values, func(i, j int) bool { return values[i].Path < values[j].Path })
}

func sortArtifactChanges(values []model.ArtifactChange) {
	sort.Slice(values, func(i, j int) bool { return values[i].Path < values[j].Path })
}

func sortServices(values []model.Service) {
	sort.Slice(values, func(i, j int) bool { return values[i].Name < values[j].Name })
}

func sortRegistry(values []model.RegistryKey) {
	sort.Slice(values, func(i, j int) bool { return values[i].Path < values[j].Path })
}
