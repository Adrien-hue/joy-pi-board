import assert from "node:assert/strict";
import { spawn } from "node:child_process";
import { createHash } from "node:crypto";
import {
  copyFileSync,
  mkdirSync,
  mkdtempSync,
  readFileSync,
  writeFileSync,
} from "node:fs";
import { createServer } from "node:net";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { setTimeout as delay } from "node:timers/promises";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const name = process.platform === "win32" ? "joy-pi-board.exe" : "joy-pi-board";
mkdirSync(join(root, "out"), { recursive: true });
const directory = mkdtempSync(join(root, "out/standalone-"));
const executable = join(directory, name);
copyFileSync(join(root, "out", name), executable);
const sha256 = (data) => createHash("sha256").update(data).digest("hex");
let healthConnections = 0;
const trap = createServer((socket) => {
  healthConnections++;
  socket.destroy();
});
let trapBound = false;
await new Promise((done) => {
  trap.once("error", (error) => {
    if (error.code !== "EADDRINUSE") throw error;
    done();
  });
  trap.listen(8080, "127.0.0.1", () => {
    trapBound = true;
    done();
  });
});
const child = spawn(executable, ["--listen", "127.0.0.1:0"], {
  cwd: directory,
  env: { ...process.env, PATH: "" },
  stdio: ["ignore", "pipe", "pipe"],
  windowsHide: true,
});
let logs = "";
child.stderr.on("data", (data) => {
  logs += data.toString();
});
let spawnError;
child.on("error", (error) => {
  spawnError = error;
});
const exited = new Promise((done) => child.once("close", done));
try {
  let origin;
  for (let i = 0; i < 100; i++) {
    if (spawnError) throw spawnError;
    const match = logs.match(/S00 shell listening on (127\.0\.0\.1:\d+)/);
    if (match) {
      origin = `http://${match[1]}`;
      break;
    }
    if (child.exitCode !== null) throw new Error(`Binary exited: ${logs}`);
    await delay(50);
  }
  assert.ok(origin, `Binary did not start: ${logs}`);
  const request = (path, options = {}) =>
    fetch(origin + path, { ...options, signal: AbortSignal.timeout(2000) });
  const index = await request("/");
  assert.equal(index.status, 200);
  assert.match(
    index.headers.get("content-security-policy"),
    /script-src 'self'/,
  );
  const body = Buffer.from(await index.arrayBuffer());
  assert.deepEqual(body, readFileSync(join(root, "web/dist/index.html")));
  const html = body.toString("utf8");
  const assets = [...html.matchAll(/(?:src|href)="(\/assets\/[^"?#]+)"/g)].map(
    (m) => m[1],
  );
  assert.ok(
    assets.some((p) => p.endsWith(".js")) &&
      assets.some((p) => p.endsWith(".css")),
  );
  assert.doesNotMatch(html, /(?:src|href)="(?:https?:)?\/\//);
  const inventory = [];
  for (const path of assets) {
    const response = await request(path);
    assert.equal(response.status, 200);
    const bytes = Buffer.from(await response.arrayBuffer());
    assert.deepEqual(
      bytes,
      readFileSync(join(root, "web/dist", path.slice(1))),
    );
    inventory.push({ path, bytes: bytes.length, sha256: sha256(bytes) });
  }
  assert.equal((await request("/api/v1/overview")).status, 404);
  assert.equal((await request("/not-a-page")).status, 404);
  assert.equal((await request("/", { method: "HEAD" })).status, 200);
  await delay(6000);
  assert.equal(healthConnections, 0);
  const report = {
    scope:
      "S00 native binary, empty working directory except binary, PATH empty; not hardware acceptance",
    platform: process.platform,
    arch: process.arch,
    binarySHA256: sha256(readFileSync(executable)),
    healthTrap: trapBound
      ? "127.0.0.1:8080, zero connections through startup/requests/6 s idle"
      : "NOT RUN: port already occupied; static source review still applies",
    assets: inventory,
  };
  writeFileSync(
    join(root, "out/standalone-smoke.json"),
    `${JSON.stringify(report, null, 2)}\n`,
  );
  console.log(JSON.stringify(report, null, 2));
} finally {
  if (child.exitCode === null) child.kill();
  await exited;
  if (trapBound) await new Promise((done) => trap.close(done));
}
