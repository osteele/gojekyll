package site

import (
	"os"
	"path/filepath"
	"time"

	"github.com/osteele/gojekyll/logger"
	"github.com/osteele/gojekyll/utils"
)

// Clean the destination. Remove files that aren't in keep_files, and resulting empty directories.
func (s *Site) Clean() error {
	// If destination directory doesn't exist, there's nothing to clean
	if _, err := os.Stat(s.DestDir()); os.IsNotExist(err) {
		return nil
	}

	// TODO PERF when called as part of build, keep files that will be re-generated
	log := logger.Default()
	removeFiles := func(filename string, info os.FileInfo, err error) error {
		if s.cfg.Verbose {
			log.Info("rm %s", filename)
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
		case s.cfg.DryRun:
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

func (s *Site) setTimeZone() error {
	log := logger.Default()
	if tz := s.cfg.Timezone; tz != "" {
		if _, err := time.LoadLocation(tz); err != nil {
			log.Error("%s; using local time zone", err)
		} else if err := os.Setenv("TZ", tz); err != nil {
			return err
		}
	}
	return nil
}
