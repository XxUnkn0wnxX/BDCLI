//go:build darwin || linux

package discord

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/betterdiscord/cli/internal/utils"
)

// validateUnixStyleInstall validates a Unix-style Discord installation path (Linux native, macOS).
// Unix Discord has a flatter structure: discord/0.0.35/modules/discord_desktop_core
func validateUnixStyleInstall(proposed string, detectFlatpak bool, detectSnap bool) *DiscordInstall {
	var finalPath = ""
	var selected = filepath.Base(proposed)

	if strings.HasPrefix(strings.ToLower(selected), "discord") {
		// Get version dir like 0.0.35
		dFiles, err := os.ReadDir(proposed)
		if err != nil {
			return nil
		}

		candidates := utils.Filter(dFiles, func(file fs.DirEntry) bool {
			return file.IsDir() && versionRegex.MatchString(file.Name())
		})
		if len(candidates) == 0 {
			return nil
		}
		sort.Slice(candidates, func(i, j int) bool { return candidates[i].Name() < candidates[j].Name() })
		versionDir := candidates[len(candidates)-1].Name()
		finalPath = filepath.Join(proposed, versionDir, "modules", "discord_desktop_core")
	}

	// Handle version directories (e.g., 0.0.35)
	if len(strings.Split(selected, ".")) == 3 {
		finalPath = filepath.Join(proposed, "modules", "discord_desktop_core")
	}

	if selected == "modules" {
		finalPath = filepath.Join(proposed, "discord_desktop_core")
	}

	if selected == "discord_desktop_core" {
		finalPath = proposed
	}

	// Verify the path and core.asar exist
	if utils.Exists(finalPath) && utils.Exists(filepath.Join(finalPath, "core.asar")) {
		isFlatpak := false
		isSnap := false

		if detectFlatpak {
			isFlatpak = strings.Contains(finalPath, "com.discordapp.")
		}
		if detectSnap {
			isSnap = strings.Contains(finalPath, "snap/")
		}

		return &DiscordInstall{
			CorePath:  finalPath,
			Channel:   GetChannel(finalPath),
			Version:   GetVersion(finalPath),
			IsFlatpak: isFlatpak,
			IsSnap:    isSnap,
		}
	}

	return nil
}
