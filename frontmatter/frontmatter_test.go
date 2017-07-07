package frontmatter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFileHasFrontMatter(t *testing.T) {
	fm := func(filename string) bool {
		fm, err := FileHasFrontMatter(filename)
		require.NoError(t, err)
		return fm
	}
	require.True(t, fm("testdata/empty_fm.md"))
	require.True(t, fm("testdata/some_fm.md"))
	require.False(t, fm("testdata/no_fm.md"))
}

func TestFrontMatter_SortedStringArray(t *testing.T) {
	sortedStringValue := func(value interface{}) []string {
		fm := map[string]interface{}{"categories": value}
		return FrontMatter(fm).SortedStringArray("categories")
	}
	require.Equal(t, []string{"a", "b"}, sortedStringValue("b a"))
	require.Equal(t, []string{"a", "b"}, sortedStringValue([]interface{}{"b", "a"}))
	require.Equal(t, []string{"a", "b"}, sortedStringValue([]string{"b", "a"}))
	require.Equal(t, []string{}, sortedStringValue(3))
}
