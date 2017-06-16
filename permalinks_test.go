package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPermalinkPattern(t *testing.T) {
	var (
		d    = VariableMap{}
		path = "/a/b/base.html"
	)

	testPermalinkPattern := func(pattern, path string, data VariableMap) (string, error) {
		vs := mergeVariableMaps(data, VariableMap{"permalink": pattern})
		p := pageFields{nil, path, "", vs}
		return p.expandPermalink()
	}

	t.Run(":output_ext", func(t *testing.T) {
		p, _ := testPermalinkPattern("/base:output_ext", path, d)
		assert.Equal(t, "/base.html", p)
	})
	t.Run(":output_ext renames markdown to .html", func(t *testing.T) {
		p, _ := testPermalinkPattern("/base:output_ext", "/a/b/base.md", d)
		assert.Equal(t, "/base.html", p)
		p, _ = testPermalinkPattern("/base:output_ext", "/a/b/base.markdown", d)
		assert.Equal(t, "/base.html", p)
	})
	t.Run(":name", func(t *testing.T) {
		p, _ := testPermalinkPattern("/name/:name", path, d)
		assert.Equal(t, "/name/base", p)
	})
	t.Run(":path", func(t *testing.T) {
		p, _ := testPermalinkPattern("/prefix:path/post", path, d)
		assert.Equal(t, "/prefix/a/b/base/post", p)
	})
	t.Run(":title", func(t *testing.T) {
		p, _ := testPermalinkPattern("/title/:title.html", path, d)
		assert.Equal(t, "/title/base.html", p)
	})
	t.Run("invalid template variable", func(t *testing.T) {
		_, err := testPermalinkPattern("/:invalid", path, d)
		assert.Error(t, err)
	})

	d["collection"] = "c"
	path = "_c/a/b/c.d"
	t.Run(":path", func(t *testing.T) {
		p, _ := testPermalinkPattern("/prefix:path/post", path, d)
		assert.Equal(t, "/prefix/a/b/c/post", p)
	})
}
