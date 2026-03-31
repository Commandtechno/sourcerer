package main

import (
	"path/filepath"
	"strings"
)

func fromUrl(ctx Context) {
	res, ok := fetch(ctx)
	if !ok {
		return
	}

	defer res.Body.Close()

	// clean up text/html; charset=utf-8
	contentType := strings.SplitN(res.Header.Get("Content-Type"), ";", 2)[0]

	switch contentType {
	case "text/html":
		fromHtml(ctx, res)

	case "text/css":
		fromCss(ctx, res)

	case "application/javascript":
		fromJs(ctx, res)

	case "application/octet-stream":
		fromSourceMap(ctx, res)

	case "application/json":
		fromSourceMap(ctx, res)

	default:
		ext := filepath.Ext(ctx.Url.Path)
		switch ext {
		case ".htm", ".html":
			fromHtml(ctx, res)

		case ".css":
			fromCss(ctx, res)

		case ".js":
			fromJs(ctx, res)

		case ".map":
			fromSourceMap(ctx, res)

		default:
			Error(ctx.Depth, "URL responded with unknown content type or extension:", contentType, ext)

		}
	}
}
