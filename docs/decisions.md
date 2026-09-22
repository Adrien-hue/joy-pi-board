# Decision register

Updated for the S01 user instruction and documentary CI closure, 2026-09-21. This closure changes no technical decision and ratifies none of P-02/P-04/P-05/P-06's deferred details. Explicit S00 foundation and S01 consumer approvals are recorded below, **not all previous documentation or every P-xx detail**. Product scope remains in [the specification](specification-v0.1.0.md).

| ID | Status and decision | Remaining review / sprint |
|---|---|---|
| P-01 | **APPROVED**: Board default `0.0.0.0:8081`; Health default remains `http://127.0.0.1:8080/v1/snapshot`. S00 shell implements `--listen`; it has no Health setting/client. | Full configuration/validation is S02. |
| P-02 | **PARTIAL / DEFERRED**: Go module `github.com/Adrien-hue/joy-pi-board`, React/TS/Vite/npm, standard Go HTTP/embed, `web/dist`, one binary are approved. Ordinary S00 shell routing/security choices are recorded in [foundation choices](s00-foundation.md). | Overview errors/method policy, full CLI/env, bounds, production lifecycle and compression remain proposals for S02/S04. No blanket ratification. |
| P-03 | **APPROVED S01 CONSUMER POLICY / S01 COMPLETE**: see the three independent concerns below. | Local and hosted S01 contract checks PASS; [CI evidence](s01-evidence.md#verified-hosted-ci). Browser-floor execution NOT RUN. Numeric convention remains distinct from compatibility and library selection. |
| P-04 | **PROPOSED, DEFERRED S02**: coalesce only simultaneously active upstream calls. [Architecture](architecture-v0.1.0.md) contains a candidate algorithm, not an implementation. | Ratify cancellation, shared deadline, atomic publication/sequence ordering, orphan completion and no completed-flight reuse; no autonomous polling or immediate retry. |
| P-05 | **PROPOSED, DEFERRED S02/S03**: elapsed-age header; candidate calculation in [API](board-api-v0.1.0.md)/[UX](ux-v0.1.0.md). Nothing implemented in S00. | Finalize moment of calculation, transit accounting, missing/invalid headers, suspension and browser expiry after 30 s. Existing detailed text remains a proposal, not an approval. |
| P-06 | **BUDGETS ESTABLISHED; DETAILS PROPOSED**: product budgets remain mandatory. Shell size observations do not validate the final dashboard. | Ratify detailed units/accessibility conventions and performance/compression/physical protocols before corresponding S03/S04 gates. |

## P-03a — adopted numerical contract

[ADR-001](adr-001-exact-json-integers.md) adopts Health's exact unquoted decimal JSON integer convention for Joy Pi Home. No lossy floating intermediary, numeric-to-string substitution or conflation of null/zero/false. Schemas define domain/bounds; rounding is presentation only. This is **approved**, independent of a parsing library. Board's TypeScript uint64 representation is bigint at every magnitude, range 0 through 18446744073709551615; floating measures remain number according to their field schema.

## P-03b — approved consumer validation and Go realization

The S01 instruction approves tolerant unknown members for schema 1.0, required exact-case known keys, duplicate rejection using decoded names (including unknown subtrees), open non-empty network state strings, uint64 decimal-only nonnegative tokens and checked bounds. Unknown members are preserved in validated raw bytes but neither satisfy required fields nor count as useful observations. Unknown versions, issue codes and nonconforming paths are never accepted by this tolerance.

Snapshot limits are 65,536 bytes inclusive and 32 containers (root object = 1). Reject trailing content, malformed UTF-8 and unpaired Unicode surrogates. The exact convention and classification order are normative in [Health integration](health-integration-v1.0.md#validator-policy-and-failure-boundary). These are Board consumption choices, not changes to Health or consequences of numerical exactness. The former strict unknown-member proposal is **replaced**.

Implemented ordinary Go choice: a standard-library Decoder.Token/UseNumber walk with per-object maps checks exact-key presence and duplicates before schema dispatch. Required absence and explicit null are distinguished before creating the typed projection. Checked uint64 conversion never uses float64. A private validated byte copy implements json.Marshaler, returning fresh bytes as a JSON object; Bytes and Snapshot also return independent data. It is not a cache or HTTP envelope.

## P-03c — selected browser parser

The approved evaluation direction is implemented with **lossless-json 4.3.1** (MIT), one necessary dependency with no production transitives. [S01 evaluation](sprint-01.md#parser-decision-and-primary-source-evaluation) records native incompatibility at the retained floors, library limitations, narrow defensive guards, original numeric tokens and actual bundle measurement. Library selection and adapter organization are ordinary implementation choices; the user approved the required behavior, not a specific upstream library.

Every uint64 field becomes bigint after runtime validation, every floating measure becomes number by schema, and decimalUInt64 returns exact decimal text. No JSON numeric rewriting, rounded-Number conversion, BigInt.prototype mutation or general-purpose owned parser is introduced. The S00 homegrown-parser preference remains superseded. P-02/P-04/P-05 and remaining P-06 choices are not ratified by S01.

Toolchain/dependency versions and browser targets selected for S00 are recorded in [foundation choices](s00-foundation.md); they are not claimed as additional historical user approvals. Foundation checks are in [S00 evidence](s00-evidence.md); contract results and limits are in [S01 evidence](s01-evidence.md).
