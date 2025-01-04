# go-fltk-diceware

A simple diceware password generator, using FLTK for extremely minimal memory usage (11-15MB, 45MB with extended word list loaded).

Features dark/light mode and portrait/landscape mode that is responsive.

## Screenshots

![Light mode landscape](./docs/light-landscape.png)

![Dark mode portrait](./docs/dark-portrait.png)

## Installation

Go to the [releases page](https://github.com/charles-m-knox/go-fltk-diceware/releases) and download the latest version there. Place it anywhere in your `$PATH` and you're good to go.

## Development setup

This repository makes use of `git lfs` for tracking its word dictionaries. Please ensure you have it working.

To build, run

```bash
make build-prod
```

To install to `~/.local/bin/`, run

```bash
make install
```

## Flatpak note

It is possible to build a flatpak distribution of this application, but I don't currently have time to deal with the extra overhead, so the only recommended installation method is by building it yourself or downloading a precompiled binary from the releases page (if available).
