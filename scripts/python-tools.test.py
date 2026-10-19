"""Controlled data test, never a physical resource measurement."""
import json,pathlib,subprocess,sys,tempfile
with tempfile.TemporaryDirectory() as d:
    root=pathlib.Path(d);raw=root/"raw.jsonl";out=root/"report.json"
    rows=[{"ticks_per_second":100}]+[{"start_ticks":5,"elapsed_s":i,"cpu_ticks":i,"rss_bytes":1000+i,"warmup":i<2} for i in range(10)]
    raw.write_text("\n".join(json.dumps(r) for r in rows))
    subprocess.run([sys.executable,"scripts/process-report.py",str(raw),str(out)],check=True)
    report=json.loads(out.read_text());assert report["cpu_mean_one_core_percent"]==1;assert report["rss_max_bytes"]==1009
    rows[-1]["start_ticks"]=6;raw.write_text("\n".join(json.dumps(r) for r in rows))
    assert subprocess.run([sys.executable,"scripts/process-report.py",str(raw),str(root/"bad.json")]).returncode!=0
for name in ["build.py","postinst","rules"]:
    assert (pathlib.Path("packaging")/name).exists()
compile(pathlib.Path("packaging/build.py").read_text(),"packaging/build.py","exec")
compile(pathlib.Path("scripts/process-samples.py").read_text(),"scripts/process-samples.py","exec")
print("PASS: external statistics and Python syntax; no system/package operation")
subprocess.run([sys.executable,"packaging/test_build_support.py"],check=True)
