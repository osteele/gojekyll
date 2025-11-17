package renderers

import (
	blackfriday "github.com/danog/blackfriday/v2"
	"github.com/osteele/gojekyll/utils"
)

const blackfridayFlags = 0 |
	blackfriday.UseXHTML |
	blackfriday.Smartypants |
	blackfriday.SmartypantsFractions |
	blackfriday.SmartypantsDashes |
	blackfriday.SmartypantsLatexDashes |
	blackfriday.FootnoteReturnLinks

const blackfridayExtensions = 0 |
	blackfriday.NoIntraEmphasis |
	blackfriday.Tables |
	blackfriday.FencedCode |
	blackfriday.Autolink |
	blackfriday.Strikethrough |
	blackfriday.SpaceHeadings |
	blackfriday.HeadingIDs |
	blackfriday.BackslashLineBreak |
	blackfriday.DefinitionLists |
	blackfriday.NoEmptyLineBeforeBlock |
	// added relative to commonExtensions
	blackfriday.AutoHeadingIDs |
	blackfriday.Footnotes

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

	params := blackfriday.HTMLRendererParameters{
		Flags: blackfridayFlags,
	}
	renderer := blackfriday.NewHTMLRenderer(params)
	html := blackfriday.Run(
		md,
		blackfriday.WithRenderer(renderer),
		blackfriday.WithExtensions(blackfridayExtensions),
	)
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
	params := blackfriday.HTMLRendererParameters{
		Flags: blackfridayFlags,
	}
	renderer := blackfriday.NewHTMLRenderer(params)
	html := blackfriday.Run(
		md,
		blackfriday.WithRenderer(renderer),
		blackfriday.WithExtensions(blackfridayExtensions),
	)
	return html, nil
}
