<h1 align="center">BetterDiscord CLI</h1>

<p align="center">
  <a href="#overview">Overview</a> |
  <a href="#compatibility">Compatibility</a> |
  <a href="#installation">Installation</a> |
  <a href="#quick-start">Quick Start</a> |
  <a href="#usage">Usage</a> |
  <a href="#faq">FAQ</a> |
  <a href="#development">Development</a>
</p>

<p align="center">
  <a href="https://github.com/XxUnkn0wnxX/BDCLI/blob/develop/go.mod">
    <img src="https://img.shields.io/badge/Go-1.24.x-00ADD8?logo=go&logoColor=white" alt="Go Version">
  </a>
  <a href="https://github.com/XxUnkn0wnxX/BDCLI/actions/workflows/release.yml">
    <img src="https://img.shields.io/github/actions/workflow/status/XxUnkn0wnxX/BDCLI/release.yml?branch=main&label=Main%20Artifact" alt="Main Artifact">
  </a>
  <a href="LICENSE">
    <img src="https://img.shields.io/github/license/XxUnkn0wnxX/BDCLI" alt="License">
  </a>
</p>

## Overview

This is a fork of BetterDiscord CLI focused on keeping the tool usable on older Intel Macs that are limited to macOS 11 Big Sur.

It installs, removes, updates, and inspects BetterDiscord for supported macOS Discord desktop installations. The README, workflows, and release assets in this fork document the macOS Big Sur Intel path only.

## Features

- Easy installation and uninstallation of BetterDiscord.
- Support for Discord Stable, PTB, and Canary.
- Discovery for Discord installs, suggested paths, and BetterDiscord addons.
- Plugin and theme management commands.
- BetterDiscord store browsing and lookup commands.
- macOS 11 Big Sur-compatible Intel Mac release artifacts.
- Local source builds with the fork's pinned Go 1.24.x toolchain.

## Compatibility

| Platform | Minimum Version | Architecture | Support Status | Notes |
| --- | --- | --- | --- | --- |
| macOS | macOS 11 Big Sur | Intel | Supported | This is the only release target for this fork. |

## Installation

### Homebrew

The recommended install path is the `xxunkn0wnxx/tap` Homebrew tap:

```bash
brew tap xxunkn0wnxx/tap
brew install --HEAD -sv xxunkn0wnxx/tap/bdcli
```

The tap formula is HEAD-only. It builds `bdcli` from this fork's `main` branch and uses the tap-local Big Sur-compatible Go formula as its build dependency.

To update or remove the Homebrew install:

```bash
brew reinstall --HEAD -sv xxunkn0wnxx/tap/bdcli
brew uninstall xxunkn0wnxx/tap/bdcli
```

### Build Locally

Use the fork's local macOS build helper, [scripts/build-macos.zsh](scripts/build-macos.zsh), from the repository root:

```bash
git clone https://github.com/XxUnkn0wnxX/BDCLI.git
cd BDCLI
./scripts/build-macos.zsh
```

The script resolves the repository root, uses the available Go toolchain from `PATH`, generates shell completions, and writes the macOS 11 Big Sur Intel binary plus a commit-named archive under `dist/`.

### Download Release

Download the latest Big Sur Intel binary from this fork's GitHub Releases page, or use the workflow artifact from the matching `main` run if you want the Actions artifact copy.

GitHub Actions for this fork build macOS Big Sur Intel binaries only. Pushes to `develop` run checks only. Pushes to `main` run checks, upload the macOS binary as a workflow artifact, and replace the single rolling GitHub release with the latest macOS binary.

## Quick Start

```bash
# Install BetterDiscord to Discord Stable
bdcli install --channel stable

# Check installation details
bdcli info

# Update BetterDiscord
bdcli update

# Uninstall from a specific channel
bdcli uninstall --channel stable
```

## Usage

Use `bdcli [command] --help` for command-specific flags and examples.

### Command Reference

| Command | Purpose | Example |
| --- | --- | --- |
| `install` | Install BetterDiscord to Discord | `bdcli install --channel stable` |
| `uninstall` | Remove BetterDiscord from Discord | `bdcli uninstall --channel stable` |
| `update` | Update BetterDiscord to latest | `bdcli update` |
| `info` | Show installation state and metadata | `bdcli info` |
| `discover` | Discover installs, paths, and addons | `bdcli discover installs` |
| `plugins` | Manage BetterDiscord plugins | `bdcli plugins list` |
| `themes` | Manage BetterDiscord themes | `bdcli themes list` |
| `store` | Search and inspect store entries | `bdcli store search <query>` |
| `completion` | Generate shell completion scripts | `bdcli completion zsh` |
| `version` | Print CLI version | `bdcli version` |

### Common Workflows

<details>
<summary>Install, update, inspect</summary>

```bash
bdcli install --channel stable
bdcli install --channel ptb
bdcli install --channel canary

bdcli update
bdcli update --check
bdcli info
bdcli version
```

</details>

<details>
<summary>Explicit paths and uninstall modes</summary>

