# Board HTTP/JSON API v0.1.0

Status: **consolidated established overview contract; not implemented**. This is normative for JPB-008–013. S02 is [PREPARED, NOT EXECUTED](sprint-02.md). Supplemental transport, routing and age rules below are explicitly **PROPOSED P-02/P-05**, awaiting ratification. [Health integration](health-integration-v1.0.md) owns `HealthSnapshot`; [UX](ux-v0.1.0.md) owns browser interpretation.

S00 implements only the page/assets [shell](s00-foundation.md); `/api/v1/overview` is deliberately absent (404). Exact numeric transport is adopted under [ADR-001](adr-001-exact-json-integers.md): uint64 remains unquoted numeric JSON; S01 Validated provides immutable raw object marshaling, never a string/base64. No functional API proposal is ratified by implementing the shell.

## Overview

`GET /api/v1/overview` returns one object with required fields:

| Field | Type and meaning |
|---|---|
| `generated_at` | RFC3339 UTC timestamp with `Z`, response generation |
| `health` | Non-null object |
| `health.availability` | `available` or `unavailable`: outcome of Health attempt |
| `health.snapshot_state` | `current`, `stale` or `none` |
| `health.last_success_at` | RFC3339 UTC with `Z` or null: Board receipt **and validation** success |
| `health.reason` | null or `timeout`, `connection_failed`, `upstream_error`, `invalid_response`, `unsupported_schema_version` |
| `health.snapshot` | HealthSnapshot object or null, never a JSON string |

`Content-Type: application/json` and `Cache-Control: no-store` are mandatory. Health failure alone remains a valid **HTTP 200** Board overview. No field named `status` replaces `availability`; there is no added Health complete/partial flag.

| Situation | availability | snapshot_state | reason | snapshot / last_success_at |
|---|---|---|---|---|
| Valid complete or partial retrieval | available | current | null | Received snapshot / new success time |
| Failure and last success age ≤30 s | unavailable | stale | Failure cause | Last valid snapshot / unchanged success time |
| Failure and last success age >30 s | unavailable | none | Failure cause | null / retained success time |
| Failure and no success since process start | unavailable | none | Failure cause | null / null |

At exactly 30 seconds a snapshot is eligible; at any greater age it is not. Compute age at response generation, not just when the failed upstream request began. Errors never advance or erase the last-success timestamp. A partial valid success is `available/current`, replaces the whole cache and resets its success anchor. Recovery occurs on the next demand attempt without restarting Board. Process restart resets the cache and last success.

`snapshot.observed_at` is produced by Health and retained unmodified; it is not `last_success_at`. The latter is Board's validated-reception time; `generated_at` is later response generation in ordinary clock conditions. Do not enforce cross-host timestamp alignment or wall-clock ordering across clock corrections. Server expiry uses elapsed time as specified in [architecture](architecture-v0.1.0.md).

Documentation-only discriminated TypeScript shape, using `HealthSnapshot` from the integration contract:

```ts
type Reason = "timeout" | "connection_failed" | "upstream_error"
  | "invalid_response" | "unsupported_schema_version";
type HealthResult =
  | { availability: "available"; snapshot_state: "current";
      last_success_at: string; reason: null; snapshot: HealthSnapshot }
  | { availability: "unavailable"; snapshot_state: "stale";
      last_success_at: string; reason: Reason; snapshot: HealthSnapshot }
  | { availability: "unavailable"; snapshot_state: "none";
      last_success_at: string | null; reason: Reason; snapshot: null };
interface Overview { generated_at: string; health: HealthResult }
```

Go response DTOs use string enums, UTC-formatted time strings and nullable timestamp/reason pointers without `omitempty`; snapshot is a validated `json.RawMessage` or literal JSON null. Input validation needs the presence-aware types described in the integration contract; output structs alone are not a validator.

## Exact failure classification

One global **1 s** attempt deadline includes resolution, connection, response headers, full body reading and the bounded validation/publication path. It starts when the demand installs its flight, not when a socket becomes available. No retry or timeout reset occurs on progress. Established reasons remain unchanged; the following **PROPOSED P-02/P-04 precedence** makes races testable:

