// Package httpui owns the embedded shell's static route selection.
package httpui

import (
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/Adrien-hue/joy-pi-board/internal/httpwire"
)

const contentSecurityPolicy = "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; base-uri 'none'; object-src 'none'; frame-ancestors 'none'; form-action 'none'"

// NewHandler does not create clients, polling loops or provider state.
func NewHandler(assets fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpwire.Headers(w)

		name := ""
		if r.URL.EscapedPath() == r.URL.Path {
			if r.URL.Path == "/" {
				name = "index.html"
			} else if strings.HasPrefix(r.URL.Path, "/assets/") {
				candidate := strings.TrimPrefix(r.URL.Path, "/")
				if fs.ValidPath(candidate) && strings.Count(candidate, "/") == 1 {
					name = candidate
				}
			}
		}
		if name == "" {
			httpwire.Error(w, r, 404)
			return
		}
		info, err := fs.Stat(assets, name)
		if err != nil || !info.Mode().IsRegular() {
			httpwire.Error(w, r, 404)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			httpwire.Error(w, r, 405)
			return
		}
		if httpwire.InvalidRequest(r) {
			w.Header().Set("Connection", "close")
			httpwire.Error(w, r, 400)
			return
		}
		body, err := fs.ReadFile(assets, name)
		if err != nil {
			httpwire.Error(w, r, 500)
			return
		}
		switch path.Ext(name) {
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Content-Security-Policy", contentSecurityPolicy)
		case ".js":
			w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		}
		// Identity/no-store resources have no conditional or range semantics in
		// S02. Avoid ServeContent's independent error and caching surface.
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(http.StatusOK)
		if r.Method != http.MethodHead {
			_, _ = w.Write(body)
		}
	})
}
