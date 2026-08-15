package site

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/osteele/gojekyll/config"
	"github.com/stretchr/testify/require"
)

// func TestSite_WatchRebuild(t *testing.T) {

func TestSite_Reloaded(t *testing.T) {
	s0 := New(config.Flags{})
	s0.cfg.Incremental = true
	s1, _ := s0.Reloaded([]string{})
	require.Equal(t, s0, s1)

	s1, _ = s0.Reloaded([]string{"_config.yml"})
	require.NotEqual(t, s0, s1)
}

// func TestSite_processFilesEvent(t *testing.T) {
// func TestSite_rebuild(t *testing.T) {

func TestSite_RequiresFullReload(t *testing.T) {
	s := New(config.Flags{})
	require.False(t, s.RequiresFullReload([]string{}))
	require.True(t, s.RequiresFullReload([]string{"file.md"}))
	require.False(t, s.RequiresFullReload([]string{".git"}))
	// require.False(t, s.RequiresFullReload([]string{"_site"}))
	// require.False(t, s.RequiresFullReload([]string{"_site/index.html"}))

	s.cfg.Incremental = true
	require.False(t, s.RequiresFullReload([]string{}))
	require.True(t, s.RequiresFullReload([]string{"file.md"}))
	require.True(t, s.RequiresFullReload([]string{"_config.yml"}))
}

func TestIncrementalReload(t *testing.T) {
	source := t.TempDir()
	filename := filepath.Join(source, "index.md")
	require.NoError(t, os.WriteFile(filename, []byte("---\npermalink: /old/\n---\nold\n"), 0o644))
	s, err := FromDirectory(source, config.Flags{Incremental: boolPointer(true)})
	require.NoError(t, err)
	require.NoError(t, s.Read())

	require.NoError(t, os.WriteFile(filename, []byte("---\npermalink: /old/\n---\nnew\n"), 0o644))
	same, err := s.Reloaded([]string{"index.md"})
	require.NoError(t, err)
	require.Same(t, s, same)

	require.NoError(t, os.WriteFile(filename, []byte("---\npermalink: /new/\n---\nnew\n"), 0o644))
	reloaded, err := same.Reloaded([]string{"index.md"})
	require.NoError(t, err)
	require.NotSame(t, same, reloaded)
	require.Contains(t, reloaded.Routes, "/new/")
	require.NotContains(t, reloaded.Routes, "/old/")
}

func TestIncrementalReloadDetectsAddedAndDeletedFiles(t *testing.T) {
	source := t.TempDir()
	index := filepath.Join(source, "index.md")
	require.NoError(t, os.WriteFile(index, []byte("---\n---\nindex\n"), 0o644))
	s, err := FromDirectory(source, config.Flags{Incremental: boolPointer(true)})
	require.NoError(t, err)
	require.NoError(t, s.Read())

	added := filepath.Join(source, "added.md")
	require.NoError(t, os.WriteFile(added, []byte("---\n---\nadded\n"), 0o644))
	reloaded, err := s.Reloaded([]string{"added.md"})
	require.NoError(t, err)
	require.NotSame(t, s, reloaded)
	require.NotNil(t, reloaded.documentForRelativePath("added.md"))

	require.NoError(t, os.Remove(index))
	reloadedAgain, err := reloaded.Reloaded([]string{"index.md"})
	require.NoError(t, err)
	require.NotSame(t, reloaded, reloadedAgain)
	require.Nil(t, reloadedAgain.documentForRelativePath("index.md"))
}

func boolPointer(value bool) *bool { return &value }

// fakeDocument is a minimal Document implementation for state-table tests.
type fakeDocument struct{ source string }

func (d fakeDocument) URL() string           { return "/" }
func (d fakeDocument) Source() string        { return d.source }
func (d fakeDocument) OutputExt() string     { return ".html" }
func (d fakeDocument) Published() bool       { return true }
func (d fakeDocument) IsStatic() bool        { return false }
func (d fakeDocument) Write(io.Writer) error { return nil }
func (d fakeDocument) Reload() error         { return nil }

