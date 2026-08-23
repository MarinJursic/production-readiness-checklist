package remediation

import (
	"bufio"
	"regexp"
	"sort"
	"strings"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
	"github.com/MarinJursic/production-readiness-checklist/scanner/provider"
)

var (
	suppressionDirective = regexp.MustCompile(`(?i)(?:#\s*(?:nosec|noqa)|//\s*nolint|eslint-disable|@ts-ignore|@ts-nocheck|pragma:\s*no cover|istanbul ignore|@SuppressWarnings|pytest\.mark\.skip|unittest\.skip|\b(?:x?it|x?describe)\.skip\s*\()`)
	vacuousAssertion     = regexp.MustCompile(`(?i)^\s*(?:assert\s+(?:true|1)|assertTrue\s*\(\s*true\s*\)|expect\s*\(\s*true\s*\)\.toBe\s*\(\s*true\s*\))\s*;?\s*$`)
	emptyGoTest          = regexp.MustCompile(`(?m)func\s+Test[A-Za-z0-9_]*\s*\([^)]*\)\s*\{\s*\}`)
)

func testPath(path string) bool {
	lower := strings.ToLower(filepathToSlash(path))
	base := lower
	if index := strings.LastIndex(lower, "/"); index >= 0 {
		base = lower[index+1:]
	}
	return strings.Contains("/"+lower+"/", "/tests/") ||
		strings.Contains("/"+lower+"/", "/test/") ||
		strings.HasSuffix(base, "_test.go") || strings.HasPrefix(base, "test_") ||
		strings.Contains(base, ".test.") || strings.Contains(base, ".spec.")
}

func filepathToSlash(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func addedPatchContent(patch string) []string {
	added := make([]string, 0)
	inHunk := false
	scanner := bufio.NewScanner(strings.NewReader(patch))
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case strings.HasPrefix(line, "diff --git "):
			inHunk = false
		case strings.HasPrefix(line, "@@ "):
			inHunk = true
		case inHunk && strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++"):
			added = append(added, strings.TrimPrefix(line, "+"))
		}
	}
	return added
}

// auditProposalAntiGaming rejects proposal shapes that can manufacture a pass
// by weakening existing verification or hiding findings. It is intentionally
// conservative because the current R2 path does not run project-specific tests.
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
	added := addedPatchContent(output.Patch)
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
