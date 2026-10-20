#!/usr/bin/env python3
"""Opt-in generated-prerm tests with stubs, only inside the verified CI container.

Never installs a package or invokes a real service manager. These tests do not
prove systemd lifecycle behavior. No script is executed without --execute.
"""
import argparse
import hashlib
import json
import os
import pathlib
import subprocess
import sys
import tempfile

sys.dont_write_bytecode = True
from build_support import verify_removal_hook


def verify_container(info, hostname):
    """Require the actual docker inspect record captured before this run."""
    host = info["HostConfig"]
    mounts = info["Mounts"]
    if (info["Name"] != "/board-package-test"
            or hostname != info["Id"][:12]
            or host["NetworkMode"] != "none"
            or host["Privileged"] or host.get("Devices")
            or host.get("PidMode") or host.get("IpcMode") == "host"
            or len(mounts) != 1 or mounts[0]["Type"] != "bind"
            or mounts[0]["Destination"] != "/src" or mounts[0]["RW"]):
        raise ValueError("Disposable container isolation record mismatch")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--execute", action="store_true")
    args = parser.parse_args()
    if not args.execute:
        print("NOT RUN: --execute requires the verified disposable CI container; no action taken")
        return 0
    if sys.platform != "linux" or not pathlib.Path("/.dockerenv").is_file():
        raise SystemExit("Disposable Linux Docker container required")
    if pathlib.Path("/proc/1/comm").read_text().strip() == "systemd":
        raise SystemExit("This is a stub test, not a real-systemd guest")
    inspect_path = pathlib.Path("/work/container-inspect.json")
    info = json.loads(inspect_path.read_text())[0]
    verify_container(info, os.uname().nodename)
    # Verify the running mount as well as the host-side Docker record.
    mounts = subprocess.check_output(
        ["findmnt", "--json", "--target", "/src", "--output", "TARGET,OPTIONS"], text=True)
    source_mount = json.loads(mounts)["filesystems"][0]
    if source_mount["target"] != "/src" or "ro" not in source_mount["options"].split(","):
        raise SystemExit("Source mount is not read-only")
    hook = pathlib.Path("/work/one/extracted/DEBIAN/prerm")
    verify_removal_hook(hook.parent, pathlib.Path(__file__).with_name("prerm"))
    runtime = pathlib.Path("/run/systemd/system")
    parent = runtime.parent
    if parent.is_symlink() or runtime.exists() or runtime.is_symlink() or (parent / "private").exists():
        raise SystemExit("Unexpected systemd runtime state; refuse to modify it")
    output = pathlib.Path("/work/prerm-stub-tests.json")
    if output.exists():
        raise SystemExit("Evidence output already exists")
    report = {
        "scope": "Generated prerm with controlled commands in disposable container; NOT systemd lifecycle",
        "containerId": info["Id"],
        "inspectionSHA256": hashlib.sha256(inspect_path.read_bytes()).hexdigest(),
        "prermSHA256": hashlib.sha256(hook.read_bytes()).hexdigest(),
        "cases": [], "result": "INCOMPLETE",
    }
    created_parent = not parent.exists()
    try:
        if created_parent:
            parent.mkdir()
        with tempfile.TemporaryDirectory(prefix="prerm-stubs-", dir="/work") as directory:
            temporary = pathlib.Path(directory)
            trace = temporary / "calls.jsonl"
            stub = '''#!/usr/bin/python3
import json, os, pathlib, sys
name = pathlib.Path(sys.argv[0]).name
with open(os.environ["TRACE"], "a") as log:
    log.write(json.dumps([name] + sys.argv[1:]) + "\\n")
if name == "deb-systemd-invoke" and sys.argv[1:] == ["stop", "joy-pi-board.service"]:
    sys.exit(int(os.environ["STOP_EXIT"]))
if name == "systemctl" and len(sys.argv) == 5 and sys.argv[1] == "show" and sys.argv[3:] == ["--value", "joy-pi-board.service"]:
    if os.environ["SHOW_EXIT"] != "0": sys.exit(int(os.environ["SHOW_EXIT"]))
    prop = sys.argv[2]
    values = {"--property=ActiveState": "STATE", "--property=MainPID": "MAIN_PID", "--property=ControlPID": "CONTROL_PID"}
    if prop in values:
        print(os.environ[values[prop]])
        sys.exit(0)
sys.exit(99)
'''
            for name in ["deb-systemd-invoke", "systemctl"]:
                path = temporary / name
                path.write_text(stub)
                path.chmod(0o755)

            def case(name, action="remove", running_systemd=True, root=None,
                     state="inactive", main_pid="0", control_pid="0",
                     stop_exit=0, show_exit=0, expected_exit=0, repeats=1):
                if running_systemd:
                    runtime.mkdir(exist_ok=True)  # Empty test sentinel, no daemon/bus.
                elif runtime.exists():
                    runtime.rmdir()
                trace.write_text("")
                env = {"PATH": str(temporary), "TRACE": str(trace), "LC_ALL": "C",
                       "STATE": state, "MAIN_PID": main_pid, "CONTROL_PID": control_pid,
                       "STOP_EXIT": str(stop_exit), "SHOW_EXIT": str(show_exit)}
                if root is not None:
                    env["DPKG_ROOT"] = root
                calls_expected = []
                if action == "remove" and running_systemd and not root:
                    calls_expected.append(["deb-systemd-invoke", "stop", "joy-pi-board.service"])
                    if not stop_exit:
                        calls_expected.append(["systemctl", "show", "--property=ActiveState", "--value", "joy-pi-board.service"])
                        if not show_exit and state in ("inactive", "failed"):
                            for prop in ["MainPID", "ControlPID"]:
                                calls_expected.append(["systemctl", "show", "--property=" + prop, "--value", "joy-pi-board.service"])
                results = []
                for _ in range(repeats):
                    result = subprocess.run(["/bin/sh", str(hook), action], env=env,
                                            capture_output=True, text=True, timeout=10)
                    results.append({"exit": result.returncode, "stdout": result.stdout, "stderr": result.stderr})
                calls = [json.loads(line) for line in trace.read_text().splitlines()]
                passed = calls == calls_expected * repeats and all(r["exit"] == expected_exit for r in results)
                report["cases"].append({"name": name, "result": "PASS" if passed else "FAIL", "calls": calls, "processes": results})
                output.write_text(json.dumps(report, indent=2) + "\n")
                if not passed:
                    raise RuntimeError("Generated prerm regression: " + name)

            case("remove-stopped")
            # Distinguish failed state from failure of the stop operation.
            case("failed-unit", state="failed")
            case("idempotent-repeat", repeats=2)
            case("empty-root", root="")
            case("offline-no-systemd", running_systemd=False)
            case("alternate-root", root="/isolated-root")
            for action in ["upgrade", "deconfigure", "failed-upgrade"]:
                case("no-action-" + action, action=action)
            case("stop-error", stop_exit=7, expected_exit=7)
            case("state-query-error", show_exit=8, expected_exit=8)
            case("policy-declined-still-active", state="active", expected_exit=1)
            case("still-deactivating", state="deactivating", expected_exit=1)
            case("unknown-state", state="", expected_exit=1)
            case("residual-main-pid", main_pid="123", expected_exit=1)
            case("residual-control-pid", control_pid="123", expected_exit=1)
        report["result"] = "PASS (stub control flow only)"
    finally:
        output.write_text(json.dumps(report, indent=2) + "\n")
        if runtime.exists():
            runtime.rmdir()  # Only the empty sentinel created above.
        if created_parent:
            parent.rmdir()
    print(json.dumps({"result": report["result"], "cases": len(report["cases"]), "evidence": str(output)}))
    return 0


if __name__ == "__main__":
    sys.exit(main())
