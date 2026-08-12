package site

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func canonicalPath(name string) (string, error) {
	abs, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}
	abs = filepath.Clean(abs)
	candidate := abs
	var suffix []string
	for {
		resolved, err := filepath.EvalSymlinks(candidate)
		if err == nil {
			parts := append([]string{resolved}, suffix...)
			return filepath.Clean(filepath.Join(parts...)), nil
		}
		if !os.IsNotExist(err) {
			return "", err
		}
		parent := filepath.Dir(candidate)
		if parent == candidate {
			return abs, nil
		}
		suffix = append([]string{filepath.Base(candidate)}, suffix...)
		candidate = parent
	}
}

func pathContains(root, candidate string) bool {
	rel, err := filepath.Rel(root, candidate)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

func (s *Site) validateDestination() error {
	source, err := canonicalPath(s.SourceDir())
	if err != nil {
		return err
	}
	destination, err := canonicalPath(s.DestDir())
	if err != nil {
		return err
	}
	volumeRoot := filepath.Clean(filepath.VolumeName(destination) + string(filepath.Separator))
	if destination == volumeRoot {
		return fmt.Errorf("unsafe destination %q: refusing to clean a filesystem root", s.DestDir())
	}
	if pathContains(destination, source) {
		return fmt.Errorf("unsafe destination %q: it contains the source directory %q", s.DestDir(), s.SourceDir())
	}
	if home, err := os.UserHomeDir(); err == nil {
		if canonicalHome, err := canonicalPath(home); err == nil && destination == canonicalHome {
			return fmt.Errorf("unsafe destination %q: refusing to clean the home directory", s.DestDir())
		}
	}
	return nil
}

func (s *Site) isDestinationPath(name string) bool {
	destination, err := canonicalPath(s.DestDir())
	if err != nil {
		return false
	}
	candidate, err := canonicalPath(name)
	return err == nil && pathContains(destination, candidate)
}
