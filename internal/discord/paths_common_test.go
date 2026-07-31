package discord

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/betterdiscord/cli/internal/models"
)

// writeAppAsar creates a resources dir seeded with an app.asar.
func writeAppAsar(t *testing.T, resourcesDir string) {
	t.Helper()
	if err := os.MkdirAll(resourcesDir, 0755); err != nil {
		t.Fatalf("Failed to create resources dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resourcesDir, "app.asar"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write app.asar: %v", err)
	}
}

// writeInjectedResources creates a resources dir in the *injected* state:
// app.asar has been renamed to betterdiscord.app.asar and a shadow app/ exists.
func writeInjectedResources(t *testing.T, resourcesDir string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(resourcesDir, "app"), 0755); err != nil {
		t.Fatalf("Failed to create app dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resourcesDir, "betterdiscord.app.asar"), []byte("preserved"), 0644); err != nil {
		t.Fatalf("Failed to write preserved asar: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resourcesDir, "app", "index.js"), []byte("// bd"), 0644); err != nil {
		t.Fatalf("Failed to write index.js: %v", err)
	}
}

// Regression: an install stays resolvable after injection (app.asar renamed to
// betterdiscord.app.asar). If it didn't, users couldn't repair or uninstall it.
func TestValidateWindowsStyleInstall_ResolvesInjected(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	resources := filepath.Join(root, "app-1.0.9002", "resources")
	writeInjectedResources(t, resources) // no app.asar, only betterdiscord.app.asar

	result := validateWindowsStyleInstall(root)
	if result == nil {
		t.Fatalf("injected install must still resolve for %s", root)
	}
	if result.ResourcesPath != resources {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resources)
	}
}

func TestResolveResources_InjectedResourcesDir(t *testing.T) {
	// Browsing/uninstalling straight to an injected resources dir must resolve.
	resources := filepath.Join(t.TempDir(), "resources")
	writeInjectedResources(t, resources)

	if got := resolveResources(resources); got != resources {
		t.Errorf("resolveResources(injected) = %q, expected %q", got, resources)
	}
}

func TestValidateWindowsStyleInstall_FromDiscordRoot(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	resources := filepath.Join(root, "app-1.0.9002", "resources")
	writeAppAsar(t, resources)

	result := validateWindowsStyleInstall(root)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if result.ResourcesPath != resources {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resources)
	}
}

func TestValidateWindowsStyleInstall_PicksLatestVersion(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	// An older leftover version dir plus the current one.
	writeAppAsar(t, filepath.Join(root, "app-1.0.9002", "resources"))
	latest := filepath.Join(root, "app-1.0.10000", "resources")
	writeAppAsar(t, latest)

	result := validateWindowsStyleInstall(root)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if result.ResourcesPath != latest {
		t.Errorf("ResourcesPath = %s, expected latest %s", result.ResourcesPath, latest)
	}
}

func TestValidateWindowsStyleInstall_FromAppFolder(t *testing.T) {
	tmpDir := t.TempDir()
	versionDir := filepath.Join(tmpDir, "Discord", "app-1.0.9002")
	resources := filepath.Join(versionDir, "resources")
	writeAppAsar(t, resources)

	result := validateWindowsStyleInstall(versionDir)
	if result == nil {
		t.Fatalf("Expected install for %s", versionDir)
	}
	if result.ResourcesPath != resources {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resources)
	}
}

func TestValidateWindowsStyleInstall_FromResourcesFolder(t *testing.T) {
	resources := filepath.Join(t.TempDir(), "resources")
	writeAppAsar(t, resources)

	result := validateWindowsStyleInstall(resources)
	if result == nil {
		t.Fatalf("Expected install for %s", resources)
	}
	if result.ResourcesPath != resources {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resources)
	}
}

func TestValidateWindowsStyleInstall_MissingAsar(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	// resources dir exists but has no app.asar.
	if err := os.MkdirAll(filepath.Join(root, "app-1.0.9002", "resources"), 0755); err != nil {
		t.Fatalf("Failed to create resources dir: %v", err)
	}

	if result := validateWindowsStyleInstall(root); result != nil {
		t.Fatalf("Expected no install when app.asar is missing")
	}
}

func TestValidateUnixStyleInstall_FromChannelRoot(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "discord")
	resources := filepath.Join(root, "app-0.0.90", "resources")
	writeAppAsar(t, resources)

	result := validateUnixStyleInstall(root, true, true)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if result.ResourcesPath != resources {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resources)
	}
	if result.IsFlatpak || result.IsSnap {
		t.Errorf("plain path should not flag flatpak/snap: %+v", result)
	}
}

func TestValidateUnixStyleInstall_FromVersionFolder(t *testing.T) {
	tmpDir := t.TempDir()
	versionDir := filepath.Join(tmpDir, "discord", "app-0.0.90")
	resources := filepath.Join(versionDir, "resources")
	writeAppAsar(t, resources)

	result := validateUnixStyleInstall(versionDir, true, true)
	if result == nil {
		t.Fatalf("Expected install for %s", versionDir)
	}
	if result.ResourcesPath != resources {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resources)
	}
}

