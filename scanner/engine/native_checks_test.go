package engine

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
	"github.com/MarinJursic/production-readiness-checklist/scanner/model"
)

func scanFixture(t *testing.T, root string) map[string]string {
	t.Helper()
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	results := map[string]string{}
	for _, result := range run.Results {
		results[result.AssertionID] = result.Assessment
	}
	return results
}

func TestWorkflowStructureChecksRejectUnsafeDefinitions(t *testing.T) {
	t.Run("invalid yaml", func(t *testing.T) {
		root := healthyRepository(t)
		writeFixture(t, root, ".github/workflows/ci.yml", "jobs: [\n")
		item, err := inventory.Build(root)
		if err != nil {
			t.Fatal(err)
		}
		run, err := scanner(t).Scan("prc/core-repository", item)
		if err != nil {
			t.Fatal(err)
		}
		result := findResult(t, run, "PRC-A-CORE-017")
		if result.Execution != "error" || result.Assessment != "unknown" {
			t.Fatalf("workflow validity result = %+v", result)
		}
	})

	t.Run("missing jobs", func(t *testing.T) {
		root := healthyRepository(t)
		writeFixture(t, root, ".github/workflows/ci.yml", "name: Empty\non: push\npermissions: {}\n")
		if got := scanFixture(t, root)["PRC-A-CORE-018"]; got != "fail" {
			t.Fatalf("jobs assertion = %s", got)
		}
	})

	t.Run("missing timeout", func(t *testing.T) {
		root := healthyRepository(t)
		path := filepath.Join(root, ".github", "workflows", "ci.yml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		writeFixture(t, root, ".github/workflows/ci.yml", strings.ReplaceAll(string(data), "    timeout-minutes: 15\n", ""))
		if got := scanFixture(t, root)["PRC-A-CORE-019"]; got != "fail" {
			t.Fatalf("timeout assertion = %s", got)
		}
	})

	t.Run("pull request target", func(t *testing.T) {
		root := healthyRepository(t)
		path := filepath.Join(root, ".github", "workflows", "ci.yml")
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		writeFixture(t, root, ".github/workflows/ci.yml", strings.ReplaceAll(string(data), "on: [push]", "on: [push, pull_request_target]"))
		if got := scanFixture(t, root)["PRC-A-CORE-020"]; got != "fail" {
			t.Fatalf("untrusted trigger assertion = %s", got)
		}
	})
}

func TestRepositoryIntegrityChecksFindConcreteFailures(t *testing.T) {
	t.Run("merge conflict markers", func(t *testing.T) {
		root := healthyRepository(t)
		writeFixture(t, root, "app.py", "<<<<<<< ours\nprint('one')\n=======\nprint('two')\n>>>>>>> theirs\n")
		if got := scanFixture(t, root)["PRC-A-CORE-021"]; got != "fail" {
			t.Fatalf("merge-conflict assertion = %s", got)
		}
	})

	t.Run("broadly writable file", func(t *testing.T) {
		root := healthyRepository(t)
		if err := os.Chmod(filepath.Join(root, "app.py"), 0o666); err != nil {
			t.Fatal(err)
		}
		if got := scanFixture(t, root)["PRC-A-CORE-022"]; got != "fail" {
			t.Fatalf("file-mode assertion = %s", got)
		}
	})

	t.Run("empty manifest and lock", func(t *testing.T) {
		root := healthyRepository(t)
		writeFixture(t, root, "requirements.txt", "")
		writeFixture(t, root, "requirements.lock.txt", "")
		results := scanFixture(t, root)
		if results["PRC-A-CORE-023"] != "fail" || results["PRC-A-CORE-024"] != "fail" {
			t.Fatalf("nonempty assertions = %v/%v", results["PRC-A-CORE-023"], results["PRC-A-CORE-024"])
		}
	})

	t.Run("missing runtime declaration", func(t *testing.T) {
		root := healthyRepository(t)
		if err := os.Remove(filepath.Join(root, ".python-version")); err != nil {
			t.Fatal(err)
		}
		if got := scanFixture(t, root)["PRC-A-CORE-025"]; got != "fail" {
			t.Fatalf("runtime-version assertion = %s", got)
		}
	})
}

func TestPrivateKeyArmorCheckIsBoundedAndRedacted(t *testing.T) {
	t.Run("recognized blocks require matching armor and payload", func(t *testing.T) {
		payload := strings.Repeat("Q", 64)
		for _, label := range privateKeyArmorLabels() {
			block := []byte("-----BEGIN " + label + "-----\r\n" + payload + "\r\n-----END " + label + "-----\r\n")
			if !containsPrivateKeyArmor(block, privateKeyArmorLabels()) {
				t.Errorf("recognized armor label was not detected: %s", label)
			}
		}
		for name, content := range map[string]string{
			"header only":       "-----BEGIN " + "PRIVATE KEY-----\n",
			"invalid payload":   "-----BEGIN " + "PRIVATE KEY-----\nnot-a-base64-payload\n-----END PRIVATE KEY-----\n",
			"mismatched footer": "-----BEGIN " + "PRIVATE KEY-----\n" + payload + "\n-----END PUBLIC KEY-----\n",
		} {
			if containsPrivateKeyArmor([]byte(content), privateKeyArmorLabels()) {
				t.Errorf("%s was treated as a private-key block", name)
			}
		}
	})

	t.Run("clean and public material pass", func(t *testing.T) {
		root := healthyRepository(t)
		writeFixture(t, root, "public.pem", "-----BEGIN "+"PUBLIC KEY-----\nfixture-public-material\n-----END PUBLIC KEY-----\n")
		if got := scanFixture(t, root)["PRC-A-CORE-031"]; got != "pass" {
			t.Fatalf("private-key armor assertion = %s", got)
		}
	})

	t.Run("private material fails without disclosure", func(t *testing.T) {
		root := healthyRepository(t)
		secret := strings.Repeat("U", 64)
		writeFixture(t, root, "keys/deploy.pem", "-----BEGIN "+"PRIVATE KEY-----\n"+secret+"\n-----END PRIVATE KEY-----\n")
		item, err := inventory.Build(root)
		if err != nil {
			t.Fatal(err)
		}
		run, err := scanner(t).Scan("prc/core-repository", item)
		if err != nil {
			t.Fatal(err)
		}
		result := findResult(t, run, "PRC-A-CORE-031")
		if result.Execution != "completed" || result.Assessment != "fail" ||
			!strings.Contains(result.Summary, "keys/deploy.pem") || len(result.EvidenceObserved) != 2 {
			t.Fatalf("private-key armor result = %+v", result)
		}
		encoded, err := json.Marshal(run)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(encoded), secret) {
			t.Fatal("scanner output retained private-key material")
		}
		foundLocation := false
		for _, finding := range run.Findings {
			if finding.AssertionID == "PRC-A-CORE-031" && len(finding.Locations) == 1 && finding.Locations[0].Path == "keys/deploy.pem" {
				foundLocation = true
			}
		}
		if !foundLocation {
			t.Fatalf("private-key finding has no safe source location: %+v", run.Findings)
		}
	})

	t.Run("oversized candidate blocks", func(t *testing.T) {
		root := healthyRepository(t)
		writeFixture(t, root, "large.pem", strings.Repeat("x", maximumSensitiveMaterialFileBytes+1))
		item, err := inventory.Build(root)
		if err != nil {
			t.Fatal(err)
		}
		run, err := scanner(t).Scan("prc/core-repository", item)
		if err != nil {
			t.Fatal(err)
		}
		result := findResult(t, run, "PRC-A-CORE-031")
		if result.Execution != "blocked" || result.Assessment != "unknown" || !strings.Contains(result.Summary, "per-file limit") {
			t.Fatalf("oversized private-key scan result = %+v", result)
		}
	})
}

