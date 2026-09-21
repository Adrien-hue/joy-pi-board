"""Developer-only corpus authoring tool; Python exact integers, no JS Number.

cases.json is the committed, authoritative input/expectation corpus consumed by
both languages. This tool is not needed to build or run Board or its tests.
Run with --check to compare bytes without rewriting. No Health files are edited.
"""
import base64
import copy
import json
from pathlib import Path
import sys

HERE = Path(__file__).resolve().parent
ROOT = HERE.parent.parent
base = json.loads((ROOT / "docs/examples/overview-current.json").read_bytes())["health"]["snapshot"]
partial = json.loads((ROOT / "docs/examples/overview-partial.json").read_bytes())["health"]["snapshot"]
cases = []
uint_keys = {"logical_cpu_count", "total_bytes", "available_bytes", "used_bytes", "uptime_seconds", "rx_bytes", "tx_bytes"}
flags = ["thermal_throttling_active", "thermal_throttling_occurred_since_boot", "undervoltage_active", "undervoltage_occurred_since_boot"]


def dump(obj):
    return json.dumps(obj, ensure_ascii=True, separators=(",", ":"), allow_nan=False)


def projection(obj, shape=base, key=""):
    if obj is None:
        return None
    if key in uint_keys:
        return "u64:" + str(obj)
    if isinstance(obj, dict):
        return {k: projection(obj[k], shape[k], k) for k in shape}
    if isinstance(obj, list):
        item_shape = {"path": "", "code": "", "message": ""} if key == "issues" else base["network"]["interfaces"][0]
        return [projection(v, item_shape) for v in obj]
    return obj


def add(name, obj, kind="valid", expected=None):
    data = obj if isinstance(obj, bytes) else (obj.encode() if isinstance(obj, str) else dump(obj).encode())
    entry = {"id": name, "input_base64": base64.b64encode(data).decode(), "classification": kind}
    if kind == "valid":
        entry["expected"] = projection(json.loads(data) if expected is None else expected)
    cases.append(entry)


def change(obj, path, value=None, delete=False):
    keys = path.strip("/").split("/")
    for k in keys[:-1]:
        obj = obj[int(k)] if isinstance(obj, list) else obj[k]
    k = int(keys[-1]) if isinstance(obj, list) else keys[-1]
    if delete:
        del obj[k]
    else:
        obj[k] = value


def mutated(name, path, value, kind="invalid_response"):
    obj = copy.deepcopy(base)
    change(obj, path, value)
    add(name, obj, kind)


def issues(obj, code="temporarily_unavailable", message="Metric unavailable."):
    paths = []
    for p in ["/uptime_seconds", "/cpu/logical_cpu_count", "/load", "/memory", "/root_filesystem", "/raspberry_pi/soc_temperature_celsius"] + ["/raspberry_pi/" + f for f in flags]:
        val = obj
        for k in p.strip("/").split("/"):
            val = val[k]
        if val is None or isinstance(val, dict) and all(v is None for v in val.values()):
            paths.append(p)
    if obj["network"] is None:
        paths.append("/network")
    else:
        for i, n in enumerate(obj["network"]["interfaces"]):
            for k in ("state", "rx_bytes", "tx_bytes"):
                if n[k] is None:
                    paths.append(f"/network/interfaces/{i}/{k}")
    if obj["cpu"]["utilization_percent"] is None:
        paths.append("/cpu/utilization_percent")
    obj["issues"] = [{"path": p, "code": code, "message": message} for p in paths]
    return obj


add("complete", base)
add("partial-documentary", partial)
for code in ["unsupported", "permission_denied", "temporarily_unavailable"]:
    for path in ["/uptime_seconds", "/cpu/logical_cpu_count", "/cpu/utilization_percent", "/network", "/raspberry_pi/soc_temperature_celsius"] + [f"/network/interfaces/0/{k}" for k in ["state", "rx_bytes", "tx_bytes"]]:
        obj = copy.deepcopy(base)
        change(obj, path, None)
        add("partial-" + code + path.replace("/", "-"), issues(obj, code))
    for group in ["load", "memory", "root_filesystem"]:
        obj = copy.deepcopy(base)
        obj[group] = dict.fromkeys(obj[group])
        add(f"partial-{code}-{group}", issues(obj, code))
    obj = copy.deepcopy(base)
    for f in flags:
        obj["raspberry_pi"][f] = None
    add(f"partial-{code}-firmware", issues(obj, code, "Firmware text need not match generic messages."))

