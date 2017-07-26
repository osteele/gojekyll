---
---

## YAML data

* {{ site.data.file_data.key }}

## JSON data

{{ site.data.json_data.key }}

## CSV data

<table>
{% for row in site.data.csv_data %}
  <tr>{% for cell in row %}<td>{{ cell }}</td>{% endfor %}</tr>
{% endfor %}
</table>
