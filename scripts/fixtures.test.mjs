import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { readFileSync } from "node:fs";
import { test } from "node:test";

const numeric = new URL("../testdata/numeric/", import.meta.url);
const manifest = JSON.parse(
  readFileSync(new URL("manifest.json", numeric), "utf8"),
);

test("Health reference fixture retains the exact pinned Git blob", () => {
  const bytes = readFileSync(new URL("health-large-integers.json", numeric));
  const blob = createHash("sha1")
    .update(`blob ${bytes.length}\0`)
    .update(bytes)
    .digest("hex");
  assert.equal(blob, "88e7b1bb4443635e9421e1cf62b55819be8e707d");
});

test("future numeric cases retain their original token bytes without Number parsing", () => {
  const ids = new Set();
  for (const entry of manifest.cases) {
    assert.ok(!ids.has(entry.id));
    ids.add(entry.id);
    assert.match(entry.file, /^[a-z0-9-]+\.json$/);
    assert.equal(
      readFileSync(new URL(entry.file, numeric), "utf8"),
      `${entry.token}\n`,
    );
  }
  assert.ok(
    ids.has("above-two-to-53") && ids.has("uint64-max") && ids.has("null"),
  );
});
