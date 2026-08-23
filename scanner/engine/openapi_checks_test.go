package engine

import (
	"strconv"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func TestOpenAPIRootCheckSupportsPublishedFeatureVersions(t *testing.T) {
	for name, testCase := range map[string]struct {
		path     string
		document string
	}{
		"3.0 YAML": {"api/openapi.yaml", "openapi: 3.0.4\ninfo: {title: Example, version: 1.0.0}\npaths: {}\n"},
		"3.1 YAML": {"api/openapi.yaml", "openapi: 3.1.2\ninfo: {title: Example, version: 1.0.0}\ncomponents: {}\n"},
		"3.2 JSON": {"api/openapi.json", `{"openapi":"3.2.0","info":{"title":"Example","version":"1.0.0"},"webhooks":{}}` + "\n"},
	} {
		t.Run(name, func(t *testing.T) {
			root := healthyRepository(t)
			writeFixture(t, root, testCase.path, testCase.document)
			result := scanFixture(t, root)["PRC-A-API-001"]
			if result != "pass" {
				t.Fatalf("OpenAPI root assertion = %s", result)
			}
		})
	}
}

func TestOpenAPIRootCheckReportsBoundedStructuralLocations(t *testing.T) {
	root := healthyRepository(t)
	writeFixture(t, root, "api/openapi.yaml", `openapi: 3.2.0
info:
  title: ""
  version: 1
paths: []
paths: {}
`)
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	result := findResult(t, run, "PRC-A-API-001")
	if result.Execution != "completed" || result.Assessment != "fail" || len(result.Locations) < 2 ||
		!strings.Contains(result.Summary, "duplicate mapping key paths") || len(result.EvidenceObserved) != 2 {
		t.Fatalf("OpenAPI root result = %+v", result)
	}
	for _, location := range result.Locations {
		if location.Path != "api/openapi.yaml" || location.Line == 0 || location.Column == 0 {
			t.Fatalf("OpenAPI location = %+v", location)
		}
	}
	found := false
	for _, finding := range run.Findings {
		if finding.AssertionID == "PRC-A-API-001" && len(finding.Locations) == len(result.Locations) {
			found = true
		}
	}
	if !found {
		t.Fatalf("OpenAPI finding locations = %+v", run.Findings)
	}
}

func TestOpenAPIRootCheckResolvesBoundedYAMLAliases(t *testing.T) {
	root := healthyRepository(t)
	writeFixture(t, root, "api/openapi.yaml", `openapi: 3.2.0
info: {title: Example, version: 1.0.0}
components: &shared {}
webhooks: *shared
`)
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	result := findResult(t, run, "PRC-A-API-001")
	if result.Execution != "completed" || result.Assessment != "pass" {
		t.Fatalf("OpenAPI YAML alias result = %+v", result)
	}
}

func TestOpenAPIRootCheckFailsClosedOnParseAndUnsupportedVersions(t *testing.T) {
	for _, testCase := range []struct {
		name     string
		document string
		summary  string
	}{
		{"malformed", "openapi: [\n", "cannot parse api/openapi.yaml"},
		{"future feature version", "openapi: 3.3.0\ninfo: {title: Future, version: 1.0.0}\npaths: {}\n", "unsupported OpenAPI feature version"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			root := healthyRepository(t)
			writeFixture(t, root, "api/openapi.yaml", testCase.document)
			item, err := inventory.Build(root)
			if err != nil {
				t.Fatal(err)
			}
			run, err := scanner(t).Scan("prc/core-repository", item)
			if err != nil {
				t.Fatal(err)
			}
			result := findResult(t, run, "PRC-A-API-001")
			if result.Execution != "error" || result.Assessment != "unknown" || !strings.Contains(result.Summary, testCase.summary) {
				t.Fatalf("OpenAPI fail-closed result = %+v", result)
			}
		})
	}
}

func TestOpenAPIRootCheckRejectsUnprovableInputBounds(t *testing.T) {
	item := model.Inventory{Components: []model.InventoryComponent{}, Files: []model.FileRecord{}}
	for index := 0; index <= maximumOpenAPIFiles; index++ {
		path := "api/contract-" + strconv.Itoa(index) + ".yaml"
		item.Components = append(item.Components, model.InventoryComponent{
			ID: "api-description:" + path, Kind: "api-description", Path: path, Ecosystem: "openapi",
		})
		item.Files = append(item.Files, model.FileRecord{Path: path, Size: 1})
	}
	if _, err := openAPIDescriptionPaths(item); err == nil || !strings.Contains(err.Error(), "more than 256") {
		t.Fatalf("OpenAPI file limit error = %v", err)
	}
	item = model.Inventory{
		Components: []model.InventoryComponent{
			{ID: "api-description:first.yaml", Kind: "api-description", Path: "first.yaml", Ecosystem: "openapi"},
			{ID: "api-description:second.yaml", Kind: "api-description", Path: "second.yaml", Ecosystem: "openapi"},
		},
		Files: []model.FileRecord{
			{Path: "first.yaml", Size: maximumOpenAPIBytes},
			{Path: "second.yaml", Size: 1},
		},
	}
	if _, err := openAPIDescriptionPaths(item); err == nil || !strings.Contains(err.Error(), "total bytes") {
		t.Fatalf("OpenAPI byte limit error = %v", err)
	}
}
