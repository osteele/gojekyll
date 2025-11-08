---
layout: default
---

Gojekyll is a Go clone of the [Jekyll](https://jekyllrb.com) static site generator.

Refer to the [README](/README.md) for more information.

This page exists mostly to dogfood site generation and themes, and also serves as an entry point to deeper documentation.

## Getting Started

{% for p in site.html_pages %}
{% if p.url != page.url and (p.name contains "benchmark" or p.name contains "plugin") %}
* [{{ p.name | remove: ".md" | capitalize }}]({{ p.url }})
{% endif %}
{% endfor %}

## Reference Documentation

* [Configuration](configuration.html) - Configuration options for `_config.yml`
* [Permalinks](PERMALINKS.html) - Permalink handling and URL patterns
{% for p in site.html_pages %}
{% if p.url != page.url and p.name != "benchmarks.md" and p.name != "plugins.md" and p.name != "configuration.md" and p.name != "PERMALINKS.md" %}
* [{{ p.name | remove: ".md" | capitalize }}]({{ p.url }})
{% endif %}
{% endfor %}
