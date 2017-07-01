package server

import (
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

func (s *Server) watchFiles() error {
	var (
		site      = s.Site
		events    = make(chan string)
		debounced = debounce(time.Second, events)
	)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				events <- event.Name
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	go func() {
		for {
			filenames := <-debounced
			// Resolve to URLS *before* reloading the site, in case the latter
			// remaps permalinks.
			urls := map[string]bool{}
			for _, filename := range filenames {
				relpath, err := filepath.Rel(site.SourceDir(), filename)
				if err != nil {
					log.Println("error:", err)
					continue
				}
				url, found := site.FilenameURLPath(relpath)
				if !found {
					// TODO don't warn re config and excluded files
					log.Println("error:", filename, "does not match a site URL")
				}
				urls[url] = true
			}
			s.reloadSite()
			for url := range urls {
				s.lr.Reload(url)
			}
		}
	}()

	return watcher.Add(site.SourceDir())
}

// debounce relays values from input to output, merging successive values within interval
// TODO consider https://github.com/ReactiveX/RxGo
func debounce(interval time.Duration, input <-chan string) <-chan []string {
	output := make(chan []string)
	var (
		pending = []string{}
		ticker  = time.Tick(interval) // nolint: staticcheck, megacheck
	)
	go func() {
		for {
			select {
			case value := <-input:
				pending = append(pending, value)
			case <-ticker:
				if len(pending) > 0 {
					output <- pending
					pending = []string{}
				}
			}
		}
	}()
	return output
}
