package plugins

import (
	"fmt"

	"github.com/osteele/gojekyll/tags"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
)

func init() {
	register("jekyll-gist", jekyllGistPlugin{})
}

type jekyllGistPlugin struct{ plugin }

func (p jekyllGistPlugin) ConfigureTemplateEngine(e *liquid.Engine) error {
	e.RegisterTag("gist", gistTag)
	return nil
}

func gistTag(ctx render.Context) (string, error) {
	argsline, err := ctx.ExpandTagArg()
	if err != nil {
		return "", err
	}
	args, err := tags.ParseArgs(argsline)
	if err != nil {
		return "", err
	}
	if len(args.Args) < 1 {
		return "", fmt.Errorf("gist tag: missing argument")
	}
	url := fmt.Sprintf("https://gist.github.com/%s.js", args.Args[0])
	if len(args.Args) >= 2 {
		url += fmt.Sprintf("?file=%s", args.Args[1])
	}
	return `<script src=` + url + `> </script>`, nil
}
