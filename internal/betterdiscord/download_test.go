package betterdiscord

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// withURLs temporarily overrides the download endpoints for a test and restores
// them on cleanup.
func withURLs(t *testing.T, website, github string) {
	t.Helper()
	origWebsite, origGithub := websiteAsarURL, githubLatestReleaseURL
	websiteAsarURL = website
	githubLatestReleaseURL = github
	t.Cleanup(func() {
		websiteAsarURL = origWebsite
		githubLatestReleaseURL = origGithub
	})
}

func newBDInstallWithDataDir(t *testing.T) *BDInstall {
	t.Helper()
	install := New(filepath.Join(t.TempDir(), "BetterDiscord"))
	if err := os.MkdirAll(install.Data(), 0o755); err != nil {
		t.Fatalf("failed to create data dir: %v", err)
	}
	return install
}

func assertFileContents(t *testing.T, path, want string) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read %s: %v", path, err)
	}
	if string(got) != want {
		t.Errorf("contents of %s = %q, expected %q", path, string(got), want)
	}
}

func TestDownload_FromWebsite(t *testing.T) {
	const body = "asar-from-website"
	website := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("x-bd-version", "1.2.3")
		fmt.Fprint(w, body)
	}))
	defer website.Close()

	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("GitHub fallback should not be called when the website succeeds")
		http.Error(w, "unexpected", http.StatusInternalServerError)
	}))
	defer github.Close()

	withURLs(t, website.URL, github.URL)

	install := newBDInstallWithDataDir(t)
	if err := install.download(); err != nil {
		t.Fatalf("download() failed: %v", err)
	}
	if !install.HasDownloaded() {
		t.Error("expected HasDownloaded() to be true")
	}
	assertFileContents(t, install.Asar(), body)
}

func TestDownload_FallsBackToGitHub(t *testing.T) {
	const body = "asar-from-github"
	asset := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, body)
	}))
	defer asset.Close()

	// Website fails, so download() should fall back to the GitHub release.
	website := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusInternalServerError)
	}))
	defer website.Close()

	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `{"tag_name":"v9.9.9","assets":[{"name":"betterdiscord.asar","url":%q}]}`, asset.URL)
	}))
	defer github.Close()

	withURLs(t, website.URL, github.URL)

	install := newBDInstallWithDataDir(t)
	if err := install.download(); err != nil {
		t.Fatalf("download() failed: %v", err)
	}
	if !install.HasDownloaded() {
		t.Error("expected HasDownloaded() to be true after GitHub fallback")
	}
	assertFileContents(t, install.Asar(), body)
}

func TestDownload_GitHubMissingAsset(t *testing.T) {
	website := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "down", http.StatusInternalServerError)
	}))
	defer website.Close()

	github := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"tag_name":"v9.9.9","assets":[{"name":"something-else.zip","url":"http://example.invalid"}]}`)
	}))
	defer github.Close()

	withURLs(t, website.URL, github.URL)

	install := newBDInstallWithDataDir(t)
	if err := install.download(); err == nil {
		t.Fatal("expected an error when the betterdiscord.asar asset is missing")
	}
}

func TestDownload_SkipsWhenAlreadyDownloaded(t *testing.T) {
	install := newBDInstallWithDataDir(t)
	install.hasDownloaded = true

	// Point the endpoints at a non-routable address to prove the network is
	// never touched when the asar is already downloaded.
	withURLs(t, "http://127.0.0.1:0", "http://127.0.0.1:0")

	if err := install.download(); err != nil {
		t.Fatalf("download() should be a no-op when already downloaded: %v", err)
	}
}