package site

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/osteele/gojekyll/pages"
	"github.com/osteele/gojekyll/utils"
)

var excludeFileRE = regexp.MustCompile(`^[#~]|^\..|~$`)

// Exclude returns a boolean indicating that the site configuration excludes a file or directory.
// It does not exclude top-level _underscore files and directories.
func (s *Site) Exclude(siteRel string) bool {
	for siteRel != "." {
		dir, base := filepath.Dir(siteRel), filepath.Base(siteRel)
		switch {
		case utils.MatchList(s.cfg.Include, siteRel):
			return false
		case utils.MatchList(s.cfg.Exclude, siteRel):
			return true
		case dir != "." && base[0] == '_':
			return true
		default:
			if excludeFileRE.MatchString(base) {
				return true
			}
		}
		siteRel = dir
	}
	return false
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
		case s.cfg.IsConfigPath(path):
			return true
		case s.Exclude(path):
			continue
		case !s.cfg.Incremental:
			return true
		case strings.HasPrefix(path, s.cfg.DataDir):
			return true
		case strings.HasPrefix(path, s.cfg.IncludesDir):
			return true
		case strings.HasPrefix(path, s.cfg.LayoutsDir):
			return true
		case strings.HasPrefix(path, s.cfg.SassDir()):
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
		case s.cfg.IsConfigPath(path):
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
		case utils.MatchList(s.cfg.Include, rel):
			return true
		case utils.MatchList(s.cfg.Exclude, rel):
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
	rel := utils.MustRel(s.SourceDir(), d.Source())
	return paths[rel]
}
