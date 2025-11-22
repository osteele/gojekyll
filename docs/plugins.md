# Gojekyll Plugin Status

Gojekyll doesn't include an extensible plugin system¹.

The functionality of some plugins is built into the core program:

| Plugin                                                       | Motivation    | Implementation Status | Missing Features                                                                                                                      |
|--------------------------------------------------------------|---------------|-----------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| [jekyll-avatar][jekyll-avatar]                               | GitHub Pages² | ✓                     |                                                                                                                                       |
| [jekyll-coffeescript][jekyll-coffeescript]                   | GitHub Pages  |                       |                                                                                                                                       |
| [jekyll-default-layout][jekyll-default-layout]               | GitHub Pages  | ✓                     |                                                                                                                                       |
| [jekyll-feed][jekyll-feed]                                   | GitHub Pages  | ✓                     |                                                                                                                                       |
| [jekyll-gist][jekyll-gist]                                   | core³         | ✓                     | `noscript` option                                                                                                                     |
| [jekyll-github-metadata][jekyll-github-metadata]             | GitHub Pages  | partial               | `contributors`, `public_repositories`, `show_downloads`, `releases`, `versions`, `wiki_url`; Octokit configuration; GitHub Enterprise |
| [jekyll-live-reload][jekyll-live-reload]                     | core          | ✓                     | always enabled (by design); no way to disable                                                                                         |
| [jekyll-mentions][jekyll-mentions]                           | GitHub Pages  | ✓                     |                                                                                                                                       |
| [jekyll-optional-front-matter][jekyll-optional-front-matter] | GitHub Pages  |                       |                                                                                                                                       |
| [jekyll-paginate][jekyll-paginate]                           | core          |                       |                                                                                                                                       |
| [jekyll-readme-index][jekyll-readme-index]                   | GitHub Pages  |                       |                                                                                                                                       |
| [jekyll-redirect_from][jekyll-redirect_from]                 | GitHub Pages  | ✓                     | user template                                                                                                                         |
| [jekyll-relative-links][jekyll-relative-links]               | GitHub Pages  | ✓                     |                                                                                                                                       |
| [jekyll-sass-converter][jekyll-sass-converter]               | core          | ✓                     | always enabled (by design); no way to disable                                                                                         |
| [jekyll-seo_tag][jekyll-seo_tag]                             | GitHub Pages  | partial               | `dateModified`, `datePublished`, `publisher`, `mainEntityOfPage`, `@type`                                                             |
| [jekyll-sitemap][jekyll-sitemap]                             | GitHub Pages  | ✓                     | file modified dates⁴                                                                                                                  |
| [jekyll-titles-from-headings][jekyll-titles-from-headings]   | GitHub Pages  |                       |                                                                                                                                       |
| [jemoji][jemoji]                                             | GitHub Pages  | ✓                     | image tag fallback                                                                                                                    |
| [GitHub pages][github-pages]                                 | GitHub Pages  | ✓                     | The plugins that github-pages *includes* are in various stages of implementation, listed above                                        |

¹ (1) The code and internal APIs are too immature for this; and (2) the [natural way](https://golang.org/pkg/plugin/) of implementing this only works on Linux.

² <https://pages.github.com/versions/>

³ “Core” plugins are referenced in the main [Jekyll documentation](https://jekyllrb.com/docs/home/).

The [Official Plugins](https://jekyllrb.com/docs/plugins/#available-plugins) section of the Jekyll documentation, and the #Official tag of [Awesome Jekyll Plugins](https://github.com/planetjekyll/awesome-jekyll-plugins), look dated; I didn't use those.

⁴ These don't seem that useful with source control and CI. (Post dates are included.)

[jekyll-avatar]: https://github.com/benbalter/jekyll-avatar
[jekyll-coffeescript]: https://github.com/jekyll/jekyll-coffeescript
[jekyll-default-layout]: https://github.com/benbalter/jekyll-default-layout
[jekyll-feed]: https://github.com/jekyll/jekyll-feed
[jekyll-gist]: https://github.com/jekyll/jekyll-gist
[jekyll-github-metadata]: https://github.com/parkr/github-metadata
[jekyll-live-reload]: https://github.com/RobertDeRose/jekyll-livereload
[jekyll-mentions]: https://github.com/jekyll/jekyll-mentions
[jekyll-optional-front-matter]: https://github.com/benbalter/jekyll-optional-front-matter
[jekyll-paginate]: https://github.com/jekyll/jekyll-paginate
[jekyll-readme-index]: https://github.com/benbalter/jekyll-readme-index
[jekyll-redirect_from]: https://github.com/jekyll/jekyll-redirect-from
[jekyll-relative-links]: https://github.com/benbalter/jekyll-relative-links
[jekyll-sass-converter]: https://github.com/jekyll/jekyll-sass-converter
[jekyll-seo_tag]: https://github.com/jekyll/jekyll-seo-tag
[jekyll-sitemap]: https://github.com/jekyll/jekyll-sitemap
[jekyll-titles-from-headings]: https://github.com/benbalter/jekyll-titles-from-headings
[jemoji]: https://github.com/jekyll/jemoji
[github-pages]: https://github.com/github/pages-gem
