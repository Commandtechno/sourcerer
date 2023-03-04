package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/flytam/filenamify"
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
		// escape parent directory references
		for strings.HasPrefix(source, "../") {
			source = strings.TrimPrefix(source, "../")
		}

		// remove query parameters
		queryIndex := strings.Index(source, "?")
		if queryIndex != -1 {
			source = source[:queryIndex]
		}

		hashIndex := strings.Index(source, "?")
		if hashIndex != -1 {
			source = source[:hashIndex]
		}

		// remove reserved characters
		source, err := filenamify.FilenamifyV2(source, func(options *filenamify.Options) {
			options.Replacement = "_"
		})

		if err != nil {
			error(ctx.Depth, "Failed to filenamify source:", err)
			continue
		}

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
