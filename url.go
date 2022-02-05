package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func fromURL(baseUrl *url.URL, urlCache map[string]struct{}) {
	if _, cached := urlCache[baseUrl.String()]; cached {
		fmt.Println("Skipping duplicate url:", baseUrl.String())
		return
	} else {
		urlCache[baseUrl.String()] = struct{}{}
	}

	fmt.Println("Downloading", baseUrl.String())
	res, err := http.Get(baseUrl.String())
	if err != nil {
		fmt.Println("Error downloading", baseUrl.String(), err)
		return
	}

	if res.StatusCode != 200 {
		fmt.Println("URL responded with", res.Status)
		return
	}

	contentType := res.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") {
		fmt.Println("Detected HTML file")
		parsedHtml, err := html.Parse(res.Body)
		if err != nil {
			fmt.Println("Error parsing HTML", err)
			return
		}

		var walk func(*html.Node)
		walk = func(node *html.Node) {
			if node.Type == html.ElementNode {
				// js source maps
				if node.Data == "script" {
					for _, attr := range node.Attr {
						if attr.Key == "src" && strings.HasSuffix(attr.Val, ".js") {
							relativeUrl, err := baseUrl.Parse(attr.Val + ".map")
							if err != nil {
								fmt.Println("Error parsing relative url", err)
								return
							}

							fromURL(relativeUrl, urlCache)
						}
					}
				}

				// css source maps
				if node.Data == "link" {
					for _, attr := range node.Attr {
						if attr.Key == "href" && strings.HasSuffix(attr.Val, ".css") {
							relativeUrl, err := baseUrl.Parse(attr.Val + ".map")
							if err != nil {
								fmt.Println("Error parsing relative url", err)
								return
							}

							fromURL(relativeUrl, urlCache)
						}
					}
				}
			}

			for child := node.FirstChild; child != nil; child = child.NextSibling {
				walk(child)
			}
		}

		walk(parsedHtml)
	} else if strings.HasPrefix(contentType, "application/javascript") {
		// TODO: get source map from js file
		// fmt.Println("Detected JavaScript file")
		// body, err := ioutil.ReadAll(res.Body)
		// if err != nil {
		// 	fmt.Println("Error reading body", err)
		// 	return
		// }

		// // get source map url from body
		// sourceMapUrl := regexp.Match(`//# sourceMappingURL=(.*)`, body)
	} else {
		parse(res.Body)
	}
}