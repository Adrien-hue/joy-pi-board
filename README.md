# Joy Pi Board

Joy Pi Board is the second microservice of Joy Pi Home: a lightweight web overview of observations supplied by Joy Pi Health, for trusted devices on the local network.

**Board presents. Health observes.**

## Current state

This repository contains **v0.1.0 software through the S04 implementation lane**: the embedded dashboard and exact overview consumer, deterministic static gzip, optional bounded request timing, Debian/systemd recipes and external acceptance tools. [S04 evidence](docs/s04-evidence.md) separates executed local checks from pending Linux packaging/lifecycle, remote CI and physical acceptance. No package is installed and no release is accepted.

S00, S01 and S02 are closed with their [S00](docs/s00-evidence.md), [S01](docs/s01-evidence.md) and [S02](docs/s02-evidence.md) exact-commit evidence. S03 is implemented and locally/CI verified within its executed scope: corrected run **36269739581**, attempt 1, for **24d527a301d66b49d86ad11fbd75393254f02094**. Its [remaining operator checks](docs/s03-evidence.md#remaining-acceptance-checklist) are deferred to S04 acceptance by explicit user decision; they remain NOT RUN and S03 is not fully validated or closed. That deferral permits independent S04 software work. Pi access does not supply Safari/iOS, historical browsers or assistive clients.

The [2026-09-27 Chrome retry](docs/s03-evidence.md#chrome-retry-with-restored-extension-access) adds real-browser exact values, degraded-state/recovery/expiry, selected keyboard checks and reviewed responsive captures. Required native zoom, lifecycle, screen-reader and other target evidence remain open.

The product target is one page (`/`), one host and one provider. React/TypeScript/Vite dashboard assets are embedded in the binary. Production targets Raspberry Pi OS 64-bit / Trixie on Linux ARM64; Board acceptance uses a Raspberry Pi 3B+.

Board serves its UI without Health. Only valid overview demand starts a provider attempt; startup, assets and idle time make no provider calls. The **approved** Board listener defaults to `0.0.0.0:8081`; the Health endpoint remains `http://127.0.0.1:8080/v1/snapshot`. The dashboard presents only validated Health observations, with bounded exact parsing, five-second refresh slots and independent local expiry. No fake observations enter production assets. `GET /api/v1/overview` returns exact JSON; unavailable Health remains HTTP 200 with an explicit provider reason.

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

The native binary is in `out/joy-pi-board` (`.exe` on Windows); the cross-built binary is `out/joy-pi-board-linux-arm64`. Run the native binary with `--listen 127.0.0.1:8081` for local development. Build commands compile frontend before Go embed; no Node runtime is required to run the binary. See [foundation commands and choices](docs/s00-foundation.md) for formatting, development, sizes, missing-dist and native race checks. Use `--help` and `--version` without binding or contacting Health. `--health-url` / `JOY_PI_BOARD_HEALTH_URL` and `--listen` / `JOY_PI_BOARD_LISTEN` follow explicit flag > present env > default. SIGTERM/SIGINT stop admission and cancel work under one 5 s budget; the delivered systemd recipe still requires isolated lifecycle and Pi acceptance. `--diagnostic-timing` is a temporary no-value opt-in flag, off by default, with no environment equivalent or public endpoint; it cannot be combined with help/version.

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
| [Sprint 03 work package](docs/sprint-03.md) | Production envelope consumer, browser lifecycle, dashboard tasks and browser/visual closure requirements |
| [Sprint 04 work package](docs/sprint-04.md) | Approved software scope and separate Pi/client/release gates |
| [S04 evidence](docs/s04-evidence.md) | Tested dirty-tree inputs, development artifacts, results and gaps |
| [S04 runbook](docs/s04-runbook.md) | Actual build/package/measurement commands and future authorized target campaign |
| [S03 evidence](docs/s03-evidence.md) | Historical local results, verified corrected CI, operator observations and outstanding closure gates |
| [S02 evidence](docs/s02-evidence.md) | Local source/results and verified exact-commit S02 CI evidence |
| [S01 evidence](docs/s01-evidence.md) | Local contract/build results, measured parsing cost and verified S01 CI closure |
| [Foundation choices](docs/s00-foundation.md) | Exact toolchains/dependencies, browsers, implemented shell and commands |
| [S00 evidence](docs/s00-evidence.md) | Local results and hashes, verified CI run/environment, S00 closure and deferred gates |
| [Exact JSON integers ADR](docs/adr-001-exact-json-integers.md) | Adopted common convention inherited from Health |
| [Shared fixtures](testdata/README.md) | Authoritative numeric data and future contract taxonomy |
| [Decision register](docs/decisions.md) | Approved, replaced and deferred choices |
| [Original documentation review](docs/documentation-review.md) | Historical documentary-only checks; not S00/product acceptance |

The [decision register](docs/decisions.md#d-2026-10-18--s04-software-implementation) records the user's approval of S04-R01–R04 and S03 operator deferral. The current instruction permits S04-01–S04-05 software work and isolated guest tests only; no primary-host/Pi installation, service operation, commit, push, tag, publication or deployment. Upstream Health remains pinned to `be7a0d824f62b94c842d8e5326110b1852c5a0bc`. All required unexecuted acceptance remains NOT RUN.

No system collection, database, history, administration, authentication or public Internet exposure belongs to v0.1.0. Health is a read-only reference: do not modify it or import its Go `internal` packages. See [AGENTS.md](AGENTS.md).

## S03 browser checks

After the pinned installation and native build, install the reviewed Playwright engines with `node web/node_modules/playwright/cli.js install --with-deps chromium firefox webkit` on supported Linux, then run `npm --prefix web run browser`. This verifies the engine revision lock, builds a separate development-only StrictMode harness outside production assets, and exercises the embedded production binary with controlled Health. `npm --prefix web run sizes` measures the real dashboard and three first-overview scenarios, counting actual encoded static gzip bodies and the first identity overview. Historical S03 offline gzip measurements remain separate. See [operator steps and remaining targets](docs/s03-evidence.md#remaining-browser-and-operator-procedure); automated WebKit is not Safari/iOS certification.
