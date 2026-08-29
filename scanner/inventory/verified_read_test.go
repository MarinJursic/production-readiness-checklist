package inventory

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func TestReadVerifiedFileAcceptsOnlySealedBytes(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "README.md")
	data := []byte("# Ready\n")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(data)
	item := model.Inventory{Root: root, Digest: hex.EncodeToString(digest[:]), Files: []model.FileRecord{{
		Path: "README.md", Size: int64(len(data)), SHA256: hex.EncodeToString(digest[:]), Mode: 0o600,
	}}}
	got, err := ReadVerifiedFile(item, "README.md", 1024)
	if err != nil || string(got) != string(data) {
		t.Fatalf("verified read = %q, %v", got, err)
	}
	if _, err := ReadVerifiedFile(item, "../README.md", 1024); err == nil {
		t.Fatal("path escape was accepted")
	}
	if _, err := ReadVerifiedFile(item, "README.md", 1); err == nil {
		t.Fatal("oversized file was accepted")
	}
	if err := os.WriteFile(path, []byte("# Changed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadVerifiedFile(item, "README.md", 1024); err == nil {
		t.Fatal("changed file was accepted")
	}
}

func TestReadVerifiedFileRejectsSymlinkReplacement(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "target")
	if err := os.WriteFile(target, []byte("safe"), 0o600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256([]byte("safe"))
	item := model.Inventory{Root: root, Digest: hex.EncodeToString(digest[:]), Files: []model.FileRecord{{
		Path: "README.md", Size: 4, SHA256: hex.EncodeToString(digest[:]), Mode: 0o600,
	}}}
	if err := os.Symlink(target, filepath.Join(root, "README.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := ReadVerifiedFile(item, "README.md", 1024); err == nil {
		t.Fatal("symlink replacement was accepted")
	}
}
