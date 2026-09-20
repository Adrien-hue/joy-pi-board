# Documentary foundation review

Date: 2026-09-20. Scope: repository documentation only, **not product acceptance**. Repository identity was verified as `C:/src/joy-pi-board`, origin `https://github.com/Adrien-hue/joy-pi-board.git`, starting commit `4fdeb4a` (`Initial commit`). Initial tree contained only a minimal README; there were no existing repository/ancestor AGENTS instructions found and no pending changes to preserve. The README was expanded, not replaced with an application scaffold.

## Review scope

- Read the six requested Health sources at commit `be7a0d824f62b94c842d8e5326110b1852c5a0bc`; pin source URLs and file hashes in [Health evidence](health-baseline.md). No newer revision or modified Health content used.
- Check local Markdown links/anchors, parse standalone JSON examples and JSON code blocks, compare example required keys/null groups/issue association and order/envelope states, and verify exact integer literals with Python integer parsing.
- Check that every established JPB requirement has planned validation and sprint ownership; review five reason values, availability/snapshot-state enums, timestamps, 30-second boundary, HTTP 200 on Health-only failure, default ports, ARM64/Trixie and Board's Pi 3B+ acceptance target across documents.
- Check the repository diff for whitespace errors and the changed-file inventory for documentary scope. Ad hoc documentation verification is not committed as executable CI or application tests.

## Contradictions resolved in Board documentation

| Earlier conceptual ambiguity | Consolidated treatment |
|---|---|
| Non-nullable metric examples or absence conflated with null/zero/false | Every required metric key is presence-checked; exact nullable types; zero/false preserved. |
| TypeScript counters as plain number | All Health integers become bigint before Number conversion; unquoted exact JSON numbers remain on the wire. |
| `status` substituted for provider availability | Only `availability` and the specified `snapshot_state` enums in overview. |
| Public “group null” phrasing | load/memory/root objects with null leaves and one group issue; actual null network group separately documented. |
| Specialized public firmware issue message | Actual emitted generic messages recorded; consumer logic uses path/code. |
| Observed network state values interpreted as a closed public type | up/down/unknown are verified implementation values; open string validator is a labeled proposal. |
| Board listener conflicting with Health | Health stays on 8080; Board 8081 explicitly proposed, never described as historically approved. |
| Cached/partial data treated as fresh or merged | Partial replaces the whole snapshot; monotonic stale expiry on both server and browser; no indefinite available state after Board failure. |
| Build/mocks implying release or hardware acceptance | All A/B/C product checks NOT RUN; release tied to actual Pi 3B+ evidence and exact candidate checksums. |

## Execution boundary and status

Only documentation and documentary JSON examples were written. No Go/React application, dependency manifest, runnable CI workflow, systemd unit, Debian package or release artifact was created. No Sprint 00 task, product test, build, benchmark, deployment or publishing operation was executed. The physical target was not accessed. Health's tests were inspected as source, not run.

The [Sprint 00 work package](sprint-00.md) is prepared; its task statuses remain NOT STARTED. [P-01–P-06](decisions.md) remain pending ratification. The current deliverable is a coherent specification and implementation/validation preparation, not a release-ready application.

Final documentary verification results: local links/anchors across 13 Markdown files, five standalone example JSON documents and their invariants, JSON code blocks, exact integer examples, all 25 requirement mappings, definitions of 23 planned validation IDs and whitespace checks passed. The final inventory contains 18 files, all Markdown or documentary JSON. This narrow result does not change any A/B/C status from NOT RUN.
