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
	d := VariableMap{
		"t": true,
		"f": false,
		"s": "ss",
	}
	assert.Equal(t, true, d.Bool("t", true))
	assert.Equal(t, true, d.Bool("t", false))
	assert.Equal(t, false, d.Bool("f", true))
	assert.Equal(t, false, d.Bool("f", true))
	assert.Equal(t, true, d.Bool("-", true))
	assert.Equal(t, false, d.Bool("-", false))
	assert.Equal(t, true, d.Bool("s", true))
	assert.Equal(t, false, d.Bool("s", false))

	assert.Equal(t, "ss", d.String("s", "-"))
	assert.Equal(t, "--", d.String("-", "--"))
	assert.Equal(t, "--", d.String("t", "--"))
}

func TestMakeVariableMap(t *testing.T) {
	input := map[interface{}]interface{}{"a": 1, 10: 2, false: 3}
	expected := VariableMap{"a": 1, "10": 2, "false": 3}
	actual := makeVariableMap(input)
	assert.Equal(t, expected, actual)
}

func TestMergeVariableMaps(t *testing.T) {
	m1 := VariableMap{"a": 1, "b": 2}
	m2 := VariableMap{"b": 3, "c": 4}
	expected := VariableMap{"a": 1, "b": 3, "c": 4}
	actual := mergeVariableMaps(m1, m2)
	assert.Equal(t, expected, actual)
}

func TestStringArrayToMap(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := map[string]bool{"a": true, "b": true, "c": true}
	actual := stringArrayToMap(input)
	assert.Equal(t, expected, actual)
}
