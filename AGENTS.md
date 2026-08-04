# AGENT OPERATING NOTES

These notes are curated for downstream agents. Treat them as binding guidance: each section funnels a different expectation for builds, style, testing, and compliance with upstream tooling. (`CLAUDE.md` is a symlink to this file.)

## 1. Project layout recap

A pure-Go Cobra CLI (`bdcli`) for installing, updating, and managing BetterDiscord. It is the command-line sibling of the BetterDiscord GUI installer — the two repos share domain concepts (release channels, Discord discovery, asar download, injection) but no code; do not assume installer file paths or APIs exist here.

- `cmd/` — Cobra commands, one file per command (`install`, `uninstall`, `update`, `info`, `discover`, `plugins`, `themes`, `store`, `version`, `completion`). `root.go` wires the global `--silent` flag, the `BDCLI_SILENT` env fallback, and version info.
- `internal/betterdiscord/` — BetterDiscord core: data-folder setup, asar download (website → GitHub fallback), install/repair, addon management, and the store client.
- `internal/discord/` — Discord install discovery and injection. Per-OS path logic lives in `paths_windows.go` / `paths_darwin.go` / `paths_linux.go` with shared logic in `paths_common.go`. Injection assets (`app_index.js`, `app_package.json`) are embedded from `internal/discord/assets/` via `//go:embed`.
- `internal/models/` — channels (`iota` constants: Stable, Canary, PTB), GitHub release shapes, options, store models.
- `internal/output/` — the single gateway for user-facing output (`Printf`, `Println`, `Blank`, `NewTableWriter`, `SetWriters`). See §4.
- `internal/utils/` — download, path, and string helpers.
- `internal/wsl/` — WSL detection and Windows-home path mapping (lazy `sync.Once`, no `init()` cost outside WSL).
- `main.go` — entry point; `-ldflags -X main.version/commit/date` feed `cmd.SetVersionInfo`. A version of `"dev"` (the default) marks a debug build (`cmd.IsDebugBuild`).
- `scripts/completions.sh` — regenerates `completions/` (bash/zsh/fish); runs automatically in the GoReleaser `before` hook.
- `Taskfile.yml` — dev/build/test/lint/release shortcuts. `.goreleaser.yaml` — release build + publishing. `.golangci.yml` — lint config (`errcheck` excludes the `fmt.Fprint*` family).
- `package.json` — npm distribution wrapper (`@betterdiscord/cli` via `@go-task/go-npm`); it downloads release binaries, contains no JS logic. The version stays `0.0.0` in-repo — CI stamps it during release. There is no frontend and no `node_modules` toolchain to install for development.
- Generated / gitignored, never hand-edit or commit: `dist/`, `debug/`, `completions/`, `notes/`, `.task/`.

## 2. Build, dev & test commands

