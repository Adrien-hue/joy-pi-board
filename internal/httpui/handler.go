// Package httpui serves only manifest-listed public resources from embed.
package httpui

import (
	"github.com/Adrien-hue/joy-pi-board/internal/httpwire"
	"io/fs"
	"net/http"
	"strconv"
)

const contentSecurityPolicy = "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; base-uri 'none'; object-src 'none'; frame-ancestors 'none'; form-action 'none'"

func NewHandler(assets fs.FS) http.Handler {
	inventory, inventoryErr := Inventory(assets)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpwire.Headers(w)
		w.Header().Set("Vary", "Accept-Encoding")
		if inventoryErr != nil {
			Error(w, r, 500)
			return
		}
		resource, ok := inventory[r.URL.Path]
		if !ok || r.URL.EscapedPath() != r.URL.Path {
			Error(w, r, 404)
			return
		}
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			Error(w, r, 405)
			return
		}
		if httpwire.InvalidRequest(r) {
			w.Header().Set("Connection", "close")
			Error(w, r, 400)
			return
		}
		identity, gzip, err := encodings(r.Header.Values("Accept-Encoding"))
		if err != nil {
			Error(w, r, 400)
			return
		}
		if identity == 0 && gzip == 0 {
			Error(w, r, 406)
			return
		}
		selected := resource.Identity
		if gzip > 0 && gzip >= identity {
			selected = resource.Gzip
			w.Header().Set("Content-Encoding", "gzip")
		}
		body, err := fs.ReadFile(assets, selected.Path)
		if err != nil {
			Error(w, r, 500)
			return
		}
		w.Header().Set("Content-Type", resource.MIME)
		if resource.URL == "/" {
			w.Header().Set("Content-Security-Policy", contentSecurityPolicy)
		}
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.WriteHeader(200)
		if r.Method != http.MethodHead {
			_, _ = w.Write(body)
		}
	})
}
