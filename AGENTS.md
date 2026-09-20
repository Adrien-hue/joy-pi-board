# Repository guidance

- This is **Joy Pi Board**, not Joy Pi Health. Confirm repository identity before changing files. Preserve existing work.
- **Board presents. Health observes.** No `/proc`, `/sys`, system commands or sensors for metric collection; no imports from Health's Go `internal` packages. Health is read-only.
- Start with [README.md](README.md) and its authority map. Product scope and stable requirement IDs live in [the specification](docs/specification-v0.1.0.md); wire behavior in the [Board](docs/board-api-v0.1.0.md) and [Health](docs/health-integration-v1.0.md) contracts.
- Pin upstream contract changes explicitly through [the evidence record](docs/health-baseline.md). Preserve nulls, exact unsigned 64-bit JSON integers and issue order. Never invent hardware model or health scores.
- Keep established decisions, verified upstream facts and [proposals](docs/decisions.md) distinct. Technical documentation is English. Do not mark planned tests PASS.
- [Sprint 00](docs/sprint-00.md) is complete, with local and verified hosted CI evidence. The current assignment is documentary closure only: do not modify application code, dependencies or workflows. S01 is ready but not started; no S01–S04 implementation is authorized (Health client/validator/parser, overview API, cache, dashboard, packaging or deployment).
- Use [foundation commands](docs/s00-foundation.md); build frontend before Go embed. Exact tools are checked by `node scripts/tasks.mjs tools` with implicit Go switching disabled. [ADR-001](docs/adr-001-exact-json-integers.md) adopts Health's numeric convention; unknown-member policy and parser choice remain separate S01 concerns. Never add an unused parser dependency or a general-purpose homegrown parser in S00.
- Update cross-references and requirement/test/sprint traceability together. Physical acceptance requires the exact candidate checksum and actual Raspberry Pi 3B+ evidence; mocks cannot satisfy Class C.
