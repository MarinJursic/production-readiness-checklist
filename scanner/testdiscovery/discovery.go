// Package testdiscovery recognizes bounded, source-backed evidence that a
// repository contains tests a conventional runner can actually discover.
package testdiscovery

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	pythonTestDeclaration = regexp.MustCompile(`(?m)^\s*(?:async\s+)?def\s+test_[A-Za-z0-9_]*\s*\(`)
	javaTestDeclaration   = regexp.MustCompile(`(?m)@(?:org\.junit(?:\.jupiter\.api)?\.)?Test\b`)
	csharpTestDeclaration = regexp.MustCompile(`(?m)\[(?:Test|Fact|Theory|TestMethod)(?:Attribute)?(?:\([^]]*\))?\]`)
	phpTestDeclaration    = regexp.MustCompile(`(?mi)(?:#\[Test\]|function\s+test[A-Za-z0-9_]*\s*\()`)
	swiftTestDeclaration  = regexp.MustCompile(`(?m)\bfunc\s+test[A-Z_a-z0-9]*\s*\(`)
	cTestDeclaration      = regexp.MustCompile(`(?m)\b(?:TEST|TEST_F|TEST_P|TEST_CASE|SCENARIO)\s*\(`)
	rustTestDeclaration   = regexp.MustCompile(`(?ms)#\s*\[\s*test\s*\]\s*(?:async\s+)?fn\s+[A-Za-z_][A-Za-z0-9_]*\s*\(`)
	rubyTestDeclaration   = regexp.MustCompile(`(?m)(?:^\s*def\s+test_[A-Za-z0-9_]*\b|^\s*(?:it|test)\s*[\("'])`)
	jsTestDeclaration     = regexp.MustCompile(`(?m)(?:^|[^A-Za-z0-9_$.])(?:it|test)\s*\(`)
	scalaTestDeclaration  = regexp.MustCompile(`(?m)(?:^|[^A-Za-z0-9_$.])(?:it|test)\s*\(`)

	pythonBehaviorCheck = regexp.MustCompile(`(?mi)(?:\bpytest\.raises\s*\(|\bself\.assert[A-Z_a-z0-9]*\s*\(|\.assert_(?:called|awaited)[A-Za-z0-9_]*\s*\()`)
	jsBehaviorCheck     = regexp.MustCompile(`(?m)(?:\bexpect\s*\(|\bassert(?:\.[A-Za-z_][A-Za-z0-9_]*)?\s*\(|\b(?:t|context)\.(?:assert|equal|deepEqual|throws|rejects|fail)\s*\()`)
	rustBehaviorCheck   = regexp.MustCompile(`(?m)\b(?:assert|assert_eq|assert_ne|debug_assert|matches)\s*!\s*\(`)
	rubyBehaviorCheck   = regexp.MustCompile(`(?m)\b(?:assert|refute|flunk|must_[A-Za-z0-9_]+|wont_[A-Za-z0-9_]+)\b`)
	javaBehaviorCheck   = regexp.MustCompile(`(?m)\b(?:assert[A-Z_a-z0-9]*|fail|verify)\s*\(`)
	csharpBehaviorCheck = regexp.MustCompile(`(?m)\b(?:Assert\.[A-Za-z_][A-Za-z0-9_]*|CollectionAssert\.[A-Za-z_][A-Za-z0-9_]*|Should\s*\()`)
	phpBehaviorCheck    = regexp.MustCompile(`(?mi)\b(?:self::|\$this->)?assert[A-Z_a-z0-9]*\s*\(`)
	swiftBehaviorCheck  = regexp.MustCompile(`(?m)\b(?:XCTAssert[A-Za-z0-9_]*|XCTFail)\s*\(`)
	cBehaviorCheck      = regexp.MustCompile(`(?m)\b(?:ASSERT|EXPECT|CHECK|REQUIRE)(?:_[A-Z0-9_]+)?\s*\(`)
)

// CandidatePath reports whether path uses a conventional test-file location
// or name. Content still has to pass RecognizedTest.
func CandidatePath(path string) bool {
	lower := strings.ToLower(filepath.ToSlash(path))
	base := filepath.Base(lower)
	parts := strings.Split(lower, "/")
	inTestDirectory := len(parts) > 1 && (parts[0] == "tests" || parts[0] == "test" || parts[0] == "__tests__")
	return inTestDirectory || strings.HasSuffix(base, "_test.go") ||
		strings.HasPrefix(base, "test_") || strings.HasSuffix(base, "_test.py") ||
		strings.Contains(base, ".test.") || strings.Contains(base, ".spec.")
}

