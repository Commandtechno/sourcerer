sourcerer is a cli tool to get a website's source code from its source maps written in go

# usage

```
sourcerer [name] ...[urls]
```

the url can either be html, css, javascript, or a source map json

output will be in the folder named the name argument

report bugs in the issues ♥

# how it work

`html` to `javascript` and `css`

`javascript` to `source map json`

`css` to `source map json`

`source map json` to `original files`