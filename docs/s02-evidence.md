# Sprint 02 execution evidence

Local implementation and verification: 2026-09-22–23, final checks 2026-09-23 (Europe/Paris). Repository: **Adrien-hue/joy-pi-board**, initially clean at `f3d8dc8959c76e78b53d2913ef068baa89618567` (S02 documentary preparation). The implementation instruction explicitly approved the prepared P-02/P-04/P-05 recommendations; the preparation itself did not. No commit, push, tag, release or deployment was made.

**Status: S02 implemented and verified locally; hosted CI closure pending.** S02 is not closed. S00 and S01 remain closed with their unchanged, separate [S00](s00-evidence.md) and [S01](s01-evidence.md) evidence. S01 hosted run 35637699325 covers only `ccd1245a64fc372faa4ebc85eb000f78a01bd396`, not its documentary successors or this working tree.

## Source identity, tools and scope

The tested uncommitted software/test/build inputs are enumerated in [s02-inputs.json](evidence/s02-inputs.json): **77 files**, aggregate SHA-256 `f58a2d94ba91b91c614994f71226a23f625ca79fc57e5b0f411d31dfd3ed7634`. The aggregate is SHA-256 of sorted UTF-8 `sha256 + two spaces + path + LF` rows. It excludes documentation/evidence/output files; executable tests' JSON examples are included. Generated embedded resources are identified separately in [the shell inventory](evidence/s02-shell-sizes.json). This is a working-tree fingerprint, not a commit or remote CI proof.

Observed tools: **Go 1.27.1, Node.js 24.21.0, npm 11.19.0**; `tasks.mjs tools` checks exact versions with `GOTOOLCHAIN=local`, GOENV/GOWORK off in build tasks. No dependency, toolchain, module or action-pin upgrade. Go uses only the standard library. Existing lossless-json 4.3.1 remains the S01 runtime dependency; no production TS consumer was added. The React shell and original Health/shared fixture bytes are unchanged. Health remains read-only at baseline `be7a0d824f62b94c842d8e5326110b1852c5a0bc`; no packages were imported from it.

Main host: Windows NT 10.0.19045, x64. Native race/signal tests used Docker Desktop's Linux amd64 environment, pinned image `golang@sha256:69a7b9788769bec032d238959b61854e9ae87f57be9029ec04e9885fabf99195`, network disabled, source bind-mounted read-only. Observed: `go version go1.27.1 linux/amd64`, Debian 12 bookworm, kernel `6.18.33.2-microsoft-standard-WSL2`, x86_64. This is neither a GitHub hosted image nor Raspberry Pi OS. No container image drift or Pi result is inferred from these checks.

## Executed commands and results

Commands below were executed from the repository root with the pinned Node distribution on PATH. PowerShell uses `npm.cmd`. The scripts' frontend-before-Go ordering is retained.

