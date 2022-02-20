package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func fromCss(ctx Context, res *http.Response) {
	info(ctx.Depth, "Processing CSS...")

	var sourceMapUrl url.URL

	str, err := ioutil.ReadAll(res.Body)
	if err != nil {
		error(ctx.Depth, "Failed to read css response body:", err)
		return
	}

	startIndex := strings.LastIndex(string(str), "/*# sourceMappingURL=")
	if startIndex == -1 {
		// fallback if there is no embedded source maps
		sourceMapUrl = ctx.Url
		sourceMapUrl.Path = sourceMapUrl.Path + ".map"
	} else {
		after := string(str[startIndex+len("/*# sourceMappingURL="):])
		endIndex := strings.LastIndex(after, "*/")
		if endIndex != -1 {
			after = after[:endIndex]
		}

		ref, err := url.Parse(after)
		if err != nil {
			error(ctx.Depth, "Failed to parse CSS source map url:", err)
			return
		}

		sourceMapUrl = *ctx.Url.ResolveReference(ref)
	}

	sourceMapCtx := ctx
	sourceMapCtx.Url = sourceMapUrl
	sourceMapCtx.Depth++

	sourceMap, ok := fetch(sourceMapCtx)
	if !ok {
		return
	}

	fromSourceMap(sourceMapCtx, sourceMap)
}