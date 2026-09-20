# Board HTTP/JSON API v0.1.0

Status: **consolidated established overview contract; not implemented**. This is normative for JPB-008–013. Supplemental routing/headers are explicitly **Proposed P-02/P-05**. [Health integration](health-integration-v1.0.md) owns `HealthSnapshot`; [UX](ux-v0.1.0.md) owns browser interpretation.

S00 implements only the page/assets [shell](s00-foundation.md); `/api/v1/overview` is deliberately absent (404). Exact numeric transport is adopted under [ADR-001](adr-001-exact-json-integers.md): uint64 remains unquoted numeric JSON; a future validated RawMessage snapshot is an object, not a string/base64. No functional API proposal is ratified by implementing the shell.

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

One global 1 s attempt deadline includes connection, response headers and reading the entire body. There is no retry. Mapping precedence is intentional:

| Observation | Board reason |
|---|---|
| Deadline expires before a complete acceptable 200 body is read; timeout while connecting, waiting for headers or stalled/slow body | `timeout` |
| Non-timeout network/transport failure before any valid HTTP response status is obtained, including DNS failure, refusal, reset or premature EOF before headers | `connection_failed` |
| Any parsed HTTP status other than 200, including redirects, 204, 500 and 503 | `upstream_error` |
| HTTP 200 but oversized body, unexpected encoding/media type under P-02, truncated non-timeout body, bad UTF-8/JSON, missing/wrong fields, bad types/ranges/nulls/issues, absent/non-string schema version | `invalid_response` |
| Well-formed JSON object with a present string schema_version other than `1.0` | `unsupported_schema_version` |

Once a non-200 status has been received, close its body and report `upstream_error` without parsing/relaying it; a later unread body timeout cannot change this classification. For 200, read/bounds/JSON integrity precede version dispatch; a syntactically broken unknown-version body is `invalid_response`. A complete valid supported payload is a success, including expected measure unavailability. Unknown-version fields are not interpreted as schema 1.0. Deadline expiration wins over a late completion; do not publish late bytes.

**Proposed P-02**: require Health media type `application/json` (optional UTF-8 charset accepted), no content encoding, body ≤64 KiB; disable redirect following. These defensive consumption rules complement the upstream emitted contract. Raw Go errors, upstream error bodies, OS paths, stack traces and network destinations never appear in public `reason` or Board messages. Invalid Board startup configuration is not a provider failure. A browser cancelling its request is not a fabricated Health failure; shared-flight semantics are in P-04.

## Browser-safe elapsed age header

**Proposed P-05**, additional header without changing the established JSON:

`X-Joy-Pi-Snapshot-Age-Ms: 12000`

When `last_success_at` is non-null, Board sends the ceiling of nonnegative monotonic elapsed milliseconds since that success, measured at response generation with the same anchor used for expiry. This also applies to expired `snapshot_state: none`. Omit the header only before first success. Encode as decimal digits; never compute it by subtracting serialized wall timestamps. No CORS expose-header configuration is needed on the same origin. Future clients validate it as a nonnegative integer, rejecting missing/invalid age when a snapshot is supplied.

The browser adds the full measured request round-trip duration and time since receiving the response, giving a conservative age bound. This includes response transit time without assuming clock synchronization and handles a first load that receives stale data. It may hide metrics slightly early, never deliberately late. [UX](ux-v0.1.0.md) defines lifecycle and failure behavior. Ratify this extension before S02/S03; do not silently substitute wall-clock subtraction if it is rejected—review another elapsed-age solution first.

## Supplemental HTTP surface

All rules in this section are **Proposed P-02**, not previously fixed Board decisions.

| Route | Methods / behavior |
|---|---|
| `/` | GET and HEAD serve embedded index; HEAD has no body; no SPA fallback for arbitrary paths |
| `/assets/` followed by an exact generated filename | GET/HEAD only; actual embedded asset and MIME type; no directory listing or traversal |
| `/api/v1/overview` | GET only; parameterless, no request body |
| Any other route | 404; never return index HTML for an unknown API path |

Route matching is exact: no trailing-slash redirect, encoded path aliases, `..` normalization or invented endpoint. Configure generated initial assets under `/assets/`; favicon if supplied uses that directory. Evaluate route first, then method, then request validity. Unknown path wins with 404. Unsupported method on a known route is 405 with `Allow: GET` for overview or `Allow: GET, HEAD` for static routes. HEAD on overview remains 405 with no body; OPTIONS has no special CORS support. Nonempty queries, forbidden bodies (positive/unknown Content-Length or transfer encoding) are 400; body rejection closes the connection without consumption. Invalid requests must not call Health. A zero-length body is permitted.

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

Proposed operational bounds: maximum header setting 8 KiB (`net/http` may reserve additional parser buffer; test actual behavior), header timeout 2 s, idle timeout 15 s, write deadline 3 s, at most 32 active application requests admitted without an unbounded queue. The upstream deadline remains 1 s, and shutdown remains ≤5 s. Capacity/bounds are reviewable P-02 choices; they do not claim protection for Internet exposure.

## Static and security headers

**Proposed P-02 concrete values:** production responses carry `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`; HTML carries:

```text
Content-Security-Policy: default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; base-uri 'none'; object-src 'none'; frame-ancestors 'none'; form-action 'none'
```

Bundle CSS/JS as external same-origin assets with no inline scripts/styles required. No external fonts/images, unsafe HTML injection or runtime remote resources. Vite development HMR may need a separately scoped dev-only policy; never weaken the production CSP to accommodate it. Test actual build output under this policy.

Index uses no-store. Fingerprinted static files may use `public, max-age=31536000, immutable`; API and Board errors always no-store with no ETag/Last-Modified. P-02 proposes precompressed deterministic gzip assets plus identity fallback, correct MIME, `Content-Encoding` and `Vary: Accept-Encoding`. S04 verifies actual delivered compressed bytes and budgets; compression at build time adds no runtime executable dependency.

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
