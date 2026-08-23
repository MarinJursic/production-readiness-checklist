package adapter

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	GrypeProtocolVersion     = "prc-adapter-grype-json-v1"
	GrypeOutputSchemaVersion = "grype.report/v0.116.1+prc-v1"
	GrypeToolVersion         = "0.116.1"
	GrypeObservationKind     = "dependency-vulnerability"
	GrypeImage               = "ghcr.io/anchore/grype@sha256:1e71065c0a4cff3e6bd3b8add525ffac4343eb4971694eb90a31cf6d4d3e85db"
	GrypeConfigSnapshotPath  = ".prc/grype-config.yaml"
	GrypeDataMountName       = "grype-db"
	grypeDataMountPath       = "/prc-inputs/grype-db"
	grypeArtifactID          = "grype-vulnerability-report"
	grypeArtifactMediaType   = "application/vnd.prc.grype.vulnerability-report+json;version=1"
	grypeMaxDatabaseAge      = 120 * time.Hour
)

var (
	grypeConfig = []byte(`check-for-app-update: false
output:
  - json
file: ""
timestamp: true
external-sources:
  enable: false
db:
  cache-dir: /prc-inputs/grype-db
  auto-update: false
  validate-by-hash-on-start: true
  validate-age: true
  max-allowed-built-age: 120h
  require-update-check: false
`)
	grypeVulnerabilityIDPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._:-]{1,127}$`)
	grypeCVEPattern             = regexp.MustCompile(`^CVE-[0-9]{4}-[0-9]{4,}$`)
	grypeDBSchemaPattern        = regexp.MustCompile(`^v6\.[0-9]+\.[0-9]+$`)
	grypeXXH64Pattern           = regexp.MustCompile(`^xxh64:[0-9a-f]{16}$`)
	grypeProviderPattern        = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
	grypeSeverities             = map[string]int{"Unknown": 0, "Negligible": 1, "Low": 2, "Medium": 3, "High": 4, "Critical": 5}
)

func grypeCommand() []string {
	return []string{
		"dir:/workspace", "--quiet", "--config=/workspace/" + GrypeConfigSnapshotPath, "--output=json",
	}
}

func validateGrypeManifest(manifest Manifest) error {
	if manifest.Protocol != GrypeProtocolVersion || manifest.OutputSchema != GrypeOutputSchemaVersion {
		return fmt.Errorf("grype adapter requires its exact protocol and output schema")
	}
	if manifest.Image != GrypeImage || !slices.Equal(manifest.Command, grypeCommand()) {
		return fmt.Errorf("grype adapter requires the reviewed immutable image and scanner-owned command")
	}
	if manifest.Tool.Name != "grype" || manifest.Tool.Version != GrypeToolVersion ||
		manifest.Tool.Upstream != "https://github.com/anchore/grype" ||
		len(manifest.Tool.Formats) != 1 || manifest.Tool.Formats[0].Name != "grype-json" ||
		!slices.Equal(manifest.Tool.Formats[0].Versions, []string{GrypeToolVersion}) {
		return fmt.Errorf("grype adapter tool identity does not match the reviewed normalizer")
	}
	if !slices.Equal(manifest.ObservationKinds, []string{GrypeObservationKind}) {
		return fmt.Errorf("grype adapter must declare only %s observations", GrypeObservationKind)
	}
	if !manifest.Capabilities.WriteScratch || !manifest.Capabilities.ChildProcesses {
		return fmt.Errorf("grype adapter requires a bounded scratch and OS-task allowance")
	}
	if !slices.Equal(manifest.DataMounts, []DataMount{{
		Name: GrypeDataMountName, Destination: grypeDataMountPath, MaxFiles: 16, MaxBytes: 3 * 1024 * 1024 * 1024,
	}}) {
		return fmt.Errorf("grype adapter requires its exact bounded vulnerability database mount")
	}
	if manifest.Resources.PIDs < 32 || manifest.Resources.PIDs > 256 ||
		manifest.Resources.TimeoutSeconds < 60 || manifest.Resources.MemoryMB < 1024 ||
		manifest.Resources.TmpfsMB < 128 || manifest.Resources.MaxStdout < 16*1024*1024 {
		return fmt.Errorf("grype adapter resource limits cannot support or cannot safely bound the reviewed command")
	}
	return nil
}

type grypeDocument struct {
	Matches         []grypeMatch      `json:"matches"`
	IgnoredMatches  []json.RawMessage `json:"ignoredMatches,omitempty"`
	AlertsByPackage []json.RawMessage `json:"alertsByPackage,omitempty"`
	Source          *grypeSource      `json:"source"`
	Distro          grypeDistro       `json:"distro"`
	Descriptor      grypeDescriptor   `json:"descriptor"`
}

type grypeSource struct {
	Type   string          `json:"type"`
	Target json.RawMessage `json:"target"`
}

type grypeDistro struct {
	Name    string   `json:"name"`
	Version string   `json:"version"`
	IDLike  []string `json:"idLike"`
}

type grypeDescriptor struct {
	Name          string                  `json:"name"`
	Version       string                  `json:"version"`
	Configuration map[string]any          `json:"configuration"`
	DB            grypeDatabaseDescriptor `json:"db"`
	Timestamp     string                  `json:"timestamp"`
}

type grypeDatabaseDescriptor struct {
	Status    grypeDatabaseStatus              `json:"status"`
	Providers map[string]grypeDatabaseProvider `json:"providers"`
}

type grypeDatabaseStatus struct {
	SchemaVersion string `json:"schemaVersion"`
	From          string `json:"from"`
	Built         string `json:"built"`
	Path          string `json:"path"`
	Valid         bool   `json:"valid"`
}

type grypeDatabaseProvider struct {
	Captured string `json:"captured"`
	Input    string `json:"input"`
}

type grypeMatch struct {
	Vulnerability          grypeVulnerability `json:"vulnerability"`
	RelatedVulnerabilities []grypeMetadata    `json:"relatedVulnerabilities"`
	MatchDetails           []grypeMatchDetail `json:"matchDetails"`
	Artifact               grypePackage       `json:"artifact"`
}

type grypeMetadata struct {
	ID             string                `json:"id"`
	DataSource     string                `json:"dataSource"`
	Namespace      string                `json:"namespace,omitempty"`
	Severity       string                `json:"severity,omitempty"`
	URLs           []string              `json:"urls"`
	Description    string                `json:"description,omitempty"`
	CVSS           []grypeCVSS           `json:"cvss"`
	KnownExploited []grypeKnownExploited `json:"knownExploited,omitempty"`
	EPSS           []grypeEPSS           `json:"epss,omitempty"`
	CWEs           []grypeCWE            `json:"cwes,omitempty"`
}

type grypeVulnerability struct {
	grypeMetadata
	Fix        grypeFix        `json:"fix"`
	Advisories []grypeAdvisory `json:"advisories"`
	Risk       float64         `json:"risk"`
}

type grypeCVSS struct {
	Source         string           `json:"source,omitempty"`
	Type           string           `json:"type"`
	Version        string           `json:"version"`
	Vector         string           `json:"vector"`
	Metrics        grypeCVSSMetrics `json:"metrics"`
	VendorMetadata json.RawMessage  `json:"vendorMetadata"`
}

type grypeCVSSMetrics struct {
	BaseScore           float64 `json:"baseScore"`
	ExploitabilityScore float64 `json:"exploitabilityScore,omitempty"`
	ImpactScore         float64 `json:"impactScore,omitempty"`
}

type grypeKnownExploited struct {
	CVE                        string   `json:"cve"`
	VendorProject              string   `json:"vendorProject,omitempty"`
	Product                    string   `json:"product,omitempty"`
	DateAdded                  string   `json:"dateAdded,omitempty"`
	RequiredAction             string   `json:"requiredAction,omitempty"`
	DueDate                    string   `json:"dueDate,omitempty"`
	KnownRansomwareCampaignUse string   `json:"knownRansomwareCampaignUse"`
	Notes                      string   `json:"notes,omitempty"`
	URLs                       []string `json:"urls,omitempty"`
	CWEs                       []string `json:"cwes,omitempty"`
}

type grypeEPSS struct {
	CVE        string  `json:"cve"`
	EPSS       float64 `json:"epss"`
	Percentile float64 `json:"percentile"`
	Date       string  `json:"date"`
}

type grypeCWE struct {
	CVE    string `json:"cve"`
	CWE    string `json:"cwe,omitempty"`
	Source string `json:"source,omitempty"`
	Type   string `json:"type,omitempty"`
}

type grypeFix struct {
	Versions  []string          `json:"versions"`
	State     string            `json:"state"`
	Available []json.RawMessage `json:"available,omitempty"`
}

type grypeAdvisory struct {
	ID   string `json:"id"`
	Link string `json:"link"`
}

type grypeMatchDetail struct {
	Type       string          `json:"type"`
	Matcher    string          `json:"matcher"`
	SearchedBy json.RawMessage `json:"searchedBy"`
	Found      json.RawMessage `json:"found"`
	Fix        *struct {
		SuggestedVersion string `json:"suggestedVersion"`
	} `json:"fix,omitempty"`
}

type grypePackage struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Version      string                 `json:"version"`
	Type         string                 `json:"type"`
	Locations    []grypePackageLocation `json:"locations"`
	Language     string                 `json:"language"`
	Licenses     []string               `json:"licenses"`
	CPEs         []string               `json:"cpes"`
	PURL         string                 `json:"purl"`
	Upstreams    []grypeUpstream        `json:"upstreams"`
	MetadataType string                 `json:"metadataType,omitempty"`
	Metadata     json.RawMessage        `json:"metadata,omitempty"`
	Annotations  map[string][]string    `json:"annotations,omitempty"`
}

type grypePackageLocation struct {
	Path        string            `json:"path"`
	AccessPath  string            `json:"accessPath"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type grypeUpstream struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type grypeNormalizedReport struct {
	Schema   string                   `json:"schema"`
	Tool     grypeNormalizedTool      `json:"tool"`
	Database grypeNormalizedDatabase  `json:"database"`
	Findings []grypeNormalizedFinding `json:"findings"`
}

