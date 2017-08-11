package commands

import (
	"fmt"
	"strings"

	"github.com/osteele/gojekyll/site"
	"github.com/osteele/gojekyll/utils"
	"github.com/osteele/liquid"
	yaml "gopkg.in/yaml.v1"
)

var variables = app.Command("variables", "Display a file or URL path's variables").Alias("v").Alias("var").Alias("vars")
var variablePath = variables.Arg("PATH", "Path, URL, site, or site...").String()

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
	if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			if b, ok := v.([]byte); ok {
				m[k] = string(b)
			}
		}
	}
	b, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	logger.label("Variables:", "")
	fmt.Println(string(b))
	return nil
}
