# Integration with Health snapshot schema 1.0

Status: **normative consumed contract; S01 validator/parser implemented and locally verified, hosted CI pending**. Normative for Board's consumed field names, units, nulls and issue semantics. Source: [pinned Health evidence H-01–H-09](health-baseline.md), never Health Go package imports. Requirements JPB-014–016.

The common exact-number convention is adopted in [ADR-001](adr-001-exact-json-integers.md). S01 implements Go validation and TypeScript lossless parsing without a Health client or UI integration. [P-03](decisions.md) records the approved consumption rules separately from inherited Health facts and ordinary implementation choices. [S01 evidence](s01-evidence.md) owns actual results and limits.

## Transport and envelope

Consume `GET /v1/snapshot` at configurable startup URL, default `http://127.0.0.1:8080/v1/snapshot`. No query/body. Complete and valid partial results are HTTP 200. Health 503 `temporarily_unavailable` means no useful snapshot (including unavailable mandatory hostname); Health 500 `internal_error` invalidates the whole cycle. Both map to Board `upstream_error`, regardless of the upstream error body's text. The [Board API](board-api-v0.1.0.md) owns timeouts and failure mapping.

All fields in the table are **required**, including nullable fields. An absent key is invalid, not shorthand for null. Objects are non-null unless explicitly stated. Arrays are never null. JSON strings must be valid UTF-8; `observed_at` is valid RFC3339 UTC ending in `Z`, fractional seconds allowed. There is no complete/partial flag, hardware model, `status` field or issue severity.

| Field | Wire type | Validity / unit |
|---|---|---|
| `schema_version` | string | Exactly `"1.0"` supported |
| `observed_at` | string | Health timestamp, preserved |
| `host.hostname` | string | Non-empty |
| `cpu.utilization_percent` | number or null | Finite, 0–100, already percent |
| `cpu.logical_cpu_count` | integer or null | uint64, positive when present |
| `load.one_minute`, `five_minutes`, `fifteen_minutes` | number or null | Finite, nonnegative, dimensionless; all present or all null |
| `memory.total_bytes`, `available_bytes`, `used_bytes` | integer or null | uint64 bytes, coherent group |
| `root_filesystem.total_bytes`, `available_bytes`, `used_bytes` | integer or null | uint64 bytes, coherent group; available to unprivileged service |
| `uptime_seconds` | integer or null | uint64 seconds |
| `network` | object or null | Null only for whole-network unavailability |
| `network.interfaces` | array | Required when network is an object, possibly empty, ≤64 interfaces |
| `network.interfaces[].name` | string | Non-empty, unique, ascending by name |
| `network.interfaces[].state` | string or null | Non-empty when present; public type is not a closed enum |
| `network.interfaces[].rx_bytes`, `tx_bytes` | integer or null | uint64 cumulative bytes, never rates |
| `raspberry_pi.soc_temperature_celsius` | number or null | Finite Celsius, independent of firmware flags |
| `raspberry_pi.thermal_throttling_active`, `thermal_throttling_occurred_since_boot`, `undervoltage_active`, `undervoltage_occurred_since_boot` | boolean or null | All four present or all four null |
| `issues` | array of Issue | Required, never null; ordered exact associations below |

For each byte group, all three leaves are null or all three are present, `available_bytes ≤ total_bytes`, and `used_bytes = total_bytes - available_bytes`. A total of zero is structurally valid; percentage display is then unavailable. Do not add temperature/load health thresholds. Zero values and false flags are meaningful available observations.

**Verified H-02:** load/memory/root filesystem remain objects when unavailable, with three null leaves and one issue for the group. The ambiguous upstream phrase “group null” does not mean these objects disappear. `network` really becomes null on enumeration/topology failure. Local interface failures leave names and independently available fields/interfaces intact. At least one useful observation is required; use H-08's exact useful-observation rule when checking the pathological all-unavailable case.

## Exact issue contract

Each emitted Issue has exactly `path`, `code`, `message`, all required and non-null. `path` is the canonical JSON Pointer to the field/group in this snapshot. `code` is one of the three values below. `message` is non-empty valid UTF-8, stable, safe user-facing text, preserved by Board.

| Code | Message emitted by the baseline coordinator |
|---|---|
| `unsupported` | `The metric is unsupported on this platform.` |
| `permission_denied` | `The metric is not accessible to the service account.` |
| `temporarily_unavailable` | `The metric is temporarily unavailable.` |

Board uses **code and path** for logic; messages are text, never HTML or a lookup key. The validator accepts valid non-empty text rather than comparing with these literals. There is no severity, identifier or issue timestamp. `internal_error` is forbidden in `issues[].code`; it is a whole-request Health error.

Build the expected path list by walking unavailable fields in exactly this order; compare to issues one-for-one, including order (no extras, duplicates or omissions):

