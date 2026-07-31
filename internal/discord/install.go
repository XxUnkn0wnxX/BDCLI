package discord

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/output"
)

type DiscordInstall struct {
	ResourcesPath  string           `json:"resourcesPath"`
	Channel   models.DiscordChannel `json:"channel"`
	Version   string                `json:"version"`
	IsFlatpak bool                  `json:"isFlatpak"`
	IsSnap    bool                  `json:"isSnap"`
}

// InstallBD installs BetterDiscord into this Discord installation
func (discord *DiscordInstall) InstallBD(options models.InstallOptions) error {
	bd, err := discord.GetBetterDiscordInstall()
	if err != nil {
		return err
	}

	// Make BetterDiscord folders
	output.Println("🛠 Preparing BetterDiscord...")
	if err := bd.Prepare(); err != nil {
		return err
	}
	output.Println("✅ BetterDiscord prepared for install")
	output.Blank()

	// Download and write betterdiscord.asar
	output.Println("📥 Downloading BetterDiscord...")
	if err := bd.Download(); err != nil {
		return err
	}
	output.Println("✅ BetterDiscord downloaded")
	output.Blank()

	// Write injection script to discord_desktop_core/index.js
	output.Println("🔌 Injecting into Discord...")
	if err := discord.inject(bd); err != nil {
		return err
	}
	output.Println("✅ Injection successful")
	output.Blank()

	if options.RestartDiscord {
		// Terminate and restart Discord if possible
		output.Printf("🔄 Restarting %s...\n", discord.Channel.Name())
		if err := discord.restart(); err != nil {
			return err
		}
		output.Blank()
	}

	return nil
}

// UninstallBD removes BetterDiscord from this Discord installation
func (discord *DiscordInstall) UninstallBD(options models.UninstallOptions) error {
	output.Println("🧹 Removing injection...")
	if err := discord.uninject(); err != nil {
		return err
	}
	output.Blank()

	if options.FullUninstall {
		install, err := discord.GetBetterDiscordInstall()
		if err != nil {
			return err
		}
		if err := install.RemoveAll(); err != nil {
			return err
		}
		output.Blank()
	}

	if options.RestartDiscord {
		output.Printf("🔄 Restarting %s...\n", discord.Channel.Name())
		if err := discord.restart(); err != nil {
			return err
		}
		output.Blank()
	}

	return nil
}

// RepairBD repairs BetterDiscord for this Discord installation
func (discord *DiscordInstall) RepairBD(options models.RepairOptions) error {
	if err := discord.UninstallBD(models.UninstallOptions{FullUninstall: false, RestartDiscord: false}); err != nil {
		return err
	}

	bd, err := discord.GetBetterDiscordInstall()
	if err != nil {
		return err
	}

	if err := bd.Repair(discord.Channel); err != nil {
		return err
	}

	return nil
}

func (discord *DiscordInstall) GetBetterDiscordInstall() (*betterdiscord.BDInstall, error) {
	// Gets the global BetterDiscord install
	bd := betterdiscord.GetInstallation()

	// Flatpaks get their own local BD folder. The resources path is in the
	// read-only deployment tree, so we can't derive the sandbox config from it;
	// instead we compute the stable ~/.var/app/{id}/config location from the
	// channel. Inside the sandbox this dir is the app's $XDG_CONFIG_HOME, which
	// is exactly where the injected index.js looks for BetterDiscord at runtime.
	if discord.IsFlatpak {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		id := "com.discordapp." + strings.ReplaceAll(discord.Channel.Name(), " ", "")
		configPath := filepath.Join(home, ".var", "app", id, "config")
		bd = betterdiscord.GetInstallation(configPath)
	}

	return bd, nil
}
