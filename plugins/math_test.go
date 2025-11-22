package plugins

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMathPlugin_InlineMath(t *testing.T) {
	plugin := mathPlugin{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple inline math",
			input:    "<p>The equation $$E=mc^2$$ is famous.</p>",
			expected: "<p>The equation \\(E=mc^2\\) is famous.</p>",
		},
		{
			name:     "inline math at start",
			input:    "<p>$$x + y = z$$ is a simple equation.</p>",
			expected: "<p>\\(x + y = z\\) is a simple equation.</p>",
		},
		{
			name:     "inline math at end",
			input:    "<p>Einstein discovered $$E=mc^2$$</p>",
			expected: "<p>Einstein discovered \\(E=mc^2\\)</p>",
		},
		{
			name:     "multiple inline math",
			input:    "<p>Compare $$a$$ and $$b$$ values.</p>",
			expected: "<p>Compare \\(a\\) and \\(b\\) values.</p>",
		},
		{
			name:     "inline math with complex expression",
			input:    "<p>The integral $$\\int_0^\\infty e^{-x^2} dx$$ converges.</p>",
			expected: "<p>The integral \\(\\int_0^\\infty e^{-x^2} dx\\) converges.</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.PostRender([]byte(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(result))
		})
	}
}

func TestMathPlugin_DisplayMath(t *testing.T) {
	plugin := mathPlugin{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standalone display math",
			input:    "<p>$$E=mc^2$$</p>",
			expected: "<p>\\[E=mc^2\\]</p>",
		},
		{
			name:     "display math with whitespace",
			input:    "<p>  $$E=mc^2$$  </p>",
			expected: "<p>  \\[E=mc^2\\]  </p>",
		},
		{
			name:     "display math with newlines",
			input:    "<p>$$\nx^2 + y^2 = z^2\n$$</p>",
			expected: "<p>\\[\nx^2 + y^2 = z^2\n\\]</p>",
		},
		{
			name:     "multiline display math",
			input:    "<p>$$\n\\begin{aligned}\nx &= y + z \\\\\na &= b + c\n\\end{aligned}\n$$</p>",
			expected: "<p>\\[\n\\begin{aligned}\nx &= y + z \\\\\na &= b + c\n\\end{aligned}\n\\]</p>",
		},
		{
			name:     "display math with tabs",
			input:    "<p>\t$$E=mc^2$$\t</p>",
			expected: "<p>\t\\[E=mc^2\\]\t</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.PostRender([]byte(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(result))
		})
	}
}

func TestMathPlugin_MixedContent(t *testing.T) {
	plugin := mathPlugin{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multiple paragraphs with math",
			input:    "<p>First paragraph with $$a=b$$.</p>\n<p>$$x^2$$</p>",
			expected: "<p>First paragraph with \\(a=b\\).</p>\n<p>\\[x^2\\]</p>",
		},
		{
			name:     "math in body without paragraph",
			input:    "<body>$$E=mc^2$$</body>",
			expected: "<body>\\(E=mc^2\\)</body>",
		},
		{
			name:     "no body tag",
			input:    "<div><p>$$E=mc^2$$</p></div>",
			expected: "<div><p>\\[E=mc^2\\]</p></div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.PostRender([]byte(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(result))
		})
	}
}

func TestMathPlugin_EdgeCases(t *testing.T) {
	plugin := mathPlugin{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no math delimiters",
			input:    "<p>This is plain text.</p>",
			expected: "<p>This is plain text.</p>",
		},
		{
			name:     "single dollar sign",
			input:    "<p>This costs $100.</p>",
			expected: "<p>This costs $100.</p>",
		},
		{
			name:     "three dollar signs",
			input:    "<p>This is $$$expensive$$$.</p>",
			expected: "<p>This is $\\(expensive\\)$.</p>",
		},
		{
			name:     "empty math",
			input:    "<p>Empty: $$$$</p>",
			expected: "<p>Empty: $$$$</p>",
		},
		{
			name:     "math in heading",
			input:    "<h1>$$E=mc^2$$</h1>",
			expected: "<h1>\\(E=mc^2\\)</h1>",
		},
		{
			name:     "nested HTML with math",
			input:    "<div><p>Inside: $$x+y$$</p></div>",
			expected: "<div><p>Inside: \\(x+y\\)</p></div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.PostRender([]byte(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(result))
		})
	}
}

func TestMathPlugin_ComplexHTML(t *testing.T) {
	plugin := mathPlugin{}

	input := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
<h1>Math Examples</h1>
<p>Inline math: $$E=mc^2$$ is Einstein's equation.</p>
<p>$$\int_0^\infty e^{-x} dx = 1$$</p>
<p>More text with $$a^2 + b^2 = c^2$$ here.</p>
</body>
</html>`

	expected := `<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body>
<h1>Math Examples</h1>
<p>Inline math: \(E=mc^2\) is Einstein's equation.</p>
<p>\[\int_0^\infty e^{-x} dx = 1\]</p>
<p>More text with \(a^2 + b^2 = c^2\) here.</p>
</body>
</html>`

	result, err := plugin.PostRender([]byte(input))
	require.NoError(t, err)
	require.Equal(t, expected, string(result))
}

func TestMathPlugin_PreserveHTMLTags(t *testing.T) {
	plugin := mathPlugin{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "math with strong tag",
			input:    "<p><strong>Important:</strong> $$E=mc^2$$</p>",
			expected: "<p><strong>Important:</strong> \\(E=mc^2\\)</p>",
		},
		{
			name:     "math with em tag",
			input:    "<p><em>Note</em>: $$x+y=z$$</p>",
			expected: "<p><em>Note</em>: \\(x+y=z\\)</p>",
		},
		{
			name:     "math with link",
			input:    "<p>See <a href=\"#ref\">reference</a> for $$E=mc^2$$</p>",
			expected: "<p>See <a href=\"#ref\">reference</a> for \\(E=mc^2\\)</p>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.PostRender([]byte(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(result))
		})
	}
}
