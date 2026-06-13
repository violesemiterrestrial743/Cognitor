package ingest

import "os"

func ReadManifest(path string) string {
	for _, suffix := range []string{".manifest", ".manifest.json"} {
		data, err := os.ReadFile(path + suffix)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
