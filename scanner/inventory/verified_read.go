package inventory

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

// ReadVerifiedFile reopens one exact regular file from an immutable inventory,
// enforces a caller-specific bound, and rejects changed bytes or identity. It
// never follows a symlink or accepts a path absent from the inventory.
func ReadVerifiedFile(item model.Inventory, relative string, maximumBytes int64) ([]byte, error) {
	if maximumBytes < 1 || item.Root == "" || item.Digest == "" {
		return nil, fmt.Errorf("verified inventory read requires a root, digest, and positive byte limit")
	}
	clean := filepath.Clean(filepath.FromSlash(relative))
	if clean == "." || filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return nil, fmt.Errorf("verified inventory read path is outside the target")
	}
	canonical := filepath.ToSlash(clean)
	index := sort.Search(len(item.Files), func(index int) bool { return item.Files[index].Path >= canonical })
	if index == len(item.Files) || item.Files[index].Path != canonical {
		return nil, fmt.Errorf("verified inventory read path is not in the sealed inventory")
	}
	record := item.Files[index]
	if record.Size < 0 || record.Size > maximumBytes {
		return nil, fmt.Errorf("verified inventory read exceeds the caller byte limit")
	}
	path := filepath.Join(item.Root, clean)
	rootWithSeparator := filepath.Clean(item.Root) + string(filepath.Separator)
	if path != filepath.Clean(item.Root) && !strings.HasPrefix(filepath.Clean(path), rootWithSeparator) {
		return nil, fmt.Errorf("verified inventory read path escaped the target")
	}
	before, err := os.Lstat(path)
	if err != nil || !before.Mode().IsRegular() || before.Mode()&os.ModeSymlink != 0 || before.Size() != record.Size {
		return nil, fmt.Errorf("verified inventory file identity changed")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open verified inventory file: %w", err)
	}
	defer file.Close()
	opened, err := file.Stat()
	if err != nil || !opened.Mode().IsRegular() || !os.SameFile(before, opened) {
		return nil, fmt.Errorf("verified inventory file changed while opening")
	}
	data, err := io.ReadAll(io.LimitReader(file, maximumBytes+1))
	if err != nil || int64(len(data)) > maximumBytes || int64(len(data)) != record.Size {
		return nil, fmt.Errorf("read verified inventory file failed or exceeded its limit")
	}
	after, err := file.Stat()
	if err != nil || !os.SameFile(opened, after) || opened.Size() != after.Size() || !opened.ModTime().Equal(after.ModTime()) {
		return nil, fmt.Errorf("verified inventory file changed while reading")
	}
	digest := sha256.Sum256(data)
	actual := hex.EncodeToString(digest[:])
	if actual != record.SHA256 {
		return nil, fmt.Errorf("verified inventory file digest changed")
	}
	return data, nil
}
