package controlbinding

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var (
	fixtureOnce sync.Once
	fixtureData []byte
	fixtureErr  error
)

func TestLoadPlainAndCompressedCatalogs(t *testing.T) {
	data := bindingFixture(t)
	digest := sha256.Sum256(data)
	expectedDigest := hex.EncodeToString(digest[:])
	for _, compressed := range []bool{false, true} {
		name := "plain"
		if compressed {
			name = "gzip"
		}
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			if compressed {
				writeCompressed(t, root, data)
			} else {
				writePlain(t, root, data)
			}
			catalog, err := Load(root)
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if catalog.Digest() != expectedDigest || catalog.BindingCount() != expectedBindingCount ||
				catalog.ClauseCount() != expectedClauseCount {
				t.Fatalf("unexpected catalog identity/counts: digest=%s bindings=%d clauses=%d",
					catalog.Digest(), catalog.BindingCount(), catalog.ClauseCount())
			}
			definitions := catalog.Definitions()
			if len(definitions) != expectedBindingCount || definitions[0].ControlID != "PRC-03-003" {
				t.Fatalf("unexpected first definition: %#v", definitions[0])
			}
			definition, ok := catalog.Definition(definitions[0].ControlID)
			if !ok || definition.SemanticSHA256 == "" || len(definition.Clauses) == 0 {
				t.Fatalf("Definition() did not return the expected binding")
			}
			definitions[0].ControlID = "mutated"
			definitions[0].Clauses[0].ResultContract.Pass.RequiresAll[0] = "mutated"
			fresh, _ := catalog.Definition(definition.ControlID)
			if fresh.ControlID != definition.ControlID || fresh.Clauses[0].ResultContract.Pass.RequiresAll[0] == "mutated" {
				t.Fatalf("Definitions() exposed mutable catalog state")
			}
			implementations := catalog.Implementations()
			implementations[0].SupportedEvidenceAuthorities[0] = "mutated"
			if catalog.Implementations()[0].SupportedEvidenceAuthorities[0] == "mutated" {
				t.Fatalf("Implementations() exposed mutable catalog state")
			}
		})
	}
}

func TestLoadRejectsMalformedUnknownAndTrailingJSON(t *testing.T) {
	data := bindingFixture(t)
	tests := []struct {
		name     string
		data     []byte
		contains string
	}{
		{name: "malformed", data: []byte("{"), contains: "parse control check bindings"},
		{
			name:     "unknown field",
			data:     bytes.Replace(data, []byte(`"schema_version":`), []byte(`"unknown":true,"schema_version":`), 1),
			contains: "unknown field",
		},
		{name: "trailing", data: append(append([]byte(nil), data...), []byte("\n{}")...), contains: "trailing JSON"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			writePlain(t, root, test.data)
			assertLoadError(t, root, test.contains)
		})
	}
}

