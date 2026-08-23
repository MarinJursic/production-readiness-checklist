// Package state maintains a transactional SQLite index over immutable scanner
// run and evidence records. The JSON records remain the canonical artifacts.
package state

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/adapter"
	"github.com/MarinJursic/production-readiness-checklist/scanner/engine"
	"github.com/MarinJursic/production-readiness-checklist/scanner/finding"
	workspaceinventory "github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	_ "modernc.org/sqlite"
)

const (
	SchemaVersion           = 2
	HistorySchema           = "prc.history/v0.1"
	CheckSchema             = "prc.state-check/v0.2"
	maxCanonicalRecordBytes = 256 * 1024 * 1024
)

type Store struct {
	root string
	path string
	db   *sql.DB
}

type RunSummary struct {
	RunID               string    `json:"run_id"`
	StartedAt           time.Time `json:"started_at"`
	CompletedAt         time.Time `json:"completed_at"`
	TargetName          string    `json:"target_name"`
	ProfileID           string    `json:"profile_id"`
	ProfileVersion      string    `json:"profile_version"`
	InventoryDigest     string    `json:"inventory_digest"`
	ConfigurationDigest string    `json:"configuration_digest,omitempty"`
	ProjectID           string    `json:"project_id,omitempty"`
	TerminalState       string    `json:"terminal_state"`
	PassCount           int       `json:"pass_count"`
	FailCount           int       `json:"fail_count"`
	BlockedCount        int       `json:"blocked_count"`
}

type Query struct {
	Limit         int
	TargetName    string
	ProfileID     string
	TerminalState string
}

type HistoryReport struct {
	SchemaVersion string       `json:"schema_version"`
	GeneratedAt   time.Time    `json:"generated_at"`
	StatePath     string       `json:"state_path"`
	Runs          []RunSummary `json:"runs"`
}

type CheckReport struct {
	SchemaVersion string    `json:"schema_version"`
	CheckedAt     time.Time `json:"checked_at"`
	StatePath     string    `json:"state_path"`
	Integrity     string    `json:"integrity"`
	Counts        Counts    `json:"counts"`
}

type Counts struct {
	Runs           int `json:"runs"`
	Results        int `json:"results"`
	Findings       int `json:"findings"`
	Evidence       int `json:"evidence"`
	InventoryFiles int `json:"inventory_files"`
	InventoryFacts int `json:"inventory_facts"`
	AuditEvents    int `json:"audit_events"`
}

func Open(ctx context.Context, root string) (*Store, error) {
	root, err := prepareRoot(root)
	if err != nil {
		return nil, err
	}
	path := filepath.Join(root, "state.sqlite")
	dsnURL := &url.URL{Scheme: "file", Path: filepath.ToSlash(path)}
	query := dsnURL.Query()
	query.Set("_busy_timeout", "5000")
	query.Set("_foreign_keys", "on")
	query.Set("_journal_mode", "DELETE")
	query.Set("_synchronous", "FULL")
	query.Set("_txlock", "immediate")
	query.Set("_dqs", "false")
	query.Set("_error_rc", "true")
	dsnURL.RawQuery = query.Encode()
	db, err := sql.Open("sqlite", dsnURL.String())
	if err != nil {
		return nil, fmt.Errorf("open scanner state database: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	store := &Store{root: root, path: path, db: db}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect scanner state database: %w", err)
	}
	if runtime.GOOS != "windows" {
		if err := os.Chmod(path, 0o600); err != nil {
			db.Close()
			return nil, fmt.Errorf("protect scanner state database: %w", err)
		}
	}
	if err := store.migrate(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return store, nil
}

func (store *Store) Close() error { return store.db.Close() }

func (store *Store) Path() string { return store.path }

func (store *Store) migrate(ctx context.Context) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin state migration: %w", err)
	}
	defer tx.Rollback()
	var version int
	if err := tx.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return fmt.Errorf("read state schema version: %w", err)
	}
	if version > SchemaVersion {
		return fmt.Errorf("state database schema %d is newer than supported schema %d", version, SchemaVersion)
	}
	if version == 0 {
		for _, statement := range schemaV1 {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("apply state schema v1: %w", err)
			}
		}
		version = 1
	}
	if version == 1 {
		for _, statement := range schemaV2 {
			if _, err := tx.ExecContext(ctx, statement); err != nil {
				return fmt.Errorf("apply state schema v2: %w", err)
			}
		}
		if _, err := tx.ExecContext(ctx, "PRAGMA user_version = 2"); err != nil {
			return fmt.Errorf("record state schema version: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit state migration: %w", err)
	}
	return nil
}

