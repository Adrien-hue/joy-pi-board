import test from "node:test";
import assert from "node:assert/strict";
import { createServer } from "node:http";
import { percentile, timingReport, load, rawHTTP } from "./measure-lib.mjs";
test("S04 tooling: nearest-rank, incomplete/corrupt diagnostic rejection", () => {
  assert.equal(percentile(Array.from({ length: 200 }, (_, i) => i + 1)), 190);
  const sample = {
    event: "request_timing",
    request: 1,
    flight: 3,
    handler_ns: 200,
    network_ns: 150,
    status: 200,
    complete: true,
    shared: false,
    cancelled: false,
    eligible: true,
  };
  const summary = {
    event: "timing_summary",
    attempts: 1,
    queued: 1,
    dropped: 0,
    truncated: 0,
  };
  assert.equal(timingReport([sample, summary]).overhead_ns.p95, 50);
  assert.throws(() => timingReport([sample, summary], 1));
  for (const rows of [
    [sample],
    [sample, { ...summary, dropped: 1 }],
    [sample, { ...summary, truncated: 1 }],
    [sample, { ...summary, attempts: 2 }],
    [{ ...sample, shared: undefined }, summary],
    [{ ...sample, network_ns: -1 }, summary],
    [sample, { ...sample, request: 2 }, { ...summary, attempts: 2, queued: 2 }],
  ])
    assert.throws(() => timingReport(rows));
});
test("S04 tooling: controlled load counts errors and retains raw samples", async () => {
  let calls = 0;
  const server = createServer((req, res) => {
    calls++;
    res.writeHead(calls % 2 ? 200 : 503);
    res.end("{}");
  });
  await new Promise((r) => server.listen(0, "127.0.0.1", r));
  try {
    const origin = "http://127.0.0.1:" + server.address().port;
    const result = await load(origin, {
      slots: 4,
      interval: 20,
      warmup: 0,
      validate: (r) => {
        if (r.status !== 200) throw Error("bad");
        return { health: { snapshot_state: "current" } };
      },
    });
    assert.equal(result.rows.length, 4);
    assert.equal(result.errors, 2);
    assert.equal(result.completeNominal, false);
    assert.ok(result.rows.every((row) => row.elapsed_ms >= 0));
    assert.equal(result.latencyByClient_ms[0].n, 2);
    await assert.rejects(rawHTTP(origin, "identity", 2000, 1));
  } finally {
    await new Promise((r) => server.close(r));
  }
});
