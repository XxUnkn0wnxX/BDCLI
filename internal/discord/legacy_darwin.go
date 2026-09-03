//go:build darwin

package discord

import (
	"fmt"
	"os"

	"github.com/betterdiscord/cli/internal/models"
)

// preflightLegacyInjection discovers legacy loaders under macOS's user config
// directory. It performs no writes.
func preflightLegacyInjection(channel models.DiscordChannel) ([]legacyInjectionState, error) {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("determine legacy Discord config directory: %w", err)
	}
	return findLegacyInjectionStates(configRoot, channel)
}
