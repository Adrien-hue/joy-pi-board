# Decision register

Updated for S02 documentary preparation, 2026-09-22. S00/S01 approvals and closure remain historical facts. This preparation recommends refinements but approves none of P-02/P-04/P-05/P-06's deferred details. Explicit S00 foundation and S01 consumer approvals are recorded below, **not all previous documentation or every P-xx detail**. Product scope remains in [the specification](specification-v0.1.0.md).

| ID | Status and decision | Remaining review / sprint |
|---|---|---|
| P-01 | **APPROVED**: Board default `0.0.0.0:8081`; Health default remains `http://127.0.0.1:8080/v1/snapshot`. S00 shell implements `--listen`; it has no Health setting/client. | Full configuration/validation is S02. |
| P-02 | **PARTIAL / DEFERRED**: Go module `github.com/Adrien-hue/joy-pi-board`, React/TS/Vite/npm, standard Go HTTP/embed, `web/dist`, one binary are approved. Ordinary S00 shell routing/security choices are recorded in [foundation choices](s00-foundation.md). | Concrete S02 recommendations: [transport/routing/serialization](board-api-v0.1.0.md), [configuration/lifecycle](architecture-v0.1.0.md#startup-configuration-and-operations). S02 identity/no-store; compression and package hardening remain S04 proposals. Ratify before S02-02/03/05/06; no blanket approval. |
| P-03 | **APPROVED S01 CONSUMER POLICY / S01 COMPLETE**: see the three independent concerns below. | Local and hosted S01 contract checks PASS; [CI evidence](s01-evidence.md#verified-hosted-ci). Browser-floor execution NOT RUN. Numeric convention remains distinct from compatibility and library selection. |
| P-04 | **PROPOSED, PREPARED FOR S02**: share only an active nonterminal flight, with lifecycle context, original deadline and atomic immutable view publication. [Architecture](architecture-v0.1.0.md#concurrent-demand-and-ordering) owns the recommended algorithm; not implemented. | Ratify orphan continuation, terminal events, delayed-response coherence and bounded worker/admission policy. S02-T06/T07; no autonomous polling, retry or completed-flight reuse. |
| P-05 | **PROPOSED, PREPARED FOR S02/S03**: X-Joy-Pi-Snapshot-Age-Ms with paired monotonic origin and a single response-generation cut. [API](board-api-v0.1.0.md#browser-safe-elapsed-age-header) owns server protocol; [UX](ux-v0.1.0.md) owns future browser behavior. Not implemented. | Ratify ceil/saturation/header presence, expired/delayed response policy and conservative round-trip/suspension handling. Server S02-T05/T08; browser B-10 in S03. No JSON-envelope change or approval implied. |
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

## S02 recommendations awaiting ratification

The [S02 work package](sprint-02.md) is **PREPARED, NOT EXECUTED**. The following summary records recommendations, rationale, consequences and test ownership; detailed rules have one authority location rather than competing copies.

| Proposal | Recommendation and rationale | Observable consequence / planned tests |
|---|---|---|
| P-02 configuration | Flags over present env over defaults; strict numeric listener and HTTP provider URL; help/version without bind/probe. Predictable startup without latent service coupling. | Safe exit 2 on config errors, exit 1 on bind/runtime failures; approved 8081/8080 defaults unchanged. S02-T01/T11. |
| P-02 transport | Dedicated no-proxy/no-decompression/no-reuse HTTP/1 client, no redirects, one selected address/dial, 1 s whole-flight deadline and 64 KiB counted body. Avoid hidden replay and unbounded input. | Next attempt only on next demand; DNS address fallback is intentionally unavailable; exact five-reason mapping and Health failure HTTP 200. S02-T02–T04. |
| P-02 response/service | Exact router and safe JSON errors, bounded admission/connections/workers; identity/no-store delivery in S02; separate 80 KiB/depth-34 envelope and non-HTML-escaped raw serialization. Preserve S01 boundary payloads and existing shell. | Overload is Board 503/transport close, not provider failure; wrong method/query/body never contacts Health. Shutdown ≤5 s; bounded transition-only logs. S02-T08–T12. Compression/packaging remain S04. |
| P-04 sharing | A lifecycle-owned active flight survives individual/all waiter departures only until ordinary completion or original deadline. Atomic terminal publication clears active; old views remain immutable and bounded. | Cancelled first caller cannot cancel others; new demand after terminal starts new work; late worker cannot overwrite state. S02-T06/T07. |
| P-05 age | Sample monotonic age at response finalization with snapshot selection; integer ceil header, present whenever success metadata exists; fail closed for unverifiable time. Browser later adds complete RTT and hides on suspension. | Failed fallback at >30 s becomes none; old successful view cannot be sent as indefinitely current. No timestamp subtraction across Pi/browser. S02-T05/T08; B-10 NOT RUN until S03. |

Implementation must record explicit ratification or replacements before dependent code/gates. P-06's existing row, budgets and deferred presentation/measurement review are unchanged. S02 preparation resolves editorial gaps, not product approvals.
