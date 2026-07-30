package discord

import (
	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/output"
	"github.com/betterdiscord/cli/internal/utils"
)

type DiscordInstall struct {
	CorePath  string                `json:"corePath"`
	Channel   models.DiscordChannel `json:"channel"`
	Version   string                `json:"version"`
	IsFlatpak bool                  `json:"isFlatpak"`
	IsSnap    bool                  `json:"isSnap"`
}

// InstallBD installs BetterDiscord into this Discord installation
func (discord *DiscordInstall) InstallBD(options models.InstallOptions) error {
	bd := discord.GetBetterDiscordInstall()

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
		install := discord.GetBetterDiscordInstall()
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
	if err := discord.UninstallBD(models.UninstallOptions{FullUninstall: false}); err != nil {
		return err
	}

	bd := discord.GetBetterDiscordInstall()

	if err := bd.Repair(discord.Channel); err != nil {
		return err
	}

	return nil
}

func (discord *DiscordInstall) GetBetterDiscordInstall() *betterdiscord.BDInstall {
	// Gets the global BetterDiscord install
	bd := betterdiscord.GetInstallation()

	// Snaps and flatpaks get their own local BD install
	if discord.IsSnap || discord.IsFlatpak {
		segment := "config"
		if discord.IsSnap {
			segment = ".config"
		}

		configPath, err := utils.FindSegment(discord.CorePath, segment)
		if err != nil {
			return nil
		}
		bd = betterdiscord.GetInstallation(configPath)
	}

	return bd
}
