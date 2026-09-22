# Validation and release strategy

Status: **S00 AND S01 COMPLETE; foundation and S01 contract subsets PASS locally and in CI, full A/B/C product acceptance NOT RUN**. [S00 evidence](s00-evidence.md) records the local results and verified run 35527458889 for `b6e139c3378cd77fdcb4c5edc610977c08bbf72a`. S01 in-memory contract subsets are verified locally and in hosted run 35637699325 attempt 1 at `ccd1245a64fc372faa4ebc85eb000f78a01bd396`; [S01 evidence](s01-evidence.md) preserves local input provenance separately from CI. S02 is ready to be prepared, NOT STARTED. S02–S04 functional/physical tests remain NOT RUN. Normative for release evidence and requirement traceability (JPB-024). Product budgets are established in [the specification](specification-v0.1.0.md); durations/sampling/statistics remain **Proposed P-06**. Shell-only asset size observations are explicitly preliminary. Historical documentary checks are in [documentation review](documentation-review.md).

Numeric expectations are adopted through [ADR-001](adr-001-exact-json-integers.md), not an open library proposal. S00 [shared fixture](../testdata/README.md) byte checks are not B-03 end-to-end validation. Future tests must cover original number tokens, uint64 bounds, nulls and schema-driven float/BigInt conversion across Health → Board → browser → display. S01 approves unknown-member tolerance under global structural checks, independently of numeric exactness.

## Evidence and release rules

Each run records test ID, PASS/FAIL/NOT RUN, UTC date, operator/CI URL, Board source commit, exact binary SHA-256, package SHA-256 if applicable, toolchain/lockfile identities, OS/kernel/architecture, configuration, Health commit/version, client/load settings, raw evidence and any exception. Record compiler flags and embedded asset manifest/checksums in build metadata. A mock is identified as a mock. Hardware model/OS are operator evidence, not fabricated Health API fields.

Release requires **Class A PASS, Class B PASS and Class C PASS**, with every required check executed and passing. A waiver or unexecuted required check does not count as PASS. Relevant P-xx choices must be ratified or replaced explicitly before their dependent gates can be accepted. An unexecuted physical check remains NOT RUN, not inferred PASS from a mock, developer workstation, cross-build or Health's own hardware release status.

Build the candidate once after source checks; identify it with SHA-256. Run binary integration and physical acceptance on those exact bytes. For architecture-specific source/race checks, attach results to the same source/toolchain inputs, then attach A-04, B binary runs and C runs to the ARM64 candidate checksum. Stage the same binary in the Debian package and prove the extracted/installed binary hash matches. Promote/copy the accepted artifacts without rebuilding. A new binary/package or embedded asset change invalidates affected acceptance and requires a new candidate identity; a same-version rebuild is not automatically accepted. Packaging/install scripts and package hash have their own S04 checks.

## Class A — build and correctness

These IDs define product gates. S00 commands now execute the foundation subsets traced below; later functional extensions and full product acceptance remain planned. A passing subset does not mark the whole product gate PASS.

| ID | Planned check and measurable PASS condition | First responsible sprint |
|---|---|---|
| A-01 | Go formatting has no diff; test/vet all first-party packages (`./cmd/... ./internal/... ./web` in S00, excluding Go files inside npm dependencies); native race tests succeed with no race; module verification/tidy leaves tracked files unchanged. No Health internal imports or metric collection. S00 subset results are in the evidence record. | S00, expanded S01/S02 |
| A-02 | Clean locked frontend install; type-check, lint, meaningful unit tests and production Vite build exit zero with no lockfile drift. Parser/types/formatting tests added by S01/S03. | S00 |
| A-03 | Clean native build and `GOOS=linux GOARCH=arm64` cross-build succeed from pinned tools in documented order; approved flags and module metadata recorded; inspect architecture and runtime dynamic dependencies. | S00 |
| A-04 | Exact candidate serves index and every initial embedded asset from an empty working directory without source/dist/Node; each MIME/byte checksum is correct; production CSP works; both compressed resource budgets pass. | S00 smoke, final S04 |
| A-05 | Configuration precedence, invalid URL/listener/port, help/version, signal shutdown and readiness without Health behave as documented; no startup provider dependency; no unsupported runtime tools required. | S02 |
| A-06 | Debian install/start/upgrade/remove review and sandbox checks; dedicated non-root service, valid unit, no Health startup dependency; installed binary checksum equals candidate; metadata/checksums complete. | S04 |

