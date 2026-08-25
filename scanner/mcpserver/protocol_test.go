package mcpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type stubService struct {
	plan       PlanResult
	scan       ScanResult
	explain    ExplainResult
	err        error
	planCalls  int
	scanCalls  int
	explainIDs []string
}

func (stub *stubService) Plan() (PlanResult, error) {
	stub.planCalls++
	return stub.plan, stub.err
}

func (stub *stubService) Scan() (ScanResult, error) {
	stub.scanCalls++
	return stub.scan, stub.err
}

func (stub *stubService) Explain(assertionID string) (ExplainResult, error) {
	stub.explainIDs = append(stub.explainIDs, assertionID)
	return stub.explain, stub.err
}

type decodedResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *rpcError       `json:"error"`
}

func serveLines(t *testing.T, service ToolService, lines ...string) []decodedResponse {
	t.Helper()
	server, err := NewServer(service, "test-version")
	if err != nil {
		t.Fatal(err)
	}
	var input, output bytes.Buffer
	for _, line := range lines {
		input.WriteString(line)
		input.WriteByte('\n')
	}
	if err := server.Serve(&input, &output); err != nil {
		t.Fatal(err)
	}
	rawLines := strings.Split(strings.TrimSuffix(output.String(), "\n"), "\n")
	if output.Len() == 0 {
		return nil
	}
	responses := make([]decodedResponse, 0, len(rawLines))
	for _, line := range rawLines {
		var response decodedResponse
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			t.Fatalf("decode response %q: %v", line, err)
		}
		responses = append(responses, response)
	}
	return responses
}

func initializeLine(version string) string {
	return `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"` + version + `","capabilities":{"roots":{"listChanged":true}},"clientInfo":{"name":"test-client","title":"Test Client","version":"1.0.0"}}}`
}

func TestLifecycleNegotiatesVersionAndListsFixedReadOnlyTools(t *testing.T) {
	responses := serveLines(t, &stubService{},
		initializeLine(ProtocolVersion20250618),
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":"tools","method":"tools/list","params":{}}`,
	)
	if len(responses) != 2 {
		t.Fatalf("responses = %d, want 2", len(responses))
	}
	var initialized initializeResult
	if err := json.Unmarshal(responses[0].Result, &initialized); err != nil {
		t.Fatal(err)
	}
	if initialized.ProtocolVersion != ProtocolVersion20250618 || initialized.Capabilities.Tools.ListChanged ||
		initialized.ServerInfo.Name != "production-readiness-scanner" || initialized.ServerInfo.Title != "Vuk" ||
		initialized.ServerInfo.Version != "test-version" {
		t.Fatalf("initialize result = %+v", initialized)
	}
	if initialized.Instructions != serverInstructions || len(initialized.Instructions) > 512 {
		t.Fatalf("instructions are not bounded: %q", initialized.Instructions)
	}
	var listed struct {
		Tools []Tool `json:"tools"`
	}
	if err := json.Unmarshal(responses[1].Result, &listed); err != nil {
		t.Fatal(err)
	}
	wantNames := []string{"prc_plan", "prc_scan", "prc_explain"}
	if len(listed.Tools) != len(wantNames) {
		t.Fatalf("tools = %+v", listed.Tools)
	}
	for index, tool := range listed.Tools {
		if tool.Name != wantNames[index] || !tool.Annotations.ReadOnlyHint || tool.Annotations.DestructiveHint ||
			!tool.Annotations.IdempotentHint || tool.Annotations.OpenWorldHint || tool.Execution.TaskSupport != "forbidden" {
			t.Fatalf("tool %d = %+v", index, tool)
		}
		if tool.InputSchema["additionalProperties"] != false || tool.OutputSchema["type"] != "object" {
			t.Fatalf("tool schemas are not closed objects: %+v", tool)
		}
	}
}

func TestUnsupportedVersionOffersLatestHandshakeRevision(t *testing.T) {
	responses := serveLines(t, &stubService{}, initializeLine("2024-11-05"))
	var initialized initializeResult
	if err := json.Unmarshal(responses[0].Result, &initialized); err != nil {
		t.Fatal(err)
	}
	if initialized.ProtocolVersion != ProtocolVersion20251125 {
		t.Fatalf("protocol version = %s", initialized.ProtocolVersion)
	}
}

