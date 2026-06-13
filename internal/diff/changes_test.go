package diff

import (
	"testing"

	"github.com/kernelstub/cognitor/pkg/model"
)

func TestSummarizeChangesIncludesArtifactsServicesAndRegistry(t *testing.T) {
	oldSnapshot := model.Snapshot{
		Binaries:  []model.Binary{{Path: "driver.sys", Name: "driver.sys", Kind: "driver", SHA256: "old", Size: 10, Imports: []string{"old"}, Strings: []string{"same"}}},
		Artifacts: []model.Artifact{{Path: "state.edb", Name: "state.edb", Kind: "edb", SHA256: "old", Size: 20, Strings: []string{"same"}}},
		Services:  []model.Service{{Name: "svc", BinaryPath: "old.sys", Permissions: "SYSTEM", StartType: "manual"}},
		Registry:  []model.RegistryKey{{Path: `HKLM\Test`, ACL: "old", Description: "policy"}},
	}
	newSnapshot := model.Snapshot{
		Binaries:  []model.Binary{{Path: "driver.sys", Name: "driver.sys", Kind: "driver", SHA256: "new", Size: 12, Imports: []string{"old", "SeAccessCheck"}, Strings: []string{"same", "new-policy"}}},
		Artifacts: []model.Artifact{{Path: "state.edb", Name: "state.edb", Kind: "edb", SHA256: "new", Size: 30, Strings: []string{"same", "AccessCheckRequired"}}},
		Services:  []model.Service{{Name: "svc", BinaryPath: "new.sys", Permissions: "SYSTEM", StartType: "auto"}},
		Registry:  []model.RegistryKey{{Path: `HKLM\Test`, ACL: "new", Description: "policy"}},
	}

	changes := SummarizeChanges(oldSnapshot, newSnapshot)
	if len(changes.ModifiedBinaries) != 1 || changes.ModifiedBinaries[0].AddedImports[0] != "SeAccessCheck" {
		t.Fatalf("unexpected binary changes %#v", changes.ModifiedBinaries)
	}
	if len(changes.ChangedArtifacts) != 1 || changes.ChangedArtifacts[0].AddedStrings[0] != "AccessCheckRequired" {
		t.Fatalf("unexpected artifact changes %#v", changes.ChangedArtifacts)
	}
	if len(changes.ChangedServices) != 1 {
		t.Fatalf("unexpected service changes %#v", changes.ChangedServices)
	}
	if len(changes.ChangedRegistry) != 1 {
		t.Fatalf("unexpected registry changes %#v", changes.ChangedRegistry)
	}
}
