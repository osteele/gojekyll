package site

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/utils"
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
	if s.RequiresFullReload(paths) {
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
	inflect := map[bool]string{true: "", false: "s"}[count == 1]
	messages <- fmt.Sprintf("wrote %d file%s in %.2fs.\n", count, inflect, elapsed.Seconds())
	return r
}

// reloads and rebuilds the site; returns a copy and count
func (s *Site) rebuild(paths []string) (r *Site, n int, err error) {
	if s.RequiresFullReload(paths) {
		r, err = s.Reloaded(paths)
		if err != nil {
			return
		}
		n, err = r.Write()
		return
	}
	r = s
	pathSet := utils.StringSet(paths)
	for _, d := range s.docs {
		if s.invalidatesDoc(pathSet, d) {
			err = d.Reload()
			if err != nil {
				return
			}
			err = s.WriteDoc(d)
			if err != nil {
				return
			}
			n++
		}
	}
	return
}

// RequiresFullReload returns true if a source file requires a full reload / rebuild.
//
// This is always true outside of incremental mode, since even a
// static asset can cause pages to change if they reference its
// variables.
//
// This function works on relative paths. It does not work for theme
// sources.
func (s *Site) RequiresFullReload(paths []string) bool {
	for _, path := range paths {
		switch {
		case s.config.IsConfigPath(path):
			return true
		case s.Exclude(path):
			continue
		case !s.config.Incremental:
			return true
		case strings.HasPrefix(path, s.config.DataDir):
			return true
		case strings.HasPrefix(path, s.config.IncludesDir):
			return true
		case strings.HasPrefix(path, s.config.LayoutsDir):
			return true
		case strings.HasPrefix(path, s.config.SassDir()):
			return true
		}
	}
	return false
}

// De-dup relative paths, and filter to those that might affect the build.
//
// Site watch uses this to decide when to send events.
func (s *Site) affectsBuildFilter(paths []string) []string {
	var (
		result = make([]string, 0, len(paths))
		seen   = map[string]bool{}
	)
loop:
	for _, path := range paths {
		switch {
		case s.config.IsConfigPath(path):
			// break
		case !s.fileAffectsBuild(path):
			continue loop
		case seen[path]:
			continue loop
		}
		result = append(result, path)
		seen[path] = true
	}
	return result
}

// Returns true if the file or a parent directory is excluded.
// Cf. Site.Exclude.
func (s *Site) fileAffectsBuild(rel string) bool {
	for rel != "" {
		switch {
		case rel == ".":
			return true
		case utils.MatchList(s.config.Include, rel):
			return true
		case utils.MatchList(s.config.Exclude, rel):
			return false
		case strings.HasPrefix(rel, "."):
			return false
		}
		rel = filepath.Dir(rel)
	}
	return true
}

// returns true if changes to the site-relative paths invalidate doc
func (s *Site) invalidatesDoc(paths map[string]bool, d pages.Document) bool {
	rel := utils.MustRel(s.SourceDir(), d.SourcePath())
	return paths[rel]
}
