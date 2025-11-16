# Jekyll TOC Behavior Investigation

## Summary

Investigation of actual Ruby Jekyll 4.4.1 behavior vs gojekyll implementation for TOC markers.

## Key Finding

**`{::toc}` is NOT a valid kramdown syntax and is NOT processed by Jekyll at all.**

## Official Kramdown Documentation

### Valid TOC Syntax

From https://kramdown.gettalong.org/converter/html.html:

> Just assign the reference name 'toc' to an ordered or unordered list by using an IAL and the list will be replaced with the actual table of contents.

**Syntax:**
```markdown
* placeholder text
{:toc}
```

The `{:toc}` must be placed immediately after a list item (either `*` or `1.`).

### Block Syntax `{::}`

From https://kramdown.gettalong.org/syntax.html:

The `{::}` syntax is for **kramdown extensions only**. Valid extensions are:
- `{::comment}` - creates ignored comments
- `{::nomarkdown}` - outputs body as-is without processing
- `{::options}` - sets global processor options

**`{::toc}` is NOT documented** as a valid kramdown extension.

## Actual Jekyll Behavior (Tested with Jekyll 4.4.1)

Test file: `toc-variations.md`

### Test Results

| Marker | Context | Jekyll Output | Expected Behavior |
|--------|---------|---------------|-------------------|
| `{:toc}` | In `<ul>` | TOC generated, replaces `<ul>` | ✅ Correct |
| `{:toc}` | In `<ol>` | **Marker remains as text** | ✅ Correct (not supported) |
| `{::toc}` | Standalone | **Marker remains as `{::toc}`** | ✅ Correct (invalid syntax) |
| `{::toc}` | In `<ul>` | **Marker remains as `{::toc}`** | ✅ Correct (invalid syntax) |

### Jekyll HTML Output Excerpt

```html
<!-- Test 2: {:toc} in ordered list -->
<ol>
  <li>This should remain</li>
</ol>

<!-- Test 3: {::toc} standalone -->
<p>{::toc}</p>

<!-- Test 4: {::toc} in list -->
<ul>
  <li>This should remain
{::toc}</li>
</ul>
```

## gojekyll Current Behavior (INCORRECT)

Test file: `toc-syntax-variations-test.html`

| Marker | Context | gojekyll Output | Status |
|--------|---------|-----------------|--------|
| `{:toc}` | In `<ul>` | TOC generated | ✅ Correct |
| `{:toc}` | In `<ol>` | Marker remains | ✅ Correct |
| `{::toc}` | Standalone | **TOC generated** | ❌ Should remain as text |
| `{::toc}` | In `<ul>` | Marker remains | ✅ Correct |

### Additional Issues Found

1. **HTML Corruption**: TOC divs appearing inside `<h2>` tags (lines 7, 13, 22, 30)
2. **{:.no_toc} not fully working**: Heading "This heading should be excluded" appears in TOC (line 45-46)
3. **{::toc} processed as standalone**: Should NOT be processed at all

## Conclusion

### What Needs to be Fixed in gojekyll

1. **Remove `{::toc}` processing entirely**
   - It's not valid kramdown syntax
   - Jekyll doesn't process it
   - Should remain as literal text in output

2. **Fix HTML corruption**
   - TOC should not appear inside heading tags
   - This is likely a bug in the DOM replacement logic

3. **Fix {:.no_toc} handling**
   - Currently including headings that should be excluded
   - The marker is being removed but heading still appears in TOC

### What the Tests Got Wrong

The test file `/Users/osteele/code/gojekyll-toc/commands/testdata/site/toc-variations.md` incorrectly states:

> **{::toc} standalone**: Inserts TOC at that location

This is **incorrect**. Jekyll does not process `{::toc}` at all.

### Correct Behavior Summary

**Only `{:toc}` in an unordered list is valid TOC syntax in Jekyll/kramdown:**

```markdown
# My Page

* Table of Contents
{:toc}

## Section 1
## Section 2
```

All other uses of TOC markers should:
- `{:toc}` in `<ol>` → remain as literal text
- `{::toc}` anywhere → remain as literal text (invalid syntax)
- `{:.no_toc}` after heading → exclude that heading from TOC

## References

1. Kramdown HTML Converter: https://kramdown.gettalong.org/converter/html.html
2. Kramdown Syntax: https://kramdown.gettalong.org/syntax.html
3. Jekyll 4.4.1 tested with: `commands/testdata/site/toc-variations.md`
