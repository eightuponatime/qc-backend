package main

import (
	"encoding/json"
	"os"
)

type ManifestEntry struct {
	File string   `json:"file"`
	CSS  []string `json:"css"`
}

type Manifest map[string]ManifestEntry

var manifest Manifest

func loadManifest() error {
	file, err := os.ReadFile("static/dist/.vite/manifest.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &manifest)
}

func asset(path string) string {
	if val, ok := manifest[path]; ok {
		return "/static/dist/" + val.File
	}
	return ""
}

func assetCSS(path string) []string {
	if val, ok := manifest[path]; ok {
		result := make([]string, 0, len(val.CSS))
		for _, cssFile := range val.CSS {
			result = append(result, "/static/dist/"+cssFile)
		}
		return result
	}
	return nil
}