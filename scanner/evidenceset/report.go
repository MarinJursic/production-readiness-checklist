package evidenceset

import (
	"fmt"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogramcatalog"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlruntime"
	"github.com/MarinJursic/production-readiness-checklist/scanner/evidencebundle"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const VerificationReportSchema = "prc.evidence-set-verification/v0.1"

type OutcomeCounts struct {
	Passed        int `json:"passed"`
	Failed        int `json:"failed"`
	NotApplicable int `json:"not_applicable"`
	Blocked       int `json:"blocked"`
}

type AuthorityVerification struct {
	Authority          string        `json:"authority"`
	BundleID           string        `json:"bundle_id"`
	BundleSHA256       string        `json:"bundle_sha256"`
	PolicySHA256       string        `json:"policy_sha256"`
	EntryCount         int           `json:"entry_count"`
	Outcomes           OutcomeCounts `json:"outcomes"`
	PolicyKeyID        string        `json:"policy_key_id"`
	EvidenceKeyID      string        `json:"evidence_key_id"`
	TrustStoreID       string        `json:"trust_store_id"`
	TrustStoreSHA256   string        `json:"trust_store_sha256"`
	SignaturesVerified bool          `json:"signatures_verified"`
}

// VerificationReport keeps signature validity separate from predicate
// outcomes. A cryptographically valid bundle can contain a failing or blocked
// observation, and neither state is rewritten as an import failure.
type VerificationReport struct {
	SchemaVersion             string                  `json:"schema_version"`
	VerifiedAt                time.Time               `json:"verified_at"`
	CatalogSHA256             string                  `json:"catalog_sha256"`
	InventorySHA256           string                  `json:"inventory_sha256"`
	CryptographicallyVerified bool                    `json:"cryptographically_verified"`
	BundleCount               int                     `json:"bundle_count"`
	EntryCount                int                     `json:"entry_count"`
	Outcomes                  OutcomeCounts           `json:"outcomes"`
	Authorities               []AuthorityVerification `json:"authorities"`
}

func SummarizeVerification(
	catalog *controlprogramcatalog.Catalog,
	inventory model.Inventory,
	executions []controlruntime.Execution,
	verifications []evidencebundle.Verification,
	verifiedAt time.Time,
) (VerificationReport, error) {
	if catalog == nil || inventory.Digest == "" || verifiedAt.IsZero() || len(verifications) == 0 {
		return VerificationReport{}, fmt.Errorf("evidence-set verification summary is incomplete")
	}
	byTemplate := make(map[string]controlruntime.Execution, len(executions))
	for _, execution := range executions {
		if !execution.Authenticated() || byTemplate[execution.TemplateID].TemplateID != "" {
			return VerificationReport{}, fmt.Errorf("evidence-set verification contains an unauthenticated or duplicate execution")
		}
		byTemplate[execution.TemplateID] = execution
	}
	report := VerificationReport{
		SchemaVersion: VerificationReportSchema, VerifiedAt: verifiedAt.UTC(), CatalogSHA256: catalog.Digest(),
		InventorySHA256: inventory.Digest, CryptographicallyVerified: true,
		BundleCount: len(verifications), Authorities: make([]AuthorityVerification, 0, len(verifications)),
	}
	seenTemplates := map[string]bool{}
	seenAuthorities := map[string]bool{}
	trustStoreID := ""
	trustStoreSHA256 := ""
	for _, verification := range verifications {
		if !evidencebundle.SupportsAuthority(controlprogram.Authority(verification.Authority)) || seenAuthorities[verification.Authority] {
			return VerificationReport{}, fmt.Errorf("evidence-set verification contains an unsupported or duplicate authority")
		}
		seenAuthorities[verification.Authority] = true
		authority := AuthorityVerification{
			Authority: verification.Authority, BundleID: verification.BundleID,
			BundleSHA256: verification.BundleSHA256, PolicySHA256: verification.PolicySHA256,
			EntryCount: verification.EntryCount, PolicyKeyID: verification.PolicySignature.KeyID,
			EvidenceKeyID:      verification.EvidenceSignature.KeyID,
			TrustStoreID:       verification.PolicySignature.TrustStoreID,
			TrustStoreSHA256:   verification.PolicySignature.TrustStoreDigest,
			SignaturesVerified: verification.PolicySignature.Verified && verification.EvidenceSignature.Verified,
		}
		if !authority.SignaturesVerified || verification.PolicySignature.TrustStoreID != verification.EvidenceSignature.TrustStoreID ||
			verification.PolicySignature.TrustStoreDigest != verification.EvidenceSignature.TrustStoreDigest ||
			verification.CatalogSHA256 != report.CatalogSHA256 || verification.InventorySHA256 != report.InventorySHA256 ||
			verification.EntryCount != len(verification.Entries) {
			return VerificationReport{}, fmt.Errorf("evidence-set verification record is inconsistent")
		}
		if trustStoreID == "" {
			trustStoreID, trustStoreSHA256 = authority.TrustStoreID, authority.TrustStoreSHA256
		} else if authority.TrustStoreID != trustStoreID || authority.TrustStoreSHA256 != trustStoreSHA256 {
			return VerificationReport{}, fmt.Errorf("evidence-set verification records use different trust stores")
		}
		for _, entry := range verification.Entries {
			execution, ok := byTemplate[entry.TemplateID]
			if !ok || seenTemplates[entry.TemplateID] || entry.Outcome != string(execution.Outcome) ||
				entry.ReasonCode != string(execution.ReasonCode) {
				return VerificationReport{}, fmt.Errorf("evidence-set verification entry %s is inconsistent", entry.TemplateID)
			}
			seenTemplates[entry.TemplateID] = true
			incrementOutcome(&authority.Outcomes, execution.Outcome)
			incrementOutcome(&report.Outcomes, execution.Outcome)
		}
		report.EntryCount += verification.EntryCount
		report.Authorities = append(report.Authorities, authority)
	}
	if report.EntryCount != len(executions) || len(seenTemplates) != len(executions) {
		return VerificationReport{}, fmt.Errorf("evidence-set verification did not account for every execution")
	}
	return report, nil
}

func incrementOutcome(counts *OutcomeCounts, outcome controlprogram.Outcome) {
	switch outcome {
	case controlprogram.OutcomePass:
		counts.Passed++
	case controlprogram.OutcomeFail:
		counts.Failed++
	case controlprogram.OutcomeNotApplicable:
		counts.NotApplicable++
	default:
		counts.Blocked++
	}
}
