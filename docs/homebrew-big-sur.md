# Homebrew Setup For Big Sur

This fork includes example Homebrew formula files under:

- `brewfiles/go.rb`
- `brewfiles/bdcli.rb`

The intended use is:

- install the `bdcli.rb` formula from your local BDCLI tap
- install the bundled `go.rb` formula from that same tap as the Big Sur-compatible Go toolchain

This guide is for local taps on macOS. It does not modify Homebrew core and it does not require publishing a tap to GitHub.

## 1. Create A Local Tap

Create a local tap skeleton with Homebrew:

```bash
brew tap-new xxunkn0wnxx/bdcli
```

That creates a repository at:

```text
$(brew --repository)/Library/Taps/xxunkn0wnxx/homebrew-bdcli
```

Make sure the `Formula` directory exists:

```bash
mkdir -p "$(brew --repository)/Library/Taps/xxunkn0wnxx/homebrew-bdcli/Formula"
```

## 2. Copy The Formula Files

From this repo, copy the bundled formula files into the tap:

```bash
cp brewfiles/go.rb "$(brew --repository)/Library/Taps/xxunkn0wnxx/homebrew-bdcli/Formula/go.rb"
cp brewfiles/bdcli.rb "$(brew --repository)/Library/Taps/xxunkn0wnxx/homebrew-bdcli/Formula/bdcli.rb"
```

At that point your local tap should look like:

```text
$(brew --repository)/Library/Taps/xxunkn0wnxx/homebrew-bdcli
└── Formula/
    ├── bdcli.rb
    └── go.rb
```

## 3. Decide Which Go Dependency To Use

The bundled `bdcli.rb` currently uses:

```ruby
depends_on "xxunkn0wnxx/bdcli/go" => :build
```

That makes `bdcli` build against the `go.rb` formula in the same local tap.

If you want to build against standard Homebrew Go instead, edit the copied `Formula/bdcli.rb` and change that line to:

```ruby
depends_on "go" => :build
```

For the Big Sur-focused setup in this repo, the default is the tap-local `go.rb` formula.

## 4. Install The Big Sur Go Toolchain

For the default tap-local setup, install the bundled Go formula first:

```bash
brew install -sv xxunkn0wnxx/bdcli/go
```

If you changed `bdcli.rb` to use standard Homebrew Go instead, install that instead:

```bash
brew install -sv go
```

Verify the active Go toolchain:

```bash
go version
```

You should see a Go `1.24.8` build.

## 5. Install BDCLI

Install `bdcli` from the local tap:

```bash
brew install -sv xxunkn0wnxx/bdcli/bdcli
```

Where:

- `-s` builds from source
- `-v` prints verbose build output

Verify the install:

```bash
bdcli version
bdcli --help
```

## 6. Rebuild Or Upgrade

Because this fork is meant for local source builds and main-branch artifacts rather than version-tagged releases, the simplest update path is to refresh your repo copy and then reinstall from source.

If your formula points at the tap-local files, update those files first, then rebuild:

```bash
cp brewfiles/go.rb "$(brew --repository)/Library/Taps/xxunkn0wnxx/homebrew-bdcli/Formula/go.rb"
cp brewfiles/bdcli.rb "$(brew --repository)/Library/Taps/xxunkn0wnxx/homebrew-bdcli/Formula/bdcli.rb"
brew reinstall -sv xxunkn0wnxx/bdcli/bdcli
```

If you changed `bdcli.rb` to depend on plain `go`, you only need to reinstall `bdcli`:

```bash
brew reinstall -sv xxunkn0wnxx/bdcli/bdcli
```

## 7. Uninstall Or Remove

Remove `bdcli`:

```bash
brew uninstall bdcli
```

If you also installed the tap-local Go formula and no longer want it:

```bash
brew uninstall xxunkn0wnxx/bdcli/go
```

If you want to remove the entire local tap:

```bash
brew untap xxunkn0wnxx/bdcli
```

## Notes

- `brewfiles/go.rb` is the preserved Go formula in this repo that still works on Big Sur.
- `brewfiles/bdcli.rb` is an example formula file for this fork, not an automatically published tap.
- If you do not want to involve Homebrew at all, the simplest local build path is still:

```bash
./scripts/build-macos-big-sur.zsh
```
