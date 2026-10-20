"""Controlled child processes only, without Debian tools or maintainer scripts."""
import contextlib
import copy
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
container_spec = importlib.util.spec_from_file_location(
    "prerm_container", pathlib.Path(__file__).with_name("test_prerm_container.py")
)
# The container module imports build_support, without executing its main().
sys.modules["build_support"] = support
container_test = importlib.util.module_from_spec(container_spec)
container_spec.loader.exec_module(container_test)


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

    def test_validation_failure_exposes_report_and_returns_status_for_evidence(self):
        with tempfile.TemporaryDirectory() as directory:
            log = pathlib.Path(directory) / "lintian.log"
            stderr = io.StringIO()
            with contextlib.redirect_stderr(stderr):
                result = support.run_logged(
                    [sys.executable, "-c",
                     "import sys; print('E: controlled validation failure'); sys.exit(2)"],
                    cwd=directory, env=os.environ.copy(), log_path=log, check=False,
                )
            self.assertEqual(result.returncode, 2)
            self.assertIn("E: controlled validation failure", stderr.getvalue())
            self.assertIn(log.read_text(encoding="utf-8"), stderr.getvalue())

    def test_validation_success_returns_status_and_retains_warnings(self):
        with tempfile.TemporaryDirectory() as directory:
            log = pathlib.Path(directory) / "lintian.log"
            result = support.run_logged(
                [sys.executable, "-c", "print('W: controlled warning')"],
                cwd=directory, env=os.environ.copy(), log_path=log, check=False,
            )
            self.assertEqual(result.returncode, 0)
            self.assertEqual(log.read_text(encoding="utf-8").strip(), "W: controlled warning")


class RemovalHookInspectionTests(unittest.TestCase):
    def test_historical_package_without_prerm_is_rejected(self):
        control = pathlib.Path(__file__).resolve().parents[1] / "docs/evidence/s04/ci-39313804438/generated-control"
        with self.assertRaisesRegex(ValueError, "missing"):
            support.verify_removal_hook(control, pathlib.Path(__file__).with_name("prerm"))

    def test_generated_hook_matches_reviewed_source(self):
        template = pathlib.Path(__file__).with_name("prerm")
        with tempfile.TemporaryDirectory() as directory:
            control = pathlib.Path(directory)
            (control / "prerm").write_text(template.read_text().replace("#DEBHELPER#", ""))
            support.verify_removal_hook(control, template)

    def test_mutated_or_appended_hook_is_rejected(self):
        template = pathlib.Path(__file__).with_name("prerm")
        original = template.read_text().replace("#DEBHELPER#", "")
        for changed in [original.replace("stop joy-pi-board", "start joy-pi-board"),
                        original + "\nsystemctl start joy-pi-board.service\n"]:
            with self.subTest(changed=changed), tempfile.TemporaryDirectory() as directory:
                control = pathlib.Path(directory)
                (control / "prerm").write_text(changed)
                with self.assertRaisesRegex(ValueError, "differs"):
                    support.verify_removal_hook(control, template)


class ContainerGuardTests(unittest.TestCase):
    def test_only_named_isolated_readonly_target_is_accepted(self):
        info = {"Name": "/board-package-test", "Id": "123456789abc" + "0" * 52,
                "HostConfig": {"NetworkMode": "none", "Privileged": False},
                "Mounts": [{"Type": "bind", "Destination": "/src", "RW": False}]}
        container_test.verify_container(info, "123456789abc")
        for field, value in [("NetworkMode", "host"), ("Privileged", True),
                             ("Devices", [{"PathOnHost": "/dev/sda"}]),
                             ("PidMode", "host"), ("IpcMode", "host")]:
            changed = copy.deepcopy(info)
            changed["HostConfig"][field] = value
            with self.subTest(field=field), self.assertRaises(ValueError):
                container_test.verify_container(changed, "123456789abc")
        for mounts in [[{"Type": "bind", "Destination": "/src", "RW": True}],
                       info["Mounts"] * 2, [],
                       [{"Type": "bind", "Destination": "/work", "RW": False}]]:
            changed = copy.deepcopy(info)
            changed["Mounts"] = mounts
            with self.subTest(mounts=mounts), self.assertRaises(ValueError):
                container_test.verify_container(changed, "123456789abc")
        with self.assertRaises(ValueError):
            container_test.verify_container(info, "different-host")
        info["Name"] = "/production"
        with self.assertRaises(ValueError):
            container_test.verify_container(info, "123456789abc")


if __name__ == "__main__":
    unittest.main()
