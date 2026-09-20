# Sprint 00 execution evidence

Date: 2026-09-20. Repository: `Adrien-hue/joy-pi-board`, starting clean commit `98503f2` (documentary baseline). These changes are an **uncommitted/unpushed S00 working tree**, not an official release candidate. Health code, documents, ports and deployment were not modified; no Health internal packages were imported.

**Conclusion: S00 is implemented and locally verified; its exit remains OPEN because no successful GitHub Actions run exists for this change.** S00-08's hosted-run completion criterion and S00-10's final closure are pending. S01–S04 have not started. Physical acceptance and functional contract/end-to-end validation remain NOT RUN.

## Inputs and environments

- Windows amd64: Go 1.27.1, Node 24.21.0, npm 11.19.0, checked from the running executables by the project runner. Node was installed only in ignored `.tools`, from the official archive; Windows archive SHA-256 `158f7685b44de51f6c0df1d153526cbcd3e1bc739a8dfc607721cef75de9e541` matched the vendor checksum list.
- Go compilation uses `GOTOOLCHAIN=local`, `GOENV=off`, `GOWORK=off`, CGO=0 for native/ARM64 binaries, `-trimpath -buildvcs=false`. No external Go modules/go.sum. Fresh task-specific Go caches were used for the repeat.
- Linux native race environment: Docker Desktop engine Linux/amd64, official `golang:1.27.1-bookworm` resolved to `sha256:69a7b9788769bec032d238959b61854e9ae87f57be9029ec04e9885fabf99195`. Go 1.27.1; GCC Debian 12.2.0-14+deb12u1; Linux `6.18.33.2-microsoft-standard-WSL2`, x86_64. Network disabled; repository mounted read-only; container/compiler caches initially empty; CGO=1. This is a native amd64 race run, not an ARM64 cross-run or Pi measurement.
- [Tool/dependency/browser choices](s00-foundation.md) record compatibility and licenses. Lockfile SHA-256: `7491e33c673bb2d4e4bea7521d7cdc76ae6b3dc397361eb295a81b87cf729b35`. [Input fingerprints](evidence/s00-inputs.json) record exact source, fixture and configuration hashes for review; documentation is excluded from this non-release fingerprint.

## Executed checks

Commands below actually ran. Except the noted first diagnostic failures, the final outcomes apply to the corrected files represented by the source fingerprints. `PASS` here is scoped to S00 foundations, never a full product A/B/C claim.

| Local ID / planned gate subset | Command or procedure actually executed | Result |
|---|---|---|
| L-01 / A-01–A-03 preparation | `node scripts/tasks.mjs tools`; npm `ci` preinstall repeats exact checks | PASS: Go 1.27.1 / Node 24.21.0 / npm 11.19.0, no automatic Go switching. |
| L-02 / A-02 | Generate package lock, then `npm --prefix web ci --cache .cache/npm-cold-install` | PASS: 177 installed packages from a fresh cache; no unused parser dependency. |
| L-03 / A-01–A-02 | `npm --prefix web run check` | PASS: npm installed-tree consistency, Prettier, strict TypeScript, ESLint, one real TSX component test, two shared-fixture integrity tests, Vite build, gofmt, `go mod tidy -diff`, `go mod verify`, vet and first-party Go tests. |
| L-04 / A-03 | `npm --prefix web run build`; `npm --prefix web run cross` | PASS: Windows/amd64 native and Linux/ARM64 binaries; frontend built before both Go compilations. |
| L-05 / A-04 | `npm --prefix web run smoke` | PASS: copied native binary in isolated directory with empty PATH served HTML/JS/CSS byte-for-byte; HEAD works, absent overview/unknown route give 404. A listening trap at 127.0.0.1:8080 saw zero connections through startup, requests and 6 s idle. |
| L-06 / A-04 | `npm --prefix web run missing-dist` | PASS: isolated copy lacking web/dist fails compilation with `pattern dist: no matching files found`; no dummy assets/fallback. |
| L-07 / A-04 preliminary | `npm --prefix web run sizes` | PASS for minimal-shell inventory only; offline gzip-9 measurements below, not final UI acceptance or actual server compression. |
| L-08 / A-01 race | `docker run --rm --network none --mount type=bind,source=C:\src\joy-pi-board,target=/src,readonly -w /src -e GOTOOLCHAIN=local -e CGO_ENABLED=1 golang@sha256:69a7b9788769bec032d238959b61854e9ae87f57be9029ec04e9885fabf99195 go test -race -count=1 ./cmd/... ./internal/... ./web` | PASS, no race report; shell/embedded assets already built before Go tests. |
| L-09 / A-03 runtime format | Inspect ARM64 ELF header/program headers using Python struct | PASS: ELF64 little-endian, machine 183 (AArch64), no PT_INTERP dynamic-loader dependency. Does not prove hardware execution. |
| L-10 / A-01–A-04 repeatability | Copy tracked + nonignored new source files to a new out/clean-s00 tree, without node_modules/dist; fresh npm/Go caches; run `ci`, `check`, `build`, `cross` there; compare SHA-256 | PASS: second clean process succeeds; both binaries identical to first builds for their respective targets. Same host/toolchains; no cross-machine universal reproducibility claim. |
| L-11 / tool enforcement | Isolated copy of runner declares Go 1.27.0, while actual Go is 1.27.1 and ambient GOTOOLCHAIN=auto | PASS: nonzero exit, `go: expected 1.27.0, got 1.27.1; automatic toolchain switching is disabled`. No toolchain auto-download. |
| L-12 / documentary consistency | Local link/anchor, decision/status/traceability and source-boundary review; `git diff --check` | PASS after S00 updates; shared fixture integrity is distinct from product schema validation. |

