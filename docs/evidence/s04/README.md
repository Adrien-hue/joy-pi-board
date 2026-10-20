# S04 retained local evidence

The later [green CI 39313804438 archive](ci-39313804438/README.md), captured on 2026-10-20, is separate from the historical local evidence below. It records software/archive PASS and prepared, unexecuted VM lifecycle work, including VM-F01.

These records identify dirty-tree DEVELOPMENT / NOT ACCEPTED software checks on 2026-10-18. See [the evidence ledger](../../s04-evidence.md) for scope, incidents and missing validation.

- `artifacts.json`: selected development build and separate comparison, exact binary/assets/input hashes; no package or ARM64 runtime claim.
- `provenance.json`, `assets.json`, `SHA256SUMS`: copies from the selected development artifact. SHA256SUMS applies to that artifact directory; its binary/notices are retained under `out/artifacts/s04-dev-20261018-c`, not checked into this evidence directory.
- `check-final.log`, `fuzz.log`, `smoke.log`, `missing-dist.log`, tool logs: actual source/native/controlled-tool checks, not Pi performance.
- `browser-summary.json` and `browser.log`: all 27 exploratory Windows engine cases, including five failures; raw report/screenshots/traces remain in `out/browser-report.json` and `out/browser-results`. No manual visual PASS is inferred.
- `sizes-summary.json`: actual encoded HTTP resource totals; full wire-body base64/header/hash inventory remains in `out/dashboard-sizes.json`, linked by checksum.
- `packaging-environment.json`: resolved base digests and proposed snapshot installation, explicitly unverified package closure. No generated package/script evidence is fabricated.
- `file-index.json`: checksums of the retained machine records and logs at capture. This README is explanatory, outside that index.

Later CI follow-ups are separate from the original file index: `ci-39308202979.json` records the first packaging failure with missing detail; [ci-39310053644.json](ci-39310053644.json), its retained build log and package inventory record the subsequent diagnosed dh_makeshlibs failure. They do not claim a successful package build or lifecycle execution. See the dated ledger entries for corrected-source validation status.

[ci-39311632646.json](ci-39311632646.json) and its retained Lintian log/package report record the first constructed and byte-verified test archive, rejected by Lintian. They are not successful full-package validation or acceptance evidence.

Logs copied from PowerShell are normalized to UTF-8; their content/results are preserved. Old S00–S03 evidence files and captures are untouched. Future accepted candidate evidence must use its own identity and retain raw artifacts; these local records cannot be promoted into Pi, systemd or official-release proof.