Race tests may run on a suitable Linux amd64 runner with its required toolchain; a successful ARM64 cross-build is not a race-detector run. Native/cross compilation and runtime dependency checks are separate gates. Review commands/build recipes for hidden local assets and unpinned downloads. Development/build Internet access is allowed by scope; runtime Internet dependency is not.

## Class B — contracts, integration and resilience

Use a Board-owned controllable fake Health HTTP server and an injected elapsed/wall clock. Do not sleep for 30 real seconds in boundary tests. Control transport reads, delayed headers, partial bodies, clock advancement and cancellation with barriers/channels. Use real sockets where HTTP parser/cancellation behavior matters. All must assert response headers, absence of leaked details and unchanged last success after failure where applicable.

| ID | Planned cases and PASS oracle | First responsible sprint |
|---|---|---|
| B-01 | Complete and partial valid schema 1.0; complete→partial atomically replaces values including nulls; zero/false preserved; exact hostname/time; no invented fields. Partial remains Health success and Board 200 available/current. | S01 validation; S02 envelope |
| B-02 | Each issue code; exact path order/associations; null-leaf load/memory/root groups; whole network null; interface-local state/RX/TX nulls; unaffected interfaces preserved; all four firmware null flags/same code+message, temperature independent; true flags with empty issues; valid empty interface array with other useful data; reject all-unavailable/no useful data. | S01 |
| B-03 | Round-trip exact `9007199254740993` and `18446744073709551615` through Health bytes → Board object bytes → browser bigint → exact display. Check 0, safe-integer boundary, uint64 max+1, negative, fractional, exponent and quoted integers; invalid cases follow ratified token policy. | S01, browser display S03 |
| B-04 | Unknown schema string; missing/non-string version; invalid/trailing JSON; wrong types/ranges; duplicate decoded members; unknown members accepted without semantics; malformed UTF-8; every required key omitted in turn versus permitted explicit null; forbidden object nulls; `issues: null`; internal_error issue code; missing/duplicate/reordered issues; unsorted/duplicate/oversized interfaces and mismatched firmware issues; invalid cases rejected with exact classification. Unknown members alone are accepted; their contents still obey global limits and Unicode/duplicate rules. | S01 |
| B-05 | 64 KiB accepted valid boundary, limit+1 rejected, oversized Content-Length and chunked body bounded; browser 80 KiB limit; nesting limit. HTTP 500/503/204/redirect → upstream_error; malformed 200 → invalid_response; wrong media/encoding, truncated 200 and unparseable status cover mapping. No raw errors/body text leak. | S02 |
| B-06 | Refused connection versus deadline expiry; DNS/transport failure; delayed headers and stalled/slow-stream body reach same 1 s global timeout; clock/read cancellation bounded. Assert exactly one outbound attempt, including reused-connection failure, no redirect follow/retry and no call before a browser demand. | S02 |
| B-07 | Never-success, success, failure at 29.999 s / 30 s / 30.001 s, expiration during a request, repeated failure, recovery and process restart. Correct current/stale/none, unchanged success time on error, reset on partial success, snapshot null after expiry. Advance wall clock both directions while elapsed age controls expiry. | S02 |
| B-08 | Simultaneous callers and data races; ratified P-04 sharing only active flight; cancellation of one/all waiters; orphan ends by deadline; no completed-flight reuse; late completion ignored; deterministic sequence/commit ordering and immutable cache copies. Next nonoverlapping demand starts new attempt. | S02 |
| B-09 | Board routing/method/query/body/error precedence, HEAD wire behavior, no SPA fallback for API errors, no-store, MIME/security headers and cache separation; Health failure produces Board 200; Board's own defect produces Board HTTP error. No idle autonomous upstream traffic; logs contain transitions, no snapshots/poll-success spam. | S02, security/assets S04 |
| B-10 | Frontend initial/nominal/partial/stale/none/recovered/Board-unreachable states; immediate load, 5 s cadence, no overlapping requests/retry; abort and late-response discard. Advance clock beyond 30 s with Board down and assert values disappear without any server response. Test header age/RTT, first-load stale, missing age, browser/Pi clock skew, local wall jumps, visibility/freeze/resume and delayed callbacks. | S03 |
| B-11 | CPU 0.5037 displays ~0.5%; bigint counters remain exact and labeled cumulative; null/zero-total percentages; false/null distinction; active/since-boot labels; issue rendering as text (including markup-like strings). 360 px no horizontal overflow, long fields, keyboard/focus, zoom, non-color state cues, no arbitrary score/threshold. | S03 |

