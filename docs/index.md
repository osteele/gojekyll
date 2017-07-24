---
layout: default
---

Gojekyll is a Go clone of the [Jekyll](https://jekyllrb.com) static site generator.

Refer to the [README](https://github.com/osteele/gojekyll/settings) for more information.

This page exists mostly to dogfood site generation and themes, although it also hosts:

{% for p in site.html_pages %}
{% if p.url != page.url %}
* [{{ p.name | remove: ".md" | capitalize }}]({{ p.url }})
{% endif %}
{% endfor %}
