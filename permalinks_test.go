package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPermalinkPattern(t *testing.T) {
	var (
		d      = map[interface{}]interface{}{}
		path   = "/a/b/c.d"
		mdPath = "/a/b/c.md"
	)

	t.Run(":ext", func(t *testing.T) {
		p, _ := expandPermalinkPattern("/ext/:ext", path, d)
		assert.Equal(t, "/ext/d", p)
	})
	t.Run(":ext", func(t *testing.T) {
		p, _ := expandPermalinkPattern("/ext/:ext", mdPath, d)
		assert.Equal(t, "/ext/md", p)
	})
	t.Run(":output_ext", func(t *testing.T) {
		p, _ := expandPermalinkPattern("/ext/:output_ext", path, d)
		assert.Equal(t, "/ext/d", p)
	})
	t.Run(":output_ext", func(t *testing.T) {
		p, _ := expandPermalinkPattern("/ext/:output_ext", mdPath, d)
		assert.Equal(t, "/ext/html", p)
	})
	t.Run(":name", func(t *testing.T) {
		p, _ := expandPermalinkPattern("/name/:name", path, d)
		assert.Equal(t, "/name/c", p)
	})
	t.Run(":path", func(t *testing.T) {
		p, _ := expandPermalinkPattern("/prefix:path/post", path, d)
		assert.Equal(t, "/prefix/a/b/c.d/post", p)
	})
	t.Run(":title", func(t *testing.T) {
		p, _ := expandPermalinkPattern("/title/:title.html", path, d)
		assert.Equal(t, "/title/c.html", p)
	})

	d["collection"] = "c"
	path = "/_c/a/b/c.d"
	t.Run(":path", func(t *testing.T) {
		p, _ := expandPermalinkPattern("/prefix:path/post", path, d)
		assert.Equal(t, "/prefix/a/b/c.d/post", p)
	})
}