1. Board shutdown/caller cancellation/admission refusal are handled separately, not mapped into Health reasons. A single departing waiter does not cancel the shared attempt.
2. At the terminal decision, elapsed time at or beyond the flight deadline wins over an otherwise uncommitted result: `timeout`. A transport timeout before that point is also `timeout`. A final outcome already committed before the deadline is never reclassified later.
3. Before a usable final HTTP status, any non-timeout transport/protocol failure is `connection_failed`: resolution, refusal, reset, header parsing/size failure or EOF. Interim 1xx does not count as a final status; terminal 101 is non-200.
4. A parsed final **non-200** committed before deadline is `upstream_error` (including 3xx, 204, 500, 503). Close without reading/parsing the error body; later body problems cannot override it. MIME/encoding is irrelevant on this branch.
5. For **200**, disallowed MIME/encoding, a declared/exceeded size limit or non-timeout truncated/body error is `invalid_response`. If the read/decode decision reaches the deadline before commitment, step 2 wins. Close on first decisive failure, without draining an unbounded body.
6. For a complete bounded 200 body, reuse S01 Decode. Global JSON/Unicode/duplicate/depth/root checks precede schema dispatch. Missing/wrong-type version or invalid schema 1.0 is `invalid_response`; a structurally valid object with a different string version is `unsupported_schema_version`, without interpreting its schema-specific fields. Unknown version never masks malformed JSON.
7. A valid complete or partial 1.0 snapshot committed before deadline is success (`available/current`, reason null). No useful-measure or health threshold is added by transport.

Public responses never expose upstream bodies, raw Go errors, URL/OS paths or stack traces. S01's bounded diagnostics may support a test assertion; the public surface still contains only the established reason. Failure does not advance/erase last_success_at. The [flight algorithm](architecture-v0.1.0.md#concurrent-demand-and-ordering) is the sole outcome publication authority. Tests: S02-T02–T07 in [validation](validation-v0.1.0.md#prepared-s02-validation-matrix).

### Proposed Health HTTP transport

**PROPOSED P-02.** Use a dedicated standard-library client/transport, not http.DefaultClient/DefaultTransport shared state. GET has no body, credentials, cookies or browser-supplied headers. Send `Accept: application/json` and `Accept-Encoding: identity`. Disable environment proxy use (`Proxy: nil`), automatic decompression (`DisableCompression: true`), keep-alive reuse (`DisableKeepAlives: true`) and HTTP/2/alternate protocol negotiation. Reject redirects with `http.ErrUseLastResponse`; do not follow Location or use a cookie jar. Do not resend after a stale/reset connection or application-level error. The extra TCP connection per flight is a deliberate correctness tradeoff, with Pi latency/CPU acceptance remaining S04.

One context/deadline applies to the entire flight. No independent longer client timeout, per-read timeout extension or retry middleware. `MaxResponseHeaderBytes` is proposed at 8,192 bytes; a 200 response is read through a limit of **65,537 bytes**, accepting at most **65,536** before S01 parsing. Reject Content-Length >65,536 immediately; absent length/chunked framing still uses the counted reader. A body exactly at the limit must reach proper EOF/framing completion before success. HTTP transfer framing is decoded by net/http; content compression is not accepted. Reject unexpected response trailers (announced or received) rather than assign them semantics.

Use `mime.ParseMediaType`: accept exactly application/json (case-insensitive media token), with no parameters except an optional single UTF-8 charset, case-insensitive. Reject missing/multiple Content-Type values, malformed parameters, other charsets, `+json` aliases, multipart and HTML. Content-Encoding must be absent or a single `identity` token; reject gzip/br/deflate, lists and ambiguous duplicate headers. These choices do not change Health's producer contract.

For hostname destinations, resolve only during the demand under the same deadline. Recommend selecting one resolved address (IPv4 first, otherwise IPv6; stable byte order), then one dial of that literal address, preserving the configured HTTP Host. No Happy Eyeballs/second-address connection fallback; a failed selected address waits for the next overview demand. Resolver DNS protocol retransmissions are not extra Health HTTP requests and remain inside the deadline. Literal IPs avoid DNS entirely. One flight has at most one connection attempt and one HTTP GET; refusal can mean zero requests reach Health. Tests count dials, accepted connections and parsed requests separately, including proxy poison, redirects and reset/reuse traps.

