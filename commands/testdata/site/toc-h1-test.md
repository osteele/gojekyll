---
title: TOC H1 Inclusion Test
layout: default
kramdown:
  toc_levels: 1..6
---

# Page Title (H1)

This page specifically tests whether H1 headings are included in the TOC.

* Table of Contents
{:toc}

## First Section (H2)

According to the GitHub issue comment, the TOC "shouldn't normally include the h1 with the page title, I think, but it does".

This is because Jekyll's default `toc_levels` is "2..6", which excludes H1 headings.

### Subsection 1.1 (H3)

Content here.

## Second Section (H2)

More content.

### Subsection 2.1 (H3)

Even more content.

## Expected Behavior

1. **With default Jekyll config** (toc_levels: 2..6):
   - H1 "Page Title" should NOT appear in TOC
   - Only H2 and H3 headings should appear

2. **With toc_levels: 1..6**:
   - H1 "Page Title" SHOULD appear in TOC
   - All heading levels appear

3. **Current gojekyll behavior** (as reported in issue):
   - H1 appears in TOC even when it shouldn't
   - This suggests gojekyll may not be respecting the default toc_levels: 2..6
