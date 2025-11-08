package site

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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

// WatchFiles returns a channel that receives FilesEvent on changes within the site directory.
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
			paths := s.affectsBuildFilter(<-debounced)
			if len(paths) > 0 {
				// Create a new timestamp. Except under pathological
				// circumstances, it will be close enough.
				filesets <- FilesEvent{time.Now(), paths}
			}
		}
	}()
	return filesets, nil
}

func (s *Site) makeFileWatcher() (<-chan string, error) {
	if s.cfg.ForcePolling {
		return s.makePollingWatcher()
	}

	// Try fsnotify first, but fall back to polling if too many directories
	filenames, err := s.tryEventWatcher()
	if err != nil {
		if s.cfg.Verbose {
			fmt.Fprintf(os.Stderr, "Event watcher unavailable (%v), using polling watcher\n", err)
		}
		return s.makePollingWatcher()
	}
	return filenames, nil
}

// shouldIgnoreDir returns true for directories that should never be watched
func (s *Site) shouldIgnoreDir(rel string) bool {
	// Always ignore version control directories
	if rel == ".git" || rel == ".svn" || rel == ".hg" || rel == ".bzr" {
		return true
	}

	// Check configured excludes
	for _, excl := range s.cfg.Exclude {
		if rel == excl || strings.HasPrefix(rel, excl+string(filepath.Separator)) {
			return true
		}
	}

	return false
}

// tryEventWatcher attempts to create an fsnotify watcher, returning an error
// if there are too many directories (to avoid exhausting file descriptors)
func (s *Site) tryEventWatcher() (<-chan string, error) {
	sourceDir := s.SourceDir()

	// Count directories first to avoid exhausting watches
	dirCount := 0
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		rel := utils.MustRel(sourceDir, path)
		if path == s.DestDir() || s.shouldIgnoreDir(rel) {
			return filepath.SkipDir
		}

		dirCount++
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Conservative limit to avoid exhausting inotify watches
	// (default limit is often 8192, but other programs use them too)
	const maxDirs = 500
	if dirCount > maxDirs {
		return nil, fmt.Errorf("directory count %d exceeds safe limit %d for fsnotify", dirCount, maxDirs)
	}

	return s.makeEventWatcher()
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

	// Add all subdirectories recursively since fsnotify doesn't watch recursively
	addRecursive := func(dir string) error {
		return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				// Skip excluded directories and destination directory
				rel := utils.MustRel(sourceDir, path)
				if path == s.DestDir() || s.shouldIgnoreDir(rel) {
					return filepath.SkipDir
				}
				if err := w.Add(path); err != nil {
					return err
				}
			}
			return nil
		})
	}

	if err := addRecursive(sourceDir); err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case event := <-w.Events:
				// When a directory is created, add it to the watcher
				if event.Op&fsnotify.Create != 0 {
					if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
						addRecursive(event.Name)
					}
				}
				filenames <- utils.MustRel(sourceDir, event.Name)
			case err := <-w.Errors:
				fmt.Fprintln(os.Stderr, "error:", err)
			}
		}
	}()
	return filenames, nil
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
	for _, path := range s.cfg.Exclude {
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
