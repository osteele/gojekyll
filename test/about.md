---
permalink: /:name
---

# About

{% assign pages = site.collection | sort: 'weight' %}
{% for p in pages %}
* {{p.title}}
{% endfor %}
