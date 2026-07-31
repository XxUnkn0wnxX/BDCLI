package discord

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/betterdiscord/cli/internal/betterdiscord"
	"github.com/betterdiscord/cli/internal/output"
	"github.com/betterdiscord/cli/internal/utils"
)

//go:embed assets/app_index.js
var appIndexScript string

//go:embed assets/app_package.json
var appPackageJSON string

// probeWritable verifies dir accepts writes before we perform any destructive
// operation, by creating and removing a throwaway file. This is the elevation
// trigger: a failure here means we abort before touching the bundle.
func probeWritable(dir string) error {
	probe := filepath.Join(dir, ".bd-write-probe")
	if err := os.WriteFile(probe, []byte{}, 0o644); err != nil {
		return err
	}
	return os.Remove(probe)
}

// inject shadows Discord's app.asar: it preserves the original as
// betterdiscord.app.asar and drops an `app/` entry directory that loads
// BetterDiscord and then the preserved app. The operation is transactional —
// any failure after the rename rolls back to the original state.
//
// bd is accepted for call-site symmetry with the install flow but unused: the
// injection script resolves the BetterDiscord folder at runtime.
func (discord *DiscordInstall) inject(bd *betterdiscord.BDInstall) error {
	resources := discord.ResourcesPath
	originalAsar := filepath.Join(resources, "app.asar")
	preservedAsar := filepath.Join(resources, "betterdiscord.app.asar")
	appDir := filepath.Join(resources, "app")

	// Probe writability before the destructive rename so we never leave a
	// half-modified bundle on a read-only/permission-denied target.
	if err := probeWritable(resources); err != nil {
		output.Printf("❌ Cannot write to %s\n", resources)
		output.Printf("   %s\n", err.Error())
		return err
	}

	// Preserve the original app.asar (idempotent, guarded).
	//
	// A live app.asar is always Discord's current app and takes priority: it
	// must be renamed away or it would shadow our app/ folder (Electron loads
	// app.asar before app/), silently disabling BetterDiscord. If a preserved
	// copy is also present — e.g. Discord repaired/reinstalled over a previous
	// injection — that copy is stale, so we discard it and re-preserve the live
	// app. Only when there is no live app.asar do we treat an existing preserved
	// copy as the (already-injected) source of truth and leave it be.
	renamed := false
	switch {
	case utils.Exists(originalAsar):
		if utils.Exists(preservedAsar) {
			if err := os.Remove(preservedAsar); err != nil {
				output.Printf("❌ Unable to replace stale %s\n", preservedAsar)
				output.Printf("   %s\n", err.Error())
				return err
			}
		}
		if err := os.Rename(originalAsar, preservedAsar); err != nil {
			output.Printf("❌ Unable to preserve app.asar in %s\n", resources)
			output.Printf("   %s\n", err.Error())
			return err
		}
		renamed = true
	case utils.Exists(preservedAsar):
		// Already preserved from a prior injection and no live app.asar; the
		// archive is correct — only the shadow app/ needs (re)writing below.
	default:
		return fmt.Errorf("no app.asar found in %s", resources)
	}

	// Roll back anything done after the rename so a partial failure never
	// leaves Discord without a loadable app.
	rollback := func() {
		os.RemoveAll(appDir)
		if renamed && !utils.Exists(originalAsar) {
			os.Rename(preservedAsar, originalAsar)
		}
	}

	if err := os.MkdirAll(appDir, 0755); err != nil {
		output.Printf("❌ Unable to create %s\n", appDir)
		output.Printf("   %s\n", err.Error())
		rollback()
		return err
	}

	if err := os.WriteFile(filepath.Join(appDir, "package.json"), []byte(appPackageJSON), 0o644); err != nil {
		output.Printf("❌ Unable to write package.json in %s\n", appDir)
		output.Printf("   %s\n", err.Error())
		rollback()
		return err
	}

	if err := os.WriteFile(filepath.Join(appDir, "index.js"), []byte(appIndexScript), 0o644); err != nil {
		output.Printf("❌ Unable to write index.js in %s\n", appDir)
		output.Printf("   %s\n", err.Error())
		rollback()
		return err
	}

	if !utils.Exists(filepath.Join(appDir, "index.js")) ||
		!utils.Exists(filepath.Join(appDir, "package.json")) ||
		!utils.Exists(preservedAsar) {
		rollback()
		return fmt.Errorf("injection verification failed in %s", resources)
	}

	output.Printf("✅ Injected into %s\n", resources)
	return nil
}

// uninject reverses inject: it removes the shadow `app/` directory and restores
// Discord's original app.asar from the preserved copy.
func (discord *DiscordInstall) uninject() error {
	resources := discord.ResourcesPath
	originalAsar := filepath.Join(resources, "app.asar")
	preservedAsar := filepath.Join(resources, "betterdiscord.app.asar")
	appDir := filepath.Join(resources, "app")

	if utils.Exists(appDir) {
		if err := os.RemoveAll(appDir); err != nil {
			output.Printf("❌ Unable to remove %s\n", appDir)
			output.Printf("   %s\n", err.Error())
			return err
		}
	}

	// Only restore when a preserved copy exists and we wouldn't clobber a live
	// app.asar (crash-recovery / partial-state safety).
	if utils.Exists(preservedAsar) && !utils.Exists(originalAsar) {
		if err := os.Rename(preservedAsar, originalAsar); err != nil {
			output.Printf("❌ Unable to restore app.asar in %s\n", resources)
			output.Printf("   %s\n", err.Error())
			return err
		}
	}

	output.Printf("✅ Removed from %s\n", discord.Channel.Name())
	return nil
}

// IsInjected reports whether this install currently has the app.asar shadow in
// place: both our `app/index.js` entry and the preserved original must exist.
func (discord *DiscordInstall) IsInjected() bool {
	resources := discord.ResourcesPath
	return utils.Exists(filepath.Join(resources, "app", "index.js")) &&
		utils.Exists(filepath.Join(resources, "betterdiscord.app.asar"))
}