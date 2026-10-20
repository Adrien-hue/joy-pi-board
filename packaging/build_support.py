"""Build subprocess logging only; never installs packages or operates services."""
import subprocess
import sys


def verify_removal_hook(control_dir, template):
    """Read-only guard on the extracted hook; never executes maintainer code."""
    hook = control_dir / "prerm"
    if not hook.is_file():
        raise ValueError("VM-F01: package is missing the removal stop hook")
    expected = template.read_text(encoding="utf-8").replace("#DEBHELPER#", "").strip()
    if hook.read_text(encoding="utf-8").strip() != expected:
        raise ValueError("Removal hook differs from reviewed source; review generated scripts")


def run_logged(command, *, cwd, env, log_path, check=True):
    """Expose failures; check=False lets callers save evidence before failing."""
    with log_path.open("w", encoding="utf-8") as log:
        result = subprocess.run(
            command, cwd=cwd, env=env, stdout=log, stderr=subprocess.STDOUT
        )
    if result.returncode:
        print(log_path.read_text(encoding="utf-8", errors="replace"),
              file=sys.stderr, end="")
        print(f"{command[0]} exited {result.returncode}; full log: {log_path}",
              file=sys.stderr)
        if check:
            raise SystemExit(result.returncode)
    return result
