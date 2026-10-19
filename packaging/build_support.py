"""Build subprocess logging only; never installs packages or operates services."""
import subprocess
import sys


def run_logged(command, *, cwd, env, log_path):
    """Retain complete output and expose the actual failure in the CI console."""
    with log_path.open("w", encoding="utf-8") as log:
        result = subprocess.run(
            command, cwd=cwd, env=env, stdout=log, stderr=subprocess.STDOUT
        )
    if result.returncode:
        print(log_path.read_text(encoding="utf-8", errors="replace"),
              file=sys.stderr, end="")
        print(f"{command[0]} exited {result.returncode}; full log: {log_path}",
              file=sys.stderr)
        raise SystemExit(result.returncode)
