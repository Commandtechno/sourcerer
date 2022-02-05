package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type SourceMap struct {
	Sources        []string `json:"sources"`
	SourcesContent []string `json:"sourcesContent"`
}

func parse(rawSourceMap io.Reader) {
	fmt.Println("Parsing source map")

	var sourceMap SourceMap
	err := json.NewDecoder(rawSourceMap).Decode(&sourceMap)
	if err != nil {
		fmt.Println("Error parsing source map", err)
		return
	}

	currentPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting current path", err)
		return
	}

	outputDir := filepath.Join(filepath.Dir(currentPath), "output")
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		fmt.Println("Error creating output directory", err)
		return
	}

	fmt.Println("Unpacking source map of", len(sourceMap.Sources), "files")
	for index, source := range sourceMap.Sources {
		path := filepath.Join(outputDir, filepath.FromSlash(strings.ReplaceAll("/"+source, "://", "")))
		dir := filepath.Dir(path)
		os.MkdirAll(dir, os.ModePerm)

		content := sourceMap.SourcesContent[index]
		os.WriteFile(path, []byte(content), os.ModePerm)
	}
}