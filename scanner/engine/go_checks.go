package engine

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

const (
	maximumGoAnalysisFiles     = 4_096
	maximumGoAnalysisBytes     = 256 * 1024 * 1024
	maximumGoReportedCallsites = 100
)

var goHTTPDefaultClientHelpers = map[string]bool{
	"Get": true, "Head": true, "Post": true, "PostForm": true,
}

type goHTTPCallsite struct {
	Path   string
	Line   int
	Column int
	Helper string
}

func goSourceFiles(inventory model.Inventory) ([]model.FileRecord, error) {
	files := make([]model.FileRecord, 0)
	var totalBytes int64
	for _, file := range inventory.Files {
		if !strings.EqualFold(filepath.Ext(file.Path), ".go") {
			continue
		}
		files = append(files, file)
		if len(files) > maximumGoAnalysisFiles {
			return nil, fmt.Errorf("%w: Go analysis found more than %d source files", errNativeInputLimit, maximumGoAnalysisFiles)
		}
		if file.Size < 0 || file.Size > maximumGoAnalysisBytes-totalBytes {
			return nil, fmt.Errorf("%w: Go analysis source exceeds %d total bytes", errNativeInputLimit, maximumGoAnalysisBytes)
		}
		totalBytes += file.Size
	}
	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	return files, nil
}

func netHTTPImportNames(file *ast.File) (map[string]bool, bool) {
	aliases := map[string]bool{}
	dotImported := false
	for _, item := range file.Imports {
		if item.Path == nil || item.Path.Value != `"net/http"` {
			continue
		}
		switch {
		case item.Name == nil:
			aliases["http"] = true
		case item.Name.Name == ".":
			dotImported = true
		case item.Name.Name != "_":
			aliases[item.Name.Name] = true
		}
	}
	return aliases, dotImported
}

func directGoHTTPCalls(path string, data []byte) ([]goHTTPCallsite, error) {
	fileset := token.NewFileSet()
	file, err := parser.ParseFile(fileset, path, data, 0)
	if err != nil {
		return nil, err
	}
	aliases, dotImported := netHTTPImportNames(file)
	if len(aliases) == 0 && !dotImported {
		return nil, nil
	}
	callsites := make([]goHTTPCallsite, 0)
	ast.Inspect(file, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}
		helper := ""
		switch function := call.Fun.(type) {
		case *ast.SelectorExpr:
			identifier, ok := function.X.(*ast.Ident)
			if ok && identifier.Obj == nil && aliases[identifier.Name] && goHTTPDefaultClientHelpers[function.Sel.Name] {
				helper = identifier.Name + "." + function.Sel.Name
			}
		case *ast.Ident:
			if dotImported && function.Obj == nil && goHTTPDefaultClientHelpers[function.Name] {
				helper = function.Name
			}
		}
		if helper != "" {
			position := fileset.Position(call.Pos())
			callsites = append(callsites, goHTTPCallsite{
				Path: path, Line: position.Line, Column: position.Column, Helper: helper,
			})
		}
		return true
	})
	return callsites, nil
}

func evaluateGoHTTPTimeouts(
	assertion model.Assertion,
	inventory model.Inventory,
	result model.AssertionResult,
	observedAt time.Time,
) model.AssertionResult {
	files, err := goSourceFiles(inventory)
	if err != nil {
		return nativeReadFailure(result, err)
	}
	callsites := make([]goHTTPCallsite, 0)
	violatingEvidence := make([]model.Evidence, 0)
	for _, file := range files {
		data, evidence, readErr := readVerifiedEvidence(
			inventory, file.Path, "go-ast-analysis", assertion.ImplementationID,
			"Parsed inventoried Go source for direct net/http package convenience calls.", observedAt,
		)
		if readErr != nil {
			return nativeReadFailure(result, readErr)
		}
		fileCallsites, parseErr := directGoHTTPCalls(file.Path, data)
		if parseErr != nil {
			return nativeReadFailure(result, fmt.Errorf("cannot parse %s as Go source: %w", file.Path, parseErr))
		}
		if len(fileCallsites) > 0 {
			callsites = append(callsites, fileCallsites...)
			if len(violatingEvidence) < maximumGoReportedCallsites {
				violatingEvidence = append(violatingEvidence, evidence)
			}
		}
	}
	scopeEvidence := inventoryEvidence(
		inventory, "go-ast-analysis", assertion.ImplementationID, ".",
		fmt.Sprintf("Bounded AST analysis inspected %d inventoried Go source file(s).", len(files)), observedAt,
	)
	result.EvidenceObserved = append([]model.Evidence{scopeEvidence}, violatingEvidence...)
	if len(callsites) == 0 {
		result.Assessment = "pass"
		result.Summary = fmt.Sprintf("Inspected %d Go source file(s); no direct net/http Get, Head, Post, or PostForm package calls backed by mutable global DefaultClient state were found.", len(files))
		return result
	}
	reported := callsites
	if len(reported) > maximumGoReportedCallsites {
		reported = reported[:maximumGoReportedCallsites]
	}
	locations := make([]string, 0, len(reported))
	for _, callsite := range reported {
		locations = append(locations, fmt.Sprintf("%s:%d:%d (%s)", callsite.Path, callsite.Line, callsite.Column, callsite.Helper))
	}
	result.Assessment = "fail"
	result.Summary = fmt.Sprintf(
		"Detected %d direct net/http package convenience call(s) whose timeout depends on mutable global DefaultClient state; reporting up to %d callsites: %s.",
		len(callsites), maximumGoReportedCallsites, strings.Join(locations, ", "),
	)
	return result
}