// RecognizedTest requires both a conventional path and a language-specific
// declaration that the corresponding test runner can collect.
func RecognizedTest(path string, data []byte) bool {
	if !CandidatePath(path) {
		return false
	}
	extension := strings.ToLower(filepath.Ext(path))
	text := codeView(extension, string(data))
	switch extension {
	case ".go":
		return recognizedGoTest(path, data)
	case ".py":
		return pythonTestDeclaration.MatchString(text)
	case ".js", ".jsx", ".ts", ".tsx":
		return jsTestDeclaration.MatchString(text)
	case ".rs":
		return rustTestDeclaration.MatchString(text)
	case ".rb":
		return rubyTestDeclaration.MatchString(text)
	case ".java", ".kt", ".kts":
		return javaTestDeclaration.MatchString(text)
	case ".cs":
		return csharpTestDeclaration.MatchString(text)
	case ".php":
		return phpTestDeclaration.MatchString(text)
	case ".swift":
		return swiftTestDeclaration.MatchString(text)
	case ".c", ".cc", ".cpp":
		return cTestDeclaration.MatchString(text)
	case ".scala":
		return scalaTestDeclaration.MatchString(text)
	default:
		return false
	}
}

// HasBehaviorCheck applies a conservative structural check used for
// scanner-generated R2 test proposals. It deliberately rejects invocation-only
// tests because the current proposal path cannot execute project test commands.
func HasBehaviorCheck(path string, data []byte) bool {
	extension := strings.ToLower(filepath.Ext(path))
	text := codeView(extension, string(data))
	switch extension {
	case ".go":
		return goBehaviorCheck(path, data)
	case ".py":
		return pythonBehaviorCheck.MatchString(text) || hasNonConstantPythonAssert(text)
	case ".js", ".jsx", ".ts", ".tsx":
		return jsBehaviorCheck.MatchString(text)
	case ".rs":
		return rustBehaviorCheck.MatchString(text)
	case ".rb":
		return rubyBehaviorCheck.MatchString(text)
	case ".java", ".kt", ".kts", ".scala":
		return javaBehaviorCheck.MatchString(text)
	case ".cs":
		return csharpBehaviorCheck.MatchString(text)
	case ".php":
		return phpBehaviorCheck.MatchString(text)
	case ".swift":
		return swiftBehaviorCheck.MatchString(text)
	case ".c", ".cc", ".cpp":
		return cBehaviorCheck.MatchString(text)
	default:
		return false
	}
}

// ManifestDeclaresTest recognizes an explicit, non-placeholder package.json
// test command. It does not infer commands from dependencies or prose.
func ManifestDeclaresTest(path string, data []byte) bool {
	if strings.ToLower(filepath.Base(path)) != "package.json" {
		return false
	}
	var document struct {
		Scripts map[string]any `json:"scripts"`
	}
	if err := json.Unmarshal(data, &document); err != nil {
		return false
	}
	value, ok := document.Scripts["test"].(string)
	if !ok {
		return false
	}
	command := strings.TrimSpace(value)
	lower := strings.ToLower(command)
	if command == "" || lower == ":" || lower == "true" || lower == "exit 0" ||
		(strings.Contains(lower, "no test specified") && strings.Contains(lower, "exit 1")) {
		return false
	}
	return true
}

func recognizedGoTest(path string, data []byte) bool {
	parsed, err := parser.ParseFile(token.NewFileSet(), path, data, 0)
	if err != nil {
		return false
	}
	aliases := goTestingAliases(parsed)
	for _, declaration := range parsed.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if ok && goTestParameter(function, aliases) != "" {
			return true
		}
	}
	return false
}

func validGoTestName(name string) bool {
	if !strings.HasPrefix(name, "Test") || len(name) == len("Test") {
		return false
	}
	next := name[len("Test")]
	return next < 'a' || next > 'z'
}