Reference upstream tests (read as evidence, not run here): `TestCoordinatorMapsExpectedReasonsInFrozenOrder`, `TestCoordinatorMapsInterfaceLocalFailure`, `TestEncodePartialSnapshotPreservesIssueOrderAndNulls`, `TestEncodeRejectsInvalidSnapshots`, `TestProviderOutcomeMapping`, plus exact-integer tests. Their pinned files and verified facts are in [Health evidence](health-baseline.md). Create Board-owned fixture expectations; do not import upstream test helpers/internal packages.

## Proposed reproducible performance protocol

Use an actual **Raspberry Pi 3B+**, Raspberry Pi OS 64-bit / Trixie, Linux ARM64 for Class C, real Health and the exact candidate binary. Record Pi cooling/power/storage, governor, ambient conditions, other services, OS/kernel, Health identity/configuration and LAN connection. Keep the environment consistent; do not remove failing samples because the machine was busy. If external interference invalidates a run, preserve it and rerun with the reason recorded. Health's hardware acceptance is a separate project obligation.

The operator's measurement tools may read process/kernel counters; that is an external validation harness, not Board system-observation functionality. Board itself does not gain `/proc` collection. Clients/load generators run on another trusted LAN machine so their CPU/RSS do not contaminate Board's process figures.

### Process scope and phases

Collect Board **process** RSS (not virtual size, not combined with Health/browser/journald) every 1 s; use a process identity/start-time check to detect restarts. Stable RSS PASS: maximum observed sample after warmup ≤50 MiB, record median/p95/max and start/end trend. This operationalizes “stable” conservatively; do not hide a growing process behind a mean. Count all Board threads' user+system CPU seconds; mean percent = `100 * CPU_seconds_delta / elapsed_seconds`, one-core denominator. Record 1 s samples and full-phase mean. Exclude only the stated warmup, never failure-recovery sections from their own results.

| Phase | Warmup / measured interval | Gate |
|---|---|---|
| No clients, Health running | 2 min / 10 min | CPU mean ≤1%, RSS ≤50 MiB; upstream request count remains zero |
| One client, exactly one overview start every 5 s | 2 min / 10 min (120 scheduled starts) | CPU mean ≤2%, RSS ≤50 MiB, nominal p95 ≤500 ms; no overlapping client requests |
| Ten clients, staggered by 0.5 s at 5 s each | 2 min / 10 min (1,200 scheduled starts) | Stable process, no Board errors/skips under nominal service, valid overviews, RSS ≤50 MiB, proposed nominal latency gate per client and aggregate |
| Ten clients, synchronized 5 s starts | 2 min / 10 min (1,200 starts) | Same capacity checks, sharing/order behavior if P-04 adopted; no unbounded queues |

For ten clients, record CPU mean/p95 and report it without inventing a new approved CPU limit. Applying the established nominal latency gate in both load phases is a proposed acceptance interpretation. Any skipped scheduled start under nominal conditions fails the capacity run; record failure/timeout counts, not just successful latency. Fake-server load tests supplement, never replace, real-Health Class C phases.

### Latency and lifecycle

Use a monotonic client timer from request dispatch to complete response body receipt/validation; record every sample, bytes/status and state. Sort N samples; p95 is the nearest-rank value at index `ceil(0.95*N)` (1-based), with no interpolation. Report N, min/median/p95/max and failure count. Do not exclude errors/timeouts from reports to improve percentiles. Transport readiness probes and performance timers are external harness traffic explicitly recorded as such.

