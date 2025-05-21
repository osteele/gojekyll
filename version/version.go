package version

import (
	"fmt"
	"os"
	"time"
)

// "make" initializes Version to the git commit hash, and BuildDate.
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
			fmt.Fprintln(os.Stderr, "invalid BuildDate", BuildDate)
		} else {
			BuildTime = bd.In(time.UTC)
		}
	}
}
