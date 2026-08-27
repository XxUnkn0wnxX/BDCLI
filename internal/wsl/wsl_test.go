package wsl

import (
	"strings"
	"sync"
	"testing"
)

func resetWSLInfo() {
	once = sync.Once{}
	info = nil
}

func TestInfo_NotWSLWhenNoSignals(t *testing.T) {
	t.Setenv("WSL_DISTRO_NAME", "")
	t.Setenv("WSL_INTEROP", "")
	resetWSLInfo()

	info := Info()
	if strings.Contains(info.KernelVersion, "microsoft") {
		if !info.IsWSL {
			t.Fatalf("Expected IsWSL true when kernel indicates WSL")
		}
	} else if info.IsWSL {
		t.Fatalf("Expected IsWSL false when no WSL signals are set")
	}
	if info.DistroName != "" {
		t.Fatalf("Expected empty DistroName, got %q", info.DistroName)
	}
	if info.InteropPath != "" {
		t.Fatalf("Expected empty InteropPath, got %q", info.InteropPath)
	}
}

func TestInfo_WSLDistroName(t *testing.T) {
	t.Setenv("WSL_DISTRO_NAME", "Ubuntu")
	t.Setenv("WSL_INTEROP", "interop-path")
	resetWSLInfo()

	info := Info()
	if !info.IsWSL {
		t.Fatalf("Expected IsWSL true when WSL_DISTRO_NAME is set")
	}
	if info.DistroName != "Ubuntu" {
		t.Fatalf("Expected DistroName Ubuntu, got %q", info.DistroName)
	}
	if info.InteropPath != "interop-path" {
		t.Fatalf("Expected InteropPath interop-path, got %q", info.InteropPath)
	}
}
