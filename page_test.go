package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMarkdown(t *testing.T) {
	assert.True(t, site.IsMarkdown("name.md"))
	assert.True(t, site.IsMarkdown("name.markdown"))
	assert.False(t, site.IsMarkdown("name.html"))
}
