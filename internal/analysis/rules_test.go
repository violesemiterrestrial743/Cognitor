package analysis

import (
	"context"
	"testing"

	"github.com/kernelstub/cognitor/pkg/model"
)

func TestDefaultEngineDetectsAccessAndBoundsChecks(t *testing.T) {
	change := model.SemanticChange{
		Binary:      model.Binary{Path: "driver.sys"},
		OldFunction: model.Function{Name: "Dispatch"},
		NewFunction: model.Function{Name: "Dispatch"},
		AddedCalls:  []string{"SeAccessCheck", "ProbeForRead"},
		AddedOps:    []string{"length check before copy"},
		Similarity:  1,
	}
	findings := DefaultEngine().Evaluate(context.Background(), []model.SemanticChange{change})
	if len(findings) < 2 {
		t.Fatalf("expected at least two findings, got %d", len(findings))
	}
	categories := map[string]bool{}
	for _, finding := range findings {
		categories[finding.Category] = true
	}
	if !categories["access-control"] {
		t.Fatalf("missing access-control finding: %#v", findings)
	}
	if !categories["memory-safety"] {
		t.Fatalf("missing memory-safety finding: %#v", findings)
	}
}

func TestDefaultEngineDetectsResearcherBoundaryRules(t *testing.T) {
	change := model.SemanticChange{
		Binary:       model.Binary{Path: "ntdll.dll"},
		OldFunction:  model.Function{Name: "NtCreateFile"},
		NewFunction:  model.Function{Name: "NtCreateFile"},
		AddedCalls:   []string{"ObReferenceObjectByHandle", "RpcBindingInqAuthClientEx", "CoInitializeSecurity"},
		AddedStrings: []string{"Ndr conformant array range check", "syscall previous mode validation"},
		AddedOps:     []string{"access mask validation", "rundown protection"},
		Similarity:   1,
	}
	findings := DefaultEngine().Evaluate(context.Background(), []model.SemanticChange{change})
	categories := map[string]bool{}
	for _, finding := range findings {
		categories[finding.Category] = true
	}
	for _, category := range []string{"native-api-boundary", "handle-validation", "object-lifetime", "rpc-hardening", "com-hardening", "marshalling-validation"} {
		if !categories[category] {
			t.Fatalf("missing %s finding: %#v", category, findings)
		}
	}
}
