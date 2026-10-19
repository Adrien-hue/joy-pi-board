package httpui

import (
	"errors"
	"github.com/Adrien-hue/joy-pi-board/internal/httpwire"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var codingToken = regexp.MustCompile("^[!#$%&'*+.^_`|~0-9A-Za-z-]+$")
var quality = regexp.MustCompile(`^(0(?:\.[0-9]{0,3})?|1(?:\.0{0,3})?)$`)

func encodings(values []string) (int, int, error) {
	joined := strings.Join(values, ",")
	if strings.TrimSpace(joined) == "" {
		return 1000, 0, nil
	}
	qs := map[string]int{}
	for _, item := range strings.Split(joined, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		parts := strings.Split(item, ";")
		token := strings.ToLower(strings.TrimSpace(parts[0]))
		q := 1000
		if !codingToken.MatchString(token) || len(parts) > 2 {
			return 0, 0, errors.New("invalid encoding")
		}
		if len(parts) == 2 {
			k, v, ok := strings.Cut(strings.TrimSpace(parts[1]), "=")
			if !ok || !strings.EqualFold(strings.TrimSpace(k), "q") || !quality.MatchString(strings.TrimSpace(v)) {
				return 0, 0, errors.New("invalid quality")
			}
			n, _ := strconv.ParseFloat(strings.TrimSpace(v), 64)
			q = int(n*1000 + 0.5)
		}
		if old, ok := qs[token]; !ok || q < old {
			qs[token] = q
		}
	}
	identity := 1000
	if q, ok := qs["identity"]; ok {
		identity = q
	} else if q, ok := qs["*"]; ok && q == 0 {
		identity = 0
	}
	gzip := qs["*"]
	if q, ok := qs["gzip"]; ok {
		gzip = q
	}
	return identity, gzip, nil
}

// Error preserves existing error precedence without forcing an excluded identity body.
func Error(w http.ResponseWriter, r *http.Request, status int) {
	w.Header().Del("Content-Encoding")
	w.Header().Del("Content-Length")
	w.Header().Del("Content-Security-Policy")
	w.Header().Set("Vary", "Accept-Encoding")
	identity, _, err := encodings(r.Header.Values("Accept-Encoding"))
	if status == 406 || identity == 0 || err != nil {
		httpwire.Headers(w)
		w.Header().Del("X-Joy-Pi-Snapshot-Age-Ms")
		w.Header().Del("Content-Type")
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(status)
		return
	}
	httpwire.Error(w, r, status)
}
