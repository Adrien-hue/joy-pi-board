# Sprint 00 — work package and execution status

Historical S00 record: subsequent explicitly authorized S01 work and its current status are tracked separately in [S01 scope](sprint-01.md) and [S01 evidence](s01-evidence.md). Statements about S01 not being started below describe the S00 closure, not the later working tree.

Status: **COMPLETE — locally verified and CI verified on 2026-09-20**. The implementation instruction explicitly authorized S00; the subsequent closure instruction authorizes documentation updates only. [S00 evidence](s00-evidence.md) records local checks and successful run 35527458889 for exact commit `b6e139c3378cd77fdcb4c5edc610977c08bbf72a`. **S01 is ready, NOT STARTED**. No S01–S04 work is authorized or executed.

Objective: establish a reproducible minimum build and test foundation so S01 can implement the pinned Health contract with no uncertainty about toolchain, assets or test ownership. [Architecture](architecture-v0.1.0.md) owns the proposed package layout; [roadmap](roadmap-v0.1.0.md) owns dependencies; [validation](validation-v0.1.0.md) owns gate IDs.

## Entry conditions and decisions

- Confirm `Adrien-hue/joy-pi-board`, preserve current documentation and verify the working tree before any future scaffolding.
- Review the requirement/source-of-truth map and pinned Health revision. No moving-branch schema or copying Health internals.
- Apply the S00 instruction's approvals: P-01/defaults, module/stack/npm/embed layout and exact-integer convention. Record P-02's limited shell choices and P-03's separated concerns; unknown-member policy and parser selection remain S01 and do not block S00. See [decisions](decisions.md).
- Identify native build/race runner and target Linux ARM64. Verify available supported Go/Node toolchains at execution time, then lock exact versions. Do not infer contemporary tool versions from this document.
- Make absence of hardware nonblocking for S00, while leaving all Class C gates NOT RUN.

## Ordered work packages

The original task definitions below retain their stable IDs. Implementation and applicable exit criteria are complete; the execution ledger records local and hosted evidence. Dependencies were followed; closure relies on the verified CI run, not merely the workflow file.

| Task | Defined work and concrete output | Depends on | Completion evidence |
|---|---|---|---|
| S00-01 | Record foundation choices: approved/replaced/deferred P decisions, approved module `github.com/Adrien-hue/joy-pi-board`, exact tool versions, supported dev OS/CI platform, npm and browser bigint targets. | Entry review | Decisions separate established/proposed facts; no unresolved choice hidden in a manifest. |
| S00-02 | Establish exact Go version (module/toolchain and CI enforcement), exact Node + package-manager versions, React/TS/Vite and minimal lint/test dependency versions with lockfile. Record licenses and selection rationale. | S00-01 | Tool-version mismatch fails clearly; no implicit tool auto-upgrade; explicit OS label and actual hosted image identity recorded (contents may drift); chosen versions verified compatible by a clean install. |
| S00-03 | Create only needed directories/files in the proposed `cmd`, coherent `internal` packages and `web` layout. Add module/manifests/ignore rules, formatter/linter config, concise contributor commands. Do not create empty future packages for architecture symmetry. | S00-02 | Tree builds with clear ownership, no Health internal dependency, no collectors, no unnecessary framework/global store/UI kit. |
| S00-04 | Add a minimal frontend shell and Go entry point proving production asset embedding; build frontend into `web/dist`, embed via package `web` child path, serve shell and actual generated assets. This is a build smoke surface, not a functional dashboard. | S00-03 | Clean frontend build precedes Go compile; absent dist fails rather than falling back to local files; binary in empty directory serves exact assets. |
| S00-05 | Establish frontend install/check/build and Go format/test/vet/module/build commands, including native and Linux ARM64 builds. Declare order/dependencies in a small build script/task entry point suited to CI and documented local use. | S00-04 | Clean checkout, no Node modules/dist/Go cache assumptions, succeeds twice from locked inputs with no source/lockfile changes; record artifact hashes and tool versions, without claiming bit-for-bit reproducibility unless actually compared. |
| S00-06 | Set up useful Go and frontend test harnesses around shell delivery and asset integration, fake-server and injectable-clock interfaces/conventions. Define fixture locations, ownership, taxonomy and expected-result metadata for S01; seed positive documentary examples as test data only with reviewed provenance. | S00-03, S00-04 | Test runner finds genuine assertions, not empty PASS jobs; example seeds parse exactly; no tests pretend to validate unimplemented provider/cache functionality. |
| S00-07 | Specify startup configuration/help/version contract and test plan in developer docs, including precedence, invalid input, port conflict, Health-absent startup and fixed resilience constants. Minimal build shell may bind the ratified listener; full config/lifecycle behavior remains S02. | S00-01, S00-04 | Config table matches P-01/P-02; frontend sees no Health URL; no startup Health request; A-05 unimplemented cases explicitly NOT RUN. |
| S00-08 | Add CI jobs for reproducible frontend install/type-check/lint/tests/build, followed by Go formatting/tests/vet/module checks/native and ARM64 builds; suitable separate native race job also builds assets first. Publish only internal check artifacts/logs, not a release. Pin actions/tools/runner image identities as practicable and record unavoidable hosted-image drift. | S00-05, S00-06 | Clean CI run with no hidden dist dependency, immutable dependency/action references, read-only default token permissions, no deploy/package/release job and no required secrets. |
| S00-09 | Review build inventory and initial compressed size; confirm no runtime remote resource, no source maps/secrets in deployable assets and no production Node dependency. Capture preliminary A evidence, links and remaining gaps. | S00-08 | A-01–A-03 foundational checks and A-04 shell embed smoke PASS; asset budgets measured for shell explicitly, not claimed for final UI. |
| S00-10 | S00 exit review: reconcile docs, tree and CI, update traceability/status per actual evidence, hand off S01 fixtures/parser tasks and unresolved functional proposals. | S00-09 | Review record names exact commit/artifacts and test runs; all deferred B/C checks retain NOT RUN. |

