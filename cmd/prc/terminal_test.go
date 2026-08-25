package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestTerminalTextNeutralizesControlAndBidirectionalCharacters(t *testing.T) {
	input := "safe\x1b]8;;https://attacker.invalid\aFAKE PASS\x1b]8;;\a\r\nnext\b\u202ehidden"
	output := terminalText(input)
	if strings.ContainsAny(output, "\x1b\a\r\n\b") || strings.ContainsRune(output, '\u202e') ||
		!strings.Contains(output, `\u001B`) || !strings.Contains(output, `\u202E`) {
		t.Fatalf("unsafe terminal text: %q", output)
	}
}

func TestTerminalColorPolicy(t *testing.T) {
	t.Setenv("NO_COLOR", "")
	t.Setenv("TERM", "xterm-256color")
	var output bytes.Buffer
	if newTerminalStyle("auto", &output).color || newTerminalStyle("never", &output).color {
		t.Fatal("redirected or never output unexpectedly enabled color")
	}
	if !newTerminalStyle("always", &output).color {
		t.Fatal("explicit color did not enable color")
	}
	t.Setenv("NO_COLOR", "1")
	if newTerminalStyle("always", &output).color {
		t.Fatal("NO_COLOR did not disable explicit terminal color")
	}
}

func TestAssessmentLabelsKeepPlainMeaningWithAndWithoutColor(t *testing.T) {
	cases := []struct{ assessment, execution, word string }{
		{"pass", "completed", "PASS"}, {"fail", "completed", "FAIL"},
		{"unknown", "blocked", "BLOCKED"}, {"manual_review", "completed", "MANUAL"},
		{"not_applicable", "completed", "N/A"}, {"unknown", "error", "ERROR"},
	}
	for _, item := range cases {
		plain := assessmentLabel(item.assessment, item.execution, terminalStyle{})
		colored := assessmentLabel(item.assessment, item.execution, terminalStyle{color: true})
		if !strings.Contains(plain, item.word) || strings.Contains(plain, "\x1b[") ||
			!strings.Contains(colored, item.word) || !strings.Contains(colored, "\x1b[") {
			t.Fatalf("assessment %s/%s plain=%q colored=%q", item.assessment, item.execution, plain, colored)
		}
	}
}

func TestWolfBannerWorksWithAndWithoutColor(t *testing.T) {
	for name, style := range map[string]terminalStyle{
		"plain": {}, "colored": {color: true},
	} {
		t.Run(name, func(t *testing.T) {
			var output bytes.Buffer
			printBrandBanner(&output, style)
			text := output.String()
			if !strings.Contains(text, "/\\       /\\") || !strings.Contains(text, "VUK") ||
				!strings.Contains(text, "Know what's left before you ship.") {
				t.Fatalf("missing wolf banner content: %q", text)
			}
			if style.color != strings.Contains(text, "\x1b[") {
				t.Fatalf("unexpected color state: %q", text)
			}
		})
	}
}
