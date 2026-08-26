package main

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
)

const (
	ansiReset  = "\x1b[0m"
	ansiGreen  = "\x1b[32m"
	ansiRed    = "\x1b[31m"
	ansiYellow = "\x1b[33m"
	ansiBlue   = "\x1b[34m"
	ansiCyan   = "\x1b[36m"
	ansiDim    = "\x1b[2m"
)

type terminalStyle struct {
	color     bool
	hyperlink bool
}

func newTerminalStyle(mode string, output io.Writer) terminalStyle {
	if os.Getenv("TERM") == "dumb" {
		return terminalStyle{}
	}
	interactive := false
	if file, ok := output.(*os.File); ok {
		information, err := file.Stat()
		interactive = err == nil && information.Mode()&os.ModeCharDevice != 0
	}
	style := terminalStyle{hyperlink: interactive}
	if os.Getenv("NO_COLOR") != "" {
		return style
	}
	if mode == "always" {
		style.color = true
		return style
	}
	if mode != "auto" {
		return style
	}
	style.color = interactive
	return style
}

func (style terminalStyle) paint(code, value string) string {
	if !style.color {
		return value
	}
	return code + value + ansiReset
}

// fileLink returns a clickable OSC 8 link only for interactive terminals.
// The visible label is neutralized separately from the percent-encoded file URL
// so a crafted path cannot inject terminal control sequences.
func (style terminalStyle) fileLink(value string) (string, bool) {
	label := style.paint(ansiCyan, terminalText(value))
	if !style.hyperlink {
		return label, false
	}
	target, err := terminalFileURL(value)
	if err != nil {
		return label, false
	}
	return "\x1b]8;;" + target + "\x1b\\" + label + "\x1b]8;;\x1b\\", true
}

func terminalFileURL(value string) (string, error) {
	absolute, err := filepath.Abs(value)
	if err != nil {
		return "", err
	}
	path := filepath.ToSlash(absolute)
	if runtime.GOOS == "windows" {
		if strings.HasPrefix(path, "//") {
			parts := strings.SplitN(strings.TrimPrefix(path, "//"), "/", 2)
			filePath := "/"
			if len(parts) == 2 {
				filePath += parts[1]
			}
			return (&url.URL{Scheme: "file", Host: parts[0], Path: filePath}).String(), nil
		}
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
	}
	return (&url.URL{Scheme: "file", Path: path}).String(), nil
}

func terminalToneColor(tone string) string {
	switch tone {
	case "great":
		return ansiGreen
	case "good":
		return ansiCyan
	case "review":
		return ansiYellow
	default:
		return ansiRed
	}
}

func terminalSeverityColor(severity string) string {
	switch severity {
	case "critical", "high":
		return ansiRed
	case "medium":
		return ansiYellow
	case "low":
		return ansiBlue
	default:
		return ansiDim
	}
}

func localCheckBar(percentage, width int) string {
	if width < 1 {
		return ""
	}
	if percentage < 0 {
		percentage = 0
	}
	if percentage > 100 {
		percentage = 100
	}
	filled := (percentage*width + 50) / 100
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func printProductBanner(output io.Writer, style terminalStyle) {
	fmt.Fprintln(output, style.paint(ansiBlue, "  ╭────────────────────────────────────────────────────────────────────╮"))
	fmt.Fprintln(output, "  │        ╭──────╮                                                    │")
	fmt.Fprintf(output, "  │       ╱        ╲   %s                  │\n",
		style.paint(ansiCyan, "PRODUCTION READINESS CHECKLIST"))
	fmt.Fprintf(output, "  │      │    %s     │  Know what's ready and what still needs work.    │\n",
		style.paint(ansiGreen, "✓"))
	fmt.Fprintln(output, "  │       ╲        ╱                                                   │")
	fmt.Fprintln(output, "  │        ╰──────╯                                                    │")
	fmt.Fprintln(output, style.paint(ansiBlue, "  ╰────────────────────────────────────────────────────────────────────╯"))
	fmt.Fprintln(output)
}

// terminalText makes untrusted repository text one printable terminal line.
// It prevents filenames and diagnostics from injecting ANSI/OSC controls,
// cursor movement, fake lines, backspaces, or bidirectional display changes.
func terminalText(value string) string {
	var result strings.Builder
	for _, character := range value {
		if terminalUnsafe(character) {
			if character <= 0xffff {
				fmt.Fprintf(&result, "\\u%04X", character)
			} else {
				fmt.Fprintf(&result, "\\U%08X", character)
			}
			continue
		}
		result.WriteRune(character)
	}
	return result.String()
}

func terminalBrief(value string, limit int) string {
	value = terminalText(strings.Join(strings.Fields(value), " "))
	characters := []rune(value)
	if limit <= 0 || len(characters) <= limit {
		return value
	}
	return strings.TrimSpace(string(characters[:limit])) + "…"
}

func terminalUnsafe(character rune) bool {
	if unicode.IsControl(character) || character == '\u007f' || character >= '\u0080' && character <= '\u009f' {
		return true
	}
	switch {
	case character == '\u061c', character == '\u200e', character == '\u200f':
		return true
	case character >= '\u202a' && character <= '\u202e':
		return true
	case character >= '\u2066' && character <= '\u2069':
		return true
	}
	return false
}

func assessmentLabel(assessment, execution string, style terminalStyle) string {
	if execution == "error" {
		return style.paint(ansiRed, "× ERROR  ")
	}
	switch assessment {
	case "pass":
		return style.paint(ansiGreen, "✓ PASS   ")
	case "fail":
		return style.paint(ansiRed, "✗ FAIL   ")
	case "unknown":
		return style.paint(ansiYellow, "! BLOCKED")
	case "manual_review":
		return style.paint(ansiBlue, "? MANUAL ")
	case "not_applicable":
		return style.paint(ansiDim, "– N/A    ")
	default:
		return style.paint(ansiRed, "× ERROR  ")
	}
}
