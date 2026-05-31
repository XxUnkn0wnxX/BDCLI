# Contributing to BDCLI

This fork tracks BetterDiscord CLI while keeping the build and release flow focused on Intel Macs that are limited to macOS 11 Big Sur.

Contributions should preserve that scope unless a maintainer explicitly decides otherwise.

## Development Scope

- Keep the fork buildable with the Go version pinned in `go.mod`.
- Keep release artifacts macOS Big Sur Intel focused.
- Preserve the `scripts/build-macos.zsh` local build path.
- Preserve the GitHub Actions workflow shape used by this fork.
- Avoid adding package-manager publishing flows unless they are explicitly for this fork.

Upstream fixes for shared CLI logic are welcome when they do not raise the Go toolchain requirement or replace this fork's macOS-focused release workflow.

## Setup

Prerequisites:

- Go matching `go.mod`.
- Git.
- Task, if you want to use `Taskfile.yml`.

Common commands:

```sh
go test ./...
task test
task build
./scripts/build-macos.zsh
```

## Pull Requests

Before opening a pull request:

1. Rebase or merge from the current `develop` branch.
2. Run `go test ./...`.
3. Confirm `go.mod`, workflow files, and release config still preserve the fork's Go pin and macOS Big Sur focus.
4. Keep changes scoped to the bug fix, upstream sync, or feature being proposed.

## Style

- Run `gofmt` on Go files you touch.
- Prefer small, direct changes over broad refactors.
- Add or update tests for path handling, install detection, and behavior that can regress across platforms.
- Keep commit messages descriptive and focused on the user-visible reason for the change.

## Releases

This fork publishes macOS Big Sur Intel artifacts through the existing fork workflow. Do not replace that with upstream's cross-platform release, npm, winget, or Homebrew-tap publishing setup without explicit maintainer approval.
