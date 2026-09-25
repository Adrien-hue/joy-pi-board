# Overview UX contract

Status: **S03 IMPLEMENTED; CLOSURE PENDING**. The current user implementation instruction approves the S03 scheduling and presentation/accessibility recommendations prepared on 2026-09-23. Established P-05 and server behavior remain fixed. See [Sprint 03](sprint-03.md) for tasks and closure evidence. Normative for JPB-017–019. [Board API](board-api-v0.1.0.md) defines server state; [Health integration](health-integration-v1.0.md) defines metric meaning. Do not derive a global health score from either.

## Layout and reading order

One page `/`, responsive from a 360 CSS-pixel viewport, with no horizontal overview scrolling. Use this order at all widths:

1. Health-provided hostname and availability. Before a snapshot exists, label the page “Joy Pi Board”; do not invent a hostname or model.
2. CPU utilization/count, memory, root filesystem and SoC temperature.
3. 1/5/15-minute load averages and uptime.
4. Network interfaces: name, operational state and cumulative RX/TX.
5. Raspberry Pi indicators, grouped into “Active at observation” and “Occurred since boot”.
6. Service information (Board and Health only) and freshness/timestamps.

Services are availability information, not service-management controls or placeholders for future providers. On narrow screens, stack cards and interface details, wrap long host/interface/message strings and exact counters, preserve readable units. Tables may become labeled lists; do not hide important status behind color or horizontal scrolling.

## Three separate state dimensions

- **Board reachability** is observed by this browser. A successful valid overview establishes reachable; network failure, timeout, Board non-200 or an invalid overview establishes a failed Board fetch with an appropriate local message. A received HTTP error proves HTTP reachability but no usable Board service response; distinguish that text from a network-unreachable message.
- **Health availability** is Board's latest attempted provider outcome. A received unavailable overview does not mean Board is down. If Board becomes unreachable, Health's current availability is unknown; any old value must be labeled “last reported”, never left as a live available badge.
- **Metric availability/freshness** comes from snapshot state, nulls/issues and elapsed age. A valid partial snapshot is current data with unavailable measures, not a failed service. Stale is prior data explicitly labeled; none hides all metrics.

## State behavior

| State | Visible content and transition |
|---|---|
| Initial loading | Page shell, loading text/skeleton, no fake zero metrics or presumed available services; request immediately. |
| Nominal | Current snapshot, hostname, available Health, explicit current/freshness label; no good-health score. |
| Partial | Same current/available service state; null measures show “Unavailable” with associated issue; other fields remain visible. No fallback to older non-null values. |
| Stale | Health unavailable reason label; previous snapshot clearly marked “Stale”, with age and last-success time. Its old nulls/issues are still old observations. |
| No snapshot | Hide metric values, percentages and old hostname as live data; show no presentable data, reason and last success if known. A remembered hostname may only be shown explicitly as last observed. No success ever: no timestamp invented. |
| Health recovered | Next valid refresh replaces the whole snapshot and issues, available/current restores automatically. A brief textual recovery announcement may be used; no restart or user action required. |
| Board unreachable / fetch failed | State explicitly that Board cannot currently be reached, or its response failed. Health state is unknown/last reported. Retain any still-eligible metrics only with stale labeling and age, then hide at expiry even without another server response. |

**Approved S03 P-06 text mapping:** timeout → “Health did not respond in time”; connection_failed → “Health could not be reached”; upstream_error → “Health returned an error”; invalid_response → “Health returned invalid data”; unsupported_schema_version → “Health uses an unsupported schema version”. Do not display raw networking errors. Metric issues retain the provider's safe message as text and the associated measure label.

## Scheduling and browser freshness

