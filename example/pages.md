---
---

# Pages

{% for page in site.pages %}
* {{ page.path }}{% endfor %}

## Page Variables

{{ site.pages | first }}
