import assert from "node:assert/strict";
import { readFileSync, writeFileSync, mkdirSync } from "node:fs";
import { createHash } from "node:crypto";
import { gzipSync } from "node:zlib";
import { dirname, resolve } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import { build } from "../web/node_modules/vite/dist/node/index.js";
import { parse } from "../web/node_modules/lossless-json/lib/esm/index.js";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const out = resolve(root, "out/s01-parser");
// Exports remain live so the parser cannot disappear as it does from the shell.
await build({
  configFile: false,
  root: resolve(root, "web"),
  build: {
    outDir: out,
    emptyOutDir: true,
    sourcemap: false,
    minify: true,
    target: ["chrome111", "edge111", "firefox115", "safari16.4"],
    lib: {
      entry: resolve(root, "web/contract-entry.ts"),
      formats: ["es"],
      fileName: () => "contract.mjs",
    },
  },
});
const bytes = readFileSync(resolve(out, "contract.mjs"));
const measurement = {
  scope:
    "Dedicated S01 parser + runtime validator + decimal formatter, not shell/dashboard; all public exports retained",
  node: process.versions.node,
  zlib: process.versions.zlib,
  bytes: bytes.length,
  gzip9: gzipSync(bytes, { level: 9 }).length,
  sha256: createHash("sha256").update(bytes).digest("hex"),
};
writeFileSync(
  resolve(root, "out/s01-parser-sizes.json"),
  `${JSON.stringify(measurement, null, 2)}\n`,
);
console.log(measurement);
const { parseSnapshot, parseDocument, NumericToken, decimalUInt64 } =
  await import(pathToFileURL(resolve(out, "contract.mjs")));
const normalize = (value) =>
  JSON.parse(
    JSON.stringify(value, (_key, v) =>
      typeof v === "bigint" ? `u64:${decimalUInt64(v)}` : v,
    ),
  );
const tokenTree = (v) => {
  if (NumericToken.is(v)) return ["number-token", v.text];
  if (Array.isArray(v)) return v.map(tokenTree);
  if (v !== null && typeof v === "object")
    return Object.fromEntries(
      Object.entries(v).map(([k, item]) => [k, tokenTree(item)]),
    );
  return v;
};
const corpus = JSON.parse(
  readFileSync(resolve(root, "testdata/s01/cases.json"), "utf8"),
).cases;
const dir = resolve(root, "out/s01-interop");
const results = JSON.parse(readFileSync(resolve(dir, "results.json"), "utf8"));
assert.equal(results.length, corpus.length);
let valid = 0;
for (const [i, row] of results.entries()) {
  assert.equal(row.id, corpus[i].id);
  assert.equal(row.classification, corpus[i].classification);
  if (row.classification !== "valid") continue;
  const snapshot = readFileSync(resolve(dir, `${row.file}.snapshot.json`));
  assert.deepEqual(
    normalize(parseSnapshot(snapshot)),
    corpus[i].expected,
    row.id,
  );
  assert.deepEqual(normalize(row.expected), corpus[i].expected, row.id);
  // Test-generated envelope adds one container around a valid depth-32 snapshot.
  // Snapshot limits must not be repurposed as the future S02 envelope policy.
  const envelope = parse(
    readFileSync(resolve(dir, `${row.file}.envelope.json`), "utf8"),
    undefined,
    { parseNumber: (token) => new NumericToken(token) },
  );
  assert.equal(typeof envelope.snapshot, "object");
  assert.notEqual(envelope.snapshot, null);
  assert.deepEqual(
    tokenTree(envelope.snapshot),
    tokenTree(parseDocument(snapshot)),
    row.id,
  );
  valid++;
}
const mutations = JSON.parse(
  readFileSync(resolve(dir, "mutations.json"), "utf8"),
);
let mutationValid = 0;
for (const c of mutations) {
  let kind = "valid",
    actual;
  try {
    actual = normalize(parseSnapshot(Buffer.from(c.input_base64, "base64")));
  } catch (e) {
    kind = e.kind;
  }
  assert.equal(kind, c.classification, c.id);
  if (kind === "valid") {
    assert.deepEqual(actual, normalize(c.expected), c.id);
    mutationValid++;
  }
}
const evidence = {
  corpus: corpus.length,
  acceptedGoOutputsParsedInTS: valid,
  seededMutations: mutations.length,
  seed: 20260920,
  acceptedMutations: mutationValid,
  result: "PASS",
  scope: "in-memory contract and test-only envelope, no HTTP/cache/dashboard",
};
mkdirSync(resolve(root, "out"), { recursive: true });
writeFileSync(
  resolve(root, "out/s01-interop.json"),
  `${JSON.stringify(evidence, null, 2)}\n`,
);
console.log(evidence);
