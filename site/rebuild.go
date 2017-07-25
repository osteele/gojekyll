package site

import (
	"fmt"
	"os"
	"strings"
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

// Reloaded returns the same or a new site reading the same source directory, configuration file, and load flags.
// build --incremental and site --incremental use this.
func (s *Site) Reloaded(paths []string) (*Site, error) {
	if s.requiresFullReload(paths) {
		copy, err := FromDirectory(s.SourceDir(), s.flags)
		if err != nil {
			return nil, err
		}
		s = copy
	}
	return s, s.Read()
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

func (s *Site) requiresFullReload(paths []string) bool {
	for _, path := range paths {
		switch {
		case path == "_config.yml":
			return true
		case strings.HasPrefix(path, s.config.DataDir):
			return true
		case strings.HasPrefix(path, s.config.LayoutsDir):
			return true
		}
	}
	return false
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
