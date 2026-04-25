# BetterDiscord CLI

[![Go Version](https://img.shields.io/badge/Go-1.24.x-00ADD8?logo=go&logoColor=white)](https://github.com/XxUnkn0wnxX/BDCLI/blob/develop/go.mod)
[![Main Artifact](https://img.shields.io/github/actions/workflow/status/XxUnkn0wnxX/BDCLI/release.yml?branch=main&label=Main%20Artifact)](https://github.com/XxUnkn0wnxX/BDCLI/actions/workflows/release.yml)
[![License](https://img.shields.io/github/license/XxUnkn0wnxX/BDCLI)](LICENSE)

A fork of BetterDiscord CLI focused on keeping the tool usable on older Intel Macs that are limited to macOS 11 Big Sur.

This fork is for older Intel Macs that are limited to macOS 11 Big Sur. It does not target Windows or Linux.

## Features

- 🚀 Easy installation and uninstallation of BetterDiscord
- 🔄 Support for multiple Discord channels (Stable, PTB, Canary)
- 🧭 Discover Discord installs and suggested paths
- 🧩 Manage plugins and themes (list, install, update, remove)
- 🛒 Browse and search the BetterDiscord store
- 🍎 Focused on macOS Big Sur-compatible Intel Mac builds
- ⚡ Fast and lightweight Go binary

## Installation

### Install With Homebrew

The recommended install path is the `xxunkn0wnxx/tap` Homebrew tap:

```bash
brew tap xxunkn0wnxx/tap
brew install --HEAD -sv xxunkn0wnxx/tap/bdcli
```

The tap formula is HEAD-only. It builds `bdcli` from this fork's `main` branch, builds from source, and uses the tap-local Big Sur-compatible Go formula as its build dependency.

To update or remove the Homebrew install:

```bash
brew reinstall --HEAD -sv xxunkn0wnxx/tap/bdcli
brew uninstall xxunkn0wnxx/tap/bdcli
```

### Clone And Build Locally

```bash
git clone https://github.com/XxUnkn0wnxX/BDCLI.git
cd BDCLI
./scripts/build-macos.zsh
```

That script builds the macOS 11 Big Sur Intel binary and a commit-named archive under `dist/`.

### Download Release

Download the latest Big Sur Intel binary from the GitHub Releases page, or use the workflow artifact from the matching `main` run if you want the Actions artifact copy.

GitHub Actions for this fork build macOS Big Sur Intel binaries only.
Pushes to `develop` run checks only.
Pushes to `main` run checks, upload the macOS binary as a workflow artifact, and replace the single rolling GitHub release with the latest macOS binary. GitHub also provides the source code archives for that release tag.

## Usage

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

## Supported Platforms

- **macOS 11 Big Sur** (x64 / Intel)

## Development

### Prerequisites

- [Go](https://go.dev/) 1.24.x
- [Task](https://taskfile.dev/) (optional, for task automation)
- [GoReleaser](https://goreleaser.com/) (for releases)

### Setup

Clone the repository and install dependencies:

```bash
git clone https://github.com/XxUnkn0wnxX/BDCLI.git
cd BDCLI
task setup  # Or: go mod download
```

### Available Tasks

Run `task --list-all` to see all available tasks:

```bash
# Development
task run             # Run the CLI (pass args with: task run -- install stable)

# Building
task build           # Build for current platform
task build:all       # Build macOS Big Sur release artifacts (GoReleaser)

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
task release          # Build a local Big Sur artifact bundle

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

# Build macOS Big Sur release artifacts
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

### Main Branch Artifacts

1. Push code changes to `main`
2. GitHub Actions will run checks, upload the macOS binary as a workflow artifact, and replace the rolling GitHub release
3. Download the macOS binary from the GitHub Releases page or from the workflow artifact

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
- [Issue Tracker](https://github.com/XxUnkn0wnxX/BDCLI/issues)

## Acknowledgments

Built with:

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [GoReleaser](https://goreleaser.com/) - Release automation
- [Task](https://taskfile.dev/) - Task runner

---

Originally based on BetterDiscord CLI.
