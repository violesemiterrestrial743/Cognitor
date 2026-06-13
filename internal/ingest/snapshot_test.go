package ingest

import (
	"context"
	"testing"
)

func TestScanFixture(t *testing.T) {
	snapshot, err := Scan(context.Background(), Options{Name: "new", Path: "../../testdata/snapshots/new", Workers: 2, StringMinLength: 5})
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Binaries) != 2 {
		t.Fatalf("expected two binaries, got %d", len(snapshot.Binaries))
	}
	if len(snapshot.Artifacts) != 1 {
		t.Fatalf("expected one artifact, got %d", len(snapshot.Artifacts))
	}
	if snapshot.Artifacts[0].Kind != "edb" {
		t.Fatalf("expected edb artifact, got %#v", snapshot.Artifacts[0])
	}
	byName := map[string]int{}
	byKind := map[string]string{}
	for _, binary := range snapshot.Binaries {
		byName[binary.Name] = len(binary.Functions)
		byKind[binary.Name] = binary.Kind
	}
	if byName["driver.sys"] != 1 {
		t.Fatalf("expected driver function, got %#v", byName)
	}
	if byName["ntdll.dll"] != 1 {
		t.Fatalf("expected ntdll function, got %#v", byName)
	}
	if byKind["ntdll.dll"] != "library" {
		t.Fatalf("expected ntdll to be a library, got %#v", byKind)
	}
}
