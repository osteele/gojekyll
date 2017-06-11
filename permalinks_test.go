package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPermalinkPattern(t *testing.T) {
	d := map[interface{}]interface{}{}
	path := "/a/b/c.d"
	t.Run(":name", func(t *testing.T) {
		p := expandPermalinkPattern("/name/:name", d, path)
		assert.Equal(t, "/name/c-d", p)
	})
	t.Run(":path", func(t *testing.T) {
		p := expandPermalinkPattern("/prefix:path/post", d, path)
		assert.Equal(t, "/prefix/a/b/c.d/post", p)
	})
	t.Run(":title", func(t *testing.T) {
		pl := expandPermalinkPattern("/title/:title.html", d, path)
		assert.Equal(t, "/title/c.html", pl)
	})

	d["collection"] = "c"
	path = "/_c/a/b/c.d"
	t.Run(":path", func(t *testing.T) {
		p := expandPermalinkPattern("/prefix:path/post", d, path)
		assert.Equal(t, "/prefix/a/b/c.d/post", p)
	})
}
