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

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const maxSnapshotBytes int64 = 4 * 1024 * 1024 * 1024

// Snapshot is a private, content-verified copy of the regular files in one
// sealed inventory. Symlinks and scanner-excluded paths are intentionally not
// materialized, so an adapter cannot observe bytes outside its declared
// assessment subject.
type Snapshot struct {
	Path   string `json:"path"`
	Digest string `json:"digest"`
	Files  int    `json:"files"`
	Bytes  int64  `json:"bytes"`
	parent string
}

type snapshotRecord struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

// PrepareSnapshot verifies every inventoried file while copying it into a
// scanner-owned temporary directory. Callers must defer Close after success.
func PrepareSnapshot(item model.Inventory) (*Snapshot, error) {
	return prepareSnapshot(item, nil, nil)
}

// PrepareSnapshotForManifest applies only protocol-owned input remappings that
// are required to keep target files from changing analyzer policy. The original
// bytes remain in the snapshot and are still scanned.
func PrepareSnapshotForManifest(item model.Inventory, manifest Manifest) (*Snapshot, error) {
	if err := manifest.ValidateForCurrentEngine(); err != nil {
		return nil, err
	}
	remapped := map[string]string(nil)
	injected := map[string][]byte(nil)
	if manifest.Protocol == GitleaksProtocolVersion {
		remapped = map[string]string{gitleaksIgnoreSourcePath: gitleaksIgnoreSnapshotPath}
	} else if manifest.Protocol == SyftProtocolVersion {
		injected = map[string][]byte{SyftConfigSnapshotPath: syftConfig}
	} else if manifest.Protocol == GrypeProtocolVersion {
		injected = map[string][]byte{GrypeConfigSnapshotPath: grypeConfig}
	}
	return prepareSnapshot(item, remapped, injected)
}

func prepareSnapshot(item model.Inventory, remapped map[string]string, injected map[string][]byte) (*Snapshot, error) {
	if item.Root == "" || item.FileCount != len(item.Files) {
		return nil, fmt.Errorf("adapter snapshot requires a rooted inventory with a consistent file count")
	}
	if err := inventory.VerifyIdentity(item); err != nil {
		return nil, fmt.Errorf("verify adapter inventory identity: %w", err)
	}
	if err := validateSnapshotRecords(item.Files); err != nil {
		return nil, err
	}

	root, err := os.OpenRoot(item.Root)
	if err != nil {
		return nil, fmt.Errorf("open adapter source root: %w", err)
	}
	defer root.Close()

	path, err := os.MkdirTemp("", "prc-adapter-snapshot-")
	if err != nil {
		return nil, fmt.Errorf("create adapter snapshot: %w", err)
	}
	snapshot := &Snapshot{Path: path, Files: len(item.Files), parent: filepath.Dir(path)}
	removeOnError := true
	defer func() {
		if removeOnError {
			_ = snapshot.Close()
		}
	}()

	destinations := map[string]bool{}
	for _, record := range item.Files {
		if snapshot.Bytes > maxSnapshotBytes-record.Size {
			return nil, fmt.Errorf("adapter snapshot exceeds %d bytes", maxSnapshotBytes)
		}
		destination := record.Path
		if mapped, ok := remapped[record.Path]; ok {
			destination = mapped
		}
		if err := validateRelativePath(destination); err != nil || destinations[destination] {
			return nil, fmt.Errorf("adapter snapshot has an invalid or duplicate materialized path")
		}
		destinations[destination] = true
		if err := copySnapshotFile(root, path, record, destination); err != nil {
			return nil, err
		}
		snapshot.Bytes += record.Size
	}
	for destination, content := range injected {
		if err := validateRelativePath(destination); err != nil || destinations[destination] || len(content) == 0 {
			return nil, fmt.Errorf("adapter snapshot has an invalid or duplicate protocol-owned path")
		}
		destinations[destination] = true
		if err := writeSnapshotOwnedFile(path, destination, content); err != nil {
			return nil, err
		}
	}
	digest, err := snapshotDigest(path)
	if err != nil {
		return nil, err
	}
	snapshot.Digest = digest
	if err := makeSnapshotReadOnly(path); err != nil {
		return nil, err
	}
	removeOnError = false
	return snapshot, nil
}

func writeSnapshotOwnedFile(destinationRoot, destinationPath string, content []byte) error {
	destination := filepath.Join(destinationRoot, filepath.FromSlash(destinationPath))
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return fmt.Errorf("create adapter protocol directory: %w", err)
	}
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o400)
	if err != nil {
		return fmt.Errorf("create adapter protocol file: %w", err)
	}
	written, writeErr := output.Write(content)
	closeErr := output.Close()
	if writeErr != nil {
		return fmt.Errorf("write adapter protocol file: %w", writeErr)
	}
	if written != len(content) {
		return fmt.Errorf("write adapter protocol file: wrote %d of %d bytes", written, len(content))
	}
	if closeErr != nil {
		return fmt.Errorf("close adapter protocol file: %w", closeErr)
	}
	return nil
}