func TestGoHTTPTimeoutCheckUsesSyntaxAwarePackageResolution(t *testing.T) {
	for name, source := range map[string]string{
		"default import": `package service

import "net/http"

func fetch() { _, _ = http.Get("https://example.invalid") }
`,
		"aliased import": `package service

import web "net/http"

func fetch() { _, _ = web.Post("https://example.invalid", "text/plain", nil) }
`,
		"dot import": `package service

import . "net/http"

func fetch() { _, _ = Head("https://example.invalid") }
`,
		"post form": `package service

import "net/http"

func fetch() { _, _ = http.PostForm("https://example.invalid", nil) }
`,
	} {
		t.Run(name, func(t *testing.T) {
			root := healthyRepository(t)
			writeFixture(t, root, "go.mod", "module example.invalid/service\n\ngo 1.27\n")
			writeFixture(t, root, "client.go", source)
			item, err := inventory.Build(root)
			if err != nil {
				t.Fatal(err)
			}
			run, err := scanner(t).Scan("prc/core-repository", item)
			if err != nil {
				t.Fatal(err)
			}
			result := findResult(t, run, "PRC-A-GO-001")
			if result.Execution != "completed" || result.Assessment != "fail" ||
				!strings.Contains(result.Summary, "client.go:") || len(result.EvidenceObserved) != 2 {
				t.Fatalf("Go HTTP timeout result = %+v", result)
			}
			foundLocation := false
			for _, finding := range run.Findings {
				if finding.AssertionID == "PRC-A-GO-001" && len(finding.Locations) == 1 &&
					finding.Locations[0].Path == "client.go" && finding.Locations[0].Line > 0 && finding.Locations[0].Column > 0 {
					foundLocation = true
				}
			}
			if !foundLocation {
				t.Fatalf("Go HTTP timeout finding has no source path: %+v", run.Findings)
			}
		})
	}

	t.Run("explicit client and shadowed names", func(t *testing.T) {
		root := healthyRepository(t)
		writeFixture(t, root, "go.mod", "module example.invalid/service\n\ngo 1.27\n")
		writeFixture(t, root, "client.go", `package service

import (
	"net/http"
	"time"
)

func fetch() {
	client := http.Client{Timeout: 10 * time.Second}
	_, _ = client.Get("https://example.invalid")
	http := struct { Get func(string) (*http.Response, error) }{}
	_, _ = http.Get("https://example.invalid")
}
`)
		result := scanFixture(t, root)["PRC-A-GO-001"]
		if result != "pass" {
			t.Fatalf("Go HTTP timeout assertion = %s", result)
		}
	})

	t.Run("unparseable source fails closed", func(t *testing.T) {
		root := healthyRepository(t)
		writeFixture(t, root, "go.mod", "module example.invalid/service\n\ngo 1.27\n")
		writeFixture(t, root, "client.go", "package service\nfunc broken( {\n")
		item, err := inventory.Build(root)
		if err != nil {
			t.Fatal(err)
		}
		run, err := scanner(t).Scan("prc/core-repository", item)
		if err != nil {
			t.Fatal(err)
		}
		result := findResult(t, run, "PRC-A-GO-001")
		if result.Execution != "error" || result.Assessment != "unknown" || !strings.Contains(result.Summary, "cannot parse client.go") {
			t.Fatalf("unparseable Go result = %+v", result)
		}
	})
}

