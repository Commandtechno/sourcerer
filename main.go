package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Context struct {
	Url   url.URL
	Dir   string
	Depth int
	Cache map[string]struct{}
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Println("Usage: sourcerer [name] ...[urls]")
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		error(0, "Failed to get current path", err)
		return
	}

	name := args[0]
	dir := filepath.Join(cwd, name)

	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		error(0, "Failed to create output directory", err)
		return
	}

	for _, arg := range args {
		if strings.HasPrefix(arg, "http") {
			url, err := url.Parse(arg)
			if err != nil {
				fmt.Println("Error parsing url:", err)
				return
			}

			ctx := Context{
				Dir:   dir,
				Url:   *url,
				Cache: make(map[string]struct{}),
				Depth: 0,
			}

			fromUrl(ctx)
		}
	}
}