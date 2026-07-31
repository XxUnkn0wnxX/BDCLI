package discord

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/betterdiscord/cli/internal/models"
	"github.com/betterdiscord/cli/internal/utils"
)

// newResourcesDir creates a resources dir seeded with an app.asar of known
// content and returns the dir plus the original content.
func newResourcesDir(t *testing.T) (string, []byte) {
	t.Helper()
	resources := t.TempDir()
	content := []byte("original discord app.asar")
	if err := os.WriteFile(filepath.Join(resources, "app.asar"), content, 0o644); err != nil {
		t.Fatalf("failed to seed app.asar: %v", err)
	}
	return resources, content
}

func TestIsInjected(t *testing.T) {
	resources := t.TempDir()
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}

	if install.IsInjected() {
		t.Fatal("expected IsInjected false for a bare resources dir")
	}

	// Only the app/ entry, no preserved asar → not injected.
	if err := os.MkdirAll(filepath.Join(resources, "app"), 0o755); err != nil {
		t.Fatalf("mkdir app: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resources, "app", "index.js"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write index.js: %v", err)
	}
	if install.IsInjected() {
		t.Fatal("expected IsInjected false without a preserved app.asar")
	}

	// Add the preserved asar → injected.
	if err := os.WriteFile(filepath.Join(resources, "betterdiscord.app.asar"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write preserved asar: %v", err)
	}
	if !install.IsInjected() {
		t.Fatal("expected IsInjected true with app/index.js + preserved asar")
	}
}

func TestInject_Clean(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}

	if err := install.inject(nil); err != nil {
		t.Fatalf("inject() failed: %v", err)
	}

	if utils.Exists(filepath.Join(resources, "app.asar")) {
		t.Error("expected original app.asar to be renamed away")
	}
	preserved, err := os.ReadFile(filepath.Join(resources, "betterdiscord.app.asar"))
	if err != nil {
		t.Fatalf("preserved asar missing: %v", err)
	}
	if string(preserved) != string(original) {
		t.Errorf("preserved asar content = %q, expected %q", preserved, original)
	}
	if !utils.Exists(filepath.Join(resources, "app", "index.js")) {
		t.Error("app/index.js not written")
	}
	if !utils.Exists(filepath.Join(resources, "app", "package.json")) {
		t.Error("app/package.json not written")
	}
	if !install.IsInjected() {
		t.Error("expected IsInjected true after inject()")
	}

	// index.js must reference the preserved app and the BD asar.
	index, _ := os.ReadFile(filepath.Join(resources, "app", "index.js"))
	if want := "../betterdiscord.app.asar"; !strings.Contains(string(index), want) {
		t.Errorf("index.js missing %q", want)
	}
}

func TestInject_Idempotent(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}

	if err := install.inject(nil); err != nil {
		t.Fatalf("first inject() failed: %v", err)
	}
	// Corrupt the shadow index.js so we can confirm the second inject rewrites it
	// without re-renaming (which would clobber the real, already-preserved asar).
	if err := os.WriteFile(filepath.Join(resources, "app", "index.js"), []byte("stale"), 0o644); err != nil {
		t.Fatalf("corrupt index.js: %v", err)
	}

	if err := install.inject(nil); err != nil {
		t.Fatalf("second inject() failed: %v", err)
	}

	preserved, _ := os.ReadFile(filepath.Join(resources, "betterdiscord.app.asar"))
	if string(preserved) != string(original) {
		t.Errorf("preserved asar was clobbered on re-inject: got %q", preserved)
	}
	index, _ := os.ReadFile(filepath.Join(resources, "app", "index.js"))
	if string(index) == "stale" {
		t.Error("expected index.js to be rewritten on re-inject")
	}
	if utils.Exists(filepath.Join(resources, "app.asar")) {
		t.Error("re-inject must not recreate a live app.asar")
	}
}

