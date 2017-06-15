package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHyphenateNonAlphaSequence(t *testing.T) {
	assert.Equal(t, "abc", hyphenateNonAlphaSequence("abc"))
	assert.Equal(t, "ab-c", hyphenateNonAlphaSequence("ab.c"))
	assert.Equal(t, "ab-c", hyphenateNonAlphaSequence("ab-c"))
	assert.Equal(t, "ab-c", hyphenateNonAlphaSequence("ab()[]c"))
	assert.Equal(t, "ab123-cde-f-g", hyphenateNonAlphaSequence("ab123(cde)[]f.g"))
}
