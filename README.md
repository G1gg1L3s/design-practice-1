[![Build Status](https://travis-ci.org/G1gg1L3s/design-practice-1.svg?branch=master)](https://travis-ci.org/G1gg1L3s/design-practice-1)
# Bood

This is an educational build system developed on top of the google/blueprint.

## Config

The main config is the `build.bood`. You can find one in the `build` or `example` directories.
It has two modules with the following keys
- `go_binary` - building and testing packages:
  - `name` - Module (and output) name.
  - `pkg` - Package to build.
  - `testPkg` - package to test.
  - `srcs` - sources to track.
  - `srcsExclude` - sources to exclude.
  - `testSrcs` - test sources to track.
  - `testSrcsExclude` - test sources to exclude.
- `go_doc` - generating docs for packages.
  - `name` - Module name.
  - `pkg` - package for which documentation will be generated.
  - `srcs` - sources to track.

## Building
To build all packages you first need to build the bood itself. Execute:

```
$ cd build
$ build ./cmd/bood
```

And bood executable will be places to your working directory. And now, you can boostrap bood using itself:
`$ ./bood`

If everyting is okay you should see something like this:
```
INFO 2021/03/10 11:19:44 Ninja build file is generated at out/build.ninja
INFO 2021/03/10 11:19:44 Starting the build now
[3/3] Testing bood
```
At this point, your bood executable will be at `out/bin/bood` and documentation will be at `out/docs`.


