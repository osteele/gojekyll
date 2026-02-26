package utils

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeftPad(t *testing.T) {
	require.Equal(t, "abc", LeftPad("abc", 0))
	require.Equal(t, "abc", LeftPad("abc", 3))
	require.Equal(t, "   abc", LeftPad("abc", 6))
}

func TestSafeReplaceAllStringFunc(t *testing.T) {
	re := regexp.MustCompile(`\w+`)
	out, err := SafeReplaceAllStringFunc(re, "1 > 0", func(m string) (string, error) {
		return fmt.Sprint(m == "1"), nil
	})
	require.NoError(t, err)
	require.Equal(t, "true > false", out)

	_, err = SafeReplaceAllStringFunc(re, "1 > 0", func(m string) (string, error) {
		return "", fmt.Errorf("an expected error")
	})
	require.Error(t, err)
	require.Equal(t, "an expected error", err.Error())
}

func TestSlugify(t *testing.T) {
	require.Equal(t, "abc", Slugify("abc"))
	require.Equal(t, "ab-c", Slugify("ab.c"))
	require.Equal(t, "ab-c", Slugify("ab-c"))
	require.Equal(t, "ab-c", Slugify("ab()[]c"))
	require.Equal(t, "ab123-cde-f-g", Slugify("ab123(cde)[]f.g"))
	require.Equal(t, "abc", Slugify("abc?"))

	// Unicode: Chinese characters preserved
	require.Equal(t, "白法", Slugify("白法"))
	require.Equal(t, "hello-白法-world", Slugify("hello 白法 world"))

	// Unicode: accented characters preserved
	require.Equal(t, "café", Slugify("café"))
	require.Equal(t, "naïve-résumé", Slugify("naïve résumé"))

	// Slugify lowercases
	require.Equal(t, "hello-world", Slugify("Hello World"))
}

func TestSlugifyPermalink(t *testing.T) {
	// Preserves case
	require.Equal(t, "Hello-World", SlugifyPermalink("Hello World"))
	require.Equal(t, "MyPage", SlugifyPermalink("MyPage"))

	// Preserves Unicode
	require.Equal(t, "白法", SlugifyPermalink("白法"))
	require.Equal(t, "Hello-白法-World", SlugifyPermalink("Hello 白法 World"))
	require.Equal(t, "Café", SlugifyPermalink("Café"))

	// Replaces non-alphanumeric with hyphens
	require.Equal(t, "ab-c", SlugifyPermalink("ab.c"))
	require.Equal(t, "ab-c", SlugifyPermalink("ab()[]c"))

	// Strips leading/trailing hyphens
	require.Equal(t, "abc", SlugifyPermalink("abc?"))
	require.Equal(t, "abc", SlugifyPermalink("?abc"))
}

func TestStringArrayToMap(t *testing.T) {
	input := []string{"a", "b", "c"}
	expected := map[string]bool{"a": true, "b": true, "c": true}
	actual := StringArrayToMap(input)
	require.Equal(t, expected, actual)
}
