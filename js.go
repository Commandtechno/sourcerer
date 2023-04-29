package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func fromJs(ctx Context, res *http.Response) {
	Info(ctx.Depth, "Processing JavaScript...")

	var sourceMapUrl url.URL

	str, err := ioutil.ReadAll(res.Body)
	if err != nil {
		Error(ctx.Depth, "Failed to read body:", err)
		return
	}

	startIndex := strings.LastIndex(string(str), "//# sourceMappingURL=")
	if startIndex == -1 {
		// fallback if there is no embedded source maps
		sourceMapUrl = ctx.Url
		sourceMapUrl.Path = sourceMapUrl.Path + ".map"
	} else {
		after := string(str[startIndex+len("//# sourceMappingURL="):])
		ref, err := url.Parse(after)
		if err != nil {
			Error(ctx.Depth, "Failed to parse JavaScript source map url:", err)
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
