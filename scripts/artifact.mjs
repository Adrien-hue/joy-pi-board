// Builds new, never-overwritten artifact identities; packaging never invokes this.
import { execFileSync } from "node:child_process";
import {
  readFileSync,
  writeFileSync,
  mkdirSync,
  existsSync,
  readdirSync,
  copyFileSync,
} from "node:fs";
import { resolve, join } from "node:path";
import { sha256 } from "./assets.mjs";
const [mode, id, ciPath, packageVersion = "0.1.0-1"] = process.argv.slice(2);
if (
  !["development", "candidate"].includes(mode) ||
  !id ||
  !/^[a-z0-9][a-z0-9.-]{0,63}$/.test(id)
)
  throw Error(
    "Usage: artifact.mjs development ID | candidate ID VERIFIED_CI_JSON PACKAGE_VERSION",
  );
const root = resolve("."),
  dir = join(root, "out/artifacts", id);
if (existsSync(dir)) throw Error("Identity already exists; choose a new ID");
const env = {
  ...process.env,
  GOENV: "off",
  GOWORK: "off",
  GOTOOLCHAIN: "local",
  CGO_ENABLED: "0",
  GOFLAGS: "",
  GOCACHE: resolve(".cache/go-build"),
  GOMODCACHE: resolve(".cache/go-mod"),
};
const run = (exe, args, options = {}) =>
  execFileSync(exe, args, {
    cwd: root,
    env,
    encoding: "utf8",
    stdio: "pipe",
    windowsHide: true,
    ...options,
  });
const revision = run("git", ["rev-parse", "HEAD"]).trim(),
  status = run("git", ["status", "--porcelain=v1", "--untracked-files=all"]);
let ci = null;
if (mode === "candidate") {
  if (!process.argv[5] || !/^0\.1\.0-[1-9][0-9]*$/.test(packageVersion))
    throw Error("Candidate requires explicit unique Debian revision 0.1.0-N");
  const parent = join(root, "out/artifacts");
  if (existsSync(parent))
    for (const other of readdirSync(parent)) {
      const prior = join(parent, other, "provenance.json");
      if (existsSync(prior)) {
        const m = JSON.parse(readFileSync(prior, "utf8"));
        if (
          m.status === "FROZEN CANDIDATE / NOT ACCEPTED" &&
          m.packageVersion === packageVersion
        )
          throw Error(
            "Debian revision already frozen; choose a new revision and ID",
          );
      }
    }
  if (status) throw Error("Candidate requires a clean committed tree");
  if (!ciPath) throw Error("Verified exact-source CI record required");
  ci = JSON.parse(readFileSync(ciPath, "utf8"));
  if (
    ci.sha !== revision ||
    ci.conclusion !== "success" ||
    ci.verified !== true ||
    !/^https:\/\/github.com\/Adrien-hue\/joy-pi-board\/actions\/runs\/[0-9]+$/.test(
      ci.url,
    )
  )
    throw Error("Invalid exact-source CI record");
}
run(process.execPath, ["scripts/tasks.mjs", "tools"], { stdio: "inherit" });
const names = run("git", [
  "ls-files",
  "--cached",
  "--others",
  "--exclude-standard",
  "-z",
])
  .split("\0")
  .filter(Boolean)
  .filter(
    (p) =>
      (!p.startsWith("docs/") || p.startsWith("docs/examples/")) &&
      (!p.endsWith(".md") || p.startsWith("packaging/")),
  )
  .sort();
const inputs = () =>
  names.map((path) => ({ path, sha256: sha256(readFileSync(path)) }));
const before = inputs(),
  epoch = Number(run("git", ["show", "-s", "--format=%ct", "HEAD"]).trim());
mkdirSync(dir, { recursive: true });
writeFileSync(join(dir, "INCOMPLETE"), "Build not finalized\n", { flag: "wx" });
run(process.execPath, ["scripts/tasks.mjs", "frontend"], { stdio: "inherit" });
const assets = JSON.parse(readFileSync("web/dist/manifest.json", "utf8"));
const built = new Date(epoch * 1000).toISOString(),
  version = "0.1.0",
  binary = join(dir, "joy-pi-board-linux-arm64");
