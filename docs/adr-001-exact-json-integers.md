# ADR-001 — Exact JSON integers across Joy Pi Home

Status: **ADOPTED by the S00 user instruction, 2026-09-20**. Scope of this change: Board documentation and future implementation; no other repository is changed.

## Decision and origin

Joy Pi Home adopts Health's existing wire convention: exact integers are unquoted decimal JSON numbers regardless of magnitude. This is a contract and result convention, not a requirement to use the same library in every language.

The Board baseline remains Health commit `be7a0d824f62b94c842d8e5326110b1852c5a0bc`. Sources: [HTTP numeric/text rules](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/docs/http-api-v0.1.md), [requirements](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/docs/requirements-v0.1.md), [snapshot encoder](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/snapshot/snapshot.go), [encoder tests](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/snapshot/snapshot_test.go) and [large-integer fixture](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/snapshot/testdata/large-integers.json).

The user reports a subsequent read-only local comparison at `581b16f1a2b42db0484d62c08e0bf760cb557d84`: HTTP contract, encoder, tests and large-integer fixture were unchanged. This reported comparison is not a replacement baseline or a claim of a comparison executed during S00. See [Health evidence](health-baseline.md).

Health uses nullable `*uint64` fields, standard `encoding/json`, and exact-digit assertions for `9007199254740993` and `18446744073709551615`. Its decoding tests use `Decoder.UseNumber()` / `json.Number` to avoid generic float64 conversion. Health did not choose a big-number dependency or an owned general-purpose JSON parser.

## Common rules

- Each schema defines signedness, bounds and nullability for each integer field.
- Exact integers remain unquoted decimal JSON numbers; converting a numeric field to a JSON string changes the contract.
- No producer, intermediary or consumer may route exact integers through a floating representation that can round them.
- Decimal measures retain their floating-point semantics; token spelling alone does not decide a field's domain.
- Null is distinct from zero and false. Required-field absence remains a separate validation concern.
- Presentation may round a summary but must retain the original exact value.

## Board application, implemented in S01–S03

S01 Go validation uses original `json.Number` tokens followed by checked uint64 conversion and range validation. `UseNumber()` retains a token; it does **not** validate signedness, integer syntax or range. Retaining validated `json.RawMessage` is compatible: embed it as an object, never as quoted JSON or base64. S00 had no validator; S01 implements it without provider transport.

TypeScript snapshot uint64 fields use `bigint` at every magnitude, bounded to `0n..18446744073709551615n`. Floating fields remain `number`, including a CPU percentage or temperature written as an integer-looking token. Conversion is driven by the schema. `response.json()` followed by bigint conversion of an already-rounded Number is forbidden.

S01 selected lossless-json 4.3.1 after evaluating the native source-token facility and unchanged browser floors; see [the evaluation](sprint-01.md#parser-decision-and-primary-source-evaluation). Evaluate compatibility, license, maintenance, invalid inputs, duplicate keys, limits and actual production bundle contribution. Do not automatically select a package. No unused parsing dependency or general-purpose homemade parser is introduced in S00. See [foundation choices](s00-foundation.md) for browser targets.

Unknown-member policy is independent: **S01 approves tolerance** under the structural/duplicate limits in P-03, not as a numerical requirement or an inherited Health decision. Required fields, nullability, types, bounds and exactness remain mandatory whichever compatibility policy is selected.

## Verification and consequences

The authoritative [shared fixtures](../testdata/README.md) preserve numeric tokens without a JavaScript numeric parse/stringify round trip. S00 checks fixture integrity, not business validation. Planned B-03 must cover Health snapshot → Board object envelope → browser parse → exact display, including uint64 extremes. S01 now verifies Go validated bytes → test-only object envelope and TypeScript parsing → exact decimal primitives in memory. The later [S02 evidence](s02-evidence.md) records real simulated Health HTTP → Board HTTP → TypeScript parsing checks; target-browser execution and dashboard display remain NOT RUN. The original [S01 evidence](s01-evidence.md) retains its in-memory scope.

This ADR supersedes P-03's initial preference for an owned tokenizer. It does not ratify unrelated validation, HTTP, concurrency, freshness or release measurement proposals. References: [integration](health-integration-v1.0.md), [Board API](board-api-v0.1.0.md), [validation](validation-v0.1.0.md), [decision register](decisions.md).
