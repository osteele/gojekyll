// Command generate creates a large Jekyll site for benchmarking template caching.
//
// Each post triggers ~10 template parses (body + 3 layout levels + ~6 includes).
// With the default 2000 posts that means ~20,000 ParseTemplateLocation calls,
// but only ~15 unique templates — making the case for caching obvious in profiles.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var (
	numPosts  = flag.Int("posts", 2000, "number of posts to generate")
	numPages  = flag.Int("pages", 20, "number of standalone pages to generate")
	outputDir = flag.String("output", "benchmarks/testsite", "output directory")
)

func main() {
	flag.Parse()

	if err := generate(*outputDir, *numPosts, *numPages); err != nil {
		log.Fatal(err)
	}
}

func generate(dir string, posts, pages int) error {
	start := time.Now()

	if err := os.RemoveAll(dir); err != nil {
		return err
	}

	dirs := []string{
		filepath.Join(dir, "_layouts"),
		filepath.Join(dir, "_includes"),
		filepath.Join(dir, "_posts"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}

	// _config.yml
	if err := writeFile(filepath.Join(dir, "_config.yml"), configYML); err != nil {
		return err
	}

	// Layouts
	for name, content := range layouts {
		if err := writeFile(filepath.Join(dir, "_layouts", name), content); err != nil {
			return err
		}
	}

	// Includes
	for name, content := range includes {
		if err := writeFile(filepath.Join(dir, "_includes", name), content); err != nil {
			return err
		}
	}

	// index.html
	if err := writeFile(filepath.Join(dir, "index.html"), indexHTML); err != nil {
		return err
	}

	// Posts
	categories := []string{"tech", "science", "travel", "food", "music"}
	tags := []string{"go", "ruby", "python", "javascript", "rust", "benchmarks", "performance", "caching"}
	funcMap := template.FuncMap{"title": titleCase}
	postTmpl := template.Must(template.New("post").Funcs(funcMap).Parse(postTemplate))
	for i := range posts {
		date := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, i)
		data := postData{
			N:        i,
			Date:     date.Format("2006-01-02"),
			Category: categories[i%len(categories)],
			Tags:     []string{tags[i%len(tags)], tags[(i+3)%len(tags)]},
		}
		filename := fmt.Sprintf("%s-post-%04d.md", data.Date, i)
		path := filepath.Join(dir, "_posts", filename)
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		if err := postTmpl.Execute(f, data); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}

	// Standalone pages
	pageTmpl := template.Must(template.New("page").Parse(pageTemplate))
	for i := range pages {
		path := filepath.Join(dir, fmt.Sprintf("page-%03d.html", i))
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		if err := pageTmpl.Execute(f, i); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}

	elapsed := time.Since(start)
	fmt.Printf("Generated benchmark site in %s:\n", dir)
	fmt.Printf("  %d posts, %d pages\n", posts, pages)
	fmt.Printf("  ~%d template parses per build (only ~15 unique templates)\n", posts*10+pages*5)
	fmt.Printf("  Generation took %v\n", elapsed.Round(time.Millisecond))
	fmt.Println()
	fmt.Println("To benchmark:")
	fmt.Println("  go build && ./gojekyll benchmark -s", dir)
	fmt.Println("To profile:")
	fmt.Println("  go build && ./gojekyll benchmark -s", dir, "--profile")
	fmt.Println("  go tool pprof gojekyll.prof")
	return nil
}

type postData struct {
	N        int
	Date     string
	Category string
	Tags     []string
}

func titleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

// ---------- Static content ----------

const configYML = `title: Benchmark Site
description: A generated site for benchmarking gojekyll template caching
url: "http://example.com"
baseurl: ""
permalink: /:categories/:year/:month/:day/:title/

collections_dir: .
defaults:
  - scope:
      path: ""
      type: "posts"
    values:
      layout: "post"
  - scope:
      path: ""
    values:
      layout: "page"
`

var layouts = map[string]string{
	"base.html": `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>{{ page.title | default: site.title }}</title>
</head>
<body>
{% include header.html %}
<main>
  {{ content }}
</main>
{% include footer.html %}
</body>
</html>
`,
	"default.html": `---
layout: base
---
{% include sidebar.html %}
{% include breadcrumbs.html %}
<article>
  {{ content }}
</article>
`,
	"post.html": `---
layout: default
---
{% include post-meta.html %}
<div class="post-content">
  {{ content }}
</div>
{% include tag-list.html %}
{% include related-posts.html %}
{% include social-share.html %}
{% include pagination.html %}
`,
	"page.html": `---
layout: default
---
{% include toc.html %}
<div class="page-content">
  {{ content }}
</div>
`,
}

