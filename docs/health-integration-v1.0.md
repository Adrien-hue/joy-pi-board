# Integration with Health snapshot schema 1.0

Status: **consolidated consumed contract**, with implementation proposals explicitly marked. Normative for Board's consumed field names, units, nulls and issue semantics. Source: [pinned Health evidence H-01–H-09](health-baseline.md), never Health Go package imports. Requirements JPB-014–016.

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

**Proposed P-03 realization.** Board owns independent DTOs. Plain `*uint64` alone cannot distinguish missing from explicit null on input. Decode every required member with presence tracking, including required object/array nodes. A documentation-only type sketch:

```go
type Field[T any] struct {
    Seen  bool // set only if the exact JSON key occurs
    Null  bool // explicit JSON null; valid only for nullable fields
    Value T
}
type Issue struct { Path, Code, Message Field[string] }
type Host struct { Hostname Field[string] }
type CPU struct {
    UtilizationPercent Field[float64]
    LogicalCPUCount    Field[uint64]
}
type Load struct { OneMinute, FiveMinutes, FifteenMinutes Field[float64] }
type ByteGroup struct { TotalBytes, AvailableBytes, UsedBytes Field[uint64] }
type NetworkInterface struct {
    Name, State Field[string]
    RXBytes, TXBytes Field[uint64]
}
type Network struct { Interfaces Field[[]NetworkInterface] }
type RaspberryPi struct {
    SoCTemperatureCelsius Field[float64]
    ThermalThrottlingActive, ThermalThrottlingOccurredSinceBoot Field[bool]
    UndervoltageActive, UndervoltageOccurredSinceBoot Field[bool]
}
type HealthSnapshot struct {
    SchemaVersion, ObservedAt Field[string]
    Host Field[Host]
    CPU Field[CPU]
    Load Field[Load]
    Memory, RootFilesystem Field[ByteGroup]
    UptimeSeconds Field[uint64]
    Network Field[Network]
    RaspberryPi Field[RaspberryPi]
    Issues Field[[]Issue]
}
```

The sketch is not executable decoder code: the future decoder must bind **exact snake_case keys from the table**, through explicit tags or key dispatch, without Go's permissive case-insensitive matching. All `Seen` flags must be true; check `Null` against each field's permitted nullability before reading `Value`. Do not allow the zero value of `Value` to act as a missing-field marker. Reject null object/array nodes except `network`. Use `uint64` decoding or `json.Number` plus `strconv.ParseUint` from the original token, never a `float64` intermediary for integers.

After validation, an immutable copy of the original snapshot bytes (`json.RawMessage`) can be embedded as the **object** value of `health.snapshot` and retained in the one-slot cache. Do not marshal bytes as base64 or a quoted JSON string. Original values, issues and numeric tokens remain unchanged; object whitespace/order need not be preserved contractually. Never forward unvalidated raw JSON.

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

**Proposed P-03:** a small bounded parser of the JSON grammar retains every number as a `NumberToken { raw: string }` until schema validation. Fetch the body as bytes (bounded streaming read), decode UTF-8 strictly, then parse. `response.text()` with a separate enforced byte limit is also viable; **`response.json()` is not**. `JSON.parse` with an ordinary reviver after Number rounding is equally insufficient. Never rewrite numbers with a regular expression: digits in strings, escapes, exponents and nested structures require proper tokenization.

At a uint64 field, accept only an unquoted nonnegative decimal integer token (`0` or a nonzero digit followed by digits), then `BigInt(raw)` and range-check. Fractional/exponent notation and negative tokens at integer fields are rejected under P-03's strict emitted-token policy. At a floating metric, convert its raw number token to `Number`, check finite/range, and retain the documented floating-point semantics. Strings such as `"9007199254740993"` are invalid at numeric fields. The browser must never call `Number(counter)` before formatting or calculating a ratio. The same fixtures drive Go and TypeScript checks.

This owned tokenizer/parser is realizable without dependencies but creates correctness work: grammar, escapes/Unicode, duplicate keys, bounds and differential/fuzz tests are S01 deliverables. A focused lossless JSON library is a permissible **alternative proposal**, only after its raw-token API, license, pinned version, duplicate/UTF-8 behavior and compressed size have been reviewed. No library or version is installed/selected now. ES bigint support must be included in S00's browser target and type-check settings.

Formatting preserves integers: exact decimal text from `bigint.toString()`; binary unit summaries and percentages use bigint division/remainder with explicit rounding. For tenths of a percent, when total > 0, use `(used * 1000n + total / 2n) / total`, then format the resulting integer tenths; no unsafe cast of the operands. Display the exact cumulative byte number alongside, or in an accessible expandable detail, so a rounded GiB summary never replaces the exact counter. See [UX](ux-v0.1.0.md).

## Validator policy and failure boundary

**Proposed P-03 supplemental strictness:** bound Health body to 65,536 bytes, read at most limit+1, JSON nesting to 32, and a Board overview body to 80 KiB in the browser. Reject duplicate or unknown members, trailing JSON/text, malformed UTF-8, invalid strings, non-object roots, missing/null-required members, wrong types, out-of-range values and incoherent issues. Empty `issues` and interface arrays are allowed under their semantic constraints. Do not normalize, repair, sort or merge malformed snapshots.

Validate a syntactically well-formed root and a present string `schema_version` before interpreting version-specific metric fields. A string other than `1.0` is `unsupported_schema_version`; an absent/non-string version is `invalid_response`. Remaining schema 1.0 violations are `invalid_response`. A structurally malformed document is always `invalid_response`, even if a fragment contains an unknown version. A future compatible schema must be reviewed explicitly rather than guessed.

Network state is deliberately an open non-empty string in this proposed Board validator; recognize up/down/unknown in presentation but render other text neutrally. The stricter three-value coordinator behavior is H-06, not an invented public enum. Validation of text messages ensures non-empty UTF-8 and firmware-group equality; no general message-literal matching.

## Documentary examples and future fixtures

[Complete](examples/overview-current.json) includes exact RX `9007199254740993` and TX `18446744073709551615`, plus CPU `0.5037` and an active true firmware flag with no issue. [Partial](examples/overview-partial.json) includes all three issue codes, null group leaves, independent temperature, four null firmware flags and an interface-local failure. [Stale](examples/overview-stale.json), [expired](examples/overview-expired.json) and [never successful](examples/overview-never-success.json) exercise Board envelopes. These are documentation examples, not an executed test suite.

Future B-01–B-05 must also cover network null with `/network`, each group's all-null form, all-unavailable rejection, empty/unsorted/duplicate interfaces, issue reordering, absent versus null, exact numeric boundaries and HTTP provider outcomes. The exact upstream tests establishing those rules are listed in [Health evidence](health-baseline.md).
