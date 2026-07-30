package discord

import (
	"os"
	"path/filepath"

	"github.com/betterdiscord/cli/internal/models"
)

func init() {
	home, _ := os.UserHomeDir()

	// On macOS the app.asar lives inside the application bundle
	// (Discord.app/Contents/Resources), not under Application Support. Search the
	// standard install locations for each channel's bundle.
	bases := []string{
		filepath.Join("/", "Applications"),
		filepath.Join(home, "Applications"),
	}

	for _, channel := range models.Channels {
		bundle := channel.Name() + ".app"
		for _, base := range bases {
			searchPaths = append(searchPaths, filepath.Join(base, bundle))
		}
	}

	allDiscordInstalls = GetAllInstalls()
}

func Validate(proposed string) *DiscordInstall {
	return validateUnixStyleInstall(proposed, false, false)
}
