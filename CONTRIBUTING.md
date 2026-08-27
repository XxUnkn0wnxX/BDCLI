# Contributing to BetterDiscord CLI

Thanks for taking the time to contribute!

The following is a set of guidelines for contributing to BetterDiscord CLI. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request. These guidelines have been adapted from [Atom](https://github.com/atom/atom/blob/master/CONTRIBUTING.md).

#### Table Of Contents

[Code of Conduct](#code-of-conduct)

[What should I know before I get started?](#what-should-i-know-before-i-get-started)
  * [Development Setup](#development-setup)

[How Can I Contribute?](#how-can-i-contribute)
  * [Reporting Bugs](#reporting-bugs)
  * [Suggesting Enhancements](#suggesting-enhancements)
  * [Your First Code Contribution](#your-first-code-contribution)
  * [Pull Requests](#pull-requests)

[AI-assisted contributions](#ai-assisted-contributions)

[Styleguides](#styleguides)
  * [Git Commit Messages](#git-commit-messages)
  * [Go Styleguide](#go-styleguide)

[Additional Notes](#additional-notes)
  * [Issue Labels](#issue-labels)

## Code of Conduct

This project and everyone participating in it is governed by the [Code of Conduct from the Contributor Covenant](https://www.contributor-covenant.org/version/1/4/code-of-conduct.html). By participating, you are expected to uphold this code. Please report unacceptable behavior.

## What should I know before I get started?

BetterDiscord CLI is a cross-platform command-line tool for installing and managing BetterDiscord across Windows, macOS, and Linux (including Flatpak support).

The repository is organized into:

```
.
├── cmd                    // CLI commands (install, themes, plugins, etc.)
├── internal/
│   ├── betterdiscord/     // BetterDiscord core logic and installation
│   ├── discord/           // Discord detection and injection
│   ├── models/            // Data structures
│   ├── output/            // Output formatting
│   ├── utils/             // Shared utilities
│   └── wsl/               // WSL support
├── completions/           // Generated shell completions (bash, fish, zsh; gitignored)
├── main.go                // CLI entry point
├── go.mod / go.sum        // Dependency management
└── README.md              // Project documentation
```

### Development Setup

Prerequisites:

* [Go](https://go.dev/) (matches `go.mod`)
* [Git](https://git-scm.com)

Common commands:

* `go run main.go` - run the CLI locally.
* `go test ./...` - run all tests.
* `gofmt -w .` - format Go files you touch.
* `go build -o bdcli .` - build a local binary.

### Building & Releases

The project uses [Taskfile](https://taskfile.dev/) for common tasks:

* `task run -- <args>` - run the CLI locally with dev version info.
* `task build` - build a binary for your current platform into `dist/`.
* `task build:all` - build binaries for all platforms (requires GoReleaser).
* `task test` - run all tests (use `task coverage` for coverage).
* `task check` - format, vet, lint, and test in one shot.

See `Taskfile.yml` for the full list of available tasks.

## How Can I Contribute?

### Reporting Bugs

Please search for existing issues first. If you find a match, add any extra context you have (logs, OS details, platform/distro info).

#### Before Submitting A Bug Report

* Reproduce on the latest `main` branch or a recent release.
* Collect output and error messages from the CLI.
* Note your OS, architecture (x64/arm64), and relevant platform details (e.g., Snap/Flatpak/native on Linux).

#### How Do I Submit A (Good) Bug Report?

Include clear steps to reproduce, expected vs. actual behavior, OS/platform info, and any error output shown by the CLI.

### Suggesting Enhancements

Open an issue describing the problem you are trying to solve, the proposed solution, and any alternatives considered.

#### Before Submitting An Enhancement Suggestion

Confirm the change aligns with the CLI's scope and doesn't require new platform support without a clear implementation plan.

#### How Do I Submit A (Good) Enhancement Suggestion?

Provide context, use cases, and a minimal spec of expected behavior. If the feature involves a new command or flag, include the proposed usage pattern.

### Your First Code Contribution

Unsure where to begin contributing? You can start by looking through `help-wanted` issues or any issues labelled `good-first-issue`.

### Pull Requests

Please follow these steps to have your contribution considered by the maintainers:

1. Use a pull request template, if one exists.
2. Follow the [styleguides](#styleguides)
3. After you submit your pull request, verify that all [status checks](https://help.github.com/articles/about-status-checks/) are passing <details><summary>What if the status checks are failing?</summary>If a status check is failing, and you believe that the failure is unrelated to your change, please leave a comment on the pull request explaining why you believe the failure is unrelated. A maintainer will re-run the status check for you. If we conclude that the failure was a false positive, then we will open an issue to track that problem with our check suite.</details>

While the prerequisites above must be satisfied prior to having your pull request reviewed, the reviewer(s) may ask you to complete additional design work, tests, or other changes before your pull request can be ultimately accepted.

## AI-assisted contributions

AI tools are fine to use. Submitting output with little or no personal review is not.

- **Review the diff yourself.** If you can't explain why each changed line is correct, the PR isn't ready.
- **Run the required checks locally.** Don't rely on CI to catch problems you could catch before pushing.
- **Disclose AI involvement in your PR description.** If an AI wrote a meaningful portion of the code, say so and briefly describe what you reviewed. This isn't meant as a penalty, it's useful context for the reviewer.
- **If you used an autonomous agent with minimal personal review of the output, say that explicitly.** PRs that appear to be unreviewed agent output and don't disclose this will be closed without detailed feedback.

## Styleguides

### Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line
* When only changing documentation, include `[ci skip]` in the commit title

### Go Styleguide

* Run `gofmt` on any Go files you touch.
* Keep standard library imports first, then a blank line, then all remaining imports (external and internal packages share one alphabetized block).
* Prefer early returns for error handling.
* Route all user-facing output through `internal/output` (never `fmt.Print*` to stdout directly) so `--silent` and output-capturing tests keep working.
* Use clear, descriptive names for functions and variables.
* Add comments to exported functions and types.
* Write tests for new functionality (see existing `*_test.go` files for patterns). Tests must not hit the real network, point package-level endpoint vars at `httptest` servers instead.

## Additional Notes

### Releases

Releases are tag-driven. Create a tag like `vX.Y.Z` and push it to trigger the release workflow. The workflow will:

* Build binaries for Windows, macOS, and Linux.
* Publish to package managers (npm, Homebrew, winget, etc.).
* Generate release notes on GitHub.

Tag format: `vX.Y.Z` (e.g., `v2.3.0`)

If you build locally for a release, use:

```sh
go build -ldflags="-X main.version=vX.Y.Z" -o bdcli
```

### Issue Labels

Use labels to clarify impact (`bug`, `enhancement`, `maintenance`) and scope (`cli`, `linux`, `macos`, `windows`, `flatpak`, `wsl`) or type (`documentation`, `testing`).

#### Type of Issue and Issue State

If you can, include an expected fix timeline or milestone to help triage.