| Evidence | Command / procedure | Observed result and scope |
|---|---|---|
| S02-L01 | `npm --prefix web ci --offline=false --cache .cache/npm-s02-clean`; `node scripts/tasks.mjs tools` | PASS: locked install, 178 packages, exact tools. Fresh npm cache; no lockfile drift. Initial sandbox download failure and successful retry are recorded below. |
| S02-L02 | `npm --prefix web run format`; Go `gofmt -w` on edited packages; `node scripts/tasks.mjs check` | PASS: final format, strict TS types, lint, installed dependency consistency, 271 frontend tests, two reference/numeric fixture checks, Vite build, Go formatting/test/vet (the final two counting/wave test refinements were additionally checked with `go test -count=1 ./internal/health ./internal/overview` and `go vet` on those packages, then the full Linux race run), `go mod tidy -diff`, `go mod verify`. Frontend built before embed. |
| S02-L03 | `node scripts/tasks.mjs contract` (also executed by check) | PASS: unchanged S01 corpus 264 cases, 76 accepted Go outputs parsed in TS, 512 seeded mutations (seed 20260920), then S02 real HTTP corpus below. Dedicated parser build retains actual parser exports. |
| S02-L04 | `node scripts/tasks.mjs fuzz` | PASS: bounded 10 s request, two workers; recorded 484,470 executions, 257 distinct initial seeds, 20 additional interesting inputs, 11.129 s package duration. This exercises the unchanged S01 decoder, not a dashboard or transport fuzzer. |
| S02-L05 | `go test -race -count=1 -v -timeout=90s ./cmd/... ./internal/... ./web` inside the pinned Linux container, `CGO_ENABLED=1` | PASS: no race report. [Machine-readable summary](evidence/s02-race-summary.json) retains 43 top-level passed test names, timing lines and raw-log hash. `TestExportInterop` and `TestHTTPTypeScriptInterop` intentionally skip without their export/Node variables in this job; both execute through L03 in the foundation path. Real Linux SIGTERM subprocess test executes here, not on Windows. |
| S02-L06 | `node scripts/tasks.mjs build`; `node scripts/tasks.mjs cross` | PASS: native Windows x64 and Linux ARM64. Both compile frontend first, then Go with CGO disabled, `-trimpath -buildvcs=false`. Initial pipeline used new `.cache/go-s02-clean`; an intermediate corrected source was repeated with new `.cache/go-s02-final-empty`. The final routing/UTC corrections were then built with new `.cache/go-s02-final-route-cold`. Local build identities are below. |
| S02-L07 | `node scripts/tasks.mjs smoke` | PASS: final native binary copied to an isolated directory, PATH empty, exact embedded HTML/JS/CSS bytes and headers, HEAD, missing route, and actual unavailable-Health overview HTTP 200. Ephemeral reset trap saw zero startup/asset requests, one explicit overview connection, no additional call over six seconds idle. [Standalone report](evidence/s02-standalone-smoke.json). |
| S02-L08 | `node scripts/tasks.mjs missing-dist`; `node scripts/tasks.mjs sizes` | PASS: isolated source without dist fails clearly at embed; offline gzip-9 inventory remains shell-only, identity bytes are served. JS+CSS 68,495 bytes gzip-9; HTML+JS+CSS 68,790. No final dashboard or first-overview resource budget claim. |
| S02-L09 | `readelf -h -l` on final ARM64 output; `go version -m` | PASS: ELF64 little-endian AArch64 executable, no INTERP; Go1.27.1, CGO_ENABLED=0, linux/arm64, GOARM64=v8.0. This is inspection/cross-compilation, not ARM64 execution. |
| S02-L10 | Internal Markdown link/anchor audit, evidence/input hash consistency, scope/status/trace review, `git diff --check` | Documentary checks recorded at the end of this ledger. No result is promoted to hosted CI or physical acceptance. |