type grypeNormalizedTool struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type grypeNormalizedDatabase struct {
	SchemaVersion string `json:"schema_version"`
	Built         string `json:"built"`
	ArchiveSHA256 string `json:"archive_sha256"`
	EPSCaptured   string `json:"epss_captured"`
	KEVCaptured   string `json:"kev_captured"`
	NVDCaptured   string `json:"nvd_captured"`
}

type grypeNormalizedFinding struct {
	ID                 string   `json:"id"`
	Aliases            []string `json:"aliases"`
	PackageName        string   `json:"package_name"`
	PackageVersion     string   `json:"package_version"`
	PackageType        string   `json:"package_type"`
	PackageURL         string   `json:"package_url"`
	Locations          []string `json:"locations"`
	Severity           string   `json:"severity"`
	CVSSScore          float64  `json:"cvss_score"`
	CVSSVector         string   `json:"cvss_vector,omitempty"`
	EPSS               float64  `json:"epss"`
	EPSSPercentile     float64  `json:"epss_percentile"`
	EPSSDate           string   `json:"epss_date,omitempty"`
	KnownExploited     bool     `json:"known_exploited"`
	KnownRansomwareUse bool     `json:"known_ransomware_use"`
	Risk               float64  `json:"risk"`
	FixState           string   `json:"fix_state"`
	FixedVersions      []string `json:"fixed_versions"`
	Description        string   `json:"description,omitempty"`
}