Prereqs: Go (matches `go.mod`), plus optionally [Task](https://taskfile.dev/), golangci-lint, and GoReleaser. No CGO, no platform build tags for local dev.

- `go run main.go <command>` or `task run -- <args>` — run the CLI locally (`task run` injects dev version ldflags).
- `task build` — local binary at `dist/bdcli`.
- `task test` (or `go test ./...`) — all tests.
- `task check` — `go fix` + `go fmt` + `go vet` + `golangci-lint` + tests, in one shot.
- `task ci` — what CI actually runs: deps + fix + fmt + vet + coverage + build. The golangci-lint GitHub action runs separately in `ci.yml`.
- `task build:all` / `task release:snapshot` — GoReleaser snapshot cross-builds (all OS/arch from one machine).

Before you claim a change is validated, run at minimum `gofmt`, `go vet ./...`, and `go test ./...`; prefer `task check` when golangci-lint is available.

## 3. Running isolated tests

- Narrow with the `-run` regex, e.g. `go test ./internal/models -run TestDiscordChannel` or `go test ./internal/discord -run TestInject`.
- Many tests gate on `runtime.GOOS` with `t.Skipf` when the OS doesn't match — run the subset aligned with your OS; these are guards, not cross-platform stubs. A green run on Linux does not prove the Windows/macOS paths.
- Network-dependent code is tested against `httptest` servers: endpoint URLs are declared as package-level `var`s (see `internal/betterdiscord/download.go`) precisely so tests can repoint them. Keep new endpoints in that pattern; never write a test that hits the real network.
- Coverage/profiling artifacts go to `debug/` (`task coverage:html`, `task bench:cpu`); that directory is gitignored.

## 4. Style rules

1. **Imports**
   - Standard library first, blank line, then everything else (external deps and `github.com/betterdiscord/cli/internal/...` share one alphabetized block in existing files). Match the surrounding file; `gofmt` is canonical.
2. **Formatting & naming**
   - Always run `gofmt`; tabs are canonical indentation.
   - Exported identifiers are PascalCase with doc comments; private helpers stay lowercase. Channel constants are grouped `iota` blocks.
   - Comments in this codebase explain *why* (see `internal/discord/injection.go`) — keep that bar for non-obvious logic, especially anything transactional or platform-specific.
3. **Output protocol** — the CLI equivalent of the installer's event protocol:
   - All user-facing output goes through `internal/output`, never `fmt.Print*` to stdout directly. This is what makes `--silent` / `BDCLI_SILENT` (which swap in `io.Discard`) and output-capturing tests work.
   - Status lines are emoji-prefixed: `✅` success, `❌` failure, `🔁` retry/fallback. Follow-up detail lines are indented with three spaces (`output.Printf("   %s\n", err.Error())`).
   - Tabular output uses `output.NewTableWriter()`; versions are normalized with `output.FormatVersion`.
4. **Error handling**
   - Prefer early returns over nested conditionals.
   - On failure: print the human-readable `❌` message via `output`, then `return` the error up to the Cobra `RunE` — `cmd.Execute()` prints it to stderr and exits 1. Don't both print and re-wrap the same message at every level.
   - Swallow errors only when a fallback genuinely handles them (e.g. website → GitHub download fallback), and always log the fallback so users can see what happened.
5. **Platform code**
   - OS-specific logic belongs in `paths_*.go`-style files or explicit `runtime.GOOS` switches, with WSL handled through `internal/wsl`. Don't sprinkle ad-hoc OS conditionals through command code.
6. **Injection safety**
   - The `app.asar` shadow injection is transactional: the original asar is preserved as `betterdiscord.app.asar` and any failure after the rename rolls back. Writability is probed (`probeWritable`) before any destructive step, and Snap installs are rejected up-front by design. Preserve all three properties when touching install/uninstall/repair flows.

## 5. Release & CI context

Releases are **tag-driven via GoReleaser**, all from a single Ubuntu runner (unlike the installer's per-OS matrix — everything here cross-compiles with `CGO_ENABLED=0`).

- **Trigger.** Pushing a `vX.Y.Z` tag runs `.github/workflows/release.yml`: `task ci`, then GoReleaser (linux/windows/darwin × amd64/arm64), then an npm publish. The `nightly` tag is ignored by GoReleaser; prereleases are auto-detected (`prerelease: auto`).
- **Version.** Injected via `-ldflags "-X main.version={{ .Version }}"` plus commit/date. CI also runs `npm version <tag>` before publishing, which is why `package.json` stays at `0.0.0` in the repo.
- **Artifacts.** `bdcli_<version>_<os>_<arch>.tar.gz` (`.zip` on Windows) containing the binary, README, LICENSE, and shell completions, plus `bdcli_checksums.txt`. Completions are generated by `scripts/completions.sh` in the `before` hook.
- **Publishing.** Homebrew cask pushed directly to `BetterDiscord/homebrew-tap` `main`; winget manifest PR'd from the `betterdiscord/winget-pkgs` fork to `microsoft/winget-pkgs` — both need `GH_PAT` (`GITHUB_TOKEN` can't push cross-repo). npm publish uses OIDC provenance (`id-token: write`), so there is no npm token secret.
- **CI on PRs/push** (`ci.yml`): `task ci` + the golangci-lint action, Linux only. Cross-platform correctness relies on the OS-gated tests and review — flag platform-specific risk in PRs since CI won't catch it.
- Keep `Taskfile.yml`, `.goreleaser.yaml`, and `release.yml` in sync when changing build flags or artifact names; the npm wrapper's `goBinary.url` template in `package.json` must keep matching GoReleaser's archive naming.

## 6. Documentation & collaborator expectations

- `README.md` is the user-facing reference (command table, compatibility matrix, FAQ). If you add/rename a command or flag, update the README command reference and help-output snippet in the same PR.
- `CONTRIBUTING.md` covers commit style (present-tense imperative, 72-char subject) and the Go styleguide summary. If it conflicts with these notes or with actual code, escalate via an issue rather than silently picking one.
- Keep AGENTS focused on machine-readable operating guidance; long-form docs belong in the README.

## 7. Cursor / Copilot / assistant config

- `CLAUDE.md` is a symlink to this file, so Claude Code reads the same guidance.
- There are no `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md` files; default assistant settings apply.
- If future agents add such rules, append them to this section so every agent reads the new constraints first.

## 8. Security, secrets, and environment

- **Never commit credentials, tokens, or `.env` files.** Release publishing relies solely on GitHub Actions secrets: `GITHUB_TOKEN` (release create/upload) and `GH_PAT` (homebrew tap push, winget cross-repo PR). npm uses OIDC provenance — no npm token exists. No secret is needed to build locally.
- **Supply chain.** The `betterdiscord.asar` is fetched over TLS from the official site with a GitHub-release fallback; there is no checksum verification of that asar yet — treat it as a known gap, not something to silently remove or quietly "fix" without discussion.
- **Dev builds must not fall back.** The `canary` development asar is GitHub-only; a failure there must never fall back to stable (a developer silently receiving stable is a near-undetectable footgun — this is documented in `download.go`; preserve it).
- **Snap is unsupported on purpose.** The read-only squashfs mount can't host the injection; the code rejects Snap installs with an actionable message before touching a running client. Don't add speculative Snap support.
- Keep tracked config (`go.mod`, `package.json`, `.goreleaser.yaml`) clean; do not add untracked credential copies to the tree.

## 9. Help & escalation

- If you find conflicting conventions between AGENTS, `CONTRIBUTING.md`, and actual code, escalate by opening an issue so the maintainers can harmonize them.
- Use the BetterDiscord Discord (linked from the README) when platform-integration questions exceed local knowledge.
- Document any new tooling (e.g., a new Task target or script) within this file so future agents immediately know how to operate the repo.

## 10. AI self-disclosure

If you are an autonomous agent submitting a PR where the human operator did not personally review the output before submission, say so explicitly in the PR description. For example:

> This PR was generated by [agent/tool]. The human operator [reviewed the diff and ran the checks locally / did not review the output before submission].

This is a transparency expectation, not a penalty against AI-assisted work. It helps calibrate review effort. PRs that appear to be unreviewed agent output without this disclosure will be closed without detailed feedback.
