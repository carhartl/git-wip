# git-wip

[![CI](https://github.com/carhartl/git-wip/actions/workflows/ci.yml/badge.svg)](https://github.com/carhartl/git-wip/actions/workflows/ci.yml)

Command-line utility for listing your work in progress Git repositories.

## Installation

With [Go](https://golang.org/):

```bash
go install github.com/carhartl/git-wip
```

Via [Homebrew](https://brew.sh/):

```bash
brew install carhartl/tap/git-wip
```

## Usage

```bash
git-wip
```

By default uses the current directory as the starting point for the search. You can also specify a different directory:

```bash
git-wip /path/to/directory
```

Tip:

Due to the executable's naming you can also call it like so:

```bash
git wip
```

E.g. use with your favorite Git alias...