func validateSnapshotRecords(records []model.FileRecord) error {
	seen := make(map[string]bool, len(records))
	for _, record := range records {
		if err := validateRelativePath(record.Path); err != nil {
			return fmt.Errorf("adapter inventory file: %w", err)
		}
		if seen[record.Path] || record.Size < 0 || !hexDigestPattern.MatchString(record.SHA256) {
			return fmt.Errorf("adapter inventory contains a duplicate path or invalid file identity")
		}
		seen[record.Path] = true
	}
	return nil
}

func copySnapshotFile(root *os.Root, destinationRoot string, record model.FileRecord, destinationPath string) error {
	relative := filepath.FromSlash(record.Path)
	before, err := root.Lstat(relative)
	if err != nil || !before.Mode().IsRegular() {
		return fmt.Errorf("adapter source file changed or is not regular: %s", record.Path)
	}
	input, err := root.Open(relative)
	if err != nil {
		return fmt.Errorf("open adapter source file %s: %w", record.Path, err)
	}
	defer input.Close()
	opened, err := input.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
		return fmt.Errorf("adapter source file changed while opening: %s", record.Path)
	}

	destination := filepath.Join(destinationRoot, filepath.FromSlash(destinationPath))
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return fmt.Errorf("create adapter snapshot directory: %w", err)
	}
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return fmt.Errorf("create adapter snapshot file %s: %w", record.Path, err)
	}
	hasher := sha256.New()
	written, copyErr := io.Copy(io.MultiWriter(output, hasher), io.LimitReader(input, record.Size+1))
	closeErr := output.Close()
	if copyErr != nil {
		return fmt.Errorf("copy adapter snapshot file %s: %w", record.Path, copyErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close adapter snapshot file %s: %w", record.Path, closeErr)
	}
	if written != record.Size || hex.EncodeToString(hasher.Sum(nil)) != record.SHA256 {
		return fmt.Errorf("adapter source file changed after inventory: %s", record.Path)
	}
	if err := os.Chmod(destination, 0o400); err != nil {
		return fmt.Errorf("protect adapter snapshot file %s: %w", record.Path, err)
	}
	return nil
}

func snapshotDigest(root string) (string, error) {
	records := []snapshotRecord{}
	var total int64
	rootHandle, err := os.OpenRoot(root)
	if err != nil {
		return "", fmt.Errorf("open adapter snapshot root: %w", err)
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
		relativePath := filepath.FromSlash(relative)
		before, err := rootHandle.Lstat(relativePath)
		if err != nil || !before.Mode().IsRegular() {
			return fmt.Errorf("adapter snapshot contains a non-regular entry")
		}
		if total > maxSnapshotBytes-before.Size() {
			return fmt.Errorf("adapter snapshot exceeds %d bytes", maxSnapshotBytes)
		}
		file, err := rootHandle.Open(relativePath)
		if err != nil {
			return err
		}
		opened, statErr := file.Stat()
		if statErr != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
			_ = file.Close()
			return fmt.Errorf("adapter snapshot changed while opening")
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
			return fmt.Errorf("adapter snapshot changed while hashing")
		}
		records = append(records, snapshotRecord{
			Path: relative, Size: opened.Size(), SHA256: hex.EncodeToString(hasher.Sum(nil)),
		})
		total += opened.Size()
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("hash adapter snapshot: %w", err)
	}
	sort.Slice(records, func(left, right int) bool { return records[left].Path < records[right].Path })
	payload, err := json.Marshal(records)
	if err != nil {
		return "", fmt.Errorf("encode adapter snapshot identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func makeSnapshotReadOnly(root string) error {
	directories := []string{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			directories = append(directories, path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("protect adapter snapshot: %w", err)
	}
	rootHandle, err := os.OpenRoot(root)
	if err != nil {
		return fmt.Errorf("open adapter snapshot for protection: %w", err)
	}
	defer rootHandle.Close()
	sort.Slice(directories, func(left, right int) bool {
		return strings.Count(directories[left], string(filepath.Separator)) >
			strings.Count(directories[right], string(filepath.Separator))
	})
	for _, directory := range directories {
		relative, relativeErr := filepath.Rel(root, directory)
		if relativeErr != nil {
			return fmt.Errorf("resolve adapter snapshot directory: %w", relativeErr)
		}
		if err := rootHandle.Chmod(relative, 0o500); err != nil {
			return fmt.Errorf("protect adapter snapshot directory: %w", err)
		}
	}
	return nil
}

// Close removes only the scanner-owned temporary snapshot. It is idempotent.
func (snapshot *Snapshot) Close() error {
	if snapshot == nil || snapshot.Path == "" {
		return nil
	}
	path := filepath.Clean(snapshot.Path)
	if filepath.Dir(path) != filepath.Clean(snapshot.parent) ||
		!strings.HasPrefix(filepath.Base(path), "prc-adapter-snapshot-") {
		return fmt.Errorf("refusing to remove an invalid adapter snapshot path")
	}
	if rootHandle, openErr := os.OpenRoot(path); openErr == nil {
		_ = filepath.WalkDir(path, func(current string, entry fs.DirEntry, walkErr error) error {
			if walkErr == nil && entry.IsDir() {
				relative, relativeErr := filepath.Rel(path, current)
				if relativeErr == nil {
					_ = rootHandle.Chmod(relative, 0o700)
				}
			}
			return nil
		})
		_ = rootHandle.Close()
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove adapter snapshot: %w", err)
	}
	snapshot.Path = ""
	return nil
}
