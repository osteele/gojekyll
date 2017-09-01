package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeStringSet(t *testing.T) {
	ss := MakeStringSet([]string{})
	require.Len(t, ss, 0)

	ss = MakeStringSet([]string{"a", "b"})
	require.Len(t, ss, 2)
	require.True(t, ss["a"])
	require.True(t, ss["b"])
}

func TestStringSet_AddStrings(t *testing.T) {
	ss := MakeStringSet([]string{"a", "b"})
	ss.AddStrings([]string{"b", "c", "d"})
	require.True(t, ss["a"])
	require.True(t, ss["b"])
	require.True(t, ss["c"])
	require.True(t, ss["d"])
}
