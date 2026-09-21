# Sprint 01 execution evidence

Local work: 2026-09-20–21 (Europe/Paris). Repository: `Adrien-hue/joy-pi-board`, initially clean at `5f0d887bfc43e35711bc5a52f08164cf5197f3d5` (`docs: close Sprint 00 with verified CI evidence`). The S01 changes are an **uncommitted working tree**, identified by the [software/test input fingerprints](evidence/s01-inputs.json). They have not been pushed or tested in hosted CI. Health remains read-only at baseline `be7a0d824f62b94c842d8e5326110b1852c5a0bc`; its snapshot source was re-read, not imported, modified or executed.

**Status: implemented and verified locally; CI closure pending.** The detailed task boundary is [Sprint 01](sprint-01.md). No S02–S04 behavior is implemented. A successful S01 hosted run is still needed for full sprint closure; S00 run 35527458889 is historical evidence for another commit and is not reused here.

## Inputs, parser choice and policy

Exact executing tools remained **Go 1.27.1, Node 24.21.0, npm 11.19.0**. No general upgrade occurred. Go still has no external module dependencies. The only added production npm dependency is **lossless-json 4.3.1**, MIT, no production transitives. The lockfile SHA-256 is `50d760f57f4ec8de8eca4efbc6b4b75051214386dc1f2ac2a49573a38914200e`. Its registry metadata, versioned source, declarations and installed license were inspected; [the evaluation](sprint-01.md#parser-decision-and-primary-source-evaluation) explains the native browser-floor incompatibility, raw-token callback and required defensive guards.

P-03's S01 policies are implemented: exact-case presence, nullable groups and field domains, lossless uint64/bigint, unknown-member tolerance without semantics, decoded duplicate rejection at every object, open network states, exact issue associations/order/firmware equality, real useful-observation rule, 65,536-byte and 32-container inclusive limits, complete JSON/UTF-8/Unicode checking before version dispatch. `unsupported_schema_version` and `invalid_response` remain distinct; diagnostics never echo the payload. No literal matching of generic issue messages is used. The numeric convention remains inherited from Health; compatibility and library selection remain Board decisions.

Go's validated value owns copied original bytes and returns independent data on every access. Raw payloads, unknown members and numeric tokens survive object serialization. TypeScript's known projection is established only after runtime validation; every uint64 is bigint, every floating metric is number. The exact decimal primitive is implemented; P-06 visual units/rounding remain deferred.

## Executed local checks

`PASS` below is scoped to S01/foundation subsets, not full product A/B/C. Commands ran on Windows amd64 unless the native Linux container is named. Frontend was built before Go checks needing embed. The task runner enforces exact tools, `GOTOOLCHAIN=local`, `GOENV=off`, `GOWORK=off` and native/cross CGO=0; race uses CGO=1.

| ID / subset | Command or actual procedure | Result |
|---|---|---|
| S01-L01 / A-01–A-03 inputs | `node scripts/tasks.mjs tools`, npm preinstall; exact executing versions checked | PASS; same S00 toolchain pins, no implicit Go switching or npm major upgrade. |
| S01-L02 / A-02 install | `npm --prefix web ci --offline=false --cache .cache/npm-s01-cold --no-audit --no-fund` | PASS; 178 installed packages from a new task-specific npm cache and the committed-format lockfile. No pre-existing node_modules prerequisite. |
| S01-L03 / A-01–A-02 and B-01–B-04 subsets | `npm --prefix web run format`, `npm --prefix web run check` | PASS; installed-tree consistency, Prettier, strict TS, ESLint, **271 frontend tests**, two original fixture-integrity tests, Vite build, gofmt, module tidy/verify, vet, first-party Go tests and actual interop. The final Go presence-diagnostic adjustment and test-only bounded-roundtrip regression were additionally checked with `go test -count=1 ./internal/healthschema` and `go vet ./internal/healthschema`. |
| S01-L04 / B-01–B-04 shared contract | Go `TestSharedCorpus`, `TestRequiredFieldPresence`, `TestImmutableAndObjectTransport`, `TestPreservesUnknownAndOriginalTokens`, `TestZeroValueAndDiagnostics`, `TestRoundtripAtLimitWithMarkup`; TS `snapshot.test.ts` | PASS; **264 common cases, 76 valid**. Exact known projections/classifications agree with the common oracle. Missing members are diagnosed separately from explicit null, including schema_version. |
| S01-L05 / B-03 interoperability and B-04 differential | `npm --prefix web run contract`: Go `TestExportInterop` → dedicated bundled TS consumer | PASS; **76 actual Go-serialized snapshots** parsed in TypeScript and compared to the common expectations; test-only envelope contains an object with unchanged number tokens/unknown members. **512 seeded mutations**, seed 20260920, agree on classification and successful projection (one accepted mutation). [Report](evidence/s01-interop.json). |
| S01-L06 / B-04 adversarial | `npm --prefix web run fuzz` → `go test -run=^$ -fuzz=^FuzzDecode$ -fuzztime=10s -parallel=2 ./internal/healthschema` | PASS, no failing input/panic. Expanded corpus with a fresh `.cache/go-fuzz-s01-final` cache: 257 distinct seeds, 384,654 executions. After the final presence-diagnostic change: 139,079 executions, 280 starting seeds including the prior 23 interesting cached inputs. After the final roundtrip test-harness correction: 40,828 executions, 287 starting seeds, PASS. Fuzz scheduling/counts are not deterministic; committed seeds and the separate differential generator are reproducible. |
| S01-L07 / A-01 native race | Native Linux amd64 container command below | PASS on final inputs after restoring Docker: cmd 1.017 s, healthschema 1.774 s, httpui 1.020 s; no race report. Web has no Go test files, not a placeholder success test. |
| S01-L08 / A-03 build | `npm --prefix web run build`; `npm --prefix web run cross` | PASS; Windows amd64 and Linux ARM64 binaries, frontend first. Source changes after these builds affected only contract diagnostics/tests/docs, not the shell dependency graph; S01 contract code is not wired into the application. |
| S01-L09 / A-04 shell regression | `npm --prefix web run smoke`; `npm --prefix web run missing-dist`; `npm --prefix web run sizes` | PASS; independent directory/empty PATH serves exact compiled assets, zero Health trap connections through startup/requests/6 s idle, absent overview still 404, missing dist fails clearly. [Standalone report](evidence/s01-standalone-smoke.json), [shell report](evidence/s01-shell-sizes.json). |
| S01-L10 / fixture authorship | `python testdata/s01/generate.py --check`; original Git-blob integrity test | PASS; corpus regeneration matches committed-format bytes; pinned Health numeric fixture unchanged. Generator uses exact Python integers and byte strings, not JavaScript Number. Python is authoring-only, not needed by CI/tests/build. |
| S01-L11 / documentary and source integrity | Link/anchor and status review, input-fingerprint verification, `git diff --check` | PASS; 19 Markdown documents / 172 internal links and anchors; all 59 software/test fingerprints match. Existing S00 evidence JSON, shell sources, Go module, toolchain pins and pinned Health fixture checked unchanged. |

Native race command (frontend already built):

```text
docker run --rm --network none --mount type=bind,source=C:\src\joy-pi-board,target=/src,readonly -w /src -e GOTOOLCHAIN=local -e CGO_ENABLED=1 golang@sha256:69a7b9788769bec032d238959b61854e9ae87f57be9029ec04e9885fabf99195 go test -race -count=1 ./cmd/... ./internal/... ./web
```

The image is the same pinned Linux/amd64 Go 1.27.1 / Debian bookworm compiler environment used in S00. It runs without network and with the repository read-only; it is not ARM64 execution or a Raspberry Pi measurement. An earlier S01 run passed cmd, httpui and healthschema under race. On resumption the Docker engine was stopped; the first final-state attempt could not connect to its named pipe. Docker Desktop was restarted and the final retry passed on the recorded inputs; the failed environment attempt is retained separately.

## Artifact provenance and measured parsing cost

[S01 input fingerprints](evidence/s01-inputs.json) cover 59 software/test input files, excluding documentary evidence and generated outputs. Manifest SHA-256: `1b86b786db5979859c86b12488a921c16e62557ffbf06272683a4a212923c7e0`. Common corpus SHA-256: `50a329c04602b86af6339875afd8b405ee79253d3ed62d044ac2708f398c6161`. These identify a local working tree, not a hosted commit or release candidate.

| Local artifact / measurement | Actual result |
|---|---|
| `out/s01-parser/contract.mjs` | **17,400 bytes minified; 5,550 bytes gzip-9**. SHA-256 `c5d9f605b355731b04cc6231741d1b03ba3078967febd69625ebc84363f5cf70`. |
| Module scope | Dedicated Vite entry imports the actual lossless parser, defensive checks, schema validator and decimal primitive, retaining public exports; no React or fixtures. Node 24.21.0, zlib `1.3.2.1-motley-8002e91`. [Measurement report](evidence/s01-parser-sizes.json). |
| Minimal shell | JS+CSS gzip-9 **68,495 bytes**, all initial resources **68,790 bytes**. Identity bytes 221,644. Parser is unused by the shell and tree-shaken; that absence is not its measured future cost. |
| Windows shell binary | SHA-256 `c75c48cc1558fa1ecdafba017fd2f20bcf3fd502d48efb2be9c2cc4ba626069d` |
| Linux ARM64 shell binary | SHA-256 `d44a4765b719a7669d709d0faaed7b5a5d4d907f57dc8ca2c50fbd0cde1e6c0f` |

Shell hashes were produced again by the local S01 build commands and match the historical local S00 hashes, as expected for unchanged reachable application code/assets. This does not mean S00 CI tested S01, nor establish universal bit-for-bit reproducibility. Old S00 evidence JSON and its provenance are preserved. Compression is measured offline; no HTTP compression or final dashboard size/CPU/RSS/latency acceptance is claimed. A future integrated bundle must be measured again rather than adding standalone compressed sizes mechanically.

## Corrected diagnostics and execution limitations

- Initial registry access encountered sandbox offline/cache restrictions; the authorized online registry query/install succeeded. An npm upgrade notice was ignored; pins remain unchanged.
- Initial Vitest startup was blocked by sandbox `spawn EPERM`; permitted execution outside that sandbox succeeded. The first test fixture URL used one excess parent segment; correcting the test path made the suite run.
- Interop initially compared Go's signed floating zero to normalized JSON test metadata as different; test normalization was corrected, with an explicit TS sign-preservation test added. Original numeric tokens remain independently asserted. No contract precision was weakened.
- The test-only envelope initially reused the snapshot's depth limit, rejecting a valid depth-32 snapshot after adding its wrapper. Envelope tests now parse the Go-generated wrapper separately; the original snapshot still goes through the full bounded TS validator. Production envelope limits remain S02.
- Boundary review corrected the fuzz roundtrip invariant: Go's default JSON HTML escaping can expand a valid 65,536-byte snapshot beyond the input limit. The test-only roundtrip encoder now disables optional HTML escaping and removes its framing newline; a maximum-size markup regression covers this distinction. Default object serialization and original numeric-token preservation remain tested. Production encoder choices and full-envelope limits remain S02.
- Review added an explicit Go schema_version presence lookup and a required-field diagnostic test. Invalid classification was already correct; absence is now visibly distinct from null at this field too.

Actual JavaScript execution: Node 24.21.0 Vitest and the browser-targeted bundled module. **Chrome/Edge 111, Firefox 115, Safari/iOS 16.4 execution: NOT RUN.** No browser-floor increase occurred. Compatibility rests on reviewed APIs/source and successful targeted compilation; this is not full browser acceptance. Interactive visual checks, mobile acceptance and all Raspberry Pi physical checks remain **NOT RUN**.

## CI and exit status

The workflow now names its scope `Board foundation and contracts`, retains immutable action SHAs, exact toolchain pins, explicit Ubuntu 24.04 and actual image recording. `check` includes the shared corpus and Go→TS tests; a separate bounded fuzz step was added. The native race job includes healthschema automatically. The internal artifact now includes contract reports, the dedicated module and test-only Go serialization outputs. There is no deployment/release job.

**Hosted S01 CI: NOT RUN.** Workflow edits and local success do not constitute a hosted pass. No new commit, push, dispatch, release or deployment was performed. S01 is implemented and locally verified, with CI closure pending. Browser-family execution is an explicit outstanding validation limit; actual dashboard/browser display and HTTP/cache assertions await S02/S03. Full Class A/B/C product acceptance is not claimed; all B-05–B-11 and C-01–C-06 remain NOT RUN.