zero = copy.deepcopy(base)
zero["cpu"]["utilization_percent"] = 0
zero["uptime_seconds"] = 0
for g in ["memory", "root_filesystem", "load"]:
    zero[g] = dict.fromkeys(zero[g], 0)
for n in zero["network"]["interfaces"]:
    n["rx_bytes"] = n["tx_bytes"] = 0
for f in flags:
    zero["raspberry_pi"][f] = False
zero["raspberry_pi"]["soc_temperature_celsius"] = 0
add("zero-and-false", zero)
true = copy.deepcopy(base)
for f in flags:
    true["raspberry_pi"][f] = True
add("all-flags-true-no-issues", true)
mutated("open-network-state", "/network/interfaces/0/state", "dormant", "valid")
mutated("empty-interfaces-with-cpu", "/network/interfaces", [], "valid")

empty = copy.deepcopy(base)
empty["uptime_seconds"] = None
empty["cpu"] = dict.fromkeys(empty["cpu"])
for group in ["load", "memory", "root_filesystem", "raspberry_pi"]:
    empty[group] = dict.fromkeys(empty[group])
empty["network"] = None
issues(empty)
add("no-useful-null-network", empty, "invalid_response")
empty["network"] = {"interfaces": []}
issues(empty)
add("empty-interfaces-not-useful", empty, "invalid_response")
empty["extra_metric"] = 42
add("unknown-not-useful", empty, "invalid_response")
empty["network"]["interfaces"] = [{"name": "eth0", "state": None, "rx_bytes": None, "tx_bytes": None}]
issues(empty)
add("interface-name-alone-useful", empty)


def paths(obj, prefix=""):
    if isinstance(obj, dict):
        for k, v in obj.items():
            yield prefix + "/" + k
            yield from paths(v, prefix + "/" + k)
    if isinstance(obj, list):
        for i, v in enumerate(obj):
            yield from paths(v, prefix + "/" + str(i))


for p in paths(base):
    obj = copy.deepcopy(base)
    change(obj, p, delete=True)
    add("missing" + p.replace("/", "-"), obj, "invalid_response")
    obj = copy.deepcopy(base)
    # Non-null structural/string positions; metric nulls need matching issues.
    change(obj, p, None)
    add("null-without-required-issues" + p.replace("/", "-"), obj, "invalid_response")
for k in ["path", "code", "message"]:
    for variant in ["missing", "null"]:
        obj = copy.deepcopy(partial)
        change(obj, f"/issues/0/{k}", None, delete=variant == "missing")
        add(f"issue-{variant}-{k}", obj, "invalid_response")
for name, value in [("code-internal", "internal_error"), ("code-unknown", "other")]:
    obj = copy.deepcopy(partial); obj["issues"][0]["code"] = value
    add(name, obj, "invalid_response")
for name, value in [("path-root", "/"), ("path-noncanonical", "/network/interfaces/00/state"), ("path-wrong-leaf", "/memory/total_bytes")]:
    obj = copy.deepcopy(partial); obj["issues"][0]["path"] = value
    add(name, obj, "invalid_response")
for name in ["reordered", "extra", "omitted", "duplicate", "empty-message", "firmware-code", "firmware-message"]:
    obj = copy.deepcopy(partial)
    if name == "reordered": obj["issues"].reverse()
    if name == "extra": obj["issues"].append({"path": "/extra", "code": "unsupported", "message": "No."})
    if name == "omitted": obj["issues"].pop()
    if name == "duplicate": obj["issues"][1] = copy.deepcopy(obj["issues"][0])
    if name == "empty-message": obj["issues"][0]["message"] = ""
    if name.startswith("firmware-"):
        index = next(i for i, issue in enumerate(obj["issues"]) if issue["path"].endswith(flags[1]))
        obj["issues"][index][name.split("-")[1]] = "unsupported" if name == "firmware-code" else "different"
    add("issues-" + name, obj, "invalid_response")

for group in ["memory", "root_filesystem", "load"]:
    obj = copy.deepcopy(base); obj[group][next(iter(obj[group]))] = None
    add("mixed-" + group, obj, "invalid_response")
for group in ["memory", "root_filesystem"]:
    mutated("available-over-total-"+group, f"/{group}/available_bytes", 18446744073709551615)
    mutated("used-inconsistent-"+group, f"/{group}/used_bytes", 18446744073709551615)
for f in flags:
    mutated("mixed-firmware-"+f, "/raspberry_pi/"+f, None)
