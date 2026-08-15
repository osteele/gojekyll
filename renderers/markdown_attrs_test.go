package renderers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRenderInnerMarkdownModeMatrix verifies the decision from the markdown
// attribute value to the processing mode.
func TestRenderInnerMarkdownModeMatrix(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantContains []string
		wantAbsent   []string
	}{
		{
			name:         "markdown=1 enables block processing",
			input:        `<div markdown="1">*a*</div>`,
			wantContains: []string{"<em>a</em>"},
			wantAbsent:   []string{"*a*"},
		},
		{
			name:         "markdown=block enables block processing",
			input:        `<div markdown="block">*a*</div>`,
			wantContains: []string{"<em>a</em>"},
			wantAbsent:   []string{"*a*"},
		},
		{
			name:         "markdown=span enables inline processing",
			input:        `<div markdown="span">*a*</div>`,
			wantContains: []string{"<em>a</em>"},
			wantAbsent:   []string{"<p><em>a</em></p>"},
		},
		{
			name:         "markdown=0 disables processing",
			input:        `<div markdown="0">*a*</div>`,
			wantContains: []string{"*a*"},
			wantAbsent:   []string{"<em>"},
		},
		{
			name:         "invalid markdown value disables processing",
			input:        `<div markdown="invalid">*a*</div>`,
			wantContains: []string{"*a*", `markdown="invalid"`},
			wantAbsent:   []string{"<em>"},
		},
		{
			name:         "no markdown attribute leaves content unchanged",
			input:        `<div>*a*</div>`,
			wantContains: []string{"*a*"},
			wantAbsent:   []string{"<em>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := renderInnerMarkdown([]byte(tt.input))
			require.NoError(t, err)
			outStr := string(out)
			for _, s := range tt.wantContains {
				require.Contains(t, outStr, s)
			}
			for _, s := range tt.wantAbsent {
				require.NotContains(t, outStr, s)
			}
		})
	}
}

// TestRenderInnerMarkdownDepthStateTable verifies that void elements and nested
// tags correctly affect the depth count while scanning for the matching end tag.
func TestRenderInnerMarkdownDepthStateTable(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantErr      bool
		wantContains []string
	}{
		{
			name:         "void element br",
			input:        `<div markdown="1"><br></div>`,
			wantContains: []string{"<br"},
		},
		{
			name:         "void element hr",
			input:        `<div markdown="1"><hr></div>`,
			wantContains: []string{"<hr"},
		},
		{
			name:         "self-closing br",
			input:        `<div markdown="1"><br/></div>`,
			wantContains: []string{"<br"},
		},
		{
			name:         "img element",
			input:        `<div markdown="1"><img src="x"></div>`,
			wantContains: []string{"img"},
		},
		{
			name:         "nested non-void tags keep depth",
			input:        `<div markdown="1"><div><span>x</span></div></div>`,
			wantContains: []string{"<div><span>x</span></div>"},
		},
		{
			name:    "unclosed tag returns error",
			input:   `<div markdown="1"><span>*a*</div>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := renderInnerMarkdown([]byte(tt.input))
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			outStr := string(out)
			for _, s := range tt.wantContains {
				require.Contains(t, outStr, s)
			}
		})
	}
}
