# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- **Liquid Engine**: Updated liquid template engine from v1.6.0 to v1.8.0. This brings Unicode identifier support, `LaxFilters` option, date timestamp support, and Jekyll-specific extensions including dot notation in assign tags.

### Added

- **Math Support** (#110): Added MathJax/KaTeX compatibility for mathematical expressions using `$$...$$` delimiters, compatible with Jekyll/kramdown syntax
- **`jekyll-relative-links` Plugin** (#104, #25): Converts relative markdown links to their rendered equivalents
- **`jekyll-readme-index` Plugin** (#106, #29): Remaps README files to index pages
- **`jekyll-gist` Noscript** (#105, #27): Added `noscript` option for the `jekyll-gist` plugin
- **`--baseurl` and `--config` CLI Flags** (#103, #17, #18): Added support for `--baseurl` to override the site base URL and `--config` to specify alternate config files
- **`sassify` Filter** (#109): Implemented the `sassify` Liquid filter for converting indented Sass syntax to CSS
- **Table of Contents (TOC) Support** (#76, #101, #62): Added Kramdown-style TOC generation with `{:toc}` and `{::toc}` markers, including support for Jekyll's `toc_levels` configuration and heading exclusion with `{:.no_toc}`. Thanks [@tekknolagi](https://github.com/tekknolagi) for requesting
- **Permalink Timezone Configuration** (#67): Added `permalink_timezone` configuration option to control timezone for permalink date generation
- **Markdown Attributes Support** (#85, #64): Added support for full Kramdown markdown attribute syntax (`markdown=1`, `markdown=0`, `markdown=block`, `markdown=span`) in HTML blocks

### Fixed

- **`page.date` for Non-Posts** (#116, #115): `page.date` is now only defined for posts and collection documents, or when explicitly set in frontmatter; previously it was unconditionally set to the file modification time. Thanks [@sampsyo](https://github.com/sampsyo) for reporting
- **Indented HTML Rendered as Code** (#117, #113): Fixed indented HTML block-level elements inside list items being erroneously rendered as code blocks after the switch to Goldmark. Thanks [@tekknolagi](https://github.com/tekknolagi) for reporting
- **`{:.no_toc}` Paragraphs** (#112): Fixed `{:.no_toc}` attribute markers being left as visible paragraphs in the HTML output instead of being removed
- **Sass Error Handling** (#99, #95): Fixed "connection is shut down" error when compiling SCSS by using a global singleton for the Sass transpiler; added helpful error message when wrong Sass package is installed
- **TOC List Replacement** (#93, #89): Fixed TOC to replace adjacent lists correctly, matching Jekyll's exact behavior. Thanks [@tekknolagi](https://github.com/tekknolagi) for reporting
- **SCSS Compilation Error** (#92, #90): Fixed "connection is shut down" error when compiling SCSS. Thanks [@tekknolagi](https://github.com/tekknolagi) for reporting
- **Custom Permalink Handling** (#82, #81): Fixed issue where `index.md` was not being rendered when custom permalink patterns were set in `_config.yml`. Custom permalink patterns now only apply to posts, not pages. Thanks [@tekknolagi](https://github.com/tekknolagi) for reporting
- **Canonical URL in SEO Plugin** (#72, #70): Fixed jekyll-seo-tag plugin to respect page's `canonical_url` front matter instead of always auto-generating. Thanks [@tekknolagi](https://github.com/tekknolagi) for reporting
- **Page Permalink Configuration** (#73, #71, #74, #61): Fixed pages to respect global permalink configuration from `_config.yml`, with proper handling of directory-style permalinks and URL routing without trailing slashes. Thanks [@tekknolagi](https://github.com/tekknolagi) for requesting
- **File Watching Issues** (#84): Fixed multiple critical bugs in file watching, dry-run mode, and live-reload including stale site references, missing render during dry-run, stale Sass partials, and spurious live-reload with `--no-watch`
- **First Parse Error Handling** (#79, #51): Changed build and serve commands to collect all rendering errors instead of stopping at the first error, making it easier to identify and fix all issues at once. Thanks [@manastungare](https://github.com/manastungare) for reporting
- **Symlink Preservation** (#80, #48): Fixed issue where `_site` directory symlinks were replaced with regular directories instead of following the symlink target. Thanks [@edgan](https://github.com/edgan) for reporting
- **URL Routing** (#74, #52): Fixed server to correctly handle URLs without trailing slashes for directory-style permalinks. Thanks [@abhijeetbodas2001](https://github.com/abhijeetbodas2001) for reporting
- **Layout Handling** (#78): Fixed pages with `layout: none` or `layout: null` in front matter to skip layout rendering instead of causing errors
- **First Build Crash**: Fixed Clean function crash when destination directory doesn't exist on first run
- **Windows Support** (#96): Fixed URL routing, path handling, and test failures on Windows
- **Nested Directory Watching**: Fixed file watcher to recursively watch nested directories and detect changes in subdirectories
- **Config Updates**: Fixed `Config.Set` to properly update YAML MapSlice so template changes are observed correctly

### Changed

- **Error Handling** (#97): Replaced `log.Fatal` calls with `panic` and `fmt.Errorf` for proper error propagation
- **Logging System** (#75, #35): Replaced scattered `fmt.Printf` statements with centralized logging package supporting proper log levels (Debug, Info, Warning, Error) and quiet mode
- **File Watcher**: Improved file watcher with automatic fallback to polling when directory count exceeds 500, preventing file descriptor exhaustion on large sites
- **Error Messages**: Enhanced markdown renderer error messages for common issues (e.g., suggesting `<br/>` instead of `<br>`)

### Maintenance

- **Go 1.25/1.26** (#119): Updated CI to test against Go 1.25 and 1.26; updated minimum Go version to 1.25
- **golangci-lint v2** (#108): Updated golangci-lint configuration for v2
- **GitHub Actions** (#87): Updated CI workflows to test on Ubuntu, macOS, and Windows; updated actions to latest versions
- **Tests** (#102): Added tests for `jekyll-default-layout` plugin
- **Code Quality**: Fixed lint issues, ran go fmt for consistent formatting
- **Documentation**: Improved documentation structure and clarity, added configuration documentation
- **.gitignore**: Updated to exclude Go build cache and macOS-specific files

## [0.2.16] - 2025-06-01

### Changed

- Updated liquid template engine dependency

## [0.2.15] - 2025-06-01

### Fixed

- Fixed linter errors (#69)
- Fixed tests to pass in all environments

### Maintenance

- Updated GitHub Actions workflow to enforce strict linting

## [0.2.14] - 2024-10-28

### Maintenance

- Tidied dependencies

## [0.2.13] - 2024-10-28

### Changed

- Updated dependencies (#59). Thanks [@danog](https://github.com/danog)

## [0.2.12] - 2024-10-17

### Maintenance

- Improved build script

## [0.2.11] - 2024-10-17

### Maintenance

- Fixed CI workflow
- Updated dependencies
- Bumped versions

## [0.2.10] - 2024-04-21

### Added

- **Footnote Support** (#58): Added support for footnotes in markdown files to match Jekyll behavior. Thanks [@sirwart](https://github.com/sirwart)

### Fixed

- **Slug Generation** (#56): Fixed slug generation to remove leading and trailing hyphens to match Jekyll behavior. Thanks [@sirwart](https://github.com/sirwart)

### Maintenance

- Updated golangci-lint version
- Improved build process

## [0.2.9] - 2023-11-17

### Added

- **Extensionless URL Serving** (#54): Server now serves extensionless URLs like `/some-url` from files like `/some-url.html`. Thanks [@chimbori](https://github.com/chimbori)

### Fixed

- **Test Coverage** (#55): Fixed test that wasn't calling the function being tested (`mustMarkdownString`). Thanks [@chimbori](https://github.com/chimbori)

### Maintenance

- Improved documentation and README

## [0.2.8] - 2023-08-26

### Added

- Added Docker image support with multi-architecture builds (amd64, armv5, armv6, armv7)

### Maintenance

- Improved release process

## [0.2.7] - 2023-08-26

### Changed

- Switched from Ruby Sass to Dart Sass

### Maintenance

- Improved GitHub Actions workflows
- Updated test infrastructure

## [0.2.6] - 2023-08-23

### Changed

- Updated dependencies (#49). Thanks [@danog](https://github.com/danog)

## Earlier Releases

For releases prior to v0.2.6, please see the [GitHub Releases page](https://github.com/osteele/gojekyll/releases).

Notable earlier releases:
- **v0.2.5** (2017-08-18): Renamed pipeline to renderers
- **v0.2.4** (2017-08-10): Render non-collection pages
- **v0.2.3** (2017-08-08): Better reload functionality
- **v0.2.2** (2017-08-03): Fixed race condition
- **v0.2.1** (2017-07-26): Tweaked in-page error display
- **v0.2.0** (2017-07-25): Created PageEmbed feature
- **v0.1.1** (2017-07-19): Updated goreleaser version varname target
- **v0.1.0** (2017-07-17): Push site build errors to open web pages

[Unreleased]: https://github.com/osteele/gojekyll/compare/v0.2.16...HEAD
[0.2.16]: https://github.com/osteele/gojekyll/compare/v0.2.15...v0.2.16
[0.2.15]: https://github.com/osteele/gojekyll/compare/v0.2.14...v0.2.15
[0.2.14]: https://github.com/osteele/gojekyll/compare/v0.2.13...v0.2.14
[0.2.13]: https://github.com/osteele/gojekyll/compare/v0.2.12...v0.2.13
[0.2.12]: https://github.com/osteele/gojekyll/compare/v0.2.11...v0.2.12
[0.2.11]: https://github.com/osteele/gojekyll/compare/v0.2.10...v0.2.11
[0.2.10]: https://github.com/osteele/gojekyll/compare/v0.2.9...v0.2.10
[0.2.9]: https://github.com/osteele/gojekyll/compare/v0.2.8...v0.2.9
[0.2.8]: https://github.com/osteele/gojekyll/compare/v0.2.7...v0.2.8
[0.2.7]: https://github.com/osteele/gojekyll/compare/v0.2.6...v0.2.7
[0.2.6]: https://github.com/osteele/gojekyll/releases/tag/v0.2.6
