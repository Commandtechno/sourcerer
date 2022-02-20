package main

import (
	"net/http"
)

func fetch(ctx Context) (*http.Response, bool) {
	if _, cached := ctx.Cache[ctx.Url.String()]; cached {
		return nil, false
	} else {
		ctx.Cache[ctx.Url.String()] = struct{}{}
	}

	info(ctx.Depth, "Fetching URL:", ctx.Url.String())
	res, err := http.Get(ctx.Url.String())
	if err != nil {
		error(ctx.Depth, "Failed to fetch URL:", err)
		return nil, false
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		warn(ctx.Depth, "URL responded with status:", res.Status)
		return nil, false
	}

	success(ctx.Depth, "URL responded with status:", res.Status)
	return res, true
}