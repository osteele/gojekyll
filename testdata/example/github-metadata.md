---
title: GitHub Metadata
---

# {{ page.title }}

## Hash Style

<table>
{% for k in site.github %}
  <tr><th style="text-align: left">{{ k[0] }}</th><td>{{ k[1] }}</td></tr>
{% endfor %}
</table>

## List-of-Keys Style

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
