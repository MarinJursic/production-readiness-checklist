package controlprogram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// DecodeProgram rejects oversized, duplicate-key, unknown-field, trailing,
// and structurally invalid JSON before returning a program.
func DecodeProgram(data []byte) (Program, error) {
	var program Program
	if err := decodeStrict(data, &program); err != nil {
		return Program{}, err
	}
	if err := ValidateProgram(program); err != nil {
		return Program{}, err
	}
	return program, nil
}

// DecodeEvidence applies the same closed JSON boundary as DecodeProgram.
func DecodeEvidence(data []byte) (Evidence, error) {
	var evidence Evidence
	if err := decodeStrict(data, &evidence); err != nil {
		return Evidence{}, err
	}
	if err := ValidateEvidence(evidence); err != nil {
		return Evidence{}, err
	}
	return evidence, nil
}

func decodeStrict(data []byte, destination any) error {
	if len(data) == 0 || len(data) > MaxDocumentBytes {
		return fmt.Errorf("document size is outside the allowed bounds")
	}
	if err := rejectDuplicateKeys(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("decode bounded document: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return fmt.Errorf("document contains trailing JSON")
	}
	return nil
}

func rejectDuplicateKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := walkJSON(decoder, 1); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return fmt.Errorf("document contains trailing JSON")
	}
	return nil
}

func walkJSON(decoder *json.Decoder, depth int) error {
	if depth > MaxJSONDepth {
		return fmt.Errorf("document exceeds JSON nesting depth")
	}
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("read JSON token: %w", err)
	}
	delimiter, composite := token.(json.Delim)
	if !composite {
		return nil
	}
	switch delimiter {
	case '{':
		seen := map[string]struct{}{}
		for decoder.More() {
			keyToken, err := decoder.Token()
			if err != nil {
				return fmt.Errorf("read JSON object key: %w", err)
			}
			key, ok := keyToken.(string)
			if !ok {
				return fmt.Errorf("JSON object key is not a string")
			}
			if _, duplicate := seen[key]; duplicate {
				return fmt.Errorf("duplicate JSON object key %q", key)
			}
			seen[key] = struct{}{}
			if err := walkJSON(decoder, depth+1); err != nil {
				return err
			}
		}
		if end, err := decoder.Token(); err != nil || end != json.Delim('}') {
			return fmt.Errorf("invalid JSON object terminator")
		}
	case '[':
		for decoder.More() {
			if err := walkJSON(decoder, depth+1); err != nil {
				return err
			}
		}
		if end, err := decoder.Token(); err != nil || end != json.Delim(']') {
			return fmt.Errorf("invalid JSON array terminator")
		}
	default:
		return fmt.Errorf("unexpected JSON delimiter")
	}
	return nil
}