func goBehaviorCheck(path string, data []byte) bool {
	parsed, err := parser.ParseFile(token.NewFileSet(), path, data, 0)
	if err != nil {
		return false
	}
	aliases := goTestingAliases(parsed)
	for _, declaration := range parsed.Decls {
		function, ok := declaration.(*ast.FuncDecl)
		if !ok {
			continue
		}
		parameter := goTestParameter(function, aliases)
		if parameter == "" || function.Body == nil {
			continue
		}
		found := false
		ast.Inspect(function.Body, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}
			selector, ok := call.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			receiver, _ := selector.X.(*ast.Ident)
			if receiver != nil && receiver.Name == parameter {
				switch selector.Sel.Name {
				case "Error", "Errorf", "Fatal", "Fatalf", "Fail", "FailNow":
					found = true
				}
			} else {
				switch selector.Sel.Name {
				case "Equal", "NotEqual", "NoError", "ErrorIs", "Contains", "True", "False", "Nil", "NotNil", "Panics":
					found = true
				}
			}
			return !found
		})
		if found {
			return true
		}
	}
	return false
}

func goTestingAliases(file *ast.File) map[string]bool {
	aliases := map[string]bool{}
	for _, imported := range file.Imports {
		if imported.Path == nil || strings.Trim(imported.Path.Value, `"`) != "testing" {
			continue
		}
		name := "testing"
		if imported.Name != nil {
			name = imported.Name.Name
		}
		if name != "." && name != "_" {
			aliases[name] = true
		}
	}
	return aliases
}

func goTestParameter(function *ast.FuncDecl, aliases map[string]bool) string {
	if function.Recv != nil || !validGoTestName(function.Name.Name) || function.Type.Results != nil ||
		function.Type.Params == nil || len(function.Type.Params.List) != 1 {
		return ""
	}
	parameter := function.Type.Params.List[0]
	if len(parameter.Names) != 1 {
		return ""
	}
	pointer, ok := parameter.Type.(*ast.StarExpr)
	if !ok {
		return ""
	}
	selector, ok := pointer.X.(*ast.SelectorExpr)
	if !ok || selector.Sel.Name != "T" {
		return ""
	}
	pkg, ok := selector.X.(*ast.Ident)
	if !ok || !aliases[pkg.Name] {
		return ""
	}
	return parameter.Names[0].Name
}

func codeView(extension, text string) string {
	lineHash := extension == ".py" || extension == ".rb" || extension == ".php"
	lineSlash := extension != ".py" && extension != ".rb"
	block := lineSlash
	var builder strings.Builder
	builder.Grow(len(text))
	quote := byte(0)
	lineComment, blockComment, escaped := false, false, false
	for index := 0; index < len(text); index++ {
		character := text[index]
		next := byte(0)
		if index+1 < len(text) {
			next = text[index+1]
		}
		if lineComment {
			if character == '\n' {
				lineComment = false
				builder.WriteByte(character)
			} else {
				builder.WriteByte(' ')
			}
			continue
		}
		if blockComment {
			if character == '*' && next == '/' {
				builder.WriteString("  ")
				index++
				blockComment = false
			} else if character == '\n' {
				builder.WriteByte('\n')
			} else {
				builder.WriteByte(' ')
			}
			continue
		}
		if quote != 0 {
			if character == '\n' && quote != '`' {
				quote, escaped = 0, false
				builder.WriteByte('\n')
				continue
			}
			builder.WriteByte(' ')
			if escaped {
				escaped = false
			} else if character == '\\' {
				escaped = true
			} else if character == quote {
				quote = 0
			}
			continue
		}
		if lineHash && character == '#' && !(extension == ".php" && next == '[') {
			lineComment = true
			builder.WriteByte(' ')
			continue
		}
		if lineSlash && character == '/' && next == '/' {
			lineComment = true
			builder.WriteString("  ")
			index++
			continue
		}
		if block && character == '/' && next == '*' {
			blockComment = true
			builder.WriteString("  ")
			index++
			continue
		}
		if character == '\'' || character == '"' || character == '`' {
			quote = character
			builder.WriteByte(' ')
			continue
		}
		builder.WriteByte(character)
	}
	return builder.String()
}

func hasNonConstantPythonAssert(text string) bool {
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "assert ") {
			continue
		}
		expression := strings.TrimSpace(strings.TrimPrefix(trimmed, "assert "))
		if comma := strings.Index(expression, ","); comma >= 0 {
			expression = strings.TrimSpace(expression[:comma])
		}
		switch strings.ToLower(expression) {
		case "true", "false", "1", "0", "none":
			continue
		}
		if expression != "" {
			return true
		}
	}
	return false
}
