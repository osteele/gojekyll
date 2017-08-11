package commands

import (
	"fmt"
	"os"
	"time"
)

// Make initializes Version to the git commit hash, and BuildDate.
var (
	Version   string
	BuildDate string
	BuildTime time.Time
)

func init() {
	if Version == "" {
		Version = "develop"
	}
	if BuildDate != "" {
		bd, err := time.Parse("2006-01-02T15:04:05-0700", BuildDate)
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid BuildDate", BuildDate) // nolint: gas
		} else {
			BuildTime = bd.In(time.UTC)
		}
	}
}

var versionCmd = app.Command("version", "Print the name and version")

func versionCommand() error {
	var d string
	if !BuildTime.IsZero() {
		d = BuildTime.Format(" (Build time: 2006-01-02T15:04)")
	}
	fmt.Printf("gojekyll version %s%s\n", Version, d)
	return nil
}
