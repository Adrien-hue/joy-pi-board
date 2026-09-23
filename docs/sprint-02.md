# Sprint 02 — Overview API, Health transport and resilience

Status: **COMPLETE**, 2026-09-23, with preserved local results and verified hosted CI for `a11f07c85ff4071a13ca3ae9227920e9f1c7dda4` (run 35908782107, attempt 1). S00/S01 remain closed with separate historical proofs. The earlier implementation instruction approved the concrete P-02/P-04/P-05 recommendations prepared on 2026-09-22; this documentary closure changes no technical decision. S03 is ready for preparation, NOT STARTED; S04 is NOT STARTED. The current task permits documentation only, with no application build/test, commit, push, tag, release or deployment.

## Starting point and authority

Implementation started from clean Board HEAD `f3d8dc8959c76e78b53d2913ef068baa89618567` (`docs: prepare Sprint 02 overview API and resilience`). Its predecessor `f7febfe9763cec8feb0be1446a1fcb9fb7a5f17b` closed S01 documentation. S01 implementation `ccd1245a64fc372faa4ebc85eb000f78a01bd396` is the only commit covered by hosted run 35637699325 attempt 1; that run does not validate S02. [S00 evidence](s00-evidence.md) and [S01 evidence](s01-evidence.md), including original hashes and incidents, remain unchanged.

