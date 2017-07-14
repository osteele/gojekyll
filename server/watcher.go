package server

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/osteele/gojekyll/utils"
)

func (s *Server) watchAndReload() error {
	site := s.Site
	return s.watchFiles(func(filenames []string) {
		// This resolves filenames to URLS *before* reloading the site, in case the latter
		// remaps permalinks.
		urls := map[string]bool{}
		for _, relpath := range filenames {
			url, ok := site.FilenameURLPath(relpath)
			if ok {
				urls[url] = true
			}
		}
		fmt.Println(filenames)
		s.reloadSite(len(filenames))
		for url := range urls {
			s.lr.Reload(url)
		}
	})
}

// calls fn repeatedly in a goroutine
func (s *Server) watchFiles(fn func([]string)) error {
	var (
		site      = s.Site
		sourceDir = site.SourceDir()
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
			relpaths := make([]string, 0, len(filenames))
			seen := map[string]bool{}
			for _, filename := range filenames {
				relpath := utils.MustRel(sourceDir, filename)
				if relpath == "_config.yml" || !site.Exclude(relpath) {
					if !seen[relpath] {
						seen[relpath] = true
						relpaths = append(relpaths, relpath)
					}
				}
			}
			fn(relpaths)
		}
	}()
	return watcher.Add(sourceDir)
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
