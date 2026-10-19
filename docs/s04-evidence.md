# Sprint 04 software evidence

Status: **SOFTWARE IMPLEMENTED; VALIDATION INCOMPLETE; NOT CLOSED; release NO-GO**. Date: 2026-10-18. The user's current instruction approves S04-R01–R04 and S04-01–S04-05 software and isolated guest tests, replacing documentary-only restrictions. No primary-host/Pi package installation, account/service creation or operation, Pi connection, network change, commit, push, tag, release or deployment was performed.

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

## Final repository checks

437 local Markdown links/anchors checked with no missing target; requirement/task/gate references reviewed together. `git diff --check` passes. All 122 selected development input hashes still match the working files. Existing application toolchains/lockfiles, fixture bytes, S00/S01/S02 evidence and tracked historical evidence artifacts are unchanged. Packaging text uses LF endings. The diff is confined to authorized S04 runtime seams, tests, external tools, packaging/CI and English documentation/evidence; specification changes only the authorization row, without budget changes. No commit was created.