The [specification](specification-v0.1.0.md) owns requirements; [Board API](board-api-v0.1.0.md) owns HTTP/envelope/age rules; [architecture](architecture-v0.1.0.md) owns concurrency/configuration/lifecycle. [Health integration](health-integration-v1.0.md), [ADR-001](adr-001-exact-json-integers.md) and the upstream baseline are unchanged. [Decisions](decisions.md) records current approval provenance and deferred P-06/S04 choices. [S02 evidence](s02-evidence.md) owns actual results and input identity; the [validation matrix](validation-v0.1.0.md#prepared-s02-validation-matrix) retains stable S02-T01–T12 assertions.

## Scope and invariants

Delivered: demand-only Health HTTP integration, `GET /api/v1/overview`, one immutable latest-valid snapshot in RAM, safe reason classification, shared active calls, elapsed-time fallback, precise routing/configuration, bounded resources/shutdown and transition logs. The existing React page remains sufficient and has no metric data or Fetch/polling consumer.

Health failure produces a valid Board HTTP 200 overview. Board defects/admission failures remain distinct HTTP errors. Partial success replaces all prior values, including nulls. On failure the fallback is eligible through 30 seconds inclusive, then snapshot is null while last_success_at is retained. Restart loses RAM state. observed_at, last_success_at and generated_at keep distinct meanings; no wall-time subtraction controls freshness. No polling, retry, completed-result window, persistence, collectors or Health internal imports. Defaults remain Board `0.0.0.0:8081` and Health `http://127.0.0.1:8080/v1/snapshot`.

## Fit to the existing implementation

| Existing seam | Actual integration |
|---|---|
| [healthschema.Validated](../internal/healthschema/types.go) | Sole schema authority. The client bounds the complete HTTP body then calls Decode. Cache/views retain private validated bytes; typed projection is used only to derive safe issue path/code transition keys. |
| [httpui](../internal/httpui/handler.go) | Retains ownership of embedded asset routes. [httpapi](../internal/httpapi/handler.go) composes overview; [httpwire](../internal/httpwire/errors.go) supplies fixed errors/common headers. No duplicate asset router. |
| [main](../cmd/joy-pi-board/main.go) | Pure [config](../internal/config/config.go) parsing, idle client construction, listener, signal lifecycle and bounded [eventlog](../internal/eventlog/log.go). No startup Health request. |
| [overview](../internal/overview/coordinator.go) | One active logical flight, separate 32-worker bound, original deadline, checked paired clock, atomic terminal view and private cache. Network, projection/compaction and log delivery run outside the lock. |
| [response](../internal/overview/response.go) | Original snapshot bytes plus bounded JSON wrapper. Age/state/timestamps are finalized coherently; at most three assemblies if a known pause changes age rounding/eligibility. |
| [HTTP interop harness](../scripts/http-contract.mjs) | Test-only real HTTP reader reuses S01 lossless parsing/decimal primitives. No production envelope Fetch consumer, UI state or browser scheduler. |

### Clock and test seams

Clock exposes paired wall/elapsed reads and a cancellable deadline timer. Production keeps monotonic elapsed from its original Go time value; fake clocks fire explicitly and move wall time independently. Regression/invalid readings fail closed. Go's monotonic clock may pause during server suspension; this limitation remains explicit and is not repaired with wall time. Private barriers before publication and around finalization exercise delayed candidates and response preparation without arbitrary 30-second sleeps.

## Tasks and dependency order

| Task | Dependencies | Delivered result | Status / evidence |
|---|---|---|---|
| S02-01 — ratify and freeze | S01 closure | Implementation-instruction P-02/P-04/P-05 approval recorded; P-06/S04 remain deferred. | DONE, decision register. |
| S02-02 — configuration and seams | S02-01 | Presence-aware flag/env/default parsing, safe actions/errors, paired clock, timer/test barriers. | LOCAL AND CI PASS: S02-T01. |
| S02-03 — Health transport | S02-02 | Dedicated bounded HTTP/1 attempt, one selected address/dial, no proxy/reuse/redirect/decompression; S01 Decode. | LOCAL AND CI PASS: S02-T02/T03/T04. |
| S02-04 — coordinator and RAM state | S02-03 | Active sharing, one cache, immutable views, atomic completion/deadline/shutdown, late-worker rejection. | LOCAL AND CI PASS: S02-T05/T06/T07 and Linux race. |
| S02-05 — HTTP overview and age | S02-04 | Exact object envelope, fixed errors, age header, single-view finalization and byte/depth boundaries. | LOCAL AND CI PASS: S02-T08/T09. |
| S02-06 — compose server and operations | S02-02/S02-05 | Static/API composition, admission/connection limits, signals, bounded logs and unchanged page. | LOCAL AND CI PASS: S02-T10/T11. |
| S02-07 — real HTTP interop and regression | S02-03–S02-06 | Shared corpus through Health HTTP → Board HTTP → TypeScript; socket/deadline tests and standalone smoke. | LOCAL AND CI PASS: S02-T01–T11. |
| S02-08 — evidence and CI closure | S02-07 | Local ledger, source fingerprint, regression/build artifacts and existing CI integration. | DONE: local ledger plus S02-CI-01–S02-CI-06, exact-commit hosted run and documentary exit review. |

## Real envelope and exact-number regression

Snapshot input remains **65,536 bytes / 32 containers**. The full overview is **81,920 identity bytes / 34 containers**, including its newline; wrapper overhead is asserted ≤1,024 bytes. The extra two containers are Board root and health. Unknown members, integer tokens, nulls and issue order are retained without float64 conversion or Health model re-encoding.

The implementation compacts private validated bytes before the clock cut, encodes the bounded wrapper with `SetEscapeHTML(false)`, then inserts the object at the fixed final snapshot field. This avoids optional HTML escaping and avoids a second full-object scan during finalization. Output is buffered and bounded before headers; malformed local state yields Board 500, never a truncated 200.

The real HTTP corpus adds byte-derived maximum `<>&`/Unicode and depth-32 cases to the unchanged shared corpus. A first test-adapter attempt using lossless-json parse/stringify was rejected by the existing prototype corpus: library duck typing could mistake an unknown inherited `isLosslessNumber` member for a number. The final test adapter extracts original snapshot text from the known fixed wrapper and passes it to S01 parseSnapshot; it is not a general JSON parser or production consumer. The library's original-number callback still parses the envelope for field/token assertions. Server-side raw-byte checks independently establish preservation. S03 must implement its own hardened envelope consumer rather than reuse this test framing assumption.

## Test approach, evidence and closure

[S02 evidence](s02-evidence.md) maps every matrix row to actual test names, commands, results, durations and limitations. Fake elapsed time covers 29.999 / 30 / 30.001 seconds, exact deadline priority and wall jumps. Channels synchronize ten callers, cancellation, late workers, blocked writes and bounded permits. A small separate real-time group exercises one-second header/body/progress cancellation, two/three-second server deadlines, refusal and process SIGTERM. These are workstation tests, not Raspberry Pi latency/resource acceptance.

Existing `check` now includes the S02 Go suites and real HTTP → TypeScript contract harness after the dedicated parser build. Existing `race` covers all new Go packages. CI action/tool pins and OS label are unchanged; workflow labels explain expanded coverage. `build` and `cross` still build frontend before embed; smoke now expects unavailable-Health HTTP 200 while checking identical embedded assets, empty PATH and no idle provider calls.

**Closure satisfied:** [verified hosted CI](s02-evidence.md#verified-hosted-ci) records run 35908782107 attempt 1, successful foundation/race jobs, exact source, observed environments and confirmed S02-T01–T12 coverage. The implementation commit matches local HEAD and remote main at verification; the 77 historical local input hashes also match its Git blobs. No S00/S01 run is reused. Target browsers, dashboard B-10/B-11, visual/mobile acceptance, Debian/systemd, real Health physical C-01–C-06 and release remain NOT RUN. S02 local and CI subsets do not imply full A/B/C product acceptance.

## Review points before later sprints

No unresolved choice blocks the delivered S02 software scope. S03 requires separate authorization and P-06 presentation review; its P-05 browser obligations are fixed but unimplemented. S04 requires the still-deferred compression/packaging/performance protocols and actual Pi 3B+ evidence. Suspend-inclusive server timing, if deployment needs it, requires an explicit platform decision before that acceptance.

## Historical preparation review

At the documentary-only preparation on 2026-09-22 (before implementation approval): 20 Markdown documents / 213 internal links and anchors checked; JSON documentation blocks parsed; eight S02 tasks and twelve planned validation rows traced; `git diff --check` clean. Diff review confirms only eight existing Markdown files plus this new work package. S00/S01 evidence and hashes, P-06's deferred decision row, application code, tests, fixtures, dependencies, tools and workflows are unchanged. The specification received only current-status/JPB-025 authorization clarification. No application build/test, commit, push, tag, release or deployment was performed. These documentary checks do not execute or pass any S02 test.
