package server

import (
	"fmt"
	"os"
	"time"

	"github.com/osteele/gojekyll/site"
)

// Create a goroutine that rebuilds the site when files change.
func (s *Server) watchReload() error {
	site := s.Site
	// FIXME reload swaps in a new site but we're still watching the old one
	changes, err := site.WatchFiles()
	if err != nil {
		return err
	}
	go func() {
		for change := range changes {
			// Resolves filenames to URLS *before* reloading the site, in case the latter
			// changes the url -> filename routes.
			urls := map[string]bool{}
			for _, rel := range change.Paths {
				url, ok := site.FilenameURLPath(rel)
				if ok {
					urls[url] = true
				}
			}
			// reload the site
			s.reload(change)
			// tell the pages their files (may have) changed
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

	// similar code to site.WatchRebuild
	fmt.Printf("Re-reading: %v...", change)
	start := time.Now()
	site, err := s.Site.Reloaded(change.Paths)
	if err != nil {
		fmt.Println()
		fmt.Fprintln(os.Stderr, err.Error())
		s.lr.Alert(fmt.Sprintf("Error reading site configuration: %s", err))
		return
	}
	s.Site = site
	s.Site.SetAbsoluteURL("")
	fmt.Printf("done (%.2fs)\n", time.Since(start).Seconds())
}
