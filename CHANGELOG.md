# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Table of Contents (TOC) Support** (#76, #62): Added Kramdown-style TOC generation with `{:toc}` and `{::toc}` markers, including support for Jekyll's `toc_levels` configuration and heading exclusion with `{:.no_toc}`. Thanks [@tekknolagi](https://github.com/tekknolagi) for requesting
- **Permalink Timezone Configuration** (#67): Added `permalink_timezone` configuration option to control timezone for permalink date generation
- **Markdown Attributes Support** (#85, #64): Added support for full Kramdown markdown attribute syntax (`markdown=1`, `markdown=0`, `markdown=block`, `markdown=span`) in HTML blocks

### Fixed

- **Jekyll Timezone Compatibility** (#63): Fixed permalink generation to honor standard Jekyll `timezone` config first before falling back to `permalink_timezone`, ensuring full Jekyll compatibility

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
- **Windows Routing**: Fixed URL routing on Windows by using forward slashes consistently in path handling
- **Nested Directory Watching**: Fixed file watcher to recursively watch nested directories and detect changes in subdirectories
- **Config Updates**: Fixed `Config.Set` to properly update YAML MapSlice so template changes are observed correctly

### Changed

- **Logging System** (#75, #35): Replaced scattered `fmt.Printf` statements with centralized logging package supporting proper log levels (Debug, Info, Warning, Error) and quiet mode
- **File Watcher**: Improved file watcher with automatic fallback to polling when directory count exceeds 500, preventing file descriptor exhaustion on large sites
- **Error Messages**: Enhanced markdown renderer error messages for common issues (e.g., suggesting `<br/>` instead of `<br>`)

### Maintenance

- **GitHub Actions** (#87): Updated CI workflows to test on Ubuntu, macOS, and Windows; updated actions to latest versions; added golangci-lint configuration
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
