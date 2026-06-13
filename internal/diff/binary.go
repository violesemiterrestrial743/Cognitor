package diff

import (
	"path/filepath"
	"strings"

	"github.com/kernelstub/cognitor/pkg/model"
)

func MatchBinaries(oldSnapshot model.Snapshot, newSnapshot model.Snapshot) [][2]model.Binary {
	oldByPath := map[string]model.Binary{}
	for _, binary := range oldSnapshot.Binaries {
		oldByPath[windowsPathKey(binary.Path)] = binary
	}
	var pairs [][2]model.Binary
	for _, newBinary := range newSnapshot.Binaries {
		if oldBinary, ok := oldByPath[windowsPathKey(newBinary.Path)]; ok {
			pairs = append(pairs, [2]model.Binary{oldBinary, newBinary})
		}
	}
	return pairs
}

func windowsPathKey(path string) string {
	return strings.ToLower(filepath.ToSlash(filepath.Clean(path)))
}
