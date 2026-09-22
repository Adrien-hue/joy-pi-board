# Sprint 02 — Overview API, Health transport and resilience

Status: **PREPARED, NOT EXECUTED**, 2026-09-22. S00 and S01 are closed. This assignment authorizes documentation only, replacing the temporary S01-closure restriction. No S02 application, fixture, dependency, toolchain or workflow change is made. All S02 tests and closure gates below are **NOT RUN**. S03 dashboard/polling/visual behavior and S04 packaging/physical acceptance remain outside this sprint.

## Starting point and authority

Preparation starts from clean Board HEAD `f7febfe9763cec8feb0be1446a1fcb9fb7a5f17b` (`docs: close Sprint 01 with verified CI evidence`). Its difference from S01 implementation commit `ccd1245a64fc372faa4ebc85eb000f78a01bd396` is nine Markdown files only. [S01 evidence](s01-evidence.md) records run 35637699325 attempt 1 for **ccd1245**, not f7febfe or future S02 changes. [S00 evidence](s00-evidence.md) retains run 35527458889 for `b6e139c3378cd77fdcb4c5edc610977c08bbf72a`, separately from its documentary closure commit `5f0d887bfc43e35711bc5a52f08164cf5197f3d5`. Historical fingerprints, incidents, measurements and NOT RUN limits are untouched. This preparation reads that evidence; it does not rerun CI or application tests.

