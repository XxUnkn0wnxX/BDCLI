package discord

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/betterdiscord/cli/internal/models"
)

func TestLegacyDiscovery_DirectAndWrappedLayouts(t *testing.T) {
	configRoot := t.TempDir()
	direct := writeLegacyCandidate(t, configRoot, models.Canary, "0.0.2", false, legacyInjectionScript)
	wrapped := writeLegacyCandidate(t, configRoot, models.Canary, "0.0.10", true, legacyInjectionScript)
	writeLegacyCandidate(t, configRoot, models.Canary, "not-a-version", false, legacyInjectionScript)

	states, err := findLegacyInjectionStates(configRoot, models.Canary)
	if err != nil {
		t.Fatalf("findLegacyInjectionStates() failed: %v", err)
	}
	if len(states) != 2 {
		t.Fatalf("findLegacyInjectionStates() returned %d states, want 2", len(states))
	}
	if states[0].indexPath != wrapped || states[1].indexPath != direct {
		t.Fatalf("states ordered as %q, %q; want %q, %q", states[0].indexPath, states[1].indexPath, wrapped, direct)
	}
}

func TestLegacyDiscovery_AcceptsAppVersionDirectory(t *testing.T) {
	configRoot := t.TempDir()
	indexPath := writeLegacyCandidate(t, configRoot, models.Stable, "app-0.0.90", false, legacyInjectionScript)

	states, err := findLegacyInjectionStates(configRoot, models.Stable)
	if err != nil {
		t.Fatalf("findLegacyInjectionStates() failed: %v", err)
	}
	if len(states) != 1 || states[0].indexPath != indexPath {
		t.Fatalf("findLegacyInjectionStates() = %#v, want one app-version candidate", states)
	}
}

func TestLegacyCleanupAndRestoration(t *testing.T) {
	configRoot := t.TempDir()
	indexPath := writeLegacyCandidate(t, configRoot, models.Stable, "0.0.1", false, legacyInjectionScript)
	const originalMode = 0o640
	if err := os.Chmod(indexPath, originalMode); err != nil {
		t.Fatalf("chmod() failed: %v", err)
	}

	states, err := findLegacyInjectionStates(configRoot, models.Stable)
	if err != nil {
		t.Fatalf("findLegacyInjectionStates() failed: %v", err)
	}
	if err := cleanupLegacyInjectionStates(states); err != nil {
		t.Fatalf("cleanupLegacyInjectionStates() failed: %v", err)
	}
	assertFile(t, indexPath, []byte(canonicalLegacyLoader), originalMode)

	if err := restoreLegacyInjectionStates(states); err != nil {
		t.Fatalf("restoreLegacyInjectionStates() failed: %v", err)
	}
	assertFile(t, indexPath, legacyInjectionScript, originalMode)
}

