package httpui

import (
	"bytes"
	"compress/gzip"
	"github.com/Adrien-hue/joy-pi-board/web"
	"io"
	"io/fs"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
)

// S04-T02: q semantics, public inventory, precedence, representation headers, HEAD and errors.
func TestCompression(t *testing.T) {
	assets := web.Assets()
	h := NewHandler(assets)
	original, _ := fs.ReadFile(assets, "index.html")
	for _, c := range []struct {
		ae       string
		status   int
		encoding string
	}{
		{"", 200, ""}, {"identity", 200, ""}, {"gzip", 200, "gzip"}, {"GZip;q=1, identity;q=0.4", 200, "gzip"}, {"gzip;q=0.5", 200, ""},
		{"gzip;q=0.5,identity;q=0.2", 200, "gzip"}, {"*;q=0,gzip;q=1", 200, "gzip"}, {"*;q=1,gzip;q=0", 200, ""},
		{"gzip;q=1,GZIP;q=0,identity;q=0", 406, ""}, {"br,identity;q=0", 406, ""}, {"*;q=0", 406, ""},
		{"gzip;q=1.001", 400, ""}, {"gzip;q=-1", 400, ""}, {"gzip;q=.5", 400, ""}, {"gzip;q=0.1234", 400, ""}, {"gzip;foo=1", 400, ""}, {"gzip;q=1;q=0", 400, ""}, {"g zip", 400, ""},
	} {
		for _, method := range []string{"GET", "HEAD"} {
			t.Run(method+c.ae, func(t *testing.T) {
				r := httptest.NewRequest(method, "/", nil)
				r.Header.Set("Accept-Encoding", c.ae)
				r.Header.Set("Range", "bytes=0-1")
				r.Header.Set("If-None-Match", "*")
				w := httptest.NewRecorder()
				h.ServeHTTP(w, r)
				if w.Code != c.status || w.Header().Get("Content-Encoding") != c.encoding || w.Header().Get("Vary") != "Accept-Encoding" {
					t.Fatal(w.Code, w.Header())
				}
				if method == "HEAD" && w.Body.Len() != 0 {
					t.Fatal("HEAD body")
				}
				if c.status == 200 {
					if w.Header().Get("Content-Type") != "text/html; charset=utf-8" || w.Header().Get("Cache-Control") != "no-store" {
						t.Fatal(w.Header())
					}
					if method == "GET" {
						if w.Header().Get("Content-Length") != strconv.Itoa(w.Body.Len()) {
							t.Fatal("length")
						}
						body := w.Body.Bytes()
						if c.encoding == "gzip" {
							z, e := gzip.NewReader(bytes.NewReader(body))
							if e != nil {
								t.Fatal(e)
							}
							body, e = io.ReadAll(z)
							if e != nil {
								t.Fatal(e)
							}
							z.Close()
						}
						if !bytes.Equal(body, original) {
							t.Fatal("representation differs")
						}
					}
				} else if w.Header().Get("Content-Security-Policy") != "" {
					t.Fatal("stale CSP")
				}
			})
		}
	}
	inv, _ := Inventory(assets)
	for _, path := range []string{"/manifest.json", "/index.html.gz", "/assets/missing.js", "/assets/../index.html"} {
		r := httptest.NewRequest("POST", path+"?x", nil)
		r.Header.Set("Accept-Encoding", "*;q=0")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 404 || w.Body.Len() != 0 {
			t.Fatal(path, w.Code)
		}
	}
	for _, asset := range inv {
		r := httptest.NewRequest("GET", "/"+asset.Gzip.Path, nil)
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 404 {
			t.Fatal("sidecar public")
		}
	}
	for _, c := range []struct {
		method, path, body string
		status             int
	}{{"POST", "/", "", 405}, {"GET", "/?", "", 400}, {"GET", "/", "x", 400}} {
		r := httptest.NewRequest(c.method, c.path, strings.NewReader(c.body))
		r.Header.Set("Accept-Encoding", "*;q=0")
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != c.status {
			t.Fatal(w.Code)
		}
	}
	// Independent header lines combine, duplicate minimum wins.
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Accept-Encoding", "gzip")
	r.Header.Add("Accept-Encoding", "gzip;q=0")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Header().Get("Content-Encoding") != "" {
		t.Fatal("duplicate")
	}
	broken := fstest.MapFS{}
	_ = fs.WalkDir(assets, ".", func(p string, d fs.DirEntry, e error) error {
		if e == nil && !d.IsDir() {
			b, _ := fs.ReadFile(assets, p)
			broken[p] = &fstest.MapFile{Data: b}
		}
		return e
	})
	broken["index.html"].Data = []byte("corrupt")
	if _, e := Inventory(broken); e == nil {
		t.Fatal("hash mismatch accepted")
	}
	w = httptest.NewRecorder()
	NewHandler(broken).ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 500 {
		t.Fatal(w.Code)
	}
}
