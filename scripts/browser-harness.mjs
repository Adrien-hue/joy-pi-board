// Test-only controlled provider. No fixture or control route enters the binary.
import { createServer } from "node:http";
import { spawn } from "node:child_process";
import { readFileSync, mkdirSync, mkdtempSync } from "node:fs";
import { resolve } from "node:path";
import { fileURLToPath } from "node:url";
export const root = fileURLToPath(new URL("../", import.meta.url));
export const corpus = JSON.parse(
  readFileSync(resolve(root, "testdata/s01/cases.json"), "utf8"),
).cases;
export function fixture(id = "complete") {
  const c = corpus.find((c) => c.id === id);
  if (!c) throw Error(`Unknown fixture ${id}`);
  return Buffer.from(c.input_base64, "base64");
}
export function markup() {
  const prefix = fixture().toString().trim().slice(0, -1) + ',"extra":"';
  const remaining = 65536 - Buffer.byteLength(prefix) - 2;
  return Buffer.from(
    prefix +
      "<>&".repeat(Math.floor(remaining / 3)) +
      " ".repeat(remaining % 3) +
      '"}',
  );
}
export async function startBoard() {
  const state = { body: fixture(), status: 200, calls: 0 };
  const health = createServer((_req, res) => {
    state.calls++;
    res.writeHead(state.status, { "Content-Type": "application/json" });
    res.end(state.body);
  });
  await new Promise((done) => health.listen(0, "127.0.0.1", done));
  mkdirSync(resolve(root, "out"), { recursive: true });
  const cwd = mkdtempSync(resolve(root, "out/browser-standalone-"));
  const binary =
    process.env.S03_BINARY ||
    resolve(
      root,
      process.platform === "win32"
        ? "out/joy-pi-board.exe"
        : "out/joy-pi-board",
    );
  const child = spawn(
    binary,
    [
      "--listen",
      "127.0.0.1:0",
      "--health-url",
      `http://127.0.0.1:${health.address().port}/v1/snapshot`,
    ],
    {
      cwd,
      env: { ...process.env, PATH: "" },
      stdio: ["ignore", "ignore", "pipe"],
      windowsHide: true,
    },
  );
  const exited = new Promise((done) => child.once("close", done));
  let logs = "";
  const origin = await new Promise((resolve, reject) => {
    const timer = setTimeout(
      () => reject(Error("Board readiness deadline")),
      5000,
    );
    child.once("error", (error) => {
      clearTimeout(timer);
      reject(error);
    });
    child.stderr.on("data", (data) => {
      logs += data;
      const match = logs.match(
        /"event":"listening"[^\n]*"address":"(127\.0\.0\.1:\d+)"/,
      );
      if (match) {
        clearTimeout(timer);
        resolve(`http://${match[1]}`);
      }
    });
  });
  return {
    origin,
    state,
    binary,
    stop: async () => {
      if (child.exitCode === null) child.kill();
      await exited;
      await new Promise((done) => health.close(done));
    },
    kill: async () => {
      child.kill();
      await exited;
    },
  };
}
