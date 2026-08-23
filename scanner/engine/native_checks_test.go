package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/MarinJursic/production-readiness-checklist/scanner/inventory"
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