func TestValidateUnixStyleInstall_FlatpakDetection(t *testing.T) {
	tmpDir := t.TempDir()
	// Flatpak deployment layout: files/{channel-}/resources (no app-* segment).
	resources := filepath.Join(tmpDir, "com.discordapp.Discord", "files", "discord", "resources")
	writeAppAsar(t, resources)

	result := validateUnixStyleInstall(resources, true, false)
	if result == nil {
		t.Fatalf("Expected install for %s", resources)
	}
	if !result.IsFlatpak {
		t.Fatalf("Expected flatpak detection for %s", resources)
	}
	if result.IsSnap {
		t.Fatalf("Did not expect snap detection")
	}
}

func TestValidateUnixStyleInstall_MacOSBundle(t *testing.T) {
	// macOS: app.asar lives in {Bundle}.app/Contents/Resources; channel/version
	// come from build_info.json (the bundle name has a space and no version).
	tmpDir := t.TempDir()
	bundle := filepath.Join(tmpDir, "Discord Canary.app")
	resources := filepath.Join(bundle, "Contents", "Resources")
	writeAppAsar(t, resources)
	if err := os.WriteFile(filepath.Join(resources, "build_info.json"),
		[]byte(`{"releaseChannel":"canary","version":"1.0.1"}`), 0644); err != nil {
		t.Fatalf("write build_info: %v", err)
	}

	// Resolving from the bundle path (as a user browsing to Discord.app would).
	result := validateUnixStyleInstall(bundle, false, false)
	if result == nil {
		t.Fatalf("Expected install for bundle %s", bundle)
	}
	if result.ResourcesPath != resources {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resources)
	}
	if result.Channel != models.Canary {
		t.Errorf("Channel = %v, expected Canary", result.Channel)
	}
	if result.Version != "1.0.1" {
		t.Errorf("Version = %q, expected 1.0.1", result.Version)
	}
}

func TestValidateUnixStyleInstall_SnapSegmentDetection(t *testing.T) {
	tmpDir := t.TempDir()

	// A real "/snap/" path segment is detected.
	snapRes := filepath.Join(tmpDir, "snap", "discord", "current", "resources")
	writeAppAsar(t, snapRes)
	if got := validateUnixStyleInstall(snapRes, false, true); got == nil || !got.IsSnap {
		t.Errorf("expected IsSnap true for a /snap/ path, got %+v", got)
	}

	// "snap" as a suffix of another segment must not false-positive.
	fakeRes := filepath.Join(tmpDir, "mysnap", "discord", "resources")
	writeAppAsar(t, fakeRes)
	if got := validateUnixStyleInstall(fakeRes, false, true); got == nil || got.IsSnap {
		t.Errorf("expected IsSnap false for a .../mysnap/... path, got %+v", got)
	}
}

func TestReadBuildInfo(t *testing.T) {
	t.Run("present", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "build_info.json"),
			[]byte(`{"releaseChannel":"canary","version":"1.0.1234"}`), 0644); err != nil {
			t.Fatalf("write build_info: %v", err)
		}
		info, ok := readBuildInfo(dir)
		if !ok {
			t.Fatal("expected ok=true for a present build_info.json")
		}
		if info.ReleaseChannel != "canary" || info.Version != "1.0.1234" {
			t.Errorf("parsed = %+v", info)
		}
	})

	t.Run("absent", func(t *testing.T) {
		if _, ok := readBuildInfo(t.TempDir()); ok {
			t.Error("expected ok=false when build_info.json is absent")
		}
	})

	t.Run("malformed", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "build_info.json"), []byte("{not json"), 0644); err != nil {
			t.Fatalf("write: %v", err)
		}
		if _, ok := readBuildInfo(dir); ok {
			t.Error("expected ok=false for malformed build_info.json")
		}
	})
}

func TestNewResourcesInstall_PrefersBuildInfo(t *testing.T) {
	// Path segments say stable/no-version, but build_info.json says canary/1.0.5.
	resources := filepath.Join(t.TempDir(), "discord", "app-0.0.1", "resources")
	writeAppAsar(t, resources)
	if err := os.WriteFile(filepath.Join(resources, "build_info.json"),
		[]byte(`{"releaseChannel":"canary","version":"1.0.5"}`), 0644); err != nil {
		t.Fatalf("write build_info: %v", err)
	}

	install := newResourcesInstall(resources)
	if install.Channel != models.Canary {
		t.Errorf("Channel = %v, expected Canary from build_info", install.Channel)
	}
	if install.Version != "1.0.5" {
		t.Errorf("Version = %q, expected 1.0.5 from build_info", install.Version)
	}
}

func TestNewResourcesInstall_FallsBackToPath(t *testing.T) {
	// No build_info.json → channel/version come from the path.
	resources := filepath.Join(t.TempDir(), "discordcanary", "app-0.0.90", "resources")
	writeAppAsar(t, resources)

	install := newResourcesInstall(resources)
	if install.Channel != models.Canary {
		t.Errorf("Channel = %v, expected Canary from path", install.Channel)
	}
	if install.Version != "0.0.90" {
		t.Errorf("Version = %q, expected 0.0.90 from path", install.Version)
	}
}