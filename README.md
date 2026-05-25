<h1 align="center">BetterDiscord CLI</h1>

<p align="center">
  <a href="#overview">Overview</a> |
  <a href="#compatibility-matrix">Compatibility</a> |
  <a href="#installation">Installation</a> |
  <a href="#quick-start">Quick Start</a> |
  <a href="#usage">Usage</a> |
  <a href="#faq">FAQ</a> |
  <a href="#development">Development</a>
</p>


<p align="center">
   <img alt="Preview" width="743" height="322" src="https://i.imgur.com/LXUEJrW.png" />
   <br />
   A cross-platform command-line interface for installing, updating, and managing <a href="https://github.com/BetterDiscord/BetterDiscord">BetterDiscord</a>.
   <br />
   <br />
   <a href="https://betterdiscord.app/invite" target="_blank">
      <img src="https://img.shields.io/badge/discord-join-green?labelColor=0c0d10&color=3a71c1&style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHdpZHRoPSI0OCIgaGVpZ2h0PSI0OCIgdmlld0JveD0iMCAwIDQ4IDQ4IiBmaWxsPSJub25lIj4NCjxwYXRoIGQ9Ik0xNi41MzUzIDUuNDUwNTNDMzMuODIzMyAtMS40NjgwOSA1MC44ODE1IDE2Ljg4NjMgNDEuNTkyNSAzMy42MDY2QzM3LjM3MzIgNDEuMjAxMiAyNi44OTA0IDQ3LjMxNyAxNC42ODQyIDQxLjUzMjZMNi4xOTk5IDQzLjk1NjdDNC44ODM0NiA0NC4zMzI4IDMuNjY2OTMgNDMuMTIxMiA0LjAzNjM0IDQxLjgwMzdDNC41NDI3MiAzOS45OTc2IDUuNzQyNTcgMzUuNzM5OCA2LjQ0NDEgMzMuNDMyM0MxLjE4Mjc5IDI0LjA0NCA0LjczMDUgMTAuMTc0OCAxNi41MzUzIDUuNDUwNTNaTTE1Ljk5NTQgMjAuMjQ5NkMxNS45OTU0IDIwLjkzOTkgMTYuNTU1IDIxLjQ5OTYgMTcuMjQ1NCAyMS40OTk2SDMwLjc0OThDMzEuNDQwMSAyMS40OTk2IDMxLjk5OTggMjAuOTM5OSAzMS45OTk4IDIwLjI0OTZDMzEuOTk5OCAxOS41NTkyIDMxLjQ0MDEgMTguOTk5NiAzMC43NDk4IDE4Ljk5OTZIMTcuMjQ1NEMxNi41NTUgMTguOTk5NiAxNS45OTU0IDE5LjU1OTIgMTUuOTk1NCAyMC4yNDk2Wk0xNy4yNDk4IDI2LjQ3NDZDMTYuNTU5NCAyNi40NzQ2IDE1Ljk5OTggMjcuMDM0MiAxNS45OTk4IDI3LjcyNDZDMTUuOTk5OCAyOC40MTQ5IDE2LjU1OTQgMjguOTc0NiAxNy4yNDk4IDI4Ljk3NDZIMjYuNzQ5OEMyNy40NDAxIDI4Ljk3NDYgMjcuOTk5OCAyOC40MTQ5IDI3Ljk5OTggMjcuNzI0NkMyNy45OTk4IDI3LjAzNDIgMjcuNDQwMSAyNi40NzQ2IDI2Ljc0OTggMjYuNDc0NkgxNy4yNDk4WiIgZmlsbD0iIzNhNzFjMSIvPg0KPC9zdmc+" alt="Chat"/>
   </a>
   <a href="https://github.com/BetterDiscord/cli/releases/" target="_blank">
      <img src="https://img.shields.io/github/downloads/BetterDiscord/cli/total?labelColor=0c0d10&color=3a71c1&style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNDgiIGhlaWdodD0iNDgiIHZpZXdCb3g9IjAgMCA0OCA0OCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTEyLjI1IDM4LjVIMzUuNzVDMzYuNzE2NSAzOC41IDM3LjUgMzkuMjgzNSAzNy41IDQwLjI1QzM3LjUgNDEuMTY4MiAzNi43OTI5IDQxLjkyMTIgMzUuODkzNSA0MS45OTQyTDM1Ljc1IDQySDEyLjI1QzExLjI4MzUgNDIgMTAuNSA0MS4yMTY1IDEwLjUgNDAuMjVDMTAuNSAzOS4zMzE4IDExLjIwNzEgMzguNTc4OCAxMi4xMDY1IDM4LjUwNThMMTIuMjUgMzguNUgzNS43NUgxMi4yNVpNMjMuNjA2NSA2LjI1NThMMjMuNzUgNi4yNUMyNC42NjgyIDYuMjUgMjUuNDIxMiA2Ljk1NzExIDI1LjQ5NDIgNy44NTY0N0wyNS41IDhWMjkuMzMzTDMwLjI5MzEgMjQuNTQwN0MzMC45NzY1IDIzLjg1NzMgMzIuMDg0NiAyMy44NTczIDMyLjc2OCAyNC41NDA3QzMzLjQ1MTQgMjUuMjI0MiAzMy40NTE0IDI2LjMzMjIgMzIuNzY4IDI3LjAxNTZMMjQuOTg5OCAzNC43OTM4QzI0LjMwNjQgMzUuNDc3MiAyMy4xOTg0IDM1LjQ3NzIgMjIuNTE1IDM0Ljc5MzhMMTQuNzM2OCAyNy4wMTU2QzE0LjA1MzQgMjYuMzMyMiAxNC4wNTM0IDI1LjIyNDIgMTQuNzM2OCAyNC41NDA3QzE1LjQyMDIgMjMuODU3MyAxNi41MjgyIDIzLjg1NzMgMTcuMjExNyAyNC41NDA3TDIyIDI5LjMyOVY4QzIyIDcuMDgxODMgMjIuNzA3MSA2LjMyODgxIDIzLjYwNjUgNi4yNTU4TDIzLjc1IDYuMjVMMjMuNjA2NSA2LjI1NThaIiBmaWxsPSIjM2E3MWMxIi8+Cjwvc3ZnPgo=" alt="Downloads"/>
   </a>
   <a href="https://github.com/BetterDiscord/cli/blob/main/LICENSE" target="_blank">
      <img src="https://img.shields.io/github/license/BetterDiscord/cli?labelColor=0c0d10&color=3a71c1&style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTEwLjk2ODQgMi4zMjQ2NUMxMS41ODMgMS44NzYxNiAxMi40MTcgMS44NzYxNiAxMy4wMzE2IDIuMzI0NjVMMjAuNDUzNCA3Ljc0MDZDMjEuNDI5OSA4LjQ1MzE1IDIwLjkyNjggOS45OTgzNSAxOS43MTg5IDEwLjAwMDNINC4yODEwOEMzLjA3MzE4IDkuOTk4MzUgMi41NzAxMSA4LjQ1MzE1IDMuNTQ2NTcgNy43NDA2TDEwLjk2ODQgMi4zMjQ2NVpNMTMgNi4yNTAzNEMxMyA1LjY5ODA1IDEyLjU1MjMgNS4yNTAzNCAxMiA1LjI1MDM0QzExLjQ0NzcgNS4yNTAzNCAxMSA1LjY5ODA1IDExIDYuMjUwMzRDMTEgNi44MDI2MiAxMS40NDc3IDcuMjUwMzQgMTIgNy4yNTAzNEMxMi41NTIzIDcuMjUwMzQgMTMgNi44MDI2MiAxMyA2LjI1MDM0WiIgZmlsbD0iIzNhNzFjMSIvPgo8cGF0aCBkPSJNMTEuMjUgMTYuMDAwM0g5LjI1VjExLjAwMDNIMTEuMjVWMTYuMDAwM1oiIGZpbGw9IiMzYTcxYzEiLz4KPHBhdGggZD0iTTE0Ljc1IDE2LjAwMDNIMTIuNzVWMTEuMDAwM0gxNC43NVYxNi4wMDAzWiIgZmlsbD0iIzNhNzFjMSIvPgo8cGF0aCBkPSJNMTguNSAxNi4wMDAzSDE2LjI1VjExLjAwMDNIMTguNVYxNi4wMDAzWiIgZmlsbD0iIzNhNzFjMSIvPgo8cGF0aCBkPSJNMTguNzUgMTcuMDAwM0g1LjI1QzQuMDA3MzYgMTcuMDAwMyAzIDE4LjAwNzcgMyAxOS4yNTAzVjE5Ljc1MDNDMyAyMC4xNjQ1IDMuMzM1NzkgMjAuNTAwMyAzLjc1IDIwLjUwMDNIMjAuMjVDMjAuNjY0MiAyMC41MDAzIDIxIDIwLjE2NDUgMjEgMTkuNzUwM1YxOS4yNTAzQzIxIDE4LjAwNzcgMTkuOTkyNiAxNy4wMDAzIDE4Ljc1IDE3LjAwMDNaIiBmaWxsPSIjM2E3MWMxIi8+CjxwYXRoIGQ9Ik03Ljc1IDE2LjAwMDNINS41VjExLjAwMDNINy43NVYxNi4wMDAzWiIgZmlsbD0iIzNhNzFjMSIvPgo8L3N2Zz4K" alt="License"/>
   </a>
   <a href="https://www.npmjs.com/package/@betterdiscord/cli" target="_blank">
      <img src="https://img.shields.io/npm/v/@betterdiscord/cli?labelColor=0c0d10&color=3a71c1&style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPHBhdGggZD0iTTEwLjk2ODQgMi4zMjQ2NUMxMS41ODMgMS44NzYxNiAxMi40MTcgMS44NzYxNiAxMy4wMzE2IDIuMzI0NjVMMjAuNDUzNCA3Ljc0MDZDMjEuNDI5OSA4LjQ1MzE1IDIwLjkyNjggOS45OTgzNSAxOS43MTg5IDEwLjAwMDNINC4yODEwOEMzLjA3MzE4IDkuOTk4MzUgMi41NzAxMSA4LjQ1MzE1IDMuNTQ2NTcgNy43NDA2TDEwLjk2ODQgMi4zMjQ2NVpNMTMgNi4yNTAzNEMxMyA1LjY5ODA1IDEyLjU1MjMgNS4yNTAzNCAxMiA1LjI1MDM0QzExLjQ0NzcgNS4yNTAzNCAxMSA1LjY5ODA1IDExIDYuMjUwMzRDMTEgNi44MDI2MiAxMS40NDc3IDcuMjUwMzQgMTIgNy4yNTAzNEMxMi41NTIzIDcuMjUwMzQgMTMgNi44MDI2MiAxMyA2LjI1MDM0WiIgZmlsbD0iIzNhNzFjMSIvPgo8cGF0aCBkPSJNMTEuMjUgMTYuMDAwM0g5LjI1VjExLjAwMDNIMTEuMjVWMTYuMDAwM1oiIGZpbGw9IiMzYTcxYzEiLz4KPHBhdGggZD0iTTE0Ljc1IDE2LjAwMDNIMTIuNzVWMTEuMDAwM0gxNC43NVYxNi4wMDAzWiIgZmlsbD0iIzNhNzFjMSIvPgo8cGF0aCBkPSJNMTguNSAxNi4wMDAzSDE2LjI1VjExLjAwMDNIMTguNVYxNi4wMDAzWiIgZmlsbD0iIzNhNzFjMSIvPgo8cGF0aCBkPSJNMTguNzUgMTcuMDAwM0g1LjI1QzQuMDA3MzYgMTcuMDAwMyAzIDE4LjAwNzcgMyAxOS4yNTAzVjE5Ljc1MDNDMyAyMC4xNjQ1IDMuMzM1NzkgMjAuNTAwMyAzLjc1IDIwLjUwMDNIMjAuMjVDMjAuNjY0MiAyMC41MDAzIDIxIDIwLjE2NDUgMjEgMTkuNzUwM1YxOS4yNTAzQzIxIDE4LjAwNzcgMTkuOTkyNiAxNy4wMDAzIDE4Ljc1IDE3LjAwMDNaIiBmaWxsPSIjM2E3MWMxIi8+CjxwYXRoIGQ9Ik03Ljc1IDE2LjAwMDNINS41VjExLjAwMDNINy43NVYxNi4wMDAzWiIgZmlsbD0iIzNhNzFjMSIvPgo8L3N2Zz4K" alt="License"/>
   </a>
</p>


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

For information on contributing to this project, please see [CONTRIBUTING.md](/CONTRIBUTING.md).

<a href="https://github.com/betterdiscord/cli/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=betterdiscord/cli" />
</a>

---

Made with ❤️ by the BetterDiscord Team
