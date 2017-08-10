---
---

## Collections

<dl>
{% for c in site.collections %}
  <dt> {{ c.label }}</dt>
  {% for k in c %}{% if k != 'docs' %}
  <dd>{{k}}={{c[k]}}
  {% endif %}{% endfor %}
  <dd>docs: <ul>
    {% for p in c.docs %}
      <li>{{ p.path }}
      properties: {% for k in p %}{{k}} {% endfor %}
      </li>
      <pre>{{p.content | escape}}</pre>
    {% endfor %}
  </ul></dd>
{% endfor %}
</dl>