func (store *Store) IndexRun(ctx context.Context, run model.RunResult) error {
	if err := validateRun(run); err != nil {
		return err
	}
	if err := store.verifyCanonicalRecords(run); err != nil {
		return err
	}
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin run index transaction: %w", err)
	}
	defer tx.Rollback()
	runPath := filepath.ToSlash(filepath.Join("runs", run.RunID+".json"))
	_, err = tx.ExecContext(ctx, `
		INSERT INTO runs (
			run_id, schema_version, started_at, completed_at, target_name,
			profile_id, profile_version, inventory_digest, configuration_digest,
			project_id, terminal_state, record_path
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NULLIF(?, ''), NULLIF(?, ''), ?, ?)
		ON CONFLICT(run_id) DO UPDATE SET
			schema_version=excluded.schema_version, started_at=excluded.started_at,
			completed_at=excluded.completed_at, target_name=excluded.target_name,
			profile_id=excluded.profile_id, profile_version=excluded.profile_version,
			inventory_digest=excluded.inventory_digest,
			configuration_digest=excluded.configuration_digest,
			project_id=excluded.project_id, terminal_state=excluded.terminal_state,
			record_path=excluded.record_path`,
		run.RunID, run.SchemaVersion, timestamp(run.StartedAt), timestamp(run.CompletedAt),
		run.Inventory.TargetName, run.Plan.ProfileID, run.Plan.ProfileVersion, run.Inventory.Digest,
		run.Plan.ConfigurationDigest, run.Plan.ProjectID, run.TerminalState, runPath)
	if err != nil {
		return fmt.Errorf("index run: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM results WHERE run_id = ?", run.RunID); err != nil {
		return fmt.Errorf("replace indexed run results: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM inventory_files WHERE inventory_digest = ?", run.Inventory.Digest); err != nil {
		return fmt.Errorf("replace indexed inventory files: %w", err)
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM inventory_facts WHERE inventory_digest = ?", run.Inventory.Digest); err != nil {
		return fmt.Errorf("replace indexed inventory facts: %w", err)
	}
	for _, record := range run.Inventory.Files {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO inventory_files
			(inventory_digest, path, sha256, size, mode) VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(inventory_digest, path) DO UPDATE SET
				sha256=excluded.sha256, size=excluded.size, mode=excluded.mode`,
			run.Inventory.Digest, record.Path, record.SHA256, record.Size, record.Mode)
		if err != nil {
			return fmt.Errorf("index inventory file %s: %w", record.Path, err)
		}
	}
	for _, fact := range run.Inventory.Facts {
		limitations, marshalErr := json.Marshal(fact.Limitations)
		if marshalErr != nil {
			return fmt.Errorf("encode inventory fact limitations: %w", marshalErr)
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO inventory_facts
			(inventory_digest, key, value, source, detector, detector_version, confidence, scope_path, limitations_json)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(inventory_digest, key, value, source, detector) DO UPDATE SET
				detector_version=excluded.detector_version, confidence=excluded.confidence,
				scope_path=excluded.scope_path, limitations_json=excluded.limitations_json`,
			run.Inventory.Digest, fact.Key, fact.Value, fact.Source, fact.Detector,
			fact.DetectorVersion, fact.Confidence, fact.ScopePath, string(limitations))
		if err != nil {
			return fmt.Errorf("index inventory fact %s: %w", fact.Key, err)
		}
	}
	for _, result := range run.Results {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO results
			(run_id, assertion_id, assessment, execution, severity, gate_state, remediation_class, summary)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			run.RunID, result.AssertionID, result.Assessment, result.Execution,
			result.Severity, result.Gate, result.RemediationClass, result.Summary)
		if err != nil {
			return fmt.Errorf("index result %s: %w", result.AssertionID, err)
		}
		for _, evidence := range result.EvidenceObserved {
			evidencePath := filepath.ToSlash(filepath.Join("evidence", evidence.ID[:2], evidence.ID+".json"))
			_, err = tx.ExecContext(ctx, `
				INSERT INTO evidence
				(evidence_id, schema_version, kind, authority, producer, target_digest, source,
				 content_sha256, size, observed_at, summary, record_path)
				VALUES (?, ?, ?, ?, ?, ?, ?, NULLIF(?, ''), ?, ?, ?, ?)
				ON CONFLICT(evidence_id) DO UPDATE SET
					schema_version=excluded.schema_version, kind=excluded.kind,
					authority=excluded.authority, producer=excluded.producer,
					target_digest=excluded.target_digest, source=excluded.source,
					content_sha256=excluded.content_sha256, size=excluded.size,
					observed_at=excluded.observed_at, summary=excluded.summary,
					record_path=excluded.record_path`,
				evidence.ID, evidence.SchemaVersion, evidence.Kind, evidence.Authority, evidence.Producer,
				evidence.TargetDigest, evidence.Source, evidence.ContentSHA256, evidence.Size,
				timestamp(evidence.ObservedAt), evidence.Summary, evidencePath)
			if err != nil {
				return fmt.Errorf("index evidence %s: %w", evidence.ID, err)
			}
			_, err = tx.ExecContext(ctx, `
				INSERT OR IGNORE INTO run_evidence (run_id, assertion_id, evidence_id) VALUES (?, ?, ?)`,
				run.RunID, result.AssertionID, evidence.ID)
			if err != nil {
				return fmt.Errorf("link evidence %s: %w", evidence.ID, err)
			}
		}
	}
	for _, item := range run.Findings {
		locations, marshalErr := json.Marshal(item.Locations)
		if marshalErr != nil {
			return fmt.Errorf("encode finding locations: %w", marshalErr)
		}
		evidenceIDs, marshalErr := json.Marshal(item.EvidenceIDs)
		if marshalErr != nil {
			return fmt.Errorf("encode finding evidence IDs: %w", marshalErr)
		}
		_, err = tx.ExecContext(ctx, `
			INSERT INTO findings
			(run_id, finding_id, fingerprint, assertion_id, subject_kind, subject_id,
			 inventory_digest, severity, gate_state, remediation_class, title, summary,
			 locations_json, evidence_ids_json)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			run.RunID, item.ID, item.Fingerprint, item.AssertionID, item.Subject.Kind,
			item.Subject.ID, item.Subject.InventoryDigest, item.Severity, item.Gate,
			item.RemediationClass, item.Title, item.Summary, string(locations), string(evidenceIDs))
		if err != nil {
			return fmt.Errorf("index finding %s: %w", item.ID, err)
		}
		for _, location := range item.Locations {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO finding_locations (run_id, finding_id, path, line, column)
				VALUES (?, ?, ?, ?, ?)`, run.RunID, item.ID, location.Path, location.Line, location.Column)
			if err != nil {
				return fmt.Errorf("index finding location %s: %w", location.Path, err)
			}
		}
		for _, evidenceID := range item.EvidenceIDs {
			_, err = tx.ExecContext(ctx, `
				INSERT INTO finding_evidence (run_id, finding_id, evidence_id)
				VALUES (?, ?, ?)`, run.RunID, item.ID, evidenceID)
			if err != nil {
				return fmt.Errorf("link finding evidence %s: %w", evidenceID, err)
			}
		}
	}
	details, err := json.Marshal(map[string]string{"record_path": runPath})
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT OR IGNORE INTO audit_events (occurred_at, event_type, run_id, details_json)
		VALUES (?, 'run.indexed', ?, ?)`, timestamp(time.Now().UTC()), run.RunID, string(details))
	if err != nil {
		return fmt.Errorf("record state audit event: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit run index transaction: %w", err)
	}
	return nil
}

func (store *Store) ListRuns(ctx context.Context, query Query) ([]RunSummary, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}
	if query.Limit < 1 || query.Limit > 1000 {
		return nil, fmt.Errorf("history limit must be between 1 and 1000")
	}
	statement := `
		SELECT r.run_id, r.started_at, r.completed_at, r.target_name, r.profile_id,
		       r.profile_version, r.inventory_digest, COALESCE(r.configuration_digest, ''),
		       COALESCE(r.project_id, ''), r.terminal_state,
		       SUM(CASE WHEN x.assessment = 'pass' THEN 1 ELSE 0 END),
		       SUM(CASE WHEN x.assessment = 'fail' THEN 1 ELSE 0 END),
		       SUM(CASE WHEN x.assessment NOT IN ('pass', 'fail', 'not_applicable') THEN 1 ELSE 0 END)
		FROM runs r
		LEFT JOIN results x ON x.run_id = r.run_id
		WHERE (? = '' OR r.target_name = ?)
		  AND (? = '' OR r.profile_id = ?)
		  AND (? = '' OR r.terminal_state = ?)
		GROUP BY r.run_id
		ORDER BY r.completed_at DESC, r.run_id ASC
		LIMIT ?`
	rows, err := store.db.QueryContext(ctx, statement,
		query.TargetName, query.TargetName, query.ProfileID, query.ProfileID,
		query.TerminalState, query.TerminalState, query.Limit)
	if err != nil {
		return nil, fmt.Errorf("query run history: %w", err)
	}
	defer rows.Close()
	result := []RunSummary{}
	for rows.Next() {
		var item RunSummary
		var startedAt, completedAt string
		if err := rows.Scan(
			&item.RunID, &startedAt, &completedAt, &item.TargetName, &item.ProfileID,
			&item.ProfileVersion, &item.InventoryDigest, &item.ConfigurationDigest,
			&item.ProjectID, &item.TerminalState, &item.PassCount, &item.FailCount, &item.BlockedCount,
		); err != nil {
			return nil, fmt.Errorf("decode run history: %w", err)
		}
		item.StartedAt, err = time.Parse(time.RFC3339Nano, startedAt)
		if err != nil {
			return nil, fmt.Errorf("decode run start time: %w", err)
		}
		item.CompletedAt, err = time.Parse(time.RFC3339Nano, completedAt)
		if err != nil {
			return nil, fmt.Errorf("decode run completion time: %w", err)
		}
		result = append(result, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate run history: %w", err)
	}
	return result, nil
}

func (store *Store) LoadRun(ctx context.Context, runID string) (model.RunResult, error) {
	if !digest(runID) {
		return model.RunResult{}, fmt.Errorf("run ID must be a lowercase SHA-256 digest")
	}
	var relative string
	if err := store.db.QueryRowContext(ctx, "SELECT record_path FROM runs WHERE run_id = ?", runID).Scan(&relative); err != nil {
		if err == sql.ErrNoRows {
			return model.RunResult{}, fmt.Errorf("run %s is not indexed", runID)
		}
		return model.RunResult{}, fmt.Errorf("query run %s: %w", runID, err)
	}
	data, err := readCanonicalRecord(store.root, relative)
	if err != nil {
		return model.RunResult{}, fmt.Errorf("read canonical run record: %w", err)
	}
	var run model.RunResult
	if err := json.Unmarshal(data, &run); err != nil {
		return model.RunResult{}, fmt.Errorf("decode canonical run record: %w", err)
	}
	if run.RunID != runID {
		return model.RunResult{}, fmt.Errorf("canonical run record identity mismatch")
	}
	if err := validateRun(run); err != nil {
		return model.RunResult{}, err
	}
	return run, nil
}

func (store *Store) Counts(ctx context.Context) (Counts, error) {
	var result Counts
	queries := []struct {
		name  string
		value *int
	}{
		{"runs", &result.Runs}, {"results", &result.Results}, {"findings", &result.Findings},
		{"evidence", &result.Evidence},
		{"inventory_files", &result.InventoryFiles}, {"inventory_facts", &result.InventoryFacts},
		{"audit_events", &result.AuditEvents},
	}
	for _, query := range queries {
		if err := store.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+query.name).Scan(query.value); err != nil {
			return Counts{}, fmt.Errorf("count state table %s: %w", query.name, err)
		}
	}
	return result, nil
}

func (store *Store) IntegrityCheck(ctx context.Context) error {
	rows, err := store.db.QueryContext(ctx, "PRAGMA integrity_check")
	if err != nil {
		return fmt.Errorf("run SQLite integrity check: %w", err)
	}
	defer rows.Close()
	seenOK := false
	for rows.Next() {
		var result string
		if err := rows.Scan(&result); err != nil {
			return err
		}
		if result != "ok" {
			return fmt.Errorf("SQLite integrity check failed: %s", result)
		}
		seenOK = true
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if !seenOK {
		return fmt.Errorf("SQLite integrity check returned no result")
	}
	foreignRows, err := store.db.QueryContext(ctx, "PRAGMA foreign_key_check")
	if err != nil {
		return fmt.Errorf("run SQLite foreign-key check: %w", err)
	}
	defer foreignRows.Close()
	if foreignRows.Next() {
		return fmt.Errorf("SQLite foreign-key check reported a violation")
	}
	return foreignRows.Err()
}

func (store *Store) verifyCanonicalRecords(run model.RunResult) error {
	runRelative := filepath.ToSlash(filepath.Join("runs", run.RunID+".json"))
	data, err := readCanonicalRecord(store.root, runRelative)
	if err != nil {
		return fmt.Errorf("read immutable run record before indexing: %w", err)
	}
	var stored model.RunResult
	if err := json.Unmarshal(data, &stored); err != nil {
		return fmt.Errorf("decode immutable run record before indexing: %w", err)
	}
	if stored.RunID != run.RunID {
		return fmt.Errorf("immutable run record identity mismatch")
	}
	if err := validateRun(stored); err != nil {
		return err
	}
	for _, result := range run.Results {
		for _, evidence := range result.EvidenceObserved {
			relative := filepath.ToSlash(filepath.Join("evidence", evidence.ID[:2], evidence.ID+".json"))
			data, err := readCanonicalRecord(store.root, relative)
			if err != nil {
				return fmt.Errorf("read immutable evidence record %s: %w", evidence.ID, err)
			}
			var storedEvidence model.Evidence
			if err := json.Unmarshal(data, &storedEvidence); err != nil {
				return fmt.Errorf("decode immutable evidence record %s: %w", evidence.ID, err)
			}
			if storedEvidence.ID != evidence.ID || evidenceIdentity(storedEvidence) != evidence.ID {
				return fmt.Errorf("immutable evidence record %s has invalid identity", evidence.ID)
			}
		}
	}
	return nil
}

func validateRun(run model.RunResult) error {
	if !digest(run.RunID) || runIdentity(run) != run.RunID {
		return fmt.Errorf("run ID does not match record content")
	}
	if !((run.SchemaVersion == model.RunSchema && run.Plan.SchemaVersion == model.PlanSchema) ||
		(run.SchemaVersion == "prc.run/v0.7" && run.Plan.SchemaVersion == model.PlanSchema) ||
		(run.SchemaVersion == "prc.run/v0.6" && run.Plan.SchemaVersion == "prc.plan/v0.5") ||
		(run.SchemaVersion == "prc.run/v0.5" && run.Plan.SchemaVersion == "prc.plan/v0.5") ||
		(run.SchemaVersion == "prc.run/v0.4" && run.Plan.SchemaVersion == "prc.plan/v0.4") ||
		(run.SchemaVersion == "prc.run/v0.3" && run.Plan.SchemaVersion == "prc.plan/v0.3")) {
		return fmt.Errorf("unsupported or mismatched run and plan schemas %q and %q", run.SchemaVersion, run.Plan.SchemaVersion)
	}
	if run.Inventory.SchemaVersion != model.InventorySchema {
		return fmt.Errorf("unsupported inventory schema %q", run.Inventory.SchemaVersion)
	}
	if !digest(run.Inventory.Digest) || run.Plan.InventoryDigest != run.Inventory.Digest ||
		run.Plan.TargetName != run.Inventory.TargetName || run.Plan.TargetCommit != run.Inventory.GitCommit {
		return fmt.Errorf("run inventory identity is invalid")
	}
	if err := workspaceinventory.VerifyIdentity(run.Inventory); err != nil {
		return err
	}
	if !digest(run.Plan.Digest) || planIdentity(run.Plan) != run.Plan.Digest {
		return fmt.Errorf("plan digest does not match record content")
	}
	if run.Plan.SchemaVersion == model.PlanSchema || run.Plan.SchemaVersion == "prc.plan/v0.5" || run.Plan.SchemaVersion == "prc.plan/v0.4" {
		if run.Plan.EngineVersion == "" || !digest(run.Plan.ProfileDigest) {
			return fmt.Errorf("bound plan lacks engine or profile binding")
		}
		for _, planned := range run.Plan.Assertions {
			if planned.AssertionRevision < 1 || !digest(planned.DefinitionDigest) {
				return fmt.Errorf("current plan assertion %s lacks a definition binding", planned.AssertionID)
			}
		}
	}
	if run.Plan.SchemaVersion == model.PlanSchema && !digest(run.Plan.CatalogDigest) {
		return fmt.Errorf("current plan lacks catalog binding")
	}
	if run.Plan.SchemaVersion == model.PlanSchema {
		if run.Plan.ExecutionMode == "" || run.Plan.Implementations == nil || run.Plan.Adapters == nil || run.Plan.Nodes == nil {
			return fmt.Errorf("current plan lacks an execution mode, implementation registry, adapters, or DAG nodes")
		}
		if err := engine.ValidateExecutionPlan(run.Plan); err != nil {
			return fmt.Errorf("current plan execution contract is invalid: %w", err)
		}
	}
	expectedExecutionSchema := "prc.adapter-execution/v0.1"
	if run.SchemaVersion == model.RunSchema {
		expectedExecutionSchema = model.AdapterExecutionSchema
	}
	seenExecutions := map[string]bool{}
	for _, execution := range run.AdapterExecutions {
		if execution.SchemaVersion != expectedExecutionSchema {
			return fmt.Errorf("run %s requires adapter execution schema %s", run.SchemaVersion, expectedExecutionSchema)
		}
		if err := adapter.ValidateExecution(execution); err != nil {
			return fmt.Errorf("invalid adapter execution %s: %w", execution.ExecutionID, err)
		}
		if seenExecutions[execution.ExecutionID] || execution.Subject.TargetName != run.Inventory.TargetName ||
			execution.Subject.TargetCommit != run.Inventory.GitCommit ||
			execution.Subject.InventoryDigest != run.Inventory.Digest {
			return fmt.Errorf("adapter execution %s is duplicate or bound to a different inventory", execution.ExecutionID)
		}
		seenExecutions[execution.ExecutionID] = true
	}
	planned := map[string]bool{}
	for _, item := range run.Plan.Assertions {
		if item.AssertionID == "" || planned[item.AssertionID] {
			return fmt.Errorf("plan contains an empty or duplicate assertion ID")
		}
		planned[item.AssertionID] = true
	}
	results := map[string]bool{}
	for _, result := range run.Results {
		if !planned[result.AssertionID] || results[result.AssertionID] {
			return fmt.Errorf("run contains an unplanned or duplicate result for %s", result.AssertionID)
		}
		results[result.AssertionID] = true
		for _, evidence := range result.EvidenceObserved {
			if !digest(evidence.ID) || evidenceIdentity(evidence) != evidence.ID {
				return fmt.Errorf("evidence ID does not match record content")
			}
			if evidence.TargetDigest != run.Inventory.Digest {
				return fmt.Errorf("evidence %s is bound to a different inventory", evidence.ID)
			}
		}
	}
	if len(results) != len(planned) {
		return fmt.Errorf("run does not contain exactly one result for every planned assertion")
	}
	if run.SchemaVersion == model.RunSchema || run.SchemaVersion == "prc.run/v0.7" || run.SchemaVersion == "prc.run/v0.6" {
		if run.Findings == nil {
			return fmt.Errorf("current run findings must encode as an array")
		}
		failureResults := map[string]model.AssertionResult{}
		for _, result := range run.Results {
			if result.Assessment == "fail" {
				failureResults[result.AssertionID] = result
			}
		}
		seenFindings := map[string]bool{}
		seenAssertions := map[string]bool{}
		for _, item := range run.Findings {
			if err := finding.Validate(item); err != nil {
				return fmt.Errorf("invalid finding %s: %w", item.ID, err)
			}
			result, ok := failureResults[item.AssertionID]
			if !ok || seenFindings[item.ID] || seenAssertions[item.AssertionID] {
				return fmt.Errorf("finding %s is duplicate or does not map one-to-one to a failed result", item.ID)
			}
			seenFindings[item.ID], seenAssertions[item.AssertionID] = true, true
			if item.Subject.InventoryDigest != run.Inventory.Digest || item.Summary != result.Summary ||
				item.Severity != result.Severity || item.Gate != result.Gate ||
				item.RemediationClass != result.RemediationClass || !sameStringSet(item.ControlIDs, result.ControlIDs) {
				return fmt.Errorf("finding %s does not match its failed result", item.ID)
			}
			evidenceIDs := make([]string, 0, len(result.EvidenceObserved))
			for _, evidence := range result.EvidenceObserved {
				evidenceIDs = append(evidenceIDs, evidence.ID)
			}
			if !sameStringSet(item.EvidenceIDs, evidenceIDs) {
				return fmt.Errorf("finding %s evidence does not match its failed result", item.ID)
			}
		}
		if len(seenAssertions) != len(failureResults) {
			return fmt.Errorf("run does not contain exactly one finding for every failed assertion")
		}
	}
	return nil
}

func runIdentity(run model.RunResult) string {
	run.RunID = ""
	payload, _ := json.Marshal(run)
	value := sha256.Sum256(payload)
	return hex.EncodeToString(value[:])
}

func planIdentity(plan model.Plan) string {
	plan.Digest = ""
	payload, _ := json.Marshal(plan)
	value := sha256.Sum256(payload)
	return hex.EncodeToString(value[:])
}

func evidenceIdentity(evidence model.Evidence) string {
	evidence.ID = ""
	payload, _ := json.Marshal(evidence)
	value := sha256.Sum256(payload)
	return hex.EncodeToString(value[:])
}

func digest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, character := range value {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}

func sameStringSet(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	values := map[string]int{}
	for _, value := range left {
		values[value]++
	}
	for _, value := range right {
		values[value]--
	}
	for _, count := range values {
		if count != 0 {
			return false
		}
	}
	return true
}

func timestamp(value time.Time) string { return value.UTC().Format(time.RFC3339Nano) }

func prepareRoot(root string) (string, error) {
	if strings.TrimSpace(root) == "" {
		return "", fmt.Errorf("state directory is required")
	}
	absolute, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve state directory: %w", err)
	}
	if err := os.MkdirAll(absolute, 0o700); err != nil {
		return "", fmt.Errorf("create state directory: %w", err)
	}
	if resolved, resolveErr := filepath.EvalSymlinks(absolute); resolveErr == nil {
		absolute = resolved
	}
	info, err := os.Stat(absolute)
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("state path is not an accessible directory")
	}
	if runtime.GOOS != "windows" && info.Mode().Perm()&0o077 != 0 {
		return "", fmt.Errorf("state directory must use mode 0700 or stricter")
	}
	return absolute, nil
}

func readCanonicalRecord(root, relative string) ([]byte, error) {
	if relative == "" || filepath.IsAbs(relative) || strings.Contains(relative, "\\") {
		return nil, fmt.Errorf("unsafe indexed record path %q", relative)
	}
	clean := filepath.Clean(filepath.FromSlash(relative))
	if clean == "." || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) || filepath.ToSlash(clean) != relative {
		return nil, fmt.Errorf("unsafe indexed record path %q", relative)
	}
	rootHandle, err := os.OpenRoot(root)
	if err != nil {
		return nil, fmt.Errorf("open state root: %w", err)
	}
	defer rootHandle.Close()
	file, err := rootHandle.Open(clean)
	if err != nil {
		return nil, fmt.Errorf("open indexed record inside state root: %w", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxCanonicalRecordBytes {
		return nil, fmt.Errorf("indexed record must be a regular file no larger than %d bytes", maxCanonicalRecordBytes)
	}
	data, err := io.ReadAll(io.LimitReader(file, maxCanonicalRecordBytes+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxCanonicalRecordBytes {
		return nil, fmt.Errorf("indexed record exceeds %d bytes", maxCanonicalRecordBytes)
	}
	return data, nil
}

var schemaV1 = []string{
	`CREATE TABLE runs (
		run_id TEXT PRIMARY KEY CHECK(length(run_id) = 64),
		schema_version TEXT NOT NULL,
		started_at TEXT NOT NULL,
		completed_at TEXT NOT NULL,
		target_name TEXT NOT NULL,
		profile_id TEXT NOT NULL,
		profile_version TEXT NOT NULL,
		inventory_digest TEXT NOT NULL CHECK(length(inventory_digest) = 64),
		configuration_digest TEXT CHECK(configuration_digest IS NULL OR length(configuration_digest) = 64),
		project_id TEXT,
		terminal_state TEXT NOT NULL,
		record_path TEXT NOT NULL UNIQUE
	) STRICT`,
	`CREATE TABLE results (
		run_id TEXT NOT NULL REFERENCES runs(run_id) ON DELETE CASCADE,
		assertion_id TEXT NOT NULL,
		assessment TEXT NOT NULL,
		execution TEXT NOT NULL,
		severity TEXT NOT NULL,
		gate_state TEXT NOT NULL,
		remediation_class TEXT NOT NULL,
		summary TEXT NOT NULL,
		PRIMARY KEY (run_id, assertion_id)
	) STRICT, WITHOUT ROWID`,
	`CREATE TABLE evidence (
		evidence_id TEXT PRIMARY KEY CHECK(length(evidence_id) = 64),
		schema_version TEXT NOT NULL,
		kind TEXT NOT NULL,
		authority TEXT NOT NULL,
		producer TEXT NOT NULL,
		target_digest TEXT NOT NULL CHECK(length(target_digest) = 64),
		source TEXT NOT NULL,
		content_sha256 TEXT CHECK(content_sha256 IS NULL OR length(content_sha256) = 64),
		size INTEGER NOT NULL CHECK(size >= 0),
		observed_at TEXT NOT NULL,
		summary TEXT NOT NULL,
		record_path TEXT NOT NULL UNIQUE
	) STRICT`,
	`CREATE TABLE run_evidence (
		run_id TEXT NOT NULL,
		assertion_id TEXT NOT NULL,
		evidence_id TEXT NOT NULL REFERENCES evidence(evidence_id),
		PRIMARY KEY (run_id, assertion_id, evidence_id),
		FOREIGN KEY (run_id, assertion_id) REFERENCES results(run_id, assertion_id) ON DELETE CASCADE
	) STRICT, WITHOUT ROWID`,
	`CREATE TABLE inventory_files (
		inventory_digest TEXT NOT NULL CHECK(length(inventory_digest) = 64),
		path TEXT NOT NULL,
		sha256 TEXT NOT NULL CHECK(length(sha256) = 64),
		size INTEGER NOT NULL CHECK(size >= 0),
		mode INTEGER NOT NULL CHECK(mode >= 0),
		PRIMARY KEY (inventory_digest, path)
	) STRICT, WITHOUT ROWID`,
	`CREATE TABLE inventory_facts (
		inventory_digest TEXT NOT NULL CHECK(length(inventory_digest) = 64),
		key TEXT NOT NULL,
		value TEXT NOT NULL,
		source TEXT NOT NULL,
		detector TEXT NOT NULL,
		detector_version TEXT NOT NULL,
		confidence REAL NOT NULL CHECK(confidence >= 0 AND confidence <= 1),
		scope_path TEXT NOT NULL,
		limitations_json TEXT NOT NULL CHECK(json_valid(limitations_json)),
		PRIMARY KEY (inventory_digest, key, value, source, detector)
	) STRICT, WITHOUT ROWID`,
	`CREATE TABLE audit_events (
		event_id INTEGER PRIMARY KEY,
		occurred_at TEXT NOT NULL,
		event_type TEXT NOT NULL,
		run_id TEXT REFERENCES runs(run_id) ON DELETE CASCADE,
		details_json TEXT NOT NULL CHECK(json_valid(details_json)),
		UNIQUE (event_type, run_id)
	) STRICT`,
	`CREATE INDEX runs_completed_idx ON runs(completed_at DESC, run_id)`,
	`CREATE INDEX runs_target_idx ON runs(target_name, completed_at DESC)`,
	`CREATE INDEX results_assessment_idx ON results(assessment, assertion_id)`,
	`CREATE INDEX evidence_target_idx ON evidence(target_digest, kind)`,
}

var schemaV2 = []string{
	`CREATE TABLE findings (
		run_id TEXT NOT NULL,
		finding_id TEXT NOT NULL CHECK(length(finding_id) = 64),
		fingerprint TEXT NOT NULL CHECK(length(fingerprint) = 64),
		assertion_id TEXT NOT NULL,
		subject_kind TEXT NOT NULL,
		subject_id TEXT NOT NULL,
		inventory_digest TEXT NOT NULL CHECK(length(inventory_digest) = 64),
		severity TEXT NOT NULL,
		gate_state TEXT NOT NULL,
		remediation_class TEXT NOT NULL,
		title TEXT NOT NULL,
		summary TEXT NOT NULL,
		locations_json TEXT NOT NULL CHECK(json_valid(locations_json)),
		evidence_ids_json TEXT NOT NULL CHECK(json_valid(evidence_ids_json)),
		PRIMARY KEY (run_id, finding_id),
		UNIQUE (run_id, assertion_id),
		FOREIGN KEY (run_id, assertion_id) REFERENCES results(run_id, assertion_id) ON DELETE CASCADE
	) STRICT, WITHOUT ROWID`,
	`CREATE TABLE finding_locations (
		run_id TEXT NOT NULL,
		finding_id TEXT NOT NULL,
		path TEXT NOT NULL,
		line INTEGER NOT NULL CHECK(line >= 0),
		column INTEGER NOT NULL CHECK(column >= 0),
		PRIMARY KEY (run_id, finding_id, path, line, column),
		FOREIGN KEY (run_id, finding_id) REFERENCES findings(run_id, finding_id) ON DELETE CASCADE
	) STRICT, WITHOUT ROWID`,
	`CREATE TABLE finding_evidence (
		run_id TEXT NOT NULL,
		finding_id TEXT NOT NULL,
		evidence_id TEXT NOT NULL REFERENCES evidence(evidence_id),
		PRIMARY KEY (run_id, finding_id, evidence_id),
		FOREIGN KEY (run_id, finding_id) REFERENCES findings(run_id, finding_id) ON DELETE CASCADE
	) STRICT, WITHOUT ROWID`,
	`CREATE INDEX findings_fingerprint_idx ON findings(fingerprint, run_id)`,
	`CREATE INDEX findings_severity_idx ON findings(severity, gate_state, assertion_id)`,
}
