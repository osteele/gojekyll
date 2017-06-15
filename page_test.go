package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMarkdown(t *testing.T) {
	assert.True(t, isMarkdown("name.md"))
	assert.True(t, isMarkdown("name.markdown"))
	assert.False(t, isMarkdown("name.html"))
}
