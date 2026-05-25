# BetterDiscord CLI

[![Go Version](https://img.shields.io/github/go-mod/go-version/BetterDiscord/cli)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/BetterDiscord/cli)](https://github.com/BetterDiscord/cli/releases)
[![License](https://img.shields.io/github/license/BetterDiscord/cli)](LICENSE)
[![npm](https://img.shields.io/npm/v/@betterdiscord/cli)](https://www.npmjs.com/package/@betterdiscord/cli)

A cross-platform command-line interface for installing, updating, and managing [BetterDiscord](https://betterdiscord.app/).

[Overview](#overview) | [Compatibility](#compatibility-matrix) | [Installation](#installation) | [Quick Start](#quick-start) | [Usage](#usage) | [FAQ](#faq) | [Development](#development)

## Overview

This repository contains the source code for the BetterDiscord CLI. It is a native Go command-line tool for installing, removing, and maintaining BetterDiscord across supported Discord desktop installations.

## Features

- Easy installation and uninstallation of BetterDiscord
- Support for multiple Discord channels (Stable, PTB, Canary)
- Discover Discord installs and suggested paths
- Manage plugins and themes (list, install, update, remove)
- Browse and search the BetterDiscord store
- Cross-platform support (Windows, macOS, Linux)
- Available via npm for easy distribution
- Fast and lightweight Go binary

## Compatibility Matrix

| Platform | Minimum Version / Install Type | Support Status | Notes |
| --- | --- | --- | --- |
| Windows | Windows 11+ | ✅ | x64, ARM64, and x86 builds are available. |
| macOS | macOS 14+ | ✅ | x64 and ARM64 builds are available. |
| Linux | Native Discord install | ✅ | Standard package-manager installs are supported. |
| Linux | Flatpak Discord install | ✅ | Flatpak-based Discord installs are supported. |
| Linux | Snap Discord install | ❌ | Unsupported due to upstream Snap packaging/runtime changes. |

## Installation

### Via npm (Recommended)

```bash
npm install -g @betterdiscord/cli
```

### Via Go

```bash
go install github.com/betterdiscord/cli@latest
```

### Via winget (Windows)

```bash
winget install betterdiscord.cli
```

### Via Homebrew/Linuxbrew

```bash
brew install betterdiscord/tap/bdcli
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/BetterDiscord/cli/releases).

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

### Detailed Examples

<details>
<summary>Common Workflows</summary>

```bash
# Install / update / inspect
bdcli install --channel stable
bdcli update
bdcli info

# Uninstall from one channel
bdcli uninstall --channel stable

# Discover installed targets
bdcli discover installs

# Basic addon workflow
bdcli plugins list
bdcli themes list
```

</details>

<details>
<summary>Advanced Workflows</summary>

```bash
# Install or uninstall by explicit path
bdcli install --path /path/to/Discord
bdcli uninstall --path /path/to/Discord

# Global removal modes
bdcli uninstall --all
bdcli uninstall --full

# Check-only updates
bdcli update --check
bdcli plugins update <name|id> --check
bdcli themes update <name|id> --check

# Store browsing
bdcli store search <query>
bdcli store plugins show <id|name>
bdcli store themes show <id|name>

# Shell + automation
bdcli completion zsh
bdcli --silent install --channel stable
BDCLI_SILENT=1 bdcli update
```

</details>

### CLI Help Output

<details>
<summary>Show Full Root Help Output</summary>

```
A cross-platform CLI for installing, updating, and managing BetterDiscord.

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

### Does the CLI support Flatpak Discord on Linux?

Yes. Flatpak Discord installs are supported.

### Why is Snap Discord unsupported on Linux?

Discord Snap packaging/runtime changes prevent the CLI from supporting Snap installs.

### How can I use the global BetterDiscord folder with Flatpak?

1. Grant the Flatpak app access to the BetterDiscord config directory:

```sh
flatpak --user override com.discordapp.Discord --filesystem=xdg-config/BetterDiscord:rw
```

2. Symlink the BetterDiscord folder into the Flatpak app config:

```sh
ln -s "${XDG_CONFIG_HOME:-$HOME/.config}/BetterDiscord" "$HOME/.var/app/com.discordapp.Discord/config/BetterDiscord"
```

Replace `com.discordapp.Discord` with your Discord Flatpak app ID if it differs on your system.

## Development

### Prerequisites

- [Go](https://go.dev/) 1.26 or higher
- [Task](https://taskfile.dev/) (optional)
- [GoReleaser](https://goreleaser.com/) (optional, release builds)

### Quick Dev Loop

```bash
# Install dependencies
task setup

# Run locally
task run -- install stable

# Test + build
task test
task build
```

### Setup

Clone the repository and install dependencies:

```bash
git clone https://github.com/BetterDiscord/cli.git
cd cli
task setup  # Or: go mod download
```

<details>
<summary>Additional Tasks And Release Flow</summary>

```bash
# Task catalog
task --list-all

# Useful checks
task coverage
task coverage:html
task check

# Multi-platform builds
task build:all

# Release helpers
task release:snapshot
task release
```

Release outline:

1. Create and push a new tag.
2. GitHub Actions builds artifacts and creates a draft release.
3. Publish the GitHub release.
4. Publish the npm package.

</details>

## Project Structure

<details>
<summary>Show Project Structure</summary>

```py
.
├── cmd/                  # Cobra commands
│   ├── install.go       # Install command
│   ├── update.go        # Update command
│   ├── info.go          # Info command
│   ├── discover.go      # Discover command
│   ├── plugins.go       # Plugins commands
│   ├── themes.go        # Themes commands
│   ├── store.go         # Store commands
│   ├── uninstall.go     # Uninstall command
│   ├── version.go       # Version command
│   └── root.go          # Root command
├── internal/            # Internal packages
│   ├── betterdiscord/  # BetterDiscord installation logic
│   ├── discord/        # Discord path resolution and injection
│   ├── models/         # Data models
│   └── utils/          # Utility functions
├── main.go             # Entry point
├── Taskfile.yml        # Task automation
└── .goreleaser.yaml    # Release configuration
```

</details>

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

<a href="https://github.com/betterdiscord/cli/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=betterdiscord/cli" />
</a>

---

Made with ❤️ by the BetterDiscord Team
