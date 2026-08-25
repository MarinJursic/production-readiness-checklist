package provider

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestStoredAuthenticationUsesPrivateProviderDirectory(t *testing.T) {
	original := locateUserConfigDirectory
	root := t.TempDir()
	locateUserConfigDirectory = func() (string, error) { return root, nil }
	t.Cleanup(func() { locateUserConfigDirectory = original })

	for _, name := range []string{"codex", "claude"} {
		if _, err := storedAuthenticationDirectory(name); err == nil ||
			!strings.Contains(err.Error(), "prc login "+name) {
			t.Fatalf("missing %s login error = %v", name, err)
		}
		if err := MarkStoredAuthentication(name); err != nil {
			t.Fatal(err)
		}
		path, err := storedAuthenticationDirectory(name)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(path, filepath.Join(root, "prc", "provider-auth")) {
			t.Fatalf("authentication path %q escaped root %q", path, root)
		}
		if runtime.GOOS != "windows" {
			information, err := os.Stat(path)
			if err != nil || information.Mode().Perm() != 0o700 {
				t.Fatalf("authentication directory mode=%v err=%v", information, err)
			}
			marker, err := os.Stat(filepath.Join(path, authenticationMarker))
			if err != nil || marker.Mode().Perm() != 0o600 {
				t.Fatalf("authentication marker mode=%v err=%v", marker, err)
			}
		}
		if err := ClearStoredAuthentication(name); err != nil {
			t.Fatal(err)
		}
		if _, err := storedAuthenticationDirectory(name); err == nil {
			t.Fatalf("cleared %s login remained usable", name)
		}
	}
}

func TestStoredAuthenticationRejectsSymlinkRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink setup requires platform-specific privileges on Windows")
	}
	original := locateUserConfigDirectory
	root := t.TempDir()
	locateUserConfigDirectory = func() (string, error) { return root, nil }
	t.Cleanup(func() { locateUserConfigDirectory = original })

	path, err := AuthenticationDirectory("codex")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), path); err != nil {
		t.Fatal(err)
	}
	if _, err := PrepareAuthenticationDirectory("codex"); err == nil ||
		!strings.Contains(err.Error(), "not a regular directory") {
		t.Fatalf("symlink authentication root was accepted: %v", err)
	}
}
