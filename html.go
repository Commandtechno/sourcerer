package main

import (
	"net/http"
	"net/url"

	"golang.org/x/net/html"
)

func fromHtml(ctx Context, res *http.Response) {
	info(ctx.Depth, "Processing HTML...")

	parsedHtml, err := html.Parse(res.Body)
	if err != nil {
		error(ctx.Depth, "Failed to parse HTML:", err)
		return
	}

	var walk func(*html.Node)
	walk = func(node *html.Node) {
		if node.Type == html.ElementNode {
			// js source maps
			if node.Data == "script" {
				for _, attr := range node.Attr {
					if attr.Key == "src" {
						ref, err := ctx.Url.Parse(attr.Val)
						if err != nil {
							error(ctx.Depth, "Failed to parse JavaScript url:", err)
							return
						}

						jsCtx := ctx
						jsCtx.Url = *ctx.Url.ResolveReference(ref)
						jsCtx.Depth++

						js, ok := fetch(jsCtx)
						if !ok {
							return
						}

						fromJs(jsCtx, js)
					}
				}
			}

			// css source maps
			if node.Data == "link" {
				isStyle := false
				href := ""
				for _, attr := range node.Attr {
					if attr.Key == "rel" && attr.Val == "stylesheet" {
						isStyle = true
					}

					if attr.Key == "as" && attr.Val == "style" {
						isStyle = true
					}

					if attr.Key == "href" {
						href = attr.Val
					}
				}

				if !isStyle || href == "" {
					return
				}

				for _, attr := range node.Attr {
					if attr.Key == "href" {
						ref, err := url.Parse(attr.Val)
						if err != nil {
							error(ctx.Depth, "Failed to parse CSS url:", err)
							return
						}

						cssCtx := ctx
						cssCtx.Url = *ctx.Url.ResolveReference(ref)
						cssCtx.Depth++

						css, ok := fetch(cssCtx)
						if !ok {
							return
						}

						fromCss(cssCtx, css)
					}
				}
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}

	walk(parsedHtml)
}