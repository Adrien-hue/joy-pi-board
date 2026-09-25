import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";
import { parseOverview } from "./contract";
import { requestOverview, RequestFailure } from "./request";

const corpus: {
  cases: { id: string; input_base64: string; classification: string }[];
} = JSON.parse(
  readFileSync(
    new URL("../../../testdata/s01/cases.json", import.meta.url),
    "utf8",
  ),
);
const time = '"2026-09-23T12:00:00Z"';
const sample = Buffer.from(corpus.cases[0]!.input_base64, "base64").toString();
function wrap(s: string) {
  return `{"generated_at":${time},"health":{"availability":"available","snapshot_state":"current","last_success_at":${time},"reason":null,"snapshot":${s}}}`;
}
describe("S03-T02/T03 original snapshot and shared corpus", () => {
  it.each(["__proto__", "__pro\\u0074o__", "constructor", "prototype"])(
    "rejects extra decoded wrapper key %s regardless of library projection",
    (key) => {
      for (const value of ["null", "0", "{}", '{"reason":null}']) {
        const extra = `"${key}":${value},`;
        expect(() =>
          parseOverview(wrap(sample).replace("{", "{" + extra), "0"),
        ).toThrow();
        expect(() =>
          parseOverview(
            wrap(sample).replace('"health":{', '"health":{' + extra),
            "0",
          ),
        ).toThrow();
      }
      expect(() =>
        parseOverview(
          wrap(sample).replace('"reason":null', `"${key}":{"reason":null}`),
          "0",
        ),
      ).toThrow();
    },
  );
  for (const c of corpus.cases)
    it(c.id, () => {
      const bytes = Buffer.concat([
        Buffer.from(wrap("").split('"snapshot":')[0]! + '"snapshot":'),
        c.classification === "valid"
          ? Buffer.from(Buffer.from(c.input_base64, "base64").toString().trim())
          : Buffer.from(c.input_base64, "base64"),
        Buffer.from("}}"),
      ]);
      if (c.classification === "valid")
        expect(parseOverview(bytes, "0").health.snapshot).not.toBeNull();
      else expect(() => parseOverview(bytes, "0")).toThrow();
    });
  it("is independent of escaped property order and marker strings", () => {
    for (const position of [0, 1, 2, 3, 4]) {
      const fields = [
        '"availability":"available"',
        '"snapshot_state":"current"',
        `"last_success_at":${time}`,
        '"reason":null',
      ];
      fields.splice(
        position,
        0,
        `"snap\\u0073hot":${sample.slice(0, -1)},"extra":{"snapshot":"} \\"snapshot\\": {","__proto__":{"isLosslessNumber":true},"health":{"snapshot":[]}}}`,
      );
      const o = parseOverview(
        `{"he\\u0061lth":{${fields.join(",")}},"generated_at":${time}}`,
        "1",
      );
      expect(o.health.snapshot!.network!.interfaces[0]!.rx_bytes).toBe(
        9007199254740993n,
      );
      expect(o.health.snapshot!.network!.interfaces[0]!.tx_bytes).toBe(
        18446744073709551615n,
      );
      expect(Object.isFrozen(o.health.snapshot!.network!.interfaces)).toBe(
        true,
      );
    }
  });
  it("keeps maximum markup bytes and wrapper boundary independent", () => {
    const prefix = sample.slice(0, -1) + ',"extra":"';
    const s =
      prefix +
      "<>&".repeat(Math.floor((65536 - prefix.length - 2) / 3)) +
      " ".repeat((65536 - prefix.length - 2) % 3) +
      '"}';
    expect(Buffer.byteLength(s)).toBe(65536);
    const body = wrap(s);
    const overhead = body.length - s.length;
    expect(
      parseOverview(body + " ".repeat(1024 - overhead), "0").health.snapshot,
    ).not.toBeNull();
    expect(() =>
      parseOverview(body + " ".repeat(1025 - overhead), "0"),
    ).toThrow();
    expect(() =>
      parseOverview(wrap(s.replace('"extra":"', '"extra":" ')), "0"),
    ).toThrow();
  });
  it.each([
    null,
    "-0",
    "00",
    "1.0",
    "1e1",
    " 1",
    "1, 1",
    "30001",
    "9007199254740992",
  ])("rejects inconsistent age %s", (age) =>
    expect(() => parseOverview(wrap(sample), age)).toThrow(),
  );
  it("accepts all four states, rejects extras and decoded duplicates", () => {
    expect(parseOverview(wrap(sample), "30000").age).toBe(30000);
    expect(
      parseOverview(
        wrap(sample)
          .replace('"available"', '"unavailable"')
          .replace('"current"', '"stale"')
          .replace('"reason":null', '"reason":"timeout"'),
        "2",
      ).health.snapshot_state,
    ).toBe("stale");
    const none = wrap("null")
      .replace('"available"', '"unavailable"')
      .replace('"current"', '"none"')
      .replace('"reason":null', '"reason":"timeout"');
    expect(parseOverview(none, "9007199254740991").health.snapshot).toBeNull();
    expect(
      parseOverview(
        none.replace(`"last_success_at":${time}`, '"last_success_at":null'),
        null,
      ).age,
    ).toBeNull();
    for (const bad of [
      wrap(sample).replace("{", '{"x":0,'),
      wrap(sample).replace('"health":{', '"health":{"x":0,'),
      wrap(sample).replace(
        '"reason":null',
        '"reason":null,"rea\\u0073on":null',
      ),
      wrap(sample) + "null",
      wrap(sample).replace("2026-09-23", "2026-02-30"),
    ])
      expect(() => parseOverview(bad, "0")).toThrow();
  });
});
describe("S03-T01 bounded Fetch", () => {
  const headers = {
    "content-type": "application/json",
    "cache-control": "no-store",
    "x-joy-pi-snapshot-age-ms": "0",
  };
  it("reads split UTF-8 and enforces same-origin options", async () => {
    const bytes = new TextEncoder().encode(
      wrap(sample.replace("joy-pi", "été")),
    );
    const response = new Response(
      new ReadableStream({
        start(c) {
          for (const b of bytes) c.enqueue(Uint8Array.of(b));
          c.close();
        },
      }),
      { headers },
    );
    const result = await requestOverview(
      new AbortController().signal,
      async (url, opts) => {
        expect(url).toBe("/api/v1/overview");
        expect(opts).toMatchObject({
          cache: "no-store",
          redirect: "error",
          mode: "same-origin",
          method: "GET",
        });
        return response;
      },
    );
    expect(result.health.snapshot!.host.hostname).toBe("été");
  });
  it.each([81919, 81920, 81921])("counts stream bytes at %i", async (size) => {
    let cancelled = false;
    const response = new Response(
      new ReadableStream({
        start(c) {
          c.enqueue(new Uint8Array(size));
          if (size <= 81920) c.close();
        },
        cancel() {
          cancelled = true;
        },
      }),
      { headers },
    );
    const signal = new AbortController();
    await expect(
      requestOverview(signal.signal, async () => response),
    ).rejects.toBeInstanceOf(RequestFailure);
    expect(cancelled).toBe(size > 81920);
  });
  it.each([
    { "content-type": "text/html" },
    { "content-type": "application/json;charset=utf-8;charset=utf-8" },
    { "cache-control": "max-age=5" },
    { "content-encoding": "gzip" },
    { "content-length": "81921" },
    { "content-length": "01" },
    { "content-length": "1" },
  ])("rejects unsafe headers %j", async (extra) => {
    await expect(
      requestOverview(
        new AbortController().signal,
        async () =>
          new Response(wrap(sample), { headers: { ...headers, ...extra } }),
      ),
    ).rejects.toMatchObject({ kind: "contract" });
  });
  it("separates HTTP, network and malformed contract errors", async () => {
    await expect(
      requestOverview(
        new AbortController().signal,
        async () => new Response("private", { status: 503 }),
      ),
    ).rejects.toMatchObject({ kind: "http", status: 503 });
    await expect(
      requestOverview(new AbortController().signal, async () => {
        throw Error("private");
      }),
    ).rejects.toMatchObject({ kind: "network" });
    for (const bytes of [
      Uint8Array.of(0xff),
      new TextEncoder().encode(wrap(sample) + "x"),
      Uint8Array.of(0xef, 0xbb, 0xbf),
    ])
      await expect(
        requestOverview(
          new AbortController().signal,
          async () => new Response(bytes, { headers }),
        ),
      ).rejects.toMatchObject({ kind: "contract" });
  });
});

