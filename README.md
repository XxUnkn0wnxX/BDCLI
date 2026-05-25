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

### Global Options

```bash
bdcli --silent <command>   # Suppress non-error output
```

You can also set `BDCLI_SILENT=1` to silence output in automation.

### Install BetterDiscord

Install BetterDiscord to a specific Discord channel:

```bash
bdcli install --channel stable   # Install to Discord Stable
bdcli install --channel ptb      # Install to Discord PTB
bdcli install --channel canary   # Install to Discord Canary
```

Install BetterDiscord by providing a Discord install path:

```bash
bdcli install --path /path/to/Discord
```

### Uninstall BetterDiscord

Uninstall BetterDiscord from a specific Discord channel:

```bash
bdcli uninstall --channel stable   # Uninstall from Discord Stable
bdcli uninstall --channel ptb      # Uninstall from Discord PTB
bdcli uninstall --channel canary   # Uninstall from Discord Canary
```

Uninstall BetterDiscord by providing a Discord install path:

```bash
bdcli uninstall --path /path/to/Discord
```

Uninject BetterDiscord from all detected Discord installations (without deleting data):

```bash
bdcli uninstall --all
```

Fully uninstall BetterDiscord from all Discord installations and remove all BetterDiscord folders:

```bash
bdcli uninstall --full
```

### Check Version

```bash
bdcli version
```

### Update BetterDiscord

```bash
bdcli update
bdcli update --check
```

### Show BetterDiscord Info

```bash
bdcli info
```

### Discover Discord Installs

```bash
bdcli discover installs
bdcli discover paths
bdcli discover addons
```

### Manage Plugins

```bash
bdcli plugins list
bdcli plugins info <name>
bdcli plugins install <name|id|url>
bdcli plugins update <name|id|url>
bdcli plugins update <name|id> --check    # Check for updates without installing
bdcli plugins remove <name|id>
```

### Manage Themes

```bash
bdcli themes list
bdcli themes info <name>
bdcli themes install <name|id|url>
bdcli themes update <name|id|url>
bdcli themes update <name|id> --check     # Check for updates without installing
bdcli themes remove <name|id>
```

### Browse the Store

```bash
bdcli store search <query>
bdcli store show <id|name>

bdcli store plugins search <query>
bdcli store plugins show <id|name>

bdcli store themes search <query>
bdcli store themes show <id|name>
```

### Shell Completions

```bash
bdcli completion bash
bdcli completion zsh
bdcli completion fish
```

### Help

```bash
bdcli --help
bdcli [command] --help
```

### Automation

For scripts and CI jobs, you can suppress non-error output:

```bash
# One-off command
bdcli --silent install --channel stable

# Environment variable (applies to all commands)
BDCLI_SILENT=1 bdcli update
```

### CLI Help Output

<details>
<summary>Show full root help output</summary>

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

## Platform Notes

- Linux Snap Discord installs are not supported.
- Flatpak Discord installs are supported.

### Unsupported Configurations

- Linux Snap Discord installs are not supported due to upstream Snap packaging/runtime changes.

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
- [Task](https://taskfile.dev/) (optional, for task automation)
- [GoReleaser](https://goreleaser.com/) (for releases)

### Setup

Clone the repository and install dependencies:

```bash
git clone https://github.com/BetterDiscord/cli.git
cd cli
task setup  # Or: go mod download
```

### Available Tasks

Run `task --list-all` to see all available tasks:

```bash
# Development
task run             # Run the CLI (pass args with: task run -- install stable)

# Building
task build           # Build for current platform
task build:all       # Build for all platforms (GoReleaser)

# Testing
task test            # Run tests
task test:verbose    # Run tests with verbose output
task coverage        # Run tests with coverage summary
task coverage:html   # Generate HTML coverage report

# Code Quality
task fmt             # Format Go files
task vet             # Run go vet
task lint            # Run golangci-lint
task check           # Run fix, fmt, vet, lint, test

# Release
task release:snapshot # Test release build
task release          # Create release (requires tag)

# Cleaning
task clean           # Remove build and debug artifacts
```

### Running Locally

```bash
# Run directly
go run main.go install stable

# Or use Task
task run -- install stable
```

### Building

```bash
# Build for current platform
task build

# Build for all platforms
task build:all

# Output will be in ./dist/
```

### Testing

```bash
# Run all tests
task test

# Run with coverage
task coverage
```

### Releasing

1. Create and push a new tag:

   ```bash
   git tag -a v0.2.0 -m "Release v0.2.0"
   git push origin v0.2.0
   ```

2. GitHub Actions will automatically build and create a draft release

3. Edit the release notes and publish

4. Publish to npm:

   ```bash
   npm publish
   ```

## Project Structure

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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Links

- [BetterDiscord Website](https://betterdiscord.app/)
- [BetterDiscord Documentation](https://docs.betterdiscord.app/)
- [Issue Tracker](https://github.com/BetterDiscord/cli/issues)
- [npm Package](https://www.npmjs.com/package/@betterdiscord/cli)

## Acknowledgments

Built with:

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [GoReleaser](https://goreleaser.com/) - Release automation
- [Task](https://taskfile.dev/) - Task runner

---

Made with ❤️ by the BetterDiscord Team
