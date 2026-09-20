# Sprint 00 — detailed preparation

Status: **PREPARED, NOT EXECUTED**. This document is the future implementation work package, not evidence that tasks have been done or authorization to run them during this documentation mission. No application directories, dependency manifests, CI workflows or packages have been created as part of this sprint now.

Objective: establish a reproducible minimum build and test foundation so S01 can implement the pinned Health contract with no uncertainty about toolchain, assets or test ownership. [Architecture](architecture-v0.1.0.md) owns the proposed package layout; [roadmap](roadmap-v0.1.0.md) owns dependencies; [validation](validation-v0.1.0.md) owns gate IDs.

## Entry conditions and decisions

- Confirm `Adrien-hue/joy-pi-board`, preserve current documentation and verify the working tree before any future scaffolding.
- Review the requirement/source-of-truth map and pinned Health revision. No moving-branch schema or copying Health internals.
- Ratify or explicitly replace foundation-affecting P-01/P-02 (ports, CLI/env and build/static layout) and P-03 (strict validation and exact parsing direction). Record owner/review evidence in the decision record when a real review occurs; none is invented now.
- Identify native build/race runner and target Linux ARM64. Verify available supported Go/Node toolchains at execution time, then lock exact versions. Do not infer contemporary tool versions from this document.
- Make absence of hardware nonblocking for S00, while leaving all Class C gates NOT RUN.

## Ordered work packages

All rows are **NOT STARTED**; each is intended to be a reviewable future change. Dependencies are explicit to avoid parallel tasks silently choosing incompatible conventions.

| Task | Future work and concrete output | Depends on | Completion evidence |
|---|---|---|---|
| S00-01 | Record short foundation decision notes: approved/replaced P choices, Go module identity proposed as `github.com/Adrien-hue/joy-pi-board`, exact tool versions, supported dev OS/CI platform, frontend package manager and browser bigint target. | Entry review | Decisions separate established/proposed facts; no unresolved choice hidden in a manifest. |
| S00-02 | Establish exact Go version (module/toolchain and CI enforcement), exact Node + package-manager versions, React/TS/Vite and minimal lint/test dependency versions with lockfile. Record licenses and selection rationale. | S00-01 | Tool-version mismatch fails clearly; no implicit auto-upgrade/latest/floating CI image; chosen versions verified compatible by a clean install. |
| S00-03 | Create only needed directories/files in the proposed `cmd`, coherent `internal` packages and `web` layout. Add module/manifests/ignore rules, formatter/linter config, concise contributor commands. Do not create empty future packages for architecture symmetry. | S00-02 | Tree builds with clear ownership, no Health internal dependency, no collectors, no unnecessary framework/global store/UI kit. |
| S00-04 | Add a minimal frontend shell and Go entry point proving production asset embedding; build frontend into `web/dist`, embed via package `web` child path, serve shell and actual generated assets. This is a build smoke surface, not a functional dashboard. | S00-03 | Clean frontend build precedes Go compile; absent dist fails rather than falling back to local files; binary in empty directory serves exact assets. |
| S00-05 | Establish frontend install/check/build and Go format/test/vet/module/build commands, including native and Linux ARM64 builds. Declare order/dependencies in a small build script/task entry point suited to CI and documented local use. | S00-04 | Clean checkout, no Node modules/dist/Go cache assumptions, succeeds twice from locked inputs with no source/lockfile changes; record artifact hashes and tool versions, without claiming bit-for-bit reproducibility unless actually compared. |
| S00-06 | Set up useful Go and frontend test harnesses around shell delivery and asset integration, fake-server and injectable-clock interfaces/conventions. Define fixture locations, ownership, taxonomy and expected-result metadata for S01; seed positive documentary examples as test data only with reviewed provenance. | S00-03, S00-04 | Test runner finds genuine assertions, not empty PASS jobs; example seeds parse exactly; no tests pretend to validate unimplemented provider/cache functionality. |
| S00-07 | Specify startup configuration/help/version contract and test plan in developer docs, including precedence, invalid input, port conflict, Health-absent startup and fixed resilience constants. Minimal build shell may bind the ratified listener; full config/lifecycle behavior remains S02. | S00-01, S00-04 | Config table matches P-01/P-02; frontend sees no Health URL; no startup Health request; A-05 unimplemented cases explicitly NOT RUN. |
| S00-08 | Add CI jobs for reproducible frontend install/type-check/lint/tests/build, followed by Go formatting/tests/vet/module checks/native and ARM64 builds; suitable separate native race job also builds assets first. Publish only internal check artifacts/logs, not a release. Pin actions/tools/runner image identities as practicable and record unavoidable hosted-image drift. | S00-05, S00-06 | Clean CI run with no hidden dist dependency, immutable dependency/action references, read-only default token permissions, no deploy/package/release job and no required secrets. |
| S00-09 | Review build inventory and initial compressed size; confirm no runtime remote resource, no source maps/secrets in deployable assets and no production Node dependency. Capture preliminary A evidence, links and remaining gaps. | S00-08 | A-01–A-03 foundational checks and A-04 shell embed smoke PASS; asset budgets measured for shell explicitly, not claimed for final UI. |
| S00-10 | S00 exit review: reconcile docs, tree and CI, update traceability/status per actual evidence, hand off S01 fixtures/parser tasks and unresolved functional proposals. | S00-09 | Review record names exact commit/artifacts and test runs; all deferred B/C checks retain NOT RUN. |