## Toolchain and build details to lock

The choices below describe the preparation intent. Exact selections and implemented commands are now recorded in [foundation choices](s00-foundation.md), which supersedes the earlier unselected-toolchain text. In particular, the npm choice is approved, tool versions are locked, and first-party Go package selection excludes npm dependency source files.

Go: exact supported patch version locked and explicitly checked; module identity is Board's. Standard library only in S00; any later module needs a specific purpose. Module consistency checks run against the implemented tree. Native/cross target and CGO=0 build choices are recorded in foundation/evidence; race uses a native compiler with CGO=1. Runtime independence is verified, not inferred from the language.

Frontend: approved npm is pinned exactly with a committed lockfile and installed through `npm ci`. TypeScript strict null checking and bigint-capable browser targets are configured. Avoid `any`-based assertions of future fetched payloads. The S01 browser parser will recover original numeric tokens, not rounded values from standard `response.json()`.

Build chain: clean locked frontend install → type-check/lint/tests → Vite production build → generated asset inventory/precompression if adopted → Go checks requiring embed → native/cross binaries → actual embedded-asset smoke. Use package `web` with directive `dist`, never a parent traversal. The production command must build assets before Go; a Go-only convenience command must fail clearly if assets are missing. Generated output should not create unrelated source diffs or leak dev-server endpoints.

No bit-for-bit reproducibility promise is established merely by dependency locking. S00 establishes repeatable steps/inputs; if byte-identical builds are desired, compare two clean artifacts and document metadata/timestamps/compression influences. Later release acceptance is always tied to the actual resulting checksum.

## Tests and fixture preparation

Implemented S00 choice: repository-level [testdata](../testdata/README.md) is the authoritative shared numeric source, with no duplicated Go/TypeScript copies. No empty internal/health package is created. JSON examples under [docs/examples](examples/overview-current.json) remain the single source for explanatory Board envelopes; future S01 extraction must preserve raw snapshot tokens. Do not use a lossy JavaScript numeric parse/stringify cycle.

S00 defines fixture IDs/expected outcomes; S01 writes the complete cases and validators. Required taxonomy:

