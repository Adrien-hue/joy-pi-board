# Verified green S04 CI — 39313804438

Captured and independently inspected on 2026-10-20. Run date: 2026-10-19 UTC.
[GitHub run](https://github.com/Adrien-hue/joy-pi-board/actions/runs/39313804438),
attempt 1, exact source `5bc101e22a3093f0da3c0e6a7c50fee000dfa5da`.
All three jobs passed. This is software/archive evidence, not accepted release,
ARM64 execution, real-systemd lifecycle, Pi performance or operator validation.

- [ci.json](ci.json): verified run/job/step/artifact metadata and exact-source CI attestation.
- `foundation.log`, `race.log`, `browser.log`: complete decoded job logs, LF normalized; original artifact downloads also remain locally.
- [inspection.json](inspection.json): independently verified ZIP/package/binary checksums, archive member mapping, browser totals and encoded-body checks.
- `artifact/`: original clean-source **DEVELOPMENT / NOT ACCEPTED** provenance, asset manifest, notices and checksums. Its SHA256SUMS applies to the original artifact directory, not this partial copy.
- `package-one/`, `package-two/`, [reproducibility.json](reproducibility.json): actual successful package reports, build/Lintian logs, full package inventory, contents and same-input comparison. Both archives have the same version/bytes, not two upgrade versions.
- `generated-control/` and installed-file copies: actual extracted script/control/configuration/unit/override/manual bytes. Scripts carry `.txt` here for read-only review, not execution.
- `browser/`: full 27-test JSON report, actual host/engine inventory and raw encoded HTTP resource report. Captures/traces remain in the full browser ZIP; engine tests do not certify Safari/iOS or minimum branded browsers.
- [file-index.json](file-index.json): checksums for every retained file except the index itself.

Full ZIPs, selected `.deb` and ARM64 executable are retained under
`out/s04-ci-green/` (Git-ignored; preserve/back up that directory separately).
Their checksums are recorded in ci.json/inspection.json. GitHub artifact expiry
is 2026-10-26; retaining metadata alone does not retain the binary or browser
captures. Nothing was uploaded or published as part of this archive task.

Package SHA-256: `915f35240b9f9e6194bfbcee71fa2e4dcd98d475a606b8f635918ebe3fbb55dd`.
ARM64 executable SHA-256: `16bb1670c81638329b01ad7c83cc219c1957bb5ee1765b4a1091680f8dbbcb03`.
The package is a synthetic internal test fixture with missing real distribution
metadata. Both extracted binaries match the supplied artifact; no install ran.

Review finding **VM-F01**: no prerm, preinst stops on upgrade only, and postrm
does not stop on removal. The orphan/stop requirement remains open. See
[VM validation preparation](../../../s04-vm-validation.md) and the
[current ledger](../../../s04-evidence.md#2026-10-20--verified-green-ci-and-vm-preparation).
Historical failures and S00–S03 evidence remain unchanged.
