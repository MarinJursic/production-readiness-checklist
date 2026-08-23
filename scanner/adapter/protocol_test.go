package adapter

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixturePath(name string) string {
	return filepath.Join("..", "..", "fixtures", "adapters", name)
}

func TestInputJSONLUsesThreeTypedMessages(t *testing.T) {
	data, err := InputJSONL(strings.Repeat("a", 64), Subject{
		TargetName: "example", InventoryDigest: strings.Repeat("b", 64),
	}, map[string]any{"source_files": 2}, map[string]any{})
	if err != nil {
		t.Fatal(err)
	}
	lines := bytes.Split(bytes.TrimSpace(data), []byte("\n"))
	if len(lines) != 3 {
		t.Fatalf("expected three messages, got %d", len(lines))
	}
	if !bytes.Contains(lines[0], []byte(`"protocol":"prc-adapter-jsonl-v1"`)) ||
		!bytes.Contains(lines[2], []byte(`"type":"execute"`)) {
		t.Fatalf("unexpected input protocol: %s", data)
	}
}

func TestParseOutputAcceptsValidTranscript(t *testing.T) {
	input := strings.Join([]string{
		`{"type":"log","level":"info","message":"started"}`,
		`{"type":"observation","observation":{"id":"OBS-1","kind":"secret-pattern","outcome":"not_found","summary":"No match in the authorized scope.","locations":[]}}`,
		`{"type":"artifact","artifact":{"id":"ART-1","media_type":"application/json","digest":"sha256:` + strings.Repeat("a", 64) + `","size":2,"path":"artifacts/result.json"}}`,
		`{"type":"summary","status":"completed","counts":{"observations":1}}`,
	}, "\n") + "\n"
	transcript, err := ParseOutput(strings.NewReader(input), DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	if len(transcript.Observations) != 1 || len(transcript.Artifacts) != 1 || transcript.Summary.Status != "completed" {
		t.Fatalf("unexpected transcript: %+v", transcript)
	}
}

func TestCheckedInProtocolFixtures(t *testing.T) {
	valid, err := os.Open(fixturePath("valid-output.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	defer valid.Close()
	if _, err := ParseOutput(valid, DefaultLimits()); err != nil {
		t.Fatalf("valid fixture: %v", err)
	}
	malicious, err := os.Open(fixturePath("malicious-authority-output.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	defer malicious.Close()
	if _, err := ParseOutput(malicious, DefaultLimits()); err == nil {
		t.Fatal("malicious authority fixture was accepted")
	}
}

func TestParseOutputRejectsAuthorityAndAmbiguityAttacks(t *testing.T) {
	tests := map[string]string{
		"adapter assessment": `{"type":"observation","observation":{"id":"OBS-1","kind":"x","outcome":"found","summary":"x","locations":[],"assessment":"pass"}}` + "\n" + `{"type":"summary","status":"completed","counts":{}}`,
		"duplicate key":      `{"type":"summary","type":"log","status":"completed","counts":{}}`,
		"path traversal":     `{"type":"artifact","artifact":{"id":"ART-1","media_type":"text/plain","digest":"sha256:` + strings.Repeat("a", 64) + `","size":1,"path":"../secret"}}` + "\n" + `{"type":"summary","status":"completed","counts":{}}`,
		"after summary":      `{"type":"summary","status":"completed","counts":{}}` + "\n" + `{"type":"log","level":"info","message":"late"}`,
		"missing summary":    `{"type":"log","level":"info","message":"only"}`,
		"partial no reason":  `{"type":"summary","status":"partial","counts":{}}`,
		"null counts":        `{"type":"summary","status":"completed","counts":null}`,
		"null locations":     `{"type":"observation","observation":{"id":"OBS-1","kind":"x","outcome":"found","summary":"x","locations":null}}` + "\n" + `{"type":"summary","status":"completed","counts":{}}`,
	}
	for name, input := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := ParseOutput(strings.NewReader(input+"\n"), DefaultLimits()); err == nil {
				t.Fatal("expected protocol rejection")
			}
		})
	}
}

func TestParseOutputEnforcesAllLimits(t *testing.T) {
	base := DefaultLimits()
	base.MaxStdout = 20
	if _, err := ParseOutput(strings.NewReader(strings.Repeat("x", 21)), base); err == nil {
		t.Fatal("expected stdout byte limit")
	}

	base = DefaultLimits()
	base.MaxLineBytes = 16
	if _, err := ParseOutput(strings.NewReader(strings.Repeat("x", 17)+"\n"), base); err == nil {
		t.Fatal("expected line limit")
	}

	base = DefaultLimits()
	base.MaxMessages = 1
	input := `{"type":"log","level":"info","message":"one"}` + "\n" +
		`{"type":"summary","status":"completed","counts":{}}` + "\n"
	if _, err := ParseOutput(strings.NewReader(input), base); err == nil {
		t.Fatal("expected message limit")
	}
}

func FuzzParseOutputNeverPanics(f *testing.F) {
	f.Add([]byte(`{"type":"summary","status":"completed","counts":{}}` + "\n"))
	f.Add([]byte(`{"type":"summary","type":"log"}` + "\n"))
	f.Add([]byte{0xff, 0x00, '\n'})
	f.Fuzz(func(t *testing.T, input []byte) {
		limits := DefaultLimits()
		limits.MaxStdout = 1024 * 1024
		_, _ = ParseOutput(bytes.NewReader(input), limits)
	})
}
