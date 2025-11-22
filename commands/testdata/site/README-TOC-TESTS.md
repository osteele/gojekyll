# TOC Test Pages

This directory contains test pages for verifying Table of Contents (TOC) behavior in gojekyll.

## Test Pages

### toc-test.md
Comprehensive test of TOC features:
- Basic {:toc} inline syntax
- {::toc} block syntax
- {:.no_toc} heading exclusion
- Multiple TOC markers in one document
- Different list placeholder text
- All heading levels (H1-H6)

### toc-h1-test.md
Specifically tests H1 inclusion behavior:
- Tests whether H1 headings appear in TOC
- Documents expected behavior with default toc_levels (2..6)
- Can be used to test toc_levels configuration

### toc-variations.md
Tests different TOC marker syntax variations:
- {:toc} in unordered lists (should work)
- {:toc} in ordered lists (should NOT work in Jekyll)
- {::toc} standalone (should work)
- {::toc} in lists (should NOT replace list item)
- Whitespace variations
- {:.no_toc} marker behavior

## How to Build

### With gojekyll:
```bash
# From repository root
go build -o /tmp/gojekyll .
cd commands/testdata/site
/tmp/gojekyll build --destination /tmp/gojekyll-toc-test
```

### With Ruby Jekyll (for comparison):
```bash
cd commands/testdata/site
jekyll build --destination /tmp/jekyll-toc-test
```

## Comparing Output

```bash
# View gojekyll output
open /tmp/gojekyll-toc-test/toc-test-page.html

# View Jekyll output (if available)
open /tmp/jekyll-toc-test/toc-test-page.html

# Diff the outputs
diff /tmp/gojekyll-toc-test/toc-test-page.html /tmp/jekyll-toc-test/toc-test-page.html
```

## Expected Kramdown/Jekyll Behavior

### TOC Marker Syntax

1. **{:toc}** - Inline attribute list
   - When in unordered list: replaces the list item
   - When in ordered list: NOT supported (marker ignored)
   - When standalone: generates TOC

2. **{::toc}** - Block syntax
   - Always generates TOC at that location
   - Never replaces surrounding content

3. **{:.no_toc}** - Heading exclusion
   - Place after heading text
   - Excludes heading from all TOCs

### Default Configuration

Jekyll's default `kramdown.toc_levels` is `"2..6"`, which means:
- H1 headings are excluded from TOC
- Only H2 through H6 appear

To include H1, add to `_config.yml`:
```yaml
kramdown:
  toc_levels: 1..6
```

## References

- [Kramdown TOC Documentation](https://kramdown.gettalong.org/converter/html.html#toc)
- [Jekyll Kramdown Configuration](https://jekyllrb.com/docs/configuration/markdown/)
- [GitHub Issue #62](https://github.com/osteele/gojekyll/issues/62)
