package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type SourceMap struct {
	Sources        []string `json:"sources"`
	SourcesContent []string `json:"sourcesContent"`
}

func fromSourceMap(ctx Context, res *http.Response) {
	info(ctx.Depth, "Processing source map...")

	var sourceMap SourceMap
	if err := json.NewDecoder(res.Body).Decode(&sourceMap); err != nil {
		error(ctx.Depth, "Failed to parse source map:", err)
		return
	}

	info(ctx.Depth, "Unpacking", len(sourceMap.Sources), "sources...")
	ctx.Depth++

	for index, source := range sourceMap.Sources {
		queryIndex := strings.Index(source, "?")
		if queryIndex != -1 {
			source = source[:queryIndex]
		}

		// remove reserved characters
		// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file
		source = strings.ReplaceAll(source, "<", "")
		source = strings.ReplaceAll(source, ">", "")
		source = strings.ReplaceAll(source, ":", "")
		source = strings.ReplaceAll(source, "\"", "")
		source = strings.ReplaceAll(source, "<", "")

		path := filepath.Join(filepath.Clean(ctx.Dir), source)
		dir := filepath.Dir(path)

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			error(ctx.Depth, "Failed to create source directory:", err)
			continue
		}

		content := sourceMap.SourcesContent[index]
		if err := os.WriteFile(path, []byte(content), os.ModePerm); err != nil {
			error(ctx.Depth, "Failed to write source file:", err)
			continue
		}
	}
}