Established: fetch overview on opening and every 5 seconds, without overlapping refreshes or immediate retry. **Approved S03-R02** is defined once in [the scheduling work package](sprint-03.md#scheduling-and-lifecycle-realization): immediate visible mount; mount-anchored slots, skip busy/missed slots; 2 s full-body/validation deadline; request/lifecycle generations; abort and clean effects. Hide/abort on hidden/pagehide/freeze, suppress ticks while hidden, and resume only on the next strictly future slot. Initial pageshow, BFCache and StrictMode replay have distinct handling. No completed response is reused as a refresh.

Under the P-05 protocol approved by the S02 implementation instruction (browser implementation and actual evidence recorded in [S03 evidence](s03-evidence.md)), at request start record monotonic time `m0`; on response completion record `m1`. Validate overview and the age header from [Board API](board-api-v0.1.0.md). For a supplied snapshot:

```text
age_at_receipt_upper_bound = header_age_ms + elapsed_request_ms
age_now_upper_bound = age_at_receipt_upper_bound + elapsed_since_receipt_ms
presentable = valid_snapshot AND age_now_upper_bound <= 30000
```

The full round trip deliberately overestimates transit age. Rounding up the header cannot extend display lifetime. This bound is the same for first-load stale and subsequent stale responses, independent of browser/Pi clock offset and even when the Pi wall clock changes. Never infer a new successful observation from an unchanged stale response. Failed Board requests do not reset the age anchor. A fresh valid response from Board replaces the local envelope atomically; a `none` response immediately removes metrics even if a previous local bound had time left.

Run an expiry timer independently of fetch completion. Recompute age before any render or use of retained metrics, including after failed requests. At >30 s hide all values, not just their status badge. If no successful refresh arrives, an old `available/current` response cannot stay live indefinitely: after a fetch failure show stale/unknown reachability immediately; at expiry show no presentable data even if the previous fetch said available. At exactly 30 s data is permitted; a conservative age estimate may hide slightly earlier than the server's true bound.

Browser `performance.now()` supplies elapsed time, not cross-machine wall alignment. For lifecycle robustness, the approved age protocol uses the larger of monotonic delta and a nonnegative local wall delta; a local backward wall jump or invalid delta invalidates the retained age anchor and hides metrics until a valid response. Positive clock jumps can hide data early. On `pagehide`, freeze/suspension or entering hidden state, invalidate visible freshness and hide retained metrics; on resume/visibility return, keep metrics hidden until a valid response supplies a new age bound. Resume participates in the next 5 s tick, without overlapping or catch-up requests. Treat unverifiable suspension as expired. Test skew, wall changes, hidden tabs and suspend/resume; never let timer throttling preserve a live available label. Timestamps can be displayed in local time with an explicit timezone, but are not the elapsed-time authority.

Missing/invalid age information under the approved P-05 protocol is a Board contract error: fail closed on metric display. S02 implements the server header; S03 now implements browser interpretation under the current instruction's approval. Local deterministic/engine checks are recorded in the S03 ledger; required branded/floor/device evidence remains NOT RUN. The 2 s browser deadline does not change the S02 server budget. See [approval scope](decisions.md).

## Units and formatting

Established semantics: `0.5037` CPU means approximately **0.5%**, not 50.37%. Null means unavailable, never zero/false/no problem. Load is dimensionless, not a percentage. Network numbers are cumulative bytes, not bytes per second. True firmware flags are observations, not collection issues. Active and since-boot flags must never share an ambiguous single indicator.

**Approved S03 P-06 conventions:** CPU and derived memory/storage percentages to one decimal; temperature to one decimal °C; load to two decimals; logical CPU count exact integer. Bytes use B, KiB, MiB, GiB, TiB, PiB, EiB (powers of 1024, covering all uint64 values), with at most one fractional digit and exact bytes accessible for all capacities/counters. RX/TX labels include “cumulative”; show exact decimal counters in an accessible detail using bigint arithmetic. Uptime uses integer days/hours/minutes/seconds from the received value, not a fabricated extrapolated uptime. Stale age updates in seconds with an “approximately” cue when showing the conservative upper bound. Display observation time and last-success time with distinct labels.

Byte percentage exists only when required values are non-null and total > 0. For total 0, show known byte values and percentage “Not applicable”; do not produce Infinity/NaN. Derived rounding belongs solely to presentation; API values remain unchanged. Use the integration contract's bigint ratio/division strategy, including for large totals.

Boolean display is “Yes”, “No”, or “Unavailable”; false is not missing. `unknown` network state is an available observed string, distinct from null acquisition failure. Unrecognized non-empty state strings under P-03 are shown as neutral text. Show issues alongside measures, optionally in a summary; a group issue covers its three unavailable leaves without inventing three new issues.

## Accessibility and review checks

Established states require text or meaningful icons with accessible labels in addition to color. **Approved S03 P-06 checks:** semantic heading/region structure, keyboard-operable details, visible focus, accessible status names, adequate contrast (WCAG AA), 200% text zoom and reduced-motion support. Announce state transitions through a restrained live region; do not announce every five-second metric change or steal focus on refresh. No required animations, images, remote fonts or third-party runtime resources.

Plan visual and interaction cases at 360/390 px mobile, 768 px tablet and 1280 px desktop; include long names, all-null groups, maximum counters, multiple interfaces and long issues. See B-10/B-11 and physical C-05 in [validation](validation-v0.1.0.md). Local automated and reviewed-image subsets are recorded in S03 evidence; manual/native-browser/device checks remain open.

## Proposed S03 page and component contract

**P-06 presentation/accessibility subset, APPROVED for S03 (S03-R03)** by the current implementation instruction. This section retains the prepared design and labels; the delivered wording uses the same state distinctions. The implementation uses English, system fonts, restrained dark green text on light neutral backgrounds and measured solid-color contrast. This ordinary visual choice adds no theme setting or product behavior. Use English (`lang=en`), fixed en-GB number conventions and explicit UTC timestamps for consistent tests. Health messages remain verbatim text in their supplied language; do not translate by matching their wording. No automatic theme setting, animation, notification, refresh button or service control is needed in this scope.

```text
Joy Pi Home / Joy Pi Board
<Health hostname when eligible>
Board: response state   Health: latest/last-reported outcome
Observation: Current | Partial | Stale | No presentable data

CPU                 Memory              Root filesystem       SoC temperature
utilization / CPUs  used / available / total                    degrees Celsius

Load averages: 1 min / 5 min / 15 min     Uptime at observation

Network interfaces
<name>  State: <observed text>
Received (RX, cumulative)                Transmitted (TX, cumulative)

Raspberry Pi indicators
Active at observation                   Occurred since boot
Thermal throttling / Undervoltage        Thermal throttling / Undervoltage

Board / Joy Pi Health service information
Conservative data age / Last successful retrieval / Observed by Health /
Response generated by Board / measurement availability details
```

Proposed components: OverviewHeader/ResponseStatus, MetricCard/ByteGroup, LoadAndUptime, NetworkList/InterfaceRow, FirmwareIndicators, IssueText and ServiceFreshness. The view receives validated typed data plus a guarded presentation state; it does not fetch or recompute schema validity. Keep DOM/heading order identical at every width. One-column at 360/390 px; two primary cards per row around 768 px; four around 1280 px if text fits, with wrap-safe grid columns and no minimum content width. Secondary/network/firmware sections follow sequentially. Use a bounded content width with fluid gutters, not fixed card heights or clipped values.

Show every interface (up to 64), keep Health order, use interface name for React identity but current snapshot index for issue paths. No pagination, virtualization or hidden omitted rows is needed at this bound. An empty valid array says “No interfaces reported”; network null says “Network observations unavailable” with its group issue. Vertical scrolling is acceptable. Wrap full hostnames, interface state text, counters and messages (`min-width: 0`, overflow-wrap); do not ellipsize away the only exact value. Use generated opaque DOM IDs, not untrusted names/paths as markup identifiers. User/provider strings render only as React text nodes, never innerHTML or URL/style content.

### Concrete state copy and retention

| State / trigger | Recommended visible outcome |
|---|---|
| Initial | “Loading observations”; Board “Checking”; Health “Not yet known”; no invented hostname/zero values. Shell remains navigable. |
| Usable current, no issues | “Current observation”; Board “Responding”; Health “Available”; every supplied value shown. No “Healthy” badge. |
| Usable current, issues | “Current observation — some measurements unavailable”; Board/Health still responding/available. Null leaves show “Unavailable”, associated issue text once at its actual scope. |
| Usable stale | Board “Responding”; Health “Unavailable” plus fixed reason text; “Stale observation” and conservative age; only still-eligible values shown. |
| Usable none | Board “Responding”; Health “Unavailable”; “No presentable data”. Show last successful retrieval only if supplied; clear all old metrics immediately. |
| Board HTTP error | “Board returned HTTP <status>”; Health “Unknown — last reported …”; prior eligible metrics explicitly stale, then expire. This is not “could not connect”. |
| Network failure / 2 s timeout | “Board could not be reached” / “Board response timed out”; Health unknown/last reported; eligible previous metrics stale, no age reset. |
| Contract/age/clock error | “Board data could not be verified” / “Data age could not be verified”; metrics masked immediately; no fallback to an old live badge. No raw exception or response body. |
| Hidden / resume waiting | “Observations paused; waiting for a fresh response”; no metric values or live Health availability until a new valid response with eligible age. |
| Locally expired usable response | Board “Last response received”; Health “Last reported …”; “No presentable data — observation expired”. Received current/stale is historical metadata, never an indefinitely live state. |
| Recovery | Atomically replace the view; one polite “Observations available again” announcement when an outage/paused/none episode ends. A normal successful poll does not announce recovery. |

When no metric snapshot is presentable, heading returns to Joy Pi Board. Optional last hostname belongs only in explicitly labeled historical metadata; recommend omitting it from the headline. Do not clear useful last-success metadata merely because snapshot values expire. Display Board HTTP failures separately from a valid Health failure, and availability failures separately from observed true firmware flags. An empty issues array never becomes “No health problems”.

### Exact presentation algorithms

Retain source bigint/number/null values unchanged. No counter passes through Number, including derived percentage inputs. Pure formatting functions own rounding, not components or transport.

| Field | Approved formatting and exact-value access |
|---|---|
| CPU utilization | One fractional digit plus %, directly from the supplied percentage: 0.5037 → 0.5%. No multiplication by 100, thresholds or warning color. |
| Logical CPUs | decimalUInt64 exact digits, no conversion to floating point. |
| Memory/root bytes | Show used, available and total. Binary summary B through EiB; for B use integer text, otherwise one decimal digit. A keyboard-operable native details section reveals all exact ungrouped decimal byte values with units. “Available” means the supplied unprivileged availability, not an invented free value. |
| Memory/root percentage | Only if all required leaves are non-null and total >0. For used u and total t, tenths of percent q = floor((2 × u × 1000 + t) / (2 × t)), using bigint throughout; render q/10 and q mod 10 as decimal digits. This is round-nearest, ties up for nonnegative ratios. Zero total → “Not applicable”; known 0 bytes remains visible. |
| RX/TX | Label “Received (RX, cumulative)” and “Transmitted (TX, cumulative)”; show exact ungrouped decimal bytes as visible text and optional binary summary. No rates/deltas. `9007199254740993` and `18446744073709551615` must be readable/copyable without tooltip-only access. |
| Binary summaries | Choose largest divisor 1024^k <= value (k 0..6); for one decimal use floor((2 × value × 10 + divisor)/(2 × divisor)). Carry a rounded 1024.0 into the next unit when available. No Number conversion. Preserve exact digits in details, including 0 and uint64 max. |
| SoC temperature / load | Temperature one fractional digit °C; load two fractional digits with “1 min / 5 min / 15 min”, dimensionless. Use explicit Intl.NumberFormat fraction options for floating fields; normalize displayed negative zero to zero without changing stored value. For extremely large finite floats, a scientific summary is permitted with the full Number decimal representation in details; no clipping or invented bound. |
| Uptime | bigint quotient/remainder into days, hours, minutes, seconds; all uint64 seconds accepted. No ticking extrapolation. Exact seconds remain in details; label “Uptime at observation”. |
| Firmware | “Yes” / “No” / “Unavailable”; active and since-boot rows separate. Null is not No; true is not a collection issue or global diagnosis. |
| Network state | Display supplied nonempty string as neutral text; up/down/unknown are observations, not a closed enum. Unknown string is not acquisition null. |
| Timestamps / age | Explicit labels “Observed by Health”, “Last successful retrieval by Board”, “Response generated by Board”. Show preserved UTC text (wrappable, fractional precision retained); missing time “Not yet”. Conservative age label “Data age: at most approximately … s”, rounded up to whole seconds, but eligibility uses the exact millisecond bound. Never derive age from these timestamps. |

Use path/code to associate issues. Render each issue's message **once** at its metric/group scope; use aria-describedby to link affected leaves to that text. For /load, /memory and /root_filesystem, one group issue describes three unavailable leaves. Firmware has four actual issues: retain their separate paths; sharing a visually grouped message is permitted only with all four associations preserved, not changing the issue count. Recommend one message per supplied issue for simplest auditability. Network index paths attach to the current snapshot rows; removed rows cannot retain prior issue text. A footer summary may count and link to issues, not repeat every message. Keep provider order in the summary. No message matching, invented severity or duplicate synthetic issues.

### Accessibility acceptance

Use semantic main/section/headings, definition lists or correctly labeled tables, and native details/summary for exact values. Every details control has an explicit metric-specific accessible name and visible keyboard focus. Preserve focus and open details across ordinary value refreshes by stable semantic identity; if the focused interface disappears, move focus only as needed to the network heading and announce that removal once. Do not move focus on normal refresh/recovery. On expiry remove values from both visual and accessibility presentation, including expanded details.

Target the relevant [WCAG 2.2](https://www.w3.org/TR/WCAG22/) AA checks: normal text contrast >=4.5:1, large text >=3:1, focus/control boundaries >=3:1, no color-only state, keyboard operation and meaningful labels. Verify 200% zoom and large text, responsive reflow, reduced-motion preference, reading order and no horizontal overview overflow. Avoid animations altogether. Give one polite atomic status region concise Board/Health/freshness/measurement-availability transitions; deduplicate unchanged states and coalesce related changes from one refresh. Never announce every metric/age tick, read all issues automatically, or use an assertive alert for routine loss/recovery. Screen-reader checks complement visual and DOM checks; an automated accessibility scan alone cannot certify these behaviors.

S03 owns these software/browser/visual checks; S04 repeats installed-device C-05 and physical acceptance. [S03-T08–T12](validation-v0.1.0.md#prepared-s03-validation-matrix) define proof, including 360/390, 768 and 1280 px, 200% zoom, long strings, all 64 interfaces and many issues. Executed local subsets and remaining browser/manual requirements are recorded in [S03 evidence](s03-evidence.md). Physical measurement and actual server compression protocols are a separate P-06/S04 review, not ratified by S03 approval.

Implementation note: exact-value details retain their semantic identity when a numeric measure becomes unavailable, presenting Unavailable inside the existing disclosure so focus/open state survives a partial refresh. The interface rows contain no removable interactive controls. If expiry/lifecycle masking removes the currently focused metric region, visible-page focus is safely returned to service information without scrolling. Normal refresh does not move focus. The load group links its single issue to all three leaves. No unknown network state receives invented semantics.
