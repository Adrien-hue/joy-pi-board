import { spawnSync } from "node:child_process";
import {
  cpSync,
  existsSync,
  mkdirSync,
  mkdtempSync,
  readdirSync,
  readFileSync,
  writeFileSync,
} from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { createHash } from "node:crypto";
import { gzipSync } from "node:zlib";

const root = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const web = join(root, "web");
// Go's ./... also discovers Go files shipped inside npm dependencies.
const goPackages = ["./cmd/...", "./internal/...", "./web"];
const versions = JSON.parse(
  readFileSync(join(root, "toolchains.json"), "utf8"),
);
const env = {
  ...process.env,
  GOENV: "off",
  GOWORK: "off",
  GOTOOLCHAIN: "local",
  GOFLAGS: "",
  CGO_ENABLED: "0",
  GOCACHE: process.env.GOCACHE || join(root, ".cache", "go-build"),
  GOMODCACHE: process.env.GOMODCACHE || join(root, ".cache", "go-mod"),
};
delete env.GOOS;
delete env.GOARCH;

const npmCLI = [
  process.env.npm_execpath,
  join(dirname(process.execPath), "node_modules/npm/bin/npm-cli.js"),
  resolve(dirname(process.execPath), "../lib/node_modules/npm/bin/npm-cli.js"),
].find((p) => p && existsSync(p));

function execute(command, args, options = {}) {
  const result = spawnSync(command, args, {
    cwd: root,
    env,
    stdio: "inherit",
    windowsHide: true,
    ...options,
  });
  if (result.error) throw result.error;
  if (result.status !== 0) {
    throw new Error(`${command} ${args.join(" ")} exited ${result.status}`);
  }
  return result;
}

function capture(command, args, options = {}) {
  return execute(command, args, {
    encoding: "utf8",
    stdio: "pipe",
    ...options,
  }).stdout.trim();
}

function npm(args) {
  execute(process.execPath, [npmCLI, ...args], { cwd: web });
}

function tools() {
  if (!npmCLI)
    throw new Error("Cannot locate npm CLI; use the pinned Node distribution.");
  const actual = {
    go: capture("go", ["env", "GOVERSION"]).replace(/^go/, ""),
    node: process.versions.node,
    npm: capture(process.execPath, [npmCLI, "--version"]),
  };
  for (const name of Object.keys(versions)) {
    if (actual[name] !== versions[name]) {
      throw new Error(
        `${name}: expected ${versions[name]}, got ${actual[name]}; automatic toolchain switching is disabled`,
      );
    }
  }
  console.log(`Exact toolchains: ${JSON.stringify(actual)}; GOTOOLCHAIN=local`);
}

function frontend() {
  execute(
    process.execPath,
    [join(web, "node_modules/vite/bin/vite.js"), "build"],
    { cwd: web },
  );
}

function requireDist() {
  if (!existsSync(join(web, "dist/index.html"))) {
    throw new Error(
      "Missing web/dist: run npm --prefix web run build:frontend before Go compilation.",
    );
  }
}

function lint() {
  execute(process.execPath, [
    join(web, "node_modules/eslint/bin/eslint.js"),
    "--config",
    "web/eslint.config.js",
    "web",
    "scripts",
  ]);
}

function goFiles(dir = root) {
  return readdirSync(dir, { withFileTypes: true }).flatMap((entry) => {
    if (
      [".git", ".tools", ".cache", "out", "node_modules", "dist"].includes(
        entry.name,
      )
    )
      return [];
    const p = join(dir, entry.name);
    return entry.isDirectory()
      ? goFiles(p)
      : entry.name.endsWith(".go")
        ? [p]
        : [];
  });
}

function checkGo() {
  requireDist();
  const unformatted = capture("gofmt", ["-l", ...goFiles()]);
  if (unformatted) throw new Error(`Run gofmt on:\n${unformatted}`);
  execute("go", ["mod", "tidy", "-diff"]);
  execute("go", ["mod", "verify"]);
  execute("go", ["vet", ...goPackages]);
  execute("go", ["test", "-count=1", ...goPackages]);
}

