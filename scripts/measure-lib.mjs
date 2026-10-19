import { request } from "node:http";
import { performance } from "node:perf_hooks";
import { gunzipSync } from "node:zlib";
import { sha256 } from "./assets.mjs";
import assert from "node:assert/strict";
export function percentile(values, p = 0.95) {
  if (!values.length || values.some((n) => !Number.isFinite(n) || n < 0))
    throw Error("Invalid sample population");
  return [...values].sort((a, b) => a - b)[Math.ceil(p * values.length) - 1];
}
export function statistics(values) {
  return {
    n: values.length,
    min: Math.min(...values),
    median: percentile(values, 0.5),
    p95: percentile(values),
    max: Math.max(...values),
  };
}
export function rawHTTP(
  url,
  encoding = "identity",
  timeout = 2000,
  maximum = 2 * 1024 * 1024,
) {
  const u = new URL(url);
  assert.equal(u.protocol, "http:");
  assert.equal(u.username + u.password + u.hash, "");
  return new Promise((resolve, reject) => {
    const start = performance.now();
    let timer;
    const req = request(
      u,
      { method: "GET", headers: { "Accept-Encoding": encoding }, agent: false },
      (res) => {
        const parts = [];
        let count = 0;
        res.on("data", (b) => {
          count += b.length;
          if (count > maximum) {
            reject(Error("response limit"));
            req.destroy(Error("response limit"));
            return;
          }
          parts.push(b);
        });
        res.on("error", reject);
        res.on("end", () => {
          clearTimeout(timer);
          resolve({
            status: res.statusCode,
            headers: res.headers,
            body: Buffer.concat(parts),
            elapsed_ms: performance.now() - start,
          });
        });
      },
    );
    timer = setTimeout(() => req.destroy(Error("deadline")), timeout);
    req.on("error", (e) => {
      clearTimeout(timer);
      reject(e);
    });
    req.end();
  });
}
export async function inventory(origin, manifest) {
  const files = [];
  for (const resource of manifest.resources) {
    for (const encoding of ["gzip", "identity"]) {
      const r = await rawHTTP(new URL(resource.url, origin), encoding);
      const expected = resource[encoding];
      assert.equal(r.status, 200);
      assert.equal(
        r.headers["content-encoding"],
        encoding === "gzip" ? "gzip" : undefined,
      );
      assert.equal(r.headers["content-type"], resource.mime);
      assert.equal(r.headers.vary, "Accept-Encoding");
      assert.equal(r.headers["cache-control"], "no-store");
      assert.equal(Number(r.headers["content-length"]), r.body.length);
      assert.equal(sha256(r.body), expected.sha256);
      assert.equal(r.body.length, expected.bytes);
      assert.equal(
        sha256(encoding === "gzip" ? gunzipSync(r.body) : r.body),
        resource.identity.sha256,
      );
      files.push({
        url: resource.url,
        encoding,
        bytes: r.body.length,
        sha256: sha256(r.body),
        body_base64: r.body.toString("base64"),
        headers: r.headers,
      });
    }
  }
  const initial = await rawHTTP(new URL("/", origin), "identity");
  const referenced = [
    ...initial.body.toString().matchAll(/(?:src|href)="(\/assets\/[^"?#]+)"/g),
  ].map((m) => m[1]);
  assert.ok(referenced.length > 0);
  for (const url of referenced)
    assert.ok(
      manifest.resources.some((r) => r.url === url),
      "unlisted initial dependency",
    );
  assert.doesNotMatch(
    initial.body.toString(),
    /(?:src|href)="(?:https?:)?\/\//,
  );
  const overview = await rawHTTP(new URL("/api/v1/overview", origin), "gzip");
  assert.equal(overview.status, 200);
  assert.equal(overview.headers["content-encoding"], undefined);
  const jsCSS = files
    .filter((r) => r.encoding === "gzip" && /\.(js|css)$/.test(r.url))
    .reduce((a, r) => a + r.bytes, 0);
  const total =
    files
      .filter((r) => r.encoding === "gzip")
      .reduce((a, r) => a + r.bytes, 0) + overview.body.length;
  assert.ok(jsCSS <= 250 * 1024);
  assert.ok(total <= 500 * 1024);
  return {
    scope:
      "Actual encoded HTTP bodies; static gzip plus first identity overview; not Pi acceptance",
    files,
    jsCSSGzipBytes: jsCSS,
    initialBytes: total,
    overview: {
      bytes: overview.body.length,
      sha256: sha256(overview.body),
      body_base64: overview.body.toString("base64"),
      headers: overview.headers,
    },
  };
}
export function timingReport(rows, skip = 0) {
  assert.ok(Number.isSafeInteger(skip) && skip >= 0, "invalid warmup count");
  const samples = rows.filter((r) => r.event === "request_timing"),
    summaries = rows.filter((r) => r.event === "timing_summary");
  assert.equal(summaries.length, 1, "one complete process summary required");
  const s = summaries[0];
  assert.deepEqual(Object.keys(s).sort(), [
    "attempts",
    "dropped",
    "event",
    "queued",
    "truncated",
  ]);
  for (const key of ["attempts", "dropped", "queued", "truncated"])
    assert.ok(
      Number.isSafeInteger(s[key]) && s[key] >= 0,
      "invalid summary counter",
    );
  for (const r of samples) {
    assert.deepEqual(Object.keys(r).sort(), [
      "cancelled",
      "complete",
      "eligible",
      "event",
      "flight",
      "handler_ns",
      "network_ns",
      "request",
      "shared",
      "status",
    ]);
    for (const key of ["cancelled", "complete", "eligible", "shared"])
      assert.equal(typeof r[key], "boolean", "invalid diagnostic flag");
    for (const key of [
      "flight",
      "handler_ns",
      "network_ns",
      "request",
      "status",
    ])
      assert.ok(
        Number.isSafeInteger(r[key]) && r[key] >= 0,
        "invalid diagnostic number",
      );
  }
  assert.equal(s.dropped, 0, "dropped timing data");
  assert.equal(s.truncated, 0, "truncated timing data");
  assert.equal(s.attempts, samples.length);
  assert.equal(s.queued, samples.length);
  assert.equal(new Set(samples.map((r) => r.request)).size, samples.length);
  const sorted = [...samples].sort((a, b) => a.request - b.request);
  sorted.forEach((r, i) => assert.equal(r.request, i + 1, "missing sequence"));
  const eligible = samples.filter((r) => r.eligible && r.request > skip);
  for (const r of eligible) {
    assert.ok(
      r.complete &&
        !r.shared &&
        !r.cancelled &&
        r.status === 200 &&
        r.flight > 0,
    );
    assert.ok(r.handler_ns >= r.network_ns && r.network_ns >= 0);
  }
  assert.equal(
    new Set(eligible.map((r) => r.flight)).size,
    eligible.length,
    "multiple samples for unshared flight",
  );
  return {
    scope:
      "Diagnostic tooling result; operator must select warmup/phase and candidate; not physical PASS",
    attempts: samples.length,
    warmupRequests: skip,
    excluded: samples.length - eligible.length,
    overhead_ns: statistics(eligible.map((r) => r.handler_ns - r.network_ns)),
  };
}
export async function load(
  origin,
  {
    clients = 1,
    slots = 120,
    interval = 5000,
    warmup = 120000,
    stagger = false,
    validate = () => {},
  } = {},
) {
  const rows = [],
    base = performance.now(),
    warmSlots = Math.ceil(warmup / interval);
  await Promise.all(
    Array.from({ length: clients }, async (_, client) => {
      for (let slot = 0; slot < slots + warmSlots; slot++) {
        const due =
          base +
          slot * interval +
          (stagger ? (client * interval) / clients : 0);
        await new Promise((r) =>
          setTimeout(r, Math.max(0, due - performance.now())),
        );
        const start = performance.now(),
          row = {
            client,
            slot,
            warmup: slot < warmSlots,
            scheduled_ms: due - base,
            started_ms: start - base,
            late_ms: Math.max(0, start - due),
            utc: new Date().toISOString(),
          };
        if (start >= due + interval) {
          rows.push({ ...row, outcome: "skipped" });
          continue;
        }
        try {
          const response = await rawHTTP(
            new URL("/api/v1/overview", origin),
            "identity",
            2000,
            81920,
          );
          row.wire_ms = response.elapsed_ms;
          row.status = response.status;
          row.bytes = response.body.length;
          const v = validate(response);
          row.state = v?.health?.snapshot_state;
          row.reason = v?.health?.reason;
          row.outcome = response.status === 200 ? "valid" : "http_error";
        } catch {
          row.outcome = "transport_or_contract_error";
        }
        row.elapsed_ms = performance.now() - start; // includes contract validation
        rows.push(row);
      }
    }),
  );
  const measured = rows.filter((r) => !r.warmup),
    valid = measured.filter(
      (r) => r.outcome === "valid" && r.state === "current",
    );
  const groups = {};
  for (const r of measured.filter((r) => r.outcome === "valid")) {
    const key = `${r.state ?? "unknown"}/${r.reason ?? "none"}`;
    (groups[key] ??= []).push(r.elapsed_ms);
  }
  return {
    scope: "External load tool; no physical PASS inferred",
    clients,
    interval,
    warmup,
    slots,
    rows,
    errors: measured.length - valid.length,
    outcomes: Object.fromEntries(
      [...new Set(measured.map((r) => r.outcome))].map((key) => [
        key,
        measured.filter((r) => r.outcome === key).length,
      ]),
    ),
    latencyByStateReason_ms: Object.fromEntries(
      Object.entries(groups).map(([key, values]) => [key, statistics(values)]),
    ),
    latencyByClient_ms: Object.fromEntries(
      Array.from({ length: clients }, (_, client) => {
        const values = valid
          .filter((r) => r.client === client)
          .map((r) => r.elapsed_ms);
        return [client, values.length ? statistics(values) : null];
      }),
    ),
    completeNominal: valid.length === clients * slots,
    latency_ms: valid.length
      ? statistics(valid.map((r) => r.elapsed_ms))
      : null,
  };
}
