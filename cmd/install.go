package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"

	"github.com/betterdiscord/cli/internal/discord"
	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/output"
)

func init() {
	installCmd.Flags().StringP("path", "p", "", "Path to a Discord installation")
	installCmd.Flags().StringP("channel", "c", "stable", "Discord release channel (stable|ptb|canary)")
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:     "install",
	Aliases: []string{"reinstall"},
	Short:   "Installs BetterDiscord to your Discord",
	Long:    "Install BetterDiscord by specifying either --path to a Discord install or --channel to auto-detect (default: stable).",
	RunE: func(cmd *cobra.Command, args []string) error {
		pathFlag, _ := cmd.Flags().GetString("path")
		channelFlag, _ := cmd.Flags().GetString("channel")

		pathProvided := pathFlag != ""
		channelProvided := cmd.Flags().Changed("channel")

		if pathProvided && channelProvided {
			return fmt.Errorf("--path and --channel are mutually exclusive")
		}

		var install *discord.DiscordInstall

		if pathProvided {
			install = discord.ResolvePath(pathFlag)
			if install == nil {
				return fmt.Errorf("could not find a valid Discord installation at %s", pathFlag)
			}
		} else {
			channel := models.ParseChannel(channelFlag)
			resourcesPath := discord.GetSuggestedPath(channel)
			install = discord.ResolvePath(resourcesPath)
			if install == nil {
				return fmt.Errorf("could not find a valid %s installation to install to", channelFlag)
			}
		}

		if err := install.InstallBD(models.InstallOptions{RestartDiscord: true, UseDevBuild: useDevBuild}); err != nil {
			return fmt.Errorf("installation failed: %w", err)
		}

		output.Printf("✅ BetterDiscord installed to %s\n", path.Dir(install.ResourcesPath))
		output.Blank()
		output.Printf("📋 Installation Summary:\n")
		output.Blank()
		output.Printf("   Release Channel: %s\n", install.Channel.Display())
		output.Printf("   Discord Version: %s\n", install.Version)
		output.Printf("   Install Type:    %s\n", func() string {
			if install.IsFlatpak {
				return "flatpak"
			} else if install.IsSnap {
				return "snap"
			}
			return "native"
		}())
		output.Printf("   Core Path:       %s\n", path.Dir(install.ResourcesPath))
		output.Blank()

		bdinstall, err := install.GetBetterDiscordInstall()
		if err != nil {
			output.Printf("failed to get BetterDiscord install info: %s", err.Error())
			return nil
		}
		if bdinstall == nil {
			output.Printf("BetterDiscord install info is nil")
			return nil
		}
		bdinstall.LogBuildinfo()
		return nil
	},
}
