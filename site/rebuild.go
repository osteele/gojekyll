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

// reloads and rebuilds the site; returns a copy and count
func (s *Site) rebuild(paths []string) (*Site, int, error) {
	if s.requiresFullReload(paths) {
		r, err := s.Reloaded(paths)
		if err != nil {
			return nil, 0, err
		}
		n, err := r.Build()
		return r, n, err
	}
	return s, 0, nil
}

func (s *Site) requiresFullReload(paths []string) bool {
	for _, path := range paths {
		switch {
		case s.config.IsConfigPath(path):
			return true
		case s.Exclude(path):
			return false
		case !s.config.Incremental:
			return true
		case strings.HasPrefix(path, s.config.DataDir):
			return true
		case strings.HasPrefix(path, s.config.LayoutsDir):
			return true
		}
	}
	return false
}

// relativize and de-dup filenames, and filter to those that affect the build
func (s *Site) affectsBuildFilter(filenames []string) []string {
	var (
		result = make([]string, 0, len(filenames))
		seen   = map[string]bool{}
	)
loop:
	for _, path := range filenames {
		switch {
		case s.config.IsConfigPath(path):
		case s.Exclude(path):
			continue loop
		case seen[path]:
			continue loop
		}
		seen[path] = true
		result = append(result, path)
	}
	return result
}