func TestLegacyDiscovery_IgnoresMalformedMissingAndSymlinkCandidates(t *testing.T) {
	configRoot := t.TempDir()
	channelRoot := filepath.Join(configRoot, "discordptb")
	if err := os.MkdirAll(channelRoot, 0o755); err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}

	// Missing index, missing core, and an invalid module layout are all ignored.
	writeLegacyCore(t, channelRoot, "0.0.1", false)
	writeLegacyIndexOnly(t, channelRoot, "0.0.2", false, legacyInjectionScript)
	writeLegacyCore(t, channelRoot, "0.0.3", true)
	customPath := writeLegacyCandidate(t, configRoot, models.PTB, "0.0.6", false, []byte("module.exports = require('./core.asar'); // custom\n"))

	validTargetDir := t.TempDir()
	validTarget := filepath.Join(validTargetDir, "index.js")
	if err := os.WriteFile(validTarget, legacyInjectionScript, 0o644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(validTargetDir, "core.asar"), []byte("core"), 0o644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	symlinkDir := filepath.Join(channelRoot, "0.0.4", "modules", "discord_desktop_core")
	if err := os.MkdirAll(filepath.Dir(symlinkDir), 0o755); err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}
	if err := os.Symlink(filepath.Dir(validTarget), symlinkDir); err != nil {
		t.Fatalf("Symlink() failed: %v", err)
	}

	indexSymlinkDir := filepath.Join(channelRoot, "0.0.5", "modules", "discord_desktop_core")
	if err := os.MkdirAll(indexSymlinkDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(indexSymlinkDir, "core.asar"), []byte("core"), 0o644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	if err := os.Symlink(validTarget, filepath.Join(indexSymlinkDir, "index.js")); err != nil {
		t.Fatalf("Symlink() failed: %v", err)
	}

	coreSymlinkDir := filepath.Join(channelRoot, "0.0.7", "modules", "discord_desktop_core")
	if err := os.MkdirAll(coreSymlinkDir, 0o755); err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}
	if err := os.Symlink(filepath.Join(validTargetDir, "core.asar"), filepath.Join(coreSymlinkDir, "core.asar")); err != nil {
		t.Fatalf("Symlink() failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(coreSymlinkDir, "index.js"), legacyInjectionScript, 0o644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}

	states, err := findLegacyInjectionStates(configRoot, models.PTB)
	if err != nil {
		t.Fatalf("findLegacyInjectionStates() failed: %v", err)
	}
	if len(states) != 0 {
		t.Fatalf("findLegacyInjectionStates() returned %d states, want none", len(states))
	}
	assertFile(t, customPath, []byte("module.exports = require('./core.asar'); // custom\n"), 0o644)
}

func TestLegacyPreflight_UnknownBetterDiscordLoaderBlocksWithoutWrites(t *testing.T) {
	configRoot := t.TempDir()
	exactPath := writeLegacyCandidate(t, configRoot, models.Stable, "0.0.1", false, legacyInjectionScript)
	unknown := []byte("// BetterDiscord legacy loader variant\nmodule.exports = require('./core.asar');\n")
	unknownPath := writeLegacyCandidate(t, configRoot, models.Stable, "0.0.2", false, unknown)

	states, err := findLegacyInjectionStates(configRoot, models.Stable)
	if err == nil {
		t.Fatal("findLegacyInjectionStates() succeeded for an unknown BetterDiscord loader")
	}
	if !strings.Contains(err.Error(), "unknown BetterDiscord legacy loader") {
		t.Fatalf("error = %q, want actionable unknown-loader message", err)
	}
	if states != nil {
		t.Fatalf("states = %#v, want nil when preflight blocks", states)
	}
	assertFile(t, exactPath, legacyInjectionScript, 0o644)
	assertFile(t, unknownPath, unknown, 0o644)
}

func TestLegacyCleanup_RollsBackEarlierWrites(t *testing.T) {
	configRoot := t.TempDir()
	firstPath := writeLegacyCandidate(t, configRoot, models.Stable, "0.0.1", false, legacyInjectionScript)
	secondPath := writeLegacyCandidate(t, configRoot, models.Stable, "0.0.2", false, legacyInjectionScript)
	states, err := findLegacyInjectionStates(configRoot, models.Stable)
	if err != nil {
		t.Fatalf("findLegacyInjectionStates() failed: %v", err)
	}
	if err := os.Remove(secondPath); err != nil {
		t.Fatalf("Remove() failed: %v", err)
	}

	err = cleanupLegacyInjectionStates(states)
	if err == nil {
		t.Fatal("cleanupLegacyInjectionStates() succeeded after target removal")
	}
	assertFile(t, firstPath, legacyInjectionScript, 0o644)
}

func TestLegacyInstallActionFailureRestoresLegacyBytes(t *testing.T) {
	configRoot := t.TempDir()
	indexPath := writeLegacyCandidate(t, configRoot, models.Stable, "0.0.1", false, legacyInjectionScript)
	states, err := findLegacyInjectionStates(configRoot, models.Stable)
	if err != nil {
		t.Fatalf("findLegacyInjectionStates() failed: %v", err)
	}
	actionErr := errors.New("inject failed")

	err = runLegacyInstallAction(states, func() error { return actionErr })
	if !errors.Is(err, actionErr) {
		t.Fatalf("runLegacyInstallAction() error = %v, want %v", err, actionErr)
	}
	assertFile(t, indexPath, legacyInjectionScript, 0o644)
}

func TestLegacyInstallActionSuccessLeavesCanonical(t *testing.T) {
	configRoot := t.TempDir()
	indexPath := writeLegacyCandidate(t, configRoot, models.Stable, "0.0.1", false, legacyInjectionScript)
	states, err := findLegacyInjectionStates(configRoot, models.Stable)
	if err != nil {
		t.Fatalf("findLegacyInjectionStates() failed: %v", err)
	}

	if err := runLegacyInstallAction(states, func() error { return nil }); err != nil {
		t.Fatalf("runLegacyInstallAction() failed: %v", err)
	}
	assertFile(t, indexPath, []byte(canonicalLegacyLoader), 0o644)
}

func TestLegacyRemovalActionFailureLeavesCanonical(t *testing.T) {
	configRoot := t.TempDir()
	indexPath := writeLegacyCandidate(t, configRoot, models.Stable, "0.0.1", false, legacyInjectionScript)
	states, err := findLegacyInjectionStates(configRoot, models.Stable)
	if err != nil {
		t.Fatalf("findLegacyInjectionStates() failed: %v", err)
	}
	actionErr := errors.New("uninject failed")

	err = runLegacyRemovalAction(states, func() error { return actionErr })
	if !errors.Is(err, actionErr) {
		t.Fatalf("runLegacyRemovalAction() error = %v, want %v", err, actionErr)
	}
	assertFile(t, indexPath, []byte(canonicalLegacyLoader), 0o644)
}

func writeLegacyCandidate(t *testing.T, configRoot string, channel models.DiscordChannel, version string, wrapped bool, index []byte) string {
	t.Helper()
	module := "discord_desktop_core"
	if wrapped {
		module += "-1"
	}
	moduleRoot := filepath.Join(configRoot, legacyChannelName(channel), version, "modules", module)
	if wrapped {
		moduleRoot = filepath.Join(moduleRoot, "discord_desktop_core")
	}
	if err := os.MkdirAll(moduleRoot, 0o755); err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(moduleRoot, "core.asar"), []byte("core"), 0o644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	indexPath := filepath.Join(moduleRoot, "index.js")
	if err := os.WriteFile(indexPath, index, 0o644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
	return indexPath
}

func writeLegacyCore(t *testing.T, channelRoot, version string, wrapped bool) {
	t.Helper()
	module := "discord_desktop_core"
	if wrapped {
		module += "-1"
	}
	moduleRoot := filepath.Join(channelRoot, version, "modules", module)
	if err := os.MkdirAll(moduleRoot, 0o755); err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(moduleRoot, "core.asar"), []byte("core"), 0o644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
}

func writeLegacyIndexOnly(t *testing.T, channelRoot, version string, wrapped bool, index []byte) {
	t.Helper()
	module := "discord_desktop_core"
	if wrapped {
		module += "-1"
	}
	moduleRoot := filepath.Join(channelRoot, version, "modules", module)
	if wrapped {
		moduleRoot = filepath.Join(moduleRoot, "discord_desktop_core")
	}
	if err := os.MkdirAll(moduleRoot, 0o755); err != nil {
		t.Fatalf("MkdirAll() failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(moduleRoot, "index.js"), index, 0o644); err != nil {
		t.Fatalf("WriteFile() failed: %v", err)
	}
}

func assertFile(t *testing.T, path string, want []byte, wantMode os.FileMode) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile(%q) failed: %v", path, err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("%s = %q, want %q", path, got, want)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat(%q) failed: %v", path, err)
	}
	if info.Mode().Perm() != wantMode {
		t.Errorf("%s mode = %o, want %o", path, info.Mode().Perm(), wantMode)
	}
}
