package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlugify(t *testing.T) {
	require.Equal(t, "abc", Slugify("abc"))
	require.Equal(t, "ab-c", Slugify("ab.c"))
	require.Equal(t, "ab-c", Slugify("ab-c"))
	require.Equal(t, "ab-c", Slugify("ab()[]c"))
	require.Equal(t, "ab123-cde-f-g", Slugify("ab123(cde)[]f.g"))
}
func TestLeftPad(t *testing.T) {
	require.Equal(t, "abc", LeftPad("abc", 0))
	require.Equal(t, "abc", LeftPad("abc", 3))
	require.Equal(t, "   abc", LeftPad("abc", 6))
}

func TestGetXXX(t *testing.T) {
	d := VariableMap{
		"t": true,
		"f": false,
		"s": "ss",
	}
	require.Equal(t, true, d.Bool("t", true))
	require.Equal(t, true, d.Bool("t", false))
	require.Equal(t, false, d.Bool("f", true))
	require.Equal(t, false, d.Bool("f", true))
	require.Equal(t, true, d.Bool("-", true))
	require.Equal(t, false, d.Bool("-", false))
	require.Equal(t, true, d.Bool("s", true))
	require.Equal(t, false, d.Bool("s", false))

	require.Equal(t, "ss", d.String("s", "-"))
	require.Equal(t, "--", d.String("-", "--"))
	require.Equal(t, "--", d.String("t", "--"))
}

func TestMergeVariableMaps(t *testing.T) {
	m1 := VariableMap{"a": 1, "b": 2}
	m2 := VariableMap{"b": 3, "c": 4}
	expected := VariableMap{"a": 1, "b": 3, "c": 4}
	actual := MergeVariableMaps(m1, m2)
	require.Equal(t, expected, actual)
}

func TestStringArrayToMap(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := map[string]bool{"a": true, "b": true, "c": true}
	actual := stringArrayToMap(input)
	require.Equal(t, expected, actual)
}
