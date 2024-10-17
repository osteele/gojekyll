package commands

import (
	"github.com/osteele/gojekyll/server"
	"github.com/osteele/gojekyll/site"
)

var (
	serve = app.Command("serve", "Serve your site locally").Alias("server").Alias("s")
	open  = serve.Flag("open-url", "Launch your site in a browser").Short('o').Bool()
	_     = serve.Flag("host", "Host to bind to").Short('H').Action(stringVar("host", &options.Host)).String()
	_     = serve.Flag("port", "Port to listen on").Short('P').Action(intVar("port", &options.Port)).Int()
)

func serveCommand(site *site.Site) error {
	server := server.Server{Site: site}
	return server.Run(*open, func(label, value string) {
		//nolint:govet
		logger.label(label, value)
	})
}
