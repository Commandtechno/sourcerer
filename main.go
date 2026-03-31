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
		Error(0, "Failed to get current path", err)
		return
	}

	name := args[0]
	dir := filepath.Join(cwd, name)

	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		if !os.IsExist(err) {
			Error(0, "Failed to create output directory", err)
			return
		}

		children, err := os.ReadDir(dir)
		if err == nil && len(children) > 0 {
			Error(0, "Cannot write into existing non-empty directory")
			return
		}
	}

	for _, arg := range args[1:] {
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
		} else if _, err := os.Stat(arg); err == nil {
			url := url.URL{
				Scheme: "file",
				Host:   "sourcerer",
				Path:   arg,
			}

			ctx := Context{
				Dir:   dir,
				Url:   url,
				Cache: make(map[string]struct{}),
				Depth: 0,
			}

			fromUrl(ctx)
		} else {
			Warn(0, "Unrecognized arg, expected URL or file", arg)
		}
	}
}