The [specification](specification-v0.1.0.md) owns established requirements. [Board API](board-api-v0.1.0.md) owns the envelope, HTTP policy and proposed age protocol; [architecture](architecture-v0.1.0.md) owns configuration, lifecycle and proposed shared-flight realization. [Health integration](health-integration-v1.0.md) and [ADR-001](adr-001-exact-json-integers.md) remain unchanged. [Decisions](decisions.md) distinguishes approved requirements from **PROPOSED P-02/P-04/P-05** recommendations. This document owns tasks and integration seams; the [S02 validation matrix](validation-v0.1.0.md#prepared-s02-validation-matrix) owns test IDs and oracles. P-06 is not ratified.

## Scope and invariants

Implement in a later authorized task: demand-only Health HTTP client, exact `GET /api/v1/overview`, one latest valid snapshot in RAM, immutable publication, current/stale/none, safe classifications, concurrency/cancellation, elapsed age, configuration, exact routing, bounded shutdown and transition logs. The existing React shell is sufficient and must still start without Health; there is no frontend fetch loop or fake dashboard data.

Health unavailability alone yields Board **HTTP 200**, with the established JSON and no-store. Board defects/admission failures use distinct HTTP errors. A partial valid response replaces all previous values, including nulls. After failure, fallback is eligible through **30 seconds inclusive**, then snapshot is null while last_success_at remains. Before first success both are null. Restart loses RAM state. `observed_at` is untouched Health time, `last_success_at` is successful receipt/validation time, `generated_at` is response generation. None is substituted for elapsed-time measurement.

No polling, retry, completed-result reuse, history, disk cache, system observation or Health internal imports. Health remains read-only at `be7a0d824f62b94c842d8e5326110b1852c5a0bc`. Approved endpoints remain Board `0.0.0.0:8081` and Health `http://127.0.0.1:8080/v1/snapshot`.

## Fit to the existing implementation

| Inspected seam / test | S02 integration recommendation (not implemented) |
|---|---|
| [healthschema.Decode and Validated](../internal/healthschema/types.go) | Call Decode once after a bounded complete HTTP 200 body. Map typed Error.Kind into the two existing schema failure reasons. Keep Validated, not its known-field Snapshot projection, as the transport/cache value. Validated copies privately owned bytes; Bytes/MarshalJSON return copies, so sharing the value is safe. |
| [S01 contract tests](../internal/healthschema/snapshot_test.go) | Retain TestSharedCorpus, TestImmutableAndObjectTransport, TestPreservesUnknownAndOriginalTokens and TestRoundtripAtLimitWithMarkup. Add HTTP-level regressions without altering pinned fixture bytes or weakening S01 limits. |
| [httpui.NewHandler(fs.FS)](../internal/httpui/handler.go) | Keep embedded asset ownership in httpui. Add one outer httpapi dispatcher for exact route/method/request/error policy, with an overview service dependency. Refactor static error responses only as needed to share the safe error writer; do not create a second asset router or generic framework. |
| [shell tests](../internal/httpui/handler_test.go) | Preserve exact GET/HEAD asset and CSP checks. The static-only handler may still reject overview; the composed production handler must now provide it. Replace the standalone smoke harness's production-overview-404 assertion with controlled unavailable-Health 200 assertions in S02, without changing the React page. |
| [cmd run(args, output)](../cmd/joy-pi-board/main.go) | Introduce pure startup parsing and lifecycle wiring. Existing --listen/--help tests remain relevant; production run currently blocks in Serve and exposes wrapped listen errors, so safe diagnostics, --version, environment support and signal shutdown are real S02 work. Preserve ephemeral port 0 for integration tests. |
| [parseSnapshot](../web/src/health/snapshot.ts), [parseDocument](../web/src/health/json.ts), [contract harness](../scripts/contract.mjs) | S01 accepts a standalone snapshot, not an overview. Reuse it through a test-only lossless envelope adapter for real HTTP interop; do not pass an 80 KiB/depth-34 envelope to the 64 KiB/depth-32 function. No production dashboard/envelope Fetch consumer is built in S02. |

Recommended new packages are small, Board-owned `internal/config`, `internal/health`, `internal/overview`, and `internal/httpapi`; existing `internal/httpui`, `internal/healthschema` and `web` remain. Create packages only when their tasks supply behavior. No empty packages or executable test framework are created during preparation.

### Minimal future interfaces

These are responsibilities and candidate signatures, not committed code/API guarantees:

- `config.Parse(args, lookupEnv)` returns a validated configuration or help/version action and a safe error. Environment lookup injection prevents process-global test mutations. Binding is a separate step; configuration never probes Health.
- `health.Fetch(ctx) (healthschema.Validated, AttemptError)` performs one HTTP attempt. AttemptError carries one stable provider reason or a distinct lifecycle cancellation; it does not expose raw errors to HTTP. Use a concrete standard-library client with constructor seams for RoundTripper/dial tests, not a broad provider plugin interface.
- `overview.Demand(ctx)` joins/starts a flight and returns an immutable **flight view**: sequence, outcome, snapshot and success anchor belonging to that outcome. A handler never rereads the latest cache to combine an old reason with a newer snapshot. Error return distinguishes Board shutdown/capacity from provider failure.
- Clock dependency supplies one paired reading of wall UTC and opaque elapsed tick, plus a cancellable deadline timer. Production elapsed time retains Go's monotonic component; fake clock advances wall and elapsed independently and explicitly fires due timers. Do not use fake 1970 timestamps with real context deadlines. Keep the seam local to the coordinator and deadline creation, not a general scheduler.
- `overview` finalization uses a flight view and a clock reading to select snapshot/state and age together. HTTP encoding receives an immutable prepared response; no ResponseWriter or network I/O is invoked under the coordinator mutex.
- Lifecycle owner starts Serve, closes admission/cancels provider work, drains under one 5 s shutdown budget, then forces close. A small injected event sink records bounded transition enums; no logger interface needs the snapshot body.

## Tasks and dependency order

Each task is **PREPARED / NOT EXECUTED**. Test IDs below are planned, never PASS.

| Task | Depends on | Deliverable for later implementation | Verifiable exit condition |
|---|---|---|---|
| S02-01 — ratify and freeze | S01 closure | Resolve P-02/P-04/P-05 recommendations, record approved/replaced details, preserve P-01/P-03/ADR-001 | Decision review completed before dependent behavior; no implicit blanket approval. |
| S02-02 — configuration and seams | S02-01 | Pure flags/env validation, safe startup actions, minimal clock/lifecycle/test seams | S02-T01; help/version never bind or contact Health; invalid config fails before Serve. |
| S02-03 — Health transport | S02-02 | One bounded, cancellable HTTP attempt; typed S01 validation; exact reason precedence | S02-T02/T03/T04; one logical outbound request, no redirect/replay/background work. |
| S02-04 — coordinator and RAM state | S02-03 | Single active flight; one cache entry; immutable outcome publication; cancellation/deadline/sequence handling | S02-T05/T06/T07; deterministic boundaries and ordering, race checks. |
| S02-05 — HTTP overview and age | S02-04 | Exact JSON envelope, safe HTTP errors, response finalization, age header, bounded serialization | S02-T08/T09; complete/partial/stale/none wire results and numeric/size/depth regressions. |
| S02-06 — compose server and operations | S02-02/S02-05 | Exact static/API dispatch, admission/connection limits, signals, safe transition logging, unchanged shell | S02-T10/T11; no provider boot dependency, bounded shutdown and no log leakage. |
| S02-07 — real HTTP interop and regression | S02-03–S02-06 | Shared corpus exercised through actual Go HTTP and TS S01 parser, controlled network cases, retained shell/embed proof | S02-T01–T11 executed; measured real-deadline cases separately labeled; no UI/physical claims. |
| S02-08 — evidence and CI closure | S02-07 | Local evidence ledger, existing check/contract/race/CI extensions, new hosted run at exact implementation SHA | S02-T12; all applicable S02 gates PASS locally and remotely; otherwise “implemented locally, CI closure pending” or exact outstanding failure. |

## Real envelope and exact-number regression

P-02's [serialization policy](board-api-v0.1.0.md#snapshot-and-envelope-byte-boundaries) recommends a bounded in-memory encoder with HTML escaping disabled for application/json, validated raw snapshot marshaling, no indentation and one optional final newline included in the response limit. No generic float64 decoding/re-encoding is allowed. Unknown members, numeric digits, nulls and issue order survive. Output errors are Board 500, never fabricated Health invalid_response. No JSON is embedded in HTML.

Snapshot input remains **65,536 bytes / 32 containers**. The proposed full overview limit is **81,920 bytes / 34 containers**: outer object + health object add two levels. It is not a new Health limit. The wrapper has only bounded enums and Board-generated RFC3339 timestamps; reserve at most 1,024 bytes for it, assert that bound, and assert the compact raw snapshot does not grow. Thus every accepted maximum-size snapshot fits without truncation. Check the final counted bytes including newline; do not rely on Content-Length alone.

Required regression: derive, in tests using byte operations, the exact S01 65,536-byte valid snapshot with unknown markup text containing `<`, `>` and `&`; send it from simulated Health through actual Board HTTP. Check application/json, preserved object and integer tokens, no optional six-byte HTML escape expansion, accepted snapshot extraction and correct envelope length. Add raw/escaped Unicode and depth-32 snapshots; whole envelope depth 34 passes while depth 35 fails the proposed envelope guard. Invalid Health input at 65,537 bytes is still rejected before parsing. Default json.Marshal's expansion is a known S01 test fact, not a reason to reduce Health's accepted size.

The **test-only TS adapter** may use the installed lossless-json parse/stringify pair to extract/re-emit the nested object with preserved numeric tokens, then call S01 parseSnapshot and decimalUInt64. It must never use response.json or ordinary numeric JSON.parse/stringify. Check raw HTTP unknown-member retention independently in Go with UseNumber; a TS known-field projection is not a replacement transport. Run maximum-byte/depth cases through this adapter; if its reserialization cannot preserve the accepted boundary, use a library-supported exact-token extraction seam rather than weakening S01 or writing a generic parser. A hardened production overview parser and browser state/scheduling remain S03. The S01 standalone parsing entry points and their limits must keep their regression coverage.

## Test approach, evidence and closure

The [S02 matrix](validation-v0.1.0.md#prepared-s02-validation-matrix) maps tasks to stable requirements and A-05/B-01/B-05–B-09. Prefer httptest.Server for actual HTTP, net.Pipe/custom RoundTripper for malformed/blocked transport and a fake clock for expiry and publication. Use channels for “headers reached”, “body held”, “validation ready”, “waiter joined”, “finalization reached” and explicit releases; wait for events rather than arbitrary sleeps. All test goroutines and servers have cleanup and outer watchdogs.

Only a small separately named group uses real time: real 1 s deadline with withheld headers, stalled body and slow progress; actual connection refusal on a reserved-and-closed loopback listener; SIGTERM shutdown with a held request; header/write timeouts. Assert cancellation/events and generous harness completion bounds (proposed 2 s watchdog for the 1 s network cases, recording actual durations), not workstation p95/physical budgets. Exact timeout precedence and 29.999/30/30.001 s assertions use fake elapsed time. A watchdog is not an enlargement of the production 1 s deadline. No 30 s sleeping test.

During implementation produce `docs/s02-evidence.md` only with actual results: starting/tested SHA or explicit dirty-tree fingerprint, commands/tool versions, approved proposal references, local test IDs/results, corpus provenance, Go HTTP → TS evidence, exact outbound counts, real-time timings labeled by environment, race results, native/ARM64 checksums, standalone asset proof and unresolved limits. Preserve failed attempts and corrected incidents. Keep all S00/S01 evidence unchanged. Extend existing `check`/`contract` tasks and CI foundation/race coverage only during implementation; frontend still precedes Go embed. Retain exact tools/action pins, clean install, format/type/lint/vet/modules, S01 corpus/fuzz and shell regression. No dependency/toolchain upgrade is needed by this preparation.

S02 closes only after its ratified decisions, task criteria and applicable matrix cases pass, followed by a **new hosted CI run testing the exact S02 implementation commit**. Record run/attempt/event/branch/jobs/steps, actual image/tools and artifacts separately from local hashes. S01's green run cannot close S02. Browser execution, dashboard B-10/B-11, physical C-01–C-06, packaging, release and deployment remain NOT RUN/out of scope. B-03 is extended only through HTTP and exact decimal primitives, not rendered UI. Product Class A/B/C acceptance is not implied.

## Review points before implementation

Ratification is still required for P-02's HTTP-only destination policy, exact configuration/error precedence, transport/admission limits and serialization choice; P-04's orphan/deadline/finalization rules; P-05's age rounding/transit, fail-closed behavior and explicit server-suspend limitation. Recommendations are concrete in their authority documents, not multiple interchangeable options. No question blocks this preparation. Implementation needs a new authorization and review of these proposals.

Editorial conflicts resolved by this preparation: the prior layout suggested replacing httpui wholesale, the old age prose did not define a response-generation cut, and an 80 KiB browser limit alone did not account for the two added containers or HTML escaping. The historic S00-only wording of JPB-025 is clarified without changing its stable scope-authorization requirement. No established enum, JSON field, timeout, stale threshold or upstream baseline changes.

## Preparation review

Documentary checks on 2026-09-22: 20 Markdown documents / 213 internal links and anchors checked; JSON documentation blocks parsed; eight S02 tasks and twelve planned validation rows traced; `git diff --check` clean. Diff review confirms only eight existing Markdown files plus this new work package. S00/S01 evidence and hashes, P-06's deferred decision row, application code, tests, fixtures, dependencies, tools and workflows are unchanged. The specification received only current-status/JPB-025 authorization clarification. No application build/test, commit, push, tag, release or deployment was performed. These documentary checks do not execute or pass any S02 test.