The first type-check identified missing Vite CSS import types; adding `vite/client` fixed it. The first module check requested removal of the redundant same-version toolchain directive; exact executable checks remain mandatory. A first broad Go `./...` run discovered Go files shipped inside npm's flatted dependency. First-party package selection was corrected, then rechecked in the clean repeat and race run. These corrected failures are not omitted or misrepresented as initial passes.

## Artifact and asset proof

These are **local S00 build artifacts**, not release packages. No rebuild after the comparison is assumed accepted merely because its filename/version matches.

| Artifact | Bytes | SHA-256 (same in original and clean-repeat build) |
|---|---:|---|
| `out/joy-pi-board.exe` | 9,047,040 | `c75c48cc1558fa1ecdafba017fd2f20bcf3fd502d48efb2be9c2cc4ba626069d` |
| `out/joy-pi-board-linux-arm64` | 8,297,908 | `d44a4765b719a7669d709d0faaed7b5a5d4d907f57dc8ca2c50fbd0cde1e6c0f` |

Recorded machine-readable [standalone proof](evidence/standalone-smoke.json) and [shell size inventory](evidence/shell-sizes.json) include hashes of the actual embedded HTML and compiled assets. Identity payload bytes: HTML 475, JS 220,427, CSS 742; total 221,644. Offline gzip level 9 with Node 24.21.0 / zlib `1.3.2.1-motley-8002e91`: JS+CSS **68,495 bytes (66.89 KiB)**; all three resources **68,790 bytes (67.18 KiB)**. No source maps, CDN/font/script requests or overview payload belong to this minimal shell inventory. Server compression is not implemented; reported gzip sizes are computed from served asset bytes, not asserted wire Content-Encoding.

The standalone subprocess has no Node runtime PATH and no source/dist in its working directory. Tests compare delivered payloads with build output, while the process uses only its embedded filesystem. Source review confirms no Health client or polling goroutine and no host collection. The six-second trap observation is supporting runtime evidence, not a proof of future resilience or the full idle CPU budget.

## CI, browser and acceptance limitations

The GitHub Actions workflow is created and pins official actions to immutable SHAs. It specifies Ubuntu 24.04, records the actual hosted image identity, and separates native race checks. **GitHub CI: NOT RUN.** No change was committed/pushed, no workflow was remotely dispatched, and no green run/log URL is available. A workflow file plus local checks does not satisfy S00-08's hosted-run exit criterion. After normal code review/publication of the change, its first real successful run and recorded image identity are needed to close S00; this task does not publish the change just to manufacture that evidence.

Interactive browser check: **NOT RUN**. The browser-control tool reported the in-app browser unavailable and no connected browsers. A temporary loopback shell process was started for that attempt and stopped afterward. Component rendering and actual binary HTTP/asset checks passed, but no browser screenshot, complete browser-family compatibility run or mobile visual acceptance is claimed.

Full A-04 final-dashboard/resource gate, A-05 production configuration/lifecycle, A-06 packaging, all functional B-01–B-11 and physical C-01–C-06 remain **NOT RUN**. No release/deployment/systemd/Debian task was executed. Readiness/CPU/RSS/latency/shutdown product budgets were not inferred from this workstation or the Linux test container.

## S01 handoff (prepared, not started)

- Apply adopted ADR-001: uint64/original json.Number with checked bounds in Go, bigint for every TypeScript uint64 field, number for floating measures and distinct null/presence.
- Evaluate library/native original-token parsing against the fixed browser floors; no unused parser is present. Ratify unknown-member compatibility and duplicate-key/defensive policies separately.
- Consume the one shared byte-preserving fixture source, including Health's exact integer fixture; build the full valid/invalid snapshot/issue corpus without consumer copies.
- Implement the actual validator and browser parsing tests in S01. Health HTTP client/cache is S02; exact end-to-end display is S03. Do not relabel S00 fixture-byte checks as those validations.
