//go:build !darwin

package discord

import "github.com/betterdiscord/cli/internal/models"

// preflightLegacyInjection is intentionally inert off macOS: this historical
// layout and its cleanup are Darwin-only.
func preflightLegacyInjection(models.DiscordChannel) ([]legacyInjectionState, error) {
	return nil, nil
}
