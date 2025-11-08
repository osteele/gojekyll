# Permalink Handling in Gojekyll

This document explains how gojekyll handles permalinks and the important distinctions between different document types. This behavior matches Jekyll's permalink handling for compatibility.

## Overview

Gojekyll supports Jekyll's permalink configuration system, allowing URLs to be customized via:
1. Front matter `permalink` field (highest priority)
2. Global `permalink` configuration in `_config.yml`
3. Default permalink pattern

## Document Type Distinctions

### Critical Concept: Posts vs Pages vs Collections

Jekyll (and `gojekyll`) treats different document types differently when processing permalink patterns. This distinction is intentional and required for Jekyll compatibility.

#### Posts (collection == "posts")
- **Full permalink pattern support**: All placeholders are honored
- **Available placeholders**: `:year`, `:month`, `:day`, `:categories`, `:title`, etc.
- **Example**: `permalink: pretty` → `/:categories/:year/:month/:day/:title/`
- **Result**: `/blog/2024/03/15/my-post/`

#### Pages (regular pages, not in any collection)
- **Date and category placeholders are IGNORED**
- **Available placeholders**: `:path`, `:basename`, `:output_ext`, `:title`
- **Example**: `permalink: pretty` → `/:title/` (dates/categories removed)
- **Result**: `/about/` (not `/2024/03/15/about/`)

#### Collection Documents (non-post collections)
- **Same as pages**: Date and category placeholders are ignored
- **Additional placeholder**: `:collection`
- **Example**: `authors` collection with `permalink: pretty` → `/:title/`

## Why This Distinction Exists

According to Jekyll's official documentation:
> "Pages and collections (excluding posts and drafts) don't have time and categories, so aspects of the permalink style are ignored for the output."

This makes sense because:
- **Posts** are inherently date-based content (blog posts, news articles)
- **Pages** are timeless content (About, Contact, Services pages)
- **Collections** can be either, but default to timeless

## Examples

### Config: `permalink: pretty`

| Document Type | Source File | Jekyll URL | Pattern After Processing |
|--------------|-------------|------------|-------------------------|
| Post | `_posts/2024-03-15-hello.md` | `/2024/03/15/hello/` | `/:year/:month/:day/:title/` |
| Page | `about.md` | `/about/` | `/:title/` |
| Collection | `_authors/john.md` | `/john/` | `/:title/` |

### Config: `permalink: date`

| Document Type | Source File | Jekyll URL | Pattern After Processing |
|--------------|-------------|------------|-------------------------|
| Post | `_posts/2024-03-15-hello.md` | `/2024/03/15/hello.html` | `/:year/:month/:day/:title:output_ext` |
| Page | `about.md` | `/about.html` | `/:title:output_ext` |

## Common Pitfalls

1. **Expecting dates in page URLs**: Pages don't have dates, even if you set `date` in front matter
2. **Categories on pages**: Categories are ignored for pages, only work for posts
3. **Front matter permalink in global config**: Built-in styles (pretty, date, etc.) only work in `_config.yml`, not in front matter

## Testing

When testing permalink behavior:
1. Test posts separately from pages
2. Test with and without categories
3. Test all built-in permalink styles
4. Test custom patterns with various placeholders

---

## For Maintainers

### Implementation Notes

⚠️ **DO NOT REMOVE** the distinction between posts and pages in permalink processing. It's required for Jekyll compatibility.

### References

- The Jekyll documentation at https://jekyllrb.com/docs/permalinks/

### Related Files

- `pages/permalinks.go` - Core permalink processing logic
- `pages/permalinks_test.go` - Comprehensive tests
- `config/default.go` - Default permalink configuration
- `collection/collection.go` - Collection-specific permalink handling