function build(cross) {
  requireDist();
  mkdirSync(join(root, "out"), { recursive: true });
  const name = cross
    ? "joy-pi-board-linux-arm64"
    : process.platform === "win32"
      ? "joy-pi-board.exe"
      : "joy-pi-board";
  execute(
    "go",
    [
      "build",
      "-trimpath",
      "-buildvcs=false",
      "-o",
      join(root, "out", name),
      "./cmd/joy-pi-board",
    ],
    {
      env: cross ? { ...env, GOOS: "linux", GOARCH: "arm64" } : env,
    },
  );
  console.log(`${name} sha256 ${hash(readFileSync(join(root, "out", name)))}`);
}

function hash(bytes) {
  return createHash("sha256").update(bytes).digest("hex");
}

function sizes() {
  requireDist();
  const names = [
    "index.html",
    ...readdirSync(join(web, "dist/assets")).map((n) => `assets/${n}`),
  ];
  const files = names.map((name) => {
    const body = readFileSync(join(web, "dist", name));
    return {
      name,
      bytes: body.length,
      gzip9: gzipSync(body, { level: 9 }).length,
      sha256: hash(body),
    };
  });
  const jsCSS = files
    .filter((f) => /\.(js|css)$/.test(f.name))
    .reduce((n, f) => n + f.gzip9, 0);
  const all = files.reduce((n, f) => n + f.gzip9, 0);
  if (jsCSS > 250 * 1024 || all > 500 * 1024)
    throw new Error("Shell exceeds size budgets");
  const report = {
    scope:
      "S00 shell only, no overview; gzip is measured offline, server serves identity bytes",
    node: process.versions.node,
    zlib: process.versions.zlib,
    jsCSSGzip9: jsCSS,
    allGzip9: all,
    files,
  };
  mkdirSync(join(root, "out"), { recursive: true });
  writeFileSync(
    join(root, "out/shell-sizes.json"),
    `${JSON.stringify(report, null, 2)}\n`,
  );
  console.log(JSON.stringify(report, null, 2));
}

function missingDist() {
  mkdirSync(join(root, "out"), { recursive: true });
  const isolated = mkdtempSync(join(root, "out/missing-dist-"));
  for (const name of ["go.mod", "cmd", "internal"])
    cpSync(join(root, name), join(isolated, name), { recursive: true });
  mkdirSync(join(isolated, "web"));
  cpSync(join(web, "assets.go"), join(isolated, "web/assets.go"));
  const result = spawnSync("go", ["build", "./cmd/joy-pi-board"], {
    cwd: isolated,
    env,
    encoding: "utf8",
    windowsHide: true,
  });
  if (
    result.error ||
    result.status === 0 ||
    !result.stderr.includes("pattern dist: no matching files found")
  ) {
    throw new Error(
      `Missing-dist check failed: ${result.error || result.stderr}`,
    );
  }
  console.log(
    "PASS: clean tree without dist fails at go:embed with a clear missing-files error.",
  );
}

try {
  tools();
  switch (process.argv[2]) {
    case "tools":
      break;
    case "frontend":
      frontend();
      break;
    case "lint":
      lint();
      break;
    case "check":
      capture(process.execPath, [npmCLI, "ls", "--all", "--json"], {
        cwd: web,
      });
      console.log("Installed npm dependency tree is consistent.");
      npm(["run", "format:check"]);
      npm(["run", "typecheck"]);
      lint();
      npm(["test"]);
      execute(process.execPath, ["--test", "scripts/fixtures.test.mjs"]);
      frontend();
      checkGo();
      break;
    case "build":
      frontend();
      build(false);
      break;
    case "cross":
      frontend();
      build(true);
      break;
    case "race":
      frontend();
      execute("go", ["test", "-race", "-count=1", ...goPackages], {
        env: { ...env, CGO_ENABLED: "1" },
      });
      break;
    case "smoke":
      requireDist();
      execute(process.execPath, ["scripts/smoke.mjs"]);
      break;
    case "sizes":
      sizes();
      break;
    case "missing-dist":
      missingDist();
      break;
    case "dev":
      execute(
        process.execPath,
        [join(web, "node_modules/vite/bin/vite.js"), "--host", "127.0.0.1"],
        { cwd: web },
      );
      break;
    default:
      throw new Error(
        "Expected tools, check, frontend, build, cross, race, smoke, sizes, missing-dist, lint or dev",
      );
  }
} catch (error) {
  console.error(error.message);
  process.exitCode = 1;
}