func TestMalformedInitializedNotificationDoesNotOpenToolCalls(t *testing.T) {
	responses := serveLines(t, &stubService{},
		initializeLine(ProtocolVersion20251125),
		`{"jsonrpc":"2.0","method":"notifications/initialized","params":null}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`,
	)
	if len(responses) != 2 || responses[1].Error == nil || responses[1].Error.Code != -32002 {
		t.Fatalf("malformed notification opened server: %+v", responses)
	}
}

func TestToolCallsReturnMatchingStructuredAndTextContent(t *testing.T) {
	stub := &stubService{
		plan:    PlanResult{SchemaVersion: PlanResultSchema},
		scan:    ScanResult{SchemaVersion: ScanResultSchema},
		explain: ExplainResult{SchemaVersion: ExplainResultSchema},
	}
	responses := serveLines(t, stub,
		initializeLine(ProtocolVersion20251125),
		`{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"prc_plan","arguments":{}}}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"prc_scan"}}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"prc_explain","arguments":{"assertion_id":" PRC-A-CORE-001 "}}}`,
	)
	if len(responses) != 4 || stub.planCalls != 1 || stub.scanCalls != 1 ||
		len(stub.explainIDs) != 1 || stub.explainIDs[0] != "PRC-A-CORE-001" {
		t.Fatalf("responses=%d calls=%d/%d/%v", len(responses), stub.planCalls, stub.scanCalls, stub.explainIDs)
	}
	for _, response := range responses[1:] {
		var result struct {
			Content    []textContent   `json:"content"`
			Structured json.RawMessage `json:"structuredContent"`
			IsError    bool            `json:"isError"`
		}
		if err := json.Unmarshal(response.Result, &result); err != nil {
			t.Fatal(err)
		}
		if result.IsError || len(result.Content) != 1 || result.Content[0].Type != "text" ||
			result.Content[0].Text != string(result.Structured) {
			t.Fatalf("tool result does not duplicate structured JSON in text: %s", response.Result)
		}
	}
}

func TestProtocolRejectsPathInjectionUnknownToolsAndPrematureCalls(t *testing.T) {
	stub := &stubService{}
	responses := serveLines(t, stub,
		`{"jsonrpc":"2.0","id":0,"method":"tools/list"}`,
		initializeLine(ProtocolVersion20251125),
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"prc_scan","arguments":{"target":"/tmp/escape"}}}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"filesystem_write","arguments":{}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"prc_scan","arguments":null}}`,
	)
	if len(responses) != 5 || responses[0].Error == nil || responses[0].Error.Code != -32002 {
		t.Fatalf("premature response = %+v", responses)
	}
	for _, index := range []int{2, 3, 4} {
		if responses[index].Error == nil || responses[index].Error.Code != -32602 {
			t.Fatalf("response %d = %+v", index, responses[index])
		}
	}
	if stub.scanCalls != 0 {
		t.Fatalf("injected scan reached service %d times", stub.scanCalls)
	}
}

func TestServiceFailureIsAVisibleToolError(t *testing.T) {
	stub := &stubService{err: errors.New("target changed during evidence collection")}
	responses := serveLines(t, stub,
		initializeLine(ProtocolVersion20251125),
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"prc_scan","arguments":{}}}`,
	)
	var result struct {
		Content    []textContent   `json:"content"`
		Structured json.RawMessage `json:"structuredContent"`
		IsError    bool            `json:"isError"`
	}
	if err := json.Unmarshal(responses[1].Result, &result); err != nil {
		t.Fatal(err)
	}
	if !result.IsError || len(result.Structured) != 0 || len(result.Content) != 1 ||
		result.Content[0].Text != "target changed during evidence collection" {
		t.Fatalf("tool error = %+v", result)
	}
}