1. `/uptime_seconds`
2. `/cpu/logical_cpu_count`
3. `/load`
4. `/memory`
5. `/root_filesystem`
6. `/raspberry_pi/soc_temperature_celsius`
7. `/raspberry_pi/thermal_throttling_active`
8. `/raspberry_pi/thermal_throttling_occurred_since_boot`
9. `/raspberry_pi/undervoltage_active`
10. `/raspberry_pi/undervoltage_occurred_since_boot`
11. `/network` if network is null; otherwise, for each interface in array order, `/network/interfaces/{i}/state`, `/network/interfaces/{i}/rx_bytes`, `/network/interfaces/{i}/tx_bytes`, only for null leaves
12. `/cpu/utilization_percent`

The notation `{i}` denotes the canonical decimal array index, without leading zeros; it is not literal JSON. Interfaces are sorted by name, not issue order or availability. Index is local to this snapshot, never persistent identity. Four unavailable firmware flags create four issues sharing code and message. Temperature failure is independent. An empty issue list means all measures are available; it is not a good-health verdict. Active undervoltage/throttling can be true with no issue.

## Go types and required-field presence

The implemented [Go types](../internal/healthschema/types.go) belong to Board. `Decode([]byte) (Validated, error)` first walks standard-library JSON tokens with `UseNumber`, preserving each numeric token and checking every decoded object key. A map lookup returns both the value and a **presence bit**: an absent required key fails before nullability is interpreted. Known keys are matched exactly and case-sensitively. No case-insensitive struct decoding or plain pointer alone establishes input presence.

After those checks, `Snapshot` is a known-field projection: `Host`, `CPU`, `Load`, `ByteGroup`, `Network`, `NetworkInterface`, `RaspberryPi`, `Issue`. Nullable metrics are `*uint64`, `*float64`, `*bool` or `*string`; `Network` alone is a nullable object pointer. Required strings/objects/arrays have already been checked, and integer tokens have passed lexical checks and `strconv.ParseUint` bounds. Projection pointers mean explicit null only, never absent.

`Validated` privately owns a copy of the original input bytes. `Bytes()` and `MarshalJSON()` return fresh copies; `Snapshot()` re-decodes the validated bytes into independent typed data. Mutation of the input after Decode, returned buffers, or a returned projection cannot change the validated payload. The zero value is invalid and refuses serialization. Concurrent reads are safe; callers must not modify input during Decode itself.

Its `json.Marshaler` implementation permits a future envelope to include the snapshot as an **object**, preserving unquoted number tokens and unknown members. `encoding/json` may compact whitespace or change equivalent string escaping; the private original bytes remain unchanged. Optional HTML escaping can increase serialized byte size; S02 must account for encoder behavior when defining transport and full-envelope limits. The typed projection intentionally omits unknown fields and must not replace the transport payload. This implements no cache or production overview envelope.

## TypeScript types after successful runtime validation

All integer fields use `bigint`, even when small, to avoid a size-dependent union. These are in-memory types, **not** a change to JSON number encoding. Before validation the parse result is untrusted, with explicit member-presence checks (`Object.hasOwn` or equivalent). An omitted required key fails; it must not be replaced by null or cast into these types.

```ts
type UInt64 = bigint; // runtime range: 0n .. 18446744073709551615n
type Nullable<T> = T | null;
type IssueCode = "unsupported" | "permission_denied" | "temporarily_unavailable";
interface Issue { path: string; code: IssueCode; message: string }
interface ByteGroup {
  total_bytes: Nullable<UInt64>;
  available_bytes: Nullable<UInt64>;
  used_bytes: Nullable<UInt64>;
}
interface NetworkInterface {
  name: string;
  state: Nullable<string>;
  rx_bytes: Nullable<UInt64>;
  tx_bytes: Nullable<UInt64>;
}
interface HealthSnapshot {
  schema_version: "1.0";
  observed_at: string;
  host: { hostname: string };
  cpu: { utilization_percent: Nullable<number>; logical_cpu_count: Nullable<UInt64> };
  load: {
    one_minute: Nullable<number>;
    five_minutes: Nullable<number>;
    fifteen_minutes: Nullable<number>;
  };
  memory: ByteGroup;
  root_filesystem: ByteGroup;
  uptime_seconds: Nullable<UInt64>;
  network: { interfaces: NetworkInterface[] } | null;
  raspberry_pi: {
    soc_temperature_celsius: Nullable<number>;
    thermal_throttling_active: Nullable<boolean>;
    thermal_throttling_occurred_since_boot: Nullable<boolean>;
    undervoltage_active: Nullable<boolean>;
    undervoltage_occurred_since_boot: Nullable<boolean>;
  };
  issues: Issue[];
}
```

## Lossless browser parsing and presentation

