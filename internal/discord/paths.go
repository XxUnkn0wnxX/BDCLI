package discord

import (
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/betterdiscord/cli/internal/models"
)

var searchPaths []string
var versionRegex = regexp.MustCompile(`[0-9]+\.[0-9]+\.[0-9]+`)
var allDiscordInstalls map[models.DiscordChannel][]*DiscordInstall

func GetAllInstalls() map[models.DiscordChannel][]*DiscordInstall {
	allDiscordInstalls = map[models.DiscordChannel][]*DiscordInstall{}

	for _, path := range searchPaths {
		if result := Validate(path); result != nil {
			allDiscordInstalls[result.Channel] = append(allDiscordInstalls[result.Channel], result)
		}
	}

	sortInstalls()

	return allDiscordInstalls
}

func GetVersion(proposed string) string {
	for folder := range strings.SplitSeq(filepath.ToSlash(proposed), "/") {
		if version := versionRegex.FindString(folder); version != "" {
			return version
		}
	}
	return ""
}

func GetChannel(proposed string) models.DiscordChannel {
	// Iterate from the leaf toward the root: the channel identifier always sits
	// closest to the leaf (e.g. `.../discordcanary/app-x/resources`), so scanning
	// backwards avoids false matches on a parent segment that happens to contain a
	// channel name (e.g. a home dir at `/home/discord`).
	// Normalize to forward slashes before splitting so a Windows path that mixes
	// separators (backslashes and forward slashes, which the OS treats
	// interchangeably) still segments cleanly.
	segments := strings.Split(filepath.ToSlash(proposed), "/")

	for _, segment := range slices.Backward(segments) {
		// Normalize the segment so macOS bundle names ("Discord Canary.app") and
		// flatpak channel dirs ("discord-canary") both match the channel names
		// ("discordcanary").
		normalized := strings.ToLower(segment)
		normalized = strings.TrimSuffix(normalized, ".app")
		normalized = strings.ReplaceAll(normalized, " ", "")
		normalized = strings.ReplaceAll(normalized, "-", "")
		for _, channel := range models.Channels {
			if normalized == strings.ReplaceAll(strings.ToLower(channel.Name()), " ", "") {
				return channel
			}
		}
	}
	return models.Stable
}

func GetSuggestedPath(channel models.DiscordChannel) string {
	if len(allDiscordInstalls[channel]) > 0 {
		return allDiscordInstalls[channel][0].ResourcesPath
	}
	return ""
}

func AddCustomPath(proposed string) *DiscordInstall {
	result := Validate(proposed)
	if result == nil {
		return nil
	}

	// Check if this already exists in our list and return reference
	index := slices.IndexFunc(allDiscordInstalls[result.Channel], func(d *DiscordInstall) bool { return d.ResourcesPath == result.ResourcesPath })
	if index >= 0 {
		return allDiscordInstalls[result.Channel][index]
	}

	allDiscordInstalls[result.Channel] = append(allDiscordInstalls[result.Channel], result)

	sortInstalls()

	return result
}

func ResolvePath(proposed string) *DiscordInstall {
	for channel := range allDiscordInstalls {
		index := slices.IndexFunc(allDiscordInstalls[channel], func(d *DiscordInstall) bool { return d.ResourcesPath == proposed })
		if index >= 0 {
			return allDiscordInstalls[channel][index]
		}
	}

	// If it wasn't found as an existing install, try to add it
	return AddCustomPath(proposed)
}

func sortInstalls() {
	for channel := range allDiscordInstalls {
		slices.SortFunc(allDiscordInstalls[channel], func(a, b *DiscordInstall) int {
			switch {
			case a.Version > b.Version:
				return -1
			case b.Version > a.Version:
				return 1
			}
			return 0
		})
	}
}
