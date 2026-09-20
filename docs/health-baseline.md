# Health baseline and evidence

Inspected on 2026-09-20, read-only, through commit-pinned GitHub file retrieval. Baseline: **`be7a0d824f62b94c842d8e5326110b1852c5a0bc`** in `Adrien-hue/joy-pi-health`. No newer revision was used; there is therefore no silent baseline change or latest-branch compatibility claim. Health was not modified, cloned into Board, imported or tested during this documentation task.

S00 follow-up: the user reports a read-only local comparison at `581b16f1a2b42db0484d62c08e0bf760cb557d84` confirming unchanged HTTP contract, snapshot encoder/tests and large-integer fixture. That reported comparison does not replace this baseline. S00 additionally retrieved the pinned [requirements](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/docs/requirements-v0.1.md) (blob `7435baae459374cbdf22156347d2756646d6bd2f`) and [numeric fixture](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/snapshot/testdata/large-integers.json) (blob `88e7b1bb4443635e9421e1cf62b55819be8e707d`). [ADR-001](adr-001-exact-json-integers.md) records the adopted convention; the fixture's copied bytes are verified without running Health or changing its repository.

## Inspected sources

Every link below resolves against the baseline commit, not a moving branch. Blob IDs identify the retrieved file content.

| Source | Blob SHA | Evidence |
|---|---|---|
| [docs/http-api-v0.1.md](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/docs/http-api-v0.1.md) | `ad6892ab003fe92d90b425ce045cbd82fba433b3` | Public transport/schema, nullability, exact JSON integers, size bounds and request-level failures |
| [internal/snapshot/snapshot.go](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/snapshot/snapshot.go) | `452335f646cd0e6ea9c18da6a7802d89c5827d25` | Wire structs, encoding, ranges, group coherence, expected issue paths/order |
| [internal/snapshot/snapshot_test.go](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/snapshot/snapshot_test.go) | `22059a43343031cb72b1248e9268ef8fc25608d9` | Null/order encoding, exact integers, invalid snapshots, no useful observations, maximum size |
| [internal/observe/coordinator.go](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/observe/coordinator.go) | `7c90bd98cacbfb31cc6b2ddded1d093ef03ac23a` | Expected failure mapping, fixed issue messages, collection order, state validation, hostname requirement |
| [internal/observe/coordinator_test.go](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/observe/coordinator_test.go) | `47953ea6dcdc979b8ff50ca597e973f908d636fa` | Frozen reason/order mapping, interface-local failure, shared/fresh collection cycles |
| [internal/httpapi/handler_test.go](https://github.com/Adrien-hue/joy-pi-health/blob/be7a0d824f62b94c842d8e5326110b1852c5a0bc/internal/httpapi/handler_test.go) | `a044e650a37d937108739b4277e0776bc7850143` | Provider HTTP outcomes, method/routing behavior, no-store, deadlines and HEAD semantics |

## Verified findings

| ID | Verified behavior and supporting evidence |
|---|---|
| H-01 | Health returns complete/partial snapshots with HTTP 200, no useful metrics with 503 `temporarily_unavailable`, internal failure with 500 `internal_error`. `TestProviderOutcomeMapping` in handler tests covers those outcomes. `Coordinator.finalize` rejects an unavailable hostname. |
| H-02 | `load`, `memory`, `root_filesystem` are non-null objects with either all three values or three null leaves. `buildDocument`, `validateAllOrNone`, `validateMemory` establish this; `TestEncodePartialSnapshotPreservesIssueOrderAndNulls` explicitly checks the memory object's null leaf and issue order. `network` alone can be a null group. |
| H-03 | `TestCoordinatorMapsExpectedReasonsInFrozenOrder` asserts the three codes, actual generic messages and fixed paths. `TestCoordinatorMapsInterfaceLocalFailure` asserts local state/RX issues without replacing the whole network group. |
| H-04 | Four firmware booleans are all present or all null; unavailable flags require four issues sharing code/message. Temperature is independent. `buildDocument` and `validateFirmwareIssues` enforce this. `TestEncodeRejectsInvalidSnapshots` includes partial firmware and wrong issue order. |
| H-05 | `TestEncodeCompleteSnapshotPreservesExactIntegers` and `TestLargeIntegerFixtureUsesExactJSONNumbers` preserve `9007199254740993` and `18446744073709551615` as number tokens. All integer fields in wire structs use nullable `uint64`. |
| H-06 | Coordinator network validation allows `up`, `down`, `unknown`; snapshot validation/public wire contract permit a non-empty UTF-8 state string. Both validate ordered unique names; snapshot encoding caps interfaces at 64. Board's degree of closure is proposed separately, not inferred as a closed public enum. |
| H-07 | Snapshot validation demands exactly the expected issue paths in order, allowed codes and non-empty UTF-8 messages; it does not compare all messages against the coordinator's literals. Firmware issues must agree with each other. Empty issues means observations available, not machine healthy. True firmware flags do not create issues. |
| H-08 | CPU percent range is 0–100; load is finite and nonnegative; temperature is finite without an invented health range; positive logical CPU count; byte groups require available ≤ total and used = total − available. All-null metrics without useful network observations cannot yield a successful snapshot. Encoding counts each interface name as useful even if that interface's nullable fields are unavailable; an empty interface array alone supplies no useful leaf. |
| H-09 | Successful encoded snapshot limit is 64 KiB. Public Health default is `127.0.0.1:8080`, request timeout 1 s, no compression, no request query/body. Handler tests confirm no-store and provider errors distinct from valid partial data. These are upstream facts, not automatically Board routing decisions. |

## Discrepancies and interpretation

The public API's wording “group ... null” is ambiguous for load, memory and root filesystem. The actual wire structs always encode these as objects; unavailable groups have **three null leaves and one group issue**. Board consumes that tested representation. Literal `"memory": null`, `"load": null` or `"root_filesystem": null` is invalid for this baseline. `"network": null` is valid with its corresponding issue. No correction was made to Health documentation.

The public issue example uses `Firmware status is not accessible to the service account.` The coordinator actually emits the generic `The metric is not accessible to the service account.` with the same code. Board examples use the actual emitted messages. Board's behavior depends on code/path, not a literal message comparison; it preserves and safely renders valid message text.

Public network `state` is a string; the observed coordinator emits three values. [P-03](decisions.md) proposes accepting other non-empty strings, displaying their text without inventing semantics. The observed implementation does not silently narrow the public type.

The [integration contract](health-integration-v1.0.md) translates this evidence into a Board consumer design. The [validation plan](validation-v0.1.0.md) references these tests as evidence to reproduce with Board-owned fixtures; none of their names is a Board test result.
