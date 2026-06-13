package ingest

import (
	"debug/pe"
	"path/filepath"
	"strings"

	"github.com/kernelstub/cognitor/pkg/model"
)

func ReadPEMetadata(path string) (string, []string, []string, []model.Section) {
	file, err := pe.Open(path)
	if err != nil {
		return kindFromExt(path), nil, nil, nil
	}
	defer file.Close()
	imports, _ := file.ImportedSymbols()
	var sections []model.Section
	for _, section := range file.Sections {
		sections = append(sections, model.Section{Name: section.Name, Size: int64(section.Size)})
	}
	return kindFromExt(path), imports, nil, sections
}

func kindFromExt(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".sys":
		return "driver"
	case ".dll", ".ocx", ".cpl":
		return "library"
	default:
		return "executable"
	}
}
