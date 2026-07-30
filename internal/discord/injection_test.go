package discord

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/betterdiscord/cli/internal/models"
)

const defaultIndexJS = `module.exports = require("./core.asar");`

func TestIsInjected(t *testing.T) {
	tmpDir := t.TempDir()
	resourcesPath := filepath.Join(tmpDir, "discord_desktop_core")
	if err := os.MkdirAll(resourcesPath, 0755); err != nil {
		t.Fatalf("Failed to create core path: %v", err)
	}
	indexFile := filepath.Join(resourcesPath, "index.js")

	install := &DiscordInstall{
		ResourcesPath: resourcesPath,
		Channel:  models.Stable,
	}

	if install.IsInjected() {
		t.Fatalf("Expected IsInjected to be false with missing index.js")
	}

	if err := os.WriteFile(indexFile, []byte(`module.exports = require("./core.asar");`), 0644); err != nil {
		t.Fatalf("Failed to write index.js: %v", err)
	}
	if install.IsInjected() {
		t.Fatalf("Expected IsInjected to be false for default index.js")
	}

	if err := os.WriteFile(indexFile, []byte(`// BetterDiscord injected`), 0644); err != nil {
		t.Fatalf("Failed to write injection index.js: %v", err)
	}
	if !install.IsInjected() {
		t.Fatalf("Expected IsInjected to be true when BetterDiscord is present")
	}
}

func TestInject_WritesInjectionScript(t *testing.T) {
	resourcesPath := t.TempDir()
	install := &DiscordInstall{ResourcesPath: resourcesPath, Channel: models.Stable}

	// bd is unused by inject(); nil is fine.
	if err := install.inject(nil); err != nil {
		t.Fatalf("inject() failed: %v", err)
	}

	if !install.IsInjected() {
		t.Fatal("expected IsInjected() to be true after inject()")
	}

	info, err := os.Stat(filepath.Join(resourcesPath, "index.js"))
	if err != nil {
		t.Fatalf("index.js not written: %v", err)
	}
	// The require target must not carry the executable bit.
	if runtime.GOOS != "windows" {
		if perm := info.Mode().Perm(); perm != 0o644 {
			t.Errorf("index.js mode = %o, expected 644", perm)
		}
	}
}

func TestUninject_RemovesInjection(t *testing.T) {
	resourcesPath := t.TempDir()
	install := &DiscordInstall{ResourcesPath: resourcesPath, Channel: models.Stable}
	indexFile := filepath.Join(resourcesPath, "index.js")

	seed := `require("BetterDiscord/data/betterdiscord.asar");` + "\n" + defaultIndexJS
	if err := os.WriteFile(indexFile, []byte(seed), 0o644); err != nil {
		t.Fatalf("failed to seed injected index.js: %v", err)
	}

	if err := install.uninject(); err != nil {
		t.Fatalf("uninject() failed: %v", err)
	}

	contents, err := os.ReadFile(indexFile)
	if err != nil {
		t.Fatalf("index.js missing after uninject: %v", err)
	}
	if string(contents) != defaultIndexJS {
		t.Errorf("index.js after uninject = %q, expected %q", string(contents), defaultIndexJS)
	}
	if install.IsInjected() {
		t.Error("expected IsInjected() to be false after uninject()")
	}
}

func TestUninject_LeavesUninjectedFileUntouched(t *testing.T) {
	resourcesPath := t.TempDir()
	install := &DiscordInstall{ResourcesPath: resourcesPath, Channel: models.Stable}
	indexFile := filepath.Join(resourcesPath, "index.js")

	original := `module.exports = require("./some-other-core.asar");`
	if err := os.WriteFile(indexFile, []byte(original), 0o644); err != nil {
		t.Fatalf("failed to seed index.js: %v", err)
	}

	if err := install.uninject(); err != nil {
		t.Fatalf("uninject() failed: %v", err)
	}

	contents, _ := os.ReadFile(indexFile)
	if string(contents) != original {
		t.Errorf("uninject rewrote a file with no BetterDiscord marker: got %q", string(contents))
	}
}

func TestUninject_MissingFileWritesDefault(t *testing.T) {
	resourcesPath := t.TempDir()
	install := &DiscordInstall{ResourcesPath: resourcesPath, Channel: models.Stable}

	// No index.js exists; uninject falls through and writes the default stub.
	if err := install.uninject(); err != nil {
		t.Fatalf("uninject() failed: %v", err)
	}

	contents, err := os.ReadFile(filepath.Join(resourcesPath, "index.js"))
	if err != nil {
		t.Fatalf("expected index.js to be created: %v", err)
	}
	if string(contents) != defaultIndexJS {
		t.Errorf("index.js = %q, expected %q", string(contents), defaultIndexJS)
	}
}

func TestInjectUninject_RoundTrip(t *testing.T) {
	resourcesPath := t.TempDir()
	install := &DiscordInstall{ResourcesPath: resourcesPath, Channel: models.Stable}

	if err := install.inject(nil); err != nil {
		t.Fatalf("inject() failed: %v", err)
	}
	if !install.IsInjected() {
		t.Fatal("expected injected after inject()")
	}
	if err := install.uninject(); err != nil {
		t.Fatalf("uninject() failed: %v", err)
	}
	if install.IsInjected() {
		t.Fatal("expected not injected after uninject()")
	}
}