package main

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strings"
)

func fetch(ctx Context) (*http.Response, bool) {
	if ctx.Url.Scheme == "data" {
		header, data, found := strings.Cut(ctx.Url.Opaque, ",")
		if !found {
			return nil, false
		}

		if !strings.HasPrefix(header, "application/json") {
			return nil, false
		}

		if strings.Contains(header, ";base64") {
			decoded, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				return nil, false
			}

			data = string(decoded)
		}

		return &http.Response{
			StatusCode: 200,
			Status:     "OK",
			Body:       io.NopCloser(strings.NewReader(data)),
		}, true
	}

	if ctx.Url.Scheme == "file" && ctx.Url.Host == "sourcerer" {
		f, err := os.Open(strings.TrimPrefix(ctx.Url.Path, "/"))
		if err != nil {
			Error(ctx.Depth, "Failed to read file:", err)
			return nil, false
		}

		return &http.Response{
			StatusCode: 200,
			Status:     "OK",
			Body:       io.NopCloser(f),
		}, true
	}

	if _, cached := ctx.Cache[ctx.Url.String()]; cached {
		return nil, false
	} else {
		ctx.Cache[ctx.Url.String()] = struct{}{}
	}

	Info(ctx.Depth, "Fetching URL:", ctx.Url.String())
	res, err := http.Get(ctx.Url.String())
	if err != nil {
		res.Body.Close()
		Error(ctx.Depth, "Failed to fetch URL:", err)
		return nil, false
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		res.Body.Close()
		Warn(ctx.Depth, "URL responded with status:", res.Status)
		return nil, false
	}

	// if res.Header.Get("Content-Type") == "" {
	// 	res.Body.Close()
	// 	Warn(ctx.Depth, "URL responded with no content type")
	// 	ext := filepath.Ext(ctx.Url.Path)
	// 	mimeType := mime.TypeByExtension(ext)
	// 	res.Header.Set("Content-Type", mimeType)
	// }

	Success(ctx.Depth, "URL responded with status:", res.Status)
	return res, true
}
