---
permalink: /:name
---

# Archive

{% assign pages = site.c1 | sort: 'weight' %}
{% for p in pages %}
* {{p.title}}
{% endfor %}
