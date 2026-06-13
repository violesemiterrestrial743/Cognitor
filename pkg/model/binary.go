package model

type Binary struct {
	ID         string     `json:"id"`
	SnapshotID string     `json:"snapshot_id"`
	Path       string     `json:"path"`
	Name       string     `json:"name"`
	Kind       string     `json:"kind"`
	SHA256     string     `json:"sha256"`
	Size       int64      `json:"size"`
	Version    string     `json:"version"`
	Signer     string     `json:"signer"`
	Imports    []string   `json:"imports"`
	Exports    []string   `json:"exports"`
	Sections   []Section  `json:"sections"`
	Strings    []string   `json:"strings"`
	Functions  []Function `json:"functions"`
	Manifest   string     `json:"manifest"`
}

type Artifact struct {
	ID         string   `json:"id"`
	SnapshotID string   `json:"snapshot_id"`
	Path       string   `json:"path"`
	Name       string   `json:"name"`
	Kind       string   `json:"kind"`
	SHA256     string   `json:"sha256"`
	Size       int64    `json:"size"`
	Strings    []string `json:"strings"`
}

type Section struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}
