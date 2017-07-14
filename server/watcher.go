package server

import (
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
			paths := s.sitePaths(<-debounced)
			if len(paths) > 0 {
				fn(paths)
			}
		}
	}()
	return watcher.Add(sourceDir)
}

// relativize and de-dup paths and filter to those the site source
func (s *Server) sitePaths(filenames []string) []string {
	var (
		site  = s.Site
		dir   = site.SourceDir()
		paths = make([]string, 0, len(filenames))
		seen  = map[string]bool{}
	)
	for _, filename := range filenames {
		path := utils.MustRel(dir, filename)
		if path == "_config.yml" || !site.Exclude(path) {
			if !seen[path] {
				seen[path] = true
				paths = append(paths, path)
			}
		}
	}
	return paths
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
