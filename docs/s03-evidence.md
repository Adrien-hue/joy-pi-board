# Sprint 03 implementation and evidence

Historical implementation ledger: 2026-09-25. The [2026-09-26 scheduling correction](#2026-09-26--early-cadence-callback-correction) below records the first S03 CI failure and its targeted fix; earlier measurements/manifests remain historical. **Current status (2026-09-27): implemented, CI verified; closure pending operator verification.** See the [corrected CI/operator addendum](#2026-09-27--verified-corrected-ci-and-operator-status) and [remaining checklist](#remaining-acceptance-checklist). That dated assignment did not authorize S04, commit, push, tag, release or deployment. S00/S01/S02 retain their separate historical evidence; their CI does not validate this implementation.

## 2026-10-18 transfer status

The [approved deferral](decisions.md#d-2026-09-28--s03-operator-deferral-and-s04-preparation) moves remaining S03 operator checks to S04 candidate acceptance. The subsequent [S04 implementation instruction](decisions.md#d-2026-10-18--s04-software-implementation) permits software work despite these gaps. S03 is **not closed or fully validated**. Every unexecuted row below remains NOT RUN; historical wording prohibiting the start of S04 is superseded, not the acceptance oracle. The [S04 transferred matrix](sprint-04.md#transferred-s03-operator-checks) retains client/device/reader/native zoom/BFCache dependencies and candidate-dependent replay. A Pi origin cannot substitute for these clients. No historical result/artifact is promoted to the new candidate.

## Source and authorization

Work began in the correct Board repository with a clean tree at `7d9dc989dbbb31409450d5f5ff394ccffb5e43c0` (`docs: prepare Sprint 03 dashboard and browser resilience`). That commit already contains the prior S02 closure and S03 preparation. There were no pre-existing uncommitted changes to overwrite. S02 implementation remains `a11f07c85ff4071a13ca3ae9227920e9f1c7dda4`; the later documentary commit is not attributed its CI results. At the 2026-09-25 implementation delivery, changes were uncommitted, identified by the accompanying input manifest, not by inventing a new SHA. Health remains read-only at baseline `be7a0d824f62b94c842d8e5326110b1852c5a0bc`; no Go server behavior, Health contract, shared fixture bytes or toolchain version changed.

Approval comes from the user's S03 implementation instruction (resumed 2026-09-25), not from the 2026-09-23 preparation: S03-R01/R02/R04 and the S03 P-06 presentation/accessibility/preliminary-size subset. P-02/P-04/P-05 retain their established scope. P-06 physical/packaging/compressed-delivery protocols remain deferred to S04. Server monotonic time can pause during suspension; no browser wall-time subtraction replaces the server clock.

## Delivered behavior and test ownership

- S03-01: bounded streaming overview acquisition; strict Board wrapper, canonical age and state invariants; guarded lossless envelope profile; order-independent source locator; unchanged parseSnapshot authority and deeply frozen known-field projection. Unknown Health members remain allowed. No uint64 passes through Number or stringify.
- S03-02/03: one page-local controller, initial microtask, fixed 5-second slots, 2-second full-read/validation deadline with post-synchronous check; generations, held cancellation admission, no retry/catch-up; lifecycle pause/mask and next-slot resume; paired elapsed segments and independent inclusive-30-second expiry. HTTP/network/contract/provider failures stay distinct.
- S03-04/05: English responsive overview for every prescribed metric, exact decimal counters, bigint binary summaries/ratios/uptime, null/false/zero distinctions, contextual text-only issues and three timestamps. Native details preserve focus/open state on partial refresh; expired values are removed. No invented observations, rates, scores, thresholds, hardware model or service actions.
- S03-06/07: controlled Health HTTP → real Board binary → production browser; separate malformed-envelope interception and development-only StrictMode harness; locked Playwright engines, regression scripts, preliminary dashboard measurement and a new browser CI job.
- S03-08: this ledger, provenance, screenshots and operator checklist. Full closure is explicitly incomplete.

Authoritative test data remains `testdata/s01/cases.json`: 264 cases, 76 valid. Wrappers are built around original bytes; only valid fixture outer whitespace is trimmed when deriving the embedded object (otherwise it counts as Board wrapper overhead). Maximum-markup and long-label/message cases are derived in test-only scripts, without numeric parsing/re-encoding. No fixture or development harness is imported by the production entry.

## Local commands and observed results

The pinned Windows distribution was prepended to PATH for each command. Tool check observed Go **1.27.1**, Node **24.21.0**, npm **11.19.0**; tasks set GOTOOLCHAIN=local, GOENV/GOWORK off. The npm upgrade notice was ignored. Existing React/TypeScript/Vite/Vitest/lossless-json and action pins are unchanged.

| Command / lane | Observed result and scope |
|---|---|
| `npm --prefix web ci --offline=false --cache .cache/npm-s03-fresh` | PASS: locked install, 181 packages, fresh named cache; preinstall exact-tool check passed. |
| `npm --prefix web run check` | PASS: format, strict types, lint, 584 Vitest tests in five files, two fixture-integrity tests, frontend build before Go, gofmt/module tidy-diff/verify/vet and all first-party Go tests. |
| Final test-only refinements: `npm --prefix web run typecheck`, `npm --prefix web test`, `format:check`, `lint` | PASS: 585 tests after normalizing the split-UTF-8 example and adding explicit backward-wall coverage; production artifacts unchanged. |
| Check's contract stage | PASS: 264 common cases, 76 Go-serialized valid outputs parsed in TS; 512 mutations with seed 20260920. Existing S02 real HTTP → Node/TS: 266 attempts, 78 valid/188 invalid, maximum envelope 65,738 bytes/depth 34. This remains a distinct Node lane. |
| `npm --prefix web run fuzz` | PASS: Go FuzzDecode, 10-second requested budget, two workers; 134,620 executions, 8 new interesting inputs, process duration 11.149 s. Workstation evidence, not Pi performance. |
| `npm --prefix web run build` / `cross` | PASS: native Windows amd64 and Linux ARM64, frontend built first, standalone binaries. Final identities are in the artifact inventory. |
| `go build -trimpath -buildvcs=false -o out/joy-pi-board-linux-amd64 ./cmd/joy-pi-board` with GOOS=linux/GOARCH=amd64/CGO_ENABLED=0/GOTOOLCHAIN=local | PASS: Linux binary used by the engine tests, same compiled frontend. |
| `npm --prefix web run smoke` | PASS: copied binary, empty working directory and PATH, byte-identical embedded HTML/assets; unavailable provider overview HTTP 200; zero provider calls for startup/assets, one overview connection, no additional connection during 6-second idle. This raw HTTP smoke does not execute JS. |
| `npm --prefix web run missing-dist` | PASS: isolated source without dist fails clearly at go:embed. |
| Linux `go test -race -count=1 ./cmd/... ./internal/... ./web` | PASS in the pinned Go container below, with CGO enabled and source mounted read-only; all applicable packages, including Linux lifecycle/signal tests. |
| `npm --prefix web run sizes` | PASS for preliminary S03 budgets; actual dashboard/parser and first overview measured separately from server identity transfer. See inventory. |
| Browser suite | Final engine results and reviewed screenshots are recorded below; unavailable branded/floor/device lanes are not inferred from it. |

No bit-for-bit cross-platform reproducibility is claimed. Go race checks exercise the unchanged server packages; the production-browser and standalone smoke checks bind the delivered assets to the recorded binaries. Documentation-only finishing changes do not alter those artifacts.

## Browser tooling, environment and proof lanes

Selected **@playwright/test 1.63.0**, Apache-2.0, Node >=20; its exact playwright/playwright-core dependencies and optional platform dependency are locked by npm. [Primary installation documentation](https://playwright.dev/docs/intro) supports Node 24 and Ubuntu 24.04; [browser documentation](https://playwright.dev/docs/browsers) distinguishes its engines from branded browsers. The [reviewed revisions](../web/playwright-browsers.json) are compared with installed package metadata by the browser command.

Local supported environment: Docker Desktop / WSL2 Linux x86_64, Ubuntu **24.04.4 LTS**, kernel `6.18.33.2-microsoft-standard-WSL2`. Image `mcr.microsoft.com/playwright@sha256:eff16c30e6f3f4af0a03fa4b706120d5e9b0891c344a27d64559aff5900a4a27` (v1.63.0-noble). Its bundled Node reported 24.20.0 and was **not used for test execution**. The official Node 24.21.0 Linux tarball SHA-256 `fd8e59d5a511510f6a298afb548f18c7d2b1be404d8b4a27d94fbe49f56cb2d6` matched SHASUMS256.txt; the test command explicitly invokes its node executable. Test annotations record actual runtime versions and binary hashes.

Race environment: `golang@sha256:69a7b9788769bec032d238959b61854e9ae87f57be9029ec04e9885fabf99195`, native Linux amd64, network disabled, Go 1.27.1/CGO enabled. This is not a hosted image or Raspberry Pi.

Local browser invocation (after frontend/native Linux build and `node scripts/build-browser-harness.mjs`):

```text
docker run --rm --ipc=host -v C:/src/joy-pi-board:/work -w /work/web -e S03_BINARY=/work/out/joy-pi-board-linux-amd64 mcr.microsoft.com/playwright@sha256:eff16c30e6f3f4af0a03fa4b706120d5e9b0891c344a27d64559aff5900a4a27 /work/.tools/node-v24.21.0-linux-x64/bin/node node_modules/playwright/cli.js test
```

The suite uses real sockets for production scenarios: nominal, partial replacement, upstream 503/stale/recovery, maximum markup, all 76 valid common snapshots, 64 interfaces, long text, Health absent before first success and actual Board termination followed by independent expiry. It also checks local resources/CSP and exact visible uint64 digits. Invalid Board HTTP/body tests use interception and are labeled separately. The unembedded development StrictMode entry verifies two effect setups and one cleanup, with one surviving initial request and the next request at its ordinary slot; it is not the production build.

Screenshots cover 360/390/768/1280 CSS px, keyboard disclosures, partial/stale/none/HTTP/network failures, exact counters, long data and 200% **text enlargement**. Text enlargement is not falsely called native browser zoom. Synthetic pagehide/pageshow tests are not proof of actual BFCache residency or device suspension. Linux WebKit is not Safari/iOS.

## Final local artifact and browser ledger

The [87-file input manifest](evidence/s03-inputs.json) has aggregate SHA-256 `57f2da4ee0ec963b625efd956ef60af50c07fc3f4fcf4214d70ec5b04ca94e6f`, computed from sorted `sha256 + two spaces + path + LF` lines. It includes production/test/build/CI inputs and excludes documentation, generated dist and downloaded tools. The final test-only Unicode/backward-wall edits were checked after building; they do not enter production assets. [Build identities](evidence/s03-builds.json), [embedded standalone smoke](evidence/s03-standalone-smoke.json), [S01 interop regression](evidence/s03-s01-interop.json) and [S02 HTTP interop regression](evidence/s03-http-interop.json) are local evidence, not hosted artifacts.

| Binary | Bytes | SHA-256 |
|---|---:|---|
| Windows amd64 | 11,482,624 | `91e542d4fddba40b76a82b33f618cbdb2e8c4f88aaa207452223e9a9983c4eb8` |
| Linux amd64, used by browsers | 11,340,064 | `a26c6442f4638f423608500177c4f0fd1a0f7828af7b21dea84ba5e7d843e42a` |
| Linux ARM64, cross-built only | 10,561,388 | `c7e05e1d6357e00a1556fbdf66b79017619f2001bd5efbcb15896d3a837b382d` |

[Final Playwright report](evidence/s03-browser-report.json): started **2026-09-25T09:41:52.000Z**, 266,561.714 ms total, **27 expected passes, 0 skipped, 0 unexpected, 0 flaky**, retries disabled. Every test annotation records the Linux binary SHA above and Node 24.21.0. No hosted run is represented by this report.

| Actual engine version | Reviewed package revision | Local result |
|---|---|---|
| Chromium 153.0.8010.12 | 1243 | 9/9 PASS |
| Firefox 155.0 | 1543 | 9/9 PASS |
| WebKit 26.6 | 2359, Linux | 9/9 PASS |

Each engine replays 76 valid common snapshots through controlled Health HTTP → real embedded Board HTTP → production UI, independently of the Node contract suite. The other cases cover exact digits, partial replacement, Health error/recovery/max markup, actual Board termination and independent expiration, synthetic lifecycle, malformed/HTTP-error interception, 64 interfaces/long issue text, actual development StrictMode replay, and no prior Health success. The report contains per-case durations; these are test durations on a workstation, not response-latency or Raspberry Pi measurements.

The [visual inventory](evidence/s03-visual.json) identifies all **45 screenshots** and the **14 images actually reviewed and retained**, with binary identity, case, engine, pixels, SHA-256 and review notes. DPR is 1; the viewport matrix is 360/390/768/1280 CSS px. Representative retained images:

- [Desktop exact values and keyboard focus](evidence/s03-visual/chromium-nominal-1280.png), [360 px](evidence/s03-visual/chromium-nominal-360.png), [390 px Firefox](evidence/s03-visual/firefox-nominal-390.png), [768 px WebKit](evidence/s03-visual/webkit-nominal-768.png).
- [Partial replacement with stable focus](evidence/s03-visual/chromium-partial.png), [Health stale](evidence/s03-visual/webkit-stale.png), [Board unreachable but eligible stale data](evidence/s03-visual/firefox-board-unreachable-stale.png), [expired without another server response](evidence/s03-visual/firefox-expired-without-server.png), [never-success](evidence/s03-visual/webkit-none-before-first-success.png), [separate Board HTTP failure](evidence/s03-visual/chromium-board-http.png).
- [200% text reflow](evidence/s03-visual/webkit-text-200-percent.png), [literal long issue text](evidence/s03-visual/firefox-long-text-issue.png), [64-interface structure](evidence/s03-visual/chromium-long-64-interfaces.png), [recovered maximum payload](evidence/s03-visual/webkit-recovered-maximum.png).

Agent review found the final samples readable, without cropped metric sections or horizontal overflow; exact large counters and state labels are visible. The extremely tall 64-interface image was reviewed for overall structure at reduced scale, supplemented by DOM/count/overflow assertions, not a claim of reading every pixel. [Calculated solid-color contrast pairs](evidence/s03-contrast.json) have a minimum ratio **5.179:1** across the enumerated text/focus pairs (WCAG relative luminance); this does not replace screen-reader, native zoom, actual-device or complete accessibility acceptance.

## Preliminary dashboard cost

[Measured inventory](evidence/s03-dashboard-sizes.json) uses Node 24.21.0 / zlib `1.3.2.1-motley-8002e91`, independent gzip level 9 for each response body, and byte-compares fetched resources against dist. All actual HTTP deliveries were **identity (Content-Encoding absent)**. No S04 compression is implemented.

| Measured scope | Raw / HTTP identity bytes | Offline gzip-9 bytes |
|---|---:|---:|
| Production JS, parser/validator included | 252,819 | 78,138 |
| CSS | 3,489 | 1,336 |
| JS + CSS | 256,308 | **79,474** |
| HTML + all initial static resources | **256,783** | **79,766** |
| First complete overview | 947 | 495 |
| First partial overview | 2,180 | 621 |
| First maximum-markup overview | 65,724 | 654 |
| Resources + complete overview | 257,730 | **80,261** |
| Resources + partial overview | 258,963 | **80,387** |
| Resources + maximum-markup overview | 322,507 | **80,420** |

The approved preliminary ≤250 KiB JS+CSS / ≤500 KiB all-initial offline-compressed gates **PASS** for these scenarios. Identity transfer is reported separately. The maximum markup fixture is exactly 65,536 Health bytes, SHA-256 `49533f8c6d262849aa6b3cb1ff3003a4c90f0a2280612a3ae5a074c437c29b5a`. Overview hashes identify the actual sampled timestamps and are not asserted invariant across runs. S00 shell sizes and S01's 17,499-byte/5,585-byte-gzip standalone parser remain historical/separate measurements, not components to add to this dashboard.

## Task and validation status

| Work item | Local status / remaining gate |
|---|---|
| S03-01–S03-04 | Implemented; bounded consumer, exact parsing, controller/freshness and formatting assertions PASS. S03-T01–T08 implemented deterministic/engine assertions are mapped in the validation document. |
| S03-05 | Dashboard implemented; static, engine keyboard/reflow and reviewed-image subsets PASS. Screen reader, native zoom and actual-device portions of S03-T09/T12 remain NOT RUN; see the dated operator checklist. |
| S03-06 | Three supported automated engines and real HTTP chain PASS (S03-T10). Actual branded/floor/mobile evidence S03-T11 and real lifecycle/device portions of T05/T12 remain NOT RUN. |
| S03-07 | Local regression/build/size checks PASS (T13); corrected hosted CI now verified in the 2026-09-27 addendum. |
| S03-08 | Local provenance/traceability recorded; T14 and sprint closure PENDING until all mandatory operator evidence exists; corrected CI is verified separately below. |

Passing implemented assertions does not waive an unexecuted acceptance lane. The fourteen prepared test IDs retain their full scope in the validation matrix.

## Corrected incidents and limits

1. Initial test fixture path had one excess parent segment; corrected before contract execution. Valid S01 boundary fixtures include outer whitespace: retaining it inside the Board wrapper exceeded its distinct 1,024-byte overhead. Derived wrapper tests trim only valid outer whitespace; original fixture bytes remain unchanged.
2. Initial TypeScript exact-optional props and a React ref-analysis lint error were fixed by explicit optional-value types and effect-owned masking registration. No lint rule was disabled.
3. The first Chromium run was 4 PASS/1 FAIL: the test's partial-state label had acquired an extra encoding character. The actual page text was correct. The test encoding was repaired.
4. The next three-engine run was 17 PASS/4 FAIL: the corpus's zero-value locator matched a hidden exact load disclosure before the visible RX/TX zero; it is now scoped to network counters. WebKit screenshot synchronization injects a `body {}` style tag and triggered the strict CSP. Production CSP is now asserted before screenshot instrumentation; the policy was not weakened and arbitrary CSP failures are not ignored.
5. Review caught undercount potential when local wall/monotonic deltas diverge differently during body reading and parsing. The controller now sums the larger delta of each segment, with an exact-30-second regression test. Expiry never uses Pi wall timestamps.
6. Exact disclosures now retain identity when a value becomes null; keyboard focus/open state remains on the same disclosure, showing Unavailable. At 200% text the load cards were cramped; automatic minimum-width columns now reflow instead of splitting decimal loads across narrow columns.
7. An intermediate 27-test browser run was interrupted when the local Docker engine stopped during the session transition; its partial progress is not a completed PASS report. Final rerun evidence is separate. Failed/partial runs were not hidden by test retries (retries=0).
8. The final source review found lossless-json assigns __proto__ rather than retaining an own property. The narrow locator now additionally enforces decoded wrapper key allowlists, covering literal/escaped prototype keys with null/number/object values and inherited-required-field attempts. Health extras remain tolerated. The corrected assets were rebuilt and all three engines rerun. A split-UTF-8 test string was also normalized to actual accented characters, and a backward-wall case was made explicit.
9. Native Computer Use inventory failed with `Computer Use native pipe is unavailable` (OS error 2). Browser-provider inventory returned no surfaces; opening Chrome reported `Browser is not available: chrome`. Installed-file inventory alone found Chrome 154.0.8037.57, Edge 153.0.4234.48, Firefox 156.0.1 on Windows 10.0.19045; none is claimed as a launched branded-browser validation. No desktop/security setting was changed.

## Remaining browser and operator procedure

The [remaining acceptance checklist](#remaining-acceptance-checklist) is the current source of status per environment; the older local ledger is historical. These procedures refine the existing S03 operator work, not its acceptance policy. Record **date/time UTC, test IDs, operator, exact source SHA or manifest, binary/asset SHA-256, browser patch/version page, OS/build, device model, viewport/DPR and native zoom** for every run. Name screenshots/traces by browser, case and identity. Never infer PASS from an inventory or accessibility tree.

### OP-1 - Desktop setup and exact data

Use an existing isolated/trusted environment and the pinned Node/Go/npm tools in [foundation commands](s00-foundation.md). Prefer the identified corrected binary; rebuild only if necessary with `npm --prefix web run build`, then record the new checksum and source. The interactive harness below is the existing `scripts/browser-harness.mjs`, not a production server control endpoint. From the repository root, start the pinned Node executable (on this Windows checkout: `.\.tools\node-v24.21.0-win-x64\node.exe --experimental-repl-await`):

```javascript
const h = await import('./scripts/browser-harness.mjs');
let b = await h.startBoard();
console.log(b.origin, b.binary);
```

Open only the printed loopback origin in the actual Chrome/Edge/Firefox/Safari browser and record its built-in version page (Chrome `chrome://version`, Edge `edge://version`, Firefox `about:support`, Safari About Safari and macOS version). Expect current observations, CPU **0.5%** with exact **0.5037%**, visible exact RX/TX **9007199254740993** and **18446744073709551615**, and distinct active/since-boot flags. Compare all ordinary metrics with the authoritative `complete` fixture without Number parse/stringify. Capture the native page, not only its accessibility tree.

Use actual isolated Chrome 111 / Edge 111 / Firefox 115 / macOS Safari 16.4 environments for the same critical parsing, expiry and lifecycle flows; record the actual patch. Do not overwrite daily browsers or spoof user agents. If no environment exists, record NOT RUN and ask the operator to provide an isolated image/device.

### OP-2 - Partial, failure, recovery and local expiry

In the same REPL, use only these existing controls, waiting for observable ordinary refreshes rather than forcing retries:

```javascript
b.state.body = h.fixture('partial-documentary');
b.state.status = 503;
b.state.status = 200; b.state.body = h.markup();
b.state.body = h.fixture('interfaces-64');
b.state.body = h.fixture();
```

Execute each line as a separate phase, after recording the prior result. Partial must replace previous values with Unavailable and scoped issues, never merge old values. Health 503 must leave Board reachable and produce stale then none after the inclusive 30 s eligibility window; last-success remains. Restore 200 plus markup to check recovery, 65,536-byte Health payload and unchanged exact counters. Check 64 interfaces for retained information/reflow. These are test fixtures, never real host observations.

After complete success, execute `await b.kill()`. Keep the page foreground, record the last response/age and a timestamped video or request trace: Board becomes unreachable, Health becomes unknown, and old hostname/metrics/exact disclosures disappear when the conservative age exceeds 30 s without another server response. Do not refresh the page to make them disappear. `await b.stop()` closes the controlled provider too; `b = await h.startBoard()` creates a fresh process/cache and a new printed origin for another case.

A provider `state.status=503` is **not** a Board HTTP 503. The existing intercepted Board-error case in `scripts/browser.spec.mjs` tests that distinct path under Playwright. There is no manual Board-error injection CLI in the harness. For additional branded-browser adversarial HTTP proof, an operator-controlled response-interception setup is needed; document that setup/result rather than inventing a harness flag or substituting the provider failure. No such new tool is introduced by this documentary assignment.

### OP-3 - Keyboard, native zoom, screen reader and visual review

With a working native desktop, use Tab/Shift+Tab and Enter/Space only to reach the skip link and exact-value disclosures. Keep one disclosure focused/open through several normal polls, then switch to partial and confirm focus/open state survives with Unavailable. Record any actual unexpected movement as FAIL. Check visible focus and logical headings/group associations. Review full page at 360/390, 768 and 1280 CSS px, long strings/issues and 64 interfaces; no horizontal overview overflow or loss of exact information.

Set **native browser zoom to 200% using the browser menu**, record its indicator plus viewport/DPR, inspect layout and exact values, then reset to the original zoom. CSS root-font enlargement or device viewport emulation is supplemental only. Capture readable screenshots; an inaccessible capture surface leaves the visual check NOT RUN.

Use an available real screen reader (Narrator/NVDA/VoiceOver) with audible output and an operator who can hear it. Read headings and the affected metric/issue association; induce complete -> partial -> provider failure once. Record spoken announcements, no repeated metric/age announcements over normal polls, keyboard operation and association observations. Do not mark this PASS from DOM/AX structure alone. Stop the reader and restore its original state afterwards.

### OP-4 - Real visibility, lifecycle and BFCache

Keep one Board tab open. Record the browser's Network timeline with preserved logs and `b.state.calls` before/after a real switch to a blank tab for more than two slots. Return: no polling while hidden, metrics masked until a new valid **future 5 s slot**, no immediate catch-up burst or old response restoration. Repeat leaving during a request and actual app background/foreground. Inspect successful and failed refresh cases; retain timestamps and whether the page was actually visible.

For actual back navigation, enable preserved console logs and temporarily observe real events in the test page's DevTools console (instrumentation only, no product edits):
```javascript
for (const type of ['pagehide', 'pageshow', 'visibilitychange', 'freeze', 'resume']) {
  const target = ['visibilitychange', 'freeze', 'resume'].includes(type) ? document : window;
  target.addEventListener(type, e => console.info('S03 lifecycle', type,
    performance.now(), document.visibilityState, e.persisted));
}
```
Navigate to a blank/local page and use the browser Back action. Require evidence of actual `pageshow.persisted=true` for the BFCache-restoration row, plus the same masking/next-slot behavior. A reload, synthetic event or `persisted=false` is not BFCache PASS. DevTools or platform policy may affect eligibility: record the stated reason and use an eligible isolated environment with a remote/preserved event recording method if needed. If no such method is available, keep that row NOT RUN; do not change server cache/security headers.

With the operator physically present, lock/suspend and resume the actual device, then test background/foreground on mobile. Preserve before/after event/network records; old data must remain masked until a new valid response. Do not claim a browser tab switch proves OS suspension. No OS protection needs to be disabled.

### OP-5 - Real Safari/iOS setup

Requires an identified macOS/Safari environment and real iOS device (minimum 16.4 and current-family evidence). `startBoard()` binds loopback and cannot itself serve the LAN; there is no `--lan` harness option. On an operator-controlled LAN test host, reuse `h.fixture()` as exact bytes in a temporary Node HTTP server in the REPL:

```javascript
const { createServer } = await import('node:http');
const state = { body: h.fixture(), status: 200, calls: 0 };
const provider = createServer((_req, res) => {
  state.calls++;
  res.writeHead(state.status, { 'Content-Type': 'application/json' });
  res.end(state.body);
});
await new Promise(done => provider.listen(0, '127.0.0.1', done));
console.log(provider.address().port);
```

In another terminal run the recorded native Board binary with `--listen <test-host-LAN-IP>:8081 --health-url http://127.0.0.1:<printed-port>/v1/snapshot` (replace placeholders with the host's actual address/port). Verify existing trusted LAN access first; do not expose Health, add a public tunnel, change firewall protections or install a service. If networking is unavailable, the operator must supply an already permitted test host. A different host build needs its own checksum; this is software browser testing, not Pi Class C acceptance.

Open that host in actual Safari/iOS, run OP-1 through OP-4 using `state` instead of `b.state` for provider transitions, and inspect touch disclosures, native zoom, background/app switch and back navigation. Retain exact OS/browser/device model and whether BFCache restore actually occurred. To stop Board for expiry use its terminal Ctrl+C; to clean up close the browser test tab, stop Board if still running and `await new Promise(done => provider.close(done))`. Exit the test REPL. Do not substitute Linux WebKit, a mobile viewport or a user-agent string for an Apple device.

### OP-6 - Evidence and cleanup

Retain PASS/FAIL/NOT RUN per environment and S03-T05/T09/T11/T12, screenshots/video or speech observations, event/timing traces and source/binary/fixture identities. Preserve every failed attempt and distinguish tooling failures from product defects. After each local harness run use `await b.stop()` and exit its REPL, close only task-created tabs/windows, reset changed zoom and stop any reader started for the test. Do not erase pre-existing work or personal browser tabs. The corrected hosted CI gate is already satisfied on the exact SHA below; incomplete operator rows still block S03 closure.

## Hosted CI and sprint status

The existing workflow retains immutable action pins and Ubuntu 24.04 labels; foundation/race remain, and a browser job installs exact package-selected revisions, builds the embedded app, runs all three engines and uploads results. It records hosted ImageOS/ImageVersion, OS/kernel and browser revision data; the hosted image label is not claimed immutable. **At the 2026-09-25 local delivery, S03 hosted CI was NOT RUN** for that then-uncommitted state. The subsequent run on commit 12108b60861582cab040b9d5a029279e5768bcea failed as recorded in the addendum below; at the time of that correction, new CI was still required. The 2026-09-27 addendum records its subsequent verified success. An edited workflow is not a green run.

S03's production implementation is delivered, but full validation and closure are **PENDING**, including the mandatory operator targets above. S04 packaging/systemd, server compression, Raspberry Pi/real-Health integration, physical budgets and global A/B/C acceptance are **NOT RUN**. There is no release, deployment or implicit authorization to start S04.

## Final repository review

Reviewed the implementation/configuration/CI diff and untracked S03 additions; no Go server code, shared Health fixture, toolchain pin or S00/S01/S02 evidence file changed. Documentation includes the foundation-command note because sizes now measures the actual dashboard and a browser command has been added. Source manifest verification matched all 87 inputs. Markdown links/anchors, JSON documentation examples, eight-task/fourteen-test-ID traceability and git diff --check were checked after harmonizing current and historical statuses. The final audit covered 23 Markdown files, 321 internal links/anchors and 1 JSON documentation example; no application test is inferred from these checks. A final process inventory found no remaining test Node/Board process (only the unrelated computer-use runtime); no deployed service was created.

## 2026-09-26 — Early cadence callback correction

### Starting point and verified CI incident

This targeted task started with a clean Board working tree at `12108b60861582cab040b9d5a029279e5768bcea`, matching the failed S03 implementation commit. It authorizes only the scheduling fix, regression tests and S03 evidence. No S04 work, Health/server contract change, dependency/toolchain/workflow change, commit, push, release or deployment.

GitHub was read directly through its repository connector on 2026-09-26: [run 36177274215](https://github.com/Adrien-hue/joy-pi-board/actions/runs/36177274215), **attempt 1**, workflow **Board foundation and contracts**, **push/main**, created **2026-09-25T19:03:05Z**, SHA above, conclusion **failure**. Jobs foundation `108210820287` and race `108210820317` concluded success; browser `108210819984` failed its build/browser step. Its log records **26 PASS / 1 FAIL**, specifically Chromium's `S03-T04 actual React development StrictMode effect replay`, expected 2 calls, received 3 at scripts/browser.spec.mjs:328. The 7-second assertion and retries=0 are unchanged.

[CI summary/log excerpts](evidence/s03-scheduling-2026-09-26/ci-summary.json), [original CI browser report](evidence/s03-scheduling-2026-09-26/ci-browser-report.json) and [observed hosted image](evidence/s03-scheduling-2026-09-26/ci-browser-host.txt) preserve these new facts separately from the prior local PASS. The actual image was ubuntu24 **20260920.314.1**, Ubuntu **24.04.5 LTS**, Linux **6.17.0-1022-azure** x86_64. The log's exact-tool check observed Go 1.27.1, Node 24.21.0, npm 11.19.0. Browser revisions remained Chromium 1243 / Firefox 1543 / WebKit 2359. A hosted image label is not immutable.

The downloaded browser artifact `10883565483` matched GitHub's SHA-256 `6716769536f6ca29c9ee9d704718a7d5a22daf9dc08494667a9d4a3ed589aa0d`. Its Chromium trace contains three successful overview requests starting at **19:05:39.508Z**, **19:05:44.506Z**, **19:05:44.520Z**. Starts 2 and 3 are **13.437 ms apart** in the trace's monotonic coordinates. [Trace extraction and archive identities](evidence/s03-scheduling-2026-09-26/ci-trace-summary.json) distinguish observation from inference: this is consistent with a premature callback followed by rearming the same slot after a fast response, but the trace does not contain the controller's epoch/deadline or callback clock readings. The precise early-callback mechanism is established by the deterministic reproduction, not asserted from HTTP timing alone.

### Established defect and correction

The old callback checked only `now < due + 5000`. It admitted an early callback, then recomputed the next timeout for the same still-future deadline. Once the request released admission, that deadline could launch again. The new deterministic test starts at **100.25 ms**, expects **5100.25 ms**, and explicitly delivers the callback at **5100 ms**. Before changing production code, it failed with actual calls `[100.25, 5100]` instead of `[100.25]`. The [red result](evidence/s03-scheduling-2026-09-26/s03-scheduling-red-test.log) is retained, rather than relabelled PASS.

`OverviewController.schedule` now:

- Checks timer identity, lifecycle generation, mounted and visible state before any effect; cancelled/already-consumed callbacks cannot clear or replace a newer timer.
- Checks the unrounded deadline at callback time. If early, retains that deadline and rearms one timeout without calling Fetch.
- Uses positive whole-millisecond waits (ceil, minimum 1 ms) so fractional rearming does not create a zero-delay spin. This is only wake-up quantization; admission still depends on the explicit deadline comparison, with no epsilon.
- Consumes a due slot even when busy and sets the minimum next deadline to that slot +5 seconds. The existing strictly-future epoch calculation discards missed slots. A fast completion cannot reuse the consumed slot.
- Retains one cadence timer, existing active-request cleanup admission, no overlap/retry/catch-up, next-slot foreground resume and StrictMode cleanup.

Six added deterministic cases cover repeated early delivery, exact fractional deadline, slight lateness within the slot, an entire missed slot at its boundary, several missed slots, fast successful completion and stale/cancelled callback replay. Existing tests retain busy/aborted-cleanup, visibility/resume, timeout and StrictMode assertions. The real browser test's exact `toBe(2)` assertion, StrictMode, timeout and Playwright retry policy are unchanged.

### Correction execution ledger

All executions in this correction are recorded, including infrastructure failures and the expected red test. Local tools remain Go **1.27.1**, Node **24.21.0**, npm **11.19.0**. Production changes are confined to controller.ts; controller.test.ts adds regressions.

| Execution | Result / evidence |
|---|---|
| Initial targeted Vitest launch in sandbox | Did not execute tests: Vite child-process spawn EPERM. [Launch log](evidence/s03-scheduling-2026-09-26/s03-scheduling-red.log). |
| Same early-callback regression outside sandbox, old controller | Expected FAIL: 1 failed, 21 filtered/skipped; observed request before its deadline. Red log above. |
| `npm --prefix web test -- src/overview/controller.test.ts`, corrected controller | **23 PASS**. [Targeted log](evidence/s03-scheduling-2026-09-26/s03-scheduling-targeted.log). |
| `npm --prefix web run check` | **PASS**: format, strict types, lint, **591 frontend tests**, two fixture-integrity checks, frontend-before-Go, gofmt/modules/vet/Go tests, 264-case corpus, 76 Go→TS valid outputs, 512 seeded mutations and 266 real S02 HTTP→TS cases. [Full check log](evidence/s03-scheduling-2026-09-26/s03-scheduling-check.log). |
| `npm --prefix web run build` and `cross`; Linux amd64 Go build | **PASS**, same new compiled assets in native Windows, Linux ARM64 and Linux amd64 binaries. [Native log](evidence/s03-scheduling-2026-09-26/s03-scheduling-build.log), [cross log](evidence/s03-scheduling-2026-09-26/s03-scheduling-cross.log). Linux uses GOOS=linux, GOARCH=amd64, CGO_ENABLED=0, GOTOOLCHAIN=local, GOENV/GOWORK=off, `go build -trimpath -buildvcs=false -o out/joy-pi-board-linux-amd64 ./cmd/joy-pi-board`. |
| `node scripts/build-browser-harness.mjs` | **PASS**, rebuilt unembedded development StrictMode entry from corrected controller; its separate hash is inventoried. |
| `npm --prefix web run smoke` and `sizes` | **PASS**: exact embedded assets from empty cwd/PATH, unavailable overview HTTP 200 and one Health attempt, no idle call; updated dashboard cost below. [Smoke](evidence/s03-scheduling-2026-09-26/standalone-smoke.json), [size inventory](evidence/s03-scheduling-2026-09-26/dashboard-sizes.json). |
| First Docker browser command | Did not execute tests: Docker Linux engine pipe absent. Started Docker Desktop hidden, verified engine **29.7.2**, then launched the suite once. This is not a Playwright test retry. |
| Full corrected browser suite after Docker startup | **27 PASS (9 per engine), 0 skipped/unexpected/flaky, retries=0**, one completed invocation. Chromium StrictMode retains exact toBe(2). [Report](evidence/s03-scheduling-2026-09-26/browser-report.json), [complete output](evidence/s03-scheduling-2026-09-26/browser.log). |
| CI artifact retrieval | Connector download succeeded; first local HTTP transfer failed during receive. Fresh signed download outside sandbox succeeded and archive digest matched. No trace evidence is invented from the failed transfer. |

Go race/fuzz and missing-dist were not rerun for this frontend scheduling-only diff; their earlier local/CI results remain dated historical evidence, not new correction results. No Go, embed layout, lockfile, CI policy or schema changed. The full check and real rebuilt-binary browser integration provide the relevant regression lanes.

The corrected browser run began **2026-09-26T19:52:01.644Z** and took **270,783.325 ms**. Observed engines: Chromium **153.0.8010.12**, Firefox **155.0**, WebKit **26.6**, all Linux x64 under Node **24.21.0**, using the same pinned Playwright 1.63.0 Ubuntu image described above. Every test annotation identifies corrected Linux binary SHA-256 4c6f137642478fc8e82177f28715e1e483e46abd40c38d8e48aace99d0bf1c5c. The development-only StrictMode harness was rebuilt too; this case is not misrepresented as production React behavior.

Invocation was the same documented Docker command with S03_BINARY pointing at the newly built Linux amd64 binary and the exact .tools/node-v24.21.0-linux-x64 executable, output retained in browser.log. The production-browser cases exercise the rebuilt embedded assets, while StrictMode uses its separate development entry. No extra browser test run, retry, timeout relaxation or assertion weakening was used to obtain this PASS. These automated engines do not fill any existing branded/minimum-version/iOS/operator gap.

### Corrected source and artifact provenance

[Correction input manifest](evidence/s03-scheduling-2026-09-26/inputs.json): **87 files**, aggregate SHA-256 `75e587ca8dbeadc5f4ec9a25f91de53b8cd22d2d35240321196447f72fbcf700`, same sorted-line algorithm as the original manifest, based on the clean commit above plus the uncommitted fix. [Build inventory](evidence/s03-scheduling-2026-09-26/builds.json) identifies the local corrected artifacts:

| Artifact | SHA-256 |
|---|---|
| Windows amd64 binary | `6bfbbfb96e47d1ec1baf07d0a8e19e1f3d541d774c420848346b1f1ca568dcf7` |
| Linux amd64 browser binary | `4c6f137642478fc8e82177f28715e1e483e46abd40c38d8e48aace99d0bf1c5c` |
| Linux ARM64 cross-build | `867468003ea3a4219d2f8c72051b50e1c3a92774abc7423090678f72c01abb8e` |
| Production JS | `887354a9e865afc3cce33a38e1aadf2d773f70ad69b7dbf46a48cd129250d8f3` |
| Test-only development StrictMode module | `2a58beb2c02d310134ae5ddbc54d671c8db9947ef0cb1adcbd1c975b18e6d7a6` |

JS+CSS is **79,522 bytes offline gzip-9**; initial static resources **79,817 gzip / 256,971 HTTP identity bytes**. Including the maximum-markup first overview gives **80,471 gzip / 322,695 identity bytes**. Preliminary S03 budgets still pass. No server compression or physical performance claim. Earlier hashes and screenshots above remain evidence for their earlier artifacts, not for these corrected builds.

At the 2026-09-26 local correction handoff, corrected hosted CI was **NOT RUN** (superseded by the verified run below). S03 then remained **IMPLEMENTED, LOCALLY VERIFIED WITHIN EXECUTED SCOPE; CLOSURE PENDING**, subject to the new exact-commit CI and all existing required operator/browser-floor/Safari/iOS/native-zoom/screen-reader/lifecycle proofs. This correction neither waives those gates nor starts S04.

Correction finishing checks: format/types/lint passed within the full check; git diff --check is clean. Internal links/anchors and the 87-file input manifest were rechecked. Final diff is limited to the controller, deterministic tests, S03 evidence and its validation trace; browser test, StrictMode entry, dependencies, toolchains, workflows, Go server and historical evidence artifacts are unchanged.

## 2026-09-27 — Verified corrected CI and operator status

**Current status: IMPLEMENTED, CI VERIFIED; CLOSURE PENDING OPERATOR VERIFICATION.** This assignment authorizes evidence consolidation and S03 validation only, not further product changes or S04. It started from a clean working tree. Local HEAD, the GitHub main branch read directly on 2026-09-27, and the tested SHA all equal `24d527a301d66b49d86ad11fbd75393254f02094` (`fix: prevent early and duplicate overview refresh attempts`). Its parent is the failed implementation commit `12108b60861582cab040b9d5a029279e5768bcea`. These new documentary changes are not attributed that CI.

### Corrected hosted CI

[Run 36269739581](https://github.com/Adrien-hue/joy-pi-board/actions/runs/36269739581), attempt **1**, workflow **Board foundation and contracts**, event **push**, branch **main**, started **2026-09-26T20:30:08Z**, last updated **20:35:41Z**, conclusion **success**. Metadata, all job steps, complete job logs and the browser artifact were inspected directly through GitHub. Retained job logs normalize line-end whitespace only; artifact reports/images retain their downloaded bytes. [Run/jobs/artifacts provenance](evidence/s03-operator-2026-09-27/ci-summary.json) records the exact identities.

| Job / ID | Executed evidence / confirmed scope |
|---|---|
| foundation / 108481283244 — PASS | [Log](evidence/s03-operator-2026-09-27/ci-foundation.log): exact tools, locked install, format/types/lint, **591 frontend tests in 5 files**, 2 fixture-integrity tests, frontend-before-embed, Go format/modules/vet/tests; 264 shared corpus cases, 76 Go-to-TypeScript outputs, 512 seeded mutations; **266 real HTTP-to-TypeScript cases / 266 Health attempts**. Bounded Go fuzzing reports 359,795 executions over its observed 11 s completion (requested fuzztime 10 s). Native Linux and Linux ARM64 build, ELF/static check, standalone unavailable-overview smoke, missing-dist rejection and size inventory passed. |
| race / 108481283155 — PASS | [Log](evidence/s03-operator-2026-09-27/ci-race.log): locked install, frontend build before native Linux `go test -race -count=1` with CGO enabled. Package results include command lifecycle, config, eventlog, Health, schema, HTTP API/UI and overview. Packages reported with no test files are not invented test counts. |
| browser / 108481283237 — PASS | [Log](evidence/s03-operator-2026-09-27/ci-browser.log) and [original report](evidence/s03-operator-2026-09-27/ci-browser-report.json): **27 PASS, 9 per engine; 0 skipped/unexpected/flaky**, retries=0. Exact StrictMode `toBe(2)` passed in Chromium, Firefox and WebKit. Embedded production Health HTTP → Board HTTP → browser lanes, uint64/corpus/maximum-markup cases, independent Board-down expiry, keyboard/reflow and CSP/resource checks ran. Intercepted invalid envelopes and the separate development StrictMode entry retain their explicit test-only scope. Synthetic lifecycle and CSS text enlargement are not real BFCache/device/native zoom proof. |

The unchanged workflow and scripts at the tested local HEAD were inspected to map these steps to `tasks.mjs` check/contract/fuzz/race/build/cross/smoke/missing-dist/browser/sizes and `browser.spec.mjs`. The nine browser cases per engine cover real HTTP/reflow/CSP; partial/failure/recovery/maximum markup; Board-down expiry; synthetic lifecycle; intercepted contract/HTTP errors; all 76 valid shared snapshots; long text/64 interfaces; actual development StrictMode replay; and Health failure before first success. This establishes coverage, not a blanket PASS for every manual assertion in S03-T01–T14.

**Configured versus observed:** ubuntu-24.04 was the configured label; logs and [host artifact](evidence/s03-operator-2026-09-27/ci-browser-host.txt) observe image **20260920.314.1**, Ubuntu **24.04.5 LTS**, Linux **6.17.0-1022-azure**, x64, runner **2.337.0**. Exact-tool execution confirms Go **1.27.1**, Node **24.21.0**, npm **11.19.0**; race compiler is GCC **13.3.0**. Playwright **1.63.0** and its revision lock were checked by the browser command. The report observes Chromium **153.0.8010.12**, Firefox **155.0**, WebKit **26.6**, Linux/x64, Node **24.21.0**. These are automated engines, not branded Chrome/Edge/Safari or minimum-version executions. Hosted image contents are not claimed immutable.

Browser artifact **10915800376** (`board-browser-evidence`, 9,013,792 bytes) was downloaded and its archive SHA-256 matched **828c5cc6aa6d57ab9c666ffeb4ef030bfde506e4ca71045a06a4388464e4d9e3**. Every test annotation and the CI size report identify Linux binary **4c6f137642478fc8e82177f28715e1e483e46abd40c38d8e48aace99d0bf1c5c**. That reported identity matches the independently retained local Linux correction inventory; no new CI Windows checksum is asserted. [CI size data](evidence/s03-operator-2026-09-27/ci-dashboard-sizes.json) is preserved separately from local measurements, with its own first-response timestamps/hashes. No measure is inferred from a green job alone.

Three **new CI** captures were extracted with [archive paths and hashes](evidence/s03-operator-2026-09-27/ci-image-provenance.json), then visually reviewed on 2026-09-27:
- [Chromium nominal 360 px](evidence/s03-operator-2026-09-27/ci-chromium-nominal-360.png): single-column layout, exact 9007199254740993 / 18446744073709551615 visible, exact CPU disclosure and all metric groups retained, no observed clipping.
- [Board-down stale](evidence/s03-operator-2026-09-27/ci-chromium-board-unreachable-stale.png): Board unreachable, Health unknown, explicit stale label and retained eligible counters.
- [Expired without server](evidence/s03-operator-2026-09-27/ci-chromium-expired-without-server.png): hostname and metrics absent, no presentable observations, last successful receipt retained.

This is a reviewed **CI engine-image subset**, not a new Windows screenshot or real mobile review. Earlier local captures retain their earlier binary provenance. The failed run, deterministic red reproduction and corrected local run remain unchanged above. The initial trace corroborates rapid duplicate HTTP calls; it still does not expose the exact internal early callback.

### Operator execution on 2026-09-27

The available host is a Windows x64 desktop, **Windows 10 Pro 22H2, build 19045.7725** (OS query/registry and Chrome version page agree). No Apple device or isolated historical browser environment was exposed by available tools; .tools contains the pinned Node installations, and the inspected conventional local VM locations had no provided environment. This is an availability statement, not a claim that the owner has no such device elsewhere.

Installed-file inventory found Chrome **154.0.8037.57**, Edge **154.0.4258.37**, Firefox **156.0.1**, and Windows Narrator **10.0.19041.4522**. All three browsers were launched; **only Chrome's exact executed version was independently confirmed in chrome://version**. Edge's version page and Firefox's runtime version were not confirmed; their inventory versions must not be promoted to completed browser tests. No historical browser was installed, no protection changed and no paid service used.

The existing native binary was reused: **6bfbbfb96e47d1ec1baf07d0a8e19e1f3d541d774c420848346b1f1ca568dcf7**. Its digest, all **87 source-manifest entries**, and local asset hashes were rechecked against the corrected 2026-09-26 inventory: no differences. The compiled JS is **887354a9e865afc3cce33a38e1aadf2d773f70ad69b7dbf46a48cd129250d8f3**, CSS **466c8ff6d19e85a5d317083aa2314163bdab85f50fead2116019338c147b9189**, HTML **73e4e91987858f459285fab450d0c1b344e534717ec5571cfef9870f2dc71a85**. The historical manifest was originally uncommitted; comparing its unchanged source entries with current HEAD now connects that local binary's source evidence to the committed fix without rewriting the old provenance. No build or full regression rerun was needed.

The pinned Node **24.21.0** interactive REPL imported the unchanged `scripts/browser-harness.mjs` and called `startBoard()`. It served the native embedded app from an empty working directory/PATH, using loopback-only controlled Health; origin was `http://127.0.0.1:64430`. No mock data was added to production assets. Fixture complete SHA-256 is **00865883f9477e8fe12a1aac4cd08cb1b6b71dc8da3c17e4507a51fcaaa9038e**; partial-documentary is **dfc3aded6cf6f8074b12ecc7b15fef8fca5ef1105a41b668bcc72d78905ea2cb** (authoritative common corpus).

[Session identities and limits](evidence/s03-operator-2026-09-27/operator-session.json) and [timestamped accessibility observations](evidence/s03-operator-2026-09-27/operator-observations.json) record the limited new execution:
- Chrome loaded complete observations: Board reachable, Health available, CPU 0.5%, exact cumulative counters **9007199254740993** and **18446744073709551615**, distinct active/since-boot flags.
- At **09:25:47.499Z**, controlled Health changed to partial-documentary. At **09:26:05.662Z**, Chrome's accessible page contained the partial label, null CPU/memory presentation as Unavailable, group messages and retained logical CPU count. This is PASS for the narrow real-Chrome parse/projection/text-update observation only.
- A Health 503 was selected at **09:26:13.214Z**. The subsequent accessibility snapshot still represented the earlier partial page while another browser was being opened; foreground/paint state was not reliably established. Stale/none, recovery and exact expiry are **NOT RUN to a conclusive operator result**, not falsely PASS or diagnosed as a product FAIL. The CI has separate passing assertions for them.
- Native Tab/F6 focus and screenshot attempts were inconclusive: captures failed with `FrameArrived timed out`, `window capture timed out` and `coordinate input geometry is unavailable`; a minimized-window recovery was attempted. The browser surface inventory exposed no alternate connected tabs. Windows reported an active console session, but this did not restore reliable interaction.
- Finally the computer-use policy guard stopped UI work because it could not determine the current browser URL sufficiently to enforce policy. No further UI interaction or alternate automation route was attempted. The operator subsequently reported that the Windows session appeared open; no locked-session cause is established. This answer does not restore the stopped UI tool or turn an incomplete check into PASS.
- Narrator is installed, but speech was not heard/recorded and no screen-reader interaction was verified. An accessibility tree is not speech experience evidence.

Setup incidents are retained: direct import into the automation JavaScript runtime failed with `process is not defined`; the first pinned Node REPL launch could not spawn Board in the sandbox (`EPERM`). The authorized native REPL then started successfully. The initial restricted OS query was denied, then the approved read succeeded. These are environment/tool incidents, not product test failures.

Cleanup executed `await b.stop()` and exited the Node REPL, confirming **S03 operator harness stopped**. The earlier failed REPL also exited; the automation JavaScript kernel was reset to release resources from its failed harness import. A final process check found no Board process. No Board service was installed. Because UI access was blocked, dedicated test/version tabs and newly launched browser windows could not be reliably closed; the operator should close those task tabs/windows without closing pre-existing personal tabs. No native zoom change was confirmed.

### Chrome retry with restored extension access

On **2026-09-27, 19:07–19:21 UTC**, the user-requested retry reached real installed Chrome through its restored extension connection. This supplements, and does not overwrite, the unsuccessful earlier native-tool attempts. HEAD remains **24d527a301d66b49d86ad11fbd75393254f02094** with pre-existing documentary changes. The same corrected Windows binary SHA-256 **6bfbbfb96e47d1ec1baf07d0a8e19e1f3d541d774c420848346b1f1ca568dcf7** and JS/CSS/HTML digests above were rechecked; no rebuild, code change or regression-suite rerun occurred.

The running page's `userAgentData.getHighEntropyValues` reported **Google Chrome 154.0.8037.57**, Windows, platformVersion 10.0.0, non-mobile; the previously recorded host is Windows 10 Pro 22H2 19045.7725. This is runtime brand/version evidence, not an emulated user agent or proof of Chrome 111. The version-page navigation was denied by the browser tool's URL policy; it was not circumvented. No new built-in version-page capture is claimed.

The unchanged pinned Node 24.21.0 harness served embedded Board at **http://127.0.0.1:50600**, using the authoritative complete/partial fixtures and `h.markup()` (65,536 bytes). Initial sandbox launch failed with spawn EPERM; that REPL exited, and the authorized launch succeeded. All product observations below used real Health HTTP → Board HTTP → installed Chrome, without response interception or synthetic lifecycle dispatch.

| Executed subset | Observation / result |
|---|---|
| OP-1, T09/T11 exact nominal presentation | **PASS within this subset**: CPU 0.5%, keyboard-opened exact 0.5037%, RX 9007199254740993 and TX 18446744073709551615 visibly intact. Nulls, flags and other groups remain distinct. [360 px capture](evidence/s03-operator-2026-09-27/chrome-retry-nominal-360.jpg). |
| OP-3, T09 keyboard/focus subset | **PASS**: Enter opens exact CPU disclosure, Tab reaches Memory exact-used, visible focus ring. Its focus and the CPU open disclosure survive ordinary refresh and the partial replacement. This does not establish all keyboard paths or screen-reader speech. |
| OP-2, T10/T11 partial replacement | **PASS**: selected partial-documentary at 19:08:27.252Z; null CPU/memory become Unavailable, logical CPUs remain 4, zero TX remains 0, group issues appear once. [Partial capture](evidence/s03-operator-2026-09-27/chrome-retry-partial-360.jpg). |
| OP-2, provider error and expiry | **PASS**: selected Health 503 at 19:08:42.857Z. Board remains reachable, Health unavailable, stale explicitly shown; subsequently no presentable observations and no hostname/metrics, preserving last receipt. [Stale](evidence/s03-operator-2026-09-27/chrome-retry-health-stale-390.jpg), [none](evidence/s03-operator-2026-09-27/chrome-retry-health-none.jpg). This observation does not replace deterministic boundary tests. |
| OP-2, recovery / maximum markup | **PASS**: at 19:09:18.131Z restored HTTP 200 with `h.markup()`. Ordinary refresh restored current data and exact large counters, without restart or forced retry. [768 px capture](evidence/s03-operator-2026-09-27/chrome-retry-recovered-markup-768.jpg). |
| OP-2, Board-down independent expiry | **PASS**: Board was killed at 19:10:18.399Z after its last successful receipt at 19:10:18.2797204Z. Chrome changed to Board could not be reached / Health availability unknown, retained explicitly stale values, then removed hostname/metrics without another successful server response. Receipt and generated timestamps remain unchanged. [Stale](evidence/s03-operator-2026-09-27/chrome-retry-board-down-stale.jpg), [expired](evidence/s03-operator-2026-09-27/chrome-retry-board-down-expired.jpg). |
| OP-3, T12 selected responsive images | **PASS for reviewed images** at requested widths 360, 390, 768 and [1280](evidence/s03-operator-2026-09-27/chrome-retry-nominal-1280.jpg), with height 900. Full-page captures exclude the 15 px scrollbar; exact values and issues remain readable. At 360 and 1280, measured scrollWidth equals clientWidth (345 / 1265), with no horizontal overflow. Partial load labels wrap at 360 but retain all text. Viewport overrides are desktop responsive evidence, not native zoom or a physical mobile device. |

[Timestamped DOM observations](evidence/s03-operator-2026-09-27/chrome-retry-observations.json), [network events](evidence/s03-operator-2026-09-27/chrome-retry-network.json) and [image dimensions/digests](evidence/s03-operator-2026-09-27/chrome-retry-images.json) preserve this retry independently of CI images. The network excerpt contains 19 overview starts, 15 responses and four failed loads, approximately five seconds apart; it is an excerpt, not a total session count or a no-overlap proof for all conditions. The harness counted 38 provider attempts before Board shutdown. Screenshot bytes are the tool's original JPEG output (extensions corrected without image re-encoding).

**Limits/incidents:** the accessibility API labelled the disclosure a button, but the DOM locator reported no matching button; the observed Exact utilization text locator successfully sent Enter. A read-only evaluation did not expose navigator; documented CDP runtime inspection supplied version metadata. One full-page expired capture timed out; the subsequent viewport capture succeeded. After additional tabs were created, a requested 390 override did not change the Board tab: measured width remained 1280, so the Board-down captures are recorded at their actual width. An attempted local auxiliary navigation and later direct API-document navigation were blocked by the browser client; the latter tool call took about 581 s. No bypass was attempted. Opening an auxiliary tab had left Board visible with no visibility events, so it is **not** a passed background test. A fresh lifecycle harness at 50801 was started, but no conclusive actual background/BFCache/lock-resume check followed. These tool incidents establish no product FAIL.

The actual-tab lifecycle, BFCache, suspension, native 200% zoom, full keyboard traversal, 64-interface/long-text manual review and screen-reader checks remain **NOT RUN** in this retry. Existing automated evidence is unchanged. Observation-only lifecycle listeners were attached to the test page but never triggered synthetically; the page was closed afterward. The viewport override was reset, all four agent-created tabs were closed, and both harness instances were stopped; final REPL cleanup completed at **19:21:29.005Z**. No persistent browser setting, service or production resource was changed.

Retry documentary checks: 22 Markdown files / 367 internal links and anchors checked without errors; all three new evidence JSON files parsed and all eight image hashes verified. `git diff --check` passed. This retry changes only this evidence record, the S03 scope/validation status wording, README and the new observation/image artifacts; earlier documentary changes remain preserved. No application/test/build command was rerun.

### Remaining acceptance checklist

Evidence abbreviations: **CI7** = corrected run/report and Linux binary above; **W27** = native Windows binary and earlier timestamped observations; **CR27** = restored Chrome retry above; **L26** = historical corrected local ledger. PASS applies only to the named subset; all rows are existing S03 requirements, not additional gates.

| Test / expected criterion | Required environment | Existing evidence / source | Status | Next action |
|---|---|---|---|---|
| T01–T04/T06–T08, T10/T13 automated assertions, exact parsing/expiry/regressions | Locked Node/Go and three supported automated engines | CI7; L26; exact same corrected source | PASS within implemented automated scope | Preserve logs; do not substitute for manual rows. |
| T09/T12 selected corrected images: 360 px exact values, stale and expired labels | CI Chromium | Three reviewed CI7 captures linked above | PASS for these images | Other real-browser visual checks remain below. |
| T09/T10/T11/T12 current Chrome executed subsets | Chrome 154.0.8037.57, Windows 10 19045.7725 desktop | CR27 real Health→Board→Chrome, reviewed nominal/partial/stale/none/recovery/Board-down expiry images and limited keyboard interaction | PASS for these subsets | Preserve actual scope; finish outstanding rows below. |
| T05/T09/T11/T12 remaining current Chrome flow | Same real Chrome desktop, working interactive session | CR27 improved access; full manual critical flow incomplete | NOT RUN | OP-3/4: full keyboard traversal, native zoom, long data/64 interfaces, actual lifecycle/BFCache/suspension. |
| T05/T09/T11/T12 current Edge flow | Actual Edge desktop, exact runtime patch/OS recorded | File inventory 154.0.4258.37; window launched, no completed Board flow | NOT RUN | OP-1–OP-4; confirm runtime version first. |
| T05/T09/T11/T12 current Firefox flow | Actual Firefox desktop, exact runtime patch/OS recorded | File inventory 156.0.1; window launched, no completed Board flow | NOT RUN | OP-1–OP-4; confirm runtime version first. |
| T11 current Safari desktop flow | Real Safari on macOS, exact patch/OS | No accessible Apple environment; WebKit CI7 is not Safari | NOT RUN | Apple operator runs OP-1–OP-4. |
| T11 Chrome minimum | Isolated Chrome **111**, exact patch/OS | No isolated installation/image supplied | NOT RUN | Obtain existing isolated environment, then OP-1/2/4; do not replace daily browser. |
| T11 Edge minimum | Isolated Edge **111**, exact patch/OS | Same gap | NOT RUN | Same floor procedure. |
| T11 Firefox minimum | Isolated Firefox **115**, exact patch/OS | Same gap | NOT RUN | Same floor procedure. |
| T11 Safari minimum | Safari **16.4** on matching macOS | No supplied environment | NOT RUN | Apple operator runs critical exact/expiry/lifecycle flow. |
| T11/T12 iOS Safari minimum | **Real iOS device / Safari 16.4**, model/OS patch | No connected device | NOT RUN | OP-5 plus OP-1/2/4; retain touch/zoom/background/back evidence. |
| T11/T12 current iOS Safari | Real iOS device, current exact version/model/OS | No connected device; mobile viewport is not device proof | NOT RUN | Same actual-device procedure. |
| T09 screen-reader transition/association and restrained announcements | Real supported browser + Narrator/NVDA/VoiceOver and an operator able to hear output | Narrator installed; no speech evidence | NOT RUN | OP-3, record spoken output and quiet repeated polls. |
| T09/T12 native 200% zoom and remaining keyboard/visual work | Actual desktop browser UI; target viewports 360/390/768/1280 | CR27 has reviewed current-Chrome responsive captures and Enter/Tab/focus subset; native zoom unchanged | NOT RUN for remaining scope | OP-3: full traversal, long data/64 interfaces, native zoom indicator and reviewed screenshots; do not substitute viewport overrides. |
| T05 actual tab/background/foreground, no hidden requests/early restoration | Actual browser and request-timing capture | CI7 synthetic events only; W27 attempt inconclusive | NOT RUN | OP-4, retain visibility/events and outgoing request timeline. |
| T05/T12 actual back navigation and BFCache restoration | Real browser capable of BFCache; record `pageshow.persisted` | No actual persisted=true evidence | NOT RUN | OP-4, document ineligibility and retry in an eligible isolated environment. |
| T05/T12 device suspension/lock/resume | Real desktop/mobile device and operator | No executed suspension proof | NOT RUN | OP-4/5; metrics masked until next valid slot. |
| T14 exact-commit hosted gate | GitHub run on corrected SHA | CI7 metadata/jobs/logs/report | PASS | Operator matrix still prevents whole T14/S03 closure. |

**No new product defect was established in this assignment.** The blocked/inconclusive UI attempts do not justify weakening an assertion or changing the product. Implementation and executed local/CI subsets are complete; **S03 is not closed** until the NOT RUN rows above have actual evidence. S04 remains NOT STARTED; packaging/systemd/server compression, Raspberry Pi integration/physical budgets and global A/B/C acceptance remain NOT RUN. No commit, push, tag, release or deployment occurred.

Documentary audit, 2026-09-27: internal links/anchors, fourteen S03 validation IDs, unchanged 87-entry executable-input manifest and retained local binary digests verified. Final diff contains documentation and evidence only; code, tests, fixtures, dependencies, toolchains, scripts/workflows and S00-S02 historical evidence are unchanged. Architecture/UX/specification/decision wording changed only to remove obsolete current-authorization/CI statements; no technical decision changed. No full application regression was rerun.
