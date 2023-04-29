package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
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
			Body:       ioutil.NopCloser(strings.NewReader(data)),
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
		Error(ctx.Depth, "Failed to fetch URL:", err)
		return nil, false
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		Warn(ctx.Depth, "URL responded with status:", res.Status)
		return nil, false
	}

	Success(ctx.Depth, "URL responded with status:", res.Status)
	return res, true
}
