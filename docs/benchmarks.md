# Benchmarks

`[go]jekyll build` on an Intel Xeon E5620 @ 2.40GHz, running current versions of everything as of 2022-01-29.

Disable the cache by setting the environment variable `GOJEKYLL_DISABLE_CACHE=1`.
Disable threading by setting `GOMAXPROCS=1`.

SASS conversion and Pygments (`{\% highlight \%}`) are cached.

## Jekyll Docs

This site contains only one SASS file.
It contains a few instances of `{\% highlight \%}`.
Each of these results in a call to Pygment. This dominates the un-cached times.

| Executable | Options         | Time   |
|------------|-----------------|--------|
| jekyll     |                 | 9.086s |
| gojekyll   | single-threaded | 5.35s  |
| gojekyll   | multi-threaded  | 2.50s  |


## MadelineProto Docs

This site contains 1873 markdown files, and runs a modified version of the complex [Just The Docs theme](https://pmarsceill.github.io/just-the-docs/), with many SASS files, sitemap, search index generation.

| Executable | Options         | Time             |
|------------|-----------------|------------------|
| jekyll     |                 | Timeout @ 1 hour |
| gojekyll   | single-threaded | 750.61s          |
| gojekyll   | multi-threaded  | 142.16s          |