// Anomalous pre-state: a live app.asar AND a leftover betterdiscord.app.asar +
// app/ (e.g. Discord repaired/reinstalled over a prior injection). inject() must
// treat the live app.asar as authoritative — discard the stale preserved copy,
// preserve the live app, and rename app.asar away so our app/ shadow loads
// (Electron would otherwise load the lingering app.asar and disable BD).
func TestInject_LiveAsarWinsOverStalePreserved(t *testing.T) {
	resources := t.TempDir()
	live := []byte("LIVE current app.asar")
	stale := []byte("stale old preserved app")
	if err := os.WriteFile(filepath.Join(resources, "app.asar"), live, 0o644); err != nil {
		t.Fatalf("seed live app.asar: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resources, "betterdiscord.app.asar"), stale, 0o644); err != nil {
		t.Fatalf("seed stale preserved: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(resources, "app"), 0o755); err != nil {
		t.Fatalf("seed leftover app/: %v", err)
	}

	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	if err := install.inject(nil); err != nil {
		t.Fatalf("inject: %v", err)
	}

	// app.asar must be renamed away so it can't shadow app/.
	if utils.Exists(filepath.Join(resources, "app.asar")) {
		t.Error("live app.asar should have been renamed away")
	}
	// The preserved copy must be the LIVE app, not the stale leftover.
	got, _ := os.ReadFile(filepath.Join(resources, "betterdiscord.app.asar"))
	if string(got) != string(live) {
		t.Errorf("preserved asar = %q, expected the live app %q", got, live)
	}
	if !install.IsInjected() {
		t.Error("expected IsInjected after re-injecting over a repaired install")
	}
}

func TestUninject_RestoresExactly(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}

	if err := install.inject(nil); err != nil {
		t.Fatalf("inject() failed: %v", err)
	}
	if err := install.uninject(); err != nil {
		t.Fatalf("uninject() failed: %v", err)
	}

	restored, err := os.ReadFile(filepath.Join(resources, "app.asar"))
	if err != nil {
		t.Fatalf("app.asar not restored: %v", err)
	}
	if string(restored) != string(original) {
		t.Errorf("restored app.asar = %q, expected %q", restored, original)
	}
	if utils.Exists(filepath.Join(resources, "betterdiscord.app.asar")) {
		t.Error("preserved asar should be gone after uninject")
	}
	if utils.Exists(filepath.Join(resources, "app")) {
		t.Error("shadow app/ should be removed after uninject")
	}
	if install.IsInjected() {
		t.Error("expected IsInjected false after uninject()")
	}
}

// If Discord repaired/reinstalled over an injection, uninject encounters a live
// app.asar alongside a now-stale betterdiscord.app.asar. It must remove the
// stale copy (reclaiming 100MB+) and leave the live app untouched.
func TestUninject_RemovesStalePreservedWhenLiveAsarPresent(t *testing.T) {
	resources := t.TempDir()
	if err := os.WriteFile(filepath.Join(resources, "app.asar"), []byte("live"), 0o644); err != nil {
		t.Fatalf("seed live app.asar: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resources, "betterdiscord.app.asar"), []byte("stale"), 0o644); err != nil {
		t.Fatalf("seed stale preserved: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(resources, "app"), 0o755); err != nil {
		t.Fatalf("seed app/: %v", err)
	}

	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	if err := install.uninject(); err != nil {
		t.Fatalf("uninject: %v", err)
	}

	if utils.Exists(filepath.Join(resources, "betterdiscord.app.asar")) {
		t.Error("stale preserved copy should be removed when a live app.asar exists")
	}
	got, _ := os.ReadFile(filepath.Join(resources, "app.asar"))
	if string(got) != "live" {
		t.Errorf("app.asar = %q, expected the untouched live app", got)
	}
	if utils.Exists(filepath.Join(resources, "app")) {
		t.Error("shadow app/ should be removed")
	}
	if install.IsInjected() {
		t.Error("should not report injected after uninject")
	}
}

func TestUninject_NotInjectedIsNoop(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}

	if err := install.uninject(); err != nil {
		t.Fatalf("uninject() on a clean install failed: %v", err)
	}

	// A never-injected install keeps its app.asar untouched.
	got, _ := os.ReadFile(filepath.Join(resources, "app.asar"))
	if string(got) != string(original) {
		t.Errorf("uninject touched a clean app.asar: got %q", got)
	}
}

func TestInject_NoAppAsarErrors(t *testing.T) {
	resources := t.TempDir() // empty, no app.asar
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}

	if err := install.inject(nil); err == nil {
		t.Fatal("expected an error when no app.asar is present")
	}
	if utils.Exists(filepath.Join(resources, "app")) {
		t.Error("no shadow app/ should be created when there's nothing to inject")
	}
}

func TestInject_EmptyResourcesPathErrors(t *testing.T) {
	install := &DiscordInstall{ResourcesPath: "", Channel: models.Stable}
	if err := install.inject(nil); err == nil {
		t.Fatal("expected an error for an empty resources path (must not touch the cwd)")
	}
}

func TestUninject_EmptyResourcesPathErrors(t *testing.T) {
	install := &DiscordInstall{ResourcesPath: "", Channel: models.Stable}
	if err := install.uninject(); err == nil {
		t.Fatal("expected an error for an empty resources path (must not RemoveAll the cwd)")
	}
}

// Rolling back a failed *re-injection* of an already-injected install must still
// leave Discord launchable: the preserve step is a no-op (no live app.asar), but
// rollback removes app/, so it must restore app.asar from the preserved copy.
func TestInject_RollbackRestoresLaunchableOnReinject(t *testing.T) {
	resources := t.TempDir()
	preserved := []byte("preserved discord app")
	if err := os.WriteFile(filepath.Join(resources, "betterdiscord.app.asar"), preserved, 0o644); err != nil {
		t.Fatalf("seed preserved: %v", err)
	}
	// Already-injected: app/ exists. Make index.js a directory so the index.js
	// write fails *after* the (no-op) preserve step, forcing rollback.
	if err := os.MkdirAll(filepath.Join(resources, "app", "index.js"), 0o755); err != nil {
		t.Fatalf("seed app/index.js dir: %v", err)
	}

	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	if err := install.inject(nil); err == nil {
		t.Fatal("expected inject to fail when app/index.js can't be written")
	}

	// Discord must remain launchable: app.asar restored from the preserved copy.
	restored, err := os.ReadFile(filepath.Join(resources, "app.asar"))
	if err != nil {
		t.Fatalf("app.asar not restored after rollback: %v", err)
	}
	if string(restored) != string(preserved) {
		t.Errorf("restored app.asar = %q, expected %q", restored, preserved)
	}
	if utils.Exists(filepath.Join(resources, "app")) {
		t.Error("shadow app/ should be removed by rollback")
	}
}

func TestInject_ProbeFailAbortsBeforeRename(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("chmod-based write denial is unreliable on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root bypasses directory write permissions")
	}
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}

	if err := os.Chmod(resources, 0o555); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(resources, 0o755) })

	if err := install.inject(nil); err == nil {
		t.Fatal("expected inject to fail the writability probe")
	}

	// The bundle must be untouched: app.asar still present, nothing renamed.
	_ = os.Chmod(resources, 0o755)
	got, err := os.ReadFile(filepath.Join(resources, "app.asar"))
	if err != nil {
		t.Fatalf("app.asar was disturbed by a probe-failed inject: %v", err)
	}
	if string(got) != string(original) {
		t.Errorf("app.asar content changed: got %q", got)
	}
	if utils.Exists(filepath.Join(resources, "betterdiscord.app.asar")) {
		t.Error("no rename should have happened after a probe failure")
	}
}

