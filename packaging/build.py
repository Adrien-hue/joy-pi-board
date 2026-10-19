#!/usr/bin/env python3
"""Archive construction only. Never installs, starts services or executes maintainer scripts."""
import argparse, datetime, gzip, hashlib, json, os, pathlib, re, shutil, subprocess, sys
sys.dont_write_bytecode = True
from build_support import run_logged
p=argparse.ArgumentParser()
p.add_argument("--artifact",required=True);p.add_argument("--output",required=True)
p.add_argument("--metadata");p.add_argument("--test-only",action="store_true")
a=p.parse_args()
root=pathlib.Path(__file__).resolve().parent
artifact=pathlib.Path(a.artifact).resolve();out=pathlib.Path(a.output).resolve()
if out.exists():p.error("output identity already exists")
manifest=json.loads((artifact/"provenance.json").read_text())
if not (artifact/"COMPLETE").exists():p.error("incomplete artifact")
binary=artifact/manifest["binary"]["file"]
digest=lambda b:hashlib.sha256(b).hexdigest()
if digest(binary.read_bytes())!=manifest["binary"]["sha256"]:p.error("binary hash mismatch")
for entry in manifest["inputs"]:
    if entry["path"].startswith("packaging/"):
        current=root.parent/entry["path"]
        if not current.is_file() or digest(current.read_bytes())!=entry["sha256"]:
            p.error("packaging input changed; create a new artifact identity")
for line in (artifact/"SHA256SUMS").read_text().splitlines():
    expected_hash, relative=line.split("  ",1)
    if pathlib.Path(relative).name!=relative or digest((artifact/relative).read_bytes())!=expected_hash:
        p.error("artifact checksum mismatch")
if not a.test_only and not a.metadata:p.error("BLOCKED: established Maintainer and distribution copyright/license metadata required")
if a.test_only:
    maintainer="Internal packaging test fixture <fixture@example.invalid>"
    copyright_text="INTERNAL TEST FIXTURE — NOT FOR DISTRIBUTION. No Board license is asserted.\n"
else:
    meta=json.loads(pathlib.Path(a.metadata).read_text())
    maintainer=meta["maintainer"];copyright_text=pathlib.Path(meta["copyrightFile"]).read_text()
    if not meta.get("established") or not re.fullmatch(r"[^<>\r\n]+ <[^<>\s]+@[^<>\s]+>",maintainer) or "example." in maintainer or not copyright_text.strip():p.error("established metadata missing")
expected={"dpkg":"1.22.22","dpkg-dev":"1.22.22","debhelper":"13.24.2"}
for name,version in expected.items():
    actual=subprocess.check_output(["dpkg-query","-W","-f=${Version}",name],text=True)
    if actual!=version:raise SystemExit(f"Packaging lock mismatch: {name}: {actual} != {version}")
out.mkdir(parents=True);work=out/"source";work.mkdir();debian=work/"debian";debian.mkdir()
version=manifest["packageVersion"]
if manifest["dirty"] or a.test_only:version += "+dev."+manifest["id"]
if not re.fullmatch(r"[0-9][a-z0-9.+:~\-]*",version):p.error("unsafe package version")
(debian/"control").write_text(f"""Source: joy-pi-board
Section: web
Priority: optional
Maintainer: {maintainer}
Build-Depends: debhelper-compat (= 13)
Standards-Version: 4.7.0
Rules-Requires-Root: no
Homepage: https://github.com/Adrien-hue/joy-pi-board

Package: joy-pi-board
Architecture: arm64
Depends: adduser, ${{misc:Depends}}
Description: Joy Pi Home LAN overview
 Embedded read-only presentation of Joy Pi Health observations.
 No Health availability or Node runtime is required.
""")
date=datetime.datetime.fromtimestamp(manifest["sourceDateEpoch"],datetime.timezone.utc).strftime("%a, %d %b %Y %H:%M:%S +0000")
(debian/"changelog").write_text(f"joy-pi-board ({version}) unstable; urgency=medium\n\n  * Private artifact {manifest['id']}; {manifest['status']}.\n\n -- {maintainer}  {date}\n")
shutil.copyfile(root/"rules",debian/"rules");(debian/"rules").chmod(0o755)
shutil.copyfile(root/"postinst",debian/"joy-pi-board.postinst");(debian/"joy-pi-board.postinst").chmod(0o755)
shutil.copyfile(root/"joy-pi-board.service",debian/"joy-pi-board.service")
payload=work/"payload"
def put(path,data,mode=0o644):
    target=payload/path;target.parent.mkdir(parents=True,exist_ok=True);target.write_bytes(data);target.chmod(mode)
