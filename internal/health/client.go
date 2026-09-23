// Package health performs a single bounded HTTP attempt. The coordinator owns
// the global deadline, publication and cancellation; this package never retries.
package health

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime"
	"net"
	"net/http"
	"sort"
	"strings"

	"github.com/Adrien-hue/joy-pi-board/internal/healthschema"
)

type Reason string

const (
	Timeout                  Reason = "timeout"
	ConnectionFailed         Reason = "connection_failed"
	UpstreamError            Reason = "upstream_error"
	InvalidResponse          Reason = "invalid_response"
	UnsupportedSchemaVersion Reason = "unsupported_schema_version"
)

// Result has no raw error or untrusted diagnostic. Cancelled is a lifecycle
// event, never a reason exposed as provider unavailability.
type Result struct {
	Snapshot  healthschema.Validated
	Reason    Reason
	Cancelled bool
}
type Client struct {
	endpoint string
	http     *http.Client
}

func New(endpoint string) *Client {
	return newClient(endpoint, net.DefaultResolver.LookupIPAddr, (&net.Dialer{}).DialContext)
}

func newClient(endpoint string, lookup func(context.Context, string) ([]net.IPAddr, error), dial func(context.Context, string, string) (net.Conn, error)) *Client {
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	transport := &http.Transport{
		Proxy: nil, DisableCompression: true, DisableKeepAlives: true,
		Protocols: protocols, MaxResponseHeaderBytes: 8192,
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			host, port, err := net.SplitHostPort(address)
			if err != nil {
				return nil, err
			}
			ip := net.ParseIP(host)
			if ip == nil {
				addresses, err := lookup(ctx, host)
				if err != nil {
					return nil, err
				}
				var ips []net.IP
				for _, a := range addresses {
					if a.Zone == "" && a.IP.To16() != nil {
						ips = append(ips, a.IP)
					}
				}
				if len(ips) == 0 {
					return nil, errors.New("no address")
				}
				sort.Slice(ips, func(i, j int) bool {
					a, b := ips[i].To4(), ips[j].To4()
					if (a != nil) != (b != nil) {
						return a != nil
					}
					return bytes.Compare(ips[i].To16(), ips[j].To16()) < 0
				})
				ip = ips[0]
			}
			return dial(ctx, "tcp", net.JoinHostPort(ip.String(), port))
		},
	}
	return &Client{endpoint, &http.Client{Transport: transport, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}}
}

func (c *Client) Close() { c.http.CloseIdleConnections() }

func failed(ctx context.Context, err error, fallback Reason) Result {
	if errors.Is(ctx.Err(), context.Canceled) {
		return Result{Cancelled: true}
	}
	var e net.Error
	if errors.Is(ctx.Err(), context.DeadlineExceeded) || errors.As(err, &e) && e.Timeout() {
		return Result{Reason: Timeout}
	}
	return Result{Reason: fallback}
}

func (c *Client) Fetch(ctx context.Context) Result {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint, nil)
	if err != nil {
		return Result{Reason: ConnectionFailed}
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", "joy-pi-board")
	res, err := c.http.Do(req)
	if err != nil {
		return failed(ctx, err, ConnectionFailed)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Result{Reason: UpstreamError}
	}
	if !allowedHeaders(res) || res.ContentLength > healthschema.MaximumBytes {
		return Result{Reason: InvalidResponse}
	}
	body, err := io.ReadAll(io.LimitReader(res.Body, healthschema.MaximumBytes+1))
	if err != nil {
		return failed(ctx, err, InvalidResponse)
	}
	if len(body) > healthschema.MaximumBytes || len(res.Trailer) != 0 {
		return Result{Reason: InvalidResponse}
	}
	v, err := healthschema.Decode(body)
	if err != nil {
		var schema *healthschema.Error
		if errors.As(err, &schema) && schema.Kind == healthschema.UnsupportedSchemaVersion {
			return Result{Reason: UnsupportedSchemaVersion}
		}
		return Result{Reason: InvalidResponse}
	}
	return Result{Snapshot: v}
}

func allowedHeaders(res *http.Response) bool {
	values := res.Header.Values("Content-Type")
	if len(values) != 1 {
		return false
	}
	typ, params, err := mime.ParseMediaType(values[0])
	if err != nil || !strings.EqualFold(typ, "application/json") || len(params) > 1 {
		return false
	}
	for key, val := range params {
		if key != "charset" || !strings.EqualFold(val, "utf-8") {
			return false
		}
	}
	// ParseMediaType folds identical repeated parameters; the accepted grammar
	// permits a SINGLE parameter, including when the repeats are identical.
	if strings.Count(values[0], ";") > 1 {
		return false
	}
	enc := res.Header.Values("Content-Encoding")
	if len(enc) > 1 || len(enc) == 1 && !strings.EqualFold(strings.TrimSpace(enc[0]), "identity") {
		return false
	}
	return len(res.Trailer) == 0 && len(res.Header.Values("Trailer")) == 0 && !res.Uncompressed
}
