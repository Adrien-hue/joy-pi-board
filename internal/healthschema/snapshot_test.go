package healthschema

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type corpusCase struct {
	ID             string          `json:"id"`
	Input          string          `json:"input_base64"`
	Classification string          `json:"classification"`
	Expected       json.RawMessage `json:"expected"`
}

func corpus(t testing.TB) []corpusCase {
	t.Helper()
	body, err := os.ReadFile("../../testdata/s01/cases.json")
	if err != nil {
		t.Fatal(err)
	}
	var doc struct {
		Cases []corpusCase `json:"cases"`
	}
	if err := json.Unmarshal(body, &doc); err != nil {
		t.Fatal(err)
	}
	if len(doc.Cases) < 200 {
		t.Fatal("missing contract corpus")
	}
	return doc.Cases
}
func input(t testing.TB, c corpusCase) []byte {
	t.Helper()
	b, e := base64.StdEncoding.DecodeString(c.Input)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func classification(err error) string {
	if err == nil {
		return "valid"
	}
	var e *Error
	if errors.As(err, &e) {
		return string(e.Kind)
	}
	return "unclassified"
}

// Test-only projection with tagged decimal strings. No transport uses this
// representation: it makes the shared expectation metadata lossless in both tests.
func normalized(v reflect.Value) any {
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		return normalized(v.Elem())
	}
	switch v.Kind() {
	case reflect.Uint64:
		return "u64:" + strconv.FormatUint(v.Uint(), 10)
	case reflect.Struct:
		m := map[string]any{}
		for i := 0; i < v.NumField(); i++ {
			m[v.Type().Field(i).Tag.Get("json")] = normalized(v.Field(i))
		}
		return m
	case reflect.Slice:
		a := make([]any, v.Len())
		for i := range a {
			a[i] = normalized(v.Index(i))
		}
		return a
	default:
		return v.Interface()
	}
}
func normalizedJSON(t testing.TB, v Validated) json.RawMessage {
	t.Helper()
	s, err := v.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(normalized(reflect.ValueOf(s)))
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func equalJSON(t testing.TB, got, want []byte) {
	t.Helper()
	var a, b any
	if err := json.Unmarshal(got, &a); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatalf("known-field projection differs\ngot %s\nwant %s", got, want)
	}
}

// B-01/B-02/B-03/B-04: same input bytes and expected outcomes as TypeScript.
func TestSharedCorpus(t *testing.T) {
	for _, c := range corpus(t) {
		t.Run(c.ID, func(t *testing.T) {
			v, err := Decode(input(t, c))
			if classification(err) != c.Classification {
				t.Fatalf("got %v, want %s", err, c.Classification)
			}
			if err == nil {
				equalJSON(t, normalizedJSON(t, v), c.Expected)
			}
		})
	}
}

func TestRequiredFieldPresence(t *testing.T) {
	for _, c := range corpus(t) {
		if !strings.HasPrefix(c.ID, "missing-") && !strings.HasPrefix(c.ID, "issue-missing-") {
			continue
		}
		t.Run(c.ID, func(t *testing.T) {
			_, err := Decode(input(t, c))
			var diagnostic *Error
			if !errors.As(err, &diagnostic) || diagnostic.Rule != "required field absent" {
				t.Fatalf("missing required member must be tracked separately from null: %v", err)
			}
		})
	}
}
func TestImmutableAndObjectTransport(t *testing.T) {
	c := corpus(t)[0]
	original := input(t, c)
	v, err := Decode(original)
	if err != nil {
		t.Fatal(err)
	}
	want := bytes.Clone(original)
	original[0] = '!'
	returned := v.Bytes()
	returned[0] = '!'
	marshaled, err := v.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	marshaled[0] = '!'
	s, err := v.Snapshot()
	if err != nil {
		t.Fatal(err)
	}
	s.Host.Hostname = "changed"
	*s.UptimeSeconds = 0
	s.Network.Interfaces[0].Name = "changed"
	*s.Network.Interfaces[0].RXBytes = 1
	if !bytes.Equal(v.Bytes(), want) {
		t.Fatal("mutable backing bytes escaped")
	}
	equalJSON(t, normalizedJSON(t, v), c.Expected)
	envelope, err := json.Marshal(struct {
		Snapshot Validated `json:"snapshot"`
	}{v})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(envelope, []byte(`{"snapshot":{`)) {
		t.Fatal("snapshot must be an object")
	}
	for _, token := range [][]byte{[]byte("9007199254740993"), []byte("18446744073709551615")} {
		if !bytes.Contains(envelope, token) {
			t.Fatal("exact integer lost")
		}
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				b, _ := v.MarshalJSON()
				b[0] = '!'
				if _, err := v.Snapshot(); err != nil {
					t.Error(err)
				}
			}
		}()
	}
	wg.Wait()
}
func TestPreservesUnknownAndOriginalTokens(t *testing.T) {
	for _, c := range corpus(t) {
		if c.Classification != "valid" {
			continue
		}
		t.Run(c.ID, func(t *testing.T) {
			b := input(t, c)
			v, err := Decode(b)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b, v.Bytes()) {
				t.Fatal("original bytes modified")
			}
			serialized, err := json.Marshal(v)
			if err != nil {
				t.Fatal(err)
			}
			var compact bytes.Buffer
			if err := json.Compact(&compact, b); err != nil {
				t.Fatal(err)
			}
			// json.Marshal may escape HTML, but never changes numeric tokens.
			var escaped bytes.Buffer
			json.HTMLEscape(&escaped, compact.Bytes())
			if !bytes.Equal(serialized, escaped.Bytes()) {
				t.Fatal("transport changed original tokens or unknown members")
			}
		})
	}
}
func TestZeroValueAndDiagnostics(t *testing.T) {
	var v Validated
	if v.Valid() {
		t.Fatal("zero value is not valid")
	}
	if _, err := json.Marshal(v); err == nil {
		t.Fatal("uninitialized marshal succeeded")
	}
	if _, err := v.Snapshot(); err == nil {
		t.Fatal("uninitialized projection succeeded")
	}
	_, err := Decode([]byte(`{"schema_version":"secret-internal-value","private":"payload-secret"}`))
	if err == nil || bytes.Contains([]byte(err.Error()), []byte("secret")) {
		t.Fatal("diagnostic echoes input")
	}
}
func FuzzDecode(f *testing.F) {
	for _, c := range corpus(f) {
		f.Add(input(f, c))
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		v, err := Decode(b)
		if err != nil {
			if classification(err) == "unclassified" {
				t.Fatal(err)
			}
			return
		}
		if !bytes.Equal(b, v.Bytes()) {
			t.Fatal("lost bytes")
		}
		out := compactObject(t, v)
		again, err := Decode(out)
		if err != nil {
			t.Fatal(err)
		}
		equalJSON(t, normalizedJSON(t, again), normalizedJSON(t, v))
	})
}

