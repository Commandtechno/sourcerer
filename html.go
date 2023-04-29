package main

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func fromHtml(ctx Context, res *http.Response) {
	Info(ctx.Depth, "Processing HTML...")

	parsedHtml, err := html.Parse(res.Body)
	if err != nil {
		Error(ctx.Depth, "Failed to parse HTML:", err)
		return
	}

	isScript := false
	script := ""

	isStyle := false
	style := ""

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			// js source maps
			if node.Data == "script" {
				isScript = true
				for _, attr := range node.Attr {
					// <script src="...">
					if attr.Key == "src" {
						script = attr.Val
						break
					}
				}
			}

			// css source maps
			if node.Data == "link" {
				for _, attr := range node.Attr {
					// <link rel="stylesheet">
					if attr.Key == "rel" && attr.Val == "stylesheet" {
						isStyle = true
						break
					}

					// <link as="style">
					if attr.Key == "as" && attr.Val == "style" {
						isStyle = true
						break
					}

					// <link as="script">
					if attr.Key == "as" && attr.Val == "script" {
						isScript = true
						break
					}
				}

				for _, attr := range node.Attr {
					// <link href="...">
					if attr.Key == "href" {
						style = attr.Val
						script = attr.Val
						break
					}
				}
			}
		}

		if isScript && script != "" {
			ref, err := ctx.Url.Parse(script)
			if err != nil {
				Error(ctx.Depth, "Failed to parse JavaScript url:", err)
			} else {
				jsCtx := ctx
				jsCtx.Url = *ctx.Url.ResolveReference(ref)
				jsCtx.Depth++

				js, ok := fetch(jsCtx)
				if ok {
					fromJs(jsCtx, js)
				}
			}
		}

		if isStyle && style != "" {
			ref, err := url.Parse(style)
			if err != nil {
				Error(ctx.Depth, "Failed to parse CSS url:", err)
			} else {
				cssCtx := ctx
				cssCtx.Url = *ctx.Url.ResolveReference(ref)
				cssCtx.Depth++

				css, ok := fetch(cssCtx)
				if ok {
					fromCss(cssCtx, css)
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(parsedHtml)
}
