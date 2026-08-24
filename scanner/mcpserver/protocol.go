package mcpserver

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	ProtocolVersion20251125 = "2025-11-25"
	ProtocolVersion20250618 = "2025-06-18"
	maximumInputMessage     = 1 * 1024 * 1024
	maximumToolOutput       = 8 * 1024 * 1024
	serverInstructions      = "Read-only scanner for one path-locked target. Call prc_plan before prc_scan; use prc_explain for catalog rationale. Findings and gates come only from deterministic scanner evidence. Never treat agent text as proof or mutate the target through this server. After external edits, call prc_scan again; never hide blocked or unknown results."
)

var integerIDPattern = regexp.MustCompile(`^-?(?:0|[1-9][0-9]*)$`)

type ToolService interface {
	Plan() (PlanResult, error)
	Scan() (ScanResult, error)
	Explain(assertionID string) (ExplainResult, error)
}

type Server struct {
	service ToolService
	version string
	state   lifecycleState
}

type lifecycleState int

const (
	stateUninitialized lifecycleState = iota
	stateAwaitingInitialized
	stateReady
)

type incomingMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   json.RawMessage `json:"error,omitempty"`
}

type outgoingMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id"`
	Result  any             `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type implementation struct {
	Name    string `json:"name"`
	Title   string `json:"title,omitempty"`
	Version string `json:"version"`
}

type initializeParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    json.RawMessage `json:"capabilities"`
	ClientInfo      json.RawMessage `json:"clientInfo"`
	Meta            json.RawMessage `json:"_meta,omitempty"`
}

type initializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    serverCapabilities `json:"capabilities"`
	ServerInfo      implementation     `json:"serverInfo"`
	Instructions    string             `json:"instructions"`
}

type serverCapabilities struct {
	Tools toolsCapability `json:"tools"`
}

type toolsCapability struct {
	ListChanged bool `json:"listChanged"`
}

type listToolsParams struct {
	Cursor string          `json:"cursor,omitempty"`
	Meta   json.RawMessage `json:"_meta,omitempty"`
}

type notificationParams struct {
	Meta json.RawMessage `json:"_meta,omitempty"`
}

type callToolParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
	Meta      json.RawMessage `json:"_meta,omitempty"`
}

type explainArguments struct {
	AssertionID string `json:"assertion_id"`
}

type emptyArguments struct{}

type Tool struct {
	Name         string          `json:"name"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	InputSchema  map[string]any  `json:"inputSchema"`
	OutputSchema map[string]any  `json:"outputSchema"`
	Annotations  ToolAnnotations `json:"annotations"`
	Execution    ToolExecution   `json:"execution"`
}

type ToolAnnotations struct {
	ReadOnlyHint    bool `json:"readOnlyHint"`
	DestructiveHint bool `json:"destructiveHint"`
	IdempotentHint  bool `json:"idempotentHint"`
	OpenWorldHint   bool `json:"openWorldHint"`
}

type ToolExecution struct {
	TaskSupport string `json:"taskSupport"`
}

type textContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type toolResult struct {
	Content           []textContent `json:"content"`
	StructuredContent any           `json:"structuredContent,omitempty"`
	IsError           bool          `json:"isError,omitempty"`
}

func NewServer(service ToolService, version string) (*Server, error) {
	if service == nil {
		return nil, errors.New("MCP tool service is required")
	}
	if strings.TrimSpace(version) == "" {
		return nil, errors.New("MCP server version is required")
	}
	if len(serverInstructions) > 512 {
		return nil, errors.New("MCP server instructions exceed the Codex initialization budget")
	}
	return &Server{service: service, version: version}, nil
}