Implemented parsing uses **lossless-json 4.3.1** and a narrow defensive preflight; [the S01 evaluation](sprint-01.md#parser-decision-and-primary-source-evaluation) records its primary sources, gaps, browser assumptions and bundle measurement. `parseSnapshot(Uint8Array | string)` validates byte size, Unicode, structure and schema before returning `HealthSnapshot`. Bytes use fatal UTF-8 decoding without BOM stripping; string input is checked for unpaired UTF-16 before measuring its UTF-8 byte size. No Fetch client exists yet. **`response.json()` followed by Number-to-bigint conversion is forbidden.**

At uint64 positions, accept only an unquoted token matching `0` or a nonzero digit followed by digits, then checked conversion in range 0..18446744073709551615. Reject negative tokens including `-0`, fractions, exponents, strings and booleans. At floating metric positions, convert the original numeric token to Number/float64 and enforce finitude and the field domain; integer-looking tokens are still floating metrics. Floating `-0` is preserved, and a subnormal underflow to zero follows IEEE-754 conversion in both languages. No temperature/load health threshold is added. Numbers in unknown fields are not interpreted as known measures.

The browser never converts a counter to Number. `decimalUInt64(bigint)` verifies bounds and returns exact base-10 text, including small values. It does not prescribe visual units or rounding. **Proposed P-06 presentation** (still deferred): binary-unit summaries and percentages can use bigint division/remainder; an exact value must remain available alongside rounded summaries. No such UI is implemented in S01.

## Validator policy and failure boundary

**Approved S01 Board consumption rules**, distinct from Health's producer guarantees:

- Limit the snapshot to **65,536 bytes inclusive before parsing**, including whitespace. Count encoded UTF-8 bytes, not JavaScript string length. The root object is container 1; each nested object/array adds one, including empty containers and unknown subtrees. Maximum simultaneous depth is **32**, inclusive. Strings containing brackets do not increase depth. Future HTTP bounded reads and the proposed full overview limit remain S02.
- Validate the entire JSON document, including unknown content, before version dispatch. Reject trailing documents/non-JSON whitespace, malformed syntax/UTF-8, BOM, unpaired escaped or literal surrogates, and duplicate decoded keys anywhere. Valid surrogate pairs and a genuinely encoded U+FFFD are accepted without normalization/replacement. Go's otherwise permissive Unicode replacement is explicitly guarded.
- For schema 1.0, **tolerate extra members** at any object level, including issues, and preserve them in the validated Go payload. The baseline producer emits only the documented issue members. Extras cannot replace a required known key, contribute useful measures, relax invariants, or authorize unknown issue codes/paths. TypeScript returns only known fields.
- Required keys are exact-case. Missing, explicit null and present values are distinct. Known objects/arrays are never null except network. Reject incoherent groups, invalid domains and mismatched issue count/path/order. Never sort, complete, repair or merge input.
- Network state accepts any non-empty Unicode string. `up`, `down`, `unknown` are observed producer values, not a closed wire enum. Message text is non-empty Unicode, not compared with generic literals; the four firmware issues must share code/message.

After global structural checks, a root object with a present string `schema_version` other than `1.0` yields **unsupported_schema_version**, before any version-specific interpretation (even if other fields are absent). Every malformed/global-limit case, non-object root, absent/non-string version, or schema 1.0 violation yields **invalid_response**. Diagnostics identify a bounded rule/known path without echoing values or whole payloads. No HTTP error envelope is implemented here.

`observed_at` is preserved and must be a nonzero, valid Gregorian RFC3339 UTC timestamp ending in uppercase Z; no offset, comma fraction or leap second. Health cannot encode Go's zero time. Fractional seconds are accepted; zero-time detection uses nanosecond interpretation consistently with Go. Hostname alone never makes a snapshot useful; neither does an empty interface array or an unknown field. **Each interface name counts as useful**, even with state/RX/TX all null. Otherwise any available metric suffices; no health verdict is implied.

## Documentary examples and shared fixtures

[Complete](examples/overview-current.json) includes exact RX `9007199254740993` and TX `18446744073709551615`, plus CPU `0.5037` and an active true firmware flag with no issue. [Partial](examples/overview-partial.json) includes all three issue codes, null group leaves, independent temperature, four null firmware flags and an interface-local failure. [Stale](examples/overview-stale.json), [expired](examples/overview-expired.json) and [never successful](examples/overview-never-success.json) exercise Board envelopes. The complete/partial snapshots seed the shared S01 corpus; the other envelope states remain documentary S02/S03 examples, not executed API tests.

The [shared corpus](../testdata/README.md) now covers B-01–B-04 snapshot rules, including network null, unavailable groups, useful-observation rules, interfaces, issues, presence/nulls, exact tokens, unknown members, Unicode and defensive boundaries. HTTP provider outcomes remain S02. The exact upstream tests establishing those rules are listed in [Health evidence](health-baseline.md).