## Toolchain and build details to lock

Go: pin an exact supported patch version after review; module identity remains Board's. Prefer standard library; any introduced module needs a specific purpose. Run module consistency checks after the full intended files exist, not against this documentation-only repository. Record `GOOS`, `GOARCH`, CGO setting and build metadata strategy; proposed release build disables CGO where dependencies permit and uses reproducible path/build flags. Verify runtime independence rather than assuming “Go binary” means static.

Frontend: proposed package manager is npm with an exact pinned version and committed lockfile, installed through `npm ci` (review as part of S00-01). Exact dependency versions/lock entries must resolve reproducibly. Set TypeScript strict null checking and a bigint-capable compilation/browser target; document the minimum browser versions selected during S00. Avoid `any`-based assertion of fetched payloads. A browser data parser needs raw numeric tokens, not standard `response.json()`.

Build chain: clean locked frontend install → type-check/lint/tests → Vite production build → generated asset inventory/precompression if adopted → Go checks requiring embed → native/cross binaries → actual embedded-asset smoke. Use package `web` with directive `dist`, never a parent traversal. The production command must build assets before Go; a Go-only convenience command must fail clearly if assets are missing. Generated output should not create unrelated source diffs or leak dev-server endpoints.

No bit-for-bit reproducibility promise is established merely by dependency locking. S00 establishes repeatable steps/inputs; if byte-identical builds are desired, compare two clean artifacts and document metadata/timestamps/compression influences. Later release acceptance is always tied to the actual resulting checksum.

## Tests and fixture preparation

Proposed future locations: `internal/health/testdata` for upstream wire fixtures and invalid mutations, plus `web` test data referencing the same authoritative fixture source or a checked synchronization step. Avoid two independently edited fixture copies. JSON examples under [docs/examples](examples/overview-current.json) are explanatory Board envelopes; extract their snapshots deliberately, preserving raw integer tokens. Do not run a JavaScript parse/stringify step that rounds them while generating fixtures.

S00 defines fixture IDs/expected outcomes; S01 writes the complete cases and validators. Required taxonomy:

- Complete, partial per group, all three issue codes, firmware group and independent temperature, whole-network failure and interface-local state/RX/TX failure, true/false/zero values.
- Exact integer lexical values at safe boundary, `9007199254740993`, uint64 maximum and maximum+1; invalid quoted/fractional forms and invalid UTF-8/JSON.
- Missing-versus-null every required field, null object versus null leaves, wrong issue path/order/code, unknown schema and network state under the ratified policy.
- Transport-controlled cases planned for S02: non-200, refusal, delayed/stalled body, excessive body, deadline, cancel/late-result ordering and no retry.

Clock seam must allow independent wall and elapsed-time advancement; fake Health must expose request counts, barriers and controlled headers/body chunks. S00 may define a harness interface and a simple demonstration, but must not implement cache/resilience semantics as incidental scaffolding. Use normal source tests plus real HTTP integration where transport behavior matters. Race checks cover concurrency when introduced in S02. Frontend tests need controlled fetch/time and later lifecycle/visibility events. Plan parser fuzz/differential and integer rendering tests for S01/S03, not tests that merely restate type definitions.

## Configuration and CI handoff

Document intended defaults and help text semantics before wiring real behavior: Board proposed `0.0.0.0:8081`, Health fixed existing default URL on 8080, CLI/env precedence from architecture; invalid config exits cleanly with safe diagnostics. No user-configurable freshness/timeout that silently breaks the established contract. Full configuration validation, help/version metadata, shutdown, admission and significant-transition logs are S02; Debian paths/env file/unit are S04.

CI must distinguish source checks, binary smoke tests and unavailable physical tests. Hardware jobs must not be green placeholders. Race on suitable native platform is not replaced by cross-build. CI workflow definitions are future S00 outputs; **none exists or is executed as part of this preparation**. No release publishing credentials or production deployment step is needed.

## Explicit sprint boundary

S00 ends with a reproducible build shell and validation scaffolding. It does **not** deliver a validated Health decoder/parser (S01), live provider fetch or overview/cache/recovery (S02), metric dashboard/freshness UI (S03), Debian/systemd artifacts or performance/physical acceptance (S04). A preliminary small shell bundle does not pass the final frontend budget gate. No Health code or deployment is changed at any sprint without its own scope; this project remains a consumer.

All work above remains future work after this documentary task. See [documentation review](documentation-review.md) for the much narrower checks actually performed now.
