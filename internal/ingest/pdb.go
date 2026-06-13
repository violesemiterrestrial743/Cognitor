package ingest

import (
	"encoding/json"
	"os"
)

func ReadSymbols(path string) ([]string, error) {
	data, err := os.ReadFile(path + ".symbols.json")
	if err != nil {
		return nil, nil
	}
	var symbols []string
	if err := json.Unmarshal(data, &symbols); err != nil {
		return nil, err
	}
	return symbols, nil
}
