# Sprint 03 implementation and evidence

Historical implementation ledger: 2026-09-25. The [2026-09-26 scheduling correction](#2026-09-26--early-cadence-callback-correction) below records the first S03 CI failure and its targeted fix; earlier measurements/manifests remain historical. **S03 implementation delivered; local automated verification complete; hosted CI and full sprint closure PENDING.** No S04, commit, push, tag, release or deployment. S00/S01/S02 retain their separate historical evidence; their CI does not validate this implementation.

## Source and authorization

Work began in the correct Board repository with a clean tree at `7d9dc989dbbb31409450d5f5ff394ccffb5e43c0` (`docs: prepare Sprint 03 dashboard and browser resilience`). That commit already contains the prior S02 closure and S03 preparation. There were no pre-existing uncommitted changes to overwrite. S02 implementation remains `a11f07c85ff4071a13ca3ae9227920e9f1c7dda4`; the later documentary commit is not attributed its CI results. Current changes are uncommitted, identified by the accompanying input manifest, not by inventing a new SHA. Health remains read-only at baseline `be7a0d824f62b94c842d8e5326110b1852c5a0bc`; no Go server behavior, Health contract, shared fixture bytes or toolchain version changed.

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
| S03-05 | Dashboard implemented; static, engine keyboard/reflow and reviewed-image subsets PASS. Screen reader, native zoom and actual-device portions of S03-T09/T12 remain NOT RUN. |
| S03-06 | Three supported automated engines and real HTTP chain PASS (S03-T10). Actual branded/floor/mobile evidence S03-T11 and real lifecycle/device portions of T05/T12 remain NOT RUN. |
| S03-07 | Local regression/build/size checks PASS (T13); workflow extended, new hosted CI NOT RUN. |
| S03-08 | Local provenance/traceability recorded; T14 and sprint closure PENDING until all mandatory browser/manual/CI evidence exists. |

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

**NOT RUN / closure blockers:** actual Chrome/Edge 111, Firefox 115, Safari 16.4 and iOS Safari 16.4; current branded Chrome/Edge/Firefox/Safari manual flows; actual iOS device; screen-reader transitions/associations; native 200% browser zoom; actual tab/background/device suspension and verified BFCache restore. Required targets remain S03, not silently moved to S04. Engine tests are supplemental.

Operator procedure (use an isolated trusted test environment and the recorded binary; no production Health changes):

1. Record source manifest/commit, binary and asset SHA-256, exact browser patch, OS/build, device model, DPR, viewport and zoom. Obtain actual floor-version environments; do not substitute a spoofed user agent. Build changes require new identities and appropriate reruns.
2. Install pinned tools/dependencies, build frontend then native Board, and run the engine suite above (or `npm --prefix web run browser` on supported Linux with installed engines). Retain `out/browser-report.json`, screenshots/traces and inventory. For interactive local desktop testing, start `node --experimental-repl-await`; use `const h = await import('./scripts/browser-harness.mjs'); const b = await h.startBoard(); console.log(b.origin)`. Open that loopback origin in the actual browser. Set `b.state.body=h.fixture('partial-documentary')`, `b.state.status=503`, restore 200/complete or `h.markup()`, and await the next ordinary slot. `await b.kill()` tests Board-down expiry; `await b.stop()` cleans up.
3. For an actual iOS device, use an operator-controlled LAN test host and a controlled Health server on loopback; serve the same binary with `--listen <test-host-LAN-IP>:8081 --health-url http://127.0.0.1:<controlled-provider-port>/v1/snapshot`. Reuse original fixture bytes and the same status transitions. Record the host artifact hash separately if a different native build is needed. Do not expose Health or Board to the Internet or install a service. This is software browser acceptance, not Raspberry Pi C acceptance.
4. Verify complete/partial/none/stale/recovery and separate Board HTTP/connection failure text; exact visible 9007199254740993/18446744073709551615; keyboard focus and disclosure state; 360/390/768/1280 reflow, native zoom 200%, long labels/messages/64 interfaces and copyable counters. Use a screen reader to verify group issue associations and one moderated status transition, with no announcements on every poll or age tick.
5. Hide the tab, switch apps/lock/unlock the device, navigate away/back and record whether pageshow.persisted really is true. Metrics must stay masked until a new valid next-slot response; no pre-hide result can revive them. Stop Board after success and observe values/hostname/exact details disappear by the conservative 30-second limit without a new response. Capture screenshots/video and request timing/counts.
6. Record PASS/FAIL/NOT RUN per S03-T05/T09/T11/T12 with the actual identities and reviewer. Any missing target keeps S03 closure pending. A new hosted run on the exact future S03 commit must independently pass foundation/race/browser; this mission performs no commit or push to obtain it.

## Hosted CI and sprint status

The existing workflow retains immutable action pins and Ubuntu 24.04 labels; foundation/race remain, and a browser job installs exact package-selected revisions, builds the embedded app, runs all three engines and uploads results. It records hosted ImageOS/ImageVersion, OS/kernel and browser revision data; the hosted image label is not claimed immutable. **At the 2026-09-25 local delivery, S03 hosted CI was NOT RUN** for that then-uncommitted state. The subsequent run on commit 12108b60861582cab040b9d5a029279e5768bcea failed as recorded in the addendum below; corrected code still needs new CI. An edited workflow is not a green run.

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

Corrected hosted CI is **NOT RUN**. S03 remains **IMPLEMENTED, LOCALLY VERIFIED WITHIN EXECUTED SCOPE; CLOSURE PENDING**, subject to the new exact-commit CI and all existing required operator/browser-floor/Safari/iOS/native-zoom/screen-reader/lifecycle proofs. This correction neither waives those gates nor starts S04.

Correction finishing checks: format/types/lint passed within the full check; git diff --check is clean. Internal links/anchors and the 87-file input manifest were rechecked. Final diff is limited to the controller, deterministic tests, S03 evidence and its validation trace; browser test, StrictMode entry, dependencies, toolchains, workflows, Go server and historical evidence artifacts are unchanged.
