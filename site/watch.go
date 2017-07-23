package site

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/osteele/gojekyll/utils"
	"github.com/radovskyb/watcher"
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
	filenames, err := s.makeFileWatcher()
	if err != nil {
		return nil, err
	}
	var (
		debounced = debounce(time.Second/2, filenames)
		filesets  = make(chan FilesEvent)
	)
	go func() {
		for {
			paths := s.sitePaths(<-debounced)
			if len(paths) > 0 {
				// Create a new timestamp. Except under pathological
				// circumstances, it will be close enough.
				filesets <- FilesEvent{time.Now(), paths}
			}
		}
	}()
	return filesets, nil
}

// WatchRebuild watches the directory. Each time a file changes, it
// rebuilds the site. It sends status messages and error to its output
// channel.
func (s *Site) WatchRebuild() (<-chan interface{}, error) {
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

		// similar code to server.reload
		events <- fmt.Sprintf("Regenerating: %s...", change)
		start := time.Now()
		r, count, err := s.rebuild(change.Paths)
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

func (s *Site) makePollingWatcher() (<-chan string, error) {
	var (
		sourceDir = utils.MustAbs(s.SourceDir())
		filenames = make(chan string, 100)
		w         = watcher.New()
	)
	if err := w.AddRecursive(sourceDir); err != nil {
		return nil, err
	}
	for _, path := range s.config.Exclude {
		if err := w.Ignore(filepath.Join(sourceDir, path)); err != nil {
			return nil, err
		}
	}
	if err := w.Ignore(s.DestDir()); err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case event := <-w.Event:
				filenames <- utils.MustRel(sourceDir, event.Path)
			case err := <-w.Error:
				fmt.Fprintln(os.Stderr, "error:", err)
			case <-w.Closed:
				return
			}
		}
	}()
	go func() {
		if err := w.Start(time.Millisecond * 250); err != nil {
			log.Fatal(err)
		}
	}()
	return filenames, nil
}

func (s *Site) makeEventWatcher() (<-chan string, error) {
	var (
		sourceDir = s.SourceDir()
		filenames = make(chan string, 100)
		w, err    = fsnotify.NewWatcher()
	)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case event := <-w.Events:
				filenames <- utils.MustRel(sourceDir, event.Name)
			case err := <-w.Errors:
				fmt.Fprintln(os.Stderr, "error:", err)
			}
		}
	}()
	return filenames, w.Add(sourceDir)
}

func (s *Site) makeFileWatcher() (<-chan string, error) {
	if s.config.ForcePolling {
		return s.makePollingWatcher()
	}
	return s.makeEventWatcher()
}

// reloads and rebuilds the site; returns a copy and count
func (s *Site) rebuild(paths []string) (r *Site, n int, err error) {
	r, err = s.Reloaded(paths)
	if err != nil {
		return
	}
	n, err = r.Build()
	return
}

// relativize and de-dup filenames, and filter to those that affect the build
func (s *Site) sitePaths(filenames []string) []string {
	var (
		paths = make([]string, 0, len(filenames))
		seen  = map[string]bool{}
	)
	for _, path := range filenames {
		if path == "_config.yml" || !s.Exclude(path) {
			if !seen[path] {
				seen[path] = true
				paths = append(paths, path)
			}
		}
	}
	return paths
}

// debounce relays values from input to output, merging successive values so long as they keep changing
// faster than interval
// TODO consider https://github.com/ReactiveX/RxGo
func debounce(interval time.Duration, input <-chan string) <-chan []string {
	var (
		pending = []string{}
		output  = make(chan []string)
		ticker  <-chan time.Time
	)
	go func() {
		for {
			select {
			case value := <-input:
				if value == "." {
					continue
				}
				pending = append(pending, value)
				ticker = time.After(interval) // replaces the previous ticker
			case <-ticker:
				ticker = nil
				if len(pending) > 0 {
					output <- pending
					pending = []string{}
				}
			}
		}
	}()
	return output
}
