package commands

import (
	"fmt"

	"github.com/osteele/gojekyll/version"
)

var versionCmd = app.Command("version", "Print the name and version")

func versionCommand() error {
	var d string
	if !version.BuildTime.IsZero() {
		d = version.BuildTime.Format(" (Build time: 2006-01-02T15:04)")
	}
	fmt.Printf("gojekyll version %s%s\n", version.Version, d)
	return nil
}