mutated("cpu-count-zero", "/cpu/logical_cpu_count", 0)
for p in ["/host/hostname", "/network/interfaces/0/name", "/network/interfaces/0/state"]:
    mutated("empty"+p.replace("/", "-"), p, "")
mutated("unsorted-interfaces", "/network/interfaces", [{"name": n, "state": "up", "rx_bytes": 0, "tx_bytes": 0} for n in ["z", "a"]])
mutated("duplicate-interfaces", "/network/interfaces", [base["network"]["interfaces"][0]]*2)
for count in [64, 65]:
    mutated(f"interfaces-{count}", "/network/interfaces", [{"name": f"n{i:02}", "state": "unknown", "rx_bytes": 0, "tx_bytes": 0} for i in range(count)], "valid" if count == 64 else "invalid_response")
for names, good in [(["\ue000", "\U00010000"], True), (["\U00010000", "\ue000"], False)]:
    mutated("unicode-name-order-"+str(good), "/network/interfaces", [{"name": n, "state": "up", "rx_bytes": 0, "tx_bytes": 0} for n in names], "valid" if good else "invalid_response")

for value in [0, 9007199254740991, 9007199254740992, 9007199254740993, 18446744073709551615]:
    mutated("uint-"+str(value), "/uptime_seconds", value, "valid")
raw = dump(base)
uptime = '"uptime_seconds":' + str(base["uptime_seconds"])
for token in ["18446744073709551616", "-1", "-0", "1.0", "1e0", "1E+0", '"4"', "true", "false", "00", "+1", ".1", "1.", "NaN", "Infinity", "1"*1000]:
    add("uint-invalid-"+str(len(cases)), raw.replace(uptime, '"uptime_seconds":'+token), "invalid_response")
for p, token, good in [("cpu", "50", True), ("cpu", "100", True), ("cpu", "100.0001", False), ("cpu", "-0", True), ("cpu", "-0.1", False), ("cpu", "1e9999", False), ("cpu", "1e-9999", True)]:
    add("float-cpu-"+token, raw.replace('"utilization_percent":'+str(base["cpu"]["utilization_percent"]), '"utilization_percent":'+token), "valid" if good else "invalid_response")
for token, good in [("50", True), ("-300", True), ("1e308", True), ("1e309", False), ("-1e309", False)]:
    add("temperature-"+token, raw.replace('"soc_temperature_celsius":'+str(base["raspberry_pi"]["soc_temperature_celsius"]), '"soc_temperature_celsius":'+token), "valid" if good else "invalid_response")
for field in ["one_minute", "five_minutes", "fifteen_minutes"]:
    mutated("negative-load-"+field, "/load/"+field, -0.01)
mutated("float-string", "/cpu/utilization_percent", "50")
mutated("float-bool", "/cpu/utilization_percent", True)
mutated("firmware-number", "/raspberry_pi/undervoltage_active", 0)

for version in ["2.0", "", "1.1"]:
    add("unsupported-"+version, {"schema_version": version}, "unsupported_schema_version")
for obj in [{}, {"schema_version": None}, {"schema_version": 1}, [], None, 3, "\"1.0\""]:
    add("bad-version-root-"+str(len(cases)), obj, "invalid_response")
obj = copy.deepcopy(base); obj["Schema_version"] = obj.pop("schema_version")
add("case-sensitive-version", obj, "invalid_response")
obj = copy.deepcopy(base); obj["host"]["Hostname"] = obj["host"].pop("hostname")
add("case-sensitive-hostname", obj, "invalid_response")
extra = copy.deepcopy(partial)
for p in list(paths(extra)):
    value = extra
    for k in p.strip("/").split("/"): value = value[int(k)] if isinstance(value, list) else value[k]
    if isinstance(value, dict): value["future"] = {"n": 18446744073709551615}
extra["future"] = {"__proto__": {"unexpected": True}, "n": 9007199254740993, "array": [False, None]}
add("unknown-members-every-object", extra)
for v in [None, 7, {"hostname": "inherited", "isLosslessNumber": True}]:
    obj = copy.deepcopy(base); obj["host"]["__proto__"] = v
    add("unknown-proto-"+str(len(cases)), obj)
add("unknown-float-overflow", raw[:-1]+',"future":1e9999}')
for fragment in ['"x":0,"x":0', '"x":null,"x":null', '"x":[],"x":[]', '"x":{},"x":{}', '"x":1,"\\u0078":2', '"__proto__":{},"__proto__":{}']:
    add("duplicate-unknown-"+str(len(cases)), raw[:-1]+',"extra":{'+fragment+'}}', "invalid_response")
