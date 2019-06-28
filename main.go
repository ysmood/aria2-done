package main

import (
	"os"
	"regexp"

	kit "github.com/ysmood/gokit"
	yaml "gopkg.in/yaml.v2"
)

// Rules ...
type Rules map[string][]string

// Context ...
type Context struct {
	rules    Rules
	filePath string
}

func main() {
	ctx := new()

	ctx.move()
}

func new() *Context {
	var rules Rules

	confData, err := kit.ReadFile(os.Getenv("aria2_done_conf"))
	kit.E(err)
	kit.E(yaml.Unmarshal(confData, &rules))

	return &Context{
		rules:    rules,
		filePath: os.Args[3],
	}
}

func (ctx *Context) move() {
	for pattern, tpl := range ctx.rules {
		p := regexp.MustCompile(pattern)
		from := p.ReplaceAllString(ctx.filePath, tpl[0])
		to := p.ReplaceAllString(ctx.filePath, tpl[1])

		if from != to && !kit.FileExists((to)) {
			kit.Log("[aria2-done] move:", from, "->", to)
			kit.E(kit.Move(from, to, nil))
			return
		}
	}
}
