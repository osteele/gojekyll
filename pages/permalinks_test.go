package pages

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/osteele/gojekyll/templates"
	"github.com/stretchr/testify/require"
)

type containerMock struct{ pathPrefix string }

func (c containerMock) Output() bool             { return true }
func (c containerMock) PathPrefix() string       { return c.pathPrefix }
func (c containerMock) DefaultPermalink() string { return "/:path:output_ext" }

type pathTest struct{ path, pattern, out string }

var tests = []pathTest{
	{"/a/b/base.html", "/out:output_ext", "/out.html"},
	{"/a/b/base.md", "/out:output_ext", "/out.html"},
	{"/a/b/base.markdown", "/out:output_ext", "/out.html"},
	{"/a/b/base.html", "/:path/out:output_ext", "/a/b/base/out.html"},
	{"/a/b/base.html", "/prefix/:name", "/prefix/base"},
	{"/a/b/base.html", "/prefix/:path/post", "/prefix/a/b/base/post"},
	{"/a/b/base.html", "/prefix/:title", "/prefix/base"},
	{"/a/b/base.html", "/prefix/:slug", "/prefix/base"},
	{"base", "/:categories/:name:output_ext", "/a/b/base"},

	{"base", "date", "/a/b/2006/02/03/base.html"},
	{"base", "pretty", "/a/b/2006/02/03/base/"},
	{"base", "ordinal", "/a/b/2006/34/base.html"},
	{"base", "none", "/a/b/base.html"},
}

var collectionTests = []pathTest{
	{"_c/a/b/c.d", "/prefix/:collection/post", "/prefix/c/post"},
	{"_c/a/b/c.d", "/prefix:path/post", "/prefix/a/b/c/post"},
}

func TestExpandPermalinkPattern(t *testing.T) {
	var (
		c = containerMock{}
		d = templates.VariableMap{
			"categories": "b a",
		}
	)

	testPermalinkPattern := func(pattern, path string, data templates.VariableMap) (string, error) {
		vs := templates.MergeVariableMaps(data, templates.VariableMap{"permalink": pattern})
		ext := filepath.Ext(path)
		switch ext {
		case ".md", ".markdown":
			ext = ".html"
		}
		p := pageFields{container: c, relpath: path, frontMatter: vs, outputExt: ext}
		t0, err := time.Parse(time.RFC3339, "2006-02-03T15:04:05Z")
		require.NoError(t, err)
		p.modTime = t0
		return p.expandPermalink()
	}

	runTests := func(tests []pathTest) {
		for i, test := range tests {
			t.Run(test.pattern, func(t *testing.T) {
				p, err := testPermalinkPattern(test.pattern, test.path, d)
				require.NoError(t, err)
				require.Equalf(t, test.out, p, "%d: pattern=%s", i+1, test.pattern)
			})
		}
	}

	runTests(tests)

	c = containerMock{"_c/"}
	d["collection"] = "c"
	runTests(collectionTests)

	t.Run("invalid template variable", func(t *testing.T) {
		p, err := testPermalinkPattern("/:invalid", "/a/b/base.html", d)
		require.Error(t, err)
		require.Zero(t, p)
	})
}
