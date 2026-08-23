package adapter

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func validGrypeDocument(now time.Time) grypeDocument {
	captured := now.Add(-2 * time.Hour).Format(time.RFC3339)
	metadata := grypeMetadata{
		ID: "GHSA-35jh-r3h4-6jhm", DataSource: "https://github.com/advisories/GHSA-35jh-r3h4-6jhm",
		Namespace: "github:language:javascript", Severity: "High", URLs: []string{},
		Description: "Command injection in a fixture package.",
		CVSS: []grypeCVSS{{
			Type: "Secondary", Version: "3.1", Vector: "CVSS:3.1/AV:N/AC:L/PR:H/UI:N/S:U/C:H/I:H/A:H",
			Metrics: grypeCVSSMetrics{BaseScore: 7.2}, VendorMetadata: json.RawMessage(`{}`),
		}},
		KnownExploited: []grypeKnownExploited{{
			CVE: "CVE-2021-23337", DateAdded: now.Add(-24 * time.Hour).Format(time.DateOnly),
			DueDate: now.Add(24 * time.Hour).Format(time.DateOnly), KnownRansomwareCampaignUse: "Known",
		}},
		EPSS: []grypeEPSS{{
			CVE: "CVE-2021-23337", EPSS: 0.21, Percentile: 0.97, Date: now.Add(-24 * time.Hour).Format(time.DateOnly),
		}},
		CWEs: []grypeCWE{},
	}
	return grypeDocument{
		Matches: []grypeMatch{{
			Vulnerability: grypeVulnerability{
				grypeMetadata: metadata,
				Fix:           grypeFix{Versions: []string{"4.17.21"}, State: "fixed"},
				Advisories:    []grypeAdvisory{}, Risk: 15.67,
			},
			RelatedVulnerabilities: []grypeMetadata{{
				ID: "CVE-2021-23337", DataSource: "https://nvd.nist.gov/vuln/detail/CVE-2021-23337",
				Namespace: "nvd:cpe", Severity: "Critical", URLs: []string{},
				Description: "Canonical CVE description.", CVSS: []grypeCVSS{},
			}},
			MatchDetails: []grypeMatchDetail{{
				Type: "exact-direct-match", Matcher: "javascript-matcher",
				SearchedBy: json.RawMessage(`{"package":{"name":"lodash","version":"4.17.15"}}`),
				Found:      json.RawMessage(`{"vulnerabilityID":"GHSA-35jh-r3h4-6jhm"}`),
			}},
			Artifact: grypePackage{
				ID: "fixture-package", Name: "lodash", Version: "4.17.15", Type: "npm",
				Locations: []grypePackageLocation{{
					Path: "/package-lock.json", AccessPath: "/package-lock.json",
					Annotations: map[string]string{"evidence": "primary"},
				}},
				Language: "javascript", Licenses: []string{}, CPEs: []string{},
				PURL: "pkg:npm/lodash@4.17.15", Upstreams: []grypeUpstream{},
			},
		}},
		Source: &grypeSource{Type: "directory", Target: json.RawMessage(`"/workspace"`)},
		Distro: grypeDistro{Name: "", Version: "", IDLike: nil},
		Descriptor: grypeDescriptor{
			Name: "grype", Version: GrypeToolVersion,
			Configuration: map[string]any{
				"check-for-app-update": false, "file": "", "output": []any{"json"},
				"externalSources": map[string]any{"enable": false},
				"db": map[string]any{
					"cache-dir": grypeDataMountPath, "auto-update": false,
					"validate-by-hash-on-start": true, "validate-age": true,
					"require-update-check": false, "max-allowed-built-age": int64(grypeMaxDatabaseAge),
				},
			},
			DB: grypeDatabaseDescriptor{
				Status: grypeDatabaseStatus{
					SchemaVersion: "v6.1.9",
					From:          "https://grype.anchore.io/databases/v6/vulnerability-db_v6.1.9_2026-08-23T00:17:31Z_1787465727.tar.zst?checksum=sha256%3A" + strings.Repeat("a", 64),
					Built:         now.Add(-3 * time.Hour).Format(time.RFC3339),
					Path:          grypeDataMountPath + "/6/vulnerability.db", Valid: true,
				},
				Providers: map[string]grypeDatabaseProvider{
					"epss": {Captured: captured, Input: "xxh64:1234567890abcdef"},
					"kev":  {Captured: captured, Input: "xxh64:234567890abcdef1"},
					"nvd":  {Captured: captured, Input: "xxh64:34567890abcdef12"},
				},
			},
			Timestamp: now.Format(time.RFC3339Nano),
		},
	}
}

