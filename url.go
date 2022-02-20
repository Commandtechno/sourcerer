package main

import (
	"strings"
)

func fromUrl(ctx Context) {
	res, ok := fetch(ctx)
	if !ok {
		return
	}

	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		warn(ctx.Depth, "URL responded with no content type:", contentType)
		return
	}

	// clean up text/html; charset=utf-8
	contentType = strings.SplitN(contentType, ";", 2)[0]

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
		error(ctx.Depth, "URL responded with unknown content type:", contentType)
	}
}