// TestRequiresFullReloadStateTable enumerates the configuration/path states that
// drive the incremental-vs-full-reload decision.
func TestRequiresFullReloadStateTable(t *testing.T) {
	source := t.TempDir()
	s := New(config.Flags{})
	s.cfg.Source = source
	require.NoError(t, os.WriteFile(filepath.Join(source, "index.md"), []byte("---\n---\n"), 0o644))
	s.docs = []Document{fakeDocument{source: filepath.Join(source, "index.md")}}

	tests := []struct {
		name     string
		setup    func()
		paths    []string
		expected bool
	}{
		{
			name:     "empty paths",
			paths:    []string{},
			expected: false,
		},
		{
			name:     "config path always requires full reload",
			setup:    func() { s.cfg.Incremental = true },
			paths:    []string{"_config.yml"},
			expected: true,
		},
		{
			name:     "non-incremental normal file",
			setup:    func() { s.cfg.Incremental = false },
			paths:    []string{"file.md"},
			expected: true,
		},
		{
			name:     "non-incremental excluded file",
			setup:    func() { s.cfg.Incremental = false },
			paths:    []string{".git"},
			expected: false,
		},
		{
			name:     "incremental excluded file",
			setup:    func() { s.cfg.Incremental = true },
			paths:    []string{".git"},
			expected: false,
		},
		{
			name:     "incremental data file",
			setup:    func() { s.cfg.Incremental = true },
			paths:    []string{"_data/foo.yml"},
			expected: true,
		},
		{
			name:     "incremental include file",
			setup:    func() { s.cfg.Incremental = true },
			paths:    []string{"_includes/header.html"},
			expected: true,
		},
		{
			name:     "incremental layout file",
			setup:    func() { s.cfg.Incremental = true },
			paths:    []string{"_layouts/default.html"},
			expected: true,
		},
		{
			name:     "incremental sass file",
			setup:    func() { s.cfg.Incremental = true },
			paths:    []string{"_sass/foo.scss"},
			expected: true,
		},
		{
			name:     "incremental existing document",
			setup:    func() { s.cfg.Incremental = true },
			paths:    []string{"index.md"},
			expected: false,
		},
		{
			name:     "incremental deleted document",
			setup:    func() { s.cfg.Incremental = true },
			paths:    []string{"deleted.md"},
			expected: true,
		},
		{
			name:     "incremental unknown non-document",
			setup:    func() { s.cfg.Incremental = true },
			paths:    []string{"assets/style.css"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.cfg.Incremental = false
			if tt.setup != nil {
				tt.setup()
			}
			got := s.RequiresFullReload(tt.paths)
			require.Equal(t, tt.expected, got)
		})
	}
}

// TestSite_ExcludeStateTable enumerates the include/exclude/underscore/dotfile
// states that drive Site.Exclude.
func TestSite_ExcludeStateTable(t *testing.T) {
	s := New(config.Flags{})
	s.cfg.Source = t.TempDir()
	defaults := New(config.Flags{})

	tests := []struct {
		name     string
		setup    func(*Site)
		path     string
		expected bool
	}{
		{"plain source file", nil, "index.md", false},
		{"top-level underscore directory", nil, "_posts", false},
		{"nested underscore directory", nil, "assets/_hidden/file.md", true},
		{"top-level dotfile", nil, ".git", true},
		{"default excluded file", nil, "Gemfile", true},
		{"default excluded directory prefix", nil, "node_modules/foo.js", true},
		{"default included file", nil, ".htaccess", false},
		{
			name: "include overrides exclude for exact path",
			setup: func(s *Site) {
				s.cfg.Include = []string{"secret.md"}
				s.cfg.Exclude = []string{"secret.md"}
			},
			path:     "secret.md",
			expected: false,
		},
		{
			name: "nested underscore not rescued by parent include",
			setup: func(s *Site) {
				s.cfg.Include = []string{"assets"}
			},
			path:     "assets/_hidden/file.md",
			expected: true,
		},
		{"temporary file prefix", nil, "#draft.md", true},
		{"backup suffix", nil, "file.md~", true},
		{"dotfile in subdirectory", nil, "dir/.hidden", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.cfg.Include = defaults.cfg.Include
			s.cfg.Exclude = defaults.cfg.Exclude
			if tt.setup != nil {
				tt.setup(s)
			}
			got := s.Exclude(tt.path)
			require.Equal(t, tt.expected, got)
		})
	}
}

// TestSite_fileAffectsBuildStateTable enumerates the include/exclude/dotfile/
// destination states that drive whether a changed file can affect the build.
func TestSite_fileAffectsBuildStateTable(t *testing.T) {
	s := New(config.Flags{})
	s.cfg.Source = t.TempDir()
	s.cfg.Destination = "_site"
	defaults := New(config.Flags{})

	tests := []struct {
		name     string
		setup    func(*Site)
		path     string
		expected bool
	}{
		{"plain source file", nil, "index.md", true},
		{"config file", nil, "_config.yml", true},
		{"top-level dotfile", nil, ".git", false},
		{"default excluded file", nil, "Gemfile", false},
		{"default included file", nil, ".htaccess", true},
		{"destination path", nil, "_site/index.html", false},
		{"excluded directory prefix", nil, "node_modules/foo.js", false},
		{
			name: "include overrides exclude for exact path",
			setup: func(s *Site) {
				s.cfg.Include = []string{"secret.md"}
				s.cfg.Exclude = []string{"secret.md"}
			},
			path:     "secret.md",
			expected: true,
		},
		{"dotfile in subdirectory is currently treated as affecting build", nil, "dir/.hidden", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s.cfg.Include = defaults.cfg.Include
			s.cfg.Exclude = defaults.cfg.Exclude
			if tt.setup != nil {
				tt.setup(s)
			}
			got := s.fileAffectsBuild(tt.path)
			require.Equal(t, tt.expected, got)
		})
	}
}

// func TestSite_affectsBuildFilter(t *testing.T) {
// func TestSite_fileAffectsBuild(t *testing.T) {
// func TestSite_invalidatesDoc(t *testing.T) {
