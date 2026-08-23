package report

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func Write(format string, output io.Writer, run model.RunResult) error {
	switch format {
	case "markdown":
		return writeMarkdown(output, run)
	case "sarif":
		return writeSARIF(output, run)
	case "junit":
		return writeJUnit(output, run)
	case "html":
		return writeHTML(output, run)
	default:
		return fmt.Errorf("unsupported report format %q", format)
	}
}

func assessmentCounts(results []model.AssertionResult) [][2]string {
	counts := map[string]int{}
	for _, result := range results {
		counts[result.Assessment]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	values := make([][2]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, [2]string{key, strconv.Itoa(counts[key])})
	}
	return values
}

func controlDispositionCounts(results []model.ControlResult) [][2]string {
	counts := map[string]int{}
	for _, result := range results {
		counts[result.Disposition]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	values := make([][2]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, [2]string{key, strconv.Itoa(counts[key])})
	}
	return values
}

func markdownCell(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "|", "\\|")
	return strings.Join(strings.Fields(value), " ")
}

func writeMarkdown(output io.Writer, run model.RunResult) error {
	if _, err := fmt.Fprintf(output, "# Production readiness assessment\n\n"+
		"- Run: `%s`\n- Profile: `%s@%s`\n- Target: `%s`\n- Inventory: `%s`\n- Terminal state: **%s**\n\n",
		run.RunID, run.Plan.ProfileID, run.Plan.ProfileVersion, markdownCell(run.Inventory.TargetName),
		run.Inventory.Digest, run.TerminalState); err != nil {
		return err
	}
	if run.Plan.ConfigurationDigest != "" {
		if _, err := fmt.Fprintf(output,
			"- Configuration: `%s`\n- Project: `%s`\n- Environments: %s\n- Artifacts: %s\n\n",
			run.Plan.ConfigurationDigest, run.Plan.ProjectID,
			markdownCell(strings.Join(run.Plan.TargetEnvironments, ", ")),
			markdownCell(strings.Join(run.Plan.ArtifactDigests, ", "))); err != nil {
			return err
		}
	}
	if _, err := fmt.Fprintln(output, "## Result counts\n\n| Assessment | Count |\n| --- | ---: |"); err != nil {
		return err
	}
	for _, item := range assessmentCounts(run.Results) {
		if _, err := fmt.Fprintf(output, "| %s | %s |\n", item[0], item[1]); err != nil {
			return err
		}
	}
	if run.ControlCatalog != nil {
		if _, err := fmt.Fprintf(output,
			"\n## Complete control catalog\n\n- Registry: `%s`\n- Registry digest: `%s`\n- Source digest: `%s`\n- Registered controls: **%d** (%d active)\n- Narrow profile terminal state before complete-catalog expansion: **%s**\n\n",
			run.ControlCatalog.RegistryVersion, run.ControlCatalog.RegistrySHA256,
			run.ControlCatalog.SourceSHA256, run.ControlCatalog.ControlCount,
			run.ControlCatalog.ActiveControlCount, run.ControlCatalog.ProfileTerminalState); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(output, "| Disposition | Count |\n| --- | ---: |"); err != nil {
			return err
		}
		for _, item := range controlDispositionCounts(run.ControlResults) {
			if _, err := fmt.Fprintf(output, "| %s | %s |\n", item[0], item[1]); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(output, "\n### Every registered control\n\n| Disposition | Control | Coverage | Authority | Statement | Source | Assertions |\n| --- | --- | --- | --- | --- | --- | --- |"); err != nil {
			return err
		}
		for _, result := range run.ControlResults {
			if _, err := fmt.Fprintf(output, "| %s | `%s` | %s | %s | %s | `%s:%d` | %s |\n",
				result.Disposition, result.ControlID, result.Coverage, result.Authority,
				markdownCell(result.Statement), result.Source.Path, result.Source.Line,
				markdownCell(strings.Join(result.ExecutedAssertionIDs, ", "))); err != nil {
				return err
			}
		}
	}
	if len(run.AdapterExecutions) > 0 {
		if _, err := fmt.Fprintln(output, "\n## Adapter executions\n\n| Adapter | Manifest | Authorization | Trust | Registry | Status | Execution |\n| --- | --- | --- | --- | --- | --- | --- |"); err != nil {
			return err
		}
		for _, execution := range run.AdapterExecutions {
			registry := execution.Resolution.RegistryID
			if registry == "" {
				registry = "—"
			}
			if _, err := fmt.Fprintf(output, "| `%s` | `%s` | %s | %s | %s | %s | `%s` |\n",
				execution.AdapterID, execution.ManifestSHA256, execution.Resolution.Source,
				execution.Resolution.Trust, registry, execution.Transcript.Summary.Status, execution.ExecutionID); err != nil {
				return err
			}
		}
	}
	if len(run.Findings) > 0 {
		if _, err := fmt.Fprintln(output, "\n## Findings\n\n| Severity | Finding | Assertion | Gate | Summary | Locations | Evidence |\n| --- | --- | --- | --- | --- | ---: | ---: |"); err != nil {
			return err
		}
		for _, finding := range run.Findings {
			if _, err := fmt.Fprintf(output, "| %s | `%s` | `%s` | %s | %s | %d | %d |\n",
				finding.Severity, finding.ID, finding.AssertionID, finding.Gate,
				markdownCell(finding.Summary), len(finding.Locations), len(finding.EvidenceIDs)); err != nil {
				return err
			}
		}
	}
	if _, err := fmt.Fprintln(output, "\n## Assertions\n\n| Assessment | Assertion | Severity | Gate | Summary | Evidence |\n| --- | --- | --- | --- | --- | ---: |"); err != nil {
		return err
	}
	for _, result := range run.Results {
		if _, err := fmt.Fprintf(output, "| %s | `%s` | %s | %s | %s | %d |\n",
			result.Assessment, result.AssertionID, result.Severity, result.Gate,
			markdownCell(result.Summary), len(result.EvidenceObserved)); err != nil {
			return err
		}
	}
	_, err := fmt.Fprintln(output, "\nThis report is scoped to the named profile, target inventory, and evidence set. It is not an unqualified production-readiness or compliance claim.")
	return err
}

type sarifLog struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool       sarifTool         `json:"tool"`
	Results    []sarifResult     `json:"results"`
	Properties map[string]string `json:"properties"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string      `json:"name"`
	InformationURI string      `json:"informationUri"`
	Rules          []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string         `json:"id"`
	ShortDescription sarifMessage   `json:"shortDescription"`
	Properties       map[string]any `json:"properties"`
}

type sarifResult struct {
	RuleID     string          `json:"ruleId"`
	Level      string          `json:"level"`
	Message    sarifMessage    `json:"message"`
	Locations  []sarifLocation `json:"locations,omitempty"`
	Properties map[string]any  `json:"properties"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           *sarifRegion          `json:"region,omitempty"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
}

func sarifLevel(severity string) string {
	switch severity {
	case "critical", "high":
		return "error"
	case "medium":
		return "warning"
	default:
		return "note"
	}
}

func writeSARIF(output io.Writer, run model.RunResult) error {
	rules := make([]sarifRule, 0)
	results := make([]sarifResult, 0)
	seenRules := map[string]bool{}
	for _, finding := range run.Findings {
		if !seenRules[finding.AssertionID] {
			seenRules[finding.AssertionID] = true
			rules = append(rules, sarifRule{
				ID: finding.AssertionID, ShortDescription: sarifMessage{Text: finding.Title},
				Properties: map[string]any{"severity": finding.Severity, "gate": finding.Gate, "control_ids": finding.ControlIDs},
			})
		}
		locations := make([]sarifLocation, 0, len(finding.Locations))
		for _, location := range finding.Locations {
			physical := sarifPhysicalLocation{ArtifactLocation: sarifArtifactLocation{URI: location.Path}}
			if location.Line > 0 {
				physical.Region = &sarifRegion{StartLine: location.Line, StartColumn: location.Column}
			}
			locations = append(locations, sarifLocation{PhysicalLocation: physical})
		}
		results = append(results, sarifResult{
			RuleID: finding.AssertionID, Level: sarifLevel(finding.Severity),
			Message: sarifMessage{Text: finding.Summary}, Locations: locations,
			Properties: map[string]any{
				"finding_id": finding.ID, "fingerprint": finding.Fingerprint,
				"severity": finding.Severity, "gate": finding.Gate,
				"remediation_class": finding.RemediationClass, "evidence_ids": finding.EvidenceIDs,
			},
		})
	}
	log := sarifLog{
		Version: "2.1.0", Schema: "https://json.schemastore.org/sarif-2.1.0.json",
		Runs: []sarifRun{{
			Tool:    sarifTool{Driver: sarifDriver{Name: "Production Readiness Scanner", InformationURI: "https://marinjursic.github.io/production-readiness-checklist/", Rules: rules}},
			Results: results,
			Properties: map[string]string{
				"run_id": run.RunID, "profile": run.Plan.ProfileID + "@" + run.Plan.ProfileVersion,
				"inventory_digest": run.Inventory.Digest, "terminal_state": run.TerminalState,
				"adapter_execution_count": strconv.Itoa(len(run.AdapterExecutions)),
			},
		}},
	}
	if run.Plan.ConfigurationDigest != "" {
		log.Runs[0].Properties["configuration_digest"] = run.Plan.ConfigurationDigest
		log.Runs[0].Properties["project_id"] = run.Plan.ProjectID
		log.Runs[0].Properties["target_environments"] = strings.Join(run.Plan.TargetEnvironments, ",")
		log.Runs[0].Properties["artifact_digests"] = strings.Join(run.Plan.ArtifactDigests, ",")
	}
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	return encoder.Encode(log)
}

type junitSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Name       string          `xml:"name,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Skipped    int             `xml:"skipped,attr"`
	Time       string          `xml:"time,attr"`
	Properties []junitProperty `xml:"properties>property"`
	Cases      []junitCase     `xml:"testcase"`
}

type junitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type junitCase struct {
	Name      string        `xml:"name,attr"`
	Classname string        `xml:"classname,attr"`
	Failure   *junitMessage `xml:"failure,omitempty"`
	Error     *junitMessage `xml:"error,omitempty"`
	Skipped   *junitMessage `xml:"skipped,omitempty"`
}

type junitMessage struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr,omitempty"`
	Text    string `xml:",chardata"`
}

func writeJUnit(output io.Writer, run model.RunResult) error {
	suite := junitSuite{
		Name: run.Plan.ProfileID + "@" + run.Plan.ProfileVersion, Tests: len(run.Results),
		Time: strconv.FormatFloat(run.CompletedAt.Sub(run.StartedAt).Seconds(), 'f', 3, 64),
		Properties: []junitProperty{
			{Name: "run_id", Value: run.RunID}, {Name: "inventory_digest", Value: run.Inventory.Digest},
			{Name: "terminal_state", Value: run.TerminalState},
			{Name: "adapter_execution_count", Value: strconv.Itoa(len(run.AdapterExecutions))},
		},
	}
	if run.Plan.ConfigurationDigest != "" {
		suite.Properties = append(suite.Properties,
			junitProperty{Name: "configuration_digest", Value: run.Plan.ConfigurationDigest},
			junitProperty{Name: "project_id", Value: run.Plan.ProjectID},
			junitProperty{Name: "target_environments", Value: strings.Join(run.Plan.TargetEnvironments, ",")},
			junitProperty{Name: "artifact_digests", Value: strings.Join(run.Plan.ArtifactDigests, ",")},
		)
	}
	for _, result := range run.Results {
		item := junitCase{Name: result.AssertionID, Classname: run.Plan.ProfileID}
		detail := &junitMessage{Message: result.Summary, Type: result.Assessment, Text: result.Summary}
		switch {
		case result.Assessment == "fail":
			item.Failure = detail
			suite.Failures++
		case result.Execution == "error" || result.Execution == "blocked" || result.Assessment == "unknown" || result.Assessment == "stale" || result.Assessment == "conflicting":
			item.Error = detail
			suite.Errors++
		case result.Assessment == "manual_review" || result.Assessment == "not_applicable":
			item.Skipped = detail
			suite.Skipped++
		}
		suite.Cases = append(suite.Cases, item)
	}
	if _, err := io.WriteString(output, xml.Header); err != nil {
		return err
	}
	encoder := xml.NewEncoder(output)
	encoder.Indent("", "  ")
	if err := encoder.Encode(suite); err != nil {
		return err
	}
	_, err := io.WriteString(output, "\n")
	return err
}