```bash
bdcli install --path /Applications/Discord.app
bdcli uninstall --path /Applications/Discord.app

bdcli uninstall --all
bdcli uninstall --full
```

</details>

<details>
<summary>Plugins, themes, store, and automation</summary>

```bash
bdcli discover installs
bdcli discover paths
bdcli discover addons

bdcli plugins list
bdcli plugins info <name>
bdcli plugins install <name|id|url>
bdcli plugins update <name|id|url>
bdcli plugins update <name|id> --check
bdcli plugins remove <name|id>

bdcli themes list
bdcli themes info <name>
bdcli themes install <name|id|url>
bdcli themes update <name|id|url>
bdcli themes update <name|id> --check
bdcli themes remove <name|id>

bdcli store search <query>
bdcli store plugins search <query>
bdcli store plugins show <id|name>
bdcli store themes search <query>
bdcli store themes show <id|name>

bdcli completion zsh
bdcli --silent install --channel stable
BDCLI_SILENT=1 bdcli update
```

</details>

### CLI Help Output

<details>
<summary>Show root help output</summary>

```text
A command-line interface for installing, updating, and managing BetterDiscord.

Usage:
   bdcli [flags]
   bdcli [command]

Available Commands:
   completion  Generate shell completions
   discover    Discover Discord installations and related data
   help        Help about any command
   info        Displays information about BetterDiscord installation
   install     Installs BetterDiscord to your Discord
   plugins     Manage BetterDiscord plugins
   store       Browse and search the BetterDiscord store
   themes      Manage BetterDiscord themes
   uninstall   Uninstalls BetterDiscord from your Discord
   update      Update BetterDiscord to the latest version
   version     Print the version number

Flags:
       --silent   Suppress non-error output
   -h, --help     help for bdcli

Use "bdcli [command] --help" for more information about a command.
```

</details>

## FAQ

### Is this the upstream BetterDiscord CLI?

No. This fork tracks useful upstream CLI fixes while preserving a macOS 11 Big Sur Intel build and release flow.

### Why is Go pinned to 1.24.x?

The fork uses a lower Go toolchain to keep the local Big Sur build path working on older Intel Macs.

### Where do release binaries come from?

The `main` workflow builds the macOS 11 Intel binary, uploads it as a workflow artifact, and replaces the rolling GitHub release asset.

## Development

### Prerequisites

- [Go](https://go.dev/) 1.24.x.
- [Task](https://taskfile.dev/) for task automation.
- [GoReleaser](https://goreleaser.com/) for local release artifact checks.

### Setup

```bash
git clone https://github.com/XxUnkn0wnxX/BDCLI.git
cd BDCLI
task setup  # Or: go mod download
```

### Quick Dev Loop

```bash
# Run locally
go run main.go install stable
task run -- install stable

# Test and build
task test
task build
```

<details>
<summary>Tasks and release flow</summary>

```bash
# Task catalog
task --list-all

# Testing
task test
task test:verbose
task coverage
task coverage:html

# Code quality
task fmt
task vet
task lint
task check

# Building
task build
task build:all

# Release helpers
task release:snapshot
task release

# Cleaning
task clean
```

Release outline:

1. Develop and verify changes on `develop`.
2. Promote `develop` to `main` with a fast-forward merge.
3. Push `main` so GitHub Actions builds the macOS Big Sur Intel artifact.
4. Return to `develop` and fast-forward it from `main`.

</details>

## Project Structure

<details>
<summary>Show project structure</summary>

```py
.
├── cmd/                  # Cobra commands
│   ├── discover.go      # Discover command
│   ├── info.go          # Info command
│   ├── install.go       # Install command
│   ├── plugins.go       # Plugin commands
│   ├── root.go          # Root command
│   ├── store.go         # Store commands
│   ├── themes.go        # Theme commands
│   ├── uninstall.go     # Uninstall command
│   ├── update.go        # Update command
│   └── version.go       # Version command
├── internal/            # Internal packages
│   ├── betterdiscord/  # BetterDiscord installation logic
│   ├── discord/        # Discord path resolution and injection
│   ├── models/         # Data models
│   ├── output/         # Output formatting
│   └── utils/          # Utility functions
├── scripts/             # Fork build and completion scripts
├── main.go              # Entry point
├── Taskfile.yml         # Task automation
└── .goreleaser.yaml     # macOS Big Sur release configuration
```

</details>

## Contributing

Contributions should preserve this fork's macOS Big Sur Intel focus. See [CONTRIBUTING.md](CONTRIBUTING.md) before changing Go versions, workflow files, release config, or build scripts.

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE) for details.

## Links

- [BetterDiscord Website](https://betterdiscord.app/)
- [BetterDiscord Documentation](https://docs.betterdiscord.app/)
- [Issue Tracker](https://github.com/XxUnkn0wnxX/BDCLI/issues)
- [Upstream BetterDiscord CLI](https://github.com/BetterDiscord/cli)

## Acknowledgments

Built with:

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [GoReleaser](https://goreleaser.com/) - release automation
- [Task](https://taskfile.dev/) - task runner

---

Originally based on BetterDiscord CLI.
