# Joy Pi Board v0.1.0 specification

Status: **consolidated product requirements; S00/S01 COMPLETE, S02 locally implemented; final product incomplete**. [S00 evidence](s00-evidence.md) identifies the locally and CI-verified scaffold commit, closure criteria and remaining product validations. S01 in-memory contracts are complete with local and verified hosted CI evidence; [S01 evidence](s01-evidence.md) records closure and unexecuted product gates. S02 is [implemented and verified locally; hosted CI closure pending](sprint-02.md). This is normative for established product scope, requirement IDs and budgets. Details are owned by the documents linked below, not duplicated as competing definitions.

## Authority and status

- **Established**: decisions supplied for this v0.1.0 baseline; requirements below use “must” in that sense.
- **Verified H-xx**: behavior inspected in the pinned Health source; see [evidence](health-baseline.md). Reading a test is not executing it.
- **Proposed P-xx**: complementary realization or measurement protocol in [decisions](decisions.md); review pending, not a historical decision.
- **PLANNED / NOT RUN**: future task/test without execution evidence. **PASS / FAIL** requires an identified run and evidence. Documentation completeness does not imply implementation completeness.

The source of truth for each subject is listed in the [README](../README.md). Conflicts with an established requirement must be resolved explicitly, not by treating a proposal as an override. Requirement IDs remain stable across editorial changes. The [validation traceability table](validation-v0.1.0.md#requirement-traceability) assigns each requirement planned tests and a responsible sprint.

## Purpose and responsibilities

Board is the second Joy Pi Home microservice. It makes Health observations readable on trusted LAN devices without exposing Health directly to browsers. It presents facts, acquisition availability and data freshness; it does not diagnose the machine or collect observations.

| ID | Established requirement |
|---|---|
| JPB-001 | Provide one main page at `/`, one host and only the Joy Pi Health provider. |
| JPB-002 | Present Health hostname; CPU percentage/count; memory; root filesystem; SoC temperature; 1/5/15-minute load; uptime; network interface states and cumulative RX/TX; active and since-boot throttling/undervoltage. Do not invent a hardware model. |
| JPB-003 | No direct collection through `/proc`, `/sys`, commands or sensors; no Health internal-package imports or modifications. |
| JPB-004 | Exclude history, time-series graphs, persistence/database, alerts/notifications, arbitrary health thresholds/scores, service configuration UI, start/stop/restart/reboot, Linux administration, user management/authentication, public Internet exposure, multi-host, plugins and future Data/Log/Hole/Network integrations. Startup configuration is allowed. |
| JPB-005 | React + TypeScript built with Vite; CSS/CSS Modules and system fonts; local React state and Fetch; Go standard-library-first backend with `net/http`; compiled assets through `go:embed`; one deployable binary/service. No backend framework, SSR, global state library, large UI or chart library without demonstrated need and explicit scope review. |
| JPB-006 | Browser → Board → Health; same-origin browser requests only, no CORS needed. Configurable Board listener and Health endpoint, preserving Health default `http://127.0.0.1:8080/v1/snapshot`. `.local` relies on existing name resolution, not Board mDNS. |
| JPB-007 | Start and serve the UI without Health. Browser loads immediately and refreshes every 5 seconds without overlapping refresh requests. |
| JPB-008 | Fetch Health on demand only; no autonomous Board polling, immediate retry or hidden cache hit instead of a refresh attempt. Global upstream timeout is 1 second including body reading. |
| JPB-009 | Keep only the latest valid snapshot and last-success metadata in RAM; no history/persistent cache; restart loses them. Partial valid success replaces the entire prior snapshot, including nulls. |
| JPB-010 | After Health failure, expose stale snapshot only while last success age is ≤30 seconds; at >30 seconds hide metrics/snapshot but retain last-success timestamp. Before first success both snapshot and timestamp are null. Recover automatically. |
| JPB-011 | Protect concurrent access and result ordering; no late result may regress the published observation. The realization is a documented proposal, not an implicit cache policy. |
| JPB-012 | Serve `GET /api/v1/overview` with the exact established envelope, enums and invariants in [Board API](board-api-v0.1.0.md). Health failure alone is HTTP 200; send `Content-Type: application/json` and `Cache-Control: no-store`. Separate Board HTTP errors; never expose raw network/system errors. |
| JPB-013 | Preserve Health `observed_at`; set Board `last_success_at` on successful reception and validation and `generated_at` on response generation. Expiry uses reliable elapsed time on server and browser, not browser/Pi wall-clock alignment. |
| JPB-014 | Consume schema `1.0` against the pinned baseline using typed validation, required-field presence checks and exact null semantics. Preserve names, units, values and issues; no fabricated complete/partial or `status` field. |
| JPB-015 | Preserve all unsigned 64-bit integers as unquoted JSON numbers and exactly in browser parsing/presentation; validate `9007199254740993` and `18446744073709551615`. `0` and `false` are present values. |
| JPB-016 | Enforce the verified null/group/firmware/issue associations and deterministic ordering in [Health integration](health-integration-v1.0.md); issue logic uses code/path, not message text. |
| JPB-017 | Distinguish Health unavailable, Board unreachable to the browser and valid partial observations. Browser must stop showing old metrics after the freshness limit even when Board requests fail. |
| JPB-018 | Responsive overview from 360 px without horizontal scrolling; states include textual/non-color cues. Follow [UX](ux-v0.1.0.md) hierarchy and initial, nominal, partial, stale, none, recovered and unreachable states. |
| JPB-019 | CPU value is already percent; derived memory/storage percentages handle null/zero total; RX/TX are cumulative, never rates. Distinguish active/since-boot flags, preserve unavailable states, show stale age and no arbitrary health verdict. |
| JPB-020 | Meet the established performance/resource/capacity budgets below on the Board reference target. |
| JPB-021 | No runtime Internet dependency: no CDN, external font, telemetry or remote script. Provide CSP, X-Content-Type-Options and Referrer-Policy suited to a same-origin LAN application. |
| JPB-022 | Log to stdout/stderr for journald: startup, shutdown and significant transitions. No full snapshots or per-poll success logging by default. |
| JPB-023 | Target Raspberry Pi OS 64-bit / Trixie, Linux ARM64; future autonomous ARM64 binary, Debian package, systemd unit with dedicated unprivileged user, checksums and build metadata. No production Node/npm/Vite or mandatory reverse proxy. Packaging must not depend on Health availability to start Board. |
| JPB-024 | Official release requires Classes A, B and C PASS on the exact candidate bytes/checksums. Physical tests use Raspberry Pi 3B+, distinct from Health's own hardware gates. Unexecuted physical validation is NOT RUN. |
| JPB-025 | Respect the currently authorized sprint boundary. S00/S01 implementation and closure are historical completed work with separate evidence; the current instruction authorizes S02 implementation only. Do not implement S03/S04, publish a release or deploy services. Preserve historical reviews and exact-commit evidence; preparation never implies execution. |

## Established non-functional budgets

| Budget | Limit |
|---|---:|
| Stable Board process RSS | ≤50 MiB |
| Board CPU without client | ≤1% of one core |
| Board mean CPU, one client refreshing every 5 s | ≤2% of one core |
| Readiness, also with Health absent | ≤2 s |
| Graceful shutdown | ≤5 s |
| Board overhead excluding Health call | ≤50 ms p95 |
| Overview with nominal Health | ≤500 ms p95 |
| Immediately detectable connection failure | ≤250 ms p95 |
| Entire Board → Health call | 1 s timeout |
| Initial compressed JS + CSS | ≤250 KiB |
| All initial compressed resources | ≤500 KiB |
| Planned client capacity | 10 clients, each refreshing every 5 s |

MiB = 2^20 bytes; KiB = 2^10 bytes. The 250 ms gate does **not** apply to a genuine network timeout. CPU is process CPU time divided by elapsed wall time, expressed as percent of one core, not divided by four on the Pi. No separate ten-client CPU percentage limit has been established; [validation](validation-v0.1.0.md) proposes a concrete capacity protocol without inventing one.

Operational measurement protocols and remaining supplemental transport defaults are proposals pending review; see the updated [decision register](decisions.md). Only the minimal shell's asset sizes have been measured in S00, not final product or physical budgets. The exact-integer convention is adopted in [ADR-001](adr-001-exact-json-integers.md).
