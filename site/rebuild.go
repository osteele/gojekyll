package site

import (
	"fmt"
	"reflect"
	"time"

	"github.com/osteele/gojekyll/logger"
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
		return s.fullyReloaded()
	}
	_, requiresFullReload, err := s.reloadDocuments(paths)
	if err != nil {
		return nil, err
	}
	if requiresFullReload {
		return s.fullyReloaded()
	}
	return s, nil
}

func (s *Site) fullyReloaded() (*Site, error) {
	reloaded, err := FromDirectory(s.SourceDir(), s.flags)
	if err != nil {
		return nil, err
	}
	return reloaded, reloaded.Read()
}

func (s *Site) processFilesEvent(fileset FilesEvent, messages chan<- interface{}) *Site {
	// similar code to server.reload
	log := logger.Default()
	messages <- fmt.Sprintf("Regenerating: %s...", fileset)
	start := time.Now()
	r, count, err := s.rebuild(fileset.Paths)
	if err != nil {
		log.Println()
		log.Error("%s", err.Error())
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
		r, err = s.fullyReloaded()
		if err != nil {
			return
		}
		n, err = r.Write()
		return
	}
	r = s
	var docs []Document
	docs, requiresFullReload, reloadErr := s.reloadDocuments(paths)
	if reloadErr != nil {
		err = reloadErr
		return
	}
	if requiresFullReload {
		r, err = s.fullyReloaded()
		if err == nil {
			n, err = r.Write()
		}
		return
	}
	for _, d := range docs {
		err = s.WriteDoc(d)
		if err != nil {
			return
		}
		n++
	}
	return
}

func (s *Site) reloadDocuments(paths []string) ([]Document, bool, error) {
	pathSet := utils.MakeStringSet(paths)
	var changed []Document
	for _, d := range s.docs {
		if !s.invalidatesDoc(pathSet, d) {
			continue
		}
		var oldFrontMatter interface{}
		if page, ok := d.(Page); ok {
			oldFrontMatter = page.FrontMatter()
		}
		oldURL := d.URL()
		oldPublished := d.Published()
		if err := d.Reload(); err != nil {
			return nil, false, err
		}
		changed = append(changed, d)
		if oldURL != d.URL() || oldPublished != d.Published() {
			return changed, true, nil
		}
		if page, ok := d.(Page); ok && !reflect.DeepEqual(oldFrontMatter, page.FrontMatter()) {
			return changed, true, nil
		}
	}
	return changed, false, nil
}
