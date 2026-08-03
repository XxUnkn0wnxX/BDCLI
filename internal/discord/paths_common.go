//go:build darwin || linux || windows

package discord

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/utils"
)

// buildInfo mirrors the fields we care about in Discord's resources/build_info.json.
type buildInfo struct {
	ReleaseChannel string `json:"releaseChannel"`
	Version        string `json:"version"`
}

// readBuildInfo reads resources/build_info.json. The second return is false when
// the file is absent or unparseable, so callers can fall back to path parsing.
func readBuildInfo(resourcesDir string) (buildInfo, bool) {
	data, err := os.ReadFile(filepath.Join(resourcesDir, "build_info.json"))
	if err != nil {
		return buildInfo{}, false
	}

	var info buildInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return buildInfo{}, false
	}

	return info, true
}

// hasDiscordApp reports whether dir is a Discord `resources` directory — that
// is, it contains Discord's app archive in *either* state:
//   - `app.asar` — a pristine (or freshly updated) install, and
//   - `betterdiscord.app.asar` — the original preserved after BetterDiscord
//     injects its shadow `app/` folder (at which point `app.asar` no longer
//     exists).
//
// Checking both is essential: once injected, an install would otherwise stop
// resolving, so users could no longer repair or — critically — uninstall it.
func hasDiscordApp(dir string) bool {
	return utils.Exists(filepath.Join(dir, "app.asar")) || utils.Exists(filepath.Join(dir, "betterdiscord.app.asar"))
}

// latestAppDir returns the highest-versioned `app-{version}` child of base whose
// resources dir actually holds a Discord app, or "" when none qualify. Skipping
// broken/incomplete version dirs (e.g. from an interrupted Discord update) lets
// resolution fall back to a slightly older but valid install instead of failing.
// Sorting is numeric so 1.0.10000 beats 1.0.9999.
func latestAppDir(base string) string {
	entries, err := os.ReadDir(base)
	if err != nil {
		return ""
	}

	bestName := ""
	bestVersion := ""
	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "app-") {
			continue
		}
		version := strings.TrimPrefix(entry.Name(), "app-")
		if !versionRegex.MatchString(version) {
			continue
		}
		if !hasDiscordApp(filepath.Join(base, entry.Name(), "resources")) {
			continue
		}
		if bestName == "" || utils.CompareVersions(version, bestVersion) > 0 {
			bestName, bestVersion = entry.Name(), version
		}
	}

	return bestName
}

// isSnapPath reports whether a resolved resources path lives under a Snap mount
// (/snap/… or /var/lib/snapd/snap/…). Anchoring to the mount points avoids
// false-positives on unrelated paths that merely contain a "snap" segment — e.g.
// the home directory of a user named "snap" (/home/snap/…).
func isSnapPath(path string) bool {
	sep := string(filepath.Separator)
	return strings.HasPrefix(path, "snap"+sep) ||
		strings.HasPrefix(path, sep+"snap"+sep) ||
		strings.HasPrefix(path, sep+"var"+sep+"lib"+sep+"snapd"+sep+"snap"+sep)
}

// resolveResources locates the Discord `resources` directory (holding Discord's
// app archive — see hasDiscordApp) from a variety of proposed inputs, returning
// "" when none is found:
//   - a resources dir itself (or macOS Contents/Resources) — the archive is directly inside
//   - an `app-{version}` dir — drills into its `resources`
//   - a dir that directly contains a `resources` child (flatpak files/{channel-})
//   - a base holding `app-{version}` dirs (Discord root / channel config dir) — picks latest
func resolveResources(proposed string) string {
	if proposed == "" {
		return ""
	}

	// The proposed path is already the resources dir (or macOS Contents/Resources).
	if hasDiscordApp(proposed) {
		return proposed
	}

	if strings.HasPrefix(filepath.Base(proposed), "app-") {
		if res := filepath.Join(proposed, "resources"); hasDiscordApp(res) {
			return res
		}
		return ""
	}

	// A dir with a direct `resources` child (flatpak files/{channel-}).
	if res := filepath.Join(proposed, "resources"); hasDiscordApp(res) {
		return res
	}

	// A macOS app bundle: app.asar lives in Contents/Resources.
	if res := filepath.Join(proposed, "Contents", "Resources"); hasDiscordApp(res) {
		return res
	}

	// A base containing versioned app dirs (Windows Discord root, Linux channel dir).
	if latest := latestAppDir(proposed); latest != "" {
		if res := filepath.Join(proposed, latest, "resources"); hasDiscordApp(res) {
			return res
		}
	}

	return ""
}

// newResourcesInstall builds a DiscordInstall for a resolved resources dir,
// preferring build_info.json for channel/version and falling back to the path.
func newResourcesInstall(resourcesDir string) *DiscordInstall {
	channel := GetChannel(resourcesDir)
	version := GetVersion(resourcesDir)

	if info, ok := readBuildInfo(resourcesDir); ok {
		if info.ReleaseChannel != "" {
			channel = models.ParseChannel(info.ReleaseChannel)
		}
		if info.Version != "" {
			version = info.Version
		}
	}

	return &DiscordInstall{
		ResourcesPath: resourcesDir,
		Channel:       channel,
		Version:       version,
	}
}

// validateWindowsStyleInstall validates a Windows-style install (native Windows
// and WSL pointing at Windows Discord). The new updater lays out installs as
// Discord/app-{version}/resources/app.asar.
func validateWindowsStyleInstall(proposed string) *DiscordInstall {
	resources := resolveResources(proposed)
	if resources == "" {
		return nil
	}
	return newResourcesInstall(resources)
}

// validateUnixStyleInstall validates a Unix-style install (Linux native, macOS).
// Linux native mirrors the Windows layout under the config dir
// (~/.config/{channel}/app-{version}/resources); macOS keeps app.asar directly in
// the bundle's Contents/Resources. Flatpak/Snap are flagged via the resolved path.
func validateUnixStyleInstall(proposed string, detectFlatpak bool, detectSnap bool) *DiscordInstall {
	resources := resolveResources(proposed)
	if resources == "" {
		return nil
	}

	install := newResourcesInstall(resources)

	// Heuristic: infer packaging format from the resolved path. These substring
	// checks match the real Flatpak/Snap layouts in practice.
	if detectFlatpak {
		install.IsFlatpak = strings.Contains(resources, "com.discordapp.")
	}
	if detectSnap {
		install.IsSnap = isSnapPath(resources)
	}

	return install
}
