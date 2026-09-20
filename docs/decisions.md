# Decision register

Updated for the S00 user instruction, 2026-09-20. This instruction approves the foundation choices below, **not all previous documentation or every P-xx detail**. Product scope remains in [the specification](specification-v0.1.0.md).

| ID | Status and decision | Remaining review / sprint |
|---|---|---|
| P-01 | **APPROVED**: Board default `0.0.0.0:8081`; Health default remains `http://127.0.0.1:8080/v1/snapshot`. S00 shell implements `--listen`; it has no Health setting/client. | Full configuration/validation is S02. |
| P-02 | **PARTIAL / DEFERRED**: Go module `github.com/Adrien-hue/joy-pi-board`, React/TS/Vite/npm, standard Go HTTP/embed, `web/dist`, one binary are approved. Ordinary S00 shell routing/security choices are recorded in [foundation choices](s00-foundation.md). | Overview errors/method policy, full CLI/env, bounds, production lifecycle and compression remain proposals for S02/S04. No blanket ratification. |
| P-03 | **SPLIT / INITIAL PARSER PREFERENCE REPLACED**: see the three independent concerns below. | S00 records direction and browser targets; S01 implements and reviews parser/validator choices. |
| P-04 | **PROPOSED, DEFERRED S02**: coalesce only simultaneously active upstream calls. [Architecture](architecture-v0.1.0.md) contains a candidate algorithm, not an implementation. | Ratify cancellation, shared deadline, atomic publication/sequence ordering, orphan completion and no completed-flight reuse; no autonomous polling or immediate retry. |
| P-05 | **PROPOSED, DEFERRED S02/S03**: elapsed-age header; candidate calculation in [API](board-api-v0.1.0.md)/[UX](ux-v0.1.0.md). Nothing implemented in S00. | Finalize moment of calculation, transit accounting, missing/invalid headers, suspension and browser expiry after 30 s. Existing detailed text remains a proposal, not an approval. |
| P-06 | **BUDGETS ESTABLISHED; DETAILS PROPOSED**: product budgets remain mandatory. Shell size observations do not validate the final dashboard. | Ratify detailed units/accessibility conventions and performance/compression/physical protocols before corresponding S03/S04 gates. |

## P-03a — adopted numerical contract

[ADR-001](adr-001-exact-json-integers.md) adopts Health's exact unquoted decimal JSON integer convention for Joy Pi Home. No lossy floating intermediary, numeric-to-string substitution or conflation of null/zero/false. Schemas define domain/bounds; rounding is presentation only. This is **approved**, independent of a parsing library. Board's TypeScript uint64 representation is bigint at every magnitude, range 0 through 18446744073709551615; floating measures remain number according to their field schema.

## P-03b — Go validation and transport

S01 must validate required fields/nulls/types, exact uint64 bounds and invariants. Use uint64 or original json.Number tokens with checked conversion; UseNumber alone is not validation. Validated json.RawMessage transport as an object is compatible and remains the proposed realization. Presence wrappers and decoder organization are designs, not S00 code.

**Separate compatibility decisions still pending S01:** unknown-member rejection versus tolerance; open network-state strings; duplicate-key policy and defensive bounds. The earlier strict unknown-member proposal is not inherited from Health and is not a consequence of numeric exactness. Do not block S00 on it or silently mark it approved.

## P-03c — browser implementation direction

**Approved direction; actual implementation deferred:** evaluate a focused lossless JSON library first, or a native original-token facility supported by every retained browser target. Assess compatibility, license, maintenance, limits, malformed input, duplicate keys and measured bundle cost. No automatic library selection, unused runtime dependency or general-purpose homegrown parser in S00. The prior preference for an owned tokenizer is **superseded**, not retained as the default alternative.

Toolchain/dependency versions and browser targets selected for S00 are recorded in [foundation choices](s00-foundation.md); they are not claimed as additional historical user approvals. Review results and actual checks are in [S00 evidence](s00-evidence.md).
