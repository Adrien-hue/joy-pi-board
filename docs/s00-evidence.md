# Sprint 00 execution evidence

Local implementation and CI closure review date: 2026-09-20. Repository: `Adrien-hue/joy-pi-board`. Local implementation started from clean documentary baseline `98503f2`; the scaffold was subsequently committed and pushed as `b6e139c3378cd77fdcb4c5edc610977c08bbf72a`. This is not an official release candidate. Health code, documents, ports and deployment were not modified; no Health internal packages were imported.

**Conclusion: S00 COMPLETE**, supported by the preserved local checks and the first successful hosted CI run verified below. S00-08's hosted-run criterion and S00-10's exit review are satisfied for the identified scaffold commit. **S01 is ready to start, NOT STARTED**; S02–S04 have not started. Physical acceptance and functional contract/end-to-end validation remain NOT RUN.

## Local inputs and environments

- Windows amd64: Go 1.27.1, Node 24.21.0, npm 11.19.0, checked from the running executables by the project runner. Node was installed only in ignored `.tools`, from the official archive; Windows archive SHA-256 `158f7685b44de51f6c0df1d153526cbcd3e1bc739a8dfc607721cef75de9e541` matched the vendor checksum list.
- Go compilation uses `GOTOOLCHAIN=local`, `GOENV=off`, `GOWORK=off`, CGO=0 for native/ARM64 binaries, `-trimpath -buildvcs=false`. No external Go modules/go.sum. Fresh task-specific Go caches were used for the repeat.
- Linux native race environment: Docker Desktop engine Linux/amd64, official `golang:1.27.1-bookworm` resolved to `sha256:69a7b9788769bec032d238959b61854e9ae87f57be9029ec04e9885fabf99195`. Go 1.27.1; GCC Debian 12.2.0-14+deb12u1; Linux `6.18.33.2-microsoft-standard-WSL2`, x86_64. Network disabled; repository mounted read-only; container/compiler caches initially empty; CGO=1. This is a native amd64 race run, not an ARM64 cross-run or Pi measurement.
- [Tool/dependency/browser choices](s00-foundation.md) record compatibility and licenses. Lockfile SHA-256: `7491e33c673bb2d4e4bea7521d7cdc76ae6b3dc397361eb295a81b87cf729b35`. [Input fingerprints](evidence/s00-inputs.json) record exact source, fixture and configuration hashes for review; documentation is excluded from this non-release fingerprint.

## Local executed checks

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

## Local artifact and asset proof

These are **local S00 build artifacts**, not release packages. No rebuild after the comparison is assumed accepted merely because its filename/version matches.

| Artifact | Bytes | SHA-256 (same in original and clean-repeat build) |
|---|---:|---|
| `out/joy-pi-board.exe` | 9,047,040 | `c75c48cc1558fa1ecdafba017fd2f20bcf3fd502d48efb2be9c2cc4ba626069d` |
| `out/joy-pi-board-linux-arm64` | 8,297,908 | `d44a4765b719a7669d709d0faaed7b5a5d4d907f57dc8ca2c50fbd0cde1e6c0f` |

Recorded machine-readable [standalone proof](evidence/standalone-smoke.json) and [shell size inventory](evidence/shell-sizes.json) include hashes of the actual embedded HTML and compiled assets. Identity payload bytes: HTML 475, JS 220,427, CSS 742; total 221,644. Offline gzip level 9 with Node 24.21.0 / zlib `1.3.2.1-motley-8002e91`: JS+CSS **68,495 bytes (66.89 KiB)**; all three resources **68,790 bytes (67.18 KiB)**. No source maps, CDN/font/script requests or overview payload belong to this minimal shell inventory. Server compression is not implemented; reported gzip sizes are computed from served asset bytes, not asserted wire Content-Encoding.

The standalone subprocess has no Node runtime PATH and no source/dist in its working directory. Tests compare delivered payloads with build output, while the process uses only its embedded filesystem. Source review confirms no Health client or polling goroutine and no host collection. The six-second trap observation is supporting runtime evidence, not a proof of future resilience or the full idle CPU budget.

## Verified GitHub CI evidence — 2026-09-20

