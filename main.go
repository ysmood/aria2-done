package main

import (
	"fmt"
	"os"
	"path"
	"regexp"

	g "github.com/ysmood/gokit"
	yaml "gopkg.in/yaml.v2"
)

/* Yaml Example

mod: 0777

automove:
  name_pattern:
    - /move/to/dir
    - (\d\d).mp4

*/
type config struct {
	Mod      os.FileMode
	AutoMove map[string][]string
}

type context struct {
	conf     config
	filePath string
}

func main() {
	ctx := new()

	ctx.chmod()
	ctx.autoMove()
}

func (ctx context) chmod() {
	os.Chmod(ctx.filePath, ctx.conf.Mod)
}

func (ctx context) autoMove() {
	for k, c := range ctx.conf.AutoMove {
		if ctx.match(k) {
			ctx.moveByIndex(c[0], c[1])
		}
	}
}

func new() *context {
	var cs config

	confData, err := g.ReadFile(os.Getenv("aria2_done_conf"))
	g.E(err)
	g.E(yaml.Unmarshal(confData, &cs))

	return &context{
		conf:     cs,
		filePath: os.Args[3],
	}
}

func (ctx *context) match(pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(ctx.filePath)
}

func (ctx *context) index(pattern string) string {
	return regexp.MustCompile(pattern).FindStringSubmatch(ctx.filePath)[1]
}

func (ctx *context) move(target string) {
	g.E(g.Move(ctx.filePath, target, nil))
	g.Log("aria2-done mv:", ctx.filePath, "->", target)
}

func (ctx *context) moveByIndex(dir, pattern string) {
	ctx.move(
		path.Join(
			dir,
			fmt.Sprintf("%s%s", ctx.index(pattern), path.Ext(ctx.filePath)),
		),
	)
}
