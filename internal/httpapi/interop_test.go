package httpapi

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/Adrien-hue/joy-pi-board/internal/health"
	"github.com/Adrien-hue/joy-pi-board/internal/overview"
	"github.com/Adrien-hue/joy-pi-board/web"
)

// S02-T09: Node executes while both real Go HTTP servers are running. The
// selection server and manifest exist only in this test, never in the binary.
func TestHTTPTypeScriptInterop(t *testing.T) {
	node := os.Getenv("S02_NODE")
	if node == "" {
		t.Skip("executed by tasks.mjs contract after dedicated parser build")
	}
	data, e := os.ReadFile("../../testdata/s01/cases.json")
	if e != nil {
		t.Fatal(e)
	}
	type row struct {
		ID             string `json:"id"`
		Input          string `json:"input_base64"`
		Classification string `json:"classification"`
	}
	var corpus struct {
		Cases []row `json:"cases"`
	}
	if e = json.Unmarshal(data, &corpus); e != nil {
		t.Fatal(e)
	}
	base := full(t).Bytes()
	var compact bytes.Buffer
	if e = json.Compact(&compact, base); e != nil {
		t.Fatal(e)
	}
	base = compact.Bytes()
	add := func(id string, body []byte) {
		corpus.Cases = append(corpus.Cases, row{id, base64.StdEncoding.EncodeToString(body), "valid"})
	}
	markup := []byte(`,"unknown":"<>&` + "\u2028\u2029\u00e9\U0001f600" + `"}`)
	initial := append(append([]byte{}, base[:len(base)-1]...), markup...)
	pad := 65536 - len(initial)
	markup = append(append([]byte{}, initial[:len(initial)-2]...), []byte(strings.Repeat("&", pad)+`"}`)...)
	add("s02-maximum-markup-unicode", markup)
	deep := append(append([]byte{}, base[:len(base)-1]...), []byte(`,"unknown":`+strings.Repeat("[", 31)+"0"+strings.Repeat("]", 31)+"}")...)
	add("s02-depth-32", deep)
	var index atomic.Int32
	var calls atomic.Int32
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		body, e := base64.StdEncoding.DecodeString(corpus.Cases[index.Load()].Input)
		if e != nil {
			t.Error(e)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	defer upstream.Close()
	client := health.New(upstream.URL + "/v1/snapshot")
	defer client.Close()
	coordinator := overview.New(client, nil, nil)
	defer coordinator.Stop()
	board := httptest.NewServer(New(web.Assets(), coordinator))
	defer board.Close()
	control := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/cases" {
			_ = json.NewEncoder(w).Encode(corpus.Cases)
			return
		}
		i, e := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/"))
		if e != nil || i < 0 || i >= len(corpus.Cases) {
			w.WriteHeader(400)
			return
		}
		index.Store(int32(i))
		w.WriteHeader(204)
	}))
	defer control.Close()
	root, e := filepath.Abs("../..")
	if e != nil {
		t.Fatal(e)
	}
	command := exec.Command(node, filepath.Join(root, "scripts/http-contract.mjs"), board.URL, control.URL)
	command.Dir = root
	out, e := command.CombinedOutput()
	t.Log(string(out))
	if e != nil {
		t.Fatal(e)
	}
	if calls.Load() != int32(len(corpus.Cases)) {
		t.Fatal("attempt count", calls.Load(), len(corpus.Cases))
	}
	t.Log(fmt.Sprintf("%d actual Health HTTP attempts, one per overview; no retries", calls.Load()))
}
