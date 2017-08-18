---
---

jekyll.version = {{ jekyll.version }}

## Page Variables ([reference](https://jekyllrb.com/docs/variables/#page-variables))

| Name       | Value               |
|------------|---------------------|
| categories | {{ page.categories }} |
| date       | {{ page.date }}       |
| id         | {{ page.id }}         |
| path       | {{ page.path }}       |
| tags       | {{ page.tags }}       |
| title      | {{ page.title }}      |
| url        | {{ page.url }}        |

## Site Variables ([reference](https://jekyllrb.com/docs/variables/#site-variables))

| Name        | Value                | Notes                |
|-------------|----------------------|----------------------|
| description | {{ site.description }} | (from `_config.yml`) |
| time        | {{ site.time }}        |                      |
| title       | {{ site.title }}       | (from `_config.yml`) |

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

## Jekyll Variables

| Name        | Value                    |
|-------------|--------------------------|
| environment | {{ jekyll.environment }} |
| version     | {{ jekyll.version }}     |
