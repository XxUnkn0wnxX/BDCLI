package discord

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/wsl"
)

// Snap is intentionally omitted: its read-only squashfs mount can't host
// the app.asar shadow, so the new injection method does not support it.
func init() {
	config, errConfig := os.UserConfigDir()
	home, errHome := os.UserHomeDir()

	// Flatpak (global). The app.asar lives in the read-only deployment files.
	// Example: `/var/lib/flatpak/app/com.discordapp.DiscordCanary/current/active/files/discord-canary/resources/app.asar`.
	// This has no home/config dependency, so it's always searched.
	paths := []string{
		filepath.Join("/var", "lib", "flatpak", "app", "com.discordapp.{CHANNEL}", "current", "active", "files", "{channel-}", "resources"),
	}

	// Only search config/home-relative locations when those dirs resolved;
	// otherwise the joins would produce relative paths anchored at the current
	// working directory.

	// Native. The new updater lays out versioned app dirs under `~/.config`.
	// Example: `~/.config/discordcanary`.
	// Resources: `~/.config/discordcanary/app-0.0.90/resources/app.asar`.
	if errConfig == nil && config != "" {
		paths = append(paths, filepath.Join(config, "{channel}"))
	}

	// Flatpak (user). Same layout under the per-user flatpak tree (writable).
	// Example: `~/.local/share/flatpak/app/com.discordapp.DiscordCanary/current/active/files/discord-canary/resources/app.asar`.
	if errHome == nil && home != "" {
		paths = append(paths, filepath.Join(home, ".local", "share", "flatpak", "app", "com.discordapp.{CHANNEL}", "current", "active", "files", "{channel-}", "resources"))
	}

	if wsl.IsWSL() {
		winHome, err := wsl.WindowsHome()
		if err == nil && winHome != "" {
			// WSL. Windows Discord installs under the Windows user's AppData folder.
			// Example: `/mnt/c/Users/Username/AppData/Local/DiscordCanary`.
			// Resources: `/mnt/c/Users/Username/AppData/Local/DiscordCanary/app-1.0.9218/resources/app.asar`.
			paths = append(paths, filepath.Join(winHome, "AppData", "Local", "{CHANNEL}"))
		}
	}

	for _, channel := range models.Channels {
		for _, path := range paths {
			upper := strings.ReplaceAll(channel.Name(), " ", "")
			lower := strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "")
			dash := strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "-")
			folder := strings.ReplaceAll(path, "{CHANNEL}", upper)
			folder = strings.ReplaceAll(folder, "{channel}", lower)
			folder = strings.ReplaceAll(folder, "{channel-}", dash)
			searchPaths = append(searchPaths, folder)
		}
	}

	allDiscordInstalls = GetAllInstalls()
}

// Validate validates a Discord installation path on Linux.
// For WSL environments, it uses Windows-style validation.
// For native Linux, it detects Flatpak and Snap installations.
func Validate(proposed string) *DiscordInstall {
	if wsl.IsWSL() {
		return validateWindowsStyleInstall(proposed)
	}

	// Native Linux validation with Flatpak and Snap detection
	return validateUnixStyleInstall(proposed, true, true)
}
