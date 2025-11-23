package renderers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMathDelimiters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string // what the output should contain
	}{
		{
			name:     "Inline math with $$",
			input:    "The equation $$E=mc^2$$ is famous.",
			contains: "$$E=mc^2$$",
		},
		{
			name: "Display math block",
			input: `Some text

$$
\int_0^\infty e^{-x} dx = 1
$$

More text`,
			contains: "$$",
		},
		{
			name:     "Math with underscores",
			input:    "The variable $$x_0$$ represents the initial value.",
			contains: "$$x_0$$",
		},
		{
			name:     "Math with asterisks",
			input:    "The expression $$x * y$$ shows multiplication.",
			contains: "$$x * y$$",
		},
		{
			name: "Complex math expression",
			input: `$$
\begin{bmatrix}
a & b \\
c & d
\end{bmatrix}
$$`,
			contains: "\\begin{bmatrix}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html, err := renderMarkdown([]byte(tt.input))
			require.NoError(t, err)
			htmlStr := string(html)

			// Verify that the math delimiters and content are preserved
			require.Contains(t, htmlStr, tt.contains,
				"Expected output to contain math expression")

			// Verify math is NOT converted to HTML entities or modified
			require.NotContains(t, htmlStr, "&lt;", "Math should not be HTML-escaped")
			require.NotContains(t, htmlStr, "<em>", "Underscores in math should not create emphasis")
		})
	}
}

func TestMathWithMarkdown(t *testing.T) {
	input := `# Header

Some **bold** text and $$E=mc^2$$ inline math.

$$
F = ma
$$

More _italic_ text.`

	html, err := renderMarkdown([]byte(input))
	require.NoError(t, err)
	htmlStr := string(html)

	// Verify markdown is processed
	require.Contains(t, htmlStr, "<strong>bold</strong>")
	require.Contains(t, htmlStr, "<em>italic</em>")
	require.Contains(t, htmlStr, "<h1")

	// Verify math is preserved
	require.Contains(t, htmlStr, "$$E=mc^2$$")
	require.Contains(t, htmlStr, "$$")
}

func TestMathDelimitersNotInCodeBlocks(t *testing.T) {
	input := "```\n$$E=mc^2$$\n```"

	html, err := renderMarkdown([]byte(input))
	require.NoError(t, err)
	htmlStr := string(html)

	// Math in code blocks should be treated as code, not math
	require.Contains(t, htmlStr, "<code>")
	// The $$ should be inside the code block, not treated as math passthrough
	require.True(t, strings.Contains(htmlStr, "<code") && strings.Contains(htmlStr, "$$"),
		"Code block should contain $$ as literal text")
}
