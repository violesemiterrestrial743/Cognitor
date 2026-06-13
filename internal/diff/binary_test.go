package diff

import (
	"testing"

	"github.com/kernelstub/cognitor/pkg/model"
)

func TestMatchBinariesUsesWindowsCaseInsensitivePaths(t *testing.T) {
	oldSnapshot := model.Snapshot{Binaries: []model.Binary{{Path: `System32/NTDLL.DLL`, Name: "NTDLL.DLL", Kind: "library"}}}
	newSnapshot := model.Snapshot{Binaries: []model.Binary{{Path: `system32/ntdll.dll`, Name: "ntdll.dll", Kind: "library"}}}

	pairs := MatchBinaries(oldSnapshot, newSnapshot)
	if len(pairs) != 1 {
		t.Fatalf("expected ntdll pair, got %#v", pairs)
	}
	if pairs[0][0].Name != "NTDLL.DLL" || pairs[0][1].Name != "ntdll.dll" {
		t.Fatalf("unexpected pair %#v", pairs[0])
	}
}