put("usr/bin/joy-pi-board",binary.read_bytes(),0o755)
put("etc/default/joy-pi-board",(root/"joy-pi-board.default").read_bytes())
put("usr/share/doc/joy-pi-board/copyright",copyright_text.encode())
put("usr/share/doc/joy-pi-board/build.json",(artifact/"provenance.json").read_bytes())
put("usr/share/doc/joy-pi-board/assets.json",(artifact/"assets.json").read_bytes())
put("usr/share/doc/joy-pi-board/README", (root/"README.md").read_bytes())
put("usr/share/man/man1/joy-pi-board.1", (root/"joy-pi-board.1").read_bytes())
put("usr/share/lintian/overrides/joy-pi-board", (root/"lintian-overrides").read_bytes())
if (artifact/"THIRD-PARTY-NOTICES").exists():put("usr/share/doc/joy-pi-board/THIRD-PARTY-NOTICES",(artifact/"THIRD-PARTY-NOTICES").read_bytes())
else:raise SystemExit("Missing third-party notices")
epoch=str(manifest["sourceDateEpoch"])
for path in work.rglob("*"):os.utime(path,(int(epoch),int(epoch)))
env={**os.environ,"SOURCE_DATE_EPOCH":epoch,"LC_ALL":"C.UTF-8","TZ":"UTC","DEB_BUILD_OPTIONS":"nostrip noautodbgsym"}
tools=subprocess.check_output(["dpkg-query","-W","-f=${binary:Package}\t${Version}\n"],text=True)
(out/"packages.lock").write_text(tools)
run_logged(["dpkg-buildpackage","-b","-us","-uc","-a","arm64","-d"],cwd=work,env=env,log_path=out/"build.log")
packages=list(out.glob("*.deb"))
if len(packages)!=1:raise SystemExit("Unexpected package/debug output")
package=packages[0];extract=out/"extracted"
subprocess.run(["dpkg-deb","-R",str(package),str(extract)],check=True)
installed=extract/"usr/bin/joy-pi-board"
if installed.read_bytes()!=binary.read_bytes():raise SystemExit("Packaged binary mutated")
# Capture the actual generated scripts, not just source templates.
control=extract/"DEBIAN"
for script in ["preinst","postinst","prerm","postrm"]:
    f=control/script
    if f.exists():subprocess.run(["sh","-n",str(f)],check=True)
if "/etc/default/joy-pi-board" not in (control/"conffiles").read_text():raise SystemExit("Missing conffile")
if gzip.decompress((extract/"usr/share/man/man1/joy-pi-board.1.gz").read_bytes()) != (root/"joy-pi-board.1").read_bytes():raise SystemExit("Packaged manual page mismatch")
if (extract/"usr/share/lintian/overrides/joy-pi-board").read_bytes() != (root/"lintian-overrides").read_bytes():raise SystemExit("Packaged lintian overrides mismatch")
contents=subprocess.check_output(["dpkg-deb","-c",str(package)],text=True)
(out/"contents.txt").write_text(contents)
lint=run_logged(["lintian","--allow-root","--show-overrides",str(package)],cwd=work,env=env,log_path=out/"lintian.log",check=False)
# Keep warnings and the two narrowly scoped overrides visible even on success.
if not lint.returncode:print((out/"lintian.log").read_text(encoding="utf-8"),end="")
report={"scope":"archive build/extraction only; no maintainer execution or systemd lifecycle","distribution":not a.test_only,"status":manifest["status"],"id":manifest["id"],"binarySHA256":digest(installed.read_bytes()),"package":package.name,"packageSHA256":digest(package.read_bytes()),"toolsSHA256":digest(tools.encode()),"lintianExit":lint.returncode}
(out/"package-evidence.json").write_text(json.dumps(report,indent=2)+"\n")
(out/"SHA256SUMS").write_text(report["packageSHA256"]+"  "+package.name+"\n")
print(json.dumps(report))
if lint.returncode:
    raise SystemExit("Package validation blocked by lintian; inspect retained report")
