# Investigation Report: Issue #51 - Build/Serve Termination on Parse Errors

## Issue Summary

**Title:** `build` and `serve` terminate when encountering the first file that cannot be parsed

**Reporter:** manastungare

**Issue Link:** https://github.com/osteele/gojekyll/issues/51

## Problem Description

The gojekyll tool halts completely when encountering a single unparseable file, preventing visibility into all errors within a project. The specific complaints are:

1. **Early termination:** Processing stops at the first problematic file rather than continuing through all files
2. **Cryptic error messages:** Users receive minimal diagnostic information such as `"filename.md: markdown: EOF"` without context
3. **Documentation gap:** Unsupported markdown features like `markdown="1"` are not clearly documented

### Reproduction Case

The issue manifests with HTML containing:

```html
<div markdown="1">

<br>
<br>

</div>
```

The problem is that `<br>` should be self-closing as `<br/>` in XHTML, but the HTML parser encounters an EOF when looking for the closing tag, causing the build to fail.

## Root Cause Analysis

### Error Flow Diagram

```
1. File with markdown="1" and <br> tags
   ↓
2. During rendering phase (site/render.go)
   ↓
3. Page.Render() called → computeContent() → rendererManager.Render()
   ↓
4. Markdown renderer (renderers/markdown.go)
   ↓
5. renderInnerMarkdown() → processInnerMarkdown()
   ↓
6. HTML tokenizer encounters ErrorToken (line 137-138)
   ↓
7. Returns z.Err() (EOF) immediately
   ↓
8. Error propagates up to site.render() (site/render.go:12-14, 17-19)
   ↓
9. Returns error immediately, terminating the entire build
```

### Key Code Locations

#### 1. Error Generation: `renderers/markdown.go:130-160`

```go
func processInnerMarkdown(w io.Writer, z *html.Tokenizer) error {
    // ...
    for {
        tt := z.Next()
        switch tt {
        case html.ErrorToken:
            return z.Err()  // LINE 138: Returns EOF immediately!
        // ...
        }
    }
}
```

**Issue:** When the HTML tokenizer encounters malformed HTML (like `<br>` instead of `<br/>`), it returns an EOF error immediately without providing context about what went wrong.

#### 2. Early Termination in Rendering: `site/render.go:10-22`

```go
func (s *Site) render() error {
    for _, c := range s.sortedCollections() {
        if err := c.Render(); err != nil {
            return err  // LINE 13: Returns on first error
        }
    }
    for _, c := range s.nonCollectionPages {
        if err := c.Render(); err != nil {
            return err  // LINE 18: Returns on first error
        }
    }
    return nil
}
```

**Issue:** The render function returns immediately on the first error, preventing other files from being processed.

#### 3. Early Termination in Collection Rendering: `collection/collection.go:73-81`

```go
func (c *Collection) Render() error {
    for _, p := range c.Pages() {
        err := p.Render()
        if err != nil {
            return err  // LINE 77: Returns on first error
        }
    }
    return nil
}
```

**Issue:** Similar to the site render function, collection rendering stops at the first error.

#### 4. Early Termination in File Reading: `site/read.go:62-89`

```go
func (s *Site) readFiles(dir, base string) error {
    return filepath.Walk(dir, func(filename string, info os.FileInfo, err error) error {
        // ...
        d, err := pages.NewFile(s, filename, filepath.ToSlash(rel), defaultFrontmatter)
        if err != nil {
            return utils.WrapPathError(err, filename)  // LINE 81: Returns on error
        }
        // ...
    })
}
```

**Issue:** During file reading, any error in creating a file object stops the entire walk.

### Contrast: Good Error Handling Example

The **write phase** already handles multiple errors correctly in `site/write.go:30-52`:

```go
func (s *Site) WriteFiles() (count int, err error) {
    errs := make(chan error)
    // ...
    var errList []error
    for i := 0; i < count; i++ {
        if e := <-errs; e != nil {
            errList = append(errList, e)  // Collects ALL errors
        }
    }
    return count, combineErrors(errList)  // Combines them at the end
}
```

