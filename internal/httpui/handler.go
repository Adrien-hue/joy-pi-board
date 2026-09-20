// Package httpui serves only the embedded S00 shell and its static assets.
package httpui

import (
	"bytes"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"time"
)

const contentSecurityPolicy = "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; base-uri 'none'; object-src 'none'; frame-ancestors 'none'; form-action 'none'"

// NewHandler does not create clients, polling loops or provider state.
func NewHandler(assets fs.FS) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cache-Control", "no-store")

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
			http.NotFound(w, r)
			return
		}
		info, err := fs.Stat(assets, name)
		if err != nil || !info.Mode().IsRegular() {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.RawQuery != "" || r.ContentLength != 0 || len(r.TransferEncoding) != 0 {
			w.Header().Set("Connection", "close")
			http.Error(w, "Invalid request.", http.StatusBadRequest)
			return
		}
		body, err := fs.ReadFile(assets, name)
		if err != nil {
			http.Error(w, "Unable to serve the interface.", http.StatusInternalServerError)
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
		http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(body))
	})
}
