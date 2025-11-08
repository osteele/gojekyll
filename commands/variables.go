package commands

import (
	"strings"

	"github.com/k0kubun/pp"
	"github.com/osteele/gojekyll/site"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
)

var variables = app.Command(
	"variables",
	"Print site or document variables",
).Alias("v").Alias("var").Alias("vars")

var variablePath = variables.Arg("PATH", `Filename, URL, "site", or e.g. "site.x.y"`).String()

func variablesCommand(site *site.Site) (err error) {
	var data interface{}
	switch {
	case strings.HasPrefix(*variablePath, "site"):
		data, err = utils.FollowDots(site, strings.Split(*variablePath, ".")[1:])
		if err != nil {
			return
		}
	case *variablePath != "":
		data, err = pageFromPathOrRoute(site, *variablePath)
		if err != nil {
			return
		}
	default:
		data = site
	}
	data = liquid.FromDrop(data)
	bytesToStrings(data)
	bannerLog.label("Variables:", "")
	_, err = pp.Print(data)
	return err
}

// modifies its argument
func bytesToStrings(data interface{}) {
	if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			if b, ok := v.([]byte); ok {
				s := string(b)
				if len(s) > 200 {
					s = s[:200] + "…"
				}
				m[k] = s
			}
		}
	}
}
