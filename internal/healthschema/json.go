package healthschema

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"unicode/utf8"
)

// Strict Unicode guard only: encoding/json owns the JSON grammar and tokenizing.
// Go otherwise replaces malformed UTF-8 and unpaired escaped surrogates.
func unicodeGuard(b []byte) error {
	if !utf8.Valid(b) {
		return invalid("/", "malformed UTF-8")
	}
	inString := false
	for i := 0; i < len(b); i++ {
		if b[i] == '"' {
			inString = !inString
			continue
		}
		if !inString || b[i] != '\\' {
			continue
		}
		i++
		if i >= len(b) {
			return invalid("/", "truncated escape")
		}
		if b[i] != 'u' {
			continue
		}
		if i+4 >= len(b) {
			return invalid("/", "truncated Unicode escape")
		}
		n, err := strconv.ParseUint(string(b[i+1:i+5]), 16, 16)
		if err != nil {
			return invalid("/", "invalid Unicode escape")
		}
		i += 4
		if n >= 0xdc00 && n <= 0xdfff {
			return invalid("/", "unpaired low surrogate")
		}
		if n < 0xd800 || n > 0xdbff {
			continue
		}
		if i+6 >= len(b) || b[i+1] != '\\' || b[i+2] != 'u' {
			return invalid("/", "unpaired high surrogate")
		}
		low, err := strconv.ParseUint(string(b[i+3:i+7]), 16, 16)
		if err != nil || low < 0xdc00 || low > 0xdfff {
			return invalid("/", "unpaired high surrogate")
		}
		i += 6
	}
	return nil
}

func document(b []byte) (map[string]any, error) {
	if err := unicodeGuard(b); err != nil {
		return nil, err
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	v, err := readValue(d, 0)
	if err != nil {
		return nil, err
	}
	if _, err := d.Token(); err != io.EOF {
		return nil, invalid("/", "trailing content")
	}
	root, ok := v.(map[string]any)
	if !ok {
		return nil, invalid("/", "root must be an object")
	}
	return root, nil
}

// Walk standard-library tokens, rejecting duplicate decoded names at every
// object (including unknown members). Root object counts as container 1.
func readValue(d *json.Decoder, depth int) (any, error) {
	t, err := d.Token()
	if err != nil {
		return nil, invalid("/", "malformed JSON")
	}
	delim, container := t.(json.Delim)
	if !container {
		return t, nil
	}
	if depth >= MaximumDepth {
		return nil, invalid("/", "depth exceeds 32 containers")
	}
	switch delim {
	case '{':
		obj := make(map[string]any)
		for d.More() {
			keyToken, err := d.Token()
			if err != nil {
				return nil, invalid("/", "malformed object key")
			}
			key, ok := keyToken.(string)
			if !ok {
				return nil, invalid("/", "object key must be a string")
			}
			if _, exists := obj[key]; exists {
				return nil, invalid("/", "duplicate object key")
			}
			value, err := readValue(d, depth+1)
			if err != nil {
				return nil, err
			}
			obj[key] = value
		}
		if end, err := d.Token(); err != nil || end != json.Delim('}') {
			return nil, invalid("/", "unclosed object")
		}
		return obj, nil
	case '[':
		array := make([]any, 0)
		for d.More() {
			value, err := readValue(d, depth+1)
			if err != nil {
				return nil, err
			}
			array = append(array, value)
		}
		if end, err := d.Token(); err != nil || end != json.Delim(']') {
			return nil, invalid("/", "unclosed array")
		}
		return array, nil
	default:
		return nil, invalid("/", "unexpected delimiter")
	}
}
