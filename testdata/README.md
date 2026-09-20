# Shared future contract fixtures

This directory is the **single authoritative numeric fixture source** for both future Go and TypeScript tests. Read the same files relative to the repository/test location; do not maintain separately editable copies under Go and web. Fixtures are never bundled into the production frontend.

`numeric/health-large-integers.json` is an exact copy of the read-only Health baseline fixture at `be7a0d824f62b94c842d8e5326110b1852c5a0bc`, Git blob `88e7b1bb4443635e9421e1cf62b55819be8e707d`. Provenance: [upstream source](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/snapshot/testdata/large-integers.json). The integrity test verifies the Git blob hash, including exact bytes/newlines; it does not run Health or a Board validator.

`numeric/manifest.json` lists raw-token fragments and future expected domains. Tokens are deliberately strings **in metadata only** so metadata can be parsed safely with ordinary JSON.parse; the actual fixture files contain original JSON value bytes, including unquoted large integers. `leading-zero.json` is intentionally invalid JSON. `decimal-integer`/`exponent` mark the lexical-acceptance policy still to ratify in S01. `float-written-integer` must be interpreted by its future field schema as a floating measure, not classified by token shape. Null is valid only in a nullable position; false is not an unsigned integer. No expectation here is a reported validator PASS.

S00's Node built-in tests check only byte integrity, manifest identity and unique IDs. They never parse numeric data through Number and are not a business parser. Go S01 can read the same files with RawMessage/UseNumber and checked uint64 conversion; TypeScript S01 reads text/bytes through its selected lossless parser. There is no synchronization script because there are no consumer copies.

## Snapshot and failure taxonomy for S01–S03

The existing [complete](../docs/examples/overview-current.json) and [partial](../docs/examples/overview-partial.json) documentary Board envelopes are the seed source, not duplicated here. In S01 extract `health.snapshot` as raw JSON without a floating conversion; move to an authoritative contract corpus and update doc references together if independent fixtures become necessary.

Future corpus: complete/partial, each issue code/group, network null and local state/RX/TX failures, coherent firmware and independent temperature, true/false/zero; missing/null fields; wrong issue order/path/code; unknown schema/member policy; exact numeric limits; all-unavailable; malformed UTF-8/JSON and oversize data. Future S02 transport fixtures use `httptest.Server` handlers with request counters/channels and per-test closures for wall/elapsed time. No generic clock/fake-server framework or cache is warranted before those consuming components exist.

The end-to-end oracle is Health snapshot bytes → Board object envelope bytes → browser exact parse → exact display. Only fixture integrity exists in S00; B-01–B-11 functional validation remains NOT RUN. See [ADR-001](../docs/adr-001-exact-json-integers.md) and [validation](../docs/validation-v0.1.0.md).