add("duplicate-known", raw[:-1]+',"schema_version":"1.0"}', "invalid_response")
add("duplicate-escaped-known", raw[:-1]+',"schema_versi\\u006fn":"1.0"}', "invalid_response")
for text_value in [raw[:-1], raw+"{}", raw+" garbage", raw+"\u00a0", '\ufeff'+raw, '{"schema_version":"2.0",}', '{"schema_version":"2.0","x":1,"x":1}', '{"schema_version":"2.0","x":"\\ud800"}', raw.replace('"issues":[]','"issues":[,]')]:
    add("malformed-"+str(len(cases)), text_value, "invalid_response")
for token, good in [('"\\ud800"',False), ('"\\udc00"',False), ('"\\ud800x"',False), ('"\\ud800\\u0041"',False), ('"\\ud83d\\ude00"',True), ('"\\ufffd"',True), ('"\\\\ud800"',True), ('"a\\u0000b"',True), ('"\\u12"',False)]:
    add("unicode-"+str(len(cases)), raw[:-1]+',"extra":'+token+'}', "valid" if good else "invalid_response")
add("surrogate-key", raw[:-1]+',"\\ud800":0}', "invalid_response")
for bad in [b'\xff', b'\xc0\xaf', b'\xed\xa0\x80', b'\xf4\x90\x80\x80']:
    add("utf8-"+str(len(cases)), raw[:-1].encode()+b',"extra":"'+bad+b'"}', "invalid_response")
for depth in [31,32,33,1000]:
    add(f"depth-{depth}", raw[:-1]+',"extra":'+'['*(depth-1)+'0'+']'*(depth-1)+'}', "valid" if depth<=32 else "invalid_response")
add("brackets-in-string-not-depth", raw[:-1]+',"extra":"'+'['*100+'"}')
obj = copy.deepcopy(base)
obj["host"]["hostname"] = "pi-\u00e9-\U0001f600"
add("literal-utf8", json.dumps(obj, ensure_ascii=False, separators=(",", ":")).encode())
add("duplicate-literal-escaped-unicode-key", raw[:-1]+',"extra":{"\u00e9":0,"\\u00e9":0}}', "invalid_response")
for size in [65536, 65537]:
    prefix = raw[:-1].encode()+b',"extra":"'
    padding = size-len(prefix)-2
    data = prefix + "\u00e9".encode()*(padding//2) + b'x'*(padding%2) + b'"}'
    assert len(data) == size
    add(f"utf8-size-{size}", data, "valid" if size == 65536 else "invalid_response")
for size in [65535,65536,65537]:
    add(f"size-{size}", raw.encode()+b' '*(size-len(raw.encode())), "valid" if size<=65536 else "invalid_response")
add("unknown-version-oversize", b'{"schema_version":"2.0"}'+b' '*65536, "invalid_response")
add("unknown-version-depth", '{"schema_version":"2.0","x":'+'['*32+'0'+']'*32+'}', "invalid_response")
for stamp in ["2024-02-29T12:34:56.123456789Z", "2026-09-20T00:00:00Z", "0000-01-01T00:00:00Z"]:
    mutated("timestamp-"+stamp, "/observed_at", stamp, "valid")
for stamp in ["2025-02-29T00:00:00Z", "2024-13-01T00:00:00Z", "2024-01-32T00:00:00Z", "2024-01-01T24:00:00Z", "2024-01-01T00:00:60Z", "2024-01-01T00:00:00+00:00", "2024-01-01t00:00:00z", "2024-01-01T00:00:00,1Z", "0001-01-01T00:00:00Z"]:
    mutated("bad-timestamp-"+stamp, "/observed_at", stamp)

assert len({case["id"] for case in cases}) == len(cases)
result = (json.dumps({"format": 1, "provenance": "Board-owned S01 cases; documentary snapshot seeds at 5f0d887; Health baseline unchanged", "cases": cases}, indent=2, ensure_ascii=True, allow_nan=False)+"\n").encode()
destination = HERE / "cases.json"
if "--check" in sys.argv:
    assert destination.read_bytes() == result, "Corpus differs from generator output"
else:
    destination.write_bytes(result)
print(f'{len(cases)} cases, {sum(c["classification"] == "valid" for c in cases)} valid, {len(result)} bytes')
