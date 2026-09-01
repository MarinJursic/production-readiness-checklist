package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSetupUsesCatalogDiscoveryAndDoesNotRequireAI(t *testing.T) {
	t.Chdir(filepath.Join("..", ".."))
	var stdout, stderr bytes.Buffer
	code := run([]string{"setup", t.TempDir()}, &stdout, &stderr)
	if code != exitSuccess || stderr.Len() != 0 ||
		!strings.Contains(stdout.String(), "Bundled rules:") ||
		!strings.Contains(stdout.String(), "Optional AI:") ||
		!strings.Contains(stdout.String(), "Run `prc` for the normal local scan") {
		t.Fatalf("setup exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestExplicitUpdateCheckNeverInstalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet || request.Header.Get("User-Agent") == "" {
			t.Fatalf("unexpected update request: %+v", request)
		}
		fmt.Fprint(response, `{"dist-tags":{"latest":"9.8.7"}}`)
	}))
	defer server.Close()
	originalURL, originalClient, originalVersion := npmPackageRegistryURL, updateHTTPClient, version
	npmPackageRegistryURL, updateHTTPClient, version = server.URL, server.Client(), "1.2.3"
	t.Cleanup(func() {
		npmPackageRegistryURL, updateHTTPClient, version = originalURL, originalClient, originalVersion
	})
	var stdout, stderr bytes.Buffer
	if code := run([]string{"update"}, &stdout, &stderr); code != exitSuccess || stderr.Len() != 0 ||
		!strings.Contains(stdout.String(), "Installed: 1.2.3") || !strings.Contains(stdout.String(), "Latest:    9.8.7") ||
		!strings.Contains(stdout.String(), "npm install -g @marinjursic/prc@latest") {
		t.Fatalf("update exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestReleaseVersionComparisonUsesNumbers(t *testing.T) {
	for _, test := range []struct {
		left, right string
		want        int
	}{
		{"1.9.9", "1.10.0", -1},
		{"2.0.0", "1.99.99", 1},
		{"3.4.5", "3.4.5", 0},
	} {
		if got := compareReleaseVersions(test.left, test.right); got != test.want {
			t.Fatalf("compareReleaseVersions(%q, %q) = %d, want %d", test.left, test.right, got, test.want)
		}
	}
	for _, invalid := range []string{"1.2", "v1.2.3", "1.02.3", "1234567890.2.3"} {
		if validReleaseVersion(invalid) {
			t.Fatalf("invalid release version %q was accepted", invalid)
		}
	}
}

func TestReportAndCacheCommandsStayInsideScannerCache(t *testing.T) {
	cacheRoot := t.TempDir()
	originalCache := userCacheDirectory
	userCacheDirectory = func() (string, error) { return cacheRoot, nil }
	t.Cleanup(func() { userCacheDirectory = originalCache })
	reports := filepath.Join(cacheRoot, "prc", "reports")
	ai := filepath.Join(cacheRoot, "prc", "control-reviews", "inventory", "catalog", "codex")
	if err := os.MkdirAll(reports, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(ai, 0o700); err != nil {
		t.Fatal(err)
	}
	older := filepath.Join(reports, "project-0000000000000001.html")
	newer := filepath.Join(reports, "project-0000000000000002.html")
	for path, data := range map[string]string{older: "old", newer: "new", filepath.Join(ai, "batch.json"): "resume"} {
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	oldTime := time.Now().Add(-time.Hour)
	if err := os.Chtimes(older, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if code := run([]string{"report", "path"}, &stdout, &stderr); code != exitSuccess ||
		strings.TrimSpace(stdout.String()) != newer {
		t.Fatalf("report path exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"report", "--limit", "1"}, &stdout, &stderr); code != exitConfiguration ||
		!strings.Contains(stderr.String(), "applies only") {
		t.Fatalf("report limit exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"cache", "status"}, &stdout, &stderr); code != exitSuccess ||
		!strings.Contains(stdout.String(), "Reports:") || !strings.Contains(stdout.String(), "AI resume data:") ||
		!strings.Contains(stdout.String(), "Nothing was deleted") {
		t.Fatalf("cache status exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"cache", "clean", "--reports"}, &stdout, &stderr); code != exitSuccess ||
		!strings.Contains(stdout.String(), "Removed 2 scanner cache files") {
		t.Fatalf("cache clean exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	if _, err := os.Stat(filepath.Join(ai, "batch.json")); err != nil {
		t.Fatalf("report cleanup touched AI resume data: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	if code := run([]string{"cache", "--help"}, &stdout, &stderr); code != exitSuccess ||
		!strings.Contains(stdout.String(), "prc cache clean") {
		t.Fatalf("cache help exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
}

func TestCompletionAndAIPlanAreProviderFree(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"completion", "zsh"}, &stdout, &stderr); code != exitSuccess ||
		!strings.Contains(stdout.String(), "#compdef prc") || !strings.Contains(stdout.String(), "report") ||
		!strings.Contains(stdout.String(), "verify") {
		t.Fatalf("completion exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	stdout.Reset()
	code := run([]string{
		"scan", t.TempDir(), "--catalog-root", filepath.Join("..", ".."), "--ai", "codex",
		"--review-control", "PRC-02-001", "--review-plan", "--format", "json", "--no-report",
	}, &stdout, &stderr)
	if code != exitSuccess {
		t.Fatalf("AI plan exit=%d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for _, expected := range []string{`"schema_version": "prc.ai-review-preview/v0.2"`, `"provider": "codex"`, `"controls": 1`, `"batches": 1`, `"maximum_context_bytes_per_batch": 393216`} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("AI plan omitted %q: %s", expected, stdout.String())
		}
	}
}
