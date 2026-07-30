package discord

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/betterdiscord/cli/internal/models"
)

// UninstallBD with neither full-uninstall nor restart should only de-inject the
// core's index.js — the safe path that never touches the global BD folder or
// the running Discord process.
func TestUninstallBD_UninjectOnly(t *testing.T) {
	corePath := t.TempDir()
	indexFile := filepath.Join(corePath, "index.js")
	seed := `require("BetterDiscord/data/betterdiscord.asar");` + "\n" + `module.exports = require("./core.asar");`
	if err := os.WriteFile(indexFile, []byte(seed), 0o644); err != nil {
		t.Fatalf("failed to seed injected index.js: %v", err)
	}

	install := &DiscordInstall{ResourcesPath: corePath, Channel: models.Stable}
	if err := install.UninstallBD(models.UninstallOptions{FullUninstall: false, RestartDiscord: false}); err != nil {
		t.Fatalf("UninstallBD() failed: %v", err)
	}

	if install.IsInjected() {
		t.Error("expected index.js to be de-injected after UninstallBD")
	}
	contents, _ := os.ReadFile(indexFile)
	if want := `module.exports = require("./core.asar");`; string(contents) != want {
		t.Errorf("index.js after uninstall = %q, expected %q", string(contents), want)
	}
}

func TestGetBetterDiscordInstall_Global(t *testing.T) {
	install := &DiscordInstall{ResourcesPath: "/some/discord/core", Channel: models.Stable}

	bd, err := install.GetBetterDiscordInstall()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bd == nil {
		t.Fatal("expected a non-nil global BD install")
	}
}

func TestGetBetterDiscordInstall_FlatpakResolvesConfig(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses POSIX-style flatpak paths")
	}
	// A flatpak-style core path containing a "config" segment.
	configDir := filepath.Join(t.TempDir(), "config")
	corePath := filepath.Join(configDir, "discord", "0.0.1", "modules", "discord_desktop_core")
	install := &DiscordInstall{ResourcesPath: corePath, Channel: models.Stable, IsFlatpak: true}

	bd, err := install.GetBetterDiscordInstall()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bd == nil {
		t.Fatal("expected non-nil BD install")
	}
	if want := filepath.Join(configDir, "BetterDiscord"); bd.Root() != want {
		t.Errorf("Root() = %s, expected %s", bd.Root(), want)
	}
}

func TestGetBetterDiscordInstall_SnapResolvesConfig(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses POSIX-style snap paths")
	}
	// A snap-style core path uses the ".config" segment.
	configDir := filepath.Join(t.TempDir(), ".config")
	corePath := filepath.Join(configDir, "discord", "0.0.1", "modules", "discord_desktop_core")
	install := &DiscordInstall{ResourcesPath: corePath, Channel: models.Stable, IsSnap: true}

	bd, err := install.GetBetterDiscordInstall()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := filepath.Join(configDir, "BetterDiscord"); bd.Root() != want {
		t.Errorf("Root() = %s, expected %s", bd.Root(), want)
	}
}

// Regression test for the nil-deref fix: a snap/flatpak core path missing the
// expected config segment must surface an error, not return a nil *BDInstall
// that callers would dereference and panic on.
func TestGetBetterDiscordInstall_FlatpakMissingSegment_Errors(t *testing.T) {
	install := &DiscordInstall{ResourcesPath: "/no/matching/segment/here", Channel: models.Stable, IsFlatpak: true}

	bd, err := install.GetBetterDiscordInstall()
	if err == nil {
		t.Fatal("expected an error when the config segment is missing")
	}
	if bd != nil {
		t.Errorf("expected nil BD install on error, got %+v", bd)
	}
}