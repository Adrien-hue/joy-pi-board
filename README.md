# Joy Pi Board

Joy Pi Board is the second microservice of Joy Pi Home: a lightweight web overview of observations supplied by Joy Pi Health, for trusted devices on the local network.

**Board presents. Health observes.**

## Current state

This repository contains the **v0.1.0 specifications and implemented Sprint 00 build shell**: a minimal React page served by one Go binary with embedded assets, exact tool/dependency locks, useful source/integration tests and a GitHub Actions workflow. S00 closure is pending actual CI execution; see [S00 evidence](docs/s00-evidence.md). The Health client, overview API, parser/validator, cache, dashboard, packages and production service are not implemented. Functional and physical validation remain NOT RUN; there is no release candidate.

The product target is one page (`/`), one host and one provider. React/TypeScript/Vite assets are already embedded for the S00 shell. Production targets Raspberry Pi OS 64-bit / Trixie on Linux ARM64; Board acceptance uses a Raspberry Pi 3B+.

The shell serves its UI without Health and makes no provider calls. The **approved** Board listener defaults to `0.0.0.0:8081`; the future Health endpoint remains `http://127.0.0.1:8080/v1/snapshot`. The minimal page identifies the development foundation and displays no fake health values. `/api/v1/overview` returns 404 until S02.

## Build the S00 shell

Use Go **1.27.1**, Node.js **24.21.0**, npm **11.19.0**. From the repository root:

```text
node scripts/tasks.mjs tools
npm --prefix web ci
npm --prefix web run check
npm --prefix web run build
npm --prefix web run cross
npm --prefix web run smoke
```

The native binary is in `out/joy-pi-board` (`.exe` on Windows); the cross-built binary is `out/joy-pi-board-linux-arm64`. Run the native binary with `--listen 127.0.0.1:8081` for local development. Build commands compile frontend before Go embed; no Node runtime is required to run the binary. See [foundation commands and choices](docs/s00-foundation.md) for formatting, development, sizes, missing-dist and native race checks. Production lifecycle and installation are future work.

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
| [Sprint 00 work package](docs/sprint-00.md) | Foundation tasks and execution status |
| [Foundation choices](docs/s00-foundation.md) | Exact toolchains/dependencies, browsers, implemented shell and commands |
| [S00 evidence](docs/s00-evidence.md) | Actual local results, artifact hashes, CI gap and deferred gates |
| [Exact JSON integers ADR](docs/adr-001-exact-json-integers.md) | Adopted common convention inherited from Health |
| [Shared fixtures](testdata/README.md) | Authoritative numeric data and future contract taxonomy |
| [Decision register](docs/decisions.md) | Approved, replaced and deferred choices |
| [Original documentation review](docs/documentation-review.md) | Historical documentary-only checks; not S00/product acceptance |

Only decisions explicitly adopted in the register are approved. Remaining **Proposed P-xx** sections require their sprint's review. Upstream observations are **Verified H-xx**, pinned to Health commit `be7a0d824f62b94c842d8e5326110b1852c5a0bc`; reading those sources is not executing Health tests. S00 does not change that baseline or authorize S01–S04.

No system collection, database, history, administration, authentication or public Internet exposure belongs to v0.1.0. Health is a read-only reference: do not modify it or import its Go `internal` packages. See [AGENTS.md](AGENTS.md).