- Nominal overview: use the 120 one-client samples and each ten-client run; real Health must return valid nominal responses. If Health is failing, report that scenario separately and rerun nominal measurement with the reason recorded.
- Board overhead: collect ≥200 requests at 5 s cadence after 2 min warmup. Proposed timing seam in the **same candidate binary**, opt-in diagnostics only, records handler start/end and upstream network start/complete-body-read (validation excluded from the upstream interval). Calculate handler duration minus upstream network wait; include decode/validation/serialization in Board overhead. Use a single client without shared flight for unambiguous attribution. p95 ≤50 ms. Measure real Health and a controlled responder; preserve both labels. Diagnostic overhead is included, so the estimate is conservative. Review this timing seam in S02; do not rebuild an “instrumented equivalent” after candidate acceptance or enable per-poll logs by default.
- Immediately detectable refusal: stop Health/listener, verify loopback connection refusal rather than packet drop; ≥200 requests at 5 s after 20 excluded warmup attempts, p95 ≤250 ms for complete Board response, HTTP 200 with `connection_failed`. Keep exact candidate unchanged.
- True timeout: controlled server accepts and withholds headers/body (also test unreachable packet-drop path where safe). At least 30 attempts per timeout mode; assert 1 s configured global deadline and cancellation. A proposed integration scheduling tolerance ≤100 ms beyond the deadline is recorded separately; it does not redefine the 1 s timeout or apply the 250 ms refusal gate. Fake clock tests assert the exact deadline without scheduler tolerance.
- Readiness: 20 fresh starts with Health stopped; measure process spawn to successful `/` and required embedded asset response, probing every 20 ms; every trial ≤2 s. Also run with Health present. No overview success is required for readiness.
- Shutdown: 20 SIGTERM trials, including idle and in-flight 1 s Health timeout, from signal to process exit; every trial ≤5 s, no orphan work. Repeat unit/service signal smoke via systemd.

### Resource-size protocol

KiB means 1,024 bytes. Build production/minified assets from locked sources. Proposed compression: gzip level 9, deterministic header (mtime 0/no source filename); report compressor/version. Count each initial URL once on an empty browser cache, no service worker. JS+CSS is the sum of compressed bytes of all eagerly needed JS/CSS, including transitive imports, modulepreload and inline content attributable to these types. Do not game the budget with lazy-loading essential first-page content.

All-initial resources additionally include index HTML, icons/images/fonts if any and the first overview response; use the complete documentation fixture and report its byte size (API may remain uncompressed, then count identity bytes). Include any resource loaded before the initial overview is usable. Exclude HTTP/TLS headers and later periodic overview responses. Record an additional stress size for the maximum accepted snapshot, without asserting all valid snapshots equal the nominal fixture. The inclusion of the first API response is a conservative proposed interpretation of the established 500 KiB budget.

Validate the deployed candidate's actual `Accept-Encoding: gzip` response inventory, not only a directory total. Verify identity fallback and `Vary`; compare served payload checksums with build inventory. JS+CSS ≤250 KiB and all initial resources ≤500 KiB. Record uncompressed sizes too. No external dependency may be excluded simply because it came from a CDN; runtime remote content is forbidden outright.

## Class C — operator acceptance on Pi 3B+

Every row is currently **NOT RUN**. Preserve raw samples/screenshots/log extracts in a candidate-specific future evidence directory; no fabricated result files are created now.

