package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLeftPad(t *testing.T) {
	assert.Equal(t, "abc", LeftPad("abc", 0))
	assert.Equal(t, "abc", LeftPad("abc", 3))
	assert.Equal(t, "   abc", LeftPad("abc", 6))
}

func TestGetXXX(t *testing.T) {
	d := map[interface{}]interface{}{
		"t": true,
		"f": false,
		"s": "ss",
	}
	assert.Equal(t, true, getBool(d, "t", true))
	assert.Equal(t, true, getBool(d, "t", false))
	assert.Equal(t, false, getBool(d, "f", true))
	assert.Equal(t, false, getBool(d, "f", true))
	assert.Equal(t, true, getBool(d, "-", true))
	assert.Equal(t, false, getBool(d, "-", false))
	assert.Equal(t, true, getBool(d, "s", true))
	assert.Equal(t, false, getBool(d, "s", false))

	assert.Equal(t, "ss", getString(d, "s", "-"))
	assert.Equal(t, "--", getString(d, "-", "--"))
	assert.Equal(t, "--", getString(d, "t", "--"))
}

func TestMergeMaps(t *testing.T) {
	m1 := map[interface{}]interface{}{"a": 1, "b": 2}
	m2 := map[interface{}]interface{}{"b": 3, "c": 4}
	expected := map[interface{}]interface{}{"a": 1, "b": 3, "c": 4}
	actual := mergeMaps(m1, m2)
	assert.Equal(t, expected, actual)
}

func TestStringMap(t *testing.T) {
	input := map[interface{}]interface{}{"a": 1, 10: 2, false: 3}
	expected := map[string]interface{}{"a": 1, "10": 2, "false": 3}
	actual := stringMap(input)
	assert.Equal(t, expected, actual)
}

func TestStringArrayToMap(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := map[string]bool{"a": true, "b": true, "c": true}
	actual := stringArrayToMap(input)
	assert.Equal(t, expected, actual)
}
