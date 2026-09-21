# Shared contract fixtures

This directory is the **single authoritative input and expectation source** for Go and TypeScript. Tests read the same files; there are no separately editable consumer copies. Data is never included in the production frontend.

## Pinned Health provenance

`numeric/health-large-integers.json` remains an exact byte copy of Health's fixture at `be7a0d824f62b94c842d8e5326110b1852c5a0bc`, Git blob `88e7b1bb4443635e9421e1cf62b55819be8e707d`. [Upstream source](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/snapshot/testdata/large-integers.json). The S00 integrity test continues checking its exact blob, including newlines. Health itself is not modified or tested.

`numeric/manifest.json` lists original token fragments and approved lexical expectations. Strings in this **metadata** preserve spelling; fixture files hold actual number/null/etc. bytes. `1.0` and exponent forms are now explicitly invalid at uint64 positions. They remain valid where a floating field permits them. Existing numeric files are unchanged.

## S01 authoritative corpus and oracles

[s01/cases.json](s01/cases.json) holds **264 cases**, of which **76 are valid**. Each has a stable ID, exact `input_base64` bytes, expected `classification` (`valid`, `invalid_response`, `unsupported_schema_version`) and, when valid, an expected known-field projection. Base64 is fixture storage only, allowing malformed UTF-8 to survive unchanged; it is never the API representation. Expected uint64 values use `u64:<decimal>` strings **in test metadata only**. Runtime/transport assertions separately require uint64/bigint and unquoted JSON numbers. Floating expected values keep numeric semantics; oracle normalization treats ±0 as numerically equal, while explicit tests/raw-token checks cover preservation of negative zero.

Seeds are Board's complete/partial documentary snapshots at `5f0d887bfc43e35711bc5a52f08164cf5197f3d5`, covering the two pinned Health integer extremes. [generate.py](s01/generate.py) is a developer-only authoring helper using Python standard-library arbitrary-precision integers and raw byte operations; it never routes them through JavaScript Number. The committed corpus is the build/test input. Python is not a runtime, build or CI dependency. `python testdata/s01/generate.py --check` checks regeneration bytes against the corpus; intentional oracle changes require review. Diagnostic results are not generated from either validator.

Families: complete and partial groups; three issue codes; local and whole-network failures; four coherent firmware flags and independent temperature; true/false/zero; every known required field removed and null-position checks; missing/null issue members; wrong issue codes/paths/count/order; firmware issue equality; useful interface-name versus empty-array/hostname/unknown-field insufficiency; interface ordering/uniqueness/64 maximum, including scalar versus UTF-16 sort order; all numerical boundaries and lexical errors; finite/domain floating rules; exact timestamp/calendar checks; extra members including prototype names; decoded duplicate keys (identical values too); malformed/trailing JSON; literal/escaped Unicode, invalid UTF-8 and byte-count/depth boundaries.

Go `TestSharedCorpus` and TS `snapshot.test.ts` consume the same expectations. Go also tests independent validated buffers/projections, object transport and untouched unknown fields/tokens. `FuzzDecode` seeds from the corpus. The test-only Go exporter emits actual validated bytes, marshaled object envelopes, classifications and known-field projections under ignored `out/s01-interop`; the TS consumer reads these real outputs. An additional **512 mutations**, Go seed **20260920**, compare both implementations without treating ordinary JSON.parse numbers as an exact-integer oracle. Build/fuzz outcomes live in [S01 evidence](../docs/s01-evidence.md), not in this fixture specification.

## Boundaries for following sprints

The test envelope is not a production Board API. Real HTTP limits/outcomes, controllable Health server, clock/cache semantics and UI display/scheduling belong to S02/S03. No generic fake-server/clock framework is created in S01. The full future path remains Health HTTP → Board overview → browser → dashboard. S01 proves its in-memory validation, object serialization, lossless parsing and decimal-primitive segments only. See [ADR-001](../docs/adr-001-exact-json-integers.md) and [validation](../docs/validation-v0.1.0.md).
