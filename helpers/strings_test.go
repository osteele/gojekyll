package helpers

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

func TestStringArrayToMap(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := map[string]bool{"a": true, "b": true, "c": true}
	actual := StringArrayToMap(input)
	require.Equal(t, expected, actual)
}