func parseGrypeOutput(data []byte, maxFindings int, now time.Time) (Transcript, map[string][]byte, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return Transcript{}, nil, fmt.Errorf("grype output is empty")
	}
	if !utf8.Valid(data) {
		return Transcript{}, nil, fmt.Errorf("grype output is not valid UTF-8")
	}
	if err := rejectDuplicateKeys(data); err != nil {
		return Transcript{}, nil, fmt.Errorf("grype output: %w", err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	var document grypeDocument
	if err := decoder.Decode(&document); err != nil {
		return Transcript{}, nil, fmt.Errorf("decode grype report: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return Transcript{}, nil, fmt.Errorf("grype output contains more than one JSON value")
		}
		return Transcript{}, nil, fmt.Errorf("decode trailing grype output: %w", err)
	}
	if document.Matches == nil || len(document.Matches) > maxFindings {
		return Transcript{}, nil, fmt.Errorf("grype report requires at most %d matches", maxFindings)
	}
	if len(document.IgnoredMatches) != 0 || len(document.AlertsByPackage) != 0 {
		return Transcript{}, nil, fmt.Errorf("grype report contains suppressed matches or unsupported package alerts")
	}
	database, reportTime, err := validateGrypeEnvelope(document, now)
	if err != nil {
		return Transcript{}, nil, err
	}
	_ = reportTime

	byKey := map[string]grypeNormalizedFinding{}
	for index, match := range document.Matches {
		finding, err := normalizeGrypeMatch(match)
		if err != nil {
			return Transcript{}, nil, fmt.Errorf("grype match %d: %w", index, err)
		}
		keyBytes, _ := json.Marshal(struct {
			ID, PURL  string
			Locations []string
		}{finding.ID, finding.PackageURL, finding.Locations})
		keyDigest := sha256.Sum256(keyBytes)
		key := hex.EncodeToString(keyDigest[:])
		if existing, ok := byKey[key]; ok {
			byKey[key] = mergeGrypeFinding(existing, finding)
		} else {
			byKey[key] = finding
		}
	}
	findings := make([]grypeNormalizedFinding, 0, len(byKey))
	for _, finding := range byKey {
		findings = append(findings, finding)
	}
	sort.Slice(findings, func(left, right int) bool {
		if findings[left].ID != findings[right].ID {
			return findings[left].ID < findings[right].ID
		}
		if findings[left].PackageURL != findings[right].PackageURL {
			return findings[left].PackageURL < findings[right].PackageURL
		}
		return strings.Join(findings[left].Locations, "\x00") < strings.Join(findings[right].Locations, "\x00")
	})
	report := grypeNormalizedReport{
		Schema:   "prc.grype-vulnerability-report/v1",
		Tool:     grypeNormalizedTool{Name: "grype", Version: GrypeToolVersion},
		Database: database, Findings: findings,
	}
	canonical, err := json.Marshal(report)
	if err != nil {
		return Transcript{}, nil, fmt.Errorf("encode normalized grype report: %w", err)
	}
	canonical = append(canonical, '\n')
	digest := sha256.Sum256(canonical)
	descriptor := "sha256:" + hex.EncodeToString(digest[:])
	observations := make([]Observation, 0, max(1, len(findings)))
	for _, finding := range findings {
		findingBytes, _ := json.Marshal(finding)
		findingDigest := sha256.Sum256(findingBytes)
		locations := make([]Location, 0, len(finding.Locations))
		for _, path := range finding.Locations {
			locations = append(locations, Location{Path: path})
		}
		observations = append(observations, Observation{
			ID: hex.EncodeToString(findingDigest[:]), Kind: GrypeObservationKind, Outcome: "found",
			Summary:   fmt.Sprintf("%s affects %s %s (%s severity).", finding.ID, finding.PackageName, finding.PackageVersion, strings.ToLower(finding.Severity)),
			Locations: locations,
			Data: map[string]any{
				"aliases": finding.Aliases, "artifact_digest": descriptor,
				"cvss_score": finding.CVSSScore, "cvss_vector": finding.CVSSVector,
				"database_built": database.Built, "database_schema": database.SchemaVersion,
				"epss": finding.EPSS, "epss_percentile": finding.EPSSPercentile,
				"fix_state": finding.FixState, "fixed_versions": finding.FixedVersions,
				"known_exploited": finding.KnownExploited, "known_ransomware_use": finding.KnownRansomwareUse,
				"package_name": finding.PackageName, "package_type": finding.PackageType,
				"package_url": finding.PackageURL, "package_version": finding.PackageVersion,
				"risk": finding.Risk, "severity": strings.ToLower(finding.Severity),
			},
		})
	}
	if len(observations) == 0 {
		id := sha256.Sum256([]byte(database.SchemaVersion + "\x00" + database.Built + "\x00" + database.ArchiveSHA256))
		observations = append(observations, Observation{
			ID: hex.EncodeToString(id[:]), Kind: GrypeObservationKind, Outcome: "not_found",
			Summary:   "Grype reported no known dependency vulnerabilities for the sealed repository inventory.",
			Locations: []Location{}, Data: map[string]any{
				"artifact_digest": descriptor, "database_built": database.Built,
				"database_schema": database.SchemaVersion, "finding_count": 0,
			},
		})
	}
	transcript := Transcript{
		Logs: []Log{}, Observations: observations,
		Artifacts: []Artifact{{
			ID: grypeArtifactID, MediaType: grypeArtifactMediaType,
			Digest: descriptor, Size: int64(len(canonical)), Path: "vulnerabilities.json",
		}},
		Summary: Summary{Type: "summary", Status: "completed", Counts: map[string]int{
			"logs": 0, "observations": len(observations), "artifacts": 1,
		}},
	}
	return transcript, map[string][]byte{descriptor: canonical}, nil
}

func validateGrypeEnvelope(document grypeDocument, now time.Time) (grypeNormalizedDatabase, time.Time, error) {
	if document.Source == nil || document.Source.Type != "directory" {
		return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype report source must be the sealed directory")
	}
	var target string
	if err := json.Unmarshal(document.Source.Target, &target); err != nil || target != "/workspace" {
		return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype report source target is invalid")
	}
	if document.Descriptor.Name != "grype" || document.Descriptor.Version != GrypeToolVersion {
		return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype report tool identity is invalid")
	}
	if !validOptionalBoundedText(document.Distro.Name, 1024) ||
		!validOptionalBoundedText(document.Distro.Version, 1024) || len(document.Distro.IDLike) > 64 {
		return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype report distribution metadata is invalid")
	}
	for _, value := range document.Distro.IDLike {
		if !validBoundedText(value, 1024) {
			return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype report distribution metadata is invalid")
		}
	}
	reportTime, err := time.Parse(time.RFC3339Nano, document.Descriptor.Timestamp)
	if err != nil || reportTime.Location() != time.UTC || reportTime.After(now.UTC().Add(5*time.Minute)) || now.UTC().Sub(reportTime) > time.Hour {
		return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype report timestamp is invalid")
	}
	if err := validateGrypeConfiguration(document.Descriptor.Configuration); err != nil {
		return grypeNormalizedDatabase{}, time.Time{}, err
	}
	status := document.Descriptor.DB.Status
	if !status.Valid || !grypeDBSchemaPattern.MatchString(status.SchemaVersion) || status.Path != grypeDataMountPath+"/6/vulnerability.db" {
		return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype vulnerability database status is invalid")
	}
	built, err := time.Parse(time.RFC3339Nano, status.Built)
	if err != nil || built.Location() != time.UTC || built.After(reportTime.Add(5*time.Minute)) || reportTime.Sub(built) > grypeMaxDatabaseAge {
		return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype vulnerability database is stale or has an invalid build time")
	}
	archiveDigest, err := grypeDatabaseArchiveDigest(status.From)
	if err != nil {
		return grypeNormalizedDatabase{}, time.Time{}, err
	}
	if len(document.Descriptor.DB.Providers) == 0 || len(document.Descriptor.DB.Providers) > 128 {
		return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype database provider provenance is invalid")
	}
	for name, provider := range document.Descriptor.DB.Providers {
		captured, parseErr := time.Parse(time.RFC3339Nano, provider.Captured)
		if !grypeProviderPattern.MatchString(name) || parseErr != nil || captured.Location() != time.UTC ||
			captured.After(reportTime.Add(5*time.Minute)) || !grypeXXH64Pattern.MatchString(provider.Input) {
			return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype database provider provenance is invalid")
		}
	}
	required := []string{"epss", "kev", "nvd"}
	captures := map[string]string{}
	for _, name := range required {
		provider, ok := document.Descriptor.DB.Providers[name]
		if !ok {
			return grypeNormalizedDatabase{}, time.Time{}, fmt.Errorf("grype database lacks valid %s provider provenance", name)
		}
		captures[name] = provider.Captured
	}
	return grypeNormalizedDatabase{
		SchemaVersion: status.SchemaVersion, Built: status.Built, ArchiveSHA256: archiveDigest,
		EPSCaptured: captures["epss"], KEVCaptured: captures["kev"], NVDCaptured: captures["nvd"],
	}, reportTime, nil
}

func validateGrypeConfiguration(configuration map[string]any) error {
	if configuration == nil || configuration["check-for-app-update"] != false || stringValue(configuration["file"]) != "" {
		return fmt.Errorf("grype report does not reflect the scanner-owned configuration")
	}
	output, ok := configuration["output"].([]any)
	if !ok || len(output) != 1 || stringValue(output[0]) != "json" {
		return fmt.Errorf("grype report output configuration is invalid")
	}
	external, ok := configuration["externalSources"].(map[string]any)
	if !ok || external["enable"] != false {
		return fmt.Errorf("grype external source access was not disabled")
	}
	database, ok := configuration["db"].(map[string]any)
	if !ok || stringValue(database["cache-dir"]) != grypeDataMountPath || database["auto-update"] != false ||
		database["validate-by-hash-on-start"] != true || database["validate-age"] != true ||
		database["require-update-check"] != false || numberInt64(database["max-allowed-built-age"]) != int64(grypeMaxDatabaseAge) {
		return fmt.Errorf("grype database safety configuration is invalid")
	}
	return nil
}

func grypeDatabaseArchiveDigest(source string) (string, error) {
	parsed, err := url.Parse(source)
	if err != nil || parsed.Scheme != "https" || parsed.Host != "grype.anchore.io" ||
		!strings.HasPrefix(parsed.Path, "/databases/v6/vulnerability-db_v6.") || !strings.HasSuffix(parsed.Path, ".tar.zst") {
		return "", fmt.Errorf("grype database source is invalid")
	}
	checksum := parsed.Query().Get("checksum")
	if !strings.HasPrefix(checksum, "sha256:") || !hexDigestPattern.MatchString(strings.TrimPrefix(checksum, "sha256:")) {
		return "", fmt.Errorf("grype database source lacks an immutable SHA-256 checksum")
	}
	return strings.TrimPrefix(checksum, "sha256:"), nil
}

func normalizeGrypeMatch(match grypeMatch) (grypeNormalizedFinding, error) {
	if len(match.RelatedVulnerabilities) > 128 {
		return grypeNormalizedFinding{}, fmt.Errorf("too many related vulnerabilities")
	}
	metadata := append([]grypeMetadata{match.Vulnerability.grypeMetadata}, match.RelatedVulnerabilities...)
	aliases := make([]string, 0, len(metadata))
	canonical := ""
	severity := "Unknown"
	description := match.Vulnerability.Description
	var cvssScore, epssScore, epssPercentile float64
	var cvssVector, epssDate string
	knownExploited, ransomware := false, false
	for _, item := range metadata {
		if !grypeVulnerabilityIDPattern.MatchString(item.ID) || !validOptionalBoundedText(item.Description, 64*1024) ||
			!validBoundedText(item.DataSource, 16*1024) || !validOptionalBoundedText(item.Namespace, 4096) ||
			item.URLs == nil || len(item.URLs) > 256 || item.CVSS == nil || len(item.CVSS) > 32 ||
			len(item.EPSS) > 32 || len(item.KnownExploited) > 32 || len(item.CWEs) > 256 {
			return grypeNormalizedFinding{}, fmt.Errorf("vulnerability metadata is incomplete or invalid")
		}
		for _, itemURL := range item.URLs {
			if !validBoundedText(itemURL, 16*1024) {
				return grypeNormalizedFinding{}, fmt.Errorf("vulnerability URL metadata is invalid")
			}
		}
		for _, cwe := range item.CWEs {
			if !grypeCVEPattern.MatchString(cwe.CVE) || !validOptionalBoundedText(cwe.CWE, 256) ||
				!validOptionalBoundedText(cwe.Source, 1024) || !validOptionalBoundedText(cwe.Type, 256) {
				return grypeNormalizedFinding{}, fmt.Errorf("vulnerability CWE data is invalid")
			}
		}
		aliases = append(aliases, item.ID)
		if grypeCVEPattern.MatchString(item.ID) && (canonical == "" || item.ID < canonical) {
			canonical = item.ID
		}
		if rank, ok := grypeSeverities[item.Severity]; !ok {
			return grypeNormalizedFinding{}, fmt.Errorf("vulnerability severity %q is unsupported", item.Severity)
		} else if rank > grypeSeverities[severity] {
			severity = item.Severity
		}
		for _, score := range item.CVSS {
			if !validGrypeCVSS(score) {
				return grypeNormalizedFinding{}, fmt.Errorf("vulnerability CVSS data is invalid")
			}
			if score.Metrics.BaseScore > cvssScore || (score.Metrics.BaseScore == cvssScore && score.Vector < cvssVector) {
				cvssScore, cvssVector = score.Metrics.BaseScore, score.Vector
			}
		}
		for _, epss := range item.EPSS {
			if !grypeCVEPattern.MatchString(epss.CVE) || epss.EPSS < 0 || epss.EPSS > 1 ||
				epss.Percentile < 0 || epss.Percentile > 1 || !validDateOnly(epss.Date) {
				return grypeNormalizedFinding{}, fmt.Errorf("vulnerability EPSS data is invalid")
			}
			if epss.EPSS > epssScore || (epss.EPSS == epssScore && epss.Date > epssDate) {
				epssScore, epssPercentile, epssDate = epss.EPSS, epss.Percentile, epss.Date
			}
		}
		for _, kev := range item.KnownExploited {
			if !grypeCVEPattern.MatchString(kev.CVE) || !validOptionalDateOnly(kev.DateAdded) || !validOptionalDateOnly(kev.DueDate) ||
				!validOptionalBoundedText(kev.VendorProject, 4096) || !validOptionalBoundedText(kev.Product, 4096) ||
				!validOptionalBoundedText(kev.RequiredAction, 16*1024) || !validOptionalBoundedText(kev.Notes, 16*1024) ||
				!validOptionalBoundedText(kev.KnownRansomwareCampaignUse, 256) || len(kev.URLs) > 64 || len(kev.CWEs) > 64 {
				return grypeNormalizedFinding{}, fmt.Errorf("vulnerability KEV data is invalid")
			}
			for _, value := range append(append([]string{}, kev.URLs...), kev.CWEs...) {
				if !validBoundedText(value, 16*1024) {
					return grypeNormalizedFinding{}, fmt.Errorf("vulnerability KEV data is invalid")
				}
			}
			knownExploited = true
			if strings.EqualFold(kev.KnownRansomwareCampaignUse, "known") {
				ransomware = true
			}
		}
	}
	if canonical == "" {
		canonical = match.Vulnerability.ID
	}
	aliases = uniqueSortedStrings(aliases)
	if len(match.Artifact.Locations) == 0 || len(match.Artifact.Locations) > 256 ||
		!validBoundedText(match.Artifact.ID, 4096) || !validOptionalBoundedText(match.Artifact.Language, 1024) ||
		!validBoundedText(match.Artifact.Name, 16*1024) || !validBoundedText(match.Artifact.Version, 16*1024) ||
		!observationKindPattern.MatchString(match.Artifact.Type) || !syftPackageURLPattern.MatchString(match.Artifact.PURL) ||
		len(match.Artifact.Licenses) > 256 || len(match.Artifact.CPEs) > 256 || len(match.Artifact.Upstreams) > 256 ||
		len(match.Artifact.Annotations) > 256 || len(match.Artifact.Metadata) > 1024*1024 {
		return grypeNormalizedFinding{}, fmt.Errorf("affected package identity is invalid")
	}
	for _, value := range append(append([]string{}, match.Artifact.Licenses...), match.Artifact.CPEs...) {
		if !validBoundedText(value, 16*1024) {
			return grypeNormalizedFinding{}, fmt.Errorf("affected package metadata is invalid")
		}
	}
	for _, upstream := range match.Artifact.Upstreams {
		if !validBoundedText(upstream.Name, 16*1024) || !validOptionalBoundedText(upstream.Version, 16*1024) {
			return grypeNormalizedFinding{}, fmt.Errorf("affected package upstream is invalid")
		}
	}
	for key, values := range match.Artifact.Annotations {
		if !validBoundedText(key, 4096) || len(values) > 256 {
			return grypeNormalizedFinding{}, fmt.Errorf("affected package annotation is invalid")
		}
		for _, value := range values {
			if !validBoundedText(value, 16*1024) {
				return grypeNormalizedFinding{}, fmt.Errorf("affected package annotation is invalid")
			}
		}
	}
	locations := make([]string, 0, len(match.Artifact.Locations))
	for _, location := range match.Artifact.Locations {
		if !validOptionalBoundedText(location.AccessPath, 16*1024) || len(location.Annotations) > 64 {
			return grypeNormalizedFinding{}, fmt.Errorf("grype package location metadata is invalid")
		}
		for key, value := range location.Annotations {
			if !validBoundedText(key, 4096) || !validBoundedText(value, 16*1024) {
				return grypeNormalizedFinding{}, fmt.Errorf("grype package location metadata is invalid")
			}
		}
		path, err := normalizeGrypePath(location.Path)
		if err != nil {
			return grypeNormalizedFinding{}, err
		}
		locations = append(locations, path)
	}
	locations = uniqueSortedStrings(locations)
	if len(match.MatchDetails) == 0 || len(match.MatchDetails) > 64 {
		return grypeNormalizedFinding{}, fmt.Errorf("vulnerability match lacks bounded matching evidence")
	}
	for _, detail := range match.MatchDetails {
		if !validBoundedText(detail.Type, 256) || !validBoundedText(detail.Matcher, 256) ||
			len(detail.SearchedBy) == 0 || len(detail.SearchedBy) > 1024*1024 || len(detail.Found) == 0 || len(detail.Found) > 1024*1024 ||
			(detail.Fix != nil && !validBoundedText(detail.Fix.SuggestedVersion, 4096)) {
			return grypeNormalizedFinding{}, fmt.Errorf("vulnerability match evidence is invalid")
		}
	}
	if !map[string]bool{"fixed": true, "not-fixed": true, "wont-fix": true, "unknown": true}[match.Vulnerability.Fix.State] ||
		match.Vulnerability.Fix.Versions == nil || len(match.Vulnerability.Fix.Versions) > 256 ||
		len(match.Vulnerability.Fix.Available) > 256 || len(match.Vulnerability.Advisories) > 256 ||
		math.IsNaN(match.Vulnerability.Risk) || math.IsInf(match.Vulnerability.Risk, 0) || match.Vulnerability.Risk < 0 {
		return grypeNormalizedFinding{}, fmt.Errorf("vulnerability fix or risk data is invalid")
	}
	for _, available := range match.Vulnerability.Fix.Available {
		if len(available) == 0 || len(available) > 64*1024 {
			return grypeNormalizedFinding{}, fmt.Errorf("vulnerability fix metadata is invalid")
		}
	}
	for _, advisory := range match.Vulnerability.Advisories {
		if !validBoundedText(advisory.ID, 4096) || !validBoundedText(advisory.Link, 16*1024) {
			return grypeNormalizedFinding{}, fmt.Errorf("vulnerability advisory is invalid")
		}
	}
	fixedVersions := uniqueSortedStrings(match.Vulnerability.Fix.Versions)
	for _, version := range fixedVersions {
		if !validBoundedText(version, 4096) {
			return grypeNormalizedFinding{}, fmt.Errorf("vulnerability fixed version is invalid")
		}
	}
	return grypeNormalizedFinding{
		ID: canonical, Aliases: aliases, PackageName: match.Artifact.Name, PackageVersion: match.Artifact.Version,
		PackageType: match.Artifact.Type, PackageURL: match.Artifact.PURL, Locations: locations,
		Severity: severity, CVSSScore: cvssScore, CVSSVector: cvssVector,
		EPSS: epssScore, EPSSPercentile: epssPercentile, EPSSDate: epssDate,
		KnownExploited: knownExploited, KnownRansomwareUse: ransomware, Risk: match.Vulnerability.Risk,
		FixState: match.Vulnerability.Fix.State, FixedVersions: fixedVersions, Description: description,
	}, nil
}

func validGrypeCVSS(score grypeCVSS) bool {
	if !map[string]bool{"2.0": true, "3.0": true, "3.1": true, "4.0": true}[score.Version] ||
		!validBoundedText(score.Type, 64) || !validBoundedText(score.Vector, 4096) ||
		math.IsNaN(score.Metrics.BaseScore) || math.IsInf(score.Metrics.BaseScore, 0) ||
		score.Metrics.BaseScore < 0 || score.Metrics.BaseScore > 10 {
		return false
	}
	return true
}

func normalizeGrypePath(path string) (string, error) {
	path = strings.TrimPrefix(path, "/workspace/")
	path = strings.TrimPrefix(path, "/")
	if err := validateRelativePath(path); err != nil {
		return "", fmt.Errorf("grype package location is invalid")
	}
	return path, nil
}

func mergeGrypeFinding(left, right grypeNormalizedFinding) grypeNormalizedFinding {
	left.Aliases = uniqueSortedStrings(append(left.Aliases, right.Aliases...))
	if grypeSeverities[right.Severity] > grypeSeverities[left.Severity] {
		left.Severity = right.Severity
	}
	if right.CVSSScore > left.CVSSScore || (right.CVSSScore == left.CVSSScore && right.CVSSVector < left.CVSSVector) {
		left.CVSSScore, left.CVSSVector = right.CVSSScore, right.CVSSVector
	}
	if right.EPSS > left.EPSS || (right.EPSS == left.EPSS && right.EPSSDate > left.EPSSDate) {
		left.EPSS, left.EPSSPercentile, left.EPSSDate = right.EPSS, right.EPSSPercentile, right.EPSSDate
	}
	left.KnownExploited = left.KnownExploited || right.KnownExploited
	left.KnownRansomwareUse = left.KnownRansomwareUse || right.KnownRansomwareUse
	if right.Risk > left.Risk {
		left.Risk = right.Risk
	}
	left.FixedVersions = uniqueSortedStrings(append(left.FixedVersions, right.FixedVersions...))
	return left
}

func uniqueSortedStrings(values []string) []string {
	sort.Strings(values)
	result := make([]string, 0, len(values))
	for _, value := range values {
		if len(result) == 0 || result[len(result)-1] != value {
			result = append(result, value)
		}
	}
	return result
}

func validDateOnly(value string) bool {
	parsed, err := time.Parse(time.DateOnly, value)
	return err == nil && parsed.Format(time.DateOnly) == value
}

func validOptionalDateOnly(value string) bool {
	return value == "" || validDateOnly(value)
}

func validOptionalBoundedText(value string, maximum int) bool {
	return value == "" || validBoundedText(value, maximum)
}

func numberInt64(value any) int64 {
	switch number := value.(type) {
	case json.Number:
		result, _ := number.Int64()
		return result
	case float64:
		if number >= math.MinInt64 && number <= math.MaxInt64 && math.Trunc(number) == number {
			return int64(number)
		}
	}
	return 0
}
