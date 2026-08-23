package provider

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ErrSensitiveInput lets scanner-owned orchestration classify a remote-source
// preflight stop without inspecting an error string.
var ErrSensitiveInput = errors.New("agent input contains secret-like material")

// remoteInputPatterns intentionally contains only credential shapes with
// distinctive issuer prefixes or credential-bearing URL syntax. This is a
// preflight guard for remote model processing, not a replacement for a full
// repository secret scanner. It must prefer a false stop over disclosing an
// obvious credential, while avoiding generic "password" heuristics that would
// make ordinary source impossible to process.
var remoteInputPatterns = []struct {
	id      string
	pattern *regexp.Regexp
}{
	{"aws-access-key-id", regexp.MustCompile(`(?:^|[^A-Z0-9])(?:AKIA|ASIA)[A-Z0-9]{16}(?:$|[^A-Z0-9])`)},
	{"github-token", regexp.MustCompile(`(?:^|[^A-Za-z0-9_])(?:gh[pousr]_[A-Za-z0-9_]{36,}|github_pat_[A-Za-z0-9_]{22,})(?:$|[^A-Za-z0-9_])`)},
	{"anthropic-api-key", regexp.MustCompile(`(?:^|[^A-Za-z0-9_-])sk-ant-(?:api|admin)[0-9]{2}-[A-Za-z0-9_-]{20,}(?:$|[^A-Za-z0-9_-])`)},
	{"openai-api-key", regexp.MustCompile(`(?:^|[^A-Za-z0-9_-])sk-(?:proj-|admin-)?[A-Za-z0-9_-]{20,}(?:$|[^A-Za-z0-9_-])`)},
	{"slack-token", regexp.MustCompile(`(?:^|[^A-Za-z0-9-])xox[baprs]-[A-Za-z0-9-]{20,}(?:$|[^A-Za-z0-9-])`)},
	{"credential-bearing-url", regexp.MustCompile(`(?i)\b(?:https?|postgres(?:ql)?|mysql|mongodb(?:\+srv)?)://[^\s/@:]+:[^\s/@]+@`)},
}

var privateKeyLabels = []string{
	"PRIVATE KEY", "ENCRYPTED PRIVATE KEY", "RSA PRIVATE KEY", "EC PRIVATE KEY",
	"DSA PRIVATE KEY", "OPENSSH PRIVATE KEY", "PGP PRIVATE KEY BLOCK",
}

func validateRemoteInputs(inputs []InputFile) error {
	for _, input := range inputs {
		if containsPrivateKey(input.Content) {
			return fmt.Errorf("%w: %s (private-key)", ErrSensitiveInput, input.Path)
		}
		for _, detector := range remoteInputPatterns {
			if detector.pattern.MatchString(input.Content) {
				return fmt.Errorf("%w: %s (%s)", ErrSensitiveInput, input.Path, detector.id)
			}
		}
	}
	return nil
}

func containsPrivateKey(content string) bool {
	lines := strings.Split(content, "\n")
	for _, label := range privateKeyLabels {
		begin, end := "-----BEGIN "+label+"-----", "-----END "+label+"-----"
		for index, line := range lines {
			if strings.TrimSpace(line) != begin {
				continue
			}
			payloadBytes, valid := 0, true
			for cursor := index + 1; cursor < len(lines); cursor++ {
				line := strings.TrimSpace(lines[cursor])
				if line == end {
					if valid && payloadBytes >= 32 {
						return true
					}
					break
				}
				if line == "" || strings.Contains(line, ":") || strings.HasPrefix(line, "=") {
					continue
				}
				if !isBase64Line(line) {
					valid = false
				}
				payloadBytes += len(line)
			}
		}
	}
	return false
}

func isBase64Line(line string) bool {
	return strings.IndexFunc(line, func(character rune) bool {
		return !(character >= 'A' && character <= 'Z') &&
			!(character >= 'a' && character <= 'z') &&
			!(character >= '0' && character <= '9') &&
			character != '+' && character != '/' && character != '='
	}) < 0
}