func TestMalformedAndDuplicateMessagesFailClosed(t *testing.T) {
	responses := serveLines(t, &stubService{},
		`{"jsonrpc":"2.0","id":1,"id":2,"method":"ping"}`,
		`[]`,
		`{"jsonrpc":"2.0","id":null,"method":"ping"}`,
		`{"jsonrpc":"2.0","id":4,"method":"ping","params":null}`,
		`{"jsonrpc":"2.0","method":"notifications/unknown"}`,
	)
	if len(responses) != 4 {
		t.Fatalf("responses = %d", len(responses))
	}
	if responses[0].Error == nil || responses[0].Error.Code != -32700 ||
		responses[1].Error == nil || responses[1].Error.Code != -32600 ||
		responses[2].Error == nil || responses[2].Error.Code != -32600 ||
		responses[3].Error == nil || responses[3].Error.Code != -32602 {
		t.Fatalf("unexpected errors: %+v", responses)
	}
}

func TestInvalidUTF8AndFractionalIDsFailClosed(t *testing.T) {
	server, err := NewServer(&stubService{}, "test-version")
	if err != nil {
		t.Fatal(err)
	}
	input := append([]byte{'{', '"', 'x', '"', ':', '"'}, 0xff)
	input = append(input, '"', '}', '\n')
	input = append(input, []byte(`{"jsonrpc":"2.0","id":1.5,"method":"ping"}`)...)
	input = append(input, '\n')
	var output bytes.Buffer
	if err := server.Serve(bytes.NewReader(input), &output); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(output.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("responses = %d: %s", len(lines), output.String())
	}
	var first, second decodedResponse
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(lines[1]), &second); err != nil {
		t.Fatal(err)
	}
	if first.Error == nil || first.Error.Code != -32700 || second.Error == nil || second.Error.Code != -32600 {
		t.Fatalf("errors = %+v / %+v", first.Error, second.Error)
	}
}

func TestOversizedInputTerminatesWithParseError(t *testing.T) {
	server, err := NewServer(&stubService{}, "test-version")
	if err != nil {
		t.Fatal(err)
	}
	input := strings.NewReader(strings.Repeat(" ", maximumInputMessage+1))
	var output bytes.Buffer
	if err := server.Serve(input, &output); err == nil {
		t.Fatal("oversized message did not terminate the transport")
	}
	var response decodedResponse
	if err := json.Unmarshal(bytes.TrimSpace(output.Bytes()), &response); err != nil {
		t.Fatal(err)
	}
	if response.Error == nil || response.Error.Code != -32700 || !strings.Contains(response.Error.Data.(string), "exceeds") {
		t.Fatalf("oversized response = %+v", response)
	}
}

func TestExactlyMaximumInputSizeIsAccepted(t *testing.T) {
	server, err := NewServer(&stubService{}, "test-version")
	if err != nil {
		t.Fatal(err)
	}
	message := `{"jsonrpc":"2.0","id":1,"method":"ping"}`
	message += strings.Repeat(" ", maximumInputMessage-len(message))
	var output bytes.Buffer
	if err := server.Serve(strings.NewReader(message), &output); err != nil {
		t.Fatal(err)
	}
	var response decodedResponse
	if err := json.Unmarshal(bytes.TrimSpace(output.Bytes()), &response); err != nil {
		t.Fatal(err)
	}
	if response.Error != nil || string(response.Result) != `{}` {
		t.Fatalf("maximum-sized response = %+v", response)
	}
}

func TestToolOutputIsBounded(t *testing.T) {
	stub := &stubService{plan: PlanResult{SchemaVersion: PlanResultSchema}}
	stub.plan.Plan.TargetName = strings.Repeat("x", maximumToolOutput)
	responses := serveLines(t, stub,
		initializeLine(ProtocolVersion20251125),
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"prc_plan","arguments":{}}}`,
	)
	var result struct {
		Content []textContent `json:"content"`
		IsError bool          `json:"isError"`
	}
	if err := json.Unmarshal(responses[1].Result, &result); err != nil {
		t.Fatal(err)
	}
	if !result.IsError || !strings.Contains(result.Content[0].Text, "exceeds") {
		t.Fatalf("unbounded result = %+v", result)
	}
}
