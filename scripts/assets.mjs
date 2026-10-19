import { readFileSync, readdirSync, writeFileSync } from "node:fs";
import { resolve } from "node:path";
import { createHash } from "node:crypto";
import { gzipSync, gunzipSync } from "node:zlib";
import assert from "node:assert/strict";
export const sha256 = (b) => createHash("sha256").update(b).digest("hex");
export function prepareAssets(dist) {
  assert.equal(process.versions.node, "24.21.0");
  assert.equal(process.versions.zlib, "1.3.2.1-motley-8002e91");
  const files = [
    "index.html",
    ...readdirSync(resolve(dist, "assets"))
      .sort()
      .map((n) => "assets/" + n),
  ];
  const resources = files.map((file) => {
    assert.match(file, /^(index\.html|assets\/[A-Za-z0-9_-]+\.(js|css))$/);
    const body = readFileSync(resolve(dist, file));
    const compress = () => {
      const b = gzipSync(body, { level: 9 });
      assert.equal(b[3], 0);
      b.fill(0, 4, 8);
      b[9] = 255;
      return b;
    };
    const gzip = compress();
    assert.deepEqual(gunzipSync(gzip), body);
    assert.deepEqual(compress(), gzip);
    writeFileSync(resolve(dist, file + ".gz"), gzip);
    const representation = (path, bytes) => ({
      path,
      bytes: bytes.length,
      sha256: sha256(bytes),
    });
    return {
      url: file === "index.html" ? "/" : "/" + file,
      mime: file.endsWith(".html")
        ? "text/html; charset=utf-8"
        : file.endsWith(".js")
          ? "text/javascript; charset=utf-8"
          : "text/css; charset=utf-8",
      identity: representation(file, body),
      gzip: representation(file + ".gz", gzip),
    };
  });
  const manifest = {
    schema: 1,
    node: process.versions.node,
    zlib: process.versions.zlib,
    gzip: "level=9,mtime=0,flags=0,os=255",
    resources,
  };
  writeFileSync(
    resolve(dist, "manifest.json"),
    JSON.stringify(manifest, null, 2) + "\n",
  );
  return manifest;
}