// Regression for the "invisible after injection" bug: injecting renames app.asar
// to betterdiscord.app.asar, so a resolver anchored only on app.asar would fail
// to find the install afterward — breaking repair and, critically, uninstall.
func TestInjectThenResolve_RemainsDiscoverable(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	resources := filepath.Join(root, "app-1.0.9002", "resources")
	writeAppAsar(t, resources) // pristine install

	if validateWindowsStyleInstall(root) == nil {
		t.Fatal("precondition: pristine install should resolve")
	}

	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	if err := install.inject(nil); err != nil {
		t.Fatalf("inject: %v", err)
	}

	// The fix: it must still resolve from the top-level Discord root after injection.
	resolved := validateWindowsStyleInstall(root)
	if resolved == nil {
		t.Fatal("injected install no longer resolves — uninstall would be impossible")
	}
	if resolved.ResourcesPath != resources {
		t.Errorf("ResourcesPath = %s, expected %s", resolved.ResourcesPath, resources)
	}
	if !resolved.IsInjected() {
		t.Error("expected the resolved install to report IsInjected")
	}

	// And uninstall works from the resolved install.
	if err := resolved.uninject(); err != nil {
		t.Fatalf("uninject: %v", err)
	}
	if !utils.Exists(filepath.Join(resources, "app.asar")) {
		t.Error("app.asar not restored after uninject")
	}
}

func TestInject_RollbackOnMidOpFailure(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}

	// Block mkdir(resources/app) by pre-creating a regular file at that path.
	// This forces a failure *after* the app.asar rename, exercising rollback.
	if err := os.WriteFile(filepath.Join(resources, "app"), []byte("blocker"), 0o644); err != nil {
		t.Fatalf("seed blocker file: %v", err)
	}

	if err := install.inject(nil); err == nil {
		t.Fatal("expected inject to fail when app/ can't be created")
	}

	// Rollback must restore the original app.asar and drop the preserved copy.
	restored, err := os.ReadFile(filepath.Join(resources, "app.asar"))
	if err != nil {
		t.Fatalf("app.asar not restored after rollback: %v", err)
	}
	if string(restored) != string(original) {
		t.Errorf("restored app.asar = %q, expected %q", restored, original)
	}
	if utils.Exists(filepath.Join(resources, "betterdiscord.app.asar")) {
		t.Error("preserved asar should be gone after rollback")
	}
}