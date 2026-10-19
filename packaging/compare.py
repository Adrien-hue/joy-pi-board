#!/usr/bin/env python3
"""Compare two package builds consuming the same artifact. Does not promote either."""
import hashlib,json,pathlib,sys
a,b=map(pathlib.Path,sys.argv[1:])
pa=list(a.glob("*.deb"));pb=list(b.glob("*.deb"))
if len(pa)!=1 or len(pb)!=1:raise SystemExit("one package per build required")
ha=hashlib.sha256(pa[0].read_bytes()).hexdigest();hb=hashlib.sha256(pb[0].read_bytes()).hexdigest()
if ha!=hb:raise SystemExit("package reproducibility mismatch; retain both builds")
print(json.dumps({"scope":"same-input archive comparison, not lifecycle or Pi acceptance","sha256":ha,"selected":str(pa[0]),"comparison":str(pb[0])}))
