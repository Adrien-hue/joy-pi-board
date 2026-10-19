package httpui

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/fs"
	"strings"
)

type Representation struct {
	Path   string `json:"path"`
	Bytes  int    `json:"bytes"`
	SHA256 string `json:"sha256"`
}
type Resource struct {
	URL      string         `json:"url"`
	MIME     string         `json:"mime"`
	Identity Representation `json:"identity"`
	Gzip     Representation `json:"gzip"`
}
type manifest struct {
	Schema    int        `json:"schema"`
	Resources []Resource `json:"resources"`
}

// Inventory verifies the embedded allowlist and hashes. No runtime decompressor.
func Inventory(assets fs.FS) (map[string]Resource, error) {
	fail := errors.New("invalid embedded asset inventory")
	data, err := fs.ReadFile(assets, "manifest.json")
	if err != nil || len(data) > 65536 {
		return nil, fail
	}
	var m manifest
	if json.Unmarshal(data, &m) != nil || m.Schema != 1 || len(m.Resources) < 3 || len(m.Resources) > 128 {
		return nil, fail
	}
	out := make(map[string]Resource)
	for _, r := range m.Resources {
		p := r.Identity.Path
		mime := ""
		if p == "index.html" && r.URL == "/" {
			mime = "text/html; charset=utf-8"
		} else if strings.HasPrefix(p, "assets/") && strings.Count(p, "/") == 1 && r.URL == "/"+p {
			if strings.HasSuffix(p, ".js") {
				mime = "text/javascript; charset=utf-8"
			}
			if strings.HasSuffix(p, ".css") {
				mime = "text/css; charset=utf-8"
			}
		}
		if !fs.ValidPath(p) || mime == "" || r.MIME != mime || r.Gzip.Path != p+".gz" {
			return nil, fail
		}
		if _, ok := out[r.URL]; ok {
			return nil, fail
		}
		for _, v := range []Representation{r.Identity, r.Gzip} {
			b, e := fs.ReadFile(assets, v.Path)
			sum := sha256.Sum256(b)
			if e != nil || v.Bytes != len(b) || v.Bytes <= 0 || hex.EncodeToString(sum[:]) != v.SHA256 {
				return nil, fail
			}
		}
		out[r.URL] = r
	}
	if _, ok := out["/"]; !ok {
		return nil, fail
	}
	return out, nil
}
