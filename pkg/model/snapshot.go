package model

import "time"

type Snapshot struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Path      string        `json:"path"`
	CreatedAt time.Time     `json:"created_at"`
	Binaries  []Binary      `json:"binaries"`
	Artifacts []Artifact    `json:"artifacts"`
	Services  []Service     `json:"services"`
	Registry  []RegistryKey `json:"registry"`
}

type RegistryKey struct {
	Path        string `json:"path"`
	ACL         string `json:"acl"`
	Description string `json:"description"`
}

type Service struct {
	Name        string `json:"name"`
	BinaryPath  string `json:"binary_path"`
	Permissions string `json:"permissions"`
	StartType   string `json:"start_type"`
}
