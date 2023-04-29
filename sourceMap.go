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

func filenamifyOpts(options *filenamify.Options) {
	options.Replacement = "_"
}

func pathify(path string) (string, error) {
	var result string
	for _, part := range strings.Split(filepath.Clean(path), string(os.PathSeparator)) {
		filename, err := filenamify.FilenamifyV2(part, filenamifyOpts)
		if err != nil {
			return "", err
		}

		result = filepath.Join(result, filename)
	}

	return result, nil
}

func fromSourceMap(ctx Context, res *http.Response) {
	Info(ctx.Depth, "Processing source map...")

	var sourceMap SourceMap
	if err := json.NewDecoder(res.Body).Decode(&sourceMap); err != nil {
		Error(ctx.Depth, "Failed to parse source map:", err)
		return
	}

	Info(ctx.Depth, "Unpacking", len(sourceMap.Sources), "sources...")
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
		relPath, err := pathify(source)
		if err != nil {
			Error(ctx.Depth, "Failed to normalize pathname:", err)
			continue
		}

		path := filepath.Join(ctx.Dir, relPath)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			Error(ctx.Depth, "Failed to create source directory:", err)
			continue
		}

		content := sourceMap.SourcesContent[index]
		if err := os.WriteFile(path, []byte(content), os.ModePerm); err != nil {
			Error(ctx.Depth, "Failed to write source file:", err)
			continue
		}
	}
}