| ID | Future operator procedure / PASS evidence | Sprint |
|---|---|---|
| C-01 | Verify downloaded binary/package checksum and installed binary hash, architecture, Trixie/ARM64, unprivileged systemd identity and embedded assets. Start with real Health stopped: UI/readiness gate passes, overview is unavailable/none with null success; help/version identify candidate. | S04 |
| C-02 | Execute idle, one-client and both ten-client phases with real Health using protocol above; retain CPU/RSS/latency/request-count series and all computed gates. | S04 |
| C-03 | Obtain real data, stop Health, observe Board stays reachable, stale and age then no metrics after >30 s with last success retained; restart Health and observe automatic recovery. Restart Board while Health stopped and confirm lost cache. Stop Board and confirm browser independently ages out data; restart and recover. Partial modes additionally use controlled fault fixtures, labeled synthetic. | S04 |
| C-04 | Run readiness/shutdown/refusal/true-timeout/overhead procedures with identified real or controlled provider per case; meet each budget, no timeout/refusal confusion. | S04 |
| C-05 | Mobile/tablet/desktop UI at documented widths and ≥360 px; long fields/counters, keyboard and non-color cues. Disconnect WAN while retaining LAN and Health: load with cold cache and continue refresh; no remote fetches or fonts. Existing `.local` resolution is tested only if available; access by IP must work. Verify CSP/nosniff/referrer policy. | S04 |
| C-06 | Journal contains startup/shutdown and significant failure/recovery transitions without snapshot/per-poll success logs. Re-check hashes after all runs; package install/upgrade does not replace the accepted binary with different bytes or require live Health. Complete evidence ledger and release review. | S04 |

The operator may need to stop/start test services during **future** acceptance; this document does not execute or authorize deployment in the present documentation task.

## Requirement traceability

Each row maps an established requirement to planned verification and the sprint responsible for delivery. Detailed case matrices above define the oracles. A spanning row does not mark all its delivery sprints complete. The S00 evidence mapping below records the completed foundation scope separately.

| Requirement | Planned verification | Delivery sprint(s) |
|---|---|---|
| JPB-001 | B-09, B-11 route and single-provider review | S02, S03 |
| JPB-002 | B-01, B-11, C-05 metric inventory | S01, S03 |
| JPB-003 | A-01 dependency/source boundary review | S00–S04 |
| JPB-004 | A-01 scope review, B-09/B-11 absent extra actions/routes | S00–S04 |
| JPB-005 | A-01–A-04 stack/build/asset review | S00 |
| JPB-006 | A-05, B-09, C-01/C-05 config and same origin | S00 design, S02 |
| JPB-007 | A-05, B-10, C-01/C-04 | S02, S03 |
| JPB-008 | B-06/B-08/B-09, C-02/C-04 | S02 |
| JPB-009 | B-01/B-07/B-08, C-03 | S02 |
| JPB-010 | B-07, C-03 | S02 |
| JPB-011 | A-01 race, B-08 | S02 |
| JPB-012 | B-05/B-09, C-03 | S02 |
| JPB-013 | B-07/B-10, C-03 | S02, S03 |
| JPB-014 | B-01/B-04, pinned evidence review | S01 |
| JPB-015 | B-03/B-11 | S01, S03 |
| JPB-016 | B-02/B-04 | S01 |
| JPB-017 | B-10, C-03 | S03 |
| JPB-018 | B-10/B-11, C-05 | S03 |
| JPB-019 | B-03/B-11, C-05 | S03 |
| JPB-020 | A-04, B-06, C-01/C-02/C-04 | S04, seams S02 |
| JPB-021 | A-04, B-09/B-11, C-05 | S02–S04 |
| JPB-022 | B-09, C-06 | S02 |
| JPB-023 | A-03/A-05/A-06, C-01/C-06 | S00, S04 |
| JPB-024 | A-01–A-06, B-01–B-11, C-01–C-06; checksum/release review | S04 |
| JPB-025 | Historic documentary review, plus S00 diff/boundary review and dated CI/exit evidence; S00-10 closed; S01 authorized in-memory contract evidence and scope review, no S02–S04 execution | Documentary foundation, S00, S01 |

Future test implementations should cite these IDs in case names/comments or a manifest so evidence can be audited without guessing. Keep this mapping current if a sprint or proposal changes.

### Completed S00 foundation trace

Evidence IDs refer to [S00 evidence](s00-evidence.md). CI-01–CI-05 belong only to run 35527458889 and its exact tested commit; L-xx identify the separately preserved local runs. This table adds execution status without changing requirement scope.

