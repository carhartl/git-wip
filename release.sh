#!/bin/sh

set -eu

git stash --include-untracked # make way for Goreleaser
git tag "v$1" -m "Release v$1"
git push --tags
goreleaser release --clean
git stash pop --index
