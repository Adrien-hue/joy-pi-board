import { readFileSync, writeFileSync } from "node:fs";
import { pathToFileURL } from "node:url";
import { resolve } from "node:path";
import { inventory, load, timingReport } from "./measure-lib.mjs";
const [mode, input, output, arg, slots = "120", warmup = "120000"] =
  process.argv.slice(2);
if (!output)
  throw Error(
    "Usage: measure.mjs inventory ORIGIN OUTPUT [MANIFEST] | load ORIGIN OUTPUT [one|ten|sync] [SLOTS] [WARMUP_MS] | timing JOURNAL_JSONL OUTPUT [WARMUP_REQUESTS]",
  );
let result;
if (mode === "inventory")
  result = await inventory(
    input,
    JSON.parse(readFileSync(arg || "web/dist/manifest.json", "utf8")),
  );
else if (mode === "timing")
  result = timingReport(
    readFileSync(input, "utf8")
      .trim()
      .split("\n")
      .map((l) => {
        const r = JSON.parse(l);
        return r.MESSAGE ? JSON.parse(r.MESSAGE) : r;
      }),
    Number(arg || 0),
  );
else if (mode === "load") {
  const profile = arg || "one";
  if (!["one", "ten", "sync"].includes(profile)) throw Error("Unknown profile");
  if (
    !/^[0-9]+$/.test(slots) ||
    !/^[0-9]+$/.test(warmup) ||
    Number(slots) < 1 ||
    Number(slots) > 4096 ||
    Number(warmup) > 1200000
  )
    throw Error("Invalid bounded load configuration");
  const { parseOverview } = await import(
    pathToFileURL(resolve("out/measurement/consumer.mjs"))
  );
  result = await load(input, {
    slots: Number(slots),
    warmup: Number(warmup),
    clients: profile === "one" ? 1 : 10,
    stagger: profile === "ten",
    validate: (r) => {
      if (
        r.headers["content-encoding"] !== undefined ||
        r.headers["content-type"] !== "application/json" ||
        r.headers["cache-control"] !== "no-store"
      )
        throw Error("headers");
      return parseOverview(
        r.body,
        r.headers["x-joy-pi-snapshot-age-ms"] ?? null,
      );
    },
  });
} else throw Error("Unknown measurement command");
writeFileSync(output, JSON.stringify(result, null, 2) + "\n", { flag: "wx" });
console.log("Recorded " + output);
