package plugins

import (
	"bytes"
	"fmt"
	"hash/crc32"
	"strings"
	"text/template"

	"github.com/osteele/gojekyll/tags"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
)

func init() {
	register("jekyll-avatar", jekyllAvatarPlugin{})
}

type jekyllAvatarPlugin struct{ plugin }

func (p jekyllAvatarPlugin) ConfigureTemplateEngine(e *liquid.Engine) error {
	e.RegisterTag("avatar", avatarTag)
	return nil
}

var avatarTemplate = template.Must(template.New("avatar").Parse(strings.TrimSpace(`
<img class="avatar avatar-small" src="https://{{.Subdomain}}.githubusercontent.com/{{.User}}?v=3&amp;s={{.Size}}" alt="{{.User}}" srcset="https://{{.Subdomain}}.githubusercontent.com/{{.User}}?v=3&amp;s={{.Size}} 1x, https://{{.Subdomain}}.githubusercontent.com/{{.User}}?v=3&amp;s=80 2x, https://{{.Subdomain}}.githubusercontent.com/{{.User}}?v=3&amp;s=120 3x, https://{{.Subdomain}}.githubusercontent.com/{{.User}}?v=3&amp;s=160 4x" width="{{.Size}}" height="{{.Size}}" data-proofer-ignore="true" />
`)))

func avatarTag(ctx render.Context) (string, error) {
	var (
		user string
		size interface{} = 40
	)
	argsline, err := ctx.ExpandTagArg()
	if err != nil {
		return "", err
	}
	args, err := tags.ParseArgs(argsline)
	if err != nil {
		return "", err
	}
	if len(args.Args) > 0 {
		user = args.Args[0]
	}
	options, err := args.EvalOptions(ctx)
	if err != nil {
		return "", err
	}
	for name, value := range options {
		switch name {
		case "user":
			user = fmt.Sprint(value)
		case "size":
			size = value
		default:
			return "", fmt.Errorf("unknown avatar argument: %s", name)
		}
	}
	if user == "" {
		return "", fmt.Errorf("parse error in avatar tag parameters %s", argsline)
	}
	n := crc32.Checksum([]byte(fmt.Sprintf("%s:%d", user, size)), crc32.IEEETable) % 4
	s := struct {
		User, Subdomain string
		Size            interface{}
	}{user, fmt.Sprintf("avatar%d", n), size}
	buf := new(bytes.Buffer)
	err = avatarTemplate.Execute(buf, s)
	return buf.String(), err
}
