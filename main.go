package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Please provide one or more url or file")
		return
	}

	for _, arg := range args {
		if strings.HasPrefix(arg, "http") {
			url, err := url.Parse(arg)
			if err != nil {
				fmt.Println("Error parsing url:", err)
				return
			}

			fromURL(url, map[string]struct{}{})
		}
	}
}