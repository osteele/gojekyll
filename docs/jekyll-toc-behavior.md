# Jekyll TOC Behavior - Empirical Findings

This document records the actual behavior of Jekyll's TOC generation, based on testing with real Jekyll (not documentation or assumptions).

**Test Environment:**
- Jekyll version: 4.4.1
- Markdown processor: Kramdown (Jekyll's default)

## Key Finding

**`{::toc}` is NOT a valid kramdown syntax and is NOT processed by Jekyll at all.**

## TOC Marker Processing

### Where `{:toc}` is Processed (Replaced with TOC)

✅ **Unordered lists only**
```markdown
* TOC
{:toc}
```
Result: Entire `<ul>` is replaced with `<ul id="markdown-toc">...</ul>` containing the table of contents.

### Where `{:toc}` is NOT Processed (Remains Literal or Removed)

❌ **Ordered lists**
```markdown
1. TOC
{:toc}
```
Result: Renders as `<ol><li>TOC</li></ol>` - the `{:toc}` marker is completely removed by Kramdown's IAL processing, but NO TOC is generated.

❌ **Standalone paragraphs**
```markdown
{:toc}
```
Result: Completely removed from output (Kramdown IAL processing removes IAL-only paragraphs)

✅ **Inside heading text**
```markdown
## How to use {:toc}
```
Result: Renders as `<h2>How to use {:toc}</h2>` - kept literally in heading text

✅ **Inside paragraph text**
```markdown
Some text about {:toc} markers.
```
Result: Renders as `<p>Some text about {:toc} markers.</p>` - kept literally

### `{::toc}` Block Syntax (Invalid)

❌ **All contexts**
```markdown
{::toc}
```
Result: Remains as literal text `{::toc}` - this is NOT valid kramdown syntax.

From https://kramdown.gettalong.org/syntax.html, the `{::}` syntax is for **kramdown extensions only**. Valid extensions are:
- `{::comment}` - creates ignored comments
- `{::nomarkdown}` - outputs body as-is without processing
- `{::options}` - sets global processor options

**`{::toc}` is NOT documented** as a valid kramdown extension.

## `.no_toc` Marker Behavior

### Working `.no_toc` (Excludes from TOC)

✅ **As IAL on line after heading**
```markdown
## Excluded Heading
{:.no_toc}
```
Result: Heading appears in document as `<h2 class="no_toc">Excluded Heading</h2>` but is **excluded from TOC**.

### Non-Working `.no_toc` (Does NOT Exclude)

❌ **Inline with heading text**
```markdown
## Excluded Heading {:.no_toc}
```
Result: Renders as `<h2>Excluded Heading {:.no_toc}</h2>` with literal `{:.no_toc}` text, and **appears in TOC** including the literal text!

## TOC HTML Structure

Jekyll generates this structure:
```html
<ul id="markdown-toc">
  <li><a href="#section-1" id="markdown-toc-section-1">Section 1</a>
    <ul>
      <li><a href="#subsection" id="markdown-toc-subsection">Subsection</a></li>
    </ul>
  </li>
  <li><a href="#section-2" id="markdown-toc-section-2">Section 2</a></li>
</ul>
```

Key properties:
- Root `<ul>` has `id="markdown-toc"`
- Each TOC entry has:
  - `<a href="#heading-id" id="markdown-toc-heading-id">Heading Text</a>`
  - Nested `<ul>` for child headings (no `id` on nested lists)
- Always uses `<ul>`, never `<ol>`, even if the marker was in an ordered list

## Key Differences from Kramdown Alone

1. **TOC Generation is Jekyll Plugin Behavior**
   - Kramdown just processes IAL markers `{:...}` as attributes
   - Jekyll's TOC plugin specifically looks for `{:toc}` in unordered lists and replaces them
   - This is why `{:toc}` in ordered lists gets removed (Kramdown IAL) but doesn't generate a TOC (Jekyll plugin)

2. **Standalone IAL Removal**
   - Kramdown removes IAL-only paragraphs (paragraphs containing only `{:something}`)
   - This is why standalone `{:toc}` disappears completely

3. **IAL Application Rules**
   - IAL on its own line applies to the **following** block element
   - IAL at end of line (with space) would apply to that element (but Kramdown doesn't parse `## Heading {:.class}` as IAL - it treats it as literal text)

## Emulation Challenges for Blackfriday

Blackfriday (our markdown processor) differs from Kramdown:

1. **No IAL support** - Blackfriday doesn't understand `{:...}` syntax at all
2. **IALs render as literal text** - `{:toc}` appears in output as literal text
3. **Must process HTML output** - We must:
   - Parse Blackfriday's HTML output
   - Find `{:toc}` text nodes
   - Determine context (in `<ul>`, `<ol>`, `<p>`, `<h2>`, etc.)
   - Only replace if in unordered list
   - Leave literal everywhere else

## Implementation Details

### How We Handle `{:.no_toc}`

1. **Sibling paragraph detection**: Check if heading has a next sibling `<p>` containing only `{:.no_toc}`
2. **Removal**: If found, remove that paragraph from DOM (matching Kramdown's IAL processing)
3. **Exclusion**: Mark heading for exclusion from TOC
4. **Literal preservation**: `{:.no_toc}` inside heading text is NOT removed and does NOT exclude

This matches Jekyll/Kramdown behavior where IAL markers are only processed when on their own line after an element.

## Test Coverage

Tests added to verify Jekyll-compatible behavior:
- ✅ `{:toc}` in headings remains literal
- ✅ `{:toc}` in paragraph text remains literal
- ✅ `{:toc}` in ordered lists remains literal (removed by our IAL-like processing)
- ✅ `{:toc}` in unordered lists generates TOC
- ✅ `{:.no_toc}` as sibling paragraph excludes heading from TOC
- ✅ `{:.no_toc}` inline in heading text remains literal and does NOT exclude

Note: We cannot test standalone `{:toc}` removal because Blackfriday renders it as `<p>{:toc}</p>`, not as an empty IAL like Kramdown does.

## Unresolved Questions

- [ ] What happens with `{:toc}` in other contexts (tables, blockquotes, definition lists)?
- [ ] Does TOC respect `toc_levels` configuration?
- [ ] What happens with multiple `{:toc}` markers in the same document?
- [ ] Edge case: What if `{:toc}` is in an unordered list inside a blockquote?

## Official Kramdown Documentation References

1. Kramdown HTML Converter: https://kramdown.gettalong.org/converter/html.html
2. Kramdown Syntax: https://kramdown.gettalong.org/syntax.html