// A bounded round trip must not mistake optional HTML escaping expansion (or
// Encoder's framing newline) for an invalid original snapshot. S02 owns the
// eventual HTTP/envelope size and encoder configuration.
func compactObject(t testing.TB, v Validated) []byte {
	t.Helper()
	var out bytes.Buffer
	encoder := json.NewEncoder(&out)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		t.Fatal(err)
	}
	return bytes.TrimSuffix(out.Bytes(), []byte("\n"))
}

func TestRoundtripAtLimitWithMarkup(t *testing.T) {
	b := input(t, corpus(t)[0])
	prefix := string(b[:len(b)-1]) + `,"extra":"`
	raw := []byte(prefix + strings.Repeat("<", MaximumBytes-len(prefix)-2) + `"}`)
	v, err := Decode(raw)
	if err != nil {
		t.Fatal(err)
	}
	defaultEncoding, err := json.Marshal(v)
	if err != nil || !json.Valid(defaultEncoding) || len(defaultEncoding) <= MaximumBytes {
		t.Fatal("expected valid object with optional HTML escaping expansion")
	}
	out := compactObject(t, v)
	if len(out) != MaximumBytes {
		t.Fatal("unexpected compact object size")
	}
	if _, err := Decode(out); err != nil {
		t.Fatal(err)
	}
}

// Explicit test-only handoff to TypeScript, invoked by the contract task.
// No application command/endpoint is added. Output files live in ignored out/.
func TestExportInterop(t *testing.T) {
	dir := os.Getenv("S01_INTEROP_DIR")
	if dir == "" {
		t.Skip("run npm --prefix web run contract for Go-to-TypeScript handoff")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	type result struct {
		ID             string          `json:"id"`
		Classification string          `json:"classification"`
		File           string          `json:"file,omitempty"`
		Expected       json.RawMessage `json:"expected,omitempty"`
	}
	var results []result
	for i, c := range corpus(t) {
		v, err := Decode(input(t, c))
		row := result{ID: c.ID, Classification: classification(err)}
		if err == nil {
			row.File = fmt.Sprintf("%03d", i)
			row.Expected = normalizedJSON(t, v)
			raw, e := v.MarshalJSON()
			if e != nil {
				t.Fatal(e)
			}
			if e = os.WriteFile(filepath.Join(dir, row.File+".snapshot.json"), raw, 0644); e != nil {
				t.Fatal(e)
			}
			envelope, e := json.Marshal(struct {
				Snapshot Validated `json:"snapshot"`
			}{v})
			if e != nil {
				t.Fatal(e)
			}
			if e = os.WriteFile(filepath.Join(dir, row.File+".envelope.json"), envelope, 0644); e != nil {
				t.Fatal(e)
			}
		}
		results = append(results, row)
	}
	b, err := json.Marshal(results)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "results.json"), b, 0644); err != nil {
		t.Fatal(err)
	}
	// Seeded adversarial differential cases supplement the reviewed common oracle.
	rng := rand.New(rand.NewSource(20260920))
	seed := input(t, corpus(t)[0])
	var mutations []corpusCase
	for i := 0; i < 512; i++ {
		b := bytes.Clone(seed)
		for j := 0; j < 1+rng.Intn(4); j++ {
			b[rng.Intn(len(b))] = byte(rng.Intn(256))
		}
		v, e := Decode(b)
		c := corpusCase{ID: fmt.Sprintf("mutation-%03d", i), Input: base64.StdEncoding.EncodeToString(b), Classification: classification(e)}
		if e == nil {
			c.Expected = normalizedJSON(t, v)
		}
		mutations = append(mutations, c)
	}
	b, err = json.Marshal(mutations)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "mutations.json"), b, 0644); err != nil {
		t.Fatal(err)
	}
}
