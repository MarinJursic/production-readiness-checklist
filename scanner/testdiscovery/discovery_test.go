package testdiscovery

import "testing"

func TestRecognizedTestRequiresPathAndCollectableDeclaration(t *testing.T) {
	tests := []struct {
		name, path, content string
		want                bool
	}{
		{"go", "ready_test.go", "package ready\nimport \"testing\"\nfunc TestReady(t *testing.T) {}\n", true},
		{"go lowercase suffix", "ready_test.go", "package ready\nimport \"testing\"\nfunc Testready(t *testing.T) {}\n", false},
		{"go fake testing type", "ready_test.go", "package ready\ntype fake struct{}\nfunc TestReady(t *fake.T) {}\n", false},
		{"python", "tests/test_ready.py", "def test_ready():\n    assert ready()\n", true},
		{"python prose", "tests/test_ready.py", "# def test_ready() is documented here\n", false},
		{"javascript", "ready.test.js", "test('ready', () => assert.equal(ready(), true));\n", true},
		{"javascript skipped", "ready.test.js", "test.skip('ready', () => {});\n", false},
		{"javascript comment", "ready.test.js", "// test('not real', () => {});\n", false},
		{"javascript string", "ready.test.js", "const example = \"test('not real', () => {})\";\n", false},
		{"rust", "tests/ready.rs", "#[test]\nfn ready() { assert!(true); }\n", true},
		{"ordinary source", "ready.py", "def test_ready():\n    assert ready()\n", false},
		{"empty test file", "tests/ready.py", "", false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := RecognizedTest(test.path, []byte(test.content)); got != test.want {
				t.Fatalf("RecognizedTest() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestBehaviorCheckRejectsVacuousOrInvocationOnlyTests(t *testing.T) {
	tests := []struct {
		path, content string
		want          bool
	}{
		{"test_ready.py", "def test_ready():\n    ready()\n", false},
		{"test_ready.py", "def test_ready():\n    assert True\n", false},
		{"test_ready.py", "def test_ready():\n    assert ready() is True\n", true},
		{"ready.test.js", "test('ready', () => ready());\n", false},
		{"ready.test.js", "test('ready', () => assert.equal(ready(), true));\n", true},
		{"ready.test.js", "test('ready', () => { /* assert.equal(ready(), true) */ ready(); });\n", false},
		{"ready_test.go", "package ready\nimport \"testing\"\nfunc TestReady(t *testing.T) { ready() }\n", false},
		{"ready_test.go", "package ready\nimport \"testing\"\nfunc TestReady(t *testing.T) { if !ready() { t.Fatal(\"not ready\") } }\n", true},
	}
	for _, test := range tests {
		if got := HasBehaviorCheck(test.path, []byte(test.content)); got != test.want {
			t.Errorf("HasBehaviorCheck(%s) = %v, want %v", test.path, got, test.want)
		}
	}
}

func TestManifestDeclaresOnlyMeaningfulTestCommand(t *testing.T) {
	tests := []struct {
		content string
		want    bool
	}{
		{`{"scripts":{"test":"node --test"}}`, true},
		{`{"scripts":{"test":"pytest"}}`, true},
		{`{"scripts":{"test":"echo \"Error: no test specified\" && exit 1"}}`, false},
		{`{"scripts":{"test":"true"}}`, false},
		{`{"scripts":{}}`, false},
		{`not json`, false},
	}
	for _, test := range tests {
		if got := ManifestDeclaresTest("package.json", []byte(test.content)); got != test.want {
			t.Errorf("ManifestDeclaresTest(%q) = %v, want %v", test.content, got, test.want)
		}
	}
}
