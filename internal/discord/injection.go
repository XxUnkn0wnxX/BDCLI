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
// operation, by creating and removing a unique throwaway file. This is the
// elevation trigger: a failure here means we abort before touching the bundle.
// A unique name (os.CreateTemp) avoids colliding with or clobbering an existing
// file and is safe under concurrent probes.
func probeWritable(dir string) error {
	f, err := os.CreateTemp(dir, ".bd-write-probe-*")
	if err != nil {
		return err
	}
	// Writability is already proven by the successful create; cleanup is
	// best-effort and must not turn a writable dir into a probe failure.
	_ = f.Close()
	_ = os.Remove(f.Name())
	return nil
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
	if resources == "" {
		return fmt.Errorf("cannot inject: resources path is empty")
	}
	// Snap installs live on a read-only squashfs mount that can't host the shadow.
	// Short-circuit with an actionable message instead of surfacing the generic
	// permission error the writability probe would raise below.
	if discord.IsSnap {
		output.Printf("❌ Snap installs are not supported\n")
		output.Printf("   The read-only Snap mount cannot host the BetterDiscord injection.\n")
		return fmt.Errorf("snap installs are not supported")
	}
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
			output.Printf("❌ Unable to modify app.asar in %s\n", resources)
			output.Printf("   Discord may still be running, please fully close it and try again.\n")
			output.Printf("   %s\n", err.Error())
			return err
		}
	case utils.Exists(preservedAsar):
		// Already preserved from a prior injection and no live app.asar; the
		// archive is correct — only the shadow app/ needs (re)writing below.
	default:
		return fmt.Errorf("no app.asar found in %s", resources)
	}

	// Roll back anything done after this point so a partial failure never leaves
	// Discord without a loadable app. The restore is keyed on filesystem state,
	// not on whether *this* call renamed: rollback's own RemoveAll(appDir) clears
	// app/ even when re-injecting an already-injected install, so we must still
	// restore app.asar from the preserved copy to keep Discord launchable.
	rollback := func() {
		os.RemoveAll(appDir)
		if !utils.Exists(originalAsar) && utils.Exists(preservedAsar) {
			if err := os.Rename(preservedAsar, originalAsar); err != nil {
				output.Printf("❌ Rollback failed: unable to restore app.asar in %s\n", resources)
				output.Printf("   %s\n", err.Error())
			}
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
	if resources == "" {
		return fmt.Errorf("cannot uninject: resources path is empty")
	}
	// Snap installs are never injectable (read-only mount), so there is nothing to
	// revert; report it explicitly rather than attempting filesystem mutations.
	if discord.IsSnap {
		output.Printf("❌ Snap installs are not supported\n")
		output.Printf("   The read-only Snap mount cannot host the BetterDiscord injection.\n")
		return fmt.Errorf("snap installs are not supported")
	}
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
	switch {
	case utils.Exists(preservedAsar) && !utils.Exists(originalAsar):
		// Normal revert: restore Discord's original app from the preserved copy.
		if err := os.Rename(preservedAsar, originalAsar); err != nil {
			output.Printf("❌ Unable to restore app.asar in %s\n", resources)
			output.Printf("   Discord may still be running, please fully close it and try again.\n")
			output.Printf("   %s\n", err.Error())
			return err
		}
	case utils.Exists(preservedAsar):
		// A live app.asar is already present (e.g. Discord repaired/reinstalled
		// over the injection), so the preserved copy is stale. Remove it to fully
		// revert and reclaim the space (100MB+). A failure here only leaves a
		// harmless leftover — Discord still launches — so don't fail the uninstall.
		if err := os.Remove(preservedAsar); err != nil {
			output.Printf("⚠️  Unable to remove stale %s\n", preservedAsar)
			output.Printf("   %s\n", err.Error())
		}
	}

	output.Printf("✅ Removed from %s\n", discord.Channel.Name())
	return nil
}

// IsInjected reports whether this install currently has the app.asar shadow in
// place: both our `app/index.js` entry and the preserved original must exist.
func (discord *DiscordInstall) IsInjected() bool {
	resources := discord.ResourcesPath
	if resources == "" {
		return false
	}
	return utils.Exists(filepath.Join(resources, "app", "index.js")) && utils.Exists(filepath.Join(resources, "betterdiscord.app.asar"))
}