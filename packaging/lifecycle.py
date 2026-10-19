#!/usr/bin/env python3
"""Opt-in disposable ARM64 VM test. Refuses ordinary hosts and inspection-only guests."""
import argparse,hashlib,json,os,pathlib,subprocess
p=argparse.ArgumentParser();p.add_argument("--execute",action="store_true");p.add_argument("--vm-id");p.add_argument("--package");p.add_argument("--previous");p.add_argument("--output")
a=p.parse_args()
if not a.execute:raise SystemExit("NOT RUN: requires --execute, --vm-id, two packages, output and isolated ARM64 VM marker; no action taken")
if not all([a.vm_id,a.package,a.previous,a.output]):p.error("all execution arguments required")
def run(*args,check=True):return subprocess.run(args,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,check=check)
if os.geteuid()!=0 or pathlib.Path("/proc/1/comm").read_text().strip()!="systemd":raise SystemExit("real systemd PID 1/root required")
if os.uname().machine not in ("aarch64","arm64"):raise SystemExit("exact ARM64 executable required")
if run("systemd-detect-virt","--vm",check=False).returncode:raise SystemExit("identified VM required")
marker=json.loads(pathlib.Path("/run/joy-pi-board-disposable.json").read_text())
if marker!={"id":a.vm_id,"disposable":True,"hostWritableMounts":False}:raise SystemExit("isolation attestation mismatch")
mounts=run("findmnt","-rn","-o","FSTYPE,OPTIONS").stdout.splitlines()
for line in mounts:
    kind,options=line.split(maxsplit=1)
    if kind in ("9p","virtiofs","vboxsf","fuse.vmhgfs-fuse","cifs","nfs","nfs4") and "rw" in options.split(","):raise SystemExit("writable shared filesystem forbidden")
if run("dpkg-query","-W","joy-pi-board",check=False).returncode==0:raise SystemExit("guest must start without Board package")
if run("getent","passwd","joy-pi-board",check=False).returncode==0:raise SystemExit("guest account already exists")
output=pathlib.Path(a.output)
if output.exists():raise SystemExit("output exists")
packages=[pathlib.Path(a.previous).resolve(),pathlib.Path(a.package).resolve()]
for pkg in packages:
    if run("dpkg-deb","-f",str(pkg),"Package").stdout.strip()!="joy-pi-board":raise SystemExit("wrong package")
    if run("dpkg-deb","-f",str(pkg),"Architecture").stdout.strip()!="arm64":raise SystemExit("wrong architecture")
events=[]
def action(*args):
    r=run(*args,check=False);events.append({"command":list(args),"exit":r.returncode,"output":r.stdout});save()
    if r.returncode:raise SystemExit("test operation failed; restore VM snapshot; evidence retained")
def save():output.write_text(json.dumps({"scope":"disposable real-systemd ARM64 VM; NOT Pi acceptance","vm":marker,"packageHashes":[hashlib.sha256(x.read_bytes()).hexdigest() for x in packages],"events":events},indent=2)+"\n")
def inactive():
    assert run("systemctl","is-active","--quiet","joy-pi-board",check=False).returncode!=0
def verify():
    assert run("systemctl","show","joy-pi-board","-p","User","--value").stdout.strip()=="joy-pi-board"
    text=run("systemctl","cat","joy-pi-board").stdout
    assert "joy-pi-health" not in text
def identity(pkg):
    # Extract in guest tmp only; never executes package scripts.
    import tempfile
    with tempfile.TemporaryDirectory() as d:
        action("dpkg-deb","-x",str(pkg),d)
        assert pathlib.Path(d,"usr/bin/joy-pi-board").read_bytes()==pathlib.Path("/usr/bin/joy-pi-board").read_bytes()
try:
    action("dpkg","-i",str(packages[0]));inactive();verify();identity(packages[0])
    assert run("systemctl","is-enabled","joy-pi-board",check=False).stdout.strip()=="disabled"
    conf=pathlib.Path("/etc/default/joy-pi-board");conf.write_text(conf.read_text()+"\n# administrator-preserved\nJOY_PI_BOARD_LISTEN=127.0.0.1:8081\nJOY_PI_BOARD_HEALTH_URL=http://127.0.0.1:1/v1/snapshot\n");saved=conf.read_bytes()
    action("systemctl","start","joy-pi-board");action("systemctl","is-active","--quiet","joy-pi-board")
    import urllib.request,time
    end=time.monotonic()+2
    while True:
        try:
            with urllib.request.urlopen("http://127.0.0.1:8081/",timeout=.2) as response:assert response.status==200
            break
        except OSError:
            if time.monotonic()>=end:raise
    with urllib.request.urlopen("http://127.0.0.1:8081/api/v1/overview") as response:
        value=json.load(response);assert value["health"]["snapshot"] is None
    action("systemctl","enable","joy-pi-board");action("systemctl","stop","joy-pi-board");inactive()
    action("dpkg","--force-confold","-i",str(packages[1]));inactive();assert conf.read_bytes()==saved;identity(packages[1])
    assert run("systemctl","is-enabled","joy-pi-board",check=False).stdout.strip()=="enabled"
    action("systemctl","disable","joy-pi-board");action("systemctl","mask","joy-pi-board")
    action("dpkg","--force-confold","-i",str(packages[0]));inactive();assert conf.read_bytes()==saved;identity(packages[0])
    assert run("systemctl","is-enabled","joy-pi-board",check=False).stdout.strip()=="masked"
    action("systemctl","unmask","joy-pi-board");action("systemctl","start","joy-pi-board")
    action("dpkg","--remove","joy-pi-board");inactive();assert conf.exists()
    assert run("pgrep","-x","joy-pi-board",check=False).returncode!=0
    action("dpkg","--purge","joy-pi-board");assert not conf.exists()
    action("getent","passwd","joy-pi-board")
    action("dpkg","-i",str(packages[1]));inactive();identity(packages[1]);action("dpkg","--purge","joy-pi-board")
    events.append({"result":"PASS within listed VM cases; collision/interrupted-upgrade and Pi cases separate"});save()
except BaseException:
    events.append({"result":"FAIL/incomplete; restore disposable VM snapshot, no host cleanup attempted"});save();raise
