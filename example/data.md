---
---

## YAML data

{% assign data = site.data.file_data %}
<table>
  {% for k in data %}
    <tr><td>{{ k[0] }}</td><td>{{ k[1] }}</td>
  {% endfor %}
</table>

## JSON data

{% assign data = site.data.json_data %}
<table>
  {% for k in data %}
    <tr><td>{{ k[0] }}</td><td>{{ k[1] }}</td>
  {% endfor %}
</table>

## CSV data

<table>
{% for row in site.data.csv_data %}
  <tr>{% for cell in row %}<td>{{ cell }}</td>{% endfor %}</tr>
{% endfor %}
</table>
