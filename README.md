# Joy Pi Board

Joy Pi Board is the second microservice of Joy Pi Home: a lightweight web overview of observations supplied by Joy Pi Health, for trusted devices on the local network.

**Board presents. Health observes.**

## Current state

This repository contains the **consolidated documentary foundation for v0.1.0**. The application, build manifests, CI, packages and deployment service are not implemented. Sprint 00 is prepared, **not executed**. All planned product validation is **NOT RUN**; no release candidate or hardware acceptance exists.

The target is one page (`/`), one host and one provider. React/TypeScript/Vite assets will be embedded in a Go HTTP server, shipped as one binary. Production targets Raspberry Pi OS 64-bit / Trixie on Linux ARM64; Board acceptance uses a Raspberry Pi 3B+.

Board will serve its UI even without Health. The browser will call Board on the same origin; Board will consume Health's existing `http://127.0.0.1:8080/v1/snapshot`. The proposed Board listener is `0.0.0.0:8081`, **pending ratification**. The Board port was not previously decided. There is no runnable service or installation command yet.

## Documentation and authority

| Document | Source of truth |
|---|---|
| [Specification](docs/specification-v0.1.0.md) | Product scope, requirement IDs, established budgets and status vocabulary |
| [Architecture](docs/architecture-v0.1.0.md) | Boundaries, data flow, proposed realization and build layout |
| [Board HTTP contract](docs/board-api-v0.1.0.md) | Overview envelope, errors, timestamps and resilience invariants |
| [Health integration contract](docs/health-integration-v1.0.md) | Consumed schema, exact integers, nulls, issues and typed validation |
| [Health evidence](docs/health-baseline.md) | Pinned upstream sources, inspected tests and documented discrepancies |
| [UX contract](docs/ux-v0.1.0.md) | Presentation, availability, freshness and browser failure states |
| [Validation and release](docs/validation-v0.1.0.md) | Planned tests, measurable gates, requirement-to-test-to-sprint traceability |
| [Roadmap](docs/roadmap-v0.1.0.md) | Dependency-based sprint objectives and completion criteria |
| [Sprint 00 preparation](docs/sprint-00.md) | Detailed future foundation tasks; no execution authorization implied |
| [Proposals to ratify](docs/decisions.md) | Complementary choices awaiting review |
| [Documentation review](docs/documentation-review.md) | Checks performed on these documents, distinct from product validation |

The user-established decisions are consolidated, not newly approved implementation choices. Sections marked **Proposed P-xx** require review. Upstream observations are marked **Verified H-xx** and pinned to Health commit `be7a0d824f62b94c842d8e5326110b1852c5a0bc`; they are not evidence of Board implementation or of passing upstream tests in this session.

No system collection, database, history, administration, authentication or public Internet exposure belongs to v0.1.0. Health is a read-only reference: do not modify it or import its Go `internal` packages. See [AGENTS.md](AGENTS.md).
