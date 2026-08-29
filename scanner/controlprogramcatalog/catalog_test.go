package controlprogramcatalog

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/controlbinding"
	"github.com/MarinJursic/production-readiness-checklist/scanner/controlprogram"
)

func repositoryRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate test source")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func TestLoadRealCatalogAndImmutableProgramBinding(t *testing.T) {
	catalog, err := Load(repositoryRoot(t))
	if err != nil {
		t.Fatalf("load real program catalog: %v", err)
	}
	if catalog.ControlCount() != expectedControlCount || catalog.TemplateCount() != expectedTemplateCount {
		t.Fatalf("coverage = %d/%d", catalog.ControlCount(), catalog.TemplateCount())
	}
	for _, digest := range []string{catalog.Digest(), catalog.CatalogSHA256(), catalog.ProgramSchemaSHA256(), catalog.BindingCatalogSHA256(), catalog.DefinitionSchemaSHA256(), catalog.DefinitionCorpusSHA256()} {
		if !isDigest(digest) {
			t.Fatalf("invalid exposed digest %q", digest)
		}
	}
	templates := catalog.Templates()
	if len(templates) != expectedTemplateCount {
		t.Fatalf("templates = %d", len(templates))
	}
	first := templates[0]
	if len(first.RawFactContracts) == 0 || first.Predicate.Op == "" || len(first.RequiredRuntimeOps) == 0 {
		t.Fatalf("first exact predicate is incomplete: %#v", first)
	}
	lookedUp, ok := catalog.Lookup(first.ControlID, first.ClauseID)
	if !ok || lookedUp.TemplateID != first.TemplateID {
		t.Fatal("lookup failed")
	}
	lookedUp.ReviewReason = "mutated"
	lookedUp.RawFactContracts[0].FactKey = "mutated"
	mutateFirstFact(&lookedUp.Predicate)
	again, ok := catalog.Lookup(first.ControlID, first.ClauseID)
	if !ok || again.ReviewReason == "mutated" || again.RawFactContracts[0].FactKey == "mutated" {
		t.Fatal("caller mutation escaped immutable catalog")
	}
	if _, ok := catalog.Lookup(first.ControlID, strings.Repeat("0", 64)); ok {
		t.Fatal("unknown clause lookup succeeded")
	}

	runtimeBinding := RuntimeBinding{
		SubjectID: "project", Subjects: []string{"project"}, InventorySHA256: strings.Repeat("1", 64),
		ApplicabilityProofContractSHA256: strings.Repeat("2", 64), MaximumEvidenceAgeSeconds: 60,
		ScannerInventoryParameters: map[string]controlprogram.Parameter{}, AuthenticatedPolicyParameters: map[string]controlprogram.Parameter{},
		AuthenticatedContextParameters: map[string]controlprogram.Parameter{},
	}
	for _, contract := range again.SealedParameterContracts {
		value, err := placeholderParameter(contract.ParameterType)
		if err != nil {
			t.Fatal(err)
		}
		switch contract.ValueOrigin {
		case ParameterOriginScannerInventory:
			runtimeBinding.ScannerInventoryParameters[contract.ParameterKey] = value
		case ParameterOriginPolicy:
			runtimeBinding.AuthenticatedPolicyParameters[contract.ParameterKey] = value
		case ParameterOriginContext:
			runtimeBinding.AuthenticatedContextParameters[contract.ParameterKey] = value
		default:
			t.Fatalf("unsupported test parameter origin %q", contract.ValueOrigin)
		}
	}
	program, err := again.Program(runtimeBinding)
	if err != nil {
		t.Fatalf("bind runtime program: %v", err)
	}
	if program.ControlID != again.ControlID || program.ClauseID != again.ClauseID || program.Predicate.Op != again.Predicate.Op {
		t.Fatal("runtime program lost template bindings")
	}
	program.Subjects[0] = "changed"
	program2, err := again.Program(runtimeBinding)
	if err != nil || program2.Subjects[0] != "project" {
		t.Fatal("runtime program instances share state")
	}
	if len(again.SealedParameterContracts) == 0 {
		t.Fatal("fixture template has no trust-lane parameter")
	}
	wrongLane := runtimeBinding
	wrongLane.ScannerInventoryParameters = cloneParameters(runtimeBinding.ScannerInventoryParameters)
	wrongLane.AuthenticatedPolicyParameters = cloneParameters(runtimeBinding.AuthenticatedPolicyParameters)
	wrongLane.AuthenticatedContextParameters = cloneParameters(runtimeBinding.AuthenticatedContextParameters)
	contract := again.SealedParameterContracts[0]
	var value controlprogram.Parameter
	switch contract.ValueOrigin {
	case ParameterOriginScannerInventory:
		value = wrongLane.ScannerInventoryParameters[contract.ParameterKey]
		delete(wrongLane.ScannerInventoryParameters, contract.ParameterKey)
		wrongLane.AuthenticatedPolicyParameters[contract.ParameterKey] = value
	case ParameterOriginPolicy:
		value = wrongLane.AuthenticatedPolicyParameters[contract.ParameterKey]
		delete(wrongLane.AuthenticatedPolicyParameters, contract.ParameterKey)
		wrongLane.AuthenticatedContextParameters[contract.ParameterKey] = value
	case ParameterOriginContext:
		value = wrongLane.AuthenticatedContextParameters[contract.ParameterKey]
		delete(wrongLane.AuthenticatedContextParameters, contract.ParameterKey)
		wrongLane.ScannerInventoryParameters[contract.ParameterKey] = value
	}
	if _, err := again.Program(wrongLane); err == nil {
		t.Fatal("parameter supplied through the wrong trust lane was accepted")
	}
}

