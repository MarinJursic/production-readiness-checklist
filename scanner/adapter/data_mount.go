package adapter

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BoundDataMount is the sealed identity of one external read-only input. Its
// source path is used only by the local execution plan and is not copied into
// the durable adapter execution record.
type BoundDataMount struct {
	Name        string `json:"name"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	SHA256      string `json:"sha256"`
	Files       int    `json:"files"`
	Bytes       int64  `json:"bytes"`
	MaxFiles    int    `json:"max_files"`
	MaxBytes    int64  `json:"max_bytes"`
}

func bindDataMounts(manifest Manifest, sources map[string]string) ([]BoundDataMount, error) {
	if len(sources) != len(manifest.DataMounts) {
		return nil, fmt.Errorf("adapter %s requires exactly %d data mounts", manifest.ID, len(manifest.DataMounts))
	}
	bound := make([]BoundDataMount, 0, len(manifest.DataMounts))
	for _, declaration := range manifest.DataMounts {
		source, ok := sources[declaration.Name]
		if !ok || strings.TrimSpace(source) == "" {
			return nil, fmt.Errorf("adapter %s requires data mount %q", manifest.ID, declaration.Name)
		}
		absolute, err := filepath.Abs(source)
		if err != nil {
			return nil, fmt.Errorf("resolve adapter data mount %q: %w", declaration.Name, err)
		}
		if resolved, resolveErr := filepath.EvalSymlinks(absolute); resolveErr == nil {
			absolute = resolved
		}
		if strings.Contains(absolute, ",") {
			return nil, fmt.Errorf("adapter data mount %q path cannot contain a comma", declaration.Name)
		}
		info, err := os.Stat(absolute)
		if err != nil || !info.IsDir() {
			return nil, fmt.Errorf("adapter data mount %q is not an accessible directory", declaration.Name)
		}
		digest, files, bytes, err := dataDirectoryIdentity(absolute, declaration.MaxFiles, declaration.MaxBytes)
		if err != nil {
			return nil, fmt.Errorf("bind adapter data mount %q: %w", declaration.Name, err)
		}
		bound = append(bound, BoundDataMount{
			Name: declaration.Name, Source: absolute, Destination: declaration.Destination,
			SHA256: digest, Files: files, Bytes: bytes,
			MaxFiles: declaration.MaxFiles, MaxBytes: declaration.MaxBytes,
		})
	}
	for name := range sources {
		found := false
		for _, declaration := range manifest.DataMounts {
			if declaration.Name == name {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("adapter %s does not declare data mount %q", manifest.ID, name)
		}
	}
	return bound, nil
}

func verifyBoundDataMounts(mounts []BoundDataMount) error {
	for _, mount := range mounts {
		digest, files, bytes, err := dataDirectoryIdentity(mount.Source, mount.MaxFiles, mount.MaxBytes)
		if err != nil || digest != mount.SHA256 || files != mount.Files || bytes != mount.Bytes {
			return fmt.Errorf("adapter data mount %q changed after execution planning", mount.Name)
		}
	}
	return nil
}

func dataDirectoryIdentity(root string, maxFiles int, maxBytes int64) (string, int, int64, error) {
	records := []snapshotRecord{}
	var total int64
	rootHandle, err := os.OpenRoot(root)
	if err != nil {
		return "", 0, 0, fmt.Errorf("open data directory: %w", err)
	}
	defer rootHandle.Close()
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root || entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if err := validateRelativePath(relative); err != nil {
			return err
		}
		if len(records) >= maxFiles {
			return fmt.Errorf("data directory exceeds %d files", maxFiles)
		}
		relativePath := filepath.FromSlash(relative)
		before, err := rootHandle.Lstat(relativePath)
		if err != nil || !before.Mode().IsRegular() {
			return fmt.Errorf("data directory contains a non-regular entry")
		}
		if total > maxBytes-before.Size() {
			return fmt.Errorf("data directory exceeds %d bytes", maxBytes)
		}
		file, err := rootHandle.Open(relativePath)
		if err != nil {
			return err
		}
		opened, statErr := file.Stat()
		if statErr != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
			_ = file.Close()
			return fmt.Errorf("data directory changed while opening")
		}
		hasher := sha256.New()
		written, copyErr := io.Copy(hasher, io.LimitReader(file, opened.Size()+1))
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
		if written != opened.Size() {
			return fmt.Errorf("data directory changed while hashing")
		}
		records = append(records, snapshotRecord{
			Path: relative, Size: opened.Size(), SHA256: hex.EncodeToString(hasher.Sum(nil)),
		})
		total += opened.Size()
		return nil
	})
	if err != nil {
		return "", 0, 0, err
	}
	if len(records) == 0 {
		return "", 0, 0, fmt.Errorf("data directory is empty")
	}
	sort.Slice(records, func(left, right int) bool { return records[left].Path < records[right].Path })
	payload, err := json.Marshal(records)
	if err != nil {
		return "", 0, 0, fmt.Errorf("encode data directory identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), len(records), total, nil
}
