package main

import (
	"regexp"

	kit "github.com/ysmood/gokit"
	yaml "gopkg.in/yaml.v2"
)

// Rules ...
type Rules yaml.MapSlice

// Context ...
type Context struct {
	rules    Rules
	filePath string
}

func main() {
	app := kit.TasksNew("aria2-done", "aria2 hook handler when download is done")
	app.Version("v0.0.1")
	kit.Tasks().App(app).Add(
		kit.Task("do", "").Init(func(cmd kit.TaskCmd) func() {
			cmd.Default()
			conf := cmd.Flag("conf", "yaml config file path").Envar("aria2_done_conf").Default("aria2-done.yml").String()
			cmd.Arg("gid", "").Required().String()
			cmd.Arg("number", "the number of files").Required().String()
			filePath := cmd.Arg("file-path", " file path").Required().String()

			return func() {
				ctx := new(*conf, *filePath)
				ctx.move()
			}
		}),
	).Do()
}

func new(confPath, filePath string) *Context {
	var rules Rules

	confData, err := kit.ReadFile(confPath)
	kit.E(err)
	kit.E(yaml.Unmarshal(confData, &rules))

	return &Context{
		rules:    rules,
		filePath: filePath,
	}
}

func (ctx *Context) move() {
	for _, item := range ctx.rules {
		pattern := item.Key.(string)
		tpl := item.Value.([]interface{})
		kit.Log(pattern)
		p := regexp.MustCompile(pattern)
		from := p.ReplaceAllString(ctx.filePath, tpl[0].(string))
		to := p.ReplaceAllString(ctx.filePath, tpl[1].(string))

		if from != to && !kit.FileExists((to)) {
			kit.Log("[aria2-done] move:", from, "->", to)
			// kit.E(kit.Move(from, to, nil))
			return
		}
	}
}
