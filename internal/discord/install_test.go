package discord

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/betterdiscord/cli/internal/models"
)

// UninstallBD with neither full-uninstall nor restart reverts the app.asar
// shadow without removing the global BD folder or relaunching Discord. (Discord
// isn't running in the test, so the stop() step is a no-op.)
func TestUninstallBD_UninjectOnly(t *testing.T) {
	resources := t.TempDir()
	// Seed an injected state: preserved asar + shadow app/ entry.
	original := []byte("original app.asar")
	if err := os.WriteFile(filepath.Join(resources, "betterdiscord.app.asar"), original, 0o644); err != nil {
		t.Fatalf("seed preserved asar: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(resources, "app"), 0o755); err != nil {
		t.Fatalf("mkdir app: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resources, "app", "index.js"), []byte("x"), 0o644); err != nil {
		t.Fatalf("seed index.js: %v", err)
	}

	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	if err := install.UninstallBD(models.UninstallOptions{FullUninstall: false, RestartDiscord: false}); err != nil {
		t.Fatalf("UninstallBD() failed: %v", err)
	}

	if install.IsInjected() {
		t.Error("expected the shadow to be reverted after UninstallBD")
	}
	restored, err := os.ReadFile(filepath.Join(resources, "app.asar"))
	if err != nil {
		t.Fatalf("app.asar not restored: %v", err)
	}
	if string(restored) != string(original) {
		t.Errorf("app.asar after uninstall = %q, expected %q", restored, original)
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

// Flatpak's BD folder is recomputed as ~/.var/app/{id}/config/BetterDiscord from
// the channel, independent of the (read-only deployment) resources path.
func TestGetBetterDiscordInstall_FlatpakRecomputesDataRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses POSIX-style flatpak paths")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("no home dir: %v", err)
	}

	cases := []struct {
		channel models.DiscordChannel
		id      string
	}{
		{models.Stable, "com.discordapp.Discord"},
		{models.Canary, "com.discordapp.DiscordCanary"},
		{models.PTB, "com.discordapp.DiscordPTB"},
	}
	for _, tc := range cases {
		// A resources path in the read-only deployment tree (no "config" segment).
		resources := "/var/lib/flatpak/app/" + tc.id + "/current/active/files/discord/resources"
		install := &DiscordInstall{ResourcesPath: resources, Channel: tc.channel, IsFlatpak: true}

		bd, err := install.GetBetterDiscordInstall()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := filepath.Join(home, ".var", "app", tc.id, "config", "BetterDiscord")
		if bd.Root() != want {
			t.Errorf("channel %v: Root() = %s, expected %s", tc.channel, bd.Root(), want)
		}
	}
}