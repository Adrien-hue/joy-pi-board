package httpui_test

import (
	"bytes"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Adrien-hue/joy-pi-board/internal/httpui"
	"github.com/Adrien-hue/joy-pi-board/web"
)

func TestCompiledAssetsAreServedExactly(t *testing.T) {
	assets := web.Assets()
	handler := httpui.NewHandler(assets)
	var assetCount int
	err := fs.WalkDir(assets, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}
		assetCount++
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			target := "/" + name
			if name == "index.html" {
				target = "/"
			}
			want, err := fs.ReadFile(assets, name)
			if err != nil {
				t.Fatal(err)
			}
			for _, method := range []string{http.MethodGet, http.MethodHead} {
				response := httptest.NewRecorder()
				handler.ServeHTTP(response, httptest.NewRequest(method, target, nil))
				if response.Code != http.StatusOK {
					t.Fatalf("%s %s = %d", method, target, response.Code)
				}
				if method == http.MethodGet && !bytes.Equal(response.Body.Bytes(), want) {
					t.Error("served body differs from compiled embedded bytes")
				}
				if method == http.MethodHead && response.Body.Len() != 0 {
					t.Error("HEAD response has a body")
				}
				if response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("X-Content-Type-Options") != "nosniff" || response.Header().Get("Referrer-Policy") != "no-referrer" {
					t.Error("missing shell cache/security headers")
				}
				if response.Header().Get("Content-Type") == "" {
					t.Error("missing content type")
				}
				if name == "index.html" && !strings.Contains(response.Header().Get("Content-Security-Policy"), "script-src 'self'") {
					t.Error("missing HTML CSP")
				}
			}
		})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if assetCount < 3 {
		t.Fatalf("expected compiled HTML, JS and CSS; got %d assets", assetCount)
	}
}

func TestShellRoutesDoNotPretendToImplementOverview(t *testing.T) {
	t.Parallel()
	handler := httpui.NewHandler(web.Assets())
	for _, tc := range []struct {
		method, path, body string
		status             int
	}{
		{http.MethodGet, "/api/v1/overview", "", 404},
		{http.MethodGet, "/other", "", 404},
		{http.MethodGet, "/index.html", "", 404},
		{http.MethodGet, "/assets/", "", 404},
		{http.MethodGet, "/assets/../index.html", "", 404},
		{http.MethodGet, "/assets/%2e%2e/index.html", "", 404},
		{http.MethodGet, "/assets/missing.js", "", 404},
		{http.MethodPost, "/", "", 405},
		{http.MethodGet, "/?x=1", "", 400},
		{http.MethodGet, "/", "body", 400},
	} {
		t.Run(tc.method+tc.path+tc.body, func(t *testing.T) {
			t.Parallel()
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body)))
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}
			if strings.Contains(response.Body.String(), "<html") {
				t.Error("invalid route returned a misleading SPA fallback")
			}
		})
	}
}
