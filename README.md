# go-fltk-diceware

A simple diceware password generator, using FLTK for extremely minimal memory usage (11-15MB, 45MB with extended word list loaded).

Features dark/light mode and portrait/landscape mode that is responsive.

## Screenshots

![Light mode landscape](./docs/light-landscape.png)

![Dark mode portrait](./docs/dark-portrait.png)

## Installation

This application can be installed via Flatpak:

```bash
# if you do not have flathub added as a remote, please add it first, so that
# the necessary flatpak runtimes can be acquired:
flatpak --user remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo

flatpak --user remote-add --if-not-exists cmcode https://flatpak.cmcode.dev/cmcode.flatpakrepo

flatpak --user install cmcode dev.cmcode.go-fltk-diceware
```

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

## Building the flatpak

To build the flatpak, use:

```bash
make flatpak-build-test
```

This will install a local version of the flatpak, and you can run it via

```bash
flatpak --user run dev.cmcode.go-fltk-diceware
```

Once you're satisfied with it, you can then proceed to release it, assuming the remote repository's mount point is set up correctly:

```bash
# WARNING: This will update the globally available repository!
make flatpak-release
```
