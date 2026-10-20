# S04 disposable ARM64 VM validation package

Prepared 2026-10-20. Status: **PREPARED; EXECUTION NOT RUN; A-06 OPEN**.
The user's request to archive green CI and prepare VM validation permits this
read-only review and planning. Existing authorization permits package/service
execution only in a verified disposable isolated guest. No VM identity, access,
snapshot or isolation attestation has been supplied. No guest was contacted or
created; no installation or service operation took place.

## Starting evidence and review finding

[CI 39313804438](https://github.com/Adrien-hue/joy-pi-board/actions/runs/39313804438)
passed on `5bc101e22a3093f0da3c0e6a7c50fee000dfa5da`. The
[retained record](evidence/s04/ci-39313804438/README.md) verifies two identical
Debian archives and exact input/extracted executable equality. These are
DEVELOPMENT / NOT ACCEPTED, internal test fixtures. The build container had no
real systemd lifecycle test and did not execute the ARM64 binary.

**VM-F01: removal-stop implementation gap, found by static review.** The actual
[preinst](evidence/s04/ci-39313804438/generated-control/preinst.txt) stops Board
only on upgrade. There is no prerm in the package. The actual
[postrm](evidence/s04/ci-39313804438/generated-control/postrm.txt) reloads systemd
on remove and cleans helper state on purge, but does not stop the process.
Therefore the source assertion that removal stops Board is not established by
the generated package. A running process can outlive removal of its executable;
an orphan is a risk inferred from these scripts, not an observed VM result.

Keep the existing harness's active-service removal and orphan assertions. In an
isolated snapshot, reproduce VM-06 and retain the result; if confirmed, correct
the package stop path, run exact-source CI on a new revision, and repeat the
lifecycle matrix using newly identified artifacts. Do not pre-stop Board to
make this case pass. VM-F01 blocks A-06 acceptance and candidate promotion until
resolved; it does not invalidate the recorded source/archive CI results.

Review the helper's generated unmask/was-enabled/update-state code against real
administrator masks and enablement choices. A helper call named unmask alone is
not evidence that an administrator mask is lost; VM-05 must establish behavior.
The [Debian helper reference](https://manpages.debian.org/trixie/init-system-helpers/deb-systemd-helper.1p.en.html)
describes its separate state bookkeeping. Preserve logs and symlinks before and
after each transition rather than assuming helper semantics prove the outcome.

## Required environment and inputs

| Input | Required record | Current status |
|---|---|---|
| Guest | Named disposable ARM64 VM; Trixie OS image/source/checksum; kernel; real systemd PID 1; VM provider/version and virtualization type | UNKNOWN |
| Isolation | Operator attestation: no host block devices, writable host/Pi shares or production credentials; explicit guest disks and mounts; isolated network | UNKNOWN |
| Recovery | Named powered-off baseline snapshot; owner and tested restore procedure | UNKNOWN |
| Access | Guest-only console/SSH route and privilege mechanism; no Pi connection | UNKNOWN |
| Prerequisites | Python 3, dpkg, adduser, systemd, procps, util-linux and standard shell tools; guest tool/version inventory | NOT VERIFIED |
| Baseline | No Board package/account/group/drop-in, no real Health service, required test ports unused | NOT VERIFIED |
| Packages | Two distinct versions, ordered by dpkg; immutable hashes, manifests and known source; inspect generated scripts first | NOT PREPARED |
| Distribution | Established Maintainer name/contact and copyright/license | MISSING; test fixtures only |

The two reproducibility outputs from green CI are identical versions and bytes:
**they are not an upgrade/rollback pair**. Prepare two separately identified
test versions in an isolated builder using the commands in the
[runbook](s04-runbook.md#build-and-identity-without-a-pi); do not edit or repack an
existing archive, and never relabel a development artifact as a frozen candidate.
Synthetic fixture metadata is permissible for this mechanics-only campaign.
Missing distribution metadata does not prevent independent test preparation.

## Preflight and execution procedure

1. Record guest identity, boot ID, OS, architecture, systemd version/PID 1,
   virtualization, mount inventory and snapshot identity. Inspect them before
   any mutating command. ARM64 emulation must be labelled as such; neither
   emulation nor native ARM64 VM execution satisfies Pi Class C.
2. Transfer only the named test packages, their checksum/provenance records and
   the reviewed tooling. Verify package SHA-256 before extraction. Use
   `dpkg-deb -f PACKAGE Package Architecture Version`, `dpkg-deb -e` and
   `dpkg-deb -x` in guest temporary directories for read-only inspection. Confirm
   package identity/ARM64, old version less than new using
   `dpkg --compare-versions OLD lt NEW`, and both extracted binary hashes.
3. Confirm absence of account, group, package remnants, overrides and listeners;
   verify the operator's isolation attestation externally to the guest as well.
   The harness alone cannot detect every hypervisor passthrough or unsafe mount.
4. Only after those checks, create the guest-local marker
   `/run/joy-pi-board-disposable.json` with the actual VM ID, `disposable: true`
   and `hostWritableMounts: false`. The marker is an attestation, not proof of
   isolation. Do not fabricate its values to bypass the guard.
5. Use the delivered commands below in the guest, substituting actual approved
   paths/ID. `--execute` installs, changes configuration, starts/stops/enables/
   disables/masks, upgrades/downgrades, removes and purges Board in that guest.
   Choose a new evidence path and run Python normally, never with `-O` (the
   current harness uses assertions). A nonzero result is FAIL/incomplete.

```sh
python3 packaging/lifecycle.py
sudo python3 packaging/lifecycle.py --execute --vm-id VERIFIED-VM-ID \
  --previous /tmp/VERIFIED-OLDER.deb --package /tmp/VERIFIED-NEWER.deb \
  --output /tmp/UNIQUE-lifecycle.json
```

6. Export the full JSON (including failed assertions), process/unit journals,
   hashes, unit/drop-in files and status/link snapshots before rollback. Restore
   the VM snapshot; never use host cleanup as a substitute. An interrupted run
   is incomplete, not PASS. Supplement harness output with explicit per-case
   results because its aggregate PASS covers only the listed automated cases.

## Matrix and evidence oracles

All rows are **NOT RUN**. S04-04/A-06 is the primary gate; identity work also
supports S04-05/A-01. Runtime cases exercise software behavior relevant to A-05,
B-05/B-06/B-09 without closing physical C-01/C-04/C-06.

| ID | Procedure | Required evidence / oracle |
|---|---|---|
| VM-01 | Fresh install, no explicit start/enable | Dedicated locked system account/group; expected modes/owners; disabled and inactive; valid unit; no Health dependency; exact installed binary checksum |
| VM-02 | Explicit start with controlled Health refusal | UI/readiness and unavailable overview; PID/User/Group; effective hardening, journal; exact running executable hash; VM timings labelled non-Pi |
| VM-03 | Administrator conffile edit; enabled then explicitly stopped upgrade | New exact binary; conffile retained; enabled state retained; inactive after upgrade; explicit restart uses new bytes |
| VM-04 | Fresh snapshot, disabled-state upgrade and rollback | Disabled stays disabled and no implicit start; old/new package versions and installed hashes; compatible preserved configuration |
| VM-05 | Administrator mask then downgrade; separately test disabled/masked upgrades | Mask symlink and state survive; no implicit unmask/start; conffile unchanged; explicit operator unmask only |
| VM-06 | Remove while Board is running, then purge | **VM-F01 focus**: service/process actually exits; no listening socket/orphan; remove preserves conffile/account; purge removes package files/state but retains account |
| VM-07 | Reinstall after purge | Dedicated account reused safely; no enable/start; exact bytes; second purge clean; no orphan |
| VM-08 | Fresh snapshots with incompatible user/group, supplementary groups, login/home/password state | Configuration fails with clear diagnosis; no start; unrelated account unchanged; recovery via snapshot |
| VM-09 | Interrupted configure/upgrade and administrator drop-ins | Preserve raw failure and state; safe recovery; do not erase admin data or silently change enable/mask choices |
| VM-10 | Explicit stop/restart and hardening review | Graceful exit, journal, effective unit; DNS/loopback behavior under restrictions; stopped listener; unchanged checksums before/after |

The harness covers part of VM-01/02/03/05/06/07; VM-04/08/09/10 and detailed
permissions/running-executable evidence require additional operator procedures.
Do not mark those supplemental rows PASS from the aggregate harness result.
An accepted candidate must repeat applicable checks on its own exact package;
development VM evidence cannot be promoted by changing a label.

## Exit and handoff

Archive a dated VM identity/isolation record, exact package and binary hashes,
tool inventory, per-case statuses, raw logs and recovery outcome. Resolve VM-F01
and any further finding before A-06 completion. Then prepare a unique clean-source
candidate with reviewed distribution metadata, exact-source CI and immutable
artifacts. Pi installation, Health interruption, network changes and final
operator/browser acceptance require their separately identified campaign.
No release GO or publication follows from VM success alone.
