// Package httpwire contains the shared, fixed Board HTTP error surface.
package httpwire

import (
	"io"
	"net/http"
)

func Headers(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Referrer-Policy", "no-referrer")
}

func Error(w http.ResponseWriter, r *http.Request, status int) {
	code, message := "internal_error", "Board encountered an internal error."
	switch status {
	case 400:
		code, message = "invalid_request", "The request is invalid."
	case 404:
		code, message = "not_found", "The resource was not found."
	case 405:
		code, message = "method_not_allowed", "The method is not allowed."
	case 503:
		code, message = "temporarily_unavailable", "Board is temporarily unavailable."
	default:
		status = 500
	}
	Headers(w)
	w.Header().Del("X-Joy-Pi-Snapshot-Age-Ms")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if r.Method != http.MethodHead {
		_, _ = io.WriteString(w, "{\"error\":{\"code\":\""+code+"\",\"message\":\""+message+"\"}}\n")
	}
}

func InvalidRequest(r *http.Request) bool {
	return r.URL.RawQuery != "" || r.URL.ForceQuery || r.ContentLength != 0 || len(r.TransferEncoding) != 0
}
