# S00 implementation choices and development commands

Scope: buildable development shell only. The user approved the stack/module/asset layout/default port and [exact-integer convention](adr-001-exact-json-integers.md). This document records ordinary foundation choices made during S00; it does not approve P-02 through P-06 wholesale.

## Toolchains and dependencies

Checked on 2026-09-20 against the official [Go downloads](https://go.dev/dl/?mode=json), [Go support policy](https://go.dev/doc/devel/release#policy), [Node release index](https://nodejs.org/dist/index.json) / [release policy](https://nodejs.org/en/about/previous-releases) and npm registry package manifests/peer dependencies.

| Component | Exact version | Rationale / license |
|---|---|---|
| Go | 1.27.1 | Current supported stable release, also installed locally; BSD-3-Clause |
| Node.js | 24.21.0 | Supported LTS line, development/build only; MIT (bundled third-party licenses apply) |
| npm | 11.19.0 | Version bundled with that exact Node archive, checked explicitly; Artistic-2.0 |
| React / react-dom | 19.3.0 | S00 direct production npm dependencies; MIT |
| Vite | 8.3.0 | Production build; Node engine compatible with 24.21.0; MIT |
| TypeScript | 6.0.3 | Strict type-check; current typescript-eslint peer range is <6.1, so TS 7.0.2 was not selected; Apache-2.0 |
| @types/react / @types/react-dom | 19.3.0 | Matching React types; MIT |
| @types/node | 24.13.6 | Node 24 tool/test configuration types; MIT |
| ESLint / @eslint/js | 10.11.0 / 10.0.1 | Lint application and build scripts; MIT |
| typescript-eslint | 8.70.0 | TS parser/rules, compatible with selected TS/ESLint; MIT |
| eslint-plugin-react-hooks | 7.1.1 | React correctness rules; MIT |
| globals | 17.12.0 | Explicit browser/Node lint environments; MIT |
| Prettier | 3.9.8 | Source/config formatting; MIT |
| Vitest | 5.0.1 | TSX component tests; peers include Vite 8 and Node 24; MIT |

Exact direct versions, transitive resolutions and integrity hashes are in [package.json](../web/package.json) and [package-lock.json](../web/package-lock.json). S01 adds exactly `lossless-json` 4.3.1 (MIT, no production transitive dependencies); see [its evaluation](sprint-01.md#parser-decision-and-primary-source-evaluation). The S00 toolchain pins are unchanged. No Go third-party module is required, so there is no go.sum to fabricate. No React Vite plugin is needed for this small shell; the TSX build uses Vite's configured transform. Component tests use React's static render **only inside tests**; the production application has no SSR.

[toolchains.json](../toolchains.json) is the explicit tool identity record. The runner checks the actually executing Node/npm/Go versions, not just declared engines. It forces `GOTOOLCHAIN=local`, `GOENV=off`, `GOWORK=off`, clears ambient Go target/flags for native work, and uses CGO=0 except the native race task. An unnecessary same-version toolchain directive was removed after `go mod tidy -diff`; `go 1.27.1` alone is not treated as an exact executable-version lock. CI pins the same versions and fails on a mismatch instead of upgrading implicitly.

## Browser targets and S01 parser evaluation

Build/runtime targets: **Chrome ≥111, Edge ≥111, Firefox ≥115 (including that ESR floor), Safari/iOS Safari ≥16.4**. These are compatibility floors, not advice to keep an unmaintained browser. Prefer an updated supported browser within these families. Vite targets are explicit in [vite.config.ts](../web/vite.config.ts), not an evolving default. TypeScript uses ES2022, strict null checks, exact optional properties and unchecked-index protection. The targets provide native bigint; no bigint runtime polyfill or external CDN is planned. See [Vite production browser support](https://vite.dev/guide/build#browser-compatibility).

S00 fixes these targets and the evaluation direction only. S01 must verify any native JSON original-token facility on **every minimum target** before relying on it. The existence of a token-source API in a current browser does not establish compatibility at these floors. Prefer a focused lossless library evaluation where necessary; compare raw token access, schema-driven conversion, uint64 bounds, invalid input, duplicate keys, resource limits, license/maintenance and actual bundle delta. Do not infer type from whether a floating metric's token has a decimal point. No general-purpose owned parser is authorized in S00.

## Minimal implementation and configuration

The S00 shell contains `cmd/joy-pi-board`, `internal/httpui` and the `web` embed package. S01 adds the independent `internal/healthschema` contract package; it is not wired into the shell. The binary serves exact `/` and actual `/assets/<filename>` GET/HEAD; unknown paths (including `/api/v1/overview`) are 404, never SPA fallback. Invalid methods receive 405; query/body on served shell resources receives 400. These are ordinary S00 choices, not the future overview error-envelope implementation. S00 errors are plain text; functional Board API rules remain S02.

Assets are built into `web/dist`, ignored by Git, and embedded with `//go:embed dist`. The binary has no filesystem fallback. Index and assets currently use no-store; production cache/compression policy is deferred. HTML has the proposed same-origin CSP already needed by this shell; responses have nosniff/no-referrer. JS/CSS are external local assets. No remote resource or runtime Node process is used. Offline gzip-9 size reports are **measurements**, not a claim of gzip delivery by the S00 server (which serves identity bytes).

The shell supports only `--listen` (default `0.0.0.0:8081`, `127.0.0.1:0` useful for isolated tests) and standard flag `--help`. No `--health-url`, environment precedence, `--version`, build-metadata injection, production admission or graceful signal lifecycle is implemented. These remain S02, with the exact target Health default `http://127.0.0.1:8080/v1/snapshot` unchanged. The process has no Health client, timer or cache. Startup prints one listening line; requests do not produce success logs. Header read/idle timeouts are 2/15 s for the minimal server; they do not establish the full P-02 production transport policy.

## Commands

Run from the repository root with the exact tools on PATH. Installation may need Internet; binary operation does not. If the environment contains another Go/Node/npm, install the pinned version outside the repository or use an ignored local tool directory; do not change the lock to match whichever version is already installed.

```text
node scripts/tasks.mjs tools
npm --prefix web ci
npm --prefix web run check
npm --prefix web run build
npm --prefix web run cross
npm --prefix web run smoke
npm --prefix web run missing-dist
npm --prefix web run sizes
npm --prefix web run race
```

`check` runs formatting, strict types, lint, component/contract tests, fixture-byte integrity, Vite production build, Go formatting/module consistency/vet/tests, then S01 Go-to-TypeScript interop and dedicated parsing bundle measurement. `contract` runs that last stage explicitly; `fuzz` adds ten seconds of bounded Go fuzzing. `build` and `cross` each rebuild frontend first. `race` requires a native Go-supported platform and C compiler and also builds frontend first. It fails if that environment is unavailable; do not relabel a cross-build as race testing. Build outputs are `out/joy-pi-board.exe` (Windows) or `out/joy-pi-board` (other native OS), and `out/joy-pi-board-linux-arm64`. `smoke` uses the native binary, empty working directory apart from the copied binary, empty PATH, byte comparisons and a local port-8080 traffic trap when available. The trap supplies no Health data. `missing-dist` creates a separate source copy under out, omits dist and asserts the actual Go embed compilation failure; it does not move/delete the working source tree.

The first-party Go package selection is `./cmd/... ./internal/... ./web`. A naive `go test ./...` also discovers Go source shipped inside some npm dependencies; the runner deliberately scopes checks to Board-owned packages. Direct Go commands requiring embed must follow `npm --prefix web run build:frontend`. `gofmt -w` formats changed Go files; `npm --prefix web run format` formats frontend/build configurations. `npm --prefix web run dev` starts local Vite development only; no Health/API proxy is needed before S02.

To run the built shell manually, use `out/joy-pi-board --listen 127.0.0.1:8081` on Unix or `.\out\joy-pi-board.exe --listen 127.0.0.1:8081` in PowerShell. Stop this development process normally; production graceful-shutdown acceptance has not been run. No installation, system service or release packaging is provided.

## CI and repeatability

The current [CI workflow](../.github/workflows/ci.yml), renamed `Board foundation and contracts` for S01, uses `ubuntu-24.04`, a foundation job and a separate native race job. Official checkout/setup-go/setup-node/upload-artifact actions are locked to reviewed full commit SHAs with their release tags as comments; action metadata was checked to use Node 24. Go and npm caches are not restored. The workflow records actual ImageOS/ImageVersion, kernel/OS and executing tools. The OS label is explicit but the hosted image contents are **not immutable**. No custom runner infrastructure is introduced.

CI checks sources, clean locked install, embedded binaries, ARM64 ELF/no interpreter, missing-dist failure, sizes and clean tracked diff. Artifacts/logs are short-lived CI evidence, not a release; tokens are read-only and checkout credentials are not persisted. The first hosted run, 35527458889, passed foundation and race on 2026-09-20 for `b6e139c3378cd77fdcb4c5edc610977c08bbf72a`. See [evidence](s00-evidence.md#verified-github-ci-evidence--2026-09-20) for verified logs, actual image identity and separate local/CI artifact provenance. This historical run closes the S00 CI gate, not full product acceptance or the later S01 changes. S01 adds bounded fuzzing and contract evidence to the workflow; its hosted status is in [S01 evidence](s01-evidence.md).

Repeatability means locked inputs and documented commands. `-trimpath` and `-buildvcs=false` remove path/VCS variability from S00 binaries; final release provenance remains S04. Byte equality is claimed only for the specific repeated-build comparison recorded in evidence, never for every machine, future toolchain or hosted image. No S00 measurement accepts the full dashboard or hardware budgets.
