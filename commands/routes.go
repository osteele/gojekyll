package commands

import (
	"fmt"
	"sort"

	"github.com/osteele/gojekyll/site"
)

var routes = app.Command("routes", "Display site permalinks and associated files")
var dynamicRoutes = routes.Flag("dynamic", "Only show routes to non-static files").Bool()

func routesCommand(site *site.Site) error {
	logger.label("Routes:", "")
	var urls []string
	for u, p := range site.Routes {
		if !*dynamicRoutes || !p.IsStatic() {
			urls = append(urls, u)
		}
	}
	sort.Strings(urls)
	for _, u := range urls {
		filename := site.Routes[u].Source()
		fmt.Printf("  %s -> %s\n", u, filename)
	}
	return nil
}
