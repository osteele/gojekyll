package renderers

import (
	"bytes"

	"github.com/gohugoio/hugo-goldmark-extensions/passthrough"
	"github.com/osteele/gojekyll/utils"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// createGoldmarkConverter creates a Goldmark markdown converter configured
// to match Jekyll/kramdown behavior as closely as possible
func createGoldmarkConverter() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,            // GitHub Flavored Markdown (includes tables, strikethrough, autolinks)
			extension.Footnote,       // Footnotes support
			extension.DefinitionList, // Definition lists
			extension.Typographer,    // Smart quotes and dashes (like Smartypants)
			passthrough.New(passthrough.Config{ // Math delimiters passthrough
				// Inline math: $$...$$ → preserved as-is for client-side rendering
				InlineDelimiters: []passthrough.Delimiters{
					{Open: "$$", Close: "$$"},
				},
				// Block/display math: $$...$$ on separate lines
				BlockDelimiters: []passthrough.Delimiters{
					{Open: "$$", Close: "$$"},
				},
			}),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Automatic heading IDs
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),  // Use XHTML tags (like Blackfriday's UseXHTML)
			html.WithUnsafe(), // Allow raw HTML (Jekyll/kramdown compatibility)
		),
	)
}

func renderMarkdown(md []byte) ([]byte, error) {
	return renderMarkdownWithOptions(md, nil)
}

func renderMarkdownWithOptions(md []byte, opts *TOCOptions) ([]byte, error) {
	// Set default options if not provided
	// Jekyll's default toc_levels is "2..6" to exclude H1 headings
	if opts == nil {
		opts = &TOCOptions{
			MinLevel:      2,
			MaxLevel:      6,
			UseJekyllHTML: true, // Use Jekyll-compatible HTML structure by default
		}
	}
	// Ensure valid level ranges
	if opts.MinLevel < 1 {
		opts.MinLevel = 1
	}
	if opts.MaxLevel > 6 {
		opts.MaxLevel = 6
	}
	if opts.MinLevel > opts.MaxLevel {
		opts.MinLevel = 1
		opts.MaxLevel = 6
	}

	// Create Goldmark converter and render markdown to HTML
	converter := createGoldmarkConverter()
	var buf bytes.Buffer
	if err := converter.Convert(md, &buf); err != nil {
		return nil, utils.WrapError(err, "markdown conversion")
	}
	html := buf.Bytes()

	// Process inner markdown (for nested markdown rendering)
	html, err := renderInnerMarkdown(html)
	if err != nil {
		return nil, utils.WrapError(err, "markdown")
	}

	// Process TOC markers if they exist
	// Note: Only {:toc} is valid kramdown syntax; {::toc} is not processed
	// Jekyll only processes {:toc} in unordered lists, leaving literals elsewhere
	if tocPatternInline.Match(html) && shouldProcessTOC(html) {
		html, err = processTOC(html, opts)
		if err != nil {
			return nil, utils.WrapError(err, "toc generation")
		}
	}
	return html, nil
}

func _renderMarkdown(md []byte) ([]byte, error) {
	converter := createGoldmarkConverter()
	var buf bytes.Buffer
	if err := converter.Convert(md, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
