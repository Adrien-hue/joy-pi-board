"""Controlled child processes only, without Debian tools or maintainer scripts."""
import contextlib
import importlib.util
import io
import os
import pathlib
import sys
import tempfile
import unittest

sys.dont_write_bytecode = True

spec = importlib.util.spec_from_file_location(
    "build_support", pathlib.Path(__file__).with_name("build_support.py")
)
support = importlib.util.module_from_spec(spec)
spec.loader.exec_module(support)


class BuildLogTests(unittest.TestCase):
    def test_failure_exposes_original_output_preserves_log_and_exit(self):
        with tempfile.TemporaryDirectory() as directory:
            log = pathlib.Path(directory) / "build.log"
            stderr = io.StringIO()
            with contextlib.redirect_stderr(stderr), self.assertRaises(SystemExit) as error:
                support.run_logged(
                    [sys.executable, "-c",
                     "import sys; print('build progress', flush=True); "
                     "print('controlled build failure', file=sys.stderr); sys.exit(2)"],
                    cwd=directory, env=os.environ.copy(), log_path=log,
                )
            self.assertEqual(error.exception.code, 2)
            saved = log.read_text(encoding="utf-8")
            self.assertIn("build progress", saved)
            self.assertIn("controlled build failure", saved)
            self.assertIn(saved, stderr.getvalue())
            self.assertIn("exited 2", stderr.getvalue())

    def test_success_retains_log_without_failure(self):
        with tempfile.TemporaryDirectory() as directory:
            log = pathlib.Path(directory) / "build.log"
            stderr = io.StringIO()
            with contextlib.redirect_stderr(stderr):
                support.run_logged(
                    [sys.executable, "-c", "print('controlled success')"],
                    cwd=directory, env=os.environ.copy(), log_path=log,
                )
            self.assertEqual(log.read_text(encoding="utf-8").strip(), "controlled success")
            self.assertEqual(stderr.getvalue(), "")


if __name__ == "__main__":
    unittest.main()
