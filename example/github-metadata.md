---
title: GitHub Metadata
---

# {{ page.title }}

{% if jekyll.version contains 'gojekyll' %}

<table>
{% for k in site.github %}
  <tr>
    <th style="text-align: left; vertical-align: top">{{ k }}</th>
    <td>
      <div style="max-height: 100px; overflow: scroll">{{ site.github[k] }}
    </td>
  </tr>
{% endfor %}
</table>

{% endif %}
