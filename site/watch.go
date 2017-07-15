package site

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/osteele/gojekyll/utils"
)

// FilesEvent is a list of changed or added site source files, with a single
// timestamp that approximates when they were changed.
type FilesEvent struct {
	Time  time.Time // A single time is used for all the changes
	Paths []string  // relative to site source
}

func (e FilesEvent) String() string {
	count := len(e.Paths)
	inflect := map[bool]string{true: "", false: "s"}[count == 1]
	return fmt.Sprintf("%d file%s changed at %s", count, inflect, e.Time.Format("3:04:05PM"))
}

// WatchFiles sends FilesEvent on changes within the site directory.
func (s *Site) WatchFiles() (<-chan FilesEvent, error) {
	var (
		sourceDir = s.SourceDir()
		events    = make(chan string)
		debounced = debounce(time.Second, events)
		results   = make(chan FilesEvent)
	)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
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
				// Make up a new timestamp. Except under pathological
				// circumstances, it will be close enough.
				results <- FilesEvent{time.Now(), paths}
			}
		}
	}()
	return results, watcher.Add(sourceDir)
}

// WatchRebuild watches the directory. Each time a file changes, it
// rebuilds the site. It sends status messages and error to its output
// channel.
func (s *Site) WatchRebuild(o BuildOptions) (<-chan interface{}, error) {
	var (
		mu           sync.Mutex
		events       = make(chan interface{})
		changes, err = s.WatchFiles()
	)
	if err != nil {
		return nil, err
	}
	go func(rebuild func(FilesEvent)) {
		for change := range changes {
			rebuild(change)
		}
	}(func(change FilesEvent) {
		mu.Lock()
		defer mu.Unlock()

		events <- fmt.Sprintf("Regenerating: %s...", change)
		start := time.Now()
		r, count, err := s.rebuild(o)
		if err != nil {
			fmt.Println()
			fmt.Fprintln(os.Stderr, err)
			return
		}
		// use the new site value the next time
		s = r
		elapsed := time.Since(start)
		events <- fmt.Sprintf("wrote %d files in %.2fs.\n", count, elapsed.Seconds())
	})
	return events, nil
}

// reloads and rebuilds the site; returns a copy and count
func (s *Site) rebuild(o BuildOptions) (r *Site, n int, err error) {
	r, err = s.Reloaded()
	if err != nil {
		return
	}
	n, err = r.Build(o)
	return
}

// relativize and de-dup filenames, and filter to those that affect the build
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
