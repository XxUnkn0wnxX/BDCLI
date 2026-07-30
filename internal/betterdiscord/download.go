package betterdiscord

import (
	"fmt"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/output"
	"github.com/betterdiscord/cli/internal/utils"
)

// Endpoints for fetching the BetterDiscord asar. Declared as package vars so
// tests can point them at a local httptest server.
var (
	websiteAsarURL         = "https://betterdiscord.app/Download/betterdiscord.asar"
	githubLatestReleaseURL = "https://api.github.com/repos/BetterDiscord/BetterDiscord/releases/latest"
)

func (i *BDInstall) download() error {
	if i.hasDownloaded {
		output.Printf("✅ Already downloaded to %s\n", i.asar)
		return nil
	}

	resp, err := utils.DownloadFile(websiteAsarURL, i.asar)
	if err == nil {
		version := resp.Header.Get("x-bd-version")
		if version == "" {
			output.Println("✅ Downloaded BetterDiscord from the official website")
		} else {
			output.Printf("✅ Downloaded BetterDiscord version %s from the official website\n", output.FormatVersion(version))
		}
		i.hasDownloaded = true
		return nil
	} else {
		output.Println("❌ Failed to download BetterDiscord from official website")
		output.Printf("❌ %s\n", err.Error())
		output.Blank()
		output.Println("🔁 Falling back to GitHub...")
	}

	// Get download URL from GitHub API
	apiData, err := utils.DownloadJSON[models.GitHubRelease](githubLatestReleaseURL)
	if err != nil {
		output.Println("❌ Failed to get asset url from GitHub")
		output.Printf("❌ %s\n", err.Error())
		return err
	}

	var index = -1
	for idx, asset := range apiData.Assets {
		if asset.Name == "betterdiscord.asar" {
			index = idx
			break
		}
	}

	if index == -1 {
		output.Println("❌ Failed to find the BetterDiscord asar on GitHub")
		return fmt.Errorf("failed to find betterdiscord.asar asset in GitHub release")
	}

	var downloadUrl = apiData.Assets[index].URL
	var version = apiData.TagName

	if downloadUrl != "" {
		output.Printf("✅ Found BetterDiscord: %s\n", downloadUrl)
	}

	// Download asar into the BD folder
	_, err = utils.DownloadFile(downloadUrl, i.asar)
	if err != nil {
		output.Println("❌ Failed to download BetterDiscord from GitHub")
		output.Printf("❌ %s\n", err.Error())
		return err
	}

	if version == "" {
		output.Println("✅ Downloaded BetterDiscord from GitHub")
	} else {
		output.Printf("✅ Downloaded BetterDiscord version %s from GitHub\n", output.FormatVersion(version))
	}
	i.hasDownloaded = true

	return nil
}
