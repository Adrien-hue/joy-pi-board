# Sprint 04 software evidence

Current update (2026-10-20): **exact-source software/archive CI VERIFIED; lifecycle and acceptance OPEN**. See [green CI and VM preparation](#2026-10-20--verified-green-ci-and-vm-preparation). Earlier dated results below remain historical; no S04 closure or release GO is granted.

Status: **SOFTWARE IMPLEMENTED; VALIDATION INCOMPLETE; NOT CLOSED; release NO-GO**. Date: 2026-10-18. The user's current instruction approves S04-R01–R04 and S04-01–S04-05 software and isolated guest tests, replacing documentary-only restrictions. No primary-host/Pi package installation, account/service creation or operation, Pi connection, network change, commit, push, tag, release or deployment was performed.

Update 2026-10-19: the [first S04 CI run and packaging diagnostics fix](#2026-10-19--first-s04-ci-and-packaging-failure-diagnostics) below supersede the earlier pending-CI statements for commit 4fd1ebc0 only. Original local evidence and limitations remain historical.

The [subsequent CI diagnosis](#2026-10-19--static-package-shared-library-scan) identifies the Debian failure for commit 491204eb and records a targeted recipe correction; its complete package build remains pending.

The [Lintian follow-up](#2026-10-19--lintian-policy-for-the-frozen-static-executable) records a constructed archive for commit a2cca6e6, rejected by Lintian; complete validation of the follow-up changes remains pending.

## Source and preserved history

Repository identity: `Adrien-hue/joy-pi-board`. Initial tree clean at `fcbb096e71bf6178145b7f370f5036b67a5e53eb`, the previously committed S04 preparation. Work in this assignment is uncommitted and dirty, with exact inputs recorded by development provenance. There is no new commit SHA or remote CI result to attribute to these changes.

Preserve S00/S01/S02 evidence, the S03 scheduling failure/correction, corrected successful run **36269739581 attempt 1** for **24d527a301d66b49d86ad11fbd75393254f02094**, and all real Chrome observations/captures. The Chrome evidence remains tied to Windows binary **6bfbbfb96e47d1ec1baf07d0a8e19e1f3d541d774c420848346b1f1ca568dcf7**. None becomes S04/Pi proof. The historical ARM64 build is likewise not an accepted candidate. All previously unexecuted S03 operator scope remains NOT RUN and explicitly transfers to S04 acceptance.

## Delivered software

| Task | Delivered | Outstanding |
|---|---|---|
| S04-01 | Approval/deferral records, retained application locks, resolved OCI base digests, dated Debian snapshot recipe, source/license inventory and honest artifact policy. | Full packaging transitive closure/derived image execution; real Maintainer and Board copyright/license. |
| S04-02 | Deterministic gzip-9, embedded public-resource representation manifest, build-time decompression/hash checks, runtime startup integrity check, q negotiation and static errors/GET/HEAD/security, private sidecars. API/upstream identity unchanged. | Final candidate and Pi encoded-byte/browser acceptance. |
| S04-03 | Explicit opt-in diagnostic flag, correlated flight/network-before-validation timing, bounded sanitized logging and loss/truncation summary; external load/raw-byte/process/statistics tools. | Physical calibration, Linux process sampling on target, paired same-binary diagnostic on/off and all physical budgets. |
| S04-04 | Debian recipe consumes supplied bytes, conffile, unprivileged account/collision policy, systemd unit, no auto enable/start, archive/script inspection and comparison tools, opt-in isolated VM lifecycle harness. | Actual Debian build, generated-script inspection, archive reproducibility/lintian and real-systemd lifecycle. No .deb produced or installed. |
| S04-05 | Source regressions, native/ARM64 development builds, provenance/checksums/notices, separate comparison build, CI extensions retaining old jobs, runbook. | Successful Linux race/engine/packaging results on future commit, clean frozen candidate, exact ARM64 runtime and acceptance. |

The diagnostic is off by default, accepts only bare `--diagnostic-timing`, rejects values/duplicates/help-version combinations and has no env/endpoint. The 64-entry nonblocking log queue accepts at most 4,096 request records per process; fixed counters expose drops/truncation. Network receipt ends before Health schema validation. Request/flight correlation excludes shared/cancelled/incomplete/unsuccessful rows. Missing/truncated output invalidates measurement. No runtime system collector, sensor access or persistence was added.

## Actual local checks

Host: Windows x64; Windows kernel family reported by browser engines as Windows NT 10.0. Exact tools: Go 1.27.1, Node 24.21.0, npm 11.19.0, Node zlib `1.3.2.1-motley-8002e91`; external Python 3.11.9. No application toolchain/dependency upgrades. CGO=0 native and Linux ARM64 v8.0 builds; the ARM64 executable was inspected, **not executed**. No real systemd was available.

| Check | Actual outcome / scope |
|---|---|
| Locked npm installation/tool versions/dependency consistency | PASS, 181 installed packages under the existing lock; no opportunistic version changes. |
| Formatting, types, lint, frontend unit tests | PASS, 591 frontend tests including S03 early-callback scheduling regression. |
| Go tests/vet/package dependency boundary | PASS; S01/S02 regression plus static q/header/error/integrity, diagnostic parsing/default silence/correlation/shared/cancelled/limits/blocked-sink and main integration checks. Linux-only signal assertions retain their platform condition; no Linux execution inferred. |
| S01 differential corpus | PASS: 264 corpus cases, 76 accepted Go outputs parsed in TypeScript, 512 seeded mutations. |
| S02 HTTP interoperability | PASS: 266 actual controlled Health attempts, 78 valid/188 invalid; max overview 65,738 bytes/depth 34. |
| Bounded Go fuzzing | PASS: 195,033 executions in the recorded run, two workers, five new interesting inputs; this is local parser evidence. |
| Native/ARM64 build, standalone smoke, missing-dist | PASS; gzip representations embedded before Go, standalone UI with Health unavailable. Cross-build alone does not prove ARM64 execution. |
| Measurement tool tests | PASS on controlled HTTP/data: nearest-rank statistics, missing/corrupt/lost/truncated diagnostic rejection, raw response bound, error counts and per-client results, Python PID continuity/statistics/syntax. No physical PASS. |
| Actual served resource sizes | PASS locally for current native bytes: 79,522 gzip JS+CSS bytes; total initial 80,762 complete / 81,997 partial / 145,541 maximum-markup bytes, each including first overview **identity**. Reports retain encoded bodies as base64, headers, sizes and hashes. |
| Native Windows Playwright 1.63.0 exploratory production run | **22 PASS / 5 FAIL**, all 27 attempted, retained without weakening assertions. Not supported-Linux acceptance or actual Safari/iOS. See incidents below. |
| Linux race detector | NOT RUN for S04: no working Linux execution environment. Existing CI job retained for the future commit; no prior S03 race result reused. |
| Debian build/actual generated scripts/lintian/reproducibility | NOT RUN: Docker 29.7.2 client present but Linux engine pipe absent. No service was started to make it available. |
| Real-systemd lifecycle / ARM64 guest / Pi / operator clients | NOT RUN. Archive inspection, when available, will not satisfy these lanes. |
| Remote CI | PENDING for a future committed SHA. No push/run/tag was created. |

Raw logs are retained in `out/s04/`; raw native HTTP bodies and hashes in `out/dashboard-sizes.json`; failed browser report, screenshots and traces in `out/browser-report.json` / `out/browser-results/`. A compact retained source/evidence index lives under [evidence/s04](evidence/s04/README.md). Browser failure captures are diagnostic artifacts, not reviewed visual PASS.

## Packaging lock and missing metadata

Official Debian [debhelper](https://packages.debian.org/trixie/debhelper) and [dpkg-dev](https://packages.debian.org/trixie/dpkg-dev) catalog entries identify retained versions **13.24.2** and **1.22.22**. The [dh_installsystemd reference](https://manpages.debian.org/trixie/debhelper/dh_installsystemd.1.en.html) documents separate no-enable/no-start controls. These source reads do not execute the package.

Registry metadata resolved the following real OCI identities on 2026-10-18; no digest was invented:

| Object | SHA-256 digest |
|---|---|
| Debian multi-platform index | `a29215f6a35e51e22adffa17f89e9d2ef06214e64a2bad10d765c46aea49f11f` |
| Debian linux/amd64 build base | `918311b7b6c4c6f68b232ba516584925f6c78ad82b6fd534b98979df6438e483` |
| Debian linux/arm64 base | `a2ca65377c251208de42e1a2722cb949211236c161b634f6f3663754e2d06850` |

The recipe selects signed snapshot `20261005T000000Z` and enforces dpkg/dpkg-dev 1.22.22, debhelper 13.24.2/compat 13. **Snapshot installation availability and the complete transitive package set have not been verified by an executed build.** Derived image ID/digest, adduser/systemd/binutils/lintian/Python/compressor transitives must be captured by the delivered image/build steps before closing S04-01/04. No floating tag is used as locking evidence; no replacement version is silently selected. Current status is partial lock resolution, not a proved reproducible packaging environment.

No explicitly established project distribution license/copyright or Maintainer identity/contact was found. A Git author is not a package Maintainer. Distributable packaging fails closed without supplied established metadata. Test-only fixture metadata is visibly synthetic and cannot pass for a real package identity. Third-party notices cover shipped dependencies/runtime and assert no Board license.

## Incidents and limits

1. Sandbox process spawning initially denied Go/npm and npm's initial cached-only acquisition failed. The authorized locked development installation/checks were rerun with appropriate process/network access; tool versions/lockfiles were retained.
2. A new controlled-tool test exposed oversized HTTP response handling that could resolve before rejection. `rawHTTP` now rejects immediately before destroying the request. The initial failing log is retained; the corrected test passes.
3. Exploratory Windows browser traffic was transformed by an environment intermediary: traces contain `X-Content-Encoding-Over-Network: gzip` and no Content-Encoding/Content-Length, despite direct raw Node HTTP confirming exact gzip bytes. Three new strict gzip-header assertions fail. A local proxy-bypass launch still observed the transformation. Assertions were kept unchanged; this is not claimed as a clean compressed browser run.
4. Firefox/WebKit Board-stop cases rendered “Board did not respond in time” rather than the expected immediate network-refusal wording. They still showed stale data at the captured instant; their full expiry assertions did not complete and are not PASS. Chromium's Board-stop case passed. Do not infer either Linux behavior or Safari certification from this Windows/proxy run. Re-run unchanged tests on supported Linux/direct networking.
5. Docker Linux engine was unavailable; WSL inventory showed only docker-desktop. No Linux service was started and no alternative host installation was attempted. Actual package generation, maintainer-script contents after debhelper, lifecycle, Linux race and supported-engine validation remain missing.

## Candidate and release boundary

The selected development artifact is `out/artifacts/s04-dev-20261018-c`; the separate comparison is `out/artifacts/s04-compare-20261018-c`. Exact SHA-256, input digest, sizes and comparison outcome are retained in [artifact summary](evidence/s04/artifacts.json). Both are **DEVELOPMENT / NOT ACCEPTED**, built from a dirty tree. Earlier development identities a/b and the b comparison are retained as superseded input history, not overwritten. Final tool validation also rejects missing diagnostic schema fields or invalid numeric/boolean fields; the focused tests and lint were rerun after this tool-only hardening. No Debian package, clean-source frozen candidate or accepted release exists.

For a clean candidate, first complete Linux checks and packaging metadata/environment validation, commit under later authorization, verify exact-commit remote CI, then create a fresh candidate ID/revision with the delivered recipe. Execute packaging and ARM64 runtime/real-systemd lanes, then separately authorize the named Pi campaign. Promote accepted binary/package bytes without reconstruction. Any packaging/source/asset change creates a new identity and reopens affected evidence.

S04-06–S04-10, C-01–C-06 and transferred S03-T05/T09/T11/T12/T14 operator scope remain NOT RUN. Required inputs are target/access/OS/architecture, storage/power/cooling, Health version/service state, Board collisions/configuration, downtime/restoration plan, load host and real client/device/reader availability. See [runbook](s04-runbook.md). S04 is not closed and release GO is not granted.

## 2026-10-18 final repository checks

437 local Markdown links/anchors checked with no missing target; requirement/task/gate references reviewed together. `git diff --check` passes. All 122 selected development input hashes still match the working files. Existing application toolchains/lockfiles, fixture bytes, S00/S01/S02 evidence and tracked historical evidence artifacts are unchanged. Packaging text uses LF endings. The diff is confined to authorized S04 runtime seams, tests, external tools, packaging/CI and English documentation/evidence; specification changes only the authorization row, without budget changes. No commit was created.


## 2026-10-19 — First S04 CI and packaging failure diagnostics

Inspected [run 39308202979](https://github.com/Adrien-hue/joy-pi-board/actions/runs/39308202979), attempt 1, push event, exact source `4fd1ebc0e8c825a7b573d976ccad7347d6bffad0`. The local tree was clean at that commit before this follow-up; no commit/push/rerun was performed by this assignment. GitHub reports `race` (118250650131) and `browser` (118250650311) successful; `foundation` (118250650452) passed its source/fuzz/build/smoke/measurement stages and failed the isolated archive step. Overall CI is **FAIL**, not a candidate acceptance result. These job results do not accept real browsers, devices, systemd or Pi.

The Docker image did build: actual amd64/Linux image configuration ID `sha256:f46a439676d0c7aa91356ddd6a745489ec13d0aab5a2fb1e77830b7aa21f601d` (no registry RepoDigest). The foundation log records installation of dpkg-dev 1.22.22 and debhelper 13.24.2 from the selected snapshot. Do not equate the local image ID with a published registry manifest digest or a complete retained transitive package lock.

The Python traceback exposes only `dpkg-buildpackage` exit **2**. Its detailed output went to `/work/one/build.log`. `docker start -a` then exited nonzero and GitHub's fail-fast shell stopped before `docker cp`, so the actual Debian cause was lost with the disposable runner. Downloaded artifact **12214837309**, `board-check-evidence`, SHA-256 **b5bdc046fdabbb02b4f5cb9d1f438ee68815030775a2640c890416b525186351**, matches GitHub metadata; its ZIP contains `packaging-image.json` but no packaging work directory or build.log. Root cause of the Debian build remains **UNDETERMINED**. No speculative tool/version/recipe substitution is made.

This follow-up corrects error visibility and evidence retention: the build helper preserves the log, prints original output on failure and propagates its exit code; the package inventory is captured before the build. CI captures Docker's failure status, copies `/work`, removes its own stopped test container, then fails without masking the error. A failed build will therefore supply the actual log and package versions. Two controlled Python-child tests cover success and exit-2 failure with stdout/stderr retention; no Debian/maintainer/service code executes in these tests. The actual workflow status/copy/removal fragment also passes simulated exit-0 and exit-2 cases in Git Bash with a shell-only Docker stub; no real container is used. Python tool checks, workflow formatting and `git diff --check` pass locally. [Run/job/artifact metadata](evidence/s04/ci-39308202979.json) preserves the read-only inspection; the downloaded archive and controlled workflow results remain under `out/s04-ci-fix/`.

Local Docker Linux remains unavailable, so **the Debian build itself is not fixed or revalidated yet**. Obtain the saved log from a surviving disposable container, or execute updated CI/isolated packaging using a new artifact identity (packaging inputs changed). The final 2026-10-18 input-hash statement above describes that historical delivery, not these modified packaging sources. No prior development identity is overwritten. Full package validation, generated-script execution/lifecycle, Pi and transferred operator gates remain open; no release GO.

## 2026-10-19 — Static package shared-library scan

Inspected [run 39310053644](https://github.com/Adrien-hue/joy-pi-board/actions/runs/39310053644), attempt 1, exact source `491204eb0bce909fc3bbd539f480c858e2984e0b`; the local tree was initially clean. Overall conclusion is **CANCELLED**: foundation job 118256993675 failed packaging, race job 118256994034 passed, and browser job 118256994312 was cancelled during its tests. The previous run's successful browser result remains historical and is not attributed to this run.

The logging correction worked. The retained [Debian build log](evidence/s04/ci-39310053644-build.log) identifies the cause: `dh_makeshlibs -a` tries to execute absent `aarch64-linux-gnu-objdump` on the ARM64 executable, then dpkg-buildpackage exits 2. The earlier cross-compiler warnings are not the fatal error. The [full package inventory](evidence/s04/ci-39310053644-packages.lock) confirms dpkg/dpkg-dev 1.22.22, debhelper 13.24.2 and native binutils 2.44-3; no target binutils package is installed. This provides the actual transitive environment inventory missing from the first failure, without proving archive reproducibility.

Downloaded artifact 12215442885 (`board-check-evidence`) has SHA-256 `f079c5c87ff1d936bb9fafcd6f119c30c4796b171d632ae60e00565df9bcb99b`, matching GitHub metadata. Read-only ZIP inspection confirms seven payload files, exactly one ELF (ET_EXEC, ELF64 little-endian AArch64), no PT_INTERP or PT_DYNAMIC, and no shared libraries. Its binary SHA-256 `9ddbc8896b1ee18e5ff1d9b8b37671ed1b4b3745965e71663786436adf7f5dbd` matches provenance, the input artifact, payload and debhelper-staged copy. Status remains **DEVELOPMENT / NOT ACCEPTED**. No binary or maintainer script was executed. The archive and inspection report remain under `out/s04-ci-makeshlibs/`; [the machine record](evidence/s04/ci-39310053644.json) retains inspection results and log/inventory checksums.

The recipe now has an empty `override_dh_makeshlibs`, alongside the existing shlibdeps/strip overrides. Debian documents [dh_makeshlibs](https://manpages.debian.org/trixie/debhelper/dh_makeshlibs.1.en.html) as shared-library metadata/ldconfig processing; this static executable-only payload needs neither. No tool/dependency lock changed, no binary mutation was introduced, and extracted-byte identity and lintian checks remain enforced. Future shared libraries or dynamic linking require revisiting these overrides.

Local controlled Python tool checks and both build-output retention tests pass. The changed Markdown references and `git diff --check` were checked. **The corrected Debian build, lintian and archive comparison remain NOT RUN locally**, because the Linux Docker engine is unavailable; the next exact-commit CI must confirm them with a new artifact identity. No service operation, installation, CI rerun, commit or push was performed. Real systemd lifecycle, ARM64 execution, Pi and deferred S03 operator acceptance remain NOT RUN; S04 is not closed and release remains NO-GO.

## 2026-10-19 — Lintian policy for the frozen static executable

Inspected [run 39311632646](https://github.com/Adrien-hue/joy-pi-board/actions/runs/39311632646), exact source `a2cca6e6c24c5334970f37ee76274640b2400528`; the local tree was initially clean. Foundation job 118262417341 failed at Lintian after archive construction, extraction, binary equality and generated-script syntax checks. Race job 118262417650 and browser job 118262417749 passed. Overall CI remains **FAIL**. This supersedes earlier archive-construction NOT RUN statements for this source only; it does not establish successful full packaging, comparison, installation or lifecycle validation.

Artifact 12216191475 (`board-check-evidence`) was downloaded and its SHA-256 verified as `a258af82c5c24370147a469d8c70abcf57811c55a1547516cf57d72f5b2ccc24`. Read-only inspection verifies the package hash `98fc642984f02f3a339b32113192a365659b8ab92fe88b2592bf1732f640e0d9` and extracted binary hash `64e962dd55e823f990bf7b4f21f6d4999b5584af691309e72e2f997cb17777fe`; extracted bytes match the supplied artifact. ELF64 little-endian AArch64 ET_EXEC has no PT_INTERP/PT_DYNAMIC. The package is **DEVELOPMENT / NOT ACCEPTED**, a synthetic internal test fixture, not distributable. The archive remains under `out/s04-ci-lintian/`; [machine metadata](evidence/s04/ci-39311632646.json), the [original package report](evidence/s04/ci-39311632646-package.json) and [original Lintian log](evidence/s04/ci-39311632646-lintian.log) are retained separately from historical records.

Locked Lintian 2.122.0 returns 2 for `statically-linked-binary` and `unstripped-binary-or-object`, both on `usr/bin/joy-pi-board`. Its other findings are warnings: `copyright-without-copyright-notice` and `no-manual-page`. The two errors conflict with the approved autonomous CGO=0 runtime and byte-preserving packaging of the frozen executable. The correction uses the documented [Lintian override mechanism](https://lintian.debian.org/manual/index.html), restricted to those two tags and that exact package/path, with reasons next to each entry. It neither strips the supplied binary nor changes compiler flags/tool locks. All other Lintian errors remain blocking; `--show-overrides` keeps overridden findings and reasons visible. Changing linkage or selected symbol policy requires reviewing these exceptions.

A manual page now describes the actual parser, defaults, diagnostic flag, environment precedence, packaged paths and exit behavior. Packaging installs it and verifies its decompressed extracted bytes, along with the exact shipped override file. No copyright identity or license is fabricated: the synthetic fixture warning remains visible and missing established metadata still blocks distributable packaging. Lintian output now reaches the CI console on success and failure while retaining the log and package evidence before a nonzero exit fails validation.

Local external-tool checks and four controlled Python subprocess logging tests pass, including validation exit 2 with original output retained and status returned for evidence capture, plus exit 0 with warnings retained. Changed local Markdown links/anchors, retained evidence hashes and `git diff --check` pass. Actual corrected-package Lintian, manual rendering and archive comparison remain **NOT RUN locally**: Docker's Linux engine pipe is absent; no engine/service was started. The next exact-commit CI must build a new artifact identity and run the unchanged error gate and comparison. No old frozen bytes or historical evidence were replaced. No package/maintainer script was executed, no Pi was contacted, and no commit/push/rerun occurred. Systemd lifecycle, Pi and deferred S03 operator acceptance remain NOT RUN; S04 remains open and release NO-GO.

## 2026-10-20 — Verified green CI and VM preparation

The user's “ok go”, resumed on 2026-10-20, authorizes archiving green CI and preparing the isolated VM campaign. Initial local source is clean at `5bc101e22a3093f0da3c0e6a7c50fee000dfa5da`. [Run 39313804438](https://github.com/Adrien-hue/joy-pi-board/actions/runs/39313804438), attempt 1, completed successfully on 2026-10-19: foundation 118269942467, race 118269942527 and browser 118269942304 all PASS. Source checks, bounded fuzz, native/cross builds, standalone/missing-dist checks, measurement tooling, Debian archive construction/Lintian and same-input archive comparison passed in their recorded CI scope. The browser report contains 27 expected passes, zero skips, unexpected or flaky results and no global errors. These are supported Linux automated engines, not real Safari/iOS or operator evidence.

Both GitHub ZIPs were downloaded before expiry and their hashes independently matched API metadata. Full logs, original manifests and package inventories, generated scripts and raw browser reports are retained in [the green CI evidence directory](evidence/s04/ci-39313804438/README.md); full ZIPs and selected executable/package remain in `out/s04-ci-green/`. The latter is Git-ignored and requires a separate backup for retention beyond this workstation. No publication occurred. This follow-up changes documentation/evidence only, without rebuilding or mutating the recorded artifacts.

| Identity / check | Verified result |
|---|---|
| Source and artifact status | `5bc101e22a3093f0da3c0e6a7c50fee000dfa5da`, clean at build, explicitly **DEVELOPMENT / NOT ACCEPTED**, ID `ci-39313804438`; no frozen candidate |
| ARM64 artifact, both extracted executables | SHA-256 `16bb1670c81638329b01ad7c83cc219c1957bb5ee1765b4a1091680f8dbbcb03`; exact equality, static ELF64 LE AArch64 with no PT_INTERP/PT_DYNAMIC; not executed |
| Two Debian archives | Same version `0.1.0-1+dev.ci-39313804438`, byte-identical SHA-256 `915f35240b9f9e6194bfbcee71fa2e4dcd98d475a606b8f635918ebe3fbb55dd`; reproducibility limited to these two builds with identical supplied inputs |
| Packaging environment | Actual image configuration ID `sha256:d2fddecaa22e7ab0e534a681abc892c37ec92eba95dcd4d80c36c1cdbc0f1508`, no registry RepoDigest; full package inventory hash `e0c2df777fa04ddd629697ad4dccf54d72d17e7b3ef82955aff18b7612612560` |
| Lintian | Both exit 0; two disclosed path-specific overrides; copyright warning remains visible; synthetic fixture is not distributable |
| Host and execution | Ubuntu 24.04.5 LTS / x86_64, runner image 20261012.334.1; Debian amd64 build container; no real systemd lifecycle test in container, no ARM64 runtime or Pi claim |
| Native encoded resources | Binary `ddfa1254b3c8dd478e1c4cd9cf3ccd8288be87a666b527104e8835d8e76a0052`; JS+CSS gzip 79,522 bytes; initial totals 80,768 complete / 82,001 partial / 145,545 maximum-markup bytes including first overview identity; raw base64 bodies and hashes verified; controlled host, not physical acceptance |

The generated-script review found **VM-F01**: no prerm, preinst stops only for upgrade, postrm removes/reloads bookkeeping without stopping the running process. This is a static implementation gap against the removal-stop requirement, not an observed VM failure. Keep A-06 open; reproduce active removal in the isolated guest and correct/revalidate the stop path before acceptance. Do not hide it by pre-stopping the VM-06 service. `deb-systemd-helper` mask bookkeeping also requires actual state tests rather than an assumption from the helper's name. The current lifecycle harness uses assertions and covers only a subset of the full planned matrix.

[The VM work package](s04-vm-validation.md) records VM-01–VM-10, environment/identity/isolation and recovery prerequisites, two-version requirements, commands, evidence oracles, the finding and the supplemental cases. All lifecycle cases are NOT RUN. No verified VM/access/snapshot/isolation details have been supplied; established Maintainer/contact and copyright/license remain missing. These gaps prevent execution or distribution respectively, not the completed independent archive/preparation work. Pi target/access/Health/service state/interruption window, restoration plan and required browser/assistive clients remain separate campaign prerequisites.

Local verification: ZIP/member/artifact/package/binary hashes, two-package equality, extracted manual equality, raw static/identity overview bytes, evidence index and links/anchors reviewed; `git diff --check` passes. The default lifecycle command was checked as inert (NOT RUN, no execution option). No application rebuild/test rerun, package installation, maintainer script/service execution, VM/Pi access, commit, push, tag or deployment occurred. S04-04/A-06, clean candidate freeze, C-01–C-06 and transferred S03 operator gates remain open; S04 is not closed and release remains NO-GO.
