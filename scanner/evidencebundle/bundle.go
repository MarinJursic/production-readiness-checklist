// Package evidencebundle verifies offline, authority-scoped deterministic
// evidence bundles. A bundle needs two different trusted signatures over the
// same bytes: one for scanner policy/program inputs and one for observations
// from the required evidence authority.
package evidencebundle

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/trust"
)

const (
	SchemaVersion      = "prc.authoritative-evidence-bundle/v0.1"
	PolicySchema       = "prc.authoritative-evidence-policy/v0.1"
	VerificationSchema = "prc.authoritative-evidence-verification/v0.1"
	maximumBundleBytes = 32 * 1024 * 1024
	maximumEntries     = 765
)

var (
	bundleIDPattern         = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9./@_-]{0,255}$`)
	providerPattern         = regexp.MustCompile(`^[a-z0-9][a-z0-9._/@:-]{0,255}$`)
	digestPattern           = regexp.MustCompile(`^[0-9a-f]{64}$`)
	evidenceKindByAuthority = map[controlprogram.Authority]string{
		controlprogram.AuthorityRepository:       "control-evidence-repository",
		controlprogram.AuthorityArtifact:         "control-evidence-artifact",
		controlprogram.AuthorityExecuted:         "control-evidence-executed",
		controlprogram.AuthorityEnvironment:      "control-evidence-environment",
		controlprogram.AuthorityExternalRegistry: "control-evidence-external-registry",
		controlprogram.AuthorityStructuredRecord: "control-evidence-structured-record",
	}
)

// Entry binds one materialized reviewed program to the evidence observed for
// it. ProviderID identifies the actual collector, not a compliance verdict.
type Entry struct {
	TemplateID string                  `json:"template_id"`
	ProviderID string                  `json:"provider_id"`
	Program    controlprogram.Program  `json:"program"`
	Evidence   controlprogram.Evidence `json:"evidence"`
}

// Bundle contains entries from exactly one evidence authority. This permits a
// trust-store key to be least-privilege for repository, artifact, executed,
// environment, registry, or structured-record observations.
type Bundle struct {
	SchemaVersion   string                   `json:"schema_version"`
	ID              string                   `json:"id"`
	CatalogSHA256   string                   `json:"catalog_sha256"`
	InventorySHA256 string                   `json:"inventory_sha256"`
	Authority       controlprogram.Authority `json:"authority"`
	Entries         []Entry                  `json:"entries"`
}

// Verification is retained in the scan result so imported evidence can be
// audited without relying on a terminal message.
type Verification struct {
	SchemaVersion     string                             `json:"schema_version"`
	BundleID          string                             `json:"bundle_id"`
	BundleSHA256      string                             `json:"bundle_sha256"`
	PolicySHA256      string                             `json:"policy_sha256"`
	CatalogSHA256     string                             `json:"catalog_sha256"`
	InventorySHA256   string                             `json:"inventory_sha256"`
	Authority         string                             `json:"authority"`
	EntryCount        int                                `json:"entry_count"`
	Entries           []model.AuthoritativeEvidenceEntry `json:"entries"`
	PolicySignature   trust.Verification                 `json:"policy_signature"`
	EvidenceSignature trust.Verification                 `json:"evidence_signature"`
}

type loadedBundle struct {
	bundle Bundle
	digest string
}

type policyEntry struct {
	TemplateID string                 `json:"template_id"`
	ProviderID string                 `json:"provider_id"`
	Program    controlprogram.Program `json:"program"`
}

type policyDocument struct {
	SchemaVersion   string                   `json:"schema_version"`
	ID              string                   `json:"id"`
	CatalogSHA256   string                   `json:"catalog_sha256"`
	InventorySHA256 string                   `json:"inventory_sha256"`
	Authority       controlprogram.Authority `json:"authority"`
	Entries         []policyEntry            `json:"entries"`
}

// VerifyAndEvaluate loads one immutable bundle, verifies two independent
// signatures, rechecks every program against the current catalog, and returns
// scanner-authenticated executions. It never runs bundle-provided code.
func VerifyAndEvaluate(
	catalog *controlprogramcatalog.Catalog,
	inventory model.Inventory,
	bundlePath, trustStorePath, policySignaturePath, evidenceSignaturePath string,
	verifiedAt time.Time,
) ([]controlruntime.Execution, Verification, error) {
	if catalog == nil || inventory.SchemaVersion != model.InventorySchema || inventory.Digest == "" || verifiedAt.IsZero() {
		return nil, Verification{}, fmt.Errorf("authoritative evidence requires a catalog, sealed inventory, and verification time")
	}
	loaded, err := load(bundlePath)
	if err != nil {
		return nil, Verification{}, err
	}
	if loaded.bundle.CatalogSHA256 != catalog.Digest() || loaded.bundle.InventorySHA256 != inventory.Digest {
		return nil, Verification{}, fmt.Errorf("authoritative evidence bundle does not match the current catalog and inventory")
	}
	store, err := trust.LoadStore(trustStorePath)
	if err != nil {
		return nil, Verification{}, err
	}
	policySignature, err := trust.LoadSignature(policySignaturePath)
	if err != nil {
		return nil, Verification{}, err
	}
	evidenceSignature, err := trust.LoadSignature(evidenceSignaturePath)
	if err != nil {
		return nil, Verification{}, err
	}
	if policySignature.KeyID == evidenceSignature.KeyID {
		return nil, Verification{}, fmt.Errorf("policy and evidence bundles require different trusted keys")
	}
	policyDigest, err := PolicySHA256(loaded.bundle)
	if err != nil {
		return nil, Verification{}, err
	}
	policyVerification, err := trust.Verify(
		store, policySignature, "control-policy-bundle", loaded.bundle.ID, policyDigest, verifiedAt.UTC(),
	)
	if err != nil {
		return nil, Verification{}, fmt.Errorf("verify authoritative policy signature: %w", err)
	}
	evidenceKind := evidenceKindByAuthority[loaded.bundle.Authority]
	evidenceVerification, err := trust.Verify(
		store, evidenceSignature, evidenceKind, loaded.bundle.ID, loaded.digest, verifiedAt.UTC(),
	)
	if err != nil {
		return nil, Verification{}, fmt.Errorf("verify authoritative evidence signature: %w", err)
	}

	templates := make(map[string]controlprogramcatalog.Template, catalog.TemplateCount())
	for _, template := range catalog.Templates() {
		templates[template.TemplateID] = template
	}
	executions := make([]controlruntime.Execution, 0, len(loaded.bundle.Entries))
	verifiedEntries := make([]model.AuthoritativeEvidenceEntry, 0, len(loaded.bundle.Entries))
	for _, entry := range loaded.bundle.Entries {
		template, ok := templates[entry.TemplateID]
		if !ok || template.RequiredAuthority != loaded.bundle.Authority {
			return nil, Verification{}, fmt.Errorf("authoritative evidence references an unknown or wrong-authority template %s", entry.TemplateID)
		}
		if entry.Program.InventorySHA256 != inventory.Digest || entry.Evidence.InventorySHA256 != inventory.Digest {
			return nil, Verification{}, fmt.Errorf("authoritative evidence template %s does not bind the current inventory", entry.TemplateID)
		}
		// The policy must have been sealed before collection, and the evidence
		// signer must not attest observations it could not yet have seen.
		if policySignature.IssuedAt.After(entry.Evidence.ObservedAt) {
			return nil, Verification{}, fmt.Errorf("authoritative evidence template %s was observed before its policy was sealed", entry.TemplateID)
		}
		if evidenceSignature.IssuedAt.Before(entry.Evidence.ObservedAt) {
			return nil, Verification{}, fmt.Errorf("authoritative evidence template %s was observed after its evidence signature", entry.TemplateID)
		}
		execution, err := controlruntime.EvaluateAuthenticated(
			template, entry.Program, entry.Evidence,
			entry.ProviderID+"@"+evidenceSignature.KeyID, verifiedAt.UTC(),
		)
		if err != nil {
			return nil, Verification{}, fmt.Errorf("evaluate authoritative evidence template %s: %w", entry.TemplateID, err)
		}
		executions = append(executions, execution)
		verifiedEntries = append(verifiedEntries, model.AuthoritativeEvidenceEntry{
			TemplateID: entry.TemplateID, ProviderID: entry.ProviderID, Program: entry.Program,
			EvidenceSHA256: execution.EvidenceSHA256, Outcome: string(execution.Outcome),
			ReasonCode: string(execution.ReasonCode),
		})
	}
	return executions, Verification{
		SchemaVersion: VerificationSchema, BundleID: loaded.bundle.ID, BundleSHA256: loaded.digest,
		PolicySHA256:  policyDigest,
		CatalogSHA256: loaded.bundle.CatalogSHA256, InventorySHA256: loaded.bundle.InventorySHA256,
		Authority: string(loaded.bundle.Authority), EntryCount: len(loaded.bundle.Entries),
		Entries:         verifiedEntries,
		PolicySignature: policyVerification, EvidenceSignature: evidenceVerification,
	}, nil
}

func load(path string) (loadedBundle, error) {
	data, err := readBoundedFile(path)
	if err != nil {
		return loadedBundle{}, err
	}
	var bundle Bundle
	if err := provider.DecodeStrictJSON(data, &bundle); err != nil {
		return loadedBundle{}, fmt.Errorf("decode authoritative evidence bundle: %w", err)
	}
	if err := validateBundle(bundle); err != nil {
		return loadedBundle{}, err
	}
	digest, err := BundleSHA256(bundle)
	if err != nil {
		return loadedBundle{}, err
	}
	return loadedBundle{bundle: bundle, digest: digest}, nil
}

func validateBundle(bundle Bundle) error {
	if err := validatePolicyInputs(bundle); err != nil {
		return err
	}
	for _, entry := range bundle.Entries {
		if entry.Evidence.Authority != bundle.Authority || entry.Evidence.InventorySHA256 != bundle.InventorySHA256 {
			return fmt.Errorf("authoritative evidence entry %s has an inconsistent evidence authority or inventory", entry.TemplateID)
		}
		if err := controlprogram.ValidateEvidence(entry.Evidence); err != nil {
			return fmt.Errorf("authoritative evidence entry %s has invalid evidence: %w", entry.TemplateID, err)
		}
	}
	return nil
}

func validatePolicyInputs(bundle Bundle) error {
	if bundle.SchemaVersion != SchemaVersion || !bundleIDPattern.MatchString(bundle.ID) ||
		!digestPattern.MatchString(bundle.CatalogSHA256) || !digestPattern.MatchString(bundle.InventorySHA256) ||
		evidenceKindByAuthority[bundle.Authority] == "" || len(bundle.Entries) < 1 || len(bundle.Entries) > maximumEntries {
		return fmt.Errorf("authoritative evidence bundle has an invalid envelope")
	}
	seen := map[string]bool{}
	previous := ""
	for _, entry := range bundle.Entries {
		if !digestPattern.MatchString(entry.TemplateID) || seen[entry.TemplateID] ||
			(previous != "" && previous >= entry.TemplateID) || !providerPattern.MatchString(entry.ProviderID) {
			return fmt.Errorf("authoritative evidence entries must be unique and ordered by template_id")
		}
		if entry.Program.RequiredAuthority != bundle.Authority || entry.Program.InventorySHA256 != bundle.InventorySHA256 {
			return fmt.Errorf("authoritative evidence entry %s has an inconsistent policy authority or inventory", entry.TemplateID)
		}
		if err := controlprogram.ValidateProgram(entry.Program); err != nil {
			return fmt.Errorf("authoritative evidence entry %s has an invalid program: %w", entry.TemplateID, err)
		}
		seen[entry.TemplateID], previous = true, entry.TemplateID
	}
	return nil
}

// PolicySHA256 returns the stable digest approved before collection. It binds
// the reviewed programs and every runtime input but deliberately excludes the
// observations, which do not exist yet. The evidence signature separately
// binds the complete raw bundle after collection.
func PolicySHA256(bundle Bundle) (string, error) {
	if err := validatePolicyInputs(bundle); err != nil {
		return "", err
	}
	document := policyDocument{
		SchemaVersion: PolicySchema, ID: bundle.ID, CatalogSHA256: bundle.CatalogSHA256,
		InventorySHA256: bundle.InventorySHA256, Authority: bundle.Authority,
		Entries: make([]policyEntry, 0, len(bundle.Entries)),
	}
	for _, entry := range bundle.Entries {
		document.Entries = append(document.Entries, policyEntry{
			TemplateID: entry.TemplateID, ProviderID: entry.ProviderID, Program: entry.Program,
		})
	}
	payload, err := json.Marshal(document)
	if err != nil {
		return "", fmt.Errorf("encode authoritative evidence policy identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

// BundleSHA256 returns the stable digest attested by the evidence signer. It
// hashes the strict decoded value rather than presentation whitespace, so the
// same typed bundle has one identity in every saved run.
func BundleSHA256(bundle Bundle) (string, error) {
	if err := validateBundle(bundle); err != nil {
		return "", err
	}
	payload, err := json.Marshal(bundle)
	if err != nil {
		return "", fmt.Errorf("encode authoritative evidence bundle identity: %w", err)
	}
	digest := sha256.Sum256(payload)
	return hex.EncodeToString(digest[:]), nil
}

func readBoundedFile(path string) ([]byte, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("authoritative evidence bundle path is empty")
	}
	before, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect authoritative evidence bundle: %w", err)
	}
	if before.Mode()&os.ModeSymlink != 0 || !before.Mode().IsRegular() || before.Size() > maximumBundleBytes {
		return nil, fmt.Errorf("authoritative evidence bundle must be a non-symlink regular file no larger than 32 MiB")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open authoritative evidence bundle: %w", err)
	}
	defer file.Close()
	after, err := file.Stat()
	if err != nil || !after.Mode().IsRegular() || !os.SameFile(before, after) {
		return nil, fmt.Errorf("authoritative evidence bundle changed while opening")
	}
	data, err := io.ReadAll(io.LimitReader(file, maximumBundleBytes+1))
	if err != nil || len(data) > maximumBundleBytes {
		return nil, fmt.Errorf("read authoritative evidence bundle: file exceeds 32 MiB")
	}
	final, err := file.Stat()
	if err != nil || !os.SameFile(after, final) || after.Size() != final.Size() ||
		!after.ModTime().Equal(final.ModTime()) || int64(len(data)) != final.Size() {
		return nil, fmt.Errorf("authoritative evidence bundle changed while reading")
	}
	return data, nil
}

// SortEntries gives bundle producers the canonical entry order required by
// the verifier without changing any entry content.
func SortEntries(entries []Entry) {
	sort.Slice(entries, func(left, right int) bool { return entries[left].TemplateID < entries[right].TemplateID })
}
