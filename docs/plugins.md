# Gojekyll Plugin Status

Gojekyll doesn't include an extensible plugin system, and won't for the foreseeable future.
((1) The code and internal APIs are too immature for this; and (2) The [natural way](https://golang.org/pkg/plugin/) of implementing this only works on Linux.)

The executable emulates some plugins:

| Plugin                       | Motivation    | Basic Functionality | Missing Features                        |
|------------------------------|---------------|---------------------|-----------------------------------------|
| jekyll-avatar                | GitHub Pages¹ | ✓                   | randomized hostname                     |
| jekyll-coffeescript          | GitHub Pages  |                     |                                         |
| jekyll-default-layout        | GitHub Pages  |                     |                                         |
| jekyll-feed                  | GitHub Pages  | ✓                   |                                         |
| jekyll-gist                  | core²         | ✓                   | `noscript`                              |
| jekyll-github-metadata       | GitHub Pages  |                     |                                         |
| jekyll-live-reload           | core          | ✓ always on         |                                         |
| jekyll-mentions              | GitHub Pages  | ✓                   |                                         |
| jekyll-optional-front-matter | GitHub Pages  |                     |                                         |
| jekyll-paginate              | core          |                     |                                         |
| jekyll-readme-index          | GitHub Pages  |                     |                                         |
| jekyll-redirect_from         | GitHub Pages  | ✓                   | user template                           |
| jekyll-relative-links        | GitHub Pages  |                     |                                         |
| jekyll-sass-converter        | core          | ✓ always on         |                                         |
| jekyll-seo_tag               | GitHub Pages  | ✓                   | SEO and JSON LD are not fully populated |
| jekyll-sitemap               | GitHub Pages  |                     |                                         |
| jemoji                       | GitHub Pages  | ✓                   | image tag fallback                      |

¹ <https://pages.github.com/versions/>

² “Core” plugins are referenced in the main [Jekyll documentation](https://jekyllrb.com/docs/home/).
The Jekyll documentation [Official Plugins](https://jekyllrb.com/docs/plugins/#available-plugins) / #Official tag of [Awesome Jekyll Plugins](https://github.com/planetjekyll/awesome-jekyll-plugins) look dated; I didn't use those.