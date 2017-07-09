# Benchmarks

`[go]jekyll build` on a late-2015 MacBook Pro, running current versions of everything as of 2017-07-09.

Disable the cache by setting the environment variable `GOJEKYLL_DISABLE_CACHE=1`.
Disable threading by setting `GOMAXPROCS=1`.

SASS conversion and Pygments (`highlight`) are cached.

## Jekyll Docs

This site contains only one SASS file.
It contains a few instances of {% highlight %}.
This causes calls to Pygment, which dominate the un-cached times.

| Executable | Options                     | Time          |
|------------|-----------------------------|---------------|
| jekyll     |                             | 18.53s        |
| gojekyll   | single-threaded; cold cache | 2.96s ± 0.09s |
| gojekyll   | single-threaded; warm cache | 2.51s ± 0.04s |
| gojekyll   | multi-threaded; cold cache  | 1.37s ± 0.03s |
| gojekyll   | multi-threaded; warm cache  | 0.80s ± 0.06s |

## Software Design web site

This site makes heavy use of SASS.

| Executable | Options                     | Time          |
|------------|-----------------------------|---------------|
| jekyll     |                             | 8.07s         |
| gojekyll   | single-threaded; cold cache | 1.28s ± 0.23s |
| gojekyll   | single-threaded; warm cache | 0.69s ± 0.07s |
| gojekyll   | multi-threaded; cold cache  | 1.21s ± 0.12s |
| gojekyll   | multi-threaded; warm cache  | 0.40s ± 0.07s |
