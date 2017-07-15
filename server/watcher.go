package server

import (
	"fmt"
	"os"
	"time"

	"github.com/osteele/gojekyll/site"
)

// Create a goroutine that rebuilds the site when files change
func (s *Server) watchAndReload() error {
	site := s.Site
	changes, err := site.WatchFiles()
	if err != nil {
		return err
	}
	go func() {
		for change := range changes {
			// Resolves filenames to URLS *before* reloading the site, in case the latter
			// changes the url -> filename routes.
			urls := map[string]bool{}
			for _, relpath := range change.Paths {
				url, ok := site.FilenameURLPath(relpath)
				if ok {
					urls[url] = true
				}
			}
			s.reload(change)
			for url := range urls {
				s.lr.Reload(url)
			}
		}
	}()
	return nil
}

func (s *Server) reload(change site.FilesEvent) {
	s.Lock()
	defer s.Unlock()

	// DRY w/ site.WatchRebuild
	start := time.Now()
	fmt.Printf("Re-reading: %v...", change)
	site, err := s.Site.Reloaded()
	if err != nil {
		fmt.Println()
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	s.Site = site
	s.Site.SetAbsoluteURL("")
	fmt.Printf("done (%.2fs)\n", time.Since(start).Seconds())
}
