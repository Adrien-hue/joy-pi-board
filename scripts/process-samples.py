#!/usr/bin/env python3
"""External Linux-only process observer. Never imported or invoked by Board."""
import argparse, json, os, pathlib, time
p=argparse.ArgumentParser()
p.add_argument("--pid",type=int,required=True);p.add_argument("--output",required=True)
p.add_argument("--seconds",type=int,default=720);p.add_argument("--warmup",type=int,default=120)
a=p.parse_args()
if a.pid<=0 or not 1<=a.seconds<=86400 or not 0<=a.warmup<a.seconds:p.error("invalid sampling range")
ticks=os.sysconf("SC_CLK_TCK");page=os.sysconf("SC_PAGE_SIZE")
def sample():
    fields=pathlib.Path(f"/proc/{a.pid}/stat").read_text().rsplit(")",1)[1].split()
    return {"start_ticks":int(fields[19]),"cpu_ticks":int(fields[11])+int(fields[12]),"rss_bytes":int(fields[21])*page}
start=time.monotonic();initial=sample();previous=None
with open(a.output,"x",encoding="utf8") as f:
    f.write(json.dumps({"scope":"external process observer; NOT Pi PASS","pid":a.pid,"ticks_per_second":ticks,"page_bytes":page,"warmup_seconds":a.warmup})+"\n")
    for i in range(a.seconds+1):
        time.sleep(max(0,start+i-time.monotonic()))
        row=sample();now=time.monotonic()
        if row["start_ticks"]!=initial["start_ticks"]:raise SystemExit("PID changed; run invalid")
        row.update(elapsed_s=now-start,warmup=now-start<a.warmup,utc=time.time())
        if previous:
            row["cpu_one_core_percent"]=100*(row["cpu_ticks"]-previous["cpu_ticks"])/ticks/(row["elapsed_s"]-previous["elapsed_s"])
        f.write(json.dumps(row)+"\n");f.flush();previous=row
