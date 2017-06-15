package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHyphenateNonAlphaSequence(t *testing.T) {
	assert.Equal(t, "abc", HyphenateNonAlphaSequence("abc"))
	assert.Equal(t, "ab-c", HyphenateNonAlphaSequence("ab.c"))
	assert.Equal(t, "ab-c", HyphenateNonAlphaSequence("ab-c"))
	assert.Equal(t, "ab-c", HyphenateNonAlphaSequence("ab()[]c"))
	assert.Equal(t, "ab123-cde-f-g", HyphenateNonAlphaSequence("ab123(cde)[]f.g"))
}
