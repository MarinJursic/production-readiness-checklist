package remediation

import (
	"regexp"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
	"github.com/MarinJursic/production-readiness-checklist/scanner/testdiscovery"
)

var (
	suppressionDirective = regexp.MustCompile(`(?i)(?:#\s*(?:nosec|noqa)|//\s*nolint|eslint-disable|@ts-ignore|@ts-nocheck|pragma:\s*no cover|istanbul ignore|@SuppressWarnings|pytest\.mark\.skip|unittest\.skip|\b(?:x?it|x?describe)\.skip\s*\()`)
	vacuousAssertion     = regexp.MustCompile(`(?i)^\s*(?:assert\s+(?:true|1)|assertTrue\s*\(\s*true\s*\)|expect\s*\(\s*true\s*\)\.toBe\s*\(\s*true\s*\))\s*;?\s*$`)
	emptyGoTest          = regexp.MustCompile(`(?m)func\s+Test[A-Za-z0-9_]*\s*\([^)]*\)\s*\{\s*\}`)
)

func testPath(path string) bool {
	lower := strings.ToLower(filepathToSlash(path))
	return testdiscovery.CandidatePath(lower) ||
		strings.Contains("/"+lower+"/", "/tests/") ||
		strings.Contains("/"+lower+"/", "/test/")
}

func filepathToSlash(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// auditProposalAntiGaming rejects proposal shapes that can manufacture a pass
// by weakening existing verification or hiding findings. It is intentionally
// conservative because passing one sandboxed suite does not prove behavioral
// coverage or protect against deliberate test weakening.
func auditProposalAntiGaming(baseline model.Inventory, output provider.Output) []string {
	existing := make(map[string]bool, len(baseline.Files))
	for _, record := range baseline.Files {
		existing[record.Path] = true
	}
	reasons := make([]string, 0)
	for _, path := range output.ChangedFiles {
		if existing[path] && testPath(path) {
			reasons = append(reasons, "Proposal modifies existing test or specification file "+path+".")
		}
	}
	files, parseErr := parseProviderPatch(output.Patch)
	added := make([]string, 0)
	for _, file := range files {
		for _, hunk := range file.hunks {
			for _, line := range hunk.lines {
				if line.kind == '+' {
					added = append(added, line.text)
				}
			}
		}
		if parseErr != nil || !file.newFile || !testPath(file.path) {
			continue
		}
		content, _, _, err := applyFileHunks(nil, file)
		if err != nil {
			continue
		}
		if !testdiscovery.RecognizedTest(file.path, content) {
			reasons = append(reasons, "Proposal adds test-shaped file "+file.path+" without a recognized test declaration.")
			continue
		}
		if !testdiscovery.HasBehaviorCheck(file.path, content) {
			reasons = append(reasons, "Proposal adds test file "+file.path+" without a recognized behavioral assertion or failure check.")
		}
	}
	for _, line := range added {
		if suppressionDirective.MatchString(line) {
			reasons = append(reasons, "Proposal introduces a suppression or skip directive.")
		}
		if vacuousAssertion.MatchString(line) {
			reasons = append(reasons, "Proposal introduces a constant assertion that cannot exercise target behavior.")
		}
	}
	if emptyGoTest.MatchString(strings.Join(added, "\n")) {
		reasons = append(reasons, "Proposal introduces an empty Go test.")
	}
	sort.Strings(reasons)
	return uniqueSorted(reasons)
}
