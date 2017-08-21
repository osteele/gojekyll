---
---

## Site Variables ([reference](https://jekyllrb.com/docs/variables/#site-variables))

| Name              | Value                                                            | Notes                |
|-------------------|------------------------------------------------------------------|----------------------|
| description       | {{ site.description }}                                           | (from `_config.yml`) |
| time              | {{ site.time }}                                                  |                      |
| title             | {{ site.title }}                                                 | (from `_config.yml`) |
| pages.size        | {{ site.pages.size }}                                            |                      |
| posts.size        | {{ site.posts.size }}                                            |                      |
| static_files.size | {{ site.static_files.size }}                                     |                      |
| html_pages.size   | {{ site.html_pages.size }}                                       |                      |
| html_files.size   | {{ site.html_files.size }}                                       |                      |
| collections.size  | {{ site.collections.size }}                                      |                      |
| documents.size    | {{ site.documents.size }}                                        |                      |
| categories        | {{ site.categories }} {% for k in categories %}{{k}}{% endfor %} |                      |
| tags              | {{ site.tags }}                                                  |                      |

{% if false %}
{% capture ks %}{% for k in site %}{{ k }} {% endfor %}{% endcapture %}
{% assign ks = ks | split: ' ' | sort %}

<table>{% for k in ks %}
  <tr>
    <td style="text-align: left; vertical-align: top">{{k}}</td>
    <td><pre>{{site[k]|escape|truncate: 80}}</pre></td>
  </tr>
{% endfor %}
</table>
{% endif %}

## Page Variables ([reference](https://jekyllrb.com/docs/variables/#page-variables))

| Name       | Value                 |
|------------|-----------------------|
| categories | {{ page.categories }} |
| date       | {{ page.date }}       |
| id         | {{ page.id }}         |
| path       | {{ page.path }}       |
| tags       | {{ page.tags }}       |
| title      | {{ page.title }}      |
| url        | {{ page.url }}        |

## Jekyll Variables

| Name        | Value                    |
|-------------|--------------------------|
| environment | {{ jekyll.environment }} |
| version     | {{ jekyll.version }}     |