This is new remote evidence, verified directly through GitHub's run, jobs, artifact metadata and both job logs. It supersedes the earlier lack of a hosted run; it does not replace the local evidence above or erase its corrected diagnostic failures.

| Identity | Verified value |
|---|---|
| Repository / tested commit | `Adrien-hue/joy-pi-board` / [`b6e139c3378cd77fdcb4c5edc610977c08bbf72a`](https://github.com/Adrien-hue/joy-pi-board/commit/b6e139c3378cd77fdcb4c5edc610977c08bbf72a) |
| Commit message | `build: establish Sprint 00 application scaffold and CI` |
| Workflow / file | `S00 foundation` / `.github/workflows/ci.yml` |
| Run ID / sequence / attempt | [35527458889](https://github.com/Adrien-hue/joy-pi-board/actions/runs/35527458889) / 1 / 1 |
| Event / branch | `push` / `main` |
| Created / started (UTC) | `2026-09-20T17:56:39Z` / `2026-09-20T17:56:39Z` |
| Last updated (UTC) | `2026-09-20T17:57:49Z` |
| Run status / conclusion | `completed` / `success` |
| Foundation job | [106122010474](https://github.com/Adrien-hue/joy-pi-board/actions/runs/35527458889/job/106122010474): `completed` / `success` |
| Race job | [106122010325](https://github.com/Adrien-hue/joy-pi-board/actions/runs/35527458889/job/106122010325): `completed` / `success` |

At this review, both local HEAD and GitHub `main` were exactly the tested commit; the local tree was clean before this documentation update. The run validates that commit, not these later documentation edits or any future HEAD. No new build, CI dispatch, commit, push, release or deployment was performed during this closure review.

### Actual hosted environment

Both job setup logs identify **Ubuntu 24.04.5 LTS**, image **`ubuntu-24.04`**, image version **`20260907.300.1`**, Actions runner **`2.337.0`** and image provisioner **`20260828.587`**. The logs link the [image inventory](https://github.com/actions/runner-images/blob/ubuntu24/20260907.300/images/ubuntu/Ubuntu2404-Readme.md) and [image release](https://github.com/actions/runner-images/releases/tag/ubuntu24%2F20260907.300). This records the executed image; the `ubuntu-24.04` label does not make future hosted contents immutable.

Both jobs explicitly passed the executing-tool check: Go **1.27.1**, Node **24.21.0**, npm **11.19.0**, with `GOTOOLCHAIN=local`. The foundation environment step writes its environment inventory to `out/runner-image.txt`; its successful execution and image identity were verified in logs, not by downloading that file. The race log additionally prints `ImageOS=ubuntu24`, `ImageVersion=20260907.300.1`, `RUNNER_OS=Linux`, `RUNNER_ARCH=X64`, kernel `6.17.0-1022-azure` / x86_64 and GCC `13.3.0` (`Ubuntu 13.3.0-6ubuntu2~24.04.1`). Native Linux amd64 race checks passed; this is not Raspberry Pi execution.

### Executed CI scope and provenance

| CI evidence ID / gate subset | Verified commands and outcome |
|---|---|
| CI-01 / A-01–A-03 tools and install | Both jobs passed `node scripts/tasks.mjs tools` and `npm --prefix web ci` with separate runner-temporary npm caches; no restored Go/npm dependency cache configured. |
| CI-02 / A-01–A-02 source checks | Foundation passed `npm --prefix web run check`: installed-tree consistency, formatting, strict types, lint, real component/fixture-byte tests, frontend build, Go formatting/module checks/vet/tests. |
| CI-03 / A-03 builds | Foundation passed `npm --prefix web run build` and `cross`, frontend first. Logs identify native x86-64 and ARM aarch64 statically linked ELF binaries; ARM64 interpreter check passed. |
| CI-04 / A-04 shell subset | Foundation passed `smoke`, `missing-dist`, `sizes` and `git diff --exit-code`. Logs record standalone embedded assets, zero Health-trap connections through startup/requests/6 s idle and clear missing-dist compilation failure. Offline gzip-9 is 68,495 bytes JS+CSS / 68,790 bytes all initial resources, for the minimal shell only. |
| CI-05 / A-01 race | Race passed `npm --prefix web run race`, building frontend before native Go race tests; cmd/joy-pi-board and internal/httpui passed without a race report. |

Foundation build logs report these **CI-generated binary** SHA-256 values:

| CI output | SHA-256 reported by build log |
|---|---|
| `out/joy-pi-board` (Linux amd64) | `441ba5f99871146efb391e1f0d778c43a302d857b255dda0cb080a93543349d9` |
| `out/joy-pi-board-linux-arm64` | `d44a4765b719a7669d709d0faaed7b5a5d4d907f57dc8ca2c50fbd0cde1e6c0f` |

GitHub lists the successfully uploaded [s00-foundation-evidence artifact](https://github.com/Adrien-hue/joy-pi-board/actions/runs/35527458889/artifacts/10610118517), ID `10610118517`, 9,612,938 bytes, created `2026-09-20T17:57:46Z`, expiry `2026-09-27T17:57:44Z`, not expired when inspected. Its archive digest is `sha256:7004387e211778fdcf76fa903cd34c809c7353288552f78d8bf552f151138d15`; this is **not a binary checksum**. Artifact metadata and upload logs were inspected; archive contents were not downloaded or independently rehashed. The CI ARM64 hash reported in the log matches the previously recorded local ARM64 hash by string comparison only. This is not a fresh comparison of downloaded artifact bytes or a general cross-machine reproducibility claim. The local Windows binary and all local evidence JSON retain their original provenance.

## S00 exit review

| Applicable criterion | Closure evidence / result |
|---|---|
| S00-01–S00-03, S00-07: decisions, exact inputs, minimal structure and configuration boundary | Existing task ledger, foundation choices, source review and L-01–L-03/L-11; CI-01–CI-02 confirm executing tools and checks. Satisfied. |
| S00-04–S00-06: real embedded shell, clean repeat, useful tests and shared fixture preparation | L-03–L-06/L-08–L-10 plus CI-02–CI-05. The two clean local builds and their scoped hash comparison remain the repeatability evidence. Satisfied. |
| S00-08: successful clean hosted CI | Run 35527458889, both jobs success, actual hosted image recorded above. The last missing execution gate is satisfied. |
| S00-09: foundation A subset, runtime inventory and shell-only sizes | Local artifact proof plus CI-03–CI-04. Satisfied within the minimal-shell scope. |
| S00-10: reconcile evidence, statuses, traceability and S01 handoff | This dated review and linked sprint/roadmap/validation updates close S00. Satisfied; no further mandatory S00 criterion remains open. |

The NOT RUN items below belong to later product acceptance, not an unfulfilled S00 hosted-run gate. Closure does not approve deferred P-xx choices or start S01. Documentary closure checks: **PASS**, 17 tracked Markdown files / 132 internal links and anchors checked, statuses and traceability reconciled, `git diff --check` clean. Only eight Markdown files changed; historical local evidence sections, corrected incidents, acceptance limitations and evidence JSON were checked unchanged. No application build was rerun or historical local check relabeled as newly executed.

## Browser and product acceptance limitations

Interactive browser check: **NOT RUN**. The browser-control tool reported the in-app browser unavailable and no connected browsers. A temporary loopback shell process was started for that attempt and stopped afterward. Component rendering and actual binary HTTP/asset checks passed, but no browser screenshot, complete browser-family compatibility run or mobile visual acceptance is claimed.

Full A-04 final-dashboard/resource gate, A-05 production configuration/lifecycle, A-06 packaging, all functional B-01–B-11 and physical C-01–C-06 remain **NOT RUN**. No release/deployment/systemd/Debian task was executed. Readiness/CPU/RSS/latency/shutdown product budgets were not inferred from this workstation or the Linux test container.

## S01 handoff (ready, not started)

- Apply adopted ADR-001: uint64/original json.Number with checked bounds in Go, bigint for every TypeScript uint64 field, number for floating measures and distinct null/presence.
- Evaluate library/native original-token parsing against the fixed browser floors; no unused parser is present. Ratify unknown-member compatibility and duplicate-key/defensive policies separately.
- Consume the one shared byte-preserving fixture source, including Health's exact integer fixture; build the full valid/invalid snapshot/issue corpus without consumer copies.
- Implement the actual validator and browser parsing tests in S01. Health HTTP client/cache is S02; exact end-to-end display is S03. Do not relabel S00 fixture-byte checks as those validations.
