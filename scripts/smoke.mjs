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
await new Promise((done) => trap.listen(0, "127.0.0.1", done));
const healthURL = `http://127.0.0.1:${trap.address().port}/v1/snapshot`;
const child = spawn(
  executable,
  ["--listen", "127.0.0.1:0", "--health-url", healthURL],
  {
    cwd: directory,
    env: { ...process.env, PATH: "" },
    stdio: ["ignore", "pipe", "pipe"],
    windowsHide: true,
  },
);
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
    const match = logs.match(
      /"event":"listening"[^\n]*"address":"(127\.0\.0\.1:\d+)"/,
    );
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
  assert.equal(
    healthConnections,
    0,
    "startup and assets must not contact Health",
  );
  const overview = await request("/api/v1/overview");
  assert.equal(overview.status, 200);
  assert.equal(overview.headers.get("content-type"), "application/json");
  assert.equal(overview.headers.get("cache-control"), "no-store");
  // No numbers exist in this failure envelope; this is not an exact-integer oracle.
  const unavailable = JSON.parse(await overview.text());
  assert.deepEqual(unavailable.health, {
    availability: "unavailable",
    snapshot_state: "none",
    last_success_at: null,
    reason: "connection_failed",
    snapshot: null,
  });
  assert.equal(healthConnections, 1);
  assert.equal((await request("/not-a-page")).status, 404);
  assert.equal((await request("/", { method: "HEAD" })).status, 200);
  await delay(6000);
  assert.equal(
    healthConnections,
    1,
    "no background retry during six seconds of idle",
  );
  const report = {
    scope:
      "S02 native binary, embedded S00 shell plus unavailable Health overview, empty working directory, PATH empty; not hardware acceptance",
    platform: process.platform,
    arch: process.arch,
    binarySHA256: sha256(readFileSync(executable)),
    healthTrap:
      "ephemeral loopback reset server: zero startup/asset calls, one overview attempt, no calls during 6 s idle",
    healthConnections,
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
  await new Promise((done) => trap.close(done));
}
