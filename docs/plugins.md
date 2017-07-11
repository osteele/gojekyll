# Gojekyll Plugin Status

Gojekyll doesn't include¹ an extensible plugin system, and won't for the foreseeable future.

The functionality of some plugins is built into the core program:

| Plugin                       | Motivation    | Basic Functionality | Missing Features                        |
|------------------------------|---------------|---------------------|-----------------------------------------|
| jekyll-avatar                | GitHub Pages² | ✓                   | randomized hostname                     |
| jekyll-coffeescript          | GitHub Pages  |                     |                                         |
| jekyll-default-layout        | GitHub Pages  |                     |                                         |
| jekyll-feed                  | GitHub Pages  | ✓                   |                                         |
| jekyll-gist                  | core³         | ✓                   | `noscript`                              |
| jekyll-github-metadata       | GitHub Pages  |                     |                                         |
| jekyll-live-reload           | core          | ✓ (always enabled)  |                                         |
| jekyll-mentions              | GitHub Pages  | ✓                   |                                         |
| jekyll-optional-front-matter | GitHub Pages  |                     |                                         |
| jekyll-paginate              | core          |                     |                                         |
| jekyll-readme-index          | GitHub Pages  |                     |                                         |
| jekyll-redirect_from         | GitHub Pages  | ✓                   | user template                           |
| jekyll-relative-links        | GitHub Pages  |                     |                                         |
| jekyll-sass-converter        | core          | ✓ (always enabled)  |                                         |
| jekyll-seo_tag               | GitHub Pages  | ✓                   | SEO and JSON LD are not fully populated |
| jekyll-sitemap               | GitHub Pages  |                     |                                         |
| jemoji                       | GitHub Pages  | ✓                   | image tag fallback                      |

¹ (1) The code and internal APIs are too immature for this; and (2) The [natural way](https://golang.org/pkg/plugin/) of implementing this only works on Linux.

² <https://pages.github.com/versions/>

³ “Core” plugins are referenced in the main [Jekyll documentation](https://jekyllrb.com/docs/home/).
The Jekyll documentation [Official Plugins](https://jekyllrb.com/docs/plugins/#available-plugins) / #Official tag of [Awesome Jekyll Plugins](https://github.com/planetjekyll/awesome-jekyll-plugins) look dated; I didn't use those.