//go:build darwin || linux || windows

package discord

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/betterdiscord/cli/internal/utils"
)

// validateWindowsStyleInstall validates a Windows-style Discord installation path.
// This is used for native Windows installs and WSL installs that point to Windows Discord.
// Windows Discord has a nested structure: Discord/app-1.0.9002/modules/discord_desktop_core-1/discord_desktop_core
func validateWindowsStyleInstall(proposed string) *DiscordInstall {
	var finalPath = ""
	var selected = filepath.Base(proposed)

	if strings.HasPrefix(selected, "Discord") {
		// Get version dir like app-1.0.9002
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

		// Get core wrap like discord_desktop_core-1
		dFiles, err = os.ReadDir(filepath.Join(proposed, versionDir, "modules"))
		if err != nil {
			return nil
		}
		candidates = utils.Filter(dFiles, func(file fs.DirEntry) bool {
			return file.IsDir() && strings.HasPrefix(file.Name(), "discord_desktop_core")
		})
		if len(candidates) == 0 {
			return nil
		}
		coreWrap := candidates[len(candidates)-1].Name()

		finalPath = filepath.Join(proposed, versionDir, "modules", coreWrap, "discord_desktop_core")
	}

	// Handle app-* directories (e.g., app-1.0.9002)
	if strings.HasPrefix(selected, "app-") {
		dFiles, err := os.ReadDir(filepath.Join(proposed, "modules"))
		if err != nil {
			return nil
		}

		candidates := utils.Filter(dFiles, func(file fs.DirEntry) bool {
			return file.IsDir() && strings.HasPrefix(file.Name(), "discord_desktop_core")
		})
		if len(candidates) == 0 {
			return nil
		}
		coreWrap := candidates[len(candidates)-1].Name()
		finalPath = filepath.Join(proposed, "modules", coreWrap, "discord_desktop_core")
	}

	if selected == "discord_desktop_core" {
		finalPath = proposed
	}

	// Verify the path and core.asar exist
	if utils.Exists(finalPath) && utils.Exists(filepath.Join(finalPath, "core.asar")) {
		return &DiscordInstall{
			CorePath:  finalPath,
			Channel:   GetChannel(finalPath),
			Version:   GetVersion(finalPath),
			IsFlatpak: false,
			IsSnap:    false,
		}
	}

	return nil
}

// validateUnixStyleInstall validates a Unix-style Discord installation path (Linux native, macOS).
// Unix Discord sometimes has a flatter structure: discord/0.0.35/modules/discord_desktop_core
// But sometimes it has the same pattern as Windows. This function detects both patterns and also identifies Flatpak and Snap installations if requested.
func validateUnixStyleInstall(proposed string, detectFlatpak bool, detectSnap bool) *DiscordInstall {
	var finalPath = ""
	var selected = filepath.Base(proposed)

	// Flatpak specific handling
	if strings.HasPrefix(selected, "com.discordapp") {
		channelPaths, err := os.ReadDir(filepath.Join(proposed, "config"))
		if err != nil {
			return nil
		}

		candidates := utils.Filter(channelPaths, func(file fs.DirEntry) bool {
			return file.IsDir() && strings.HasPrefix(file.Name(), "discord")
		})

		if len(candidates) == 0 {
			return nil
		}

		// Assume the first candidate is the correct one (e.g., discord or discordcanary)
		// Then set proposed and select so the remaining logic can find the core.asar
		//
		// TODO: This entire validation function could be refactored to use this fall-through logic
		// instead of trying to fully handle each pattern, but for now this is a simple way to support
		// Flatpak's extra nesting without breaking existing validations
		channelPath := candidates[0].Name()
		proposed = filepath.Join(proposed, "config", channelPath)
		selected = channelPath
	}

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

		// Get core wrap like discord_desktop_core-1
		dFiles, err = os.ReadDir(filepath.Join(proposed, versionDir, "modules"))
		if err != nil {
			return nil
		}
		candidates = utils.Filter(dFiles, func(file fs.DirEntry) bool {
			return file.IsDir() && strings.HasPrefix(file.Name(), "discord_desktop_core")
		})

		if len(candidates) == 0 {
			return nil
		}

		// If no core wrap is found, assume the structure is flatter and point directly to discord_desktop_core
		coreWrap := candidates[len(candidates)-1].Name()
		if coreWrap == "discord_desktop_core" {
			finalPath = filepath.Join(proposed, versionDir, "modules", "discord_desktop_core")
		} else {
			finalPath = filepath.Join(proposed, versionDir, "modules", coreWrap, "discord_desktop_core")
		}
	}

	// Handle version directories (e.g. app-0.0.35, 0.0.35)
	if strings.HasPrefix(selected, "app-") || versionRegex.MatchString(selected) {
		dFiles, err := os.ReadDir(filepath.Join(proposed, "modules"))
		if err != nil {
			return nil
		}

		candidates := utils.Filter(dFiles, func(file fs.DirEntry) bool {
			return file.IsDir() && strings.HasPrefix(file.Name(), "discord_desktop_core")
		})

		if len(candidates) == 0 {
			return nil
		}

		// If no core wrap is found, assume the structure is flatter and point directly to discord_desktop_core
		coreWrap := candidates[len(candidates)-1].Name()
		if coreWrap == "discord_desktop_core" {
			finalPath = filepath.Join(proposed, "modules", "discord_desktop_core")
		} else {
			finalPath = filepath.Join(proposed, "modules", coreWrap, "discord_desktop_core")
		}
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
