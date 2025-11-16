---
title: TOC Edge Cases Test
layout: default
---

# TOC Edge Cases

This page tests edge cases for TOC generation.

## Test 1: TOC in Fenced Code Block

The marker should appear literally, not generate a TOC:

```
{:toc}
```

## Test 2: TOC in Inline Code

Use `{:toc}` to generate a table of contents.

## Test 3: Actual TOC

This should work:

* Table of Contents
{:toc}

## Test 4: Empty Headings Document

Create a page with TOC but content below has no headings:

---

# Another Document Start

{:toc}

Just text, no headings to list.

---

## Test 5: Single Heading

When there's only one heading besides the page title:

### Only Subsection

This is the only subsection.

## Test 6: TOC Markers in HTML Comments

<!-- {:toc} should not generate TOC when in a comment -->

## Test 7: Escaped Markers

These should not generate TOCs (if properly escaped in markdown):

\{:toc\} - escaped braces

## Test 8: TOC After All Content

All content first:

### Section A
Content A

### Section B
Content B

### Section C
Content C

Now the TOC at the end:

{::toc}

## Summary

This page tests that TOC markers:
- In code blocks are displayed literally
- In inline code are displayed literally
- In HTML comments are ignored
- Work correctly when placed at various positions
