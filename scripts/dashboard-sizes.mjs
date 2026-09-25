import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { readdirSync, readFileSync, writeFileSync, mkdirSync } from "node:fs";
import { resolve } from "node:path";
import { gzipSync } from "node:zlib";
import { startBoard, fixture, markup, root } from "./browser-harness.mjs";
const hash = (bytes) => createHash("sha256").update(bytes).digest("hex");
const describe = (bytes) => ({
  bytes: bytes.length,
  gzip9: gzipSync(bytes, { level: 9 }).length,
  sha256: hash(bytes),
});
const board = await startBoard();
try {
  const files = [];
  const names = [
    "index.html",
    ...readdirSync(resolve(root, "web/dist/assets")).map((n) => `assets/${n}`),
  ];
  for (const name of names) {
    assert.match(name, /^(index\.html|assets\/[^/]+\.(js|css))$/);
    const response = await fetch(
      board.origin + (name === "index.html" ? "/" : "/" + name),
    );
    assert.equal(response.status, 200);
    const body = Buffer.from(await response.arrayBuffer());
    assert.deepEqual(body, readFileSync(resolve(root, "web/dist", name)));
    assert.equal(response.headers.get("content-encoding"), null);
    files.push({
      name,
      ...describe(body),
      httpBodyBytes: body.length,
      contentEncoding: "identity (header absent)",
    });
    assert.ok(
      !body.includes(Buffer.from("input_base64")) &&
        !body.includes(Buffer.from("Playwright")),
      "test corpus/tooling absent",
    );
  }
  const scenarios = [];
  for (const [name, body] of [
    ["complete", fixture()],
    ["partial", fixture("partial-documentary")],
    ["maximum-markup", markup()],
  ]) {
    board.state.body = body;
    const response = await fetch(board.origin + "/api/v1/overview");
    assert.equal(response.status, 200);
    assert.equal(response.headers.get("content-encoding"), null);
    const wire = Buffer.from(await response.arrayBuffer());
    assert.ok(wire.includes(Buffer.from('"availability":"available"')));
    scenarios.push({
      name,
      health: describe(body),
      overview: describe(wire),
      httpBodyBytes: wire.length,
      contentEncoding: "identity (header absent)",
    });
  }
  const jsCSS = files
    .filter((f) => /\.(css|js)$/.test(f.name))
    .reduce((n, f) => n + f.gzip9, 0);
  const all = files.reduce((n, f) => n + f.gzip9, 0),
    raw = files.reduce((n, f) => n + f.bytes, 0);
  assert.ok(jsCSS <= 250 * 1024);
  for (const s of scenarios) {
    s.resourcePlusOverviewGzip9 = all + s.overview.gzip9;
    s.resourcePlusOverviewHTTPBytes = raw + s.httpBodyBytes;
    assert.ok(s.resourcePlusOverviewGzip9 <= 500 * 1024);
  }
  const report = {
    scope:
      "S03 dashboard including production parser; independent offline gzip-9, actual HTTP identity; no S04 compression or physical acceptance",
    node: process.versions.node,
    zlib: process.versions.zlib,
    binarySHA256: hash(readFileSync(board.binary)),
    jsCSSGzip9: jsCSS,
    resourcesGzip9: all,
    resourcesHTTPBytes: raw,
    files,
    scenarios,
  };
  mkdirSync(resolve(root, "out"), { recursive: true });
  writeFileSync(
    resolve(root, "out/dashboard-sizes.json"),
    JSON.stringify(report, null, 2) + "\n",
  );
  console.log(JSON.stringify(report, null, 2));
} finally {
  await board.stop();
}