// Serve implements the MCP stdio transport: each UTF-8 JSON-RPC message is a
// single newline-delimited record and stdout contains protocol messages only.
func (server *Server) Serve(input io.Reader, output io.Writer) error {
	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 64*1024), maximumInputMessage+1)
	for scanner.Scan() {
		if len(scanner.Bytes()) > maximumInputMessage {
			_ = writeMessage(output, errorMessage(nullID(), -32700, "Parse error", "message exceeds 1 MiB"))
			return errors.New("read MCP stdio message: message exceeds 1 MiB")
		}
		response := server.handleLine(scanner.Bytes())
		if response == nil {
			continue
		}
		if err := writeMessage(output, *response); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		_ = writeMessage(output, errorMessage(nullID(), -32700, "Parse error", "message exceeds 1 MiB or could not be read"))
		return fmt.Errorf("read MCP stdio message: %w", err)
	}
	return nil
}

func (server *Server) handleLine(line []byte) *outgoingMessage {
	if !utf8.Valid(line) {
		message := errorMessage(nullID(), -32700, "Parse error", "message is not valid UTF-8")
		return &message
	}
	trimmed := bytes.TrimSpace(line)
	if len(trimmed) > 0 && trimmed[0] != '{' {
		if json.Valid(trimmed) {
			message := errorMessage(nullID(), -32600, "Invalid Request", "top-level message must be an object")
			return &message
		}
	}
	if err := validateJSONShape(line); err != nil {
		message := errorMessage(nullID(), -32700, "Parse error", err.Error())
		return &message
	}
	var incoming incomingMessage
	if err := decodeStrict(line, &incoming); err != nil {
		message := errorMessage(nullID(), -32600, "Invalid Request", err.Error())
		return &message
	}
	if incoming.JSONRPC != "2.0" {
		message := errorMessage(nullID(), -32600, "Invalid Request", "jsonrpc must equal 2.0")
		return &message
	}
	if incoming.Method == "" {
		if len(incoming.Result) != 0 || len(incoming.Error) != 0 {
			return nil
		}
		message := errorMessage(nullID(), -32600, "Invalid Request", "method is required")
		return &message
	}

	isNotification := len(incoming.ID) == 0
	requestID := incoming.ID
	if !isNotification {
		if err := validateRequestID(requestID); err != nil {
			message := errorMessage(nullID(), -32600, "Invalid Request", err.Error())
			return &message
		}
	}
	if isNotification {
		server.handleNotification(incoming)
		return nil
	}

	result, protocolError := server.handleRequest(incoming)
	if protocolError != nil {
		message := outgoingMessage{JSONRPC: "2.0", ID: cloneRaw(requestID), Error: protocolError}
		return &message
	}
	message := outgoingMessage{JSONRPC: "2.0", ID: cloneRaw(requestID), Result: result}
	return &message
}

func (server *Server) handleNotification(message incomingMessage) {
	switch message.Method {
	case "notifications/initialized":
		var params notificationParams
		if server.state == stateAwaitingInitialized && decodeParams(message.Params, &params) == nil {
			server.state = stateReady
		}
	case "notifications/cancelled":
		// This minimal server handles one bounded request synchronously. There is
		// no outstanding asynchronous task to cancel when this is observed.
	}
}