func TestGoHTTPTimeoutCheckRejectsUnprovableInputBounds(t *testing.T) {
	files := make([]model.FileRecord, maximumGoAnalysisFiles+1)
	for index := range files {
		files[index] = model.FileRecord{Path: filepath.ToSlash(filepath.Join("src", "file"+strings.Repeat("x", index%2)+".go")), Size: 1}
	}
	if _, err := goSourceFiles(model.Inventory{Files: files}); !errors.Is(err, errNativeInputLimit) {
		t.Fatalf("file-count limit error = %v", err)
	}
	files = []model.FileRecord{
		{Path: "first.go", Size: maximumGoAnalysisBytes},
		{Path: "second.go", Size: 1},
	}
	if _, err := goSourceFiles(model.Inventory{Files: files}); !errors.Is(err, errNativeInputLimit) {
		t.Fatalf("total-byte limit error = %v", err)
	}
}

func TestContainerChecksHavePassAndFailFixtures(t *testing.T) {
	pinned := strings.Repeat("b", 64)
	root := healthyRepository(t)
	writeFixture(t, root, "Dockerfile", "FROM example@sha256:"+pinned+"\nUSER 10001\n")
	results := scanFixture(t, root)
	if results["PRC-A-CORE-026"] != "pass" || results["PRC-A-CORE-027"] != "pass" {
		t.Fatalf("container pass assertions = %v/%v", results["PRC-A-CORE-026"], results["PRC-A-CORE-027"])
	}

	writeFixture(t, root, "Dockerfile", "FROM example:latest\nUSER root\n")
	results = scanFixture(t, root)
	if results["PRC-A-CORE-026"] != "fail" || results["PRC-A-CORE-027"] != "fail" {
		t.Fatalf("container fail assertions = %v/%v", results["PRC-A-CORE-026"], results["PRC-A-CORE-027"])
	}
}

