package util

import (
	"path/filepath"
	"strings"
)

func IsBinaryLike(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".exe", ".dll", ".sys", ".drv", ".ocx", ".cpl":
		return true
	default:
		return false
	}
}

func IsArtifactLike(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".edb", ".dat", ".log", ".evtx", ".etl", ".reg", ".json", ".xml", ".ini", ".inf", ".cfg", ".conf":
		return true
	default:
		return false
	}
}

func NormalizeSlash(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}