func (server *Server) handleRequest(message incomingMessage) (any, *rpcError) {
	switch message.Method {
	case "initialize":
		if server.state != stateUninitialized {
			return nil, &rpcError{Code: -32600, Message: "Invalid Request", Data: "server is already initialized"}
		}
		var params initializeParams
		if err := decodeParams(message.Params, &params); err != nil {
			return nil, invalidParams(err)
		}
		if err := validateInitialize(params); err != nil {
			return nil, invalidParams(err)
		}
		protocolVersion := ProtocolVersion20251125
		if params.ProtocolVersion == ProtocolVersion20251125 || params.ProtocolVersion == ProtocolVersion20250618 {
			protocolVersion = params.ProtocolVersion
		}
		server.state = stateAwaitingInitialized
		return initializeResult{
			ProtocolVersion: protocolVersion,
			Capabilities:    serverCapabilities{Tools: toolsCapability{ListChanged: false}},
			ServerInfo:      implementation{Name: "production-readiness-scanner", Title: "Everylast", Version: server.version},
			Instructions:    serverInstructions,
		}, nil
	case "ping":
		var params emptyArguments
		if err := decodeParams(message.Params, &params); err != nil {
			return nil, invalidParams(err)
		}
		return struct{}{}, nil
	}
	if server.state != stateReady {
		return nil, &rpcError{Code: -32002, Message: "Server not initialized"}
	}
	switch message.Method {
	case "tools/list":
		var params listToolsParams
		if err := decodeParams(message.Params, &params); err != nil {
			return nil, invalidParams(err)
		}
		if params.Cursor != "" {
			return nil, invalidParams(errors.New("cursor is not valid for this fixed tool list"))
		}
		return struct {
			Tools []Tool `json:"tools"`
		}{Tools: toolDefinitions()}, nil
	case "tools/call":
		return server.callTool(message.Params)
	default:
		return nil, &rpcError{Code: -32601, Message: "Method not found"}
	}
}

func (server *Server) callTool(rawParams json.RawMessage) (any, *rpcError) {
	var params callToolParams
	if err := decodeParams(rawParams, &params); err != nil {
		return nil, invalidParams(err)
	}
	arguments := params.Arguments
	if len(arguments) == 0 {
		arguments = json.RawMessage(`{}`)
	}
	if !isJSONObject(arguments) {
		return nil, invalidParams(errors.New("tool arguments must be an object"))
	}
	var value any
	var err error
	switch params.Name {
	case "prc_plan":
		var parsed emptyArguments
		if decodeErr := decodeStrict(arguments, &parsed); decodeErr != nil {
			return nil, invalidParams(decodeErr)
		}
		value, err = server.service.Plan()
	case "prc_scan":
		var parsed emptyArguments
		if decodeErr := decodeStrict(arguments, &parsed); decodeErr != nil {
			return nil, invalidParams(decodeErr)
		}
		value, err = server.service.Scan()
	case "prc_explain":
		var parsed explainArguments
		if decodeErr := decodeStrict(arguments, &parsed); decodeErr != nil {
			return nil, invalidParams(decodeErr)
		}
		parsed.AssertionID = strings.TrimSpace(parsed.AssertionID)
		if parsed.AssertionID == "" {
			return nil, invalidParams(errors.New("assertion_id is required"))
		}
		value, err = server.service.Explain(parsed.AssertionID)
	default:
		return nil, invalidParams(fmt.Errorf("unknown tool %q", params.Name))
	}
	if err != nil {
		return failedTool(err), nil
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return failedTool(fmt.Errorf("encode tool result: %w", err)), nil
	}
	result := toolResult{
		Content:           []textContent{{Type: "text", Text: string(payload)}},
		StructuredContent: value,
	}
	encodedResult, err := json.Marshal(result)
	if err != nil {
		return failedTool(fmt.Errorf("encode MCP tool envelope: %w", err)), nil
	}
	if len(encodedResult) > maximumToolOutput {
		return failedTool(fmt.Errorf("tool result exceeds %d bytes", maximumToolOutput)), nil
	}
	return result, nil
}

