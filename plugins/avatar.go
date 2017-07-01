package plugins

import (
	"fmt"
	"io"
	"strings"

	"github.com/osteele/gojekyll/tags"
	"github.com/osteele/liquid/chunks"
)

func init() {
	registerPlugin("jekyll-avatar", func(ctx PluginContext) error {
		ctx.TemplateEngine().DefineTag("avatar", avatarTag)
		return nil
	})
}

const avatarTemplate = `<img class="avatar avatar-small" src="https://avatars3.githubusercontent.com/{user}?v=3&amp;s=40" alt="{user}" srcset="https://avatars3.githubusercontent.com/{user}?v=3&amp;s=40 1x, https://avatars3.githubusercontent.com/{user}?v=3&amp;s=80 2x, https://avatars3.githubusercontent.com/{user}?v=3&amp;s=120 3x, https://avatars3.githubusercontent.com/{user}?v=3&amp;s=160 4x" width="40" height="40" data-proofer-ignore="true" />`

func avatarTag(_ string) (func(io.Writer, chunks.RenderContext) error, error) {
	return func(w io.Writer, ctx chunks.RenderContext) error {
		var (
			user string
			size interface{} = 40
		)
		argsline, err := ctx.ParseTagArgs()
		if err != nil {
			return err
		}
		args, err := tags.ParseArgs(argsline)
		if err != nil {
			return err
		}
		if len(args.Args) > 0 {
			user = args.Args[0]
		}
		options, err := args.EvalOptions(ctx)
		if err != nil {
			return err
		}
		for name, value := range options {
			switch name {
			case "user":
				user = fmt.Sprint(value)
			case "size":
				size = value
			default:
				return fmt.Errorf("unknown avatar argument: %s", name)
			}
		}
		if user == "" {
			return fmt.Errorf("parse error in avatar tag parameters %s", argsline)
		}
		s := strings.Replace(avatarTemplate, "40", fmt.Sprint(size), -1)
		s = strings.Replace(s, "{user}", user, -1)
		_, err = w.Write([]byte(s))
		return err
	}, nil
}