func TestLoadRejectsAmbiguousMissingSymlinkAndOversizeFiles(t *testing.T) {
	t.Run("ambiguous forms", func(t *testing.T) {
		root := t.TempDir()
		data := bindingFixture(t)
		writePlain(t, root, data)
		writeCompressed(t, root, data)
		assertLoadError(t, root, "ambiguous")
	})
	t.Run("missing", func(t *testing.T) {
		assertLoadError(t, t.TempDir(), "file does not exist")
	})
	t.Run("symlink", func(t *testing.T) {
		root := t.TempDir()
		catalogDirectory := filepath.Join(root, "catalog")
		if err := os.MkdirAll(catalogDirectory, 0o755); err != nil {
			t.Fatal(err)
		}
		target := filepath.Join(root, "fixture.json")
		if err := os.WriteFile(target, bindingFixture(t), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(target, filepath.Join(catalogDirectory, "control-check-bindings.json")); err != nil {
			t.Fatal(err)
		}
		assertLoadError(t, root, "regular file")
	})
	t.Run("plain oversize", func(t *testing.T) {
		root := t.TempDir()
		path := plainPath(root)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.Truncate(maximumPlainBytes + 1); err != nil {
			file.Close()
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		assertLoadError(t, root, "regular file")
	})
	t.Run("compressed oversize", func(t *testing.T) {
		root := t.TempDir()
		path := plainPath(root) + ".gz"
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			t.Fatal(err)
		}
		if err := file.Truncate(maximumCompressedBytes + 1); err != nil {
			file.Close()
			t.Fatal(err)
		}
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		assertLoadError(t, root, "regular file")
	})
}

func TestLoadRejectsBrokenOrExpandingGzip(t *testing.T) {
	t.Run("checksum", func(t *testing.T) {
		root := t.TempDir()
		encoded := gzipBytes(t, bindingFixture(t))
		encoded[len(encoded)-1] ^= 0xff
		writeCompressedBytes(t, root, encoded)
		assertLoadError(t, root, "read compressed")
	})
	t.Run("expanded size", func(t *testing.T) {
		root := t.TempDir()
		data := bytes.Repeat([]byte{' '}, int(maximumExpandedBytes+1))
		writeCompressedBytes(t, root, gzipBytes(t, data))
		assertLoadError(t, root, "expanded control check bindings exceed")
	})
	t.Run("trailing member", func(t *testing.T) {
		root := t.TempDir()
		encoded := gzipBytes(t, bindingFixture(t))
		encoded = append(encoded, gzipBytes(t, []byte("appended"))...)
		writeCompressedBytes(t, root, encoded)
		assertLoadError(t, root, "trailing gzip member or data")
	})
}

func TestLoadRejectsStaleRuntimeDescriptor(t *testing.T) {
	document := bindingDocument(t)
	implementation := &document.ImplementationRegistry[0]
	implementation.CapabilityClass += "_changed"
	digest, err := implementationContractDigest(*implementation)
	if err != nil {
		t.Fatal(err)
	}
	implementation.ImplementationContractSHA256 = digest
	for bindingIndex := range document.Bindings {
		for clauseIndex := range document.Bindings[bindingIndex].Clauses {
			clause := &document.Bindings[bindingIndex].Clauses[clauseIndex]
			if clause.CheckerFamily == implementation.CheckerFamily {
				clause.ImplementationContractSHA256 = digest
			}
		}
	}
	root := writeDocument(t, document)
	assertLoadError(t, root, "does not match the closed runtime descriptor")
}

func TestLoadRejectsDuplicateControlsAndClauses(t *testing.T) {
	t.Run("control", func(t *testing.T) {
		document := bindingDocument(t)
		document.Bindings[1].ControlID = document.Bindings[0].ControlID
		root := writeDocument(t, document)
		assertLoadError(t, root, "duplicate control binding")
	})
	t.Run("clause", func(t *testing.T) {
		document := bindingDocument(t)
		duplicate := document.Bindings[0].Clauses[0]
		duplicate.Ordinal = 2
		document.Bindings[0].Clauses = append(document.Bindings[0].Clauses, duplicate)
		root := writeDocument(t, document)
		assertLoadError(t, root, "duplicate clause")
	})
}

func TestLoadRejectsInvalidClauseDigestOrdinalAndResultContract(t *testing.T) {
	tests := []struct {
		name     string
		mutate   func(*rawDocument)
		contains string
	}{
		{
			name: "digest",
			mutate: func(document *rawDocument) {
				document.Bindings[0].Clauses[0].Statement += " changed"
			},
			contains: "digest does not match",
		},
		{
			name: "ordinal",
			mutate: func(document *rawDocument) {
				document.Bindings[0].Clauses[0].Ordinal = 2
			},
			contains: "invalid clause at ordinal",
		},
		{
			name: "result contract",
			mutate: func(document *rawDocument) {
				document.Bindings[0].Clauses[0].ResultContract.Pass.RequiresAll = []string{"full_clause_proven"}
			},
			contains: "invalid result contract",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document := bindingDocument(t)
			test.mutate(&document)
			root := writeDocument(t, document)
			assertLoadError(t, root, test.contains)
		})
	}
}

func bindingFixture(t *testing.T) []byte {
	t.Helper()
	fixtureOnce.Do(func() {
		fixtureData, fixtureErr = os.ReadFile(filepath.Join("..", "..", "catalog", "control-check-bindings.json"))
	})
	if fixtureErr != nil {
		t.Fatalf("read binding fixture: %v", fixtureErr)
	}
	return append([]byte(nil), fixtureData...)
}

func bindingDocument(t *testing.T) rawDocument {
	t.Helper()
	var document rawDocument
	if err := json.Unmarshal(bindingFixture(t), &document); err != nil {
		t.Fatalf("decode binding fixture: %v", err)
	}
	return document
}

func writeDocument(t *testing.T, document rawDocument) string {
	t.Helper()
	data, err := json.Marshal(document)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	writePlain(t, root, data)
	return root
}

func writePlain(t *testing.T, root string, data []byte) {
	t.Helper()
	path := plainPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func writeCompressed(t *testing.T, root string, data []byte) {
	t.Helper()
	writeCompressedBytes(t, root, gzipBytes(t, data))
}

func writeCompressedBytes(t *testing.T, root string, encoded []byte) {
	t.Helper()
	path := plainPath(root) + ".gz"
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, encoded, 0o600); err != nil {
		t.Fatal(err)
	}
}

func gzipBytes(t *testing.T, data []byte) []byte {
	t.Helper()
	var output bytes.Buffer
	writer, err := gzip.NewWriterLevel(&output, gzip.BestCompression)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

func plainPath(root string) string {
	return filepath.Join(root, filepath.FromSlash(catalogRelativePath))
}

func assertLoadError(t *testing.T, root, contains string) {
	t.Helper()
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), contains) {
		t.Fatalf("Load() error = %v, want text %q", err, contains)
	}
}
