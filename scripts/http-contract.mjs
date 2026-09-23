// Test-only adapter: no production Fetch consumer, polling or UI state.
import assert from "node:assert/strict";
import { writeFileSync } from "node:fs";
import { parse } from "../web/node_modules/lossless-json/lib/esm/index.js";
import {
  parseSnapshot,
  decimalUInt64,
  NumericToken,
} from "../out/s01-parser/contract.mjs";

const [board, control] = process.argv.slice(2);
const normalize = (v) =>
  JSON.stringify(v, (_k, value) =>
    typeof value === "bigint" ? `u64:${decimalUInt64(value)}` : value,
  );
// Only fixture metadata uses JSON.parse; input bytes are carried as base64.
const cases = JSON.parse(await (await fetch(control + "/cases")).text());
const decoder = new TextDecoder("utf-8", { fatal: true });
function depth(value, n = 0) {
  if (value === null || typeof value !== "object" || NumericToken.is(value))
    return n;
  return Math.max(n + 1, ...Object.values(value).map((v) => depth(v, n + 1)));
}
function envelope(bytes) {
  assert.ok(bytes.byteLength <= 81920, "envelope byte limit");
  const result = parse(decoder.decode(bytes), undefined, {
    parseNumber: (token) => new NumericToken(token),
  });
  assert.ok(depth(result) <= 34, "envelope depth limit");
  assert.deepEqual(Object.keys(result), ["generated_at", "health"]);
  assert.deepEqual(Object.keys(result.health), [
    "availability",
    "snapshot_state",
    "last_success_at",
    "reason",
    "snapshot",
  ]);
  return result;
}
let valid = 0,
  invalid = 0,
  maximumBytes = 0,
  maximumDepth = 0;
for (const [index, c] of cases.entries()) {
  assert.equal((await fetch(control + "/" + index)).status, 204);
  const res = await fetch(board + "/api/v1/overview");
  assert.equal(res.status, 200, c.id);
  assert.equal(res.headers.get("content-type"), "application/json");
  assert.equal(res.headers.get("cache-control"), "no-store");
  const bytes = new Uint8Array(await res.arrayBuffer());
  maximumBytes = Math.max(maximumBytes, bytes.length);
  const doc = envelope(bytes),
    h = doc.health;
  maximumDepth = Math.max(maximumDepth, depth(doc));
  if (c.classification === "valid") {
    assert.equal(h.availability, "available", c.id);
    assert.equal(h.snapshot_state, "current", c.id);
    assert.equal(h.reason, null);
    assert.equal(typeof h.snapshot, "object");
    assert.notEqual(h.snapshot, null);
    // This test knows Board's fixed wrapper: snapshot is its final member.
    // Extract its original bytes, avoiding the library stringifier's duck-typed
    // isLosslessNumber/prototype behavior on otherwise valid unknown members.
    // This is a fixed test-envelope slice, not a general JSON parser.
    const text = decoder.decode(bytes);
    const marker = '"snapshot":';
    assert.ok(text.endsWith("}}\n"));
    const extracted = text.slice(text.indexOf(marker) + marker.length, -3);
    const original = Buffer.from(c.input_base64, "base64");
    assert.equal(
      normalize(parseSnapshot(extracted)),
      normalize(parseSnapshot(original)),
      c.id,
    );
    const tokenTree = (v) =>
      NumericToken.is(v)
        ? ["number-token", v.text]
        : Array.isArray(v)
          ? v.map(tokenTree)
          : v !== null && typeof v === "object"
            ? Object.fromEntries(
                Object.entries(v).map(([k, item]) => [k, tokenTree(item)]),
              )
            : v;
    assert.deepEqual(
      tokenTree(
        parse(extracted, undefined, {
          parseNumber: (token) => new NumericToken(token),
        }),
      ),
      tokenTree(h.snapshot),
      c.id,
    );
    assert.ok(bytes.length - Buffer.byteLength(extracted) <= 1024);
    assert.match(
      res.headers.get("x-joy-pi-snapshot-age-ms"),
      /^(0|[1-9][0-9]*)$/,
    );
    if (c.id === "s02-maximum-markup-unicode") {
      assert.equal(original.length, 65536);
      assert.ok(!decoder.decode(bytes).includes("\\u003c"));
      assert.equal(Buffer.byteLength(extracted), 65536);
    }
    valid++;
  } else {
    assert.equal(h.availability, "unavailable", c.id);
    assert.equal(h.reason, c.classification, c.id);
    invalid++;
  }
}
assert.equal(maximumDepth, 34);
assert.throws(() => envelope(new Uint8Array(81921)), /byte limit/);
assert.throws(
  () => envelope(Buffer.from("[".repeat(35) + "0" + "]".repeat(35))),
  /depth limit/,
);
const report = {
  scope:
    "S02 real simulated Health HTTP -> Board HTTP -> S01 TypeScript parser in Node, not browser/dashboard",
  node: process.versions.node,
  cases: cases.length,
  valid,
  invalid,
  healthAttempts: cases.length,
  maximumEnvelopeBytes: maximumBytes,
  maximumEnvelopeDepth: maximumDepth,
  result: "PASS",
};
writeFileSync(
  "out/s02-http-interop.json",
  JSON.stringify(report, null, 2) + "\n",
);
console.log(report);
