import { afterEach, describe, expect, it, vi } from "vitest";
import { OverviewController, type Clock, type Runtime } from "./controller";
import { parseOverview } from "./contract";
import { RequestFailure } from "./request";
import { readFileSync } from "node:fs";

const corpus = JSON.parse(
  readFileSync(
    new URL("../../../testdata/s01/cases.json", import.meta.url),
    "utf8",
  ),
) as { cases: { id: string; input_base64: string }[] };
function overview(age = 0, id = "complete") {
  const s = Buffer.from(
    corpus.cases.find((c) => c.id === id)!.input_base64,
    "base64",
  )
    .toString()
    .trim();
  return parseOverview(
    `{"generated_at":"2026-09-23T00:00:00Z","health":{"availability":"available","snapshot_state":"current","last_success_at":"2026-09-23T00:00:00Z","reason":null,"snapshot":${s}}}`,
    String(age),
  );
}
function setup() {
  vi.useFakeTimers();
  let time: Clock = { mono: 0, wall: 100000 };
  const pending: {
    signal: AbortSignal;
    received: () => void;
    resolve: (v: ReturnType<typeof overview>) => void;
    reject: (e: unknown) => void;
  }[] = [];
  const runtime: Runtime = {
    now: () => ({ ...time }),
    later: (fn, ms) => setTimeout(fn, ms),
    cancel: clearTimeout,
    microtask: queueMicrotask,
    request: (signal, received) =>
      new Promise((resolve, reject) =>
        pending.push({ signal, received, resolve, reject }),
      ),
  };
  const changed = vi.fn(),
    mask = vi.fn(),
    c = new OverviewController(runtime, changed, mask);
  const tick = async (ms: number) => {
    time = { mono: time.mono + ms, wall: time.wall + ms };
    await vi.advanceTimersByTimeAsync(ms);
  };
  return {
    c,
    pending,
    mask,
    changed,
    tick,
    setTime: (v: Clock) => {
      time = v;
    },
    time: () => time,
  };
}
afterEach(() => vi.useRealTimers());
const settle = async () => {
  await Promise.resolve();
  await Promise.resolve();
  await Promise.resolve();
};
describe("S03-T04/T05 scheduling and lifecycle", () => {
  it("survives StrictMode replay with one initial launch, fixed slots and no retry", async () => {
    const { c, pending, tick } = setup();
    c.start(true);
    c.stop();
    c.start(true);
    await settle();
    expect(pending).toHaveLength(1);
    pending[0]!.resolve(overview());
    await settle();
    await tick(4999);
    expect(pending).toHaveLength(1);
    await tick(1);
    expect(pending).toHaveLength(2);
    pending[1]!.reject(new RequestFailure("network"));
    await settle();
    expect(pending).toHaveLength(2);
    await tick(5000);
    expect(pending).toHaveLength(3);
    c.stop();
  });
  it.each([1999, 2000, 2001])(
    "checks synchronous completion deadline at %i",
    async (ms) => {
      const { c, pending, setTime } = setup();
      c.start(true);
      await settle();
      setTime({ mono: ms, wall: 100000 + ms });
      pending[0]!.resolve(overview());
      await settle();
      expect(c.view().board).toBe(ms < 2000 ? "usable" : "timeout");
      c.stop();
    },
  );
  it("holds active permit through aborted cleanup and discards old results", async () => {
    const { c, pending, tick } = setup();
    c.start(true);
    await settle();
    await tick(2000);
    expect(pending[0]!.signal.aborted).toBe(true);
    expect(c.view().board).toBe("timeout");
    await tick(8000);
    expect(pending).toHaveLength(1);
    pending[0]!.resolve(overview());
    await settle();
    expect(c.view().board).toBe("timeout");
    await tick(5000);
    expect(pending).toHaveLength(2);
    c.stop();
  });
  it("hides, suspends, resumes only at next slot and rejects pre-hide completion", async () => {
    const { c, pending, tick, mask } = setup();
    c.start(true);
    await settle();
    c.pause();
    c.pause();
    expect(mask).toHaveBeenCalled();
    pending[0]!.resolve(overview());
    await settle();
    await tick(12000);
    expect(pending).toHaveLength(1);
    c.resume();
    c.resume();
    await tick(2999);
    expect(pending).toHaveLength(1);
    await tick(1);
    expect(pending).toHaveLength(2);
    expect(c.view().metrics).toBe("none");
    pending[1]!.resolve(overview());
    await settle();
    expect(c.view().metrics).toBe("current");
    c.stop();
  });
  it("hidden mount has no request and unmount clears timers", async () => {
    const { c, pending, tick } = setup();
    c.start(false);
    await tick(10000);
    expect(pending).toHaveLength(0);
    c.resume();
    await tick(5000);
    expect(pending).toHaveLength(1);
    c.stop();
    expect(vi.getTimerCount()).toBe(0);
    expect(pending[0]!.signal.aborted).toBe(true);
  });
});
describe("S03-T06/T07 conservative age and atomic states", () => {
  it.each([29999, 30000, 30001])(
    "guards use at age %i independently of timers",
    async (age) => {
      const { c, pending, setTime } = setup();
      c.start(true);
      await settle();
      setTime({ mono: 10, wall: 100010 });
      pending[0]!.resolve(overview(29000));
      await settle();
      setTime({ mono: age - 29000, wall: 100000 + age - 29000 });
      expect(c.view().metrics).toBe(age <= 30000 ? "current" : "none");
      c.stop();
    },
  );
  it("expires with no new response and never refreshes an old anchor on failure", async () => {
    const { c, pending, tick } = setup();
    c.start(true);
    await settle();
    pending[0]!.resolve(overview(29999));
    await settle();
    await tick(1);
    expect(c.view().metrics).toBe("current");
    await tick(1);
    expect(c.view().metrics).toBe("none");
    expect(pending).toHaveLength(1);
    c.stop();
  });
  it.each([NaN, -1])(
    "invalidates unreliable monotonic reads %s",
    async (bad) => {
      const { c, pending, setTime } = setup();
      c.start(true);
      await settle();
      pending[0]!.resolve(overview());
      await settle();
      setTime({ mono: bad, wall: 100001 });
      expect(c.view().metrics).toBe("none");
      c.stop();
    },
  );
  it.each([131000, 99999])(
    "invalidates forward expiry or backward wall reading %i without Pi ordering",
    async (wall) => {
      const { c, pending, setTime } = setup();
      c.start(true);
      await settle();
      pending[0]!.resolve(overview());
      await settle();
      setTime({ mono: 1, wall });
      expect(c.view().metrics).toBe("none");
      c.stop();
    },
  );
  it("retains eligible stale metrics on HTTP failure but clears on invalid contract and none", async () => {
    const { c, pending, tick } = setup();
    c.start(true);
    await settle();
    pending[0]!.resolve(overview());
    await settle();
    await tick(5000);
    pending[1]!.reject(new RequestFailure("http", 503));
    await settle();
    expect(c.view()).toMatchObject({
      board: "http",
      httpStatus: 503,
      metrics: "stale",
    });
    await tick(5000);
    pending[2]!.reject(new RequestFailure("contract"));
    await settle();
    expect(c.view().metrics).toBe("none");
    await tick(5000);
    pending[3]!.resolve(
      parseOverview(
        '{"generated_at":"2026-09-23T00:00:00Z","health":{"availability":"unavailable","snapshot_state":"none","last_success_at":null,"reason":"timeout","snapshot":null}}',
        null,
      ),
    );
    await settle();
    expect(c.view()).toMatchObject({ board: "usable", metrics: "none" });
    expect(c.view().envelope!.health.snapshot).toBeNull();
    c.stop();
  });
});

it("S03-T06 accumulates paired segments across body receipt and synchronous validation", async () => {
  const { c, pending, setTime } = setup();
  c.start(true);
  await settle();
  setTime({ mono: 100, wall: 100200 });
  pending[0]!.received();
  setTime({ mono: 300, wall: 100250 });
  pending[0]!.resolve(overview(29600));
  await settle();
  expect(c.view().age).toBe(30000);
  expect(c.view().metrics).toBe("current");
  setTime({ mono: 301, wall: 100251 });
  expect(c.view().metrics).toBe("none");
  c.stop();
});