func TestTerraformLockCheckHasPassAndFailFixtures(t *testing.T) {
	root := healthyRepository(t)
	writeFixture(t, root, "infra/main.tf", "terraform { required_version = \">= 1.9\" }\n")
	writeFixture(t, root, "infra/.terraform.lock.hcl", "provider \"registry.terraform.io/example/example\" {}\n")
	if got := scanFixture(t, root)["PRC-A-CORE-028"]; got != "pass" {
		t.Fatalf("terraform lock assertion = %s", got)
	}
	if err := os.Remove(filepath.Join(root, "infra", ".terraform.lock.hcl")); err != nil {
		t.Fatal(err)
	}
	if got := scanFixture(t, root)["PRC-A-CORE-028"]; got != "fail" {
		t.Fatalf("missing terraform lock assertion = %s", got)
	}
}

func TestKubernetesChecksHavePassAndFailFixtures(t *testing.T) {
	root := healthyRepository(t)
	writeFixture(t, root, "deploy/app.yaml", `apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      securityContext:
        runAsNonRoot: true
      containers:
        - name: app
          image: example@sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
          resources:
            requests:
              cpu: 100m
            limits:
              cpu: 500m
`)
	results := scanFixture(t, root)
	if results["PRC-A-CORE-029"] != "pass" || results["PRC-A-CORE-030"] != "pass" {
		t.Fatalf("kubernetes pass assertions = %v/%v", results["PRC-A-CORE-029"], results["PRC-A-CORE-030"])
	}

	writeFixture(t, root, "deploy/app.yaml", `apiVersion: apps/v1
kind: Deployment
spec:
  template:
    spec:
      containers:
        - name: app
          image: example:latest
          securityContext:
            runAsUser: 0
`)
	results = scanFixture(t, root)
	if results["PRC-A-CORE-029"] != "fail" || results["PRC-A-CORE-030"] != "fail" {
		t.Fatalf("kubernetes fail assertions = %v/%v", results["PRC-A-CORE-029"], results["PRC-A-CORE-030"])
	}
}

func TestNativeParserRejectsMutationAfterInventory(t *testing.T) {
	root := healthyRepository(t)
	item, err := inventory.Build(root)
	if err != nil {
		t.Fatal(err)
	}
	writeFixture(t, root, ".github/workflows/ci.yml", "name: Changed\non: push\njobs: {}\n")
	run, err := scanner(t).Scan("prc/core-repository", item)
	if err != nil {
		t.Fatal(err)
	}
	result := findResult(t, run, "PRC-A-CORE-017")
	if result.Execution != "error" || result.Assessment == "pass" || !strings.Contains(result.Summary, "changed after inventory") {
		t.Fatalf("native parser mutation result = %+v", result)
	}
}
