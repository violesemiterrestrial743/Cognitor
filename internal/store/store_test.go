package store

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/kernelstub/cognitor/pkg/model"
)

func TestStoreSnapshotAndFindings(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "cognitor.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	snapshot := model.Snapshot{
		ID:        "s1",
		Name:      "sample",
		Path:      "/tmp/sample",
		CreatedAt: time.Unix(0, 0).UTC(),
		Binaries: []model.Binary{{
			ID:         "b1",
			SnapshotID: "s1",
			Path:       "driver.sys",
			Name:       "driver.sys",
			Kind:       "driver",
			SHA256:     "abc",
			Functions:  []model.Function{{ID: "fn1", Name: "Dispatch"}},
		}},
		Artifacts: []model.Artifact{{
			ID:         "a1",
			SnapshotID: "s1",
			Path:       "state.edb",
			Name:       "state.edb",
			Kind:       "edb",
			SHA256:     "def",
			Strings:    []string{"Policy"},
		}},
		Services: []model.Service{{Name: "svc", BinaryPath: "driver.sys", Permissions: "SYSTEM", StartType: "auto"}},
		Registry: []model.RegistryKey{{Path: `HKLM\Software\Cognitor`, ACL: "Administrators", Description: "test"}},
	}
	if err := db.SaveSnapshot(context.Background(), snapshot); err != nil {
		t.Fatal(err)
	}
	loaded, err := db.LoadSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Binaries[0].Functions[0].Name != "Dispatch" {
		t.Fatalf("unexpected loaded snapshot %#v", loaded)
	}
	if len(loaded.Artifacts) != 1 || loaded.Artifacts[0].Kind != "edb" {
		t.Fatalf("unexpected loaded artifacts %#v", loaded.Artifacts)
	}
	if len(loaded.Services) != 1 || loaded.Services[0].Name != "svc" {
		t.Fatalf("unexpected loaded services %#v", loaded.Services)
	}
	if len(loaded.Registry) != 1 || loaded.Registry[0].Path == "" {
		t.Fatalf("unexpected loaded registry %#v", loaded.Registry)
	}
	findings := []model.Finding{{ID: "f1", Title: "title", AffectedBinary: "driver.sys", Category: "access-control", Severity: "high", Evidence: []string{"SeAccessCheck"}}}
	if err := db.SaveFindings(context.Background(), findings); err != nil {
		t.Fatal(err)
	}
	loadedFindings, err := db.LoadFindings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(loadedFindings) != 1 || loadedFindings[0].ID != "f1" {
		t.Fatalf("unexpected findings %#v", loadedFindings)
	}
	changes := model.ChangeSummary{ChangedArtifacts: []model.ArtifactChange{{Path: "state.edb", Kind: "edb"}}}
	if err := db.SaveChangeSummary(context.Background(), changes); err != nil {
		t.Fatal(err)
	}
	loadedChanges, err := db.LoadChangeSummary(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(loadedChanges.ChangedArtifacts) != 1 || loadedChanges.ChangedArtifacts[0].Path != "state.edb" {
		t.Fatalf("unexpected changes %#v", loadedChanges)
	}
}