it("S03-T02 required fields, forbidden nulls and Gregorian Board timestamps", () => {
  const fields = [
    `"availability":"available"`,
    `"snapshot_state":"current"`,
    `"last_success_at":${time}`,
    `"reason":null`,
    `"snapshot":${sample}`,
  ];
  for (let i = 0; i < fields.length; i++) {
    const missing = `{"generated_at":${time},"health":{${fields.filter((_, j) => i !== j).join(",")}}}`;
    expect(() => parseOverview(missing, "0")).toThrow();
    if (i !== 3) {
      const replaced = [...fields];
      replaced[i] = fields[i]!.slice(0, fields[i]!.indexOf(":") + 1) + "null";
      expect(() =>
        parseOverview(
          `{"generated_at":${time},"health":{${replaced.join(",")}}}`,
          "0",
        ),
      ).toThrow();
    }
  }
  for (const body of [
    "{}",
    "null",
    "[]",
    `{"health":{${fields.join(",")}}}`,
    `{"generated_at":null,"health":{${fields.join(",")}}}`,
    `{"generated_at":${time},"health":null}`,
  ])
    expect(() => parseOverview(body, "0")).toThrow();
  for (const valid of [
    "0000-02-29T00:00:00Z",
    "0001-01-01T00:00:00Z",
    "9999-12-31T23:59:59.123Z",
  ])
    expect(
      parseOverview(wrap(sample).replace(time, JSON.stringify(valid)), "0")
        .generated_at,
    ).toBe(valid);
  for (const bad of [
    "2026-02-29T00:00:00Z",
    "2026-01-01T24:00:00Z",
    "2026-01-01T00:00:60Z",
    "2026-01-01T00:00:00+00:00",
  ])
    expect(() =>
      parseOverview(wrap(sample).replace(time, JSON.stringify(bad)), "0"),
    ).toThrow();
});
it("S03-T01 cancels a hanging body when aborted and releases its reader", async () => {
  let cancelled = false;
  const body = new ReadableStream<Uint8Array>({
    cancel() {
      cancelled = true;
    },
  });
  const abort = new AbortController();
  const promise = requestOverview(
    abort.signal,
    async () =>
      new Response(body, {
        headers: {
          "content-type": "application/json",
          "cache-control": "no-store",
        },
      }),
  );
  await Promise.resolve();
  abort.abort();
  await expect(promise).rejects.toMatchObject({ kind: "network" });
  expect(cancelled).toBe(true);
  expect(body.locked).toBe(false);
});
