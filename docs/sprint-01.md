# Sprint 01 — Health contract foundation

Status: **COMPLETE, closed on 2026-09-21 with local verification and verified hosted CI**, recorded in [S01 evidence](s01-evidence.md#verified-hosted-ci). S00 remains complete. The current instruction authorizes documentary closure only and supersedes the temporary S01 implementation authorization. S02 is ready to be prepared, but NOT STARTED; S02–S04 implementation is not authorized.

## Scope and ownership

This sprint consumes Health schema 1.0 at pinned commit `be7a0d824f62b94c842d8e5326110b1852c5a0bc`, in memory only. [Health integration](health-integration-v1.0.md) owns the normative fields/invariants and approved consumer rules; [ADR-001](adr-001-exact-json-integers.md) owns numerical exactness; [P-03](decisions.md) separates Board compatibility choices from inherited Health facts.

Implemented packages: Go [internal/healthschema](../internal/healthschema/types.go), TypeScript [web/src/health](../web/src/health/snapshot.ts), authoritative [shared corpus](../testdata/s01/cases.json). The React shell, listener defaults and runtime request behavior are unchanged. No Health process, HTTP client, overview API, cache, freshness state, scheduling, dashboard, production lifecycle, packaging, release or deployment is added.

| Task | Output and dependency | Exit evidence |
|---|---|---|
| S01-01 | Adopt consumer P-03 rules and evaluate parsing against unchanged browser floors. Depends on S00 and H-01–H-09. | Version/license/source review, explicit unknown/duplicate/Unicode/limit semantics; parser selection below. |
| S01-02 | Board-owned Go types, standard-library token decoding, exact presence/type/range/issue checks and immutable original payload. Depends on S01-01. | Shared corpus, mutation isolation, object serialization, race and fuzz tests. |
| S01-03 | Token-preserving TypeScript parsing, runtime validation, bigint uint64 fields and exact decimal primitive. Depends on S01-01. | Same corpus/oracles, strict types/lint/tests and dedicated bundle measurement. Browser executions separately identified. |
| S01-04 | Shared positive/negative corpus and genuine Go → TypeScript handoff with a test-only envelope. Depends on S01-02/S01-03. | All reviewed cases agree with expectations; seeded differential mutations agree; original fixture bytes unchanged. |
| S01-05 | Integrate checks/fuzz in CI, repeat install/build/embed/race checks, update evidence and traceability. Depends on S01-04. | DONE: local checks plus S01-CI-01–S01-CI-05 PASS at `ccd1245a64fc372faa4ebc85eb000f78a01bd396`, run 35637699325 attempt 1. S00 evidence remains separate. |

## Parser decision and primary-source evaluation

Selected **lossless-json 4.3.1**, MIT, exact lockfile pin, no production transitive dependencies. The npm registry was queried during this task; publication metadata reports 2026-07-31, and the [versioned package manifest](https://github.com/josdejong/lossless-json/blob/v4.3.1/package.json) provides ESM, browser distribution and TypeScript declarations. No other toolchain/dependency version was upgraded. Installed code/license and the [versioned parser source](https://github.com/josdejong/lossless-json/blob/v4.3.1/src/parse.ts) were inspected. Upstream remains a maintenance dependency; its future releases are not automatically adopted.

Native `JSON.parse` source-context is unsuitable for the retained floors: [MDN's compatibility data](https://github.com/mdn/browser-compat-data/blob/main/javascript/builtins/JSON.json) reports Chrome/Edge 114, Firefox 135 and Safari 18.4. Board keeps Chrome/Edge 111, Firefox 115 and Safari/iOS 16.4. The selected library exposes each original numeric token to `parseNumber`; Board stores it in a private-branded token until field-specific conversion. Library ESM uses ordinary supported language features; Board's dedicated Vite build targets all four existing floors. TextDecoder fatal UTF-8, TextEncoder, bigint, Set, Object.hasOwn and private fields fit those floors. This is source/API and compilation evidence, **not execution in those browsers**.

The library alone is insufficient for Board's entire policy. It allows equal duplicate values, accepts individual escaped surrogate code units, and has no input-size/depth option. Board adds a bounded preflight: count containers outside quoted strings, validate quoted string Unicode (native JSON.parse on string tokens only), and compare decoded key names in per-object sets. This guard produces no JSON values, recognizes no numeric/keyword grammar, and does not replace the library parser. A raw numeric grammar check also guards its custom number callback. Input is never rewritten.

Unknown `__proto__` members can affect a library-created object's prototype; Board never trusts inherited fields, uses a private numeric brand, and constructs a fresh known-field projection. Such unknown members may disappear from the TypeScript projection, as allowed by the contract. Their original bytes remain intact in Go's validated transport. Tests cover prototype names, inherited-field substitution, identical/escaped duplicates, surrogate failures and overflow. No global prototype is modified. Upstream diagnostic strings are replaced by bounded Board errors because they may quote input.

The narrow guard plus the evaluated library satisfies the approved consumption rules. No general-purpose JSON parser is owned by Board, and no native source-token API or modern-only polyfill is assumed. See [S01 evidence](s01-evidence.md) for the actual minified/gzip cost of the imported consumer module; the unused parser is correctly absent from the shell bundle.

## Build, tests and handoff

Existing `npm --prefix web run check` now runs TS/component tests, fixture integrity, frontend build, Go checks, then `contract`. `contract` executes Go's test-only exporter and bundles the actual TS consumer before parsing Go output. `npm --prefix web run fuzz` runs `FuzzDecode` for 10 seconds with two workers and corpus seeds; Go's mutation schedule is not claimed deterministic. The separate differential mutation stream uses fixed seed 20260920. `race` retains a native compiler requirement. All existing native/cross/smoke/missing-dist commands remain valid.

`scripts/contract.mjs` creates a dedicated browser-targeted ESM module from `web/contract-entry.ts`, retaining its exported functions, then measures its bytes and offline gzip-9 bytes. It is not an application entry, not served by the shell, and not part of final dashboard size acceptance. The Go test envelope adds one container around the snapshot; it is parsed by the library in the test harness without imposing a premature S02 envelope limit. The actual snapshot is separately passed through Board's bounded TS validator.

S01 closes B-01 snapshot-only, B-02, B-03 exact in-memory parsing/decimal primitives and B-04 consumption-policy subsets with local and hosted CI evidence. HTTP statuses, provider failures, cache, full overview envelopes, browser scheduling and dashboard display remain NOT RUN. Future S02 can consume `healthschema.Decode`, its stable error classification and immutable JSON marshaler; no network or stateful seam is prebuilt here.
