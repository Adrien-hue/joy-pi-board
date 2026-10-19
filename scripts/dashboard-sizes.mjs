import { readFileSync, writeFileSync, mkdirSync } from "node:fs";
import { resolve } from "node:path";
import { startBoard, fixture, markup, root } from "./browser-harness.mjs";
import { inventory } from "./measure-lib.mjs";
import { sha256 } from "./assets.mjs";
const board = await startBoard();
try {
  const manifest = JSON.parse(
    readFileSync(resolve(root, "web/dist/manifest.json"), "utf8"),
  );
  const scenarios = [];
  for (const [name, body] of [
    ["complete", fixture()],
    ["partial", fixture("partial-documentary")],
    ["maximum-markup", markup()],
  ]) {
    board.state.body = body;
    scenarios.push({
      name,
      healthSHA256: sha256(body),
      ...(await inventory(board.origin, manifest)),
    });
  }
  const report = {
    scope:
      "S04 actual gzip static and identity first overview; controlled workstation, not Pi acceptance",
    binarySHA256: sha256(readFileSync(board.binary)),
    node: process.versions.node,
    zlib: process.versions.zlib,
    scenarios,
  };
  mkdirSync(resolve(root, "out"), { recursive: true });
  writeFileSync(
    resolve(root, "out/dashboard-sizes.json"),
    JSON.stringify(report, null, 2) + "\n",
  );
  console.log(
    JSON.stringify(
      {
        binarySHA256: report.binarySHA256,
        scenarios: scenarios.map((s) => ({
          name: s.name,
          jsCSSGzipBytes: s.jsCSSGzipBytes,
          initialBytes: s.initialBytes,
        })),
      },
      null,
      2,
    ),
  );
} finally {
  await board.stop();
}
