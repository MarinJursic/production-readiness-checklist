// Package evidenceset verifies a bounded set of independently signed,
// authority-scoped deterministic evidence bundles. The set is only an
// orchestration document: every program policy and every observation remains
// covered by the two signatures required by evidencebundle.
package evidenceset

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidencebundle"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

const (
	SchemaVersion        = "prc.authoritative-evidence-set/v0.1"
	maximumManifestBytes = 1024 * 1024
	maximumBundles       = 6
)

var fileNamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{0,127}$`)

// BundleReference names the three immutable files needed to verify one
// authority. Paths are deliberately sibling file names, not arbitrary paths.
type BundleReference struct {
	Authority             controlprogram.Authority `json:"authority"`
	BundleFile            string                   `json:"bundle_file"`
	PolicySignatureFile   string                   `json:"policy_signature_file"`
	EvidenceSignatureFile string                   `json:"evidence_signature_file"`
}

// Manifest groups at most one bundle for each supported evidence authority.
// The trust store and all referenced files must be non-symlink regular files
// beside this manifest.
type Manifest struct {
	SchemaVersion  string            `json:"schema_version"`
	TrustStoreFile string            `json:"trust_store_file"`
	Bundles        []BundleReference `json:"bundles"`
}

type fileSnapshot struct {
	path string
	info os.FileInfo
}

// VerifyAndEvaluate verifies the complete set before returning any execution.
// This all-or-nothing behavior prevents a malformed later authority from
// silently leaving the caller with a partial authoritative result.
func VerifyAndEvaluate(
	catalog *controlprogramcatalog.Catalog,
	inventory model.Inventory,
	manifestPath string,
	verifiedAt time.Time,
) ([]controlruntime.Execution, []evidencebundle.Verification, error) {
	if catalog == nil {
		return nil, nil, fmt.Errorf("authoritative evidence set requires an exact program catalog")
	}
	manifest, directory, err := load(manifestPath)
	if err != nil {
		return nil, nil, err
	}
	files, err := snapshotFiles(directory, manifest)
	if err != nil {
		return nil, nil, err
	}
	trustStorePath := filepath.Join(directory, manifest.TrustStoreFile)
	executions := make([]controlruntime.Execution, 0, catalog.TemplateCount())
	verifications := make([]evidencebundle.Verification, 0, len(manifest.Bundles))
	seenTemplates := map[string]bool{}
	seenBundleIDs := map[string]bool{}
	for _, reference := range manifest.Bundles {
		items, verification, verifyErr := evidencebundle.VerifyAndEvaluate(
			catalog,
			inventory,
			filepath.Join(directory, reference.BundleFile),
			trustStorePath,
			filepath.Join(directory, reference.PolicySignatureFile),
			filepath.Join(directory, reference.EvidenceSignatureFile),
			verifiedAt,
		)
		if verifyErr != nil {
			return nil, nil, fmt.Errorf("verify %s evidence set entry: %w", reference.Authority, verifyErr)
		}
		if verification.Authority != string(reference.Authority) || seenBundleIDs[verification.BundleID] {
			return nil, nil, fmt.Errorf("evidence set entry %s does not match its declared authority or has a duplicate bundle id", reference.Authority)
		}
		seenBundleIDs[verification.BundleID] = true
		for _, execution := range items {
			if seenTemplates[execution.TemplateID] {
				return nil, nil, fmt.Errorf("evidence set contains duplicate template %s", execution.TemplateID)
			}
			seenTemplates[execution.TemplateID] = true
		}
		executions = append(executions, items...)
		verifications = append(verifications, verification)
	}
	if err := verifyUnchanged(files); err != nil {
		return nil, nil, err
	}
	return executions, verifications, nil
}

func load(path string) (Manifest, string, error) {
	data, absolute, err := readManifest(path)
	if err != nil {
		return Manifest{}, "", err
	}
	var manifest Manifest
	if err := provider.DecodeStrictJSON(data, &manifest); err != nil {
		return Manifest{}, "", fmt.Errorf("decode authoritative evidence set: %w", err)
	}
	if err := validate(manifest); err != nil {
		return Manifest{}, "", err
	}
	return manifest, filepath.Dir(absolute), nil
}

func validate(manifest Manifest) error {
	if manifest.SchemaVersion != SchemaVersion || !validFileName(manifest.TrustStoreFile) ||
		len(manifest.Bundles) < 1 || len(manifest.Bundles) > maximumBundles {
		return fmt.Errorf("authoritative evidence set has an invalid envelope")
	}
	seenAuthorities := map[controlprogram.Authority]bool{}
	seenFiles := map[string]bool{manifest.TrustStoreFile: true}
	previous := ""
	for _, reference := range manifest.Bundles {
		authority := string(reference.Authority)
		if !validAuthority(reference.Authority) || seenAuthorities[reference.Authority] ||
			(previous != "" && previous >= authority) {
			return fmt.Errorf("authoritative evidence set entries must use unique authorities in ascending order")
		}
		for _, name := range []string{reference.BundleFile, reference.PolicySignatureFile, reference.EvidenceSignatureFile} {
			if !validFileName(name) || seenFiles[name] {
				return fmt.Errorf("authoritative evidence set references an invalid or duplicate sibling file")
			}
			seenFiles[name] = true
		}
		seenAuthorities[reference.Authority] = true
		previous = authority
	}
	return nil
}

func validAuthority(authority controlprogram.Authority) bool {
	switch authority {
	case controlprogram.AuthorityArtifact, controlprogram.AuthorityEnvironment,
		controlprogram.AuthorityExecuted, controlprogram.AuthorityExternalRegistry,
		controlprogram.AuthorityRepository, controlprogram.AuthorityStructuredRecord:
		return true
	default:
		return false
	}
}

func validFileName(value string) bool {
	return fileNamePattern.MatchString(value) && filepath.Base(value) == value &&
		value != "." && value != ".." && !strings.ContainsAny(value, `/\\`)
}

func snapshotFiles(directory string, manifest Manifest) ([]fileSnapshot, error) {
	names := []string{manifest.TrustStoreFile}
	for _, reference := range manifest.Bundles {
		names = append(names, reference.BundleFile, reference.PolicySignatureFile, reference.EvidenceSignatureFile)
	}
	snapshots := make([]fileSnapshot, 0, len(names))
	for _, name := range names {
		path := filepath.Join(directory, name)
		info, err := os.Lstat(path)
		if err != nil {
			return nil, fmt.Errorf("inspect authoritative evidence set file %s: %w", name, err)
		}
		if info.Mode()&os.ModeSymlink != 0 || !info.Mode().IsRegular() {
			return nil, fmt.Errorf("authoritative evidence set file %s must be a non-symlink regular file", name)
		}
		for _, existing := range snapshots {
			if os.SameFile(existing.info, info) {
				return nil, fmt.Errorf("authoritative evidence set files %s and %s refer to the same file", filepath.Base(existing.path), name)
			}
		}
		snapshots = append(snapshots, fileSnapshot{path: path, info: info})
	}
	return snapshots, nil
}

func verifyUnchanged(snapshots []fileSnapshot) error {
	for _, snapshot := range snapshots {
		current, err := os.Lstat(snapshot.path)
		if err != nil || current.Mode()&os.ModeSymlink != 0 || !current.Mode().IsRegular() ||
			!os.SameFile(snapshot.info, current) || snapshot.info.Size() != current.Size() ||
			!snapshot.info.ModTime().Equal(current.ModTime()) {
			return fmt.Errorf("authoritative evidence set file %s changed during verification", filepath.Base(snapshot.path))
		}
	}
	return nil
}

func readManifest(path string) ([]byte, string, error) {
	if strings.TrimSpace(path) == "" {
		return nil, "", fmt.Errorf("authoritative evidence set path is empty")
	}
	absolute, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("resolve authoritative evidence set: %w", err)
	}
	before, err := os.Lstat(absolute)
	if err != nil {
		return nil, "", fmt.Errorf("inspect authoritative evidence set: %w", err)
	}
	if before.Mode()&os.ModeSymlink != 0 || !before.Mode().IsRegular() || before.Size() > maximumManifestBytes {
		return nil, "", fmt.Errorf("authoritative evidence set must be a non-symlink regular file no larger than 1 MiB")
	}
	file, err := os.Open(absolute)
	if err != nil {
		return nil, "", fmt.Errorf("open authoritative evidence set: %w", err)
	}
	defer file.Close()
	after, err := file.Stat()
	if err != nil || !after.Mode().IsRegular() || !os.SameFile(before, after) {
		return nil, "", fmt.Errorf("authoritative evidence set changed while opening")
	}
	data, err := io.ReadAll(io.LimitReader(file, maximumManifestBytes+1))
	if err != nil || len(data) > maximumManifestBytes {
		return nil, "", fmt.Errorf("read authoritative evidence set: file exceeds 1 MiB")
	}
	final, err := file.Stat()
	if err != nil || !os.SameFile(after, final) || after.Size() != final.Size() ||
		!after.ModTime().Equal(final.ModTime()) || int64(len(data)) != final.Size() {
		return nil, "", fmt.Errorf("authoritative evidence set changed while reading")
	}
	return data, absolute, nil
}
