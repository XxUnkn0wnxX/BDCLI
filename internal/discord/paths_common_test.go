package discord

import (
	"os"
	"path/filepath"
	"testing"
)

func writeCoreAsar(t *testing.T, resourcesPath string) {
	t.Helper()
	if err := os.MkdirAll(resourcesPath, 0755); err != nil {
		t.Fatalf("Failed to create core path: %v", err)
	}
	if err := os.WriteFile(filepath.Join(resourcesPath, "core.asar"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write core.asar: %v", err)
	}
}

func TestValidateWindowsStyleInstall_FromDiscordRoot(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	versionDir := filepath.Join(root, "app-1.0.9002")
	coreWrap := filepath.Join(versionDir, "modules", "discord_desktop_core-1", "discord_desktop_core")

	writeCoreAsar(t, coreWrap)

	result := validateWindowsStyleInstall(root)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if result.ResourcesPath != coreWrap {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, coreWrap)
	}
}

func TestValidateWindowsStyleInstall_FromAppFolder(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	versionDir := filepath.Join(root, "app-1.0.9002")
	coreWrap := filepath.Join(versionDir, "modules", "discord_desktop_core-1", "discord_desktop_core")

	writeCoreAsar(t, coreWrap)

	result := validateWindowsStyleInstall(versionDir)
	if result == nil {
		t.Fatalf("Expected install for %s", versionDir)
	}
	if result.ResourcesPath != coreWrap {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, coreWrap)
	}
}

func TestValidateWindowsStyleInstall_FromCoreFolder(t *testing.T) {
	tmpDir := t.TempDir()
	resourcesPath := filepath.Join(tmpDir, "discord_desktop_core")
	writeCoreAsar(t, resourcesPath)

	result := validateWindowsStyleInstall(resourcesPath)
	if result == nil {
		t.Fatalf("Expected install for %s", resourcesPath)
	}
	if result.ResourcesPath != resourcesPath {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resourcesPath)
	}
}

func TestValidateWindowsStyleInstall_MissingAsar(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "Discord")
	versionDir := filepath.Join(root, "app-1.0.9002")
	coreWrap := filepath.Join(versionDir, "modules", "discord_desktop_core-1", "discord_desktop_core")

	if err := os.MkdirAll(coreWrap, 0755); err != nil {
		t.Fatalf("Failed to create core path: %v", err)
	}

	result := validateWindowsStyleInstall(root)
	if result != nil {
		t.Fatalf("Expected no install when core.asar is missing")
	}
}

func TestValidateUnixStyleInstall_FromDiscordRoot(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "discord")
	resourcesPath := filepath.Join(root, "0.0.35", "modules", "discord_desktop_core")

	writeCoreAsar(t, resourcesPath)

	result := validateUnixStyleInstall(root, true, true)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if result.ResourcesPath != resourcesPath {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resourcesPath)
	}
}

func TestValidateUnixStyleInstall_FromVersionFolder(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "discord")
	versionDir := filepath.Join(root, "0.0.35")
	resourcesPath := filepath.Join(versionDir, "modules", "discord_desktop_core")

	writeCoreAsar(t, resourcesPath)

	result := validateUnixStyleInstall(versionDir, true, true)
	if result == nil {
		t.Fatalf("Expected install for %s", versionDir)
	}
	if result.ResourcesPath != resourcesPath {
		t.Errorf("ResourcesPath = %s, expected %s", result.ResourcesPath, resourcesPath)
	}
}

func TestValidateUnixStyleInstall_FlatpakDetection(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "com.discordapp.Discord", "config", "discord")
	resourcesPath := filepath.Join(root, "0.0.35", "modules", "discord_desktop_core")

	writeCoreAsar(t, resourcesPath)

	result := validateUnixStyleInstall(root, true, false)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if !result.IsFlatpak {
		t.Fatalf("Expected flatpak detection")
	}
	if result.IsSnap {
		t.Fatalf("Did not expect snap detection")
	}
}

func TestValidateUnixStyleInstall_SnapDetection(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "snap", "discord", "current", ".config", "discord")
	resourcesPath := filepath.Join(root, "0.0.35", "modules", "discord_desktop_core")

	writeCoreAsar(t, resourcesPath)

	result := validateUnixStyleInstall(root, false, true)
	if result == nil {
		t.Fatalf("Expected install for %s", root)
	}
	if !result.IsSnap {
		t.Fatalf("Expected snap detection")
	}
	if result.IsFlatpak {
		t.Fatalf("Did not expect flatpak detection")
	}
}