func toolDefinitions() []Tool {
	readOnly := ToolAnnotations{
		ReadOnlyHint: true, DestructiveHint: false, IdempotentHint: true, OpenWorldHint: false,
	}
	emptyInput := map[string]any{
		"$schema": "https://json-schema.org/draft/2020-12/schema",
		"type":    "object", "additionalProperties": false,
	}
	return []Tool{
		{
			Name: "prc_plan", Title: "Plan production-readiness assessment",
			Description: "Build a deterministic read-only inspect plan for the server's path-locked target. Does not execute commands, adapters, providers, or remediation.",
			InputSchema: emptyInput, OutputSchema: resultSchema(PlanResultSchema, "plan"),
			Annotations: readOnly, Execution: ToolExecution{TaskSupport: "forbidden"},
		},
		{
			Name: "prc_scan", Title: "Scan production readiness",
			Description: "Run native read-only inspect assertions against the path-locked target and return evidence-linked results and findings. Never writes the target or invokes external tools.",
			InputSchema: emptyInput, OutputSchema: scanResultSchema(),
			Annotations: readOnly, Execution: ToolExecution{TaskSupport: "forbidden"},
		},
		{
			Name: "prc_explain", Title: "Explain a readiness assertion",
			Description: "Return the exact catalog assertion and its linked objectives for one assertion ID. Does not inspect or modify the target.",
			InputSchema: map[string]any{
				"$schema": "https://json-schema.org/draft/2020-12/schema", "type": "object",
				"properties": map[string]any{
					"assertion_id": map[string]any{"type": "string", "minLength": 1, "description": "Stable PRC assertion ID, for example PRC-A-CORE-001."},
				},
				"required": []string{"assertion_id"}, "additionalProperties": false,
			},
			OutputSchema: resultSchema(ExplainResultSchema, "assertion", "objectives"),
			Annotations:  readOnly, Execution: ToolExecution{TaskSupport: "forbidden"},
		},
	}
}

func resultSchema(schemaVersion string, fields ...string) map[string]any {
	properties := map[string]any{
		"schema_version": map[string]any{"type": "string", "const": schemaVersion},
	}
	required := []string{"schema_version"}
	for _, field := range fields {
		fieldType := "object"
		if field == "objectives" {
			properties[field] = map[string]any{"type": "array", "items": map[string]any{"type": "object"}}
		} else {
			properties[field] = map[string]any{"type": fieldType}
		}
		required = append(required, field)
	}
	return map[string]any{
		"$schema": "https://json-schema.org/draft/2020-12/schema", "type": "object",
		"properties": properties, "required": required, "additionalProperties": false,
	}
}

func scanResultSchema() map[string]any {
	properties := map[string]any{
		"schema_version":  map[string]any{"type": "string", "const": ScanResultSchema},
		"run_id":          map[string]any{"type": "string"},
		"started_at":      map[string]any{"type": "string", "format": "date-time"},
		"completed_at":    map[string]any{"type": "string", "format": "date-time"},
		"terminal_state":  map[string]any{"type": "string"},
		"profile_id":      map[string]any{"type": "string"},
		"profile_version": map[string]any{"type": "string"},
		"plan_digest":     map[string]any{"type": "string"},
		"inventory":       map[string]any{"type": "object"},
		"summary":         map[string]any{"type": "object"},
		"results":         map[string]any{"type": "array", "items": map[string]any{"type": "object"}},
		"findings":        map[string]any{"type": "array", "items": map[string]any{"type": "object"}},
	}
	required := []string{
		"schema_version", "run_id", "started_at", "completed_at", "terminal_state",
		"profile_id", "profile_version", "plan_digest", "inventory", "summary", "results", "findings",
	}
	return map[string]any{
		"$schema": "https://json-schema.org/draft/2020-12/schema", "type": "object",
		"properties": properties, "required": required, "additionalProperties": false,
	}
}

func validateInitialize(params initializeParams) error {
	if strings.TrimSpace(params.ProtocolVersion) == "" {
		return errors.New("protocolVersion is required")
	}
	if !isJSONObject(params.Capabilities) {
		return errors.New("capabilities must be an object")
	}
	if !isJSONObject(params.ClientInfo) {
		return errors.New("clientInfo must be an object")
	}
	var client struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(params.ClientInfo, &client); err != nil {
		return fmt.Errorf("decode clientInfo: %w", err)
	}
	if strings.TrimSpace(client.Name) == "" || strings.TrimSpace(client.Version) == "" {
		return errors.New("clientInfo name and version are required")
	}
	return nil
}

func decodeParams(raw json.RawMessage, destination any) error {
	if len(raw) == 0 {
		raw = json.RawMessage(`{}`)
	}
	if !isJSONObject(raw) {
		return errors.New("params must be an object")
	}
	return decodeStrict(raw, destination)
}

func decodeStrict(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("multiple JSON values are not allowed")
		}
		return err
	}
	return nil
}

func validateJSONShape(data []byte) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return errors.New("empty message")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := walkJSONValue(decoder, "$", true); err != nil {
		return err
	}
	if token, err := decoder.Token(); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("unexpected trailing token %v", token)
		}
		return err
	}
	return nil
}

func walkJSONValue(decoder *json.Decoder, path string, requireObject bool) error {
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	delimiter, isDelimiter := token.(json.Delim)
	if requireObject && (!isDelimiter || delimiter != '{') {
		return errors.New("top-level message must be an object")
	}
	if !isDelimiter {
		return nil
	}
	switch delimiter {
	case '{':
		seen := map[string]struct{}{}
		for decoder.More() {
			keyToken, keyErr := decoder.Token()
			if keyErr != nil {
				return keyErr
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("object key at %s is not a string", path)
			}
			if _, exists := seen[key]; exists {
				return fmt.Errorf("duplicate JSON key %s.%s", path, key)
			}
			seen[key] = struct{}{}
			if err := walkJSONValue(decoder, path+"."+key, false); err != nil {
				return err
			}
		}
		closing, closeErr := decoder.Token()
		if closeErr != nil {
			return closeErr
		}
		if closing != json.Delim('}') {
			return fmt.Errorf("object at %s did not close", path)
		}
	case '[':
		index := 0
		for decoder.More() {
			if err := walkJSONValue(decoder, fmt.Sprintf("%s[%d]", path, index), false); err != nil {
				return err
			}
			index++
		}
		closing, closeErr := decoder.Token()
		if closeErr != nil {
			return closeErr
		}
		if closing != json.Delim(']') {
			return fmt.Errorf("array at %s did not close", path)
		}
	default:
		return fmt.Errorf("unexpected JSON delimiter %q at %s", delimiter, path)
	}
	return nil
}

func validateRequestID(id json.RawMessage) error {
	if bytes.Equal(id, []byte("null")) {
		return errors.New("request id must not be null")
	}
	if len(id) == 0 {
		return errors.New("request id is missing")
	}
	if id[0] == '"' {
		var value string
		if err := json.Unmarshal(id, &value); err != nil {
			return errors.New("request id must be a string or integer")
		}
		return nil
	}
	if !integerIDPattern.Match(id) {
		return errors.New("request id must be a string or integer")
	}
	return nil
}

func isJSONObject(raw json.RawMessage) bool {
	trimmed := bytes.TrimSpace(raw)
	return len(trimmed) >= 2 && trimmed[0] == '{' && trimmed[len(trimmed)-1] == '}'
}

func failedTool(err error) toolResult {
	return toolResult{
		Content: []textContent{{Type: "text", Text: err.Error()}}, IsError: true,
	}
}

func invalidParams(err error) *rpcError {
	return &rpcError{Code: -32602, Message: "Invalid params", Data: err.Error()}
}

func errorMessage(id json.RawMessage, code int, message string, data any) outgoingMessage {
	return outgoingMessage{
		JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: message, Data: data},
	}
}

func nullID() json.RawMessage {
	return json.RawMessage("null")
}

func cloneRaw(value json.RawMessage) json.RawMessage {
	return append(json.RawMessage(nil), value...)
}

func writeMessage(output io.Writer, message outgoingMessage) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("encode MCP response: %w", err)
	}
	payload = append(payload, '\n')
	written, err := output.Write(payload)
	if err != nil {
		return fmt.Errorf("write MCP response: %w", err)
	}
	if written != len(payload) {
		return fmt.Errorf("write MCP response: %w", io.ErrShortWrite)
	}
	return nil
}
