# Joy Pi Board package operation

Private artifacts are not accepted releases. Record the exact package/binary hashes and candidate status before installation. Package construction consumes a supplied frozen binary and never compiles or strips it.

The package uses /usr/bin/joy-pi-board, /usr/lib/systemd/system/joy-pi-board.service and the dpkg conffile /etc/default/joy-pi-board. The dedicated system account has no login, home or device groups. There is no Health service/package dependency or production Node runtime.

Installation does not enable or start Board. After verifying the package and configuration, an authorized operator explicitly starts it, then optionally enables it. HTTP readiness, not systemd active status, is the readiness oracle. Defaults remain 0.0.0.0:8081 and http://127.0.0.1:8080/v1/snapshot. Flag > present environment > default; an empty value is invalid. EnvironmentFile uses systemd syntax, not shell expansion. Choose an unprivileged listener port.

For upgrade or downgrade: record running/enabled/masked state; back up conffile and administrator drop-ins; explicitly stop Board; verify exit; install the selected checksum-verified package using normal dpkg/apt tooling; preserve administrator conffile choices; verify the installed executable; explicitly start only if intended. Never unmask automatically. Replacing an executable does not replace a running process. Verify both installed and running executable hashes after restart.

Removal stops Board and retains the conffile/account. Purge removes package-managed configuration and enablement bookkeeping, retaining the account to avoid UID reuse. Administrator-owned drop-ins and shared journal are not recursively erased. Rollback restores a retained old package and a compatible reviewed configuration; no persistent snapshot migration exists. Re-check readiness, hashes and Health-absent operation after restoration.

Logs use journald; startup/shutdown and significant provider transitions only by default. --diagnostic-timing is an explicit no-value flag, incompatible with help/version; no environment variable exists. Use it temporarily via an authorized ExecStart override, retain its full process summary and remove the override afterward. It does not enable system observation or a public endpoint. Missing/dropped/truncated records invalidate timing evidence.

Systemd allows a 6 s emergency stop safeguard; acceptance still requires the application's exit within 5 s. No MemoryMax/CPUQuota or artificial ten-client CPU cap is configured.

For the complete future Pi campaign and actual tooling commands, consult docs/s04-runbook.md in the source distribution. No operation in this guide is authorized by merely reading it.
