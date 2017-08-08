# Benchmarks

`[go]jekyll build` on a late-2015 MacBook Pro, running current versions of everything as of 2017-07-09.

Disable the cache by setting the environment variable `GOJEKYLL_DISABLE_CACHE=1`.
Disable threading by setting `GOMAXPROCS=1`.

SASS conversion and Pygments (`{\% highlight \%}`) are cached.

## Jekyll Docs

This site contains only one SASS file.
It contains a few instances of `{\% highlight \%}`.
Each of these results in a call to Pygment. This dominates the un-cached times.

| Executable | Options                     | Time          |
|------------|-----------------------------|---------------|
| jekyll     |                             | 18.53s        |
| gojekyll   | single-threaded; cold cache | 3.14s ± 0.23s |
| gojekyll   | single-threaded; warm cache | 2.19s ± 0.03s |
| gojekyll   | multi-threaded; cold cache  | 1.19s ± 0.03s |
| gojekyll   | multi-threaded; warm cache  | 0.63s ± 0.03s |

## Software Design web site

This site makes heavy use of SASS.

| Executable | Options                     | Time          |
|------------|-----------------------------|---------------|
| jekyll     |                             | 8.07s         |
| gojekyll   | single-threaded; cold cache | 1.46s ± 0.21s |
| gojekyll   | single-threaded; warm cache | 0.60s ± 0.23s |
| gojekyll   | multi-threaded; cold cache  | 1.23s ± 0.10s |
| gojekyll   | multi-threaded; warm cache  | 0.35s ± 0.04s |
