# Repository guidance

- This is **Joy Pi Board**, not Joy Pi Health. Confirm repository identity before changing files. Preserve existing work.
- **Board presents. Health observes.** No `/proc`, `/sys`, system commands or sensors for metric collection; no imports from Health's Go `internal` packages. Health is read-only.
- Start with [README.md](README.md) and its authority map. Product scope and stable requirement IDs live in [the specification](docs/specification-v0.1.0.md); wire behavior in the [Board](docs/board-api-v0.1.0.md) and [Health](docs/health-integration-v1.0.md) contracts.
- Pin upstream contract changes explicitly through [the evidence record](docs/health-baseline.md). Preserve nulls, exact unsigned 64-bit JSON integers and issue order. Never invent hardware model or health scores.
- Keep established decisions, verified upstream facts and [proposals](docs/decisions.md) distinct. Technical documentation is English. Do not mark planned tests PASS.
- The current assignment is documentation only. [Sprint 00](docs/sprint-00.md) is prepared, not executed: no application skeleton, dependency manifests, executable CI, packages, deployment or release in this assignment. Later implementation requires a new task authorizing that work.
- Update cross-references and requirement/test/sprint traceability together. Physical acceptance requires the exact candidate checksum and actual Raspberry Pi 3B+ evidence; mocks cannot satisfy Class C.
