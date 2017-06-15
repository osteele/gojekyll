---
permalink: /:name
layout: default
---

## Pages

{% assign pages = site.c1 | sort: 'weight' %}
{% for p in pages %}
* [{{p.title}}]({{p.url}})
{% endfor %}

[An article]({% link _c1/c1.md %})
