---
title: TOC Test Page
layout: default
---

# Test Page for Table of Contents

This page tests various TOC behaviors to verify compatibility with Ruby Jekyll.

## Test 1: Basic TOC with inline syntax

* Table of Contents
{:toc}

### Subsection 1.1

This is subsection 1.1 content.

### Subsection 1.2

This is subsection 1.2 content.

## Test 2: Heading exclusion with .no_toc

### This section should appear in TOC
{:.no_toc}

This heading has the `{:.no_toc}` marker and should NOT appear in the table of contents.

### This section should also appear

Normal heading that should be in the TOC.

## Test 3: Block syntax

Below is a TOC using block syntax:

{::toc}

### Section 3.1

Content for section 3.1.

### Section 3.2

Content for section 3.2.

## Test 4: Different list placeholder text

* Contents will be replaced
{:toc}

### Section 4.1

Some content here.

## Test 5: H1 inclusion behavior

According to the GitHub issue comment, the TOC shouldn't normally include the H1 (page title), but gojekyll currently does. This is likely because Jekyll's default `toc_levels` is "2..6" which excludes H1 headings.

## Test 6: Multiple heading levels

### Level 3 heading
#### Level 4 heading
##### Level 5 heading
###### Level 6 heading

## Notes on Expected Behavior

Based on Jekyll/Kramdown documentation:

1. **{:toc}** - Inline attribute list syntax. When used after a list item in an unordered list, it should replace that list item with the TOC.

2. **{::toc}** - Block syntax. Should generate a TOC at that location without needing a list item.

3. **{:.no_toc}** - When placed after a heading, that heading should be excluded from the TOC.

4. **Default toc_levels**: Jekyll's default is "2..6", meaning H1 headings are excluded by default.

5. **List replacement**: Only works with `{:toc}` in unordered lists, not with `{::toc}` or in ordered lists.

## Kramdown Configuration

You can configure TOC generation in `_config.yml`:

```yaml
kramdown:
  toc_levels: 1..6  # Include all heading levels
  # or
  toc_levels: 2..3  # Only include H2 and H3
```
