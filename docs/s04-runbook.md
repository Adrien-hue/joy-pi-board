# S04 build and acceptance runbook

Status: software/archive CI verified at `5bc101e`; installed lifecycle and Pi acceptance incomplete. See [evidence](s04-evidence.md#2026-10-20--verified-green-ci-and-vm-preparation), [VM preparation](s04-vm-validation.md) and the [approved work package](sprint-04.md). Commands below are procedures, not evidence that they ran. No primary-host/Pi installation or service operation is authorized by the current assignment.

## Build and identity without a Pi

Use the existing exact Go 1.27.1, Node 24.21.0 and npm 11.19.0 on PATH. Do not run frontend construction concurrently with embed builds/tests. From repository root:

```sh
node scripts/tasks.mjs tools
npm --prefix web ci
npm --prefix web run check
npm --prefix web run fuzz
npm --prefix web run build
npm --prefix web run cross
npm --prefix web run smoke
npm --prefix web run missing-dist
npm --prefix web run sizes
node scripts/build-measurement.mjs
python3 scripts/python-tools.test.py
```

On supported Linux additionally run `npm --prefix web run race` and, after installing locked engines with `node web/node_modules/playwright/cli.js install --with-deps chromium firefox webkit`, `npm --prefix web run browser`. Record actual OS/architecture, tools, engine revisions and binary hashes. Native Windows is supplemental; current/minimum branded browsers and Apple devices are separate operator lanes.

Create a never-overwritten development identity:

```sh
node scripts/artifact.mjs development s04-dev-UNIQUE-ID
```

Use lowercase letters/digits/dot/hyphen in the actual ID. The command verifies tools, builds frontend and deterministic gzip, verifies decompression/manifest, then builds Linux ARM64 with CGO=0, GOARM64=v8.0, trimpath/buildvcs=false and explicit version/revision/source-epoch metadata. It checks all recorded inputs again after construction. The unique `out/artifacts/ID` contains binary, assets.json, THIRD-PARTY-NOTICES, provenance.json, SHA256SUMS and completion markers. COMPLETE is authoritative; INCOMPLETE preserves creation history. Never consume a directory without COMPLETE.

The provenance contains exact input hashes (including packaging), Git status, source/tool/flags, actual build time separately from deterministic linker time, and ELF/Go build information. The intended binary version 0.1.0 and private npm version are not acceptance/release status. Dirty inputs are always DEVELOPMENT / NOT ACCEPTED. Development artifacts do not satisfy clean-source candidate provenance.

After a separately authorized commit and successful exact-source CI, record a genuinely verified CI JSON with `sha`, `conclusion: "success"`, `verified: true` and the actual repository actions/run URL. Do not fabricate this attestation. With clean source:

```sh
node scripts/artifact.mjs candidate c01 path/to/verified-ci.json 0.1.0-1
```

Use an unused candidate ID and Debian revision. Local existing frozen revisions are rejected; also check the retained project-wide candidate ledger before using a new build directory. Reserve the ID/revision even for an incomplete attempt. Record the build environment, final manifests and checksums externally. A later changed source/asset/packaging input uses c02/0.1.0-2 or another unused revision, with affected evidence repeated. A comparison rebuild is separately named and never replaces selected bytes. Promotion copies accepted artifacts unchanged, without rebuilding, repackaging, stripping or editing embedded metadata.

## Locked packaging and archive checks

The original environment-lock status is historical input metadata: [CI 39313804438](evidence/s04/ci-39313804438/README.md) now retains the actual built image configuration ID, full transitive inventory and two identical test archives. No local Docker execution is inferred, and the image ID is not a registry manifest digest. Use the dated evidence record for verified closure, preserving the original input file and its checksum.

[environment-lock.json](../packaging/environment-lock.json) records the resolved Debian OCI base digests and dated snapshot. The base digest is verified registry metadata; **the derived image and full transitive package closure have not been executed/verified locally**. The Dockerfile uses the signed immutable snapshot with exact dpkg/dpkg-dev/debhelper checks. Do not fall back to a moving mirror or silently choose another version on failure. Retain apt error details; review any replacement.

In a disposable Linux build environment with a working container engine, the following performs builds and archive extraction only. Source is mounted read-only; the container has no writable host data mount:

```sh
docker build --platform linux/amd64 -f packaging/Dockerfile -t board-package-test packaging
docker image inspect board-package-test > out/packaging-image.json
docker create --name board-package-inspection --network none \
  --mount type=bind,source="$PWD",target=/src,readonly \
  board-package-test python3 /src/packaging/build.py \
  --artifact /src/out/artifacts/SELECTED-ID --output /work/package --test-only
docker start -a board-package-inspection
docker inspect --format '{{.State.ExitCode}}' board-package-inspection
docker cp board-package-inspection:/work out/packaging-inspection
docker rm board-package-inspection
```

Substitute the real retained identity. Verify exit zero and inspect package-evidence.json, build.log, contents.txt, packages.lock, lintian.log and **actual extracted DEBIAN scripts**. The complete installed package list includes transitive compressor, systemd/helper, adduser, binutils, lintian and Python versions; retain its checksum and the derived image ID/digest before claiming a locked tested packaging environment. Current unavailable execution does not establish those versions.

The test-only package is explicitly an INTERNAL TEST FIXTURE / NOT FOR DISTRIBUTION, with an example.invalid fixture identity, not a real Maintainer. It exercises mechanics only. No distributable package may be produced until the user supplies established Maintainer name/contact and Board copyright/distribution-license text. Third-party notices do not license Board. Supply an external metadata JSON containing `established: true`, `maintainer: "actual name <actual email>"` and `copyrightFile` pointing to the approved text; these descriptive examples must not be copied as real metadata. Omit `--test-only` and pass `--metadata` only once supplied. Archive evidence records the resulting checksum; retain the metadata input alongside it.

Packaging rejects mutated artifact hashes or changed packaging inputs. It copies the supplied binary, disables compile/strip/dwz/strip-nondeterminism/debug generation, extracts the output and byte-compares the executable. Debian archive compression is gzip-9, uniform, one thread, source epoch timestamps. Build two separately named output directories with identical inputs, then:

```sh
python3 packaging/compare.py /work/one /work/two
```

Archive inspection is not execution. ARM64 ELF inspection on amd64 does not prove ARM64 runtime behavior. A chroot/container without real systemd does not prove service lifecycle.

Review the two path-specific Lintian overrides shipped in `usr/share/lintian/overrides/joy-pi-board`: the approved static runtime and preservation of the frozen executable's symbols. `--show-overrides` makes these decisions visible in lintian.log and the CI console; all other errors still fail validation. The test-only copyright warning remains visible until established distribution metadata exists. Check the installed manual page under `usr/share/man/man1/joy-pi-board.1.gz`; packaging compares its decompressed bytes with the source manual. These policy exceptions do not waive lifecycle or Pi acceptance.

## Isolated service lifecycle

Start with [the prepared VM matrix and prerequisites](s04-vm-validation.md). VM-F01 found no removal stop path in the generated green-CI package; this must be resolved before A-06 acceptance. The two identical reproducibility archives are not an upgrade pair: supply two distinct, correctly ordered package versions with verified hashes. The default command below is inert; execute only after guest isolation and recovery are verified. Do not run Python with `-O`, because the current harness uses assertions.

Use a disposable **ARM64 VM** with real systemd PID 1, snapshot restoration available, no writable host/Pi data sharing, and no pre-existing Board account/package. Explicitly verify VM identity, mounts, OS/architecture and PID 1. Install only test prerequisites in that guest. A human must attest that no host block device or writable host data is exposed; filesystem checks alone cannot prove hypervisor isolation.

Create guest-only `/run/joy-pi-board-disposable.json` with its actual VM ID, `disposable: true` and `hostWritableMounts: false` after verification. Transfer two distinct package versions and hashes. The script refuses default execution, non-VM/non-systemd/non-ARM64 targets, existing package/accounts and known writable shared mounts:

```sh
python3 packaging/lifecycle.py
sudo python3 packaging/lifecycle.py --execute --vm-id VERIFIED-VM-ID \
  --previous /tmp/PREVIOUS.deb --package /tmp/SELECTED.deb \
  --output /tmp/s04-lifecycle.json
```

The first command is inert and returns NOT RUN. Only the second mutates the verified guest: install disabled/inactive, dedicated user, absent Health, configuration edits, explicit start/stop/enable/mask, upgrade/downgrade preserving choices, remove/purge/reinstall, executable equality and orphan checks. Failure retains evidence; restore the VM snapshot, not the primary host. Separately test incompatible user/group collisions, interrupted configure/upgrade and administrator drop-ins in fresh VM snapshots before A-06 acceptance. No script claims these additional cases passed automatically.

Inspect the generated helper code for enable/start/restart behavior and conffile handling. Verify `systemd-analyze verify`, DNS access under hardening, journald and the installed effective unit. Retaining the locked account after purge is intentional; administrator drop-ins/shared journal are not package-owned garbage.

## Pi prerequisites and later authorization

Obtain actual target identity/access, Pi 3B+ model, Trixie image/kernel and 64-bit ARM64 userland, storage/free space, power/cooling/ambient conditions, clocks/thermal state, existing services/ports/accounts and background load. Obtain Health version/commit/service/config and contract compatibility, Board existing state/drop-ins, credentials/privilege mechanism, trusted transfer path, outage window, backup/rollback package and restoration owner. None is assumed.

Obtain LAN/client hosts, current and minimum Chrome 111/Edge 111/Firefox 115/Safari 16.4, macOS and real iOS 16.4/current devices, reader/audio operator, native zoom, visibility/suspension/BFCache-capable clients. Agree a reversible WAN-isolation method preserving management connectivity. No Pi connection, WAN change or Health stop is authorized now.

The subsequent instruction must authorize the named target, artifact hashes, installation/configuration, Board account/unit creation and exact start/enable/stop/restart/upgrade/rollback/removal steps; separately name Health interruption and restoration, controlled provider, observers and WAN/client operations. Prepare a reviewed backup and rollback before execution.

## Future installed candidate procedure

All commands in this section require that later target-specific authorization. On the target, verify transferred SHA256SUMS using `sha256sum -c` before package extraction/install. Retain the old package, `/etc/default/joy-pi-board` and administrator drop-ins, plus `systemctl is-active`/`is-enabled`/`cat` output. Do not unmask or enable implicitly.

Install the selected package with `sudo dpkg -i ./SELECTED.deb`. Initial installation must leave Board inactive/disabled. Confirm installed `sha256sum /usr/bin/joy-pi-board` equals the standalone/extracted binary; record version, account and effective unit. Review the EnvironmentFile (systemd syntax, no shell expansion): existing listener/Health variables only, defaults 0.0.0.0:8081 and http://127.0.0.1:8080/v1/snapshot unchanged. Start explicitly with `sudo systemctl start joy-pi-board`; enable only if separately intended.

HTTP readiness is success for / and every manifest resource, not just unit active. Measure 20 fresh starts with Health present and absent, polling externally every 20 ms from monotonic spawn/start request. Separate standalone child startup from installed-unit startup timing. Stop trials measure SIGTERM-to-exit/`systemctl stop` with no second budget: 20 idle and 20 controlled in-flight timeout trials. Preserve trial logs, timestamps, exit status and orphan checks. Every readiness <=2 s and shutdown <=5 s; systemd's 6 s kill safeguard does not make 6 s acceptable.

For upgrade/rollback, explicitly stop Board, confirm exit, install the retained selected version using `dpkg --force-confold -i` only after reviewing config compatibility, verify hashes, then start according to the recorded state. An enabled unit may remain enabled but must not auto-start; disabled/masked states remain respected. Verify `/proc/MAINPID/exe` checksum externally after restart as well as installed path; old running bytes cannot count as the new package. Test remove/purge in the isolated VM before considering such operations on Pi.

## External measurements

Freeze these tools and inputs with the campaign; Node/Python are measurement-host tools, not Board production dependencies. Run diagnostics **off** for ordinary CPU/RSS/journal gates. One-core CPU denominator is process user+system CPU time / elapsed monotonic time, never divided by core count. Read only the identified Board process from the external collector:

```sh
python3 scripts/process-samples.py --pid VERIFIED-PID --seconds 720 --warmup 120 --output idle.raw.jsonl
python3 scripts/process-report.py idle.raw.jsonl idle.summary.json
node scripts/measure.mjs load http://TARGET:8081 one.raw.json one
node scripts/measure.mjs load http://TARGET:8081 ten.raw.json ten
node scripts/measure.mjs load http://TARGET:8081 sync.raw.json sync
node scripts/measure.mjs inventory http://TARGET:8081 resources.json out/artifacts/SELECTED-ID/assets.json
```

Use new output names each time. The load tool uses 2 min warmup + 120 measured five-second slots per client; ten staggers at 0.5 s or synchronizes. It retains planned/actual starts, late/skips, statuses, bytes, outcomes and durations including body validation, grouped provider states/reasons and per-client latency. Do not pass an incomplete nominal population by excluding failures. Raw byte inventories retain base64-encoded wire bodies, headers and hashes with gzip decompression comparisons; overview remains identity. A cold real browser/HAR must additionally prove the complete initial dependency set and LAN-only operation.

For >=200 five-second overhead/refusal samples, use `node scripts/measure.mjs load http://TARGET:8081 overhead.raw.json one 200 120000` or `node scripts/measure.mjs load http://TARGET:8081 refusal.raw.json one 200 100000` (20 refusal warmup attempts). Retain raw output and exact invocation. Verify an immediate refused listener separately; a packet drop/stalled body is a true timeout, not refusal. For timeout modes use `one 30 0` for each controlled mode and retain Health's global 1 s cancellation behavior independently of client latency. The external client deadline is 2 s and does not alter the Health deadline.

Temporarily enable `--diagnostic-timing` in an authorized ExecStart override of the **same binary**. Use a dedicated one-client process/run, no unrelated overview requests, 24 warmup requests then >=200 measured requests at 5 s. Stop normally within the original budget to emit the summary. Export its entire process journal, not a truncated slice:

```sh
journalctl _SYSTEMD_INVOCATION_ID=RECORDED-INVOCATION-ID --output=json > timing.journal.jsonl
node scripts/measure.mjs timing timing.journal.jsonl timing.summary.json 24
```

Record the actual invocation ID with systemctl show before stopping, and scope the journal to exactly that process/invocation; unit-wide systemd messages are not Board JSON. Lifecycle records are ignored by the timing reducer; malformed JSON, multiple/missing summary, lost/truncated rows or missing sequence fails it. It subtracts each eligible request's associated pre-validation network interval from handler duration, then computes nearest-rank p95. Shared/cancelled/incomplete/unsuccessful rows are excluded and counted; a whole run with losses is invalid. At 4,096 records the diagnostic truncates; split campaigns into separately identified processes. Enqueue/reporting cost follows the timed handler; calibrate paired on/off runs with identical binary hash. Remove the override and restart explicitly after diagnostics.

Retain raw data and compare unchanged budgets: RSS max <=50 MiB; idle CPU <=1% one core; one-client mean <=2%; ten-client resource trend with **no new CPU ceiling**; overhead p95 <=50 ms; nominal p95 <=500 ms; immediate refusal p95 <=250 ms; readiness <=2 s; shutdown <=5 s; JS+CSS gzip <=250 KiB and all initial resources <=500 KiB. Report N/min/median/p95/max, errors and late/skipped samples, not a success-only percentile. Never label controlled local tooling results physical PASS.

Finally run real Health failure/stale/expiration/recovery, independent browser expiry with Board stopped, empty cache after restart, cold LAN load with WAN disconnected, journal leak/poll-volume review and the [transferred client matrix](sprint-04.md#transferred-s03-operator-checks). Recheck package/file/running executable hashes and restore services/network/configuration. A/B/C and transferred gates determine a dated GO/NO-GO decision. Publication requires separate authorization and must promote the accepted bytes unchanged.
