package snapshot

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateEmptySnapshot(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "snapshot")
	result, err := Create(context.Background(), CreateOptions{Name: "old", Path: dir})
	if err != nil {
		t.Fatal(err)
	}
	if result.CreatedFiles != 3 {
		t.Fatalf("expected three created files, got %#v", result)
	}
	for _, name := range []string{"services.json", "registry.json", "SNAPSHOT.md"} {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			t.Fatalf("missing %s: %v", name, err)
		}
	}
}

func TestCreateSnapshotFromSource(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "source")
	dest := filepath.Join(root, "dest")
	if err := os.MkdirAll(filepath.Join(source, "drivers"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "drivers", "sample.sys"), []byte("driver"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "drivers", "sample.sys.analysis.json"), []byte(`{"functions":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "ignore.txt"), []byte("ignore"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "state.edb"), []byte("database evidence"), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := Create(context.Background(), CreateOptions{Name: "new", Path: dest, Source: source})
	if err != nil {
		t.Fatal(err)
	}
	if result.CopiedFiles != 3 {
		t.Fatalf("expected three copied files, got %#v", result)
	}
	if _, err := os.Stat(filepath.Join(dest, "drivers", "sample.sys")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dest, "ignore.txt")); err == nil {
		t.Fatal("unexpected copied ignored file")
	}
	if _, err := os.Stat(filepath.Join(dest, "state.edb")); err != nil {
		t.Fatal(err)
	}
}
