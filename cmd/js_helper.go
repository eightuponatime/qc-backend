package main

import (
	"encoding/json"
	"os"
)

type Manifest map[string]struct {
	File string `json:"file"`
}

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