This removes known default replay behavior: Go documents retry eligibility on previously successful connections, so disabling reuse needs explicit regression coverage, not merely “no retry loop” in application code. See the primary [Transport documentation](https://pkg.go.dev/net/http#Transport). No transport optimization is ratified by this preparation.

## Browser-safe elapsed age header

**PROPOSED P-05**, additional header only; the established JSON envelope does not change:

`X-Joy-Pi-Snapshot-Age-Ms: 12000`

**Origin:** on successful reception and S01 validation, capture paired wall UTC and monotonic elapsed origin at validation completion and carry them in the publication candidate. Atomic publication commits that pair without resampling a newer success time after lock contention. Preserve the monotonic-bearing value; formatting UTC is performed on a copy. Timeout/cancellation/failed validation never creates a success anchor. Wall-clock steps can reverse the ordering of displayed timestamps but never change elapsed eligibility.

**Response-generation cut:** use the immutable flight view, prepare/compact the snapshot outside the lock, then sample wall/elapsed time at finalization immediately before bounded wrapper assembly and header commitment. This single cut selects current/stale/none, last_success_at, generated_at and the age header coherently. Do not reuse age computed at request start or before a slow serializer. If a failed-flight fallback crosses 30 s during the Health attempt or preparation, omit it at finalization; exactly 30 s remains eligible. Do not reread newer shared state while finishing an older view. A known pause before commitment requires re-finalization from that same view; never repair a partially written response. Keep this bounded by the response write deadline, with no unbounded regenerate loop; inability to finish yields safe Board 503 before headers, or a closed connection after commitment.

A successful flight delayed so long that its view is already >30 s old at finalization is a **Board response-delivery failure**: proposed HTTP 503 `temporarily_unavailable`, not invented Health unavailability or an available/none envelope. This exceptional path preserves the established successful-retrieval invariant. An invalid/regressing elapsed clock likewise fails closed with Board 500 `internal_error`, not wall-clock fallback or negative age. Once headers/body are committed, socket/transit delay is handled by S03's conservative age budget; the server cannot promise freshness at the recipient's rendering instant.

**Header format:** when last_success_at is non-null, emit one canonical decimal integer `0|[1-9][0-9]*`, no sign/fraction/exponent/whitespace. Use integer arithmetic for `ceil(elapsed_ns / 1,000,000)` without addition overflow. Cap at **9007199254740991** milliseconds, which is exactly representable by JavaScript; saturation is permanently expired and must never wrap below 30,001. Expiry uses unrounded duration, not the header. At 30,000.001 ms, snapshot is absent and age is 30,001. Missing success means omit the header, with null snapshot/time. Expired none retains both last_success_at and its age header; errors outside the overview carry no age header. This cap is a consumer-protocol choice, not a new schema uint64 rule.

**S02 guarantees:** origin, integer/header formatting, single-cut selection, no wall-step effect, boundaries and immutable views are server tests S02-T05/T08. Elapsed time uses the platform monotonic clock. Go notes that some systems pause that clock during sleep; normal Pi service operation must not rely on system suspend. Negative/unverifiable readings fail closed. Suspend-inclusive server timing, if required for deployment, needs an explicit platform choice before that acceptance; no wall-time subtraction silently replaces the anchor. See [Go monotonic-clock semantics](https://pkg.go.dev/time#hdr-Monotonic_Clocks).

**S03 obligations, not implemented in S02:** bound and validate the whole envelope first, including this single header. Missing/invalid/duplicate/comma-joined/out-of-range header when success metadata exists is a contract error: hide metrics; do not infer age from Pi timestamps. An age header before any success is also invalid. For supplied metrics use `header_age + full_request_round_trip + elapsed_since_receipt` as a conservative upper bound; compare with 30,000 ms. Include fetch and body-read duration, then processing time; failed Board requests never reset the anchor. A none response hides immediately. Header ages over 30,000 forbid display even with a contradictory current/stale body.

On hidden/pagehide/freeze/suspend or unverifiable elapsed timing, hide retained data and keep it hidden until a valid refresh on resumption. Use the larger nonnegative local monotonic/wall elapsed delta conservatively; a backward wall jump invalidates the age anchor. This is not alignment with the Pi wall clock. Timer throttling must not extend availability. The [UX lifecycle contract](ux-v0.1.0.md#scheduling-and-browser-freshness) remains S03 work; protocol ratification and server conformance do not count as browser validation. Tests involving real browser clocks/rendering remain B-10/B-11 NOT RUN.

## Supplemental HTTP surface

All rules in this section are **PROPOSED P-02** refinements, not previously fixed Board decisions. Outer shutdown/admission rejection precedes route handling; for admitted requests use route → method → request validity → provider demand. An unknown embedded filename is an unknown route. Reject non-origin-form/absolute-form request targets with 400; do not provide a forward proxy. Encoded aliases and noncanonical paths remain 404, with no cleaning redirect.

| Route | Methods / behavior |
|---|---|
| `/` | GET and HEAD serve embedded index; HEAD has no body; no SPA fallback for arbitrary paths |
| `/assets/` followed by an exact generated filename | GET/HEAD only; actual embedded asset and MIME type; no directory listing or traversal |
| `/api/v1/overview` | GET only; parameterless, no request body |
| Any other route | 404; never return index HTML for an unknown API path |

Route matching is exact: no trailing-slash redirect, encoded path aliases, `..` normalization or invented endpoint. Configure generated initial assets under `/assets/`; favicon if supplied uses that directory. Evaluate route first, then method, then request validity. Unknown path wins with 404. Unsupported method on a known route is 405 with `Allow: GET` for overview or `Allow: GET, HEAD` for static routes. HEAD on overview remains 405 with no body; OPTIONS has no special CORS support. Any query delimiter (including bare `?`, via ForceQuery), nonempty query, or forbidden body (positive/unknown Content-Length or transfer encoding) is 400; body rejection requests connection close and does not read the body in the application. net/http cleanup may perform bounded transport handling; a read deadline/close must prevent draining an attacker-controlled stream. Reject a forbidden Expect: 100-continue body without an application read or 100 response. Invalid requests must not call Health. A zero-length body is permitted.

| Board condition | HTTP / safe error code |
|---|---|
| Bad query/body/request | 400 / `invalid_request` |
| Unknown route | 404 / `not_found` |
| Wrong method | 405 / `method_not_allowed` |
| Admission closed during shutdown or Board capacity exhausted | 503 / `temporarily_unavailable` |
| Board handler/encoding defect | 500 / `internal_error` |

Proposed fixed Board error bodies contain only `error.code` and `error.message`, with fixed messages respectively: `The request is invalid.`, `The resource was not found.`, `The method is not allowed.`, `Board is temporarily unavailable.`, `Board encountered an internal error.` Example:

```json
{"error":{"code":"internal_error","message":"Board encountered an internal error."}}
```

Board application errors use JSON/no-store and do not masquerade as provider overview responses. Encode before committing headers. HEAD responses have no body. Transport parser rejection (such as 431 or malformed HTTP), a closed connection or a failed socket write may have no JSON envelope; the browser treats those as a failed Board fetch. No promise is made to replace standard `net/http` transport errors.

Proposed operational bounds: MaxHeaderBytes 8,192 (`net/http` reserves additional parser buffer; test actual rejection instead of claiming an exact 8 KiB wire cap), ReadHeaderTimeout/ReadTimeout 2 s, IdleTimeout 15 s, WriteTimeout 3 s, at most 32 admitted application requests (including waiters/writes), 64 open accepted connections including idle clients, and 32 unfinished provider workers including cancelled-but-not-returned ones. Use nonblocking permits: application/worker exhaustion returns safe Board 503; listener overflow closes the newly accepted connection without an HTTP promise. No semaphore wait queue, rate-limit retry or queued background job. Release permits on completion/close exactly once, including errors. New demand still attempts Health unless rejected as a Board capacity error. OS listen backlog is finite external transport state, not an application queue. The ten-client capacity target is unchanged and must be tested. The upstream deadline remains 1 s, and shutdown remains ≤5 s. Capacity/bounds are reviewable P-02 choices; they do not claim protection for Internet exposure.

## Static and security headers

**P-02 recommendation for S02:** retain S00's existing values on static responses and apply the common nosniff/referrer headers to overview and safe errors. Responses carry `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`; HTML carries:

```text
Content-Security-Policy: default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; base-uri 'none'; object-src 'none'; frame-ancestors 'none'; form-action 'none'
```

Bundle CSS/JS as external same-origin assets with no inline scripts/styles required. No external fonts/images, unsafe HTML injection or runtime remote resources. Vite development HMR may need a separately scoped dev-only policy; never weaken the production CSP to accommodate it. Test actual build output under this policy.

S02 keeps identity bytes and no-store for index, assets, overview and Board errors; no new ETag/Last-Modified or compression negotiation. Preserve static GET/HEAD MIME and byte equality, CSP on HTML, and no remote resource references. These checks are planned for S02 and remain NOT RUN; rendered CSP/browser/mobile acceptance remains S03/S04. S04 retains the proposal for immutable caching of fingerprinted assets and precompressed deterministic gzip with identity fallback, correct MIME/Content-Encoding/Vary and actual delivered-size acceptance. No compression dependency/runtime tool or packaging is added in S02; P-06 measurement conventions remain deferred.

## Snapshot and envelope byte boundaries

**PROPOSED P-02.** Health input remains S01's 65,536 bytes inclusive and 32 containers, before validation. The full Board overview is separately limited to **81,920 identity bytes** (including any final newline) and **34 containers**: overview root 1, health 2, Health snapshot root 3. Never apply the snapshot limit to the complete envelope or increase the Health limit to accommodate it. Board-generated wrapper strings are bounded RFC3339 UTC timestamps (year 0000–9999) and fixed enums; assert wrapper overhead ≤1,024 bytes. Invalid local timestamps are a Board encoding defect.

Use the S01 Validated marshaler (or a private copy via json.RawMessage), never Snapshot re-encoding or generic float64 maps. Encode to a bounded private buffer with `json.Encoder.SetEscapeHTML(false)`, no indentation; count/remove the encoder's final newline explicitly. Preserve original number tokens, nulls, issues/order and unknown members. Compaction may change whitespace, not values. Precompact the raw object without growing it, then finalize age/state as described above. Assert final size before any headers; a violated output bound is Board 500, not a Health failure or a truncated 200.

Disabling optional HTML escaping is appropriate only for the separate application/json response, never inline HTML. [Go's encoder documentation](https://pkg.go.dev/encoding/json#Encoder.SetEscapeHTML) explains the default `<`, `>` and `&` escaping and final newline. S01's TestRoundtripAtLimitWithMarkup already demonstrates why default escaping can expand an accepted maximum-size payload. S02-T09 extends it through real HTTP, including raw/escaped Unicode, uint64 extremes and nested unknown fields. Do not assume the change to encoder options alone proves all Unicode/size cases; assert them under the pinned toolchain.

The future browser reader needs its own 80 KiB/depth-34 envelope profile and separate 64 KiB/depth-32 nested snapshot enforcement, exact-token parsing, duplicate/Unicode checks and runtime envelope invariants. S01 parseSnapshot is a standalone entry, not that envelope parser. S02's test adapter exercises actual HTTP output and reuses S01 validation without claiming S03 production browser parsing or scheduling is implemented; see [the regression plan](sprint-02.md#real-envelope-and-exact-number-regression).

## Complete examples

All example files are standalone valid JSON with real timestamps and snapshot objects:

| Example | Meaning |
|---|---|
| [overview-current.json](examples/overview-current.json) | Full success, true firmware flag, CPU already percent, exact uint64 counters |
| [overview-partial.json](examples/overview-partial.json) | Partial success replacing a full snapshot; all three issue codes |
| [overview-stale.json](examples/overview-stale.json) | Failure at age 12 s with last snapshot retained; proposed age header 12000 |
| [overview-expired.json](examples/overview-expired.json) | Failure at age 30.001 s; snapshot hidden; proposed age header 30001 |
| [overview-never-success.json](examples/overview-never-success.json) | Failure before any success; no age header |

Examples are illustrative data, not fabricated run results. Failure reason is tied to the attempted retrieval, not to any individual null measure. The response has no Board availability field: browser reachability is local state, described in the UX contract.