| Requirements / planned gates | Responsible tasks | Executed evidence / status | Remaining scope |
|---|---|---|---|
| JPB-003/JPB-004; A-01 boundary review | S00-03, S00-09 | Local source review, L-12; no Health imports, collectors or functional scope added. PASS for S00. | Repeat boundary review as S01–S04 introduce code. |
| JPB-005; A-01–A-03, A-04 shell subset | S00-02–S00-06, S00-08–S00-09 | L-01–L-11 and CI-01–CI-05: exact tools, clean builds, source/race tests, embed smoke and preliminary sizes. PASS for foundation. | Final dashboard, interactive/browser/mobile checks and production resource acceptance NOT RUN. |
| JPB-006; configuration design | S00-01, S00-07 | Approved defaults and implemented minimal flags recorded in foundation choices; L-03/L-05 and CI-02/CI-04 cover shell behavior only. DONE for S00. | A-05 full production configuration and B-09 functional API NOT RUN (S02). |
| JPB-023; A-03 standalone ARM64 build subset | S00-05, S00-09 | L-04/L-09/L-10 and CI-03: cross-build/ELF inspection. PASS for build only. | ARM64 hardware execution, A-05/A-06 packaging/lifecycle and C-01/C-06 NOT RUN. |
| JPB-025; sprint evidence and scope | S00-08, S00-10 | Verified foundation/race success, exact commit/image/provenance and documentary exit review. DONE; S00 closed on 2026-09-20. | S01 was ready/not started at S00 closure; later S01 evidence is separate. No release/deployment. |

JPB-024's release gate remains unsatisfied: S00/S01 subsets do not accept full Class A/B or any C-01–C-06 gate. Interactive visual verification, complete browser-family validation and mobile visual acceptance remain **NOT RUN**, as does Raspberry Pi physical acceptance.


### Completed S01 contract trace (local and verified CI evidence)

| Requirement / gate subset | Responsible S01 task | Actual verification | Deferred / NOT RUN |
|---|---|---|---|
| JPB-014; B-01 snapshot validity | S01-01/S01-02/S01-03 | Shared Go and TS corpus: complete/partial, hostname/timestamp, domain and presence checks. Local PASS; S01-CI-02 PASS. | Provider HTTP status and complete→partial cache replacement: S02. |
| JPB-016; B-02 issue semantics | S01-02–S01-04 | All issue families/codes, exact association/order, firmware equality, useful observations and interface-local failure. Local PASS; S01-CI-02 PASS. | Live Health HTTP integration: S02. |
| JPB-015; B-03 exact-number subset | S01-02–S01-04 | uint64/bigint extremes, independent raw copy, Go object serialization actually parsed by TS, exact decimal primitive. Local PASS; S01-CI-02 PASS. | Browser-family execution and dashboard display: S03; production overview transport: S02. |
| JPB-014/JPB-016; B-04 consumption policy | S01-01–S01-04 | Shared invalid/valid corpus, all required-field omissions, approved unknown/duplicate/Unicode/size/depth/version rules, seeded differential mutations, bounded Go fuzzing. Local PASS; S01-CI-02/S01-CI-03 PASS. | Actual browser-family execution NOT RUN. |
| JPB-003/JPB-004/JPB-005/JPB-023; A-01–A-04 foundation regression | S01-05 | Locked install, format/types/lint/test/vet/modules, native/cross build, race, embed smoke and separate parser/shell size reports. Local PASS; S01-CI-01–S01-CI-05 PASS. | Final UI budgets, production configuration/package and all physical checks. |
| JPB-025; scope and honest status | S01-05 | Preserved local working-tree fingerprint; verified exact hosted SHA, attempt, job scope/environment and closure. S00 evidence retained separately, no Health modifications or S02–S04 behavior. DONE. | S02 preparation awaits a new task; no implementation/publication/deployment authorized. |

Actual test names/commands/counts, original fixture provenance and limitations live in [S01 evidence](s01-evidence.md) and [shared fixtures](../testdata/README.md). S01's test-only envelope is not `/api/v1/overview`; its Node execution is not browser acceptance. All B-05–B-11 and C-01–C-06 remain NOT RUN.
