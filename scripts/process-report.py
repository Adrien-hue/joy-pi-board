#!/usr/bin/env python3
"""Summarize an external observer run without assigning physical PASS."""
import argparse,json,math
p=argparse.ArgumentParser();p.add_argument("input");p.add_argument("output");a=p.parse_args()
rows=[json.loads(l) for l in open(a.input,encoding="utf8")];meta=rows.pop(0)
if len({r["start_ticks"] for r in rows})!=1:raise SystemExit("process identity changed")
measured=[r for r in rows if not r["warmup"]]
if len(measured)<2:raise SystemExit("insufficient samples")
gaps=[b["elapsed_s"]-a["elapsed_s"] for a,b in zip(measured,measured[1:])]
if any(g<=0 or g>2 for g in gaps):raise SystemExit("sample gap; retain raw run and repeat")
rss=sorted(r["rss_bytes"] for r in measured)
cpu=100*(measured[-1]["cpu_ticks"]-measured[0]["cpu_ticks"])/meta["ticks_per_second"]/(measured[-1]["elapsed_s"]-measured[0]["elapsed_s"])
report={"scope":"external tooling summary; no physical PASS","samples":len(measured),"cpu_mean_one_core_percent":cpu,"rss_max_bytes":max(rss),"rss_p95_bytes":rss[math.ceil(.95*len(rss))-1],"rss_median_bytes":rss[math.ceil(.5*len(rss))-1],"rss_start_bytes":measured[0]["rss_bytes"],"rss_end_bytes":measured[-1]["rss_bytes"]}
with open(a.output,"x",encoding="utf8") as f:json.dump(report,f,indent=2)
