# Benchmarks

Disable the cache by setting the environment variable `GOJEKYLL_DISABLE_CACHE=1`.
Disable threading by setting `GOMAXPROCS=1`.

SASS conversion and Pygments (`{\% highlight \%}`) are cached.

## MadelineProto Docs

`[go]jekyll build` on an Intel Xeon E5620 @ 2.40GHz, running current versions of
everything as of 2022-01-29.

This site contains 1873 markdown files, and runs a modified version of the
complex [Just The Docs theme](https://pmarsceill.github.io/just-the-docs/), with
many SASS files, sitemap, search index generation.

| Executable | Options         | Time             |
| ---------- | --------------- | ---------------- |
| jekyll     |                 | Timeout @ 1 hour |
| gojekyll   | single-threaded | 750.61s          |
| gojekyll   | multi-threaded  | 142.16s          |

## Software Design web site

Site source: <https://github.com/sd17spring/sd17spring.github.io>
This site makes heavy use of SASS.

MacBook Pro (13", M1, 2020), macOS Monterey (12.2)
gojekyll v0.2.5
go1.17.6 darwin/arm64

| Executable | Options                         | Time               |
| ---------- | ------------------------------- | ------------------ |
| gojekyll   | single-threaded; cache disabled | 1.568 s ±  0.145 s |
| gojekyll   | single-threaded; warm cache     | 1.427 s ±  0.191 s |
| gojekyll   | multi-threaded; cache disabled  | 1.291 s ±  0.104 s |
| gojekyll   | multi-threaded; warm cache      | 1.118 s ±  0.110 s |

## Older Versions

## Jekyll Docs (gojekyll 0.2.5)

`[go]jekyll build` on a Late-2015 MacBook Pro, running current versions of
everything as of 2017-07-09.

This site contains only one SASS file.
It contains a few instances of `{\% highlight \%}`.
Each of these results in a call to Pygment. This dominates the un-cached times.

| Executable | Options         | Time   |
| ---------- | --------------- | ------ |
| jekyll     |                 | 9.086s |
| gojekyll   | single-threaded | 5.35s  |
| gojekyll   | multi-threaded  | 2.50s  |

### Software Design web site (gojekyll 0.2.5)

Site source: <https://github.com/sd17spring/sd17spring.github.io>

MacBook Pro (13", M1, 2020), macOS Monterey (12.2)
gojekyll v0.2.5
go1.17.6 darwin/arm64
Ruby 3.1.0, Jekyll 4.2.1

| Executable | Options                         | Time               |
| ---------- | ------------------------------- | ------------------ |
| jekyll     | [haven't been able to install]  |                    |
| gojekyll   | single-threaded; cache disabled | 1.417 s ±  0.140 s |
| gojekyll   | single-threaded; warm cache     | 1.297 s ±  0.145 s |
| gojekyll   | multi-threaded; cache disabled  | 1.262 s ±  0.201 s |
| gojekyll   | multi-threaded; warm cache      | 1.004 s ±  0.142 s |

MacBook Pro (15" Late-2015), running current versions of all software as of
2017-07-09
gojekyll v0.2.5
Ruby 2.4.1, Jekyll 3.4.3

| Executable | Options                     | Time          |
| ---------- | --------------------------- | ------------- |
| jekyll     |                             | 8.07s         |
| gojekyll   | single-threaded; cold cache | 1.46s ± 0.21s |
| gojekyll   | single-threaded; warm cache | 0.60s ± 0.23s |
| gojekyll   | multi-threaded; cold cache  | 1.23s ± 0.10s |
| gojekyll   | multi-threaded; warm cache  | 0.35s ± 0.04s |
