# Joy Pi Board

Joy Pi Board is the second microservice of Joy Pi Home: a lightweight web overview of observations supplied by Joy Pi Health, for trusted devices on the local network.

**Board presents. Health observes.**

## Current state

This repository contains the **v0.1.0 specifications and completed Sprint 00 build shell**: a minimal React page served by one Go binary with embedded assets, exact tool/dependency locks, useful source/integration tests and a GitHub Actions workflow. S00 closed on 2026-09-20 after [CI run 35527458889](https://github.com/Adrien-hue/joy-pi-board/actions/runs/35527458889) passed for commit `b6e139c3378cd77fdcb4c5edc610977c08bbf72a`; see [S00 evidence](docs/s00-evidence.md) for its exact scope. **S01 is complete**, with local evidence and [CI run 35637699325](https://github.com/Adrien-hue/joy-pi-board/actions/runs/35637699325), attempt 1, verified on 2026-09-21 for `ccd1245a64fc372faa4ebc85eb000f78a01bd396`. **S02 is complete**, with local verification and [CI run 35908782107](https://github.com/Adrien-hue/joy-pi-board/actions/runs/35908782107), attempt 1, verified on 2026-09-23 for `a11f07c85ff4071a13ca3ae9227920e9f1c7dda4`. See [S02 evidence](docs/s02-evidence.md) for the confirmed scope and separate provenance. S03 is [PREPARED, NOT EXECUTED](docs/sprint-03.md); S04 is NOT STARTED. Board-owned Go/TypeScript validators, lossless parsing and shared contract tests are documented in [S01 scope](docs/sprint-01.md) and [S01 evidence](docs/s01-evidence.md). The demand-only Health client, overview API, volatile cache, safe configuration, shutdown and transition logging are implemented. The dashboard, packaging and installed production service are not implemented. Interactive visual checks, full browser/mobile acceptance, later functional and physical validation remain NOT RUN; there is no release candidate, release or deployment.

The product target is one page (`/`), one host and one provider. React/TypeScript/Vite assets are already embedded for the S00 shell. Production targets Raspberry Pi OS 64-bit / Trixie on Linux ARM64; Board acceptance uses a Raspberry Pi 3B+.

The shell serves its UI without Health. Only valid overview demand starts a provider attempt; startup, assets and idle time make no provider calls. The **approved** Board listener defaults to `0.0.0.0:8081`; the Health endpoint remains `http://127.0.0.1:8080/v1/snapshot`. The minimal page identifies the development foundation and displays no fake health values. `GET /api/v1/overview` returns exact JSON; unavailable Health remains HTTP 200 with an explicit provider reason.

## Build and run Board

Use Go **1.27.1**, Node.js **24.21.0**, npm **11.19.0**. From the repository root:

```text
node scripts/tasks.mjs tools
npm --prefix web ci
npm --prefix web run check
npm --prefix web run build
npm --prefix web run cross
npm --prefix web run smoke
npm --prefix web run fuzz
```

The native binary is in `out/joy-pi-board` (`.exe` on Windows); the cross-built binary is `out/joy-pi-board-linux-arm64`. Run the native binary with `--listen 127.0.0.1:8081` for local development. Build commands compile frontend before Go embed; no Node runtime is required to run the binary. See [foundation commands and choices](docs/s00-foundation.md) for formatting, development, sizes, missing-dist and native race checks. Use `--help` and `--version` without binding or contacting Health. `--health-url` / `JOY_PI_BOARD_HEALTH_URL` and `--listen` / `JOY_PI_BOARD_LISTEN` follow explicit flag > present env > default. SIGTERM/SIGINT stop admission and cancel work under one 5 s budget; installation/systemd remain S04.

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
| [Sprint 01 scope](docs/sprint-01.md) | Implemented in-memory contracts, parser evaluation and boundaries |
| [Sprint 02 work package](docs/sprint-02.md) | Completed overview/transport/resilience tasks and closure criteria |
| [Sprint 03 work package](docs/sprint-03.md) | Proposed production envelope consumer, browser lifecycle, dashboard tasks and browser/visual closure requirements |
| [S02 evidence](docs/s02-evidence.md) | Local source/results and verified exact-commit S02 CI evidence |
| [S01 evidence](docs/s01-evidence.md) | Local contract/build results, measured parsing cost and verified S01 CI closure |
| [Foundation choices](docs/s00-foundation.md) | Exact toolchains/dependencies, browsers, implemented shell and commands |
| [S00 evidence](docs/s00-evidence.md) | Local results and hashes, verified CI run/environment, S00 closure and deferred gates |
| [Exact JSON integers ADR](docs/adr-001-exact-json-integers.md) | Adopted common convention inherited from Health |
| [Shared fixtures](testdata/README.md) | Authoritative numeric data and future contract taxonomy |
| [Decision register](docs/decisions.md) | Approved, replaced and deferred choices |
| [Original documentation review](docs/documentation-review.md) | Historical documentary-only checks; not S00/product acceptance |

Only decisions explicitly adopted in the register are approved. Remaining **Proposed P-xx** sections require their sprint's review. Upstream observations are **Verified H-xx**, pinned to Health commit `be7a0d824f62b94c842d8e5326110b1852c5a0bc`; reading those sources is not executing Health tests. The current assignment authorizes S03 documentary preparation only; S03/S04 implementation remains unauthorized. The historical S01 evidence covers its implementation commit, not later documentary commits or the S02 implementation. S02 changes neither the upstream baseline nor deferred P-06/S04 decisions.

No system collection, database, history, administration, authentication or public Internet exposure belongs to v0.1.0. Health is a read-only reference: do not modify it or import its Go `internal` packages. See [AGENTS.md](AGENTS.md).
