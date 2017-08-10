package site

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/osteele/gojekyll/utils"
)

// Clean the destination. Remove files that aren't in keep_files, and resulting empty directories.
func (s *Site) Clean() error {
	// TODO PERF when called as part of build, keep files that will be re-generated
	removeFiles := func(filename string, info os.FileInfo, err error) error {
		if s.config.Verbose {
			fmt.Println("rm", filename)
		}
		switch {
		case err != nil && os.IsNotExist(err):
			return nil
		case err != nil:
			return err
		case info.IsDir():
			return nil
		case s.KeepFile(utils.MustRel(s.DestDir(), filename)):
			return nil
		case s.config.DryRun:
			return nil
		default:
			// empirically, moving the os.Remove into a goroutine has no performance benefit
			return os.Remove(filename)
		}
	}
	if err := filepath.Walk(s.DestDir(), removeFiles); err != nil {
		return err
	}
	return utils.RemoveEmptyDirectories(s.DestDir())
}

// Build cleans the destination and create files in it.
// This sets TZ from the site config.
func (s *Site) Build() (int, error) {
	if err := s.setTimeZone(); err != nil {
		return 0, err
	}
	if err := s.ensureRendered(); err != nil {
		return 0, err
	}
	if err := s.Clean(); err != nil {
		return 0, err
	}
	return s.WriteFiles()
}

func (s *Site) setTimeZone() error {
	if tz := s.config.Timezone; tz != "" {
		if _, err := time.LoadLocation(tz); err != nil {
			fmt.Fprintf(os.Stderr, "%s; using local time zone\n", err)
		} else if err := os.Setenv("TZ", tz); err != nil {
			return err
		}
	}
	return nil
}