Detailed S02-T01–T12 → requirements → actual test names are in [the local execution trace](validation-v0.1.0.md#s02-local-execution-trace). The original test matrix remains the acceptance oracle; tests are not renamed to imply broader product completion.

## HTTP, precision and resilience proof

[HTTP interop report](evidence/s02-http-interop.json): **266 real simulated Health requests → production composed Board HTTP responses → S01 TypeScript parsing in Node**. The unchanged common corpus supplies 264 inputs; two in-memory byte-derived cases add maximum markup/Unicode and depth 32. Exactly **78 valid and 188 invalid** outcomes matched; exactly **266 outgoing Health attempts**, with no retry. Observed maximum overview size **65,737 bytes**, depth **34**; test-only envelope guard rejects 81,921 bytes and depth 35. Snapshot input limits remain 65,536/depth 32. Observed wrapper length varies slightly with RFC3339Nano fractional digits; the enforced envelope bound is unchanged.

Valid cases preserve zero, null, false, float semantics, issue order, unknown members and uint64 tokens including `9007199254740993` and `18446744073709551615`. Actual nested bytes are supplied to the S01 parser and compared with the authoritative input, then exact decimal primitives; no `response.json()`/ordinary Number oracle. Go additionally compares compact validated snapshot bytes against the actual JSON object in the response. The max-size markup/Unicode object remains 65,536 bytes after extraction and does not gain optional HTML escapes. The test adapter is not a production S03 envelope consumer.

Dedicated parser measurement remains **17,400 bytes minified / 5,550 bytes gzip-9**, SHA-256 `c5d9f605b355731b04cc6231741d1b03ba3078967febd69625ebc84363f5cf70`; [separate parser report](evidence/s02-parser-sizes.json). The unused parser is absent from the unchanged shell; no zero-cost inference is made for the future dashboard.

Deterministic fake-clock tests verify inclusive 29.999 / 30 / 30.001 s fallback, exact one-second terminal priority, full-to-partial replacement without merging, repeated reasons, recovery and cache loss in a fresh coordinator. Deadline and publication barriers preserve the validation-completion anchor and prevent old workers replacing a newer partial. Three waves of ten callers at five-second fake-clock intervals each share one active real Health GET (three GETs total); a demand after completion starts another. Individual/all cancellation does not cancel other waiters or fabricate Health failure; an orphan remains joinable until its original deadline. A counted infinite-body stub verifies at most 65,537 collected bytes, and zero body reads when Content-Length already exceeds the limit; early transport timeout is independently classified. Worker permits are held through actual late exits, capped at 32. Tests exhaust 32 requests and 64 connections and verify rejection/cleanup rather than enqueueing.

Old response views never combine their reason with a newer cache. Preparation crossing expiry removes failed fallback; delayed successful views become Board 503, invalid/regressing elapsed reads Board 500. Header/state/timestamps share finalization, with canonical conservative age and bounded reassembly. Tests cover blocked response/log writers, transition deduplication, safe diagnostics, immediate second-signal shutdown and no raw upstream destinations/errors/issue messages in logs.

## Real-time observations, not physical gates

Final native Linux race run, single observations (not p95 and not Raspberry Pi measurements). PowerShell corrupted the non-ASCII microsecond symbol in two captured Go duration strings; the summary restores the known Go `µs` unit and retains the original raw-log hash:

| Controlled case | Observed duration |
|---|---|
| Health headers withheld | 1.0002539s |
| Health body stalled | 1.0006606s |
| Health body progresses without completing | 1.0008299s |
| Loopback refusal | 305.8µs |
| Server request headers absent, net.Pipe | 2.0008679s |
| Server response writer blocked, net.Pipe | 3.0010456s |
| Real SIGTERM subprocess exit | 1.0034997s |
| Logger blocked, original stop context expires | 5.0012782s |
| Logger blocked, second signal | 133.3µs |

Fake time verifies exact semantic cutoffs; real tests assert cancellation/cleanup with watchdog tolerance. The production shutdown context has one absolute five-second deadline; the measured blocked-writer return includes scheduler/return overhead and is **not** a claim that a physical ≤5 s gate passed. The 15 s idle deadline is observed on a real Server connection through a deadline recorder without a 15 s sleeping test. No 30-second real sleep is used. Server suspension may pause the platform monotonic clock; no wall-time fallback was introduced.

## Local artifacts and reproducibility limits

[Build record](evidence/s02-builds.json), final corrected source:

| Artifact | Bytes | SHA-256 |
|---|---:|---|
| Windows x64 native binary | 11,447,808 | `dc5f60e2d73b68fea7fecc01fbffb63d49c767d1d4a4a8d218697f5206454855` |
| Linux ARM64 binary | 10,495,796 | `45bd17ddbb1e11408cf61a1b535eb3042f9ed84a3176145c5544a939e0378939` |

These are local hashes, not CI artifact hashes, release checksums or accepted Pi candidates. The final native/cross build process was run from a new empty Go cache after the last production corrections. An earlier intermediate source repeat was byte-identical across local cache modes, but those intermediate binaries are superseded. No final-source bit-for-bit reproducibility claim across builds, machines or hosted images is made. The final standalone smoke identifies its actual native binary hash; changing or rebuilding a future release candidate still requires the corresponding acceptance evidence.

## Incidents and corrections retained

- The first direct Go command used an unwritable default cache/telemetry location in the sandbox. Explicit workspace caches and the approved execution environment succeeded; the failed attempt is not a test failure hidden as PASS.
- Fresh npm installation initially failed with sandbox registry `ECONNREFUSED` through its restricted proxy. The authorized network retry installed the same lockfile successfully; no registry/dependency substitution.
- The first HTTP test adapter tried lossless-json parse/stringify. Valid unknown prototype members exposed that stringifier's duck-typed number handling. The adapter now extracts the original text of Board's fixed final snapshot field and reuses S01's private-branded parser. No production parser, global prototype change or Number conversion was added.
- The blocked-log second-signal test initially took the remaining five seconds: forced Close did not cancel the later log flush context. The code now cancels that same absolute context on a forced path; the regression test passes, including the second signal during log flushing.
- Final HTTP review found net/http's automatic `OPTIONS *` handler bypassed Board routing. It is now disabled and both direct and raw-socket tests require Board 400 for non-origin targets. Timestamp bounds are also checked after UTC conversion, with an explicit year-overflow regression.
- One local race invocation overlapped Vite rebuilding `web/dist`, producing embed missing-file build errors. The invocation was invalid because shared generated assets were changing. The final Linux run was repeated **after** frontend/native/cross completion with a stable read-only mount and passed. CI jobs have isolated checkouts and retain sequential frontend-before-Go order.
- An intermediate pre-fix native hash `f4903f8352e8339f19a17fd8d516f01955d6fa676f5263e93593106b7db6edf8` and ARM64 hash `5ab11a250b4ac4286c580b56f2c5cae4f8cb786cd0d56c73815d1efc39cab90e` are superseded by the corrected builds above. They are not the delivered evidence identities. The subsequent pre-routing-correction native `676640e39813d94918e7755812f9b298480a73fae0047a114f0f6d27bd16f9d8` and ARM64 `40edaf1d155142a6cae037e615b0f1dc47525c0e9b41768a00259de234e7f7e9` repeats are likewise superseded.

## CI and unexecuted validation

The existing **Board foundation and contracts** workflow keeps its immutable action references, ubuntu-24.04 label and exact tool pins. Foundation `check` now executes the S02 Go suites and real HTTP/TypeScript harness; race includes all new Go packages. Smoke asserts the real API instead of the historic 404. Existing fuzz/build/cross/missing-dist/size checks remain. Workflow labels were updated; **no new hosted S02 run exists in this task**, and creation/modification of CI is not a passed run. A future exact-commit run must record actual hosted image/tools, run attempt/event/SHA/jobs/steps and artifacts before S02 closure.

Still **NOT RUN**: target-browser execution, dashboard/polling/Fetch integration, visual/mobile acceptance, packaging/systemd, real Health on Raspberry Pi 3B+, physical performance/CPU/RSS/readiness/shutdown budgets and all Class C. No complete A/B/C product acceptance, release or deployment. P-06 presentation/measurement conventions and S04 extensions remain deferred. No S03/S04 implementation was performed.

## Documentary verification

Final audit on 2026-09-23: 21 Markdown documents and 236 internal file links/anchors checked, all links/anchors and JSON example blocks valid; source-input and final binary hashes match their ledgers; `git diff --check` clean. Final diff review confirms S02 scope and no S03/S04 code. Historical S00/S01 evidence files, original fixtures, dependency/tool locks and action pins are preserved. Stable requirements and S02 test/task IDs remain linked; no planned browser or physical test is labeled PASS.
