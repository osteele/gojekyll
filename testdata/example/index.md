---
permalink: /
variable: page variable
---

# Site Title

Here is a test with a {{ site.variable}}, a {{ page.variable }} and [a link]({% link archive.md %}).

## Site Variables [reference](https://jekyllrb.com/docs/variables/#site-variables)

| Name | Value | Notes |
| --- | --- | --- |
| title | {{site.title}} | (from `_config.yml`)
| description | {{site.description}} | (from `_config.yml`)
| time | {{site.time}} |

{% comment %}
TODO:
pages
posts
related_posts
static_files
html_pages
html_files
collections
data
documents
categories
tags
{% endcomment %}

## Page Variables [reference](https://jekyllrb.com/docs/variables/#page-variables)

| Name | Value |
| --- | --- |
| title | {{page.title}} |
| url | {{page.url}} |
| date | {{page.date}} |
| id | {{page.id}} |
| categories | {{page.categories}} |
| tags | {{page.tags}} |
| path | {{page.path}} |


{% comment %}
TODO:
excerpt
content
next
previous
{% endcomment %}
