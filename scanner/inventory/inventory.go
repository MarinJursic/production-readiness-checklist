package inventory

import (
	"bufio"
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

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const maxEntries = 200_000

var excludedDirectories = map[string]bool{
	".git": true, ".mypy_cache": true, ".prc": true, ".pytest_cache": true,
	".venv": true, "__pycache__": true, "node_modules": true, "site": true,
	"vendor": true,
}

var sourceExtensions = map[string]bool{
	".c": true, ".cc": true, ".cpp": true, ".cs": true, ".go": true,
	".java": true, ".js": true, ".jsx": true, ".kt": true, ".php": true,
	".py": true, ".rb": true, ".rs": true, ".scala": true, ".sh": true,
	".swift": true, ".ts": true, ".tsx": true,
}

var manifests = map[string]string{
	"package.json":     "node",
	"pyproject.toml":   "python",
	"requirements.txt": "python",
	"setup.py":         "python",
	"go.mod":           "go",
	"Cargo.toml":       "rust",
	"pom.xml":          "java",
	"build.gradle":     "java",
	"build.gradle.kts": "java",
}

var lockFiles = map[string]bool{
	"package-lock.json": true, "npm-shrinkwrap.json": true, "pnpm-lock.yaml": true,
	"yarn.lock": true, "uv.lock": true, "Pipfile.lock": true, "poetry.lock": true,
	"go.sum": true, "Cargo.lock": true, "gradle.lockfile": true,
}

type digestInput struct {
	GitCommit         string             `json:"git_commit,omitempty"`
	Files             []model.FileRecord `json:"files"`
	PackageEcosystems []string           `json:"package_ecosystems"`
	Manifests         []string           `json:"manifests"`
	LockFiles         []string           `json:"lock_files"`
}

func Build(target string) (model.Inventory, error) {
	root, err := filepath.Abs(target)
	if err != nil {
		return model.Inventory{}, fmt.Errorf("resolve target: %w", err)
	}
	info, err := os.Stat(root)
	if err != nil {
		return model.Inventory{}, fmt.Errorf("inspect target: %w", err)
	}
	if !info.IsDir() {
		return model.Inventory{}, fmt.Errorf("target is not a directory: %s", root)
	}
	if resolved, err := filepath.EvalSymlinks(root); err == nil {
		root = resolved
	}

	result := model.Inventory{SchemaVersion: model.InventorySchema, Root: root, TargetName: filepath.Base(root)}
	ecosystems := map[string]bool{}
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}
		if path == root {
			return nil
		}
		if entry.IsDir() && excludedDirectories[entry.Name()] {
			return filepath.SkipDir
		}
		if entry.IsDir() {
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		entryInfo, err := entry.Info()
		if err != nil {
			return fmt.Errorf("inspect %s: %w", path, err)
		}
		if !entryInfo.Mode().IsRegular() {
			return nil
		}
		if len(result.Files) >= maxEntries {
			return fmt.Errorf("target exceeds %d files", maxEntries)
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		fileRecord, err := hashFile(path, relative)
		if err != nil {
			return err
		}
		result.Files = append(result.Files, fileRecord)
		name := entry.Name()
		ecosystem, manifest := manifests[name]
		if strings.HasPrefix(name, "requirements") && strings.HasSuffix(name, ".txt") &&
			!strings.HasSuffix(name, ".lock.txt") {
			ecosystem, manifest = "python", true
		}
		if manifest {
			result.Manifests = append(result.Manifests, relative)
			ecosystems[ecosystem] = true
		}
		if lockFiles[name] || strings.HasSuffix(name, ".lock.txt") {
			result.LockFiles = append(result.LockFiles, relative)
		}
		if sourceExtensions[strings.ToLower(filepath.Ext(name))] {
			result.SourceFiles++
		}
		if strings.HasPrefix(relative, ".github/workflows/") &&
			(strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml")) {
			result.CI.GitHubActions = true
		}
		return nil
	})
	if err != nil {
		return model.Inventory{}, err
	}
	sort.Slice(result.Files, func(i, j int) bool { return result.Files[i].Path < result.Files[j].Path })
	sort.Strings(result.Manifests)
	sort.Strings(result.LockFiles)
	for ecosystem := range ecosystems {
		result.PackageEcosystems = append(result.PackageEcosystems, ecosystem)
	}
	sort.Strings(result.PackageEcosystems)
	result.FileCount = len(result.Files)
	result.GitCommit = gitCommit(root)

	payload, err := json.Marshal(digestInput{
		GitCommit: result.GitCommit, Files: result.Files, PackageEcosystems: result.PackageEcosystems,
		Manifests: result.Manifests, LockFiles: result.LockFiles,
	})
	if err != nil {
		return model.Inventory{}, fmt.Errorf("encode inventory digest: %w", err)
	}
	digest := sha256.Sum256(payload)
	result.Digest = hex.EncodeToString(digest[:])
	return result, nil
}

func hashFile(path, relative string) (model.FileRecord, error) {
	pathInfo, err := os.Lstat(path)
	if err != nil {
		return model.FileRecord{}, fmt.Errorf("inspect %s: %w", relative, err)
	}
	if !pathInfo.Mode().IsRegular() {
		return model.FileRecord{}, fmt.Errorf("file changed to a non-regular entry: %s", relative)
	}
	file, err := os.Open(path)
	if err != nil {
		return model.FileRecord{}, fmt.Errorf("open %s: %w", relative, err)
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return model.FileRecord{}, fmt.Errorf("inspect open file %s: %w", relative, err)
	}
	if !openedInfo.Mode().IsRegular() || !os.SameFile(pathInfo, openedInfo) {
		return model.FileRecord{}, fmt.Errorf("file changed while opening: %s", relative)
	}
	hasher := sha256.New()
	size, err := io.Copy(hasher, file)
	if err != nil {
		return model.FileRecord{}, fmt.Errorf("hash %s: %w", relative, err)
	}
	return model.FileRecord{
		Path: relative, Size: size, SHA256: hex.EncodeToString(hasher.Sum(nil)),
	}, nil
}

func gitCommit(root string) string {
	gitDirectory := filepath.Join(root, ".git")
	info, err := os.Stat(gitDirectory)
	if err != nil {
		return ""
	}
	if !info.IsDir() {
		data, err := os.ReadFile(gitDirectory)
		if err != nil {
			return ""
		}
		line := strings.TrimSpace(string(data))
		if !strings.HasPrefix(line, "gitdir: ") {
			return ""
		}
		gitDirectory = strings.TrimPrefix(line, "gitdir: ")
		if !filepath.IsAbs(gitDirectory) {
			gitDirectory = filepath.Join(root, gitDirectory)
		}
	}
	head, err := os.ReadFile(filepath.Join(gitDirectory, "HEAD"))
	if err != nil {
		return ""
	}
	value := strings.TrimSpace(string(head))
	if !strings.HasPrefix(value, "ref: ") {
		return value
	}
	ref := strings.TrimPrefix(value, "ref: ")
	if data, err := os.ReadFile(filepath.Join(gitDirectory, filepath.FromSlash(ref))); err == nil {
		return strings.TrimSpace(string(data))
	}
	file, err := os.Open(filepath.Join(gitDirectory, "packed-refs"))
	if err != nil {
		return ""
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[1] == ref {
			return fields[0]
		}
	}
	return ""
}