func mutateFirstFact(expression *controlprogram.Expression) bool {
	if expression.Fact != "" {
		expression.Fact = "mutated"
		return true
	}
	if expression.Arg != nil && mutateFirstFact(expression.Arg) {
		return true
	}
	for index := range expression.Args {
		if mutateFirstFact(&expression.Args[index]) {
			return true
		}
	}
	return false
}

func TestProgramRejectsMissingOrUnexpectedSealedParameters(t *testing.T) {
	catalog, err := Load(repositoryRoot(t))
	if err != nil {
		t.Fatal(err)
	}
	template := catalog.Templates()[0]
	if _, err := template.Program(RuntimeBinding{}); err == nil {
		t.Fatal("empty runtime binding was accepted")
	}
	parameters := map[string]controlprogram.Parameter{"unexpected": {Type: controlprogram.FactString, String: stringPointer("value")}}
	if _, err := template.Program(RuntimeBinding{SubjectID: "p", Subjects: []string{"p"}, InventorySHA256: strings.Repeat("1", 64), ApplicabilityProofContractSHA256: strings.Repeat("2", 64), MaximumEvidenceAgeSeconds: 1, AuthenticatedPolicyParameters: parameters}); err == nil {
		t.Fatal("unexpected sealed parameter was accepted")
	}
}

func stringPointer(value string) *string { return &value }

func TestDirectedPairParametersAreSupportedAndCloned(t *testing.T) {
	parameter, err := placeholderParameter(controlprogram.FactDirectedGraph)
	if err != nil || len(parameter.Edges) != 1 {
		t.Fatalf("directed pair placeholder: parameter=%#v err=%v", parameter, err)
	}
	original := map[string]controlprogram.Parameter{"pairs": parameter}
	cloned := cloneParameters(original)
	clonedParameter := cloned["pairs"]
	clonedParameter.Edges[0].From = "changed"
	if original["pairs"].Edges[0].From == "changed" {
		t.Fatal("directed pair parameter slice was not defensively copied")
	}
}

