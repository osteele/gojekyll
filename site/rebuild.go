package site

import (
	"fmt"
	"os"
	"time"
)

// WatchRebuild watches the site directory. Each time a file changes, it
// rebuilds the site. It sends status messages (strings) and errors to its output
// channel.
//
// TODO use a logger instead of a message channel?
func (s *Site) WatchRebuild() (<-chan interface{}, error) {
	var (
		messages      = make(chan interface{})
		filesets, err = s.WatchFiles()
	)
	if err != nil {
		return nil, err
	}
	go func() {
		for fileset := range filesets {
			s = s.processFilesEvent(fileset, messages)
		}
	}()
	return messages, nil
}

func (s *Site) processFilesEvent(fileset FilesEvent, messages chan<- interface{}) *Site {
	// similar code to server.reload
	messages <- fmt.Sprintf("Regenerating: %s...", fileset)
	start := time.Now()
	r, count, err := s.rebuild(fileset.Paths)
	if err != nil {
		fmt.Println()
		fmt.Fprintln(os.Stderr, err)
		return s
	}
	elapsed := time.Since(start)
	messages <- fmt.Sprintf("wrote %d files in %.2fs.\n", count, elapsed.Seconds())
	return r
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
