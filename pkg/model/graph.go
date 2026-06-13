package model

type Node struct {
	ID    string            `json:"id"`
	Type  string            `json:"type"`
	Label string            `json:"label"`
	Attrs map[string]string `json:"attrs"`
}

type Edge struct {
	From  string            `json:"from"`
	To    string            `json:"to"`
	Type  string            `json:"type"`
	Attrs map[string]string `json:"attrs"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}
