package site

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/osteele/gojekyll/utils"
)

// WatchFiles calls fn repeatedly in a goroutine when
// files change
func (s *Site) WatchFiles(fn func([]string)) error {
	var (
		sourceDir = s.SourceDir()
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
				fmt.Fprintln(os.Stderr, "error:", err)
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
func (s *Site) sitePaths(filenames []string) []string {
	var (
		dir   = s.SourceDir()
		paths = make([]string, 0, len(filenames))
		seen  = map[string]bool{}
	)
	for _, filename := range filenames {
		path := utils.MustRel(dir, filename)
		if path == "_config.yml" || !s.Exclude(path) {
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
