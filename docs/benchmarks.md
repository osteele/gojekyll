# Benchmarks

This document lists the times to run `jekyll build` and `gojekyll build` on a
few different sites.

See the Benchmarks section of CONTRIBUTING.md for instructions on how to
run a benchmark.

Notes:

- SASS conversion is cached
- Pygments (for version gojkeyll <= 0.2.5) is cached.

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

MacBook Pro (13", M1, 2020), macOS Monterey (12.2)
gojekyll v0.2.5
go1.17.6 darwin/arm64

Notes:

- This site makes heavy use of SASS.
- This site uses The GitHub metadata plugin. The results of that plugin are not
  cached. (See issue [#43](https://github.com/osteele/gojekyll/issues/43).)
  These benchmarks are run from a machine with a high latency to GitHub. I
  suspect that this latency dominates these benchmark time. Previous benchmarks
  were from the U.S.

| Executable | Options                         | Time               |
| ---------- | ------------------------------- | ------------------ |
| gojekyll   | single-threaded; cache disabled | 1.568 s ±  0.145 s |
| gojekyll   | single-threaded; warm cache     | 1.427 s ±  0.191 s |
| gojekyll   | multi-threaded; cache disabled  | 1.291 s ±  0.104 s |
| gojekyll   | multi-threaded; warm cache      | 1.118 s ±  0.110 s |
