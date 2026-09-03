package discord

import (
	"bytes"
	_ "embed"
	"errors"
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

// errIfSnap rejects Snap installs with a clear, actionable message: their
// read-only squashfs mount can't host the app.asar shadow. It's called at the
// start of the install/uninstall/repair flows (before Discord is stopped, so an
// unsupported install never needlessly kills a running client) and again in
// inject/uninject as a backstop for any direct callers.
func (discord *DiscordInstall) errIfSnap() error {
	if !discord.IsSnap {
		return nil
	}
	output.Printf("❌ Snap installs are not supported\n")
	output.Printf("   The read-only Snap mount cannot host the BetterDiscord injection.\n")
	return fmt.Errorf("snap installs are not supported")
}

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

// validateAppShadow checks an existing app/ before any injection mutation. An
// existing directory is safe to reuse only when it is exactly the shadow this
// CLI owns; arbitrary app state must never be removed by rollback.
func validateAppShadow(appDir, preservedAsar string) error {
	info, err := os.Lstat(appDir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("cannot inspect existing app shadow %s: %w", appDir, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("refusing to use %s: existing app shadow is a symlink", appDir)
	}
	if !info.IsDir() {
		return fmt.Errorf("refusing to use %s: existing app shadow is not a directory", appDir)
	}

	entries, err := os.ReadDir(appDir)
	if err != nil {
		return fmt.Errorf("cannot inspect existing app shadow %s: %w", appDir, err)
	}
	expected := map[string][]byte{
		"index.js":     []byte(appIndexScript),
		"package.json": []byte(appPackageJSON),
	}
	if len(entries) != len(expected) {
		return fmt.Errorf("refusing to use %s: existing app shadow must contain only index.js and package.json", appDir)
	}

	assetsMatch := true
	for _, entry := range entries {
		want, ok := expected[entry.Name()]
		if !ok {
			return fmt.Errorf("refusing to use %s: existing app shadow contains unexpected entry %q", appDir, entry.Name())
		}

		entryPath := filepath.Join(appDir, entry.Name())
		entryInfo, err := os.Lstat(entryPath)
		if err != nil {
			return fmt.Errorf("cannot inspect existing app shadow entry %s: %w", entryPath, err)
		}
		if entryInfo.Mode()&os.ModeSymlink != 0 || !entryInfo.Mode().IsRegular() {
			return fmt.Errorf("refusing to use %s: existing app shadow entry %q is not a regular file", appDir, entry.Name())
		}
		got, err := os.ReadFile(entryPath)
		if err != nil {
			return fmt.Errorf("cannot read existing app shadow entry %s: %w", entryPath, err)
		}
		if !bytes.Equal(got, want) {
			assetsMatch = false
		}
	}

	if assetsMatch || isRegularNonSymlinkPath(preservedAsar) {
		return nil
	}
	return fmt.Errorf("refusing to use %s: existing app shadow assets are altered and no preserved BetterDiscord app archive proves prior ownership", appDir)
}

func isRegularNonSymlinkPath(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular() && info.Mode()&os.ModeSymlink == 0
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

	// Backstop: the install flow rejects Snap before stopping Discord, but guard
	// here too so any direct caller gets the same actionable error rather than the
	// generic permission failure the writability probe would raise below.
	if err := discord.errIfSnap(); err != nil {
		return err
	}

	originalAsar := filepath.Join(resources, "app.asar")
	preservedAsar := filepath.Join(resources, "betterdiscord.app.asar")
	appDir := filepath.Join(resources, "app")

	// Rollback removes app/, so never allow an arbitrary pre-existing app path
	// into the transaction. Managed shadows are the only safe re-injection state.
	if err := validateAppShadow(appDir, preservedAsar); err != nil {
		output.Printf("❌ Cannot safely inject into %s\n", resources)
		output.Printf("   %s\n", err.Error())
		return err
	}

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
		err := os.RemoveAll(appDir)
		if err != nil {
			output.Printf("❌ Rollback failed: unable to remove %s\n", appDir)
			output.Printf("   %s\n", err.Error())
		}
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

	// Backstop for direct callers; the uninstall/repair flows reject Snap before
	// stopping Discord. Snap installs are never injectable, so there's nothing to
	// revert — report it explicitly rather than attempting filesystem mutations.
	if err := discord.errIfSnap(); err != nil {
		return err
	}

	originalAsar := filepath.Join(resources, "app.asar")
	preservedAsar := filepath.Join(resources, "betterdiscord.app.asar")
	appDir := filepath.Join(resources, "app")

	// Removal clears app/, so never allow an arbitrary pre-existing app path to
	// enter the uninject transaction. Managed shadows are safe to remove.
	if err := validateAppShadow(appDir, preservedAsar); err != nil {
		output.Printf("❌ Cannot safely uninject from %s\n", resources)
		output.Printf("   %s\n", err.Error())
		return err
	}

	// A clean install (only app.asar; no shadow app/ and no preserved copy) was
	// never injected — report a no-op instead of claiming a removal that didn't
	// happen, which would mislead anyone troubleshooting uninstall/repair.
	if !utils.Exists(appDir) && !utils.Exists(preservedAsar) {
		output.Printf("ℹ️  No injection found in %s\n", discord.Channel.Name())
		return nil
	}

	// Restore Discord's original app.asar *before* removing the shadow app/. If the
	// restore fails (e.g. a running Discord still locks the file), the injection is
	// left fully intact and loadable rather than bricked with neither app.asar nor
	// app/ present.
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

	// With both asars gone (a corrupted state discovery normally filters out),
	// there is nothing to restore. The shadow app/ can't boot either (its entry
	// point requires the missing preserved asar) but deleting it would erase the
	// evidence and report success on an unlaunchable install. Refuse so the user
	// learns Discord itself needs a reinstall.
	if !utils.Exists(originalAsar) && !utils.Exists(preservedAsar) {
		output.Printf("❌ Unable to uninject %s: missing both app.asar and betterdiscord.app.asar\n", discord.Channel.Name())
		output.Printf("   Your %s install appears corrupted! Reinstall Discord, then run this command again.\n", discord.Channel.Name())
		return fmt.Errorf("cannot uninject: missing both app.asar and betterdiscord.app.asar in %s", resources)
	}

	// An app.asar is in place (restored above, or already live); now clear the
	// shadow app/. A failure here is non-bricking since Electron prefers app.asar
	// over app/, but still surface it so the leftover can be cleaned up.
	if utils.Exists(appDir) {
		if err := os.RemoveAll(appDir); err != nil {
			output.Printf("❌ Unable to remove %s\n", appDir)
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
	if resources == "" {
		return false
	}
	return utils.Exists(filepath.Join(resources, "app", "index.js")) && utils.Exists(filepath.Join(resources, "betterdiscord.app.asar"))
}