**Good Practice:** This code collects all errors and combines them using the `combineErrors()` function from `site/errors.go`, allowing all files to be processed.

## Proposed Solutions

### Option 1: Collect Errors During Rendering (Recommended)

Modify the rendering phase to collect all errors instead of returning on the first one:

**Files to modify:**
- `site/render.go`: Collect errors from collection and page rendering
- `collection/collection.go`: Collect errors from page rendering within a collection

**Advantages:**
- Consistent with the existing error handling in the write phase
- Shows all rendering errors at once
- Minimal changes to the codebase

**Implementation approach:**
1. Change `site.render()` to collect errors in a slice
2. Change `Collection.Render()` to collect errors in a slice
3. Use the existing `combineErrors()` utility to combine them

### Option 2: Improve Error Messages

Enhance the error messages to provide more context about what went wrong:

**Files to modify:**
- `renderers/markdown.go`: Add context about the markdown="1" attribute and suggest fixes

**Advantages:**
- Helps users understand what went wrong
- Can suggest fixes (e.g., use `<br/>` instead of `<br>`)

**Implementation approach:**
1. When an EOF error is encountered in `processInnerMarkdown()`, wrap it with a more descriptive message
2. Include the line number and context if possible
3. Suggest common fixes

### Option 3: Continue Processing on Errors During File Reading

Modify the file reading phase to continue even if some files cannot be parsed:

**Files to modify:**
- `site/read.go`: Collect errors instead of returning immediately
- `collection/read.go`: Collect errors instead of returning immediately

**Advantages:**
- Catches all problems upfront, even before rendering
- Allows users to see all problematic files at once

**Disadvantages:**
- More complex changes required
- May need to handle partial site state

### Option 4: Make HTML Parsing More Lenient

Make the HTML parser more forgiving of common mistakes like `<br>` vs `<br/>`:

**Files to modify:**
- `renderers/markdown.go`: Handle EOF errors more gracefully in `processInnerMarkdown()`

**Advantages:**
- More user-friendly, matches how browsers handle HTML
- Reduces errors users encounter

**Disadvantages:**
- May hide real problems
- Could lead to unexpected behavior

## Recommended Approach

**Combination of Options 1 and 2:**

1. **Collect errors during rendering** (Option 1) - This is the most important fix that addresses the core issue
2. **Improve error messages** (Option 2) - This helps users understand and fix the problems

This approach:
- Allows users to see all errors at once
- Provides helpful error messages
- Is consistent with existing code patterns (like the write phase)
- Requires moderate code changes
- Doesn't hide errors or make the parser too lenient

## Additional Findings

### Existing Error Combining Utility

The codebase already has a `combineErrors()` function in `site/errors.go:8-21`:

```go
func combineErrors(errs []error) error {
    switch len(errs) {
    case 0:
        return nil
    case 1:
        return errs[0]
    default:
        messages := make([]string, len(errs))
        for i, e := range errs {
            messages[i] = e.Error()
        }
        return errors.New(strings.Join(messages, "\n"))
    }
}
```

This utility is already used in the write phase and can be reused for the rendering and reading phases.

### Related Issues

- Issue #64: Support for `markdown="1"` attribute
- Issue #65: Support for `markdown="0|span|block"` attributes
- Issue #66: Parsing errors with `<br>` tags

These related issues suggest that proper support for the `markdown="1"` attribute (which is currently only partially implemented) would also help address this bug.

## Testing Strategy

To validate any fixes, we should:

1. Create test files with various parsing errors:
   - Malformed HTML with `<br>` tags
   - Missing closing tags
   - Invalid markdown attributes
   - Multiple problematic files

2. Verify that:
   - All files are processed
   - All errors are reported
   - Error messages are descriptive
   - The build/serve doesn't terminate early

3. Add regression tests to prevent the issue from recurring

## Conclusion

Issue #51 is caused by early termination in the rendering phase when encountering parse errors. The fix involves collecting errors instead of returning immediately, which is already done in the write phase. This investigation has identified the root cause and proposed practical solutions that can be implemented incrementally.