- Complete, partial per group, all three issue codes, firmware group and independent temperature, whole-network failure and interface-local state/RX/TX failure, true/false/zero values.
- Exact integer lexical values at safe boundary, `9007199254740993`, uint64 maximum and maximum+1; invalid quoted/fractional forms and invalid UTF-8/JSON.
- Missing-versus-null every required field, null object versus null leaves, wrong issue path/order/code, unknown schema and network state under the ratified policy.
- Transport-controlled cases planned for S02: non-200, refusal, delayed/stalled body, excessive body, deadline, cancel/late-result ordering and no retry.

Clock seam must allow independent wall and elapsed-time advancement; fake Health must expose request counts, barriers and controlled headers/body chunks. S00 may define a harness interface and a simple demonstration, but must not implement cache/resilience semantics as incidental scaffolding. Use normal source tests plus real HTTP integration where transport behavior matters. Race checks cover concurrency when introduced in S02. Frontend tests need controlled fetch/time and later lifecycle/visibility events. Plan parser fuzz/differential and integer rendering tests for S01/S03, not tests that merely restate type definitions.

## Configuration and CI handoff

Approved defaults: Board `0.0.0.0:8081`, Health existing URL on 8080. S00 implements the minimal --listen and --help surface only; remaining CLI/env precedence is a proposal for S02. No user-configurable freshness/timeout may silently break the established contract. Full configuration validation, version metadata, shutdown, admission and significant-transition logs are S02; Debian paths/env file/unit are S04.

CI distinguishes source checks and binary smoke tests; it has no green hardware placeholder. Race on a suitable native platform is not replaced by cross-build. [The workflow](../.github/workflows/ci.yml) passed its first hosted run: **35527458889, foundation and race success**, for the exact commit identified above. The [evidence record](s00-evidence.md#verified-github-ci-evidence--2026-09-20) records actual image `ubuntu-24.04` / `20260907.300.1`, tools and logs. No release credentials or deployment step is present; later commits do not inherit this run's validation.

## Explicit sprint boundary

S00 ends with a reproducible build shell and validation scaffolding. It does **not** deliver a validated Health decoder/parser (S01), live provider fetch or overview/cache/recovery (S02), metric dashboard/freshness UI (S03), Debian/systemd artifacts or performance/physical acceptance (S04). A preliminary small shell bundle does not pass the final frontend budget gate. No Health code or deployment is changed at any sprint without its own scope; this project remains a consumer.

## Task execution ledger

| Task | Actual S00 status | Evidence / outstanding condition |
|---|---|---|
| S00-01 | DONE | Updated decision register, ADR-001, module identity and foundation/browser choices. |
| S00-02 | DONE | Exact tool checks, compatible dependency pins/lockfile, explicit no-auto-switch Go runner. |
| S00-03 | DONE | Only cmd, internal/httpui and web Go packages; no future empty packages or Go external dependencies. |
| S00-04 | DONE | Minimal React shell; compiled dist embedded from web; standalone binary byte checks. |
| S00-05 | DONE | Local native/cross builds and clean-source/fresh-cache repeat; CI-03 adds hosted native/cross evidence. See exact comparison scope in evidence. |
| S00-06 | DONE | Real shell/asset tests and byte-preserving shared fixtures; no fake provider/cache framework. |
| S00-07 | DONE | Minimal shell flags documented separately from future full configuration. |
| S00-08 | DONE; CI PASS | Run 35527458889: foundation/race success; CI-01–CI-05 and actual hosted image recorded in evidence. |
| S00-09 | DONE | Local and CI S00 A-01–A-04 subset, asset inventory and preliminary shell-only size report. |
| S00-10 | DONE; EXIT CLOSED | Dated exit review identifies tested commit, local/CI provenance and every applicable criterion; S01 ready, not started. |

The [original documentation review](documentation-review.md) is historical; [S00 evidence](s00-evidence.md) records implementation and the subsequent documentary closure. Interactive visual checks and full browser/mobile validation remain NOT RUN. Full Class A product checks, functional B and physical C acceptance are not claimed by completing S00 foundation checks locally and in CI. No release or deployment was performed.