var includes = map[string]string{
	"header.html": `<header>
  <nav>
    <a href="/">{{ site.title }}</a>
    <ul>
    {% for p in site.pages %}
      {% if p.title %}
      <li><a href="{{ p.url }}">{{ p.title | escape }}</a></li>
      {% endif %}
    {% endfor %}
    </ul>
  </nav>
</header>
`,
	"footer.html": `<footer>
  <p>&copy; {{ site.time | date: "%Y" }} {{ site.title | escape }}</p>
  <p>{{ site.description | escape }}</p>
  <p>Built with <a href="https://github.com/osteele/gojekyll">gojekyll</a></p>
</footer>
`,
	"sidebar.html": `<aside class="sidebar">
  <h3>Recent Posts</h3>
  <ul>
  {% for post in site.posts limit:10 %}
    <li>
      <a href="{{ post.url | relative_url }}">{{ post.title | escape }}</a>
      <span>{{ post.date | date: "%b %-d, %Y" }}</span>
    </li>
  {% endfor %}
  </ul>
</aside>
`,
	"tag-list.html": `{% if page.tags.size > 0 %}
<div class="tags">
  <span>Tags:</span>
  {% for tag in page.tags %}
    <a href="/tags/{{ tag | slugify }}/">{{ tag | capitalize }}</a>
    {% unless forloop.last %}, {% endunless %}
  {% endfor %}
</div>
{% endif %}
`,
	"post-meta.html": `<div class="post-meta">
  <time datetime="{{ page.date | date_to_xmlschema }}">
    {{ page.date | date: "%B %-d, %Y" }}
  </time>
  {% if page.categories.size > 0 %}
  <span class="categories">
    in {% for cat in page.categories %}
      <a href="/categories/{{ cat | slugify }}/">{{ cat | capitalize }}</a>
      {% unless forloop.last %}, {% endunless %}
    {% endfor %}
  </span>
  {% endif %}
</div>
`,
	"pagination.html": `{% if paginator %}
<nav class="pagination">
  {% if paginator.previous_page %}
    <a href="{{ paginator.previous_page_path | relative_url }}">&laquo; Newer</a>
  {% else %}
    <span class="disabled">&laquo; Newer</span>
  {% endif %}
  <span>Page {{ paginator.page }} of {{ paginator.total_pages }}</span>
  {% if paginator.next_page %}
    <a href="{{ paginator.next_page_path | relative_url }}">Older &raquo;</a>
  {% else %}
    <span class="disabled">Older &raquo;</span>
  {% endif %}
</nav>
{% endif %}
`,
	"breadcrumbs.html": `{% assign crumbs = page.url | split: "/" %}
<nav class="breadcrumbs">
  <a href="/">Home</a>
  {% for crumb in crumbs %}
    {% if crumb != "" %}
      {% assign path = crumbs | slice: 0, forloop.index | join: "/" %}
      <span>/</span>
      <a href="/{{ path }}/">{{ crumb | replace: "-", " " | capitalize }}</a>
    {% endif %}
  {% endfor %}
</nav>
`,
	"related-posts.html": `{% if site.posts.size > 1 %}
<div class="related-posts">
  <h3>Related Posts</h3>
  <ul>
  {% for post in site.related_posts limit:5 %}
    <li>
      <a href="{{ post.url | relative_url }}">{{ post.title | escape }}</a>
      <span>{{ post.date | date: "%b %-d, %Y" }}</span>
    </li>
  {% endfor %}
  </ul>
</div>
{% endif %}
`,
	"social-share.html": `<div class="social-share">
  {% if page.title %}
    {% assign share_title = page.title | url_encode %}
  {% else %}
    {% assign share_title = site.title | url_encode %}
  {% endif %}
  {% assign share_url = page.url | absolute_url | url_encode %}
  {% unless page.hide_share %}
    <a href="https://twitter.com/intent/tweet?text={{ share_title }}&url={{ share_url }}">Twitter</a>
  {% endunless %}
</div>
`,
	"toc.html": `{% if page.toc %}
<nav class="toc">
  <h3>Table of Contents</h3>
  <ul>
  {% for section in page.toc %}
    <li>
      <a href="#{{ section.title | slugify }}">{{ section.title | escape }}</a>
      {% if section.children %}
      <ul>
        {% for child in section.children %}
        <li><a href="#{{ child.title | slugify }}">{{ child.title | escape }}</a></li>
        {% endfor %}
      </ul>
      {% endif %}
    </li>
  {% endfor %}
  </ul>
</nav>
{% endif %}
`,
}

const indexHTML = `---
layout: default
title: Home
---
<h1>{{ site.title }}</h1>
<p>{{ site.description }}</p>

<h2>Posts</h2>
<ul>
{% for post in site.posts %}
  <li>
    <a href="{{ post.url | relative_url }}">{{ post.title | escape }}</a>
    <span>{{ post.date | date: "%b %-d, %Y" }}</span>
    {% if post.categories.size > 0 %}
      <span>in {{ post.categories | join: ", " }}</span>
    {% endif %}
  </li>
{% endfor %}
</ul>
`

const postTemplate = `---
title: "Post Number {{ .N }}: Exploring {{ .Category | title }} Topics"
date: {{ .Date }}
categories: [{{ .Category }}]
tags: [{{ index .Tags 0 }}, {{ index .Tags 1 }}]
excerpt: "This is post number {{ .N }} about {{ .Category }}."
---

## Introduction to Post {{ .N }}

This is a benchmark post about **{{ .Category }}**. It exists to exercise
the Liquid template engine across layouts and includes.

### Section One

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod
tempor incididunt ut labore et dolore magna aliqua.

{{ "{{ page.title }}" }}  was published on {{ "{{ page.date | date: \"%B %-d, %Y\" }}" }}.

### Section Two

{{ "{% for tag in page.tags %}{{ tag | capitalize }}{% unless forloop.last %}, {% endunless %}{% endfor %}" }}

More content to give the Markdown renderer something to chew on.
Paragraph with **bold**, *italic*, and ` + "`code`" + ` formatting.

- List item one
- List item two
- List item three

> A blockquote to add variety to the rendered output.

### Conclusion

Post {{ .N }} is part of the {{ .Category }} category and is tagged with
{{ index .Tags 0 }} and {{ index .Tags 1 }}.
`

const pageTemplate = `---
layout: page
title: "Page {{ . }}"
---

## Page {{ . }}

This is standalone page number {{ . }}. It uses the page layout,
which includes the table-of-contents include.

### First Section

Content for the first section of page {{ . }}.

### Second Section

More content to exercise the template engine with page layout rendering.
`