func TestDecodeRejectsMalformedAndTamperedCatalogs(t *testing.T) {
	root := repositoryRoot(t)
	data := validCatalogData(t)
	schemaData, _ := os.ReadFile(filepath.Join(root, programSchemaRelativePath))
	bindings, err := controlbinding.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := decodeAndValidate(data, sha256Hex(schemaData), bindings); err != nil {
		t.Fatalf("valid catalog rejected: %v", err)
	}
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"duplicate key", bytes.Replace(data, []byte(`"schema_version":`), []byte(`"schema_version":"duplicate","schema_version":`), 1), "duplicate"},
		{"unknown field", bytes.Replace(data, []byte(`"schema_version":`), []byte(`"unknown":true,"schema_version":`), 1), "unknown field"},
		{"trailing JSON", append(append([]byte(nil), data...), []byte(` {}`)...), "trailing"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := decodeAndValidate(test.data, sha256Hex(schemaData), bindings)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error=%v want=%q", err, test.want)
			}
		})
	}
	if _, err := decodeAndValidate(data, strings.Repeat("f", 64), bindings); err == nil || !strings.Contains(err.Error(), "stale envelope") {
		t.Fatalf("wrong schema digest accepted: %v", err)
	}

	var document rawDocument
	if err := decodeStrict(data, &document); err != nil {
		t.Fatal(err)
	}
	document.Templates[0], document.Templates[1] = document.Templates[1], document.Templates[0]
	if _, err := decodeAndValidate(remarshalDocument(t, document), sha256Hex(schemaData), bindings); err == nil || !strings.Contains(err.Error(), "order") {
		t.Fatalf("unordered templates accepted: %v", err)
	}

	if err := decodeStrict(data, &document); err != nil {
		t.Fatal(err)
	}
	document.BindingCatalogSHA256 = strings.Repeat("a", 64)
	if _, err := decodeAndValidate(remarshalDocument(t, document), sha256Hex(schemaData), bindings); err == nil || !strings.Contains(err.Error(), "binding catalog") {
		t.Fatalf("binding mismatch accepted: %v", err)
	}

	if err := decodeStrict(data, &document); err != nil {
		t.Fatal(err)
	}
	var first rawTemplate
	if err := decodeStrict(document.Templates[0], &first); err != nil {
		t.Fatal(err)
	}
	originalShape := first.PredicateShape
	first.PredicateShape = "execute"
	first.RequiredRuntimeOps = []controlprogram.Operation{"execute"}
	first.Predicate = json.RawMessage(`{"op":"execute","fact":"raw_value"}`)
	first.TemplateSHA256 = ""
	unsigned, _ := json.Marshal(first)
	first.TemplateSHA256, _ = unsignedDigest(unsigned, "template_sha256")
	document.Templates[0], _ = json.Marshal(first)
	document.PredicateShapeCounts[originalShape]--
	if document.PredicateShapeCounts[originalShape] == 0 {
		delete(document.PredicateShapeCounts, originalShape)
	}
	document.PredicateShapeCounts["execute"]++
	if _, err := decodeAndValidate(remarshalDocument(t, document), sha256Hex(schemaData), bindings); err == nil || !strings.Contains(err.Error(), "predicate") {
		t.Fatalf("unsupported predicate accepted: %v", err)
	}
}

func remarshalDocument(t *testing.T, document rawDocument) []byte {
	t.Helper()
	document.CatalogSHA256 = ""
	unsigned, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	document.CatalogSHA256, err = unsignedDigest(unsigned, "catalog_sha256")
	if err != nil {
		t.Fatal(err)
	}
	result, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func validCatalogData(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(repositoryRoot(t), catalogRelativePath))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func TestReadCatalogDocumentPlainGzipAndHardening(t *testing.T) {
	payload := []byte(`{"catalog":"bounded"}`)
	t.Run("plain", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "programs.json")
		if err := os.WriteFile(path, payload, 0o600); err != nil {
			t.Fatal(err)
		}
		got, err := readCatalogDocument(path)
		if err != nil || !bytes.Equal(got, payload) {
			t.Fatalf("got %q, %v", got, err)
		}
	})
	t.Run("gzip", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "programs.json")
		writeGzip(t, path+".gz", payload)
		got, err := readCatalogDocument(path)
		if err != nil || !bytes.Equal(got, payload) {
			t.Fatalf("got %q, %v", got, err)
		}
	})
	t.Run("ambiguous", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "programs.json")
		_ = os.WriteFile(path, payload, 0o600)
		writeGzip(t, path+".gz", payload)
		if _, err := readCatalogDocument(path); err == nil || !strings.Contains(err.Error(), "ambiguous") {
			t.Fatalf("ambiguous files accepted: %v", err)
		}
	})
	t.Run("symlink", func(t *testing.T) {
		directory := t.TempDir()
		target := filepath.Join(directory, "target")
		path := filepath.Join(directory, "programs.json")
		_ = os.WriteFile(target, payload, 0o600)
		if err := os.Symlink(target, path); err != nil {
			t.Skip(err)
		}
		if _, err := readCatalogDocument(path); err == nil || !strings.Contains(err.Error(), "regular file") {
			t.Fatalf("symlink accepted: %v", err)
		}
	})
	t.Run("trailing gzip", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "programs.json")
		writeGzip(t, path+".gz", payload)
		file, _ := os.OpenFile(path+".gz", os.O_APPEND|os.O_WRONLY, 0)
		_, _ = file.Write([]byte("trailing"))
		_ = file.Close()
		if _, err := readCatalogDocument(path); err == nil || !strings.Contains(err.Error(), "trailing") {
			t.Fatalf("trailing data accepted: %v", err)
		}
	})
	t.Run("expanded limit", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "programs.json")
		writeGzip(t, path+".gz", bytes.Repeat([]byte{'x'}, int(maximumExpandedBytes+1)))
		if _, err := readCatalogDocument(path); err == nil || !strings.Contains(err.Error(), "expanded") {
			t.Fatalf("oversize expansion accepted: %v", err)
		}
	})
}

func writeGzip(t *testing.T, path string, data []byte) {
	t.Helper()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	writer := gzip.NewWriter(file)
	if _, err := writer.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
}
