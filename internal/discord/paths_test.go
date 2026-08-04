package discord

import (
	"path/filepath"
	"testing"

	"github.com/betterdiscord/cli/internal/models"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		// Old core structure
		{
			name:     "Old structure: Version in middle of path",
			path:     "/usr/share/discord/0.0.35/modules",
			expected: "0.0.35",
		},
		{
			name:     "Old structure: Version at end of path",
			path:     "/home/user/.config/discord/0.0.36",
			expected: "0.0.36",
		},
		{
			name:     "Old structure: Full path with modules",
			path:     "/home/user/.config/discord/0.0.90/modules/discord_desktop_core/core.asar",
			expected: "0.0.90",
		},

		// New core structure (app-X.X.X format)
		{
			name:     "New structure: Version with app prefix",
			path:     "/home/user/.config/discord/app-0.0.90/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expected: "0.0.90",
		},
		{
			name:     "New structure: Canary channel with app prefix",
			path:     "/home/user/.config/discordcanary/app-0.0.200/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expected: "0.0.200",
		},
		{
			name:     "New structure: PTB channel with app prefix",
			path:     "/home/user/.config/discordptb/app-0.0.150/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expected: "0.0.150",
		},

		// Edge cases
		{
			name:     "Multiple versions (should return first)",
			path:     "/usr/share/1.2.3/discord/0.0.35/modules",
			expected: "1.2.3",
		},
		{
			name:     "No version in path",
			path:     "/usr/share/discord/modules",
			expected: "",
		},
		{
			name:     "Version with many digits",
			path:     "/opt/discord/123.456.789/core",
			expected: "123.456.789",
		},
		{
			name:     "Empty path",
			path:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetVersion(tt.path)
			if result != tt.expected {
				t.Errorf("GetVersion(%s) = %s, expected %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetChannel(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected models.DiscordChannel
	}{
		// Old core structure
		{
			name:     "Old structure: Stable in path (lowercase)",
			path:     "/usr/share/discord/modules",
			expected: models.Stable,
		},
		{
			name:     "Old structure: Canary in path (lowercase)",
			path:     "/usr/share/discordcanary/modules",
			expected: models.Canary,
		},
		{
			name:     "Old structure: PTB in path (lowercase)",
			path:     "/usr/share/discordptb/modules",
			expected: models.PTB,
		},
		{
			name:     "Old structure: DiscordCanary without space",
			path:     "/home/user/.config/DiscordCanary/modules",
			expected: models.Canary,
		},
		{
			name:     "Old structure: Full path with version",
			path:     "/home/user/.config/discord/0.0.90/modules/discord_desktop_core/core.asar",
			expected: models.Stable,
		},

		// New core structure (app-X.X.X format)
		{
			name:     "New structure: Stable with app-version",
			path:     "/home/user/.config/discord/app-0.0.90/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expected: models.Stable,
		},
		{
			name:     "New structure: Canary with app-version",
			path:     "/home/user/.config/discordcanary/app-0.0.200/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expected: models.Canary,
		},
		{
			name:     "New structure: PTB with app-version",
			path:     "/home/user/.config/discordptb/app-0.0.150/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expected: models.PTB,
		},

		// Edge cases
		{
			name:     "DiscordPTB without space",
			path:     "/home/user/.config/DiscordPTB/modules",
			expected: models.PTB,
		},
		{
			name:     "No channel identifier defaults to Stable",
			path:     "/some/random/path/modules",
			expected: models.Stable,
		},
		{
			name:     "Multiple Discord mentions (nearest wins)",
			path:     filepath.Join("discordcanary", "discord", "modules"),
			expected: models.Stable, // The nearest segment is "discord", which maps to Stable
		},
		{
			name:     "Empty path defaults to Stable",
			path:     "",
			expected: models.Stable,
		},

		// New injection
		{
			name:     "macOS bundle name",
			path:     filepath.Join("/Applications", "Discord Canary.app", "Contents", "Resources"),
			expected: models.Canary,
		},
		{
			name:     "macOS stable bundle name",
			path:     filepath.Join("/Applications", "Discord.app", "Contents", "Resources"),
			expected: models.Stable,
		},
		{
			name:     "flatpak dashed canary dir",
			path:     filepath.Join("/var", "lib", "flatpak", "app", "com.discordapp.DiscordCanary", "current", "active", "files", "discord-canary", "resources"),
			expected: models.Canary,
		},
		{
			name:     "flatpak dashed ptb dir",
			path:     filepath.Join("/var", "lib", "flatpak", "app", "com.discordapp.DiscordPTB", "current", "active", "files", "discord-ptb", "resources"),
			expected: models.PTB,
		},
		{
			name:     "flatpak stable dir",
			path:     filepath.Join("/var", "lib", "flatpak", "app", "com.discordapp.Discord", "current", "active", "files", "discord", "resources"),
			expected: models.Stable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetChannel(tt.path)
			if result != tt.expected {
				t.Errorf("GetChannel(%q) = %v (%s), expected %v (%s)",
					tt.path, result, result.String(), tt.expected, tt.expected.String())
			}
		})
	}
}

func TestGetChannel_CaseInsensitive(t *testing.T) {
	tests := []struct {
		path     string
		expected models.DiscordChannel
	}{
		// Old structure - case insensitive
		{"/usr/share/DISCORD/modules", models.Stable},
		{"/usr/share/Discord/modules", models.Stable},
		{"/usr/share/DISCORDCANARY/modules", models.Canary},
		{"/usr/share/DiscordCanary/modules", models.Canary},
		{"/usr/share/DISCORDPTB/modules", models.PTB},
		{"/usr/share/DiscordPTB/modules", models.PTB},

		// New structure - case insensitive (with app-version)
		{"/usr/share/DISCORD/app-0.0.90/modules/discord_desktop_core-1/discord_desktop_core/core.asar", models.Stable},
		{"/usr/share/Discord/app-0.0.90/modules/discord_desktop_core-1/discord_desktop_core/core.asar", models.Stable},
		{"/usr/share/DISCORDCANARY/app-0.0.200/modules/discord_desktop_core-1/discord_desktop_core/core.asar", models.Canary},
		{"/usr/share/DiscordCanary/app-0.0.200/modules/discord_desktop_core-1/discord_desktop_core/core.asar", models.Canary},
		{"/usr/share/DISCORDPTB/app-0.0.150/modules/discord_desktop_core-1/discord_desktop_core/core.asar", models.PTB},
		{"/usr/share/DiscordPTB/app-0.0.150/modules/discord_desktop_core-1/discord_desktop_core/core.asar", models.PTB},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := GetChannel(tt.path)
			if result != tt.expected {
				t.Errorf("GetChannel(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestGetSuggestedPath(t *testing.T) {
	// Reset allDiscordInstalls for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Test empty installs
	result := GetSuggestedPath(models.Stable)
	if result != "" {
		t.Errorf("GetSuggestedPath with no installs should return empty string, got %s", result)
	}

	// Add some test installs - mix of old and new path formats
	oldCorePath := "/usr/share/discord/0.0.35"
	newCorePath := "/home/user/.config/discord/app-0.0.35/modules/discord_desktop_core-1/discord_desktop_core/core.asar"

	allDiscordInstalls[models.Stable] = []*DiscordInstall{
		{ResourcesPath: oldCorePath, Version: "0.0.35"},
		{ResourcesPath: "/usr/share/discord/0.0.34", Version: "0.0.34"},
	}

	allDiscordInstalls[models.Canary] = []*DiscordInstall{
		{ResourcesPath: newCorePath, Version: "0.0.200"}, // New format
	}

	// Test that it returns the first install (old format)
	stableResult := GetSuggestedPath(models.Stable)
	if stableResult != oldCorePath {
		t.Errorf("GetSuggestedPath(Stable) = %s, expected %s", stableResult, oldCorePath)
	}

	// Test that it returns the first install (new format)
	canaryResult := GetSuggestedPath(models.Canary)
	if canaryResult != newCorePath {
		t.Errorf("GetSuggestedPath(Canary) = %s, expected %s", canaryResult, newCorePath)
	}

	// Test channel with no installs
	ptbResult := GetSuggestedPath(models.PTB)
	if ptbResult != "" {
		t.Errorf("GetSuggestedPath(PTB) with no PTB installs should return empty string, got %s", ptbResult)
	}
}

func TestAddCustomPath(t *testing.T) {
	// This test is limited because Validate() depends on OS-specific paths
	// We're mainly testing the logic around adding and deduplication

	// Reset for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Test with invalid path (will return nil since Validate will fail)
	result := AddCustomPath("/nonexistent/invalid/path")
	if result != nil {
		t.Error("AddCustomPath with invalid path should return nil")
	}

	// Further testing would require mocking the Validate function
	// or setting up actual Discord installation directories
}

func TestResolvePath(t *testing.T) {
	// Reset for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Add a test install with new path format
	testInstall := &DiscordInstall{
		ResourcesPath: "/home/user/.config/discord/app-1.0.0/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
		Channel:       models.Stable,
		Version:       "1.0.0",
	}
	allDiscordInstalls[models.Stable] = []*DiscordInstall{testInstall}

	// Test resolving existing path
	result := ResolvePath(testInstall.ResourcesPath)
	if result != testInstall {
		t.Error("ResolvePath should return the existing install")
	}

	// Test resolving non-existent path (will try AddCustomPath and likely return nil)
	result2 := ResolvePath("/nonexistent/path")
	if result2 != nil {
		// This might succeed or fail depending on whether Validate passes
		// In most test environments, it should return nil
		t.Log("ResolvePath returned non-nil for non-existent path (may be valid in some environments)")
	}
}

// TestOldAndNewCorePathsCompatibility ensures both old and new Discord core path structures
// are detected and handled correctly by the path resolution functions.
func TestOldAndNewCorePathsCompatibility(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		expectedChannel models.DiscordChannel
		expectedVersion string
	}{
		// Old core structure - Linux
		{
			name:            "Old: Stable core path",
			path:            "/home/user/.config/discord/0.0.90/modules/discord_desktop_core/core.asar",
			expectedChannel: models.Stable,
			expectedVersion: "0.0.90",
		},
		{
			name:            "Old: Canary core path",
			path:            "/home/user/.config/discordcanary/0.0.100/modules/discord_desktop_core/core.asar",
			expectedChannel: models.Canary,
			expectedVersion: "0.0.100",
		},
		{
			name:            "Old: PTB core path",
			path:            "/home/user/.config/discordptb/0.0.80/modules/discord_desktop_core/core.asar",
			expectedChannel: models.PTB,
			expectedVersion: "0.0.80",
		},

		// New core structure (app-version with -1 suffix) - Linux
		{
			name:            "New: Stable core path with app prefix",
			path:            "/home/user/.config/discord/app-0.0.90/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expectedChannel: models.Stable,
			expectedVersion: "0.0.90",
		},
		{
			name:            "New: Canary core path with app prefix",
			path:            "/home/user/.config/discordcanary/app-0.0.200/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expectedChannel: models.Canary,
			expectedVersion: "0.0.200",
		},
		{
			name:            "New: PTB core path with app prefix",
			path:            "/home/user/.config/discordptb/app-0.0.150/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expectedChannel: models.PTB,
			expectedVersion: "0.0.150",
		},

		// Flatpak old structure
		{
			name:            "Flatpak: Old Canary path",
			path:            "/home/user/.var/app/com.discordapp.DiscordCanary/config/discordcanary/0.0.200/modules/discord_desktop_core/core.asar",
			expectedChannel: models.Canary,
			expectedVersion: "0.0.200",
		},

		// Flatpak new structure
		{
			name:            "Flatpak: New Canary path",
			path:            "/home/user/.var/app/com.discordapp.DiscordCanary/config/discordcanary/app-0.0.200/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expectedChannel: models.Canary,
			expectedVersion: "0.0.200",
		},

		// Snap old structure
		{
			name:            "Snap: Old stable path",
			path:            "/home/user/snap/discord/current/.config/discord/0.0.90/modules/discord_desktop_core/core.asar",
			expectedChannel: models.Stable,
			expectedVersion: "0.0.90",
		},

		// Snap new structure
		{
			name:            "Snap: New stable path",
			path:            "/home/user/snap/discord/current/.config/discord/app-0.0.90/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expectedChannel: models.Stable,
			expectedVersion: "0.0.90",
		},

		// Snap canary - new structure
		{
			name:            "Snap: New Canary path",
			path:            "/home/user/snap/discord-canary/current/.config/discordcanary/app-0.0.200/modules/discord_desktop_core-1/discord_desktop_core/core.asar",
			expectedChannel: models.Canary,
			expectedVersion: "0.0.200",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			channel := GetChannel(tt.path)
			version := GetVersion(tt.path)

			if channel != tt.expectedChannel {
				t.Errorf("GetChannel: got %v, expected %v", channel, tt.expectedChannel)
			}

			if version != tt.expectedVersion {
				t.Errorf("GetVersion: got %s, expected %s", version, tt.expectedVersion)
			}
		})
	}
}

func TestSortInstalls(t *testing.T) {
	// Reset for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Add unsorted installs - mix of old and new path formats
	allDiscordInstalls[models.Stable] = []*DiscordInstall{
		{ResourcesPath: "/path1", Version: "0.0.34", Channel: models.Stable},                                                                                              // Old format
		{ResourcesPath: "/home/user/.config/discord/app-0.0.36/modules/discord_desktop_core-1/discord_desktop_core/core.asar", Version: "0.0.36", Channel: models.Stable}, // New format
		{ResourcesPath: "/path3", Version: "0.0.35", Channel: models.Stable},                                                                                              // Old format
	}

	// Sort them
	sortInstalls()

	// Verify sorted in descending order by version (format doesn't matter)
	installs := allDiscordInstalls[models.Stable]
	if len(installs) != 3 {
		t.Fatalf("Expected 3 installs, got %d", len(installs))
	}

	if installs[0].Version != "0.0.36" {
		t.Errorf("First install should have version 0.0.36, got %s", installs[0].Version)
	}
	if installs[1].Version != "0.0.35" {
		t.Errorf("Second install should have version 0.0.35, got %s", installs[1].Version)
	}
	if installs[2].Version != "0.0.34" {
		t.Errorf("Third install should have version 0.0.34, got %s", installs[2].Version)
	}
}

func TestSortInstalls_MultipleChannels(t *testing.T) {
	// Reset for testing
	allDiscordInstalls = make(map[models.DiscordChannel][]*DiscordInstall)

	// Add unsorted installs for multiple channels - mix of old and new formats
	allDiscordInstalls[models.Stable] = []*DiscordInstall{
		{ResourcesPath: "/stable1", Version: "1.0.0", Channel: models.Stable},
		{ResourcesPath: "/home/user/.config/discord/app-1.0.2/modules/discord_desktop_core-1/discord_desktop_core/core.asar", Version: "1.0.2", Channel: models.Stable}, // New format
	}

	allDiscordInstalls[models.Canary] = []*DiscordInstall{
		{ResourcesPath: "/canary1", Version: "0.0.100", Channel: models.Canary},
		{ResourcesPath: "/home/user/.config/discordcanary/app-0.0.150/modules/discord_desktop_core-1/discord_desktop_core/core.asar", Version: "0.0.150", Channel: models.Canary}, // New format
		{ResourcesPath: "/canary3", Version: "0.0.125", Channel: models.Canary},
	}

	// Sort them
	sortInstalls()

	// Verify Stable channel is sorted (format doesn't matter)
	stableInstalls := allDiscordInstalls[models.Stable]
	if stableInstalls[0].Version != "1.0.2" {
		t.Errorf("Stable: First version should be 1.0.2, got %s", stableInstalls[0].Version)
	}
	if stableInstalls[1].Version != "1.0.0" {
		t.Errorf("Stable: Second version should be 1.0.0, got %s", stableInstalls[1].Version)
	}

	// Verify Canary channel is sorted (format doesn't matter)
	canaryInstalls := allDiscordInstalls[models.Canary]
	if canaryInstalls[0].Version != "0.0.150" {
		t.Errorf("Canary: First version should be 0.0.150, got %s", canaryInstalls[0].Version)
	}
	if canaryInstalls[1].Version != "0.0.125" {
		t.Errorf("Canary: Second version should be 0.0.125, got %s", canaryInstalls[1].Version)
	}
	if canaryInstalls[2].Version != "0.0.100" {
		t.Errorf("Canary: Third version should be 0.0.100, got %s", canaryInstalls[2].Version)
	}
}
