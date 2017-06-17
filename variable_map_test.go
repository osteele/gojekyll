package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