const htmlReport = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Production readiness assessment</title>
  <style>
    :root { color-scheme: light; --ink: #17202a; --muted: #52606d; --line: #c8d0d8; --panel: #f7f9fb; --accent: #175cd3; }
    body { color: var(--ink); font: 16px/1.55 system-ui, sans-serif; margin: 2rem auto; max-width: 90rem; padding: 0 1rem; }
    code { overflow-wrap: anywhere; } table { border-collapse: collapse; width: 100%; }
    caption { font-size: 1.25rem; font-weight: 700; margin: 1rem 0; text-align: left; }
    th, td { border: 1px solid #c8d0d8; padding: .6rem; text-align: left; vertical-align: top; }
    th { background: #edf2f7; } .status { font-weight: 700; }
    .notice { background: #eef5ff; border-left: .35rem solid var(--accent); padding: .8rem 1rem; }
    .counts { display: flex; flex-wrap: wrap; gap: .75rem; list-style: none; padding: 0; }
    .counts li { background: var(--panel); border: 1px solid var(--line); border-radius: .4rem; padding: .6rem .9rem; }
    .card, details { border: 1px solid var(--line); border-radius: .5rem; margin: 1rem 0; padding: 1rem; }
    .card h3 { margin-top: 0; } .meta { color: var(--muted); }
    summary { cursor: pointer; font-weight: 700; } dt { font-weight: 700; } dd { margin-bottom: .5rem; }
    ul { padding-left: 1.4rem; }
  </style>
</head>
<body>
  <main>
    <h1>Production readiness assessment</h1>
    <dl>
      <dt>Run</dt><dd><code>{{.Run.RunID}}</code></dd>
      <dt>Profile</dt><dd><code>{{.Run.Plan.ProfileID}}@{{.Run.Plan.ProfileVersion}}</code></dd>
      <dt>Target</dt><dd>{{.Run.Inventory.TargetName}}</dd>
      <dt>Inventory</dt><dd><code>{{.Run.Inventory.Digest}}</code></dd>
      {{if .Run.Plan.ConfigurationDigest}}<dt>Configuration</dt><dd><code>{{.Run.Plan.ConfigurationDigest}}</code></dd>
      <dt>Project</dt><dd><code>{{.Run.Plan.ProjectID}}</code></dd>
      <dt>Target environments</dt><dd>{{join .Run.Plan.TargetEnvironments}}</dd>
      <dt>Artifact digests</dt><dd>{{join .Run.Plan.ArtifactDigests}}</dd>{{end}}
      <dt>Terminal state</dt><dd class="status">{{.Run.TerminalState}}</dd>
    </dl>
    <p class="notice"><strong>Report-only scan:</strong> this command assessed the project and did not apply fixes. Results are scoped to the profile, target inventory, and evidence named above.</p>
    <h2>Assessment summary</h2>
    <ul class="counts">{{range .Counts}}<li><strong>{{index . 0}}</strong>: {{index . 1}}</li>{{end}}<li><strong>findings</strong>: {{len .Run.Findings}}</li></ul>
	{{if .Run.ControlCatalog}}<h2>Complete control catalog</h2>
	<p class="notice"><strong>All {{.Run.ControlCatalog.ControlCount}} registered controls are included.</strong> Deterministic assertions prove only their exact, narrow statements. A <code>partially_verified</code> result is not a complete Pass, and an AI review is advisory only.</p>
	<dl>
	  <dt>Registry version / digest</dt><dd><code>{{.Run.ControlCatalog.RegistryVersion}}</code> / <code>{{.Run.ControlCatalog.RegistrySHA256}}</code></dd>
	  <dt>Source digest</dt><dd><code>{{.Run.ControlCatalog.SourceSHA256}}</code></dd>
	  <dt>Controls</dt><dd>{{.Run.ControlCatalog.ControlCount}} registered / {{.Run.ControlCatalog.ActiveControlCount}} active</dd>
	  <dt>Narrow profile state before expansion</dt><dd>{{.Run.ControlCatalog.ProfileTerminalState}}</dd>
	</dl>
	<ul class="counts">{{range .ControlCounts}}<li><strong>{{index . 0}}</strong>: {{index . 1}}</li>{{end}}</ul>
	<p>Use browser search to find a control ID or phrase. Open a control to see its source, coverage, linked assertions, and any advisory AI review.</p>
	{{range .Run.ControlResults}}<details {{if eq .Disposition "confirmed_failure"}}open{{end}}>
	  <summary>[{{.Disposition}}] <code>{{.ControlID}}</code> — {{.Statement}}</summary>
	  <dl>
	    <dt>Coverage / authority</dt><dd>{{.Coverage}} / {{.Authority}}</dd>
	    <dt>Source</dt><dd><code>{{.Source.Path}}:{{.Source.Line}}</code></dd>
	    <dt>Result explanation</dt><dd>{{.Summary}}</dd>
	    <dt>All known assertions</dt><dd>{{if .AssertionIDs}}<code>{{join .AssertionIDs}}</code>{{else}}None yet.{{end}}</dd>
	    <dt>Assertions executed in this profile</dt><dd>{{if .ExecutedAssertionIDs}}<code>{{join .ExecutedAssertionIDs}}</code>{{else}}None.{{end}}</dd>
	    {{if .AIReview}}<dt>Advisory AI review</dt><dd><strong>{{.AIReview.Provider}}</strong>: {{.AIReview.AssessmentCandidate}} / {{.AIReview.ApplicabilityCandidate}} / confidence {{.AIReview.Confidence}}. {{.AIReview.Reason}} {{.AIReview.Advice}} This cannot create a verified Pass.</dd>{{end}}
	  </dl>
	</details>{{end}}{{end}}
    {{if .Run.AdapterExecutions}}<table>
      <caption>Adapter executions</caption>
      <thead><tr><th scope="col">Adapter</th><th scope="col">Manifest</th><th scope="col">Authorization</th><th scope="col">Trust</th><th scope="col">Registry</th><th scope="col">Status</th><th scope="col">Execution</th></tr></thead>
      <tbody>{{range .Run.AdapterExecutions}}<tr><td><code>{{.AdapterID}}</code></td><td><code>{{.ManifestSHA256}}</code></td><td>{{.Resolution.Source}}</td><td>{{.Resolution.Trust}}</td><td><code>{{.Resolution.RegistryID}}</code></td><td>{{.Transcript.Summary.Status}}</td><td><code>{{.ExecutionID}}</code></td></tr>{{end}}</tbody>
    </table>{{end}}
    <h2>Detailed findings</h2>
    {{if .Run.Findings}}{{range .Run.Findings}}<article class="card">
      <h3>{{.Title}}</h3>
      <p class="meta"><strong>{{.Severity}}</strong> severity · {{.Gate}} gate · remediation class {{.RemediationClass}}</p>
      <p>{{.Summary}}</p>
      <dl>
        <dt>Assertion</dt><dd><code>{{.AssertionID}}</code></dd>
        <dt>Control IDs</dt><dd>{{if .ControlIDs}}<code>{{join .ControlIDs}}</code>{{else}}—{{end}}</dd>
        <dt>Next action</dt><dd>{{remediation .RemediationClass}}</dd>
        <dt>Locations</dt><dd>{{if .Locations}}<ul>{{range .Locations}}<li><code>{{location .}}</code></li>{{end}}</ul>{{else}}No source location was emitted.{{end}}</dd>
        <dt>Evidence IDs</dt><dd>{{if .EvidenceIDs}}<ul>{{range .EvidenceIDs}}<li><code>{{.}}</code></li>{{end}}</ul>{{else}}No evidence identifier was attached.{{end}}</dd>
        <dt>Finding ID</dt><dd><code>{{.ID}}</code></dd>
        <dt>Stable fingerprint</dt><dd><code>{{.Fingerprint}}</code></dd>
      </dl>
    </article>{{end}}{{else}}<p>No verified failure was converted into an actionable finding. Review incomplete and manual assertion results below.</p>{{end}}
    <h2>All assertion results</h2>
    <p>Open each result to see applicability, execution state, required evidence, observed evidence, source locations, controls, and remediation class.</p>
    {{range .Run.Results}}<details {{if eq .Assessment "fail"}}open{{end}}>
      <summary>[{{.Assessment}}] <code>{{.AssertionID}}</code> — {{.Summary}}</summary>
      <dl>
        <dt>Applicability / execution</dt><dd>{{.Applicability}} / {{.Execution}}</dd>
        <dt>Severity / gate</dt><dd>{{.Severity}} / {{.Gate}}</dd>
        <dt>Control IDs</dt><dd>{{if .ControlIDs}}<code>{{join .ControlIDs}}</code>{{else}}—{{end}}</dd>
        <dt>Remediation class</dt><dd>{{.RemediationClass}}</dd>
        <dt>Locations</dt><dd>{{if .Locations}}<ul>{{range .Locations}}<li><code>{{location .}}</code></li>{{end}}</ul>{{else}}No source location was emitted.{{end}}</dd>
        <dt>Required evidence</dt><dd>{{if .EvidenceRequired}}<ul>{{range .EvidenceRequired}}<li><strong>{{.Kind}}</strong> (minimum authority: {{.MinimumAuthority}}): {{.Description}}</li>{{end}}</ul>{{else}}No evidence requirement was declared.{{end}}</dd>
        <dt>Observed evidence</dt><dd>{{if .EvidenceObserved}}<ul>{{range .EvidenceObserved}}<li><strong>{{.Kind}}</strong> from <code>{{.Source}}</code>: {{.Summary}} <span class="meta">(<code>{{.ID}}</code>)</span></li>{{end}}</ul>{{else}}No evidence was observed.{{end}}</dd>
      </dl>
    </details>{{end}}
    <p>This report is scoped to the named profile, target inventory, and evidence set. It is not an unqualified production-readiness or compliance claim.</p>
  </main>
</body>
</html>
`

func writeHTML(output io.Writer, run model.RunResult) error {
	view := struct {
		Run           model.RunResult
		Counts        [][2]string
		ControlCounts [][2]string
	}{Run: run, Counts: assessmentCounts(run.Results), ControlCounts: controlDispositionCounts(run.ControlResults)}
	return template.Must(template.New("report").Funcs(template.FuncMap{
		"join": func(values []string) string { return strings.Join(values, ", ") },
		"location": func(value model.FindingLocation) string {
			result := value.Path
			if value.Line > 0 {
				result += ":" + strconv.Itoa(value.Line)
				if value.Column > 0 {
					result += ":" + strconv.Itoa(value.Column)
				}
			}
			return result
		},
		"remediation": func(class string) string {
			switch class {
			case "R0":
				return "Review and resolve this manually; the scanner does not author a change for this class."
			case "R1":
				return "A deterministic, behavior-preserving candidate may be available through the separate prc fix workflow."
			case "R2":
				return "Use an isolated agent-authored candidate only with independent deterministic verification."
			case "R3":
				return "Treat this as a dependency or build-behavior change requiring explicit opt-in and stronger verification."
			case "R4":
				return "Infrastructure or deployment-definition changes require human authorization and environment-specific validation."
			case "R5":
				return "External staging-system changes require a separate connector policy and approval."
			case "R6":
				return "Production, destructive, legal, financial, and risk decisions are prohibited in the general remediation loop."
			default:
				return "Review the finding and its evidence before selecting a remediation path."
			}
		},
	}).Parse(htmlReport)).Execute(output, view)
}
