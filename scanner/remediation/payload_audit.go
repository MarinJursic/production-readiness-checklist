package remediation

import (
	"path/filepath"
	"regexp"
	"strings"
)

type payloadRule struct {
	category string
	pattern  *regexp.Regexp
}

var commonPayloadRules = []payloadRule{
	{"shell execution", regexp.MustCompile(`(?i)(?:/bin/(?:ba|z|fi)?sh\b|\b(?:powershell(?:\.exe)?|cmd\.exe)\b)`)},
	{"encoded payload", regexp.MustCompile(`(?i)(?:["'][0-9a-f]{128,}["']|["'][A-Za-z0-9+/]{256,}={0,2}["'])`)},
}

var pythonPayloadRules = []payloadRule{
	{"process execution", regexp.MustCompile(`(?m)^\s*(?:from\s+(?:subprocess|multiprocessing|ctypes)\s+import|import\s+(?:subprocess|multiprocessing|ctypes)(?:\s|,|$))|\b(?:subprocess\.[A-Za-z_]+|os\.(?:system|popen)|pty\.spawn)\s*\(`)},
	{"network access", regexp.MustCompile(`(?m)^\s*(?:from\s+(?:socket|urllib(?:\.[A-Za-z0-9_]+)?|requests|httpx|ftplib|telnetlib|paramiko)\s+import|import\s+(?:socket|urllib(?:\.[A-Za-z0-9_]+)?|requests|httpx|ftplib|telnetlib|paramiko)(?:\s|,|$))|\b(?:socket\.|requests\.|httpx\.|urllib\.)`)},
	{"environment access", regexp.MustCompile(`\b(?:os\.(?:getenv|environ)|dotenv\.)`)},
	{"filesystem mutation or absolute-path access", regexp.MustCompile(`\b(?:os\.(?:remove|unlink|rmdir|removedirs|rename|renames|replace)|shutil\.rmtree)\s*\(|\bopen\s*\(\s*["']/`)},
	{"dynamic or encoded code execution", regexp.MustCompile(`\b(?:eval|exec|compile)\s*\(|\b(?:base64\.(?:b64decode|decodebytes)|marshal\.loads|pickle\.loads)\s*\(`)},
}

var javascriptPayloadRules = []payloadRule{
	{"process execution", regexp.MustCompile(`(?m)(?:require\s*\(\s*["'](?:node:)?(?:child_process|cluster|worker_threads)["']|from\s*["'](?:node:)?(?:child_process|cluster|worker_threads)["'])|\b(?:Deno\.(?:Command|run)|Bun\.spawn)\s*\(`)},
	{"network access", regexp.MustCompile(`(?m)(?:require\s*\(\s*["'](?:node:)?(?:dgram|dns|http|https|net|tls)["']|from\s*["'](?:node:)?(?:dgram|dns|http|https|net|tls)["'])|\b(?:fetch|WebSocket|EventSource)\s*\(`)},
	{"environment access", regexp.MustCompile(`\b(?:process\.env|Deno\.env|Bun\.env)\b`)},
	{"filesystem access", regexp.MustCompile(`(?m)(?:require\s*\(\s*["'](?:node:)?fs(?:/promises)?["']|from\s*["'](?:node:)?fs(?:/promises)?["'])`)},
	{"dynamic or encoded code execution", regexp.MustCompile(`\b(?:eval\s*\(|new\s+Function\s*\(|Buffer\.from\s*\([^\n]*["']base64["'])`)},
}

var goPayloadRules = []payloadRule{
	{"process execution", regexp.MustCompile(`(?m)^\s*(?:import\s+)?(?:[._A-Za-z][A-Za-z0-9_]*\s+)?"(?:os/exec|plugin|syscall|unsafe)"\s*$|\b(?:exec\.Command|syscall\.Exec|plugin\.Open)\s*\(`)},
	{"network access", regexp.MustCompile(`(?m)^\s*(?:import\s+)?(?:[._A-Za-z][A-Za-z0-9_]*\s+)?"(?:net|net/http|net/rpc)"\s*$|\b(?:http\.(?:Get|Post|PostForm)|net\.Dial|rpc\.Dial)\s*\(`)},
	{"environment access", regexp.MustCompile(`\bos\.(?:Getenv|LookupEnv|Environ)\s*\(`)},
	{"filesystem mutation or absolute-path access", regexp.MustCompile(`\bos\.(?:Create|OpenFile|Remove|RemoveAll|Rename|WriteFile)\s*\(|\bos\.(?:Open|ReadFile)\s*\(\s*"/`)},
	{"dynamic or encoded code execution", regexp.MustCompile(`\b(?:base64\.(?:NewDecoder|RawStdEncoding|StdEncoding)|hex\.DecodeString)\b`)},
}

// auditTestPayload rejects capabilities that a focused generated regression
// test never needs. The candidate still runs in a deny-by-default sandbox; this
// separate content gate prevents dangerous test code from being preserved and
// executed later in a less restricted developer or CI environment.
func auditTestPayload(path string, content []byte) []string {
	rules := append([]payloadRule(nil), commonPayloadRules...)
	switch strings.ToLower(filepath.Ext(path)) {
	case ".py":
		rules = append(rules, pythonPayloadRules...)
	case ".js", ".jsx", ".ts", ".tsx":
		rules = append(rules, javascriptPayloadRules...)
	case ".go":
		rules = append(rules, goPayloadRules...)
	}
	text := string(content)
	reasons := make([]string, 0)
	for _, rule := range rules {
		if rule.pattern.MatchString(text) {
			reasons = append(reasons, "Proposal adds test file "+path+" with prohibited "+rule.category+" capability.")
		}
	}
	return uniqueSorted(reasons)
}
