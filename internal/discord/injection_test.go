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
	// without clobbering the real, already-preserved asar.
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

func TestInject_RefusesAlteredShadowWithoutPreservedArchive(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	writeManagedApp(t, resources)
	if err := os.WriteFile(filepath.Join(resources, "app", "index.js"), []byte("altered"), 0o644); err != nil {
		t.Fatalf("alter index.js: %v", err)
	}

	if err := install.inject(nil); err == nil {
		t.Fatal("expected inject to refuse an altered app shadow without preserved asar")
	}
	assertAppAsarUnchanged(t, resources, original)
	altered, err := os.ReadFile(filepath.Join(resources, "app", "index.js"))
	if err != nil {
		t.Fatalf("altered index.js was removed: %v", err)
	}
	if string(altered) != "altered" {
		t.Errorf("altered index.js = %q, expected %q", altered, "altered")
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
	writeManagedApp(t, resources)

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
	writeManagedApp(t, resources)

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

func TestUninject_RefusesForeignAppFileBeforeMutation(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	appPath := filepath.Join(resources, "app")
	if err := os.WriteFile(appPath, []byte("foreign app state"), 0o644); err != nil {
		t.Fatalf("seed foreign app file: %v", err)
	}

	if err := install.uninject(); err == nil {
		t.Fatal("expected uninject to refuse a foreign app file")
	}
	assertAppAsarUnchanged(t, resources, original)
	got, err := os.ReadFile(appPath)
	if err != nil {
		t.Fatalf("foreign app file was removed: %v", err)
	}
	if string(got) != "foreign app state" {
		t.Errorf("foreign app file = %q, expected %q", got, "foreign app state")
	}
}

func TestUninject_RefusesAppSymlinkBeforeMutation(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "marker"), []byte("keep me"), 0o644); err != nil {
		t.Fatalf("seed symlink target: %v", err)
	}
	appPath := filepath.Join(resources, "app")
	if err := os.Symlink(target, appPath); err != nil {
		t.Fatalf("seed app symlink: %v", err)
	}

	if err := install.uninject(); err == nil {
		t.Fatal("expected uninject to refuse an app symlink")
	}
	assertAppAsarUnchanged(t, resources, original)
	info, err := os.Lstat(appPath)
	if err != nil {
		t.Fatalf("app symlink was removed: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("app path is no longer a symlink")
	}
	link, err := os.Readlink(appPath)
	if err != nil {
		t.Fatalf("read app symlink: %v", err)
	}
	if link != target {
		t.Errorf("app symlink target = %q, expected %q", link, target)
	}
}

func TestUninject_RefusesNestedAppPayloadBeforeMutation(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	preserved := []byte("preserved discord app")
	if err := os.WriteFile(filepath.Join(resources, "betterdiscord.app.asar"), preserved, 0o644); err != nil {
		t.Fatalf("seed preserved app.asar: %v", err)
	}
	nested := filepath.Join(resources, "app", "nested", "payload.bin")
	if err := os.MkdirAll(filepath.Dir(nested), 0o755); err != nil {
		t.Fatalf("seed nested app/: %v", err)
	}
	if err := os.WriteFile(nested, []byte("keep me"), 0o644); err != nil {
		t.Fatalf("seed nested payload: %v", err)
	}

	if err := install.uninject(); err == nil {
		t.Fatal("expected uninject to refuse an app directory with nested payload")
	}
	live, err := os.ReadFile(filepath.Join(resources, "app.asar"))
	if err != nil {
		t.Fatalf("app.asar was disturbed: %v", err)
	}
	if string(live) != string(original) {
		t.Errorf("app.asar content changed: got %q, expected %q", live, original)
	}
	gotPreserved, err := os.ReadFile(filepath.Join(resources, "betterdiscord.app.asar"))
	if err != nil {
		t.Fatalf("preserved app.asar was removed: %v", err)
	}
	if string(gotPreserved) != string(preserved) {
		t.Errorf("preserved app.asar changed: got %q, expected %q", gotPreserved, preserved)
	}
	got, err := os.ReadFile(nested)
	if err != nil {
		t.Fatalf("nested app payload was removed: %v", err)
	}
	if string(got) != "keep me" {
		t.Errorf("nested app payload = %q, expected %q", got, "keep me")
	}
}

// Corrupted pre-state: only the shadow app/ remains — both app.asar and the
// preserved copy are gone (external interference; discovery normally filters
// this state out). uninject must refuse rather than delete app/ and report
// success on an install with nothing left for Electron to load.
func TestUninject_RefusesWhenNoAsarPresent(t *testing.T) {
	resources := t.TempDir()
	if err := os.MkdirAll(filepath.Join(resources, "app"), 0o755); err != nil {
		t.Fatalf("seed app/: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resources, "app", "index.js"), []byte("x"), 0o644); err != nil {
		t.Fatalf("seed index.js: %v", err)
	}

	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	if err := install.uninject(); err == nil {
		t.Fatal("expected an error when both app.asar and betterdiscord.app.asar are missing")
	}
	if !utils.Exists(filepath.Join(resources, "app", "index.js")) {
		t.Error("app/ must be left untouched when refusing to uninject")
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

// An unmanaged app/ entry must be rejected before a re-injection can remove it
// during rollback.
func TestInject_RefusesUnmanagedShadowBeforeMutation(t *testing.T) {
	resources := t.TempDir()
	preserved := []byte("preserved discord app")
	if err := os.WriteFile(filepath.Join(resources, "betterdiscord.app.asar"), preserved, 0o644); err != nil {
		t.Fatalf("seed preserved: %v", err)
	}
	// A directory at app/index.js is not a managed shadow asset.
	if err := os.MkdirAll(filepath.Join(resources, "app", "index.js"), 0o755); err != nil {
		t.Fatalf("seed app/index.js dir: %v", err)
	}

	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	if err := install.inject(nil); err == nil {
		t.Fatal("expected inject to refuse an unmanaged app shadow")
	}

	if utils.Exists(filepath.Join(resources, "app.asar")) {
		t.Fatal("app.asar should not appear after preflight refusal")
	}
	gotPreserved, err := os.ReadFile(filepath.Join(resources, "betterdiscord.app.asar"))
	if err != nil {
		t.Fatalf("preserved app.asar changed after preflight refusal: %v", err)
	}
	if string(gotPreserved) != string(preserved) {
		t.Errorf("preserved app.asar = %q, expected %q", gotPreserved, preserved)
	}
	if !utils.Exists(filepath.Join(resources, "app", "index.js")) {
		t.Error("unmanaged app/ should remain after preflight refusal")
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

	// Pre-create a regular file at resources/app. This must be rejected before
	// any app.asar rename, exercising preservation of an unmanaged blocker.
	if err := os.WriteFile(filepath.Join(resources, "app"), []byte("blocker"), 0o644); err != nil {
		t.Fatalf("seed blocker file: %v", err)
	}

	if err := install.inject(nil); err == nil {
		t.Fatal("expected inject to refuse an unmanaged app path")
	}

	// Preflight must preserve the original app.asar, blocker bytes, and absence
	// of a preserved copy because no mutation is safe around unmanaged app state.
	restored, err := os.ReadFile(filepath.Join(resources, "app.asar"))
	if err != nil {
		t.Fatalf("app.asar not restored after rollback: %v", err)
	}
	if string(restored) != string(original) {
		t.Errorf("restored app.asar = %q, expected %q", restored, original)
	}
	blocker, err := os.ReadFile(filepath.Join(resources, "app"))
	if err != nil {
		t.Fatalf("blocker app file was removed: %v", err)
	}
	if string(blocker) != "blocker" {
		t.Errorf("blocker app file = %q, expected %q", blocker, "blocker")
	}
	if utils.Exists(filepath.Join(resources, "betterdiscord.app.asar")) {
		t.Error("no preserved asar should appear after preflight refusal")
	}
}

func TestInject_RefusesUnmanagedShadowDirectoryBeforeMutation(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	nested := filepath.Join(resources, "app", "nested", "payload.bin")
	if err := os.MkdirAll(filepath.Dir(nested), 0o755); err != nil {
		t.Fatalf("seed nested app/: %v", err)
	}
	if err := os.WriteFile(nested, []byte("keep me"), 0o644); err != nil {
		t.Fatalf("seed nested payload: %v", err)
	}

	if err := install.inject(nil); err == nil {
		t.Fatal("expected inject to refuse an unmanaged app directory")
	}
	assertAppAsarUnchanged(t, resources, original)
	got, err := os.ReadFile(nested)
	if err != nil {
		t.Fatalf("nested app payload was removed: %v", err)
	}
	if string(got) != "keep me" {
		t.Errorf("nested app payload = %q, expected %q", got, "keep me")
	}
}

func TestInject_RefusesAppSymlinkBeforeMutation(t *testing.T) {
	resources, original := newResourcesDir(t)
	install := &DiscordInstall{ResourcesPath: resources, Channel: models.Stable}
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "marker"), []byte("keep me"), 0o644); err != nil {
		t.Fatalf("seed symlink target: %v", err)
	}
	appDir := filepath.Join(resources, "app")
	if err := os.Symlink(target, appDir); err != nil {
		t.Fatalf("seed app symlink: %v", err)
	}

	if err := install.inject(nil); err == nil {
		t.Fatal("expected inject to refuse an app symlink")
	}
	assertAppAsarUnchanged(t, resources, original)
	info, err := os.Lstat(appDir)
	if err != nil {
		t.Fatalf("app symlink was removed: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Fatal("app path is no longer a symlink")
	}
}

func writeManagedApp(t *testing.T, resources string) {
	t.Helper()
	appDir := filepath.Join(resources, "app")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		t.Fatalf("seed managed app/: %v", err)
	}
	if err := os.WriteFile(filepath.Join(appDir, "index.js"), []byte(appIndexScript), 0o644); err != nil {
		t.Fatalf("seed managed app/index.js: %v", err)
	}
	if err := os.WriteFile(filepath.Join(appDir, "package.json"), []byte(appPackageJSON), 0o644); err != nil {
		t.Fatalf("seed managed app/package.json: %v", err)
	}
}

func assertAppAsarUnchanged(t *testing.T, resources string, original []byte) {
	t.Helper()
	got, err := os.ReadFile(filepath.Join(resources, "app.asar"))
	if err != nil {
		t.Fatalf("app.asar was disturbed: %v", err)
	}
	if string(got) != string(original) {
		t.Errorf("app.asar content changed: got %q, expected %q", got, original)
	}
	if utils.Exists(filepath.Join(resources, "betterdiscord.app.asar")) {
		t.Error("no preserved asar should appear after preflight refusal")
	}
}