func marshalGrypeDocument(t *testing.T, document grypeDocument) []byte {
	t.Helper()
	payload, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

func TestCheckedInGrypeManifestLoads(t *testing.T) {
	manifest, err := LoadManifest(filepath.Join("..", "..", "adapters", "grype-v0.116.1.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if manifest.ID != "prc.adapter.grype@0.116" || manifest.Image != GrypeImage {
		t.Fatalf("unexpected Grype manifest: %+v", manifest)
	}
	digest, err := ManifestDigest(manifest)
	if err != nil {
		t.Fatal(err)
	}
	if digest != "b5d83d172c115d2b0b3d855287582a46d2d4c40555100d196ba39631b4a72599" {
		t.Fatalf("Grype manifest digest = %s", digest)
	}
}

func TestGrypeManifestPinsReviewedExecutionContract(t *testing.T) {
	path := filepath.Join("..", "..", "adapters", "grype-v0.116.1.yaml")
	tests := map[string]func(*Manifest){
		"image":    func(value *Manifest) { value.Image = "ghcr.io/anchore/grype@sha256:" + strings.Repeat("f", 64) },
		"command":  func(value *Manifest) { value.Command = append(value.Command, "--only-fixed") },
		"tool":     func(value *Manifest) { value.Tool.Version = "0.116.0" },
		"kind":     func(value *Manifest) { value.ObservationKinds = []string{"package"} },
		"database": func(value *Manifest) { value.DataMounts[0].Destination = "/tmp/grype-db" },
		"network":  func(value *Manifest) { value.Capabilities.Network = "allow" },
		"tasks":    func(value *Manifest) { value.Capabilities.ChildProcesses = false },
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			manifest, err := LoadManifest(path)
			if err != nil {
				t.Fatal(err)
			}
			mutate(&manifest)
			if err := manifest.Validate(); err == nil {
				t.Fatal("tampered Grype manifest was accepted")
			}
		})
	}
}

func TestGrypeParserNormalizesEnrichmentAndDeduplicates(t *testing.T) {
	now := time.Date(2026, 8, 23, 14, 0, 0, 0, time.UTC)
	document := validGrypeDocument(now)
	document.Matches = append(document.Matches, document.Matches[0])
	transcript, payloads, err := parseGrypeOutput(marshalGrypeDocument(t, document), 100, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(transcript.Observations) != 1 || transcript.Observations[0].Outcome != "found" ||
		transcript.Observations[0].Data["severity"] != "critical" ||
		transcript.Observations[0].Data["known_exploited"] != true ||
		transcript.Observations[0].Data["known_ransomware_use"] != true {
		t.Fatalf("unexpected normalized observation: %+v", transcript.Observations)
	}
	artifact := transcript.Artifacts[0]
	payload := payloads[artifact.Digest]
	if int64(len(payload)) != artifact.Size || !bytes.Contains(payload, []byte(`"id":"CVE-2021-23337"`)) ||
		bytes.Contains(payload, []byte(now.Format(time.RFC3339Nano))) {
		t.Fatalf("unexpected canonical artifact: %s", payload)
	}
	var report grypeNormalizedReport
	if err := json.Unmarshal(payload, &report); err != nil || len(report.Findings) != 1 {
		t.Fatalf("canonical report=%+v err=%v", report, err)
	}
}

func TestGrypeParserProducesTimeBoundNoFindingEvidence(t *testing.T) {
	now := time.Date(2026, 8, 23, 14, 0, 0, 0, time.UTC)
	document := validGrypeDocument(now)
	document.Matches = []grypeMatch{}
	transcript, payloads, err := parseGrypeOutput(marshalGrypeDocument(t, document), 100, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(transcript.Observations) != 1 || transcript.Observations[0].Outcome != "not_found" ||
		len(payloads) != 1 {
		t.Fatalf("unexpected no-finding result: %+v", transcript)
	}
}

func TestGrypeParserFailsClosedOnDatabasePolicyAndUnsafeOutput(t *testing.T) {
	now := time.Date(2026, 8, 23, 14, 0, 0, 0, time.UTC)
	tests := map[string]func(*grypeDocument){
		"stale database": func(value *grypeDocument) {
			value.Descriptor.DB.Status.Built = now.Add(-121 * time.Hour).Format(time.RFC3339)
		},
		"invalid database": func(value *grypeDocument) { value.Descriptor.DB.Status.Valid = false },
		"missing KEV provenance": func(value *grypeDocument) {
			delete(value.Descriptor.DB.Providers, "kev")
		},
		"updates enabled": func(value *grypeDocument) {
			value.Descriptor.Configuration["db"].(map[string]any)["auto-update"] = true
		},
		"suppressed match": func(value *grypeDocument) {
			value.IgnoredMatches = []json.RawMessage{json.RawMessage(`{}`)}
		},
		"escaping location": func(value *grypeDocument) {
			value.Matches[0].Artifact.Locations[0].Path = "/../../etc/passwd"
		},
		"oversized nested metadata": func(value *grypeDocument) {
			value.Matches[0].Vulnerability.URLs = make([]string, 257)
		},
		"invalid provider provenance": func(value *grypeDocument) {
			value.Descriptor.DB.Providers["nvd"] = grypeDatabaseProvider{
				Captured: now.Add(time.Hour).Format(time.RFC3339), Input: "mutable",
			}
		},
		"too many matches": func(value *grypeDocument) {
			value.Matches = append(value.Matches, value.Matches[0])
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			document := validGrypeDocument(now)
			mutate(&document)
			limit := 100
			if name == "too many matches" {
				limit = 1
			}
			if _, _, err := parseGrypeOutput(marshalGrypeDocument(t, document), limit, now); err == nil {
				t.Fatal("unsafe Grype output was accepted")
			}
		})
	}
	duplicate := []byte(`{"matches":[],"matches":[]}`)
	if _, _, err := parseGrypeOutput(duplicate, 100, now); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("duplicate keys were not rejected: %v", err)
	}
}
