package diff

import (
	"context"
	"testing"

	"github.com/kernelstub/cognitor/pkg/model"
)

func TestAnalyzeProducesDefensiveFindings(t *testing.T) {
	oldSnapshot := model.Snapshot{Binaries: []model.Binary{{
		Path: "driver.sys",
		Functions: []model.Function{{
			ID:             "old",
			Name:           "Dispatch",
			NormalizedName: "dispatch",
			Calls:          []string{"memcpy"},
			Operations:     []string{"copy user buffer"},
		}},
	}}}
	newSnapshot := model.Snapshot{Binaries: []model.Binary{{
		Path: "driver.sys",
		Functions: []model.Function{{
			ID:             "new",
			Name:           "Dispatch",
			NormalizedName: "dispatch",
			Calls:          []string{"ProbeForRead", "SeAccessCheck", "memcpy"},
			Operations:     []string{"length check before copy", "copy user buffer"},
		}},
	}}}
	findings := Analyze(context.Background(), oldSnapshot, newSnapshot)
	if len(findings) == 0 {
		t.Fatal("expected findings")
	}
	if findings[0].Severity == "" || findings[0].RiskScore == 0 {
		t.Fatalf("expected scored finding, got %#v", findings[0])
	}
}
