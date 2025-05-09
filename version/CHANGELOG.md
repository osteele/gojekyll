# Release Notes
<!-- markdownlint-disable MD024 -->

## 0.3 (unreleased)

### Contributions

* [@danog (Daniil Gentili)](https://github.com/danog): Add native chroma
  highlighting (~30x performance increase) [PR
  #42](https://github.com/osteele/gojekyll/pull/42)

* [@mevdschee (Maurits van der Schee)](https://github.com/mevdschee): add
  support for go mod and blackfriday v2 [PR
  #39](https://github.com/osteele/gojekyll/pull/39)

* [@bep (Bjørn Erik Pedersen)](https://github.com/bep): Update Hugo matrix in
  README [PR #38](https://github.com/osteele/gojekyll/pull/38)

### Other changes

* Upgraded to liquid v1.3.0. See the liquid release notes
  [here](https://github.com/osteele/liquid/blob/main/CHANGELOG.md#130-2020-02-13).

* The binary is built with go v1.17. (The previous release, v0.2.5, was built
  with go v1.08.)

## 0.2.5 (Aug 18, 2017)

### New
* Added support for `jekyll.environment` variable from `JEKYLL_ENV`
* Added script to compare jekyll and gojekyll builds

### Improved
* Renamed internal pipeline to renderers for better code organization
* Updated documentation about auto-generated ID compatibility

## 0.2.4 (Aug 10, 2017)

### New
* Added support for page excerpts
* Added support for rendering non-collection pages

### Improved
* Implemented `page.previous` and `page.next` navigation
* Changed Page content handling from []byte to string
* Reorganized rendering order - posts collection now renders last
* Improved Travis CI configuration to skip tests and lint on tags
* Moved example directory and updated documentation

## 0.2.3 (Aug 8, 2017)

### New
* Added support for dots in destination directory paths
* Added minification for SEO tags

### Improved
* Improved reload functionality
* Enhanced variables subcommand to handle byte-to-string conversion

## 0.2.2 (Aug 3, 2017)

### Improved
* Improved file change detection accuracy
* Enhanced file exclusion logic

### Fixed
* Fixed file descriptor leak
* Fixed race condition in file handling

## 0.2.1 (Jul 26, 2017)

### New
* Added CSS minification using minify instead of cssmin
* Added Gemfile to example site for side-by-side comparison

### Improved
* Improved livereloader output clarity
* Enhanced GitHub metadata plugin with environment variables support
* Better error reporting for failed plugins
* Implemented HTML minimization for feeds and sitemaps

## 0.2.0 (Jul 25, 2017)

### Major Changes
* Implemented incremental build support
* Created PageEmbed functionality
* Major refactoring for improved performance and maintainability

## 0.1.1 (Jul 19, 2017)

### Improved
* Reorganized command structure
* Updated build and deployment configuration
* Improved package organization

## 0.1.0 (Jul 17, 2017)

### Major Changes
* Initial stable release
* Complete implementation of core Jekyll functionality

## 0.0.5 (Jul 13, 2017)

### New
* Added mutex for cache operations

### Improved
* Moved commands to cmd package
* Updated package organization

### Fixed
* Fixed race conditions

## 0.0.4 (Jul 13, 2017)

### New
* Implemented support for data files (CSV and JSON)
* Added server watch functionality
* Added support for blackfriday auto header IDs

### Improved
* Improved error display in page

## 0.0.3 (Jul 12, 2017)

### Improved
* Reorganized main.go to repository root
* Updated goreleaser configuration

### Fixed
* Fixed static file serving

## 0.0.2 (Jul 12, 2017)

### New
* Initial setup of goreleaser

### Improved
* Configuration updates for Travis CI
* Basic build system setup

## 0.0.1 (Jul 12, 2017)

### New
* Initial release
* Basic Jekyll functionality implemented

### Improved
* Setup of Travis CI and basic tooling