const flags = [
  "-trimpath",
  "-buildvcs=false",
  "-ldflags",
  `-X main.version=${version} -X main.revision=${revision} -X main.buildTime=${built}`,
];
run("go", ["build", ...flags, "-o", binary, "./cmd/joy-pi-board"], {
  env: { ...env, GOOS: "linux", GOARCH: "arm64", GOARM64: "v8.0" },
  stdio: "inherit",
});
const after = inputs();
if (JSON.stringify(before) !== JSON.stringify(after))
  throw Error("Inputs changed during build");
const b = readFileSync(binary);
if (
  b.readUInt32BE(0) !== 0x7f454c46 ||
  b[4] !== 2 ||
  b.readUInt16LE(18) !== 183
)
  throw Error("Not ELF64 ARM64");
const header = b.readBigUInt64LE(32),
  size = b.readUInt16LE(54),
  count = b.readUInt16LE(56);
for (let i = 0; i < count; i++)
  if (b.readUInt32LE(Number(header) + i * size) === 3)
    throw Error("Unexpected dynamic interpreter");
const manifest = {
  schema: 1,
  id,
  status:
    mode === "development"
      ? "DEVELOPMENT / NOT ACCEPTED"
      : "FROZEN CANDIDATE / NOT ACCEPTED",
  version,
  packageVersion,
  revision,
  dirty: !!status,
  gitStatus: status,
  ci,
  sourceDateEpoch: epoch,
  buildMetadataTime: built,
  actualBuildTime: new Date().toISOString(),
  tools: JSON.parse(readFileSync("toolchains.json", "utf8")),
  flags,
  target: { GOOS: "linux", GOARCH: "arm64", GOARM64: "v8.0", CGO_ENABLED: "0" },
  inputs: before,
  inputDigest: sha256(
    Buffer.from(before.map((x) => x.sha256 + "  " + x.path + "\n").join("")),
  ),
  assets,
  binary: {
    file: "joy-pi-board-linux-arm64",
    bytes: b.length,
    sha256: sha256(b),
  },
  goBuildInfo: run("go", ["version", "-m", binary]),
};
copyFileSync("web/dist/manifest.json", join(dir, "assets.json"));
const notices = [
  "Board distribution license is not established by these third-party notices.\n",
];
for (const name of ["react", "react-dom", "scheduler", "lossless-json"]) {
  const location = join("web/node_modules", name),
    pkg = JSON.parse(readFileSync(join(location, "package.json"), "utf8"));
  const license = existsSync(join(location, "LICENSE"))
    ? "LICENSE"
    : "LICENSE.md";
  notices.push(
    name +
      " " +
      pkg.version +
      "\n" +
      readFileSync(join(location, license), "utf8"),
  );
}
notices.push(
  "Go runtime\n" +
    readFileSync(join(run("go", ["env", "GOROOT"]).trim(), "LICENSE"), "utf8"),
);
writeFileSync(join(dir, "THIRD-PARTY-NOTICES"), notices.join("\n\n"), {
  flag: "wx",
});
writeFileSync(
  join(dir, "provenance.json"),
  JSON.stringify(manifest, null, 2) + "\n",
  { flag: "wx" },
);
writeFileSync(
  join(dir, "SHA256SUMS"),
  [
    "joy-pi-board-linux-arm64",
    "assets.json",
    "provenance.json",
    "THIRD-PARTY-NOTICES",
  ]
    .map((p) => sha256(readFileSync(join(dir, p))) + "  " + p + "\n")
    .join(""),
  { flag: "wx" },
);
// Completion marker is append-only: INCOMPLETE remains history, COMPLETE is authoritative.
writeFileSync(
  join(dir, "COMPLETE"),
  "Immutable selected artifact; package/campaign acceptance still required\n",
  { flag: "wx" },
);
console.log(
  JSON.stringify({
    directory: dir,
    status: manifest.status,
    sha256: manifest.binary.sha256,
  }),
);
