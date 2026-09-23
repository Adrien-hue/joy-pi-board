# Overview UX contract

Status: **established behavior plus proposed presentation details; not implemented**. Normative for JPB-017–019. [Board API](board-api-v0.1.0.md) defines server state; [Health integration](health-integration-v1.0.md) defines metric meaning. Do not derive a global health score from either.

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

**Proposed P-06 text mapping:** timeout → “Health did not respond in time”; connection_failed → “Health could not be reached”; upstream_error → “Health returned an error”; invalid_response → “Health returned invalid data”; unsupported_schema_version → “Health uses an unsupported schema version”. Do not display raw networking errors. Metric issues retain the provider's safe message as text and the associated measure label.

## Scheduling and browser freshness

Established: fetch overview on opening and every 5 seconds, with no overlapping refresh requests. No immediate retries. **S03 scheduling realization, still proposed**: start on a 5 s cadence measured from mount; skip a tick while an attempt is in flight, do not queue catch-up requests. Use a 2 s AbortController deadline for Board fetch including body reading; it is separate from Board's 1 s upstream deadline. Completion/error does not schedule an immediate retry; next scheduled tick is the next attempt. Cleanup timers and requests on unmount; guard publication with a request generation so a late aborted response cannot overwrite a later one. Tests must account for React development effect cleanup without intentional duplicate production requests.

Under the P-05 protocol approved by the S02 implementation instruction (browser code still NOT RUN), at request start record monotonic time `m0`; on response completion record `m1`. Validate overview and the age header from [Board API](board-api-v0.1.0.md). For a supplied snapshot:

```text
age_at_receipt_upper_bound = header_age_ms + elapsed_request_ms
age_now_upper_bound = age_at_receipt_upper_bound + elapsed_since_receipt_ms
presentable = valid_snapshot AND age_now_upper_bound <= 30000
```

The full round trip deliberately overestimates transit age. Rounding up the header cannot extend display lifetime. This bound is the same for first-load stale and subsequent stale responses, independent of browser/Pi clock offset and even when the Pi wall clock changes. Never infer a new successful observation from an unchanged stale response. Failed Board requests do not reset the age anchor. A fresh valid response from Board replaces the local envelope atomically; a `none` response immediately removes metrics even if a previous local bound had time left.

Run an expiry timer independently of fetch completion. Recompute age before any render or use of retained metrics, including after failed requests. At >30 s hide all values, not just their status badge. If no successful refresh arrives, an old `available/current` response cannot stay live indefinitely: after a fetch failure show stale/unknown reachability immediately; at expiry show no presentable data even if the previous fetch said available. At exactly 30 s data is permitted; a conservative age estimate may hide slightly earlier than the server's true bound.

Browser `performance.now()` supplies elapsed time, not cross-machine wall alignment. For lifecycle robustness, the approved age protocol uses the larger of monotonic delta and a nonnegative local wall delta; a local backward wall jump or invalid delta invalidates the retained age anchor and hides metrics until a valid response. Positive clock jumps can hide data early. On `pagehide`, freeze/suspension or entering hidden state, invalidate visible freshness and hide retained metrics; on resume/visibility return, keep metrics hidden until a valid response supplies a new age bound. Resume participates in the next 5 s tick, without overlapping or catch-up requests. Treat unverifiable suspension as expired. Test skew, wall changes, hidden tabs and suspend/resume; never let timer throttling preserve a live available label. Timestamps can be displayed in local time with an explicit timezone, but are not the elapsed-time authority.

Missing/invalid age information under the approved P-05 protocol is a Board contract error: fail closed on metric display. S02 implements the server header; browser interpretation remains S03, NOT RUN. The exact scheduling realization and 2 s browser AbortController choice above remain S03 review items; they are not needed by the S02 server. See [approval scope](decisions.md).

## Units and formatting

Established semantics: `0.5037` CPU means approximately **0.5%**, not 50.37%. Null means unavailable, never zero/false/no problem. Load is dimensionless, not a percentage. Network numbers are cumulative bytes, not bytes per second. True firmware flags are observations, not collection issues. Active and since-boot flags must never share an ambiguous single indicator.

**Proposed P-06 conventions:** CPU and derived memory/storage percentages to one decimal; temperature to one decimal °C; load to two decimals; logical CPU count exact integer. Bytes use B, KiB, MiB, GiB, TiB (powers of 1024), with at most one fractional digit and exact bytes accessible for all capacities/counters. RX/TX labels include “cumulative”; show exact decimal counters in an accessible detail using bigint arithmetic. Uptime uses integer days/hours/minutes/seconds from the received value, not a fabricated extrapolated uptime. Stale age updates in seconds with an “approximately” cue when showing the conservative upper bound. Display observation time and last-success time with distinct labels.

Byte percentage exists only when required values are non-null and total > 0. For total 0, show known byte values and percentage “Not applicable”; do not produce Infinity/NaN. Derived rounding belongs solely to presentation; API values remain unchanged. Use the integration contract's bigint ratio/division strategy, including for large totals.

Boolean display is “Yes”, “No”, or “Unavailable”; false is not missing. `unknown` network state is an available observed string, distinct from null acquisition failure. Unrecognized non-empty state strings under P-03 are shown as neutral text. Show issues alongside measures, optionally in a summary; a group issue covers its three unavailable leaves without inventing three new issues.

## Accessibility and review checks

Established states require text or meaningful icons with accessible labels in addition to color. **Proposed P-06 checks:** semantic heading/region structure, keyboard-operable details, visible focus, accessible status names, adequate contrast (WCAG AA), 200% text zoom and reduced-motion support. Announce state transitions through a restrained live region; do not announce every five-second metric change or steal focus on refresh. No required animations, images, remote fonts or third-party runtime resources.

Plan visual and interaction cases at 360/390 px mobile, 768 px tablet and 1280 px desktop; include long names, all-null groups, maximum counters, multiple interfaces and long issues. See B-10/B-11 and physical C-05 in [validation](validation-v0.1.0.md). These sizes/checks are proposed acceptance procedures, not completed design testing.
