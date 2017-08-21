---
permalink: /
variable: page variable
---

## Stats

* {{ site.collections.size }} Collections
* {{ site.documents.size }} Documents
* {{ site.html_files.size }} HTML File(s)
* {{ site.html_pages.size }} HTML Pages
* {{ site.pages.size }} Pages
* {{ site.posts.size }} Post(s)
* {{ site.static_files.size }} Static Files

## Tests

* [Archive]({% link archive.md %})
* [Collections]({% link collections.html %})
* [Data]({% link data.md %})
* [Markdown]({% link markdown.md %})
* [Pages]({% link pages.md %})
* [Plugins]({% link plugins.md %})
  * [GitHub Metadata]({% link github-metadata.md %})
* [Static file]({% link static.html %})
* [Tags]({% link tags.md %})
* [Variables]({% link variables.md %})

## Variables

| &nbsp;          | &nbsp;                       |
|-----------------|------------------------------|
| Page variable   | {{ page.variable }}          |
| Site variable   | {{ site.variable }}          |
| Data file value | {{ site.data.file_data.k1 }} |
