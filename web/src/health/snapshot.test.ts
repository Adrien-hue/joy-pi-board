import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";
import { ContractError, parseDocument } from "./json";
import { decimalUInt64, parseSnapshot } from "./snapshot";

interface Case {
  id: string;
  input_base64: string;
  classification: string;
  expected?: unknown;
}
const corpus: { cases: Case[] } = JSON.parse(
  readFileSync(
    new URL("../../../testdata/s01/cases.json", import.meta.url),
    "utf8",
  ),
);
function normalize(value: unknown): unknown {
  return JSON.parse(
    JSON.stringify(value, (_key, v: unknown) =>
      typeof v === "bigint" ? `u64:${decimalUInt64(v)}` : v,
    ),
  );
}
describe("B-01–B-04 shared Health 1.0 corpus", () => {
  for (const c of corpus.cases)
    it(c.id, () => {
      const bytes = Buffer.from(c.input_base64, "base64");
      if (c.classification === "valid") {
        const s = parseSnapshot(bytes);
        expect(normalize(s)).toEqual(c.expected);
        expect(normalize(parseSnapshot(bytes.toString("utf8")))).toEqual(
          c.expected,
        );
      } else {
        expect(() => parseSnapshot(bytes)).toThrow(ContractError);
        try {
          parseSnapshot(bytes);
        } catch (e) {
          expect(e).toBeInstanceOf(ContractError);
          expect((e as ContractError).kind).toBe(c.classification);
        }
      }
    });
});
it("keeps all uint64 values bigint and integer-looking floats number", () => {
  const c = corpus.cases.find((c) => c.id === "zero-and-false")!;
  const s = parseSnapshot(Buffer.from(c.input_base64, "base64"));
  expect(s.uptime_seconds).toBe(0n);
  expect(s.cpu.logical_cpu_count).toBe(4n);
  expect(s.cpu.utilization_percent).toBe(0);
  expect(s.raspberry_pi.soc_temperature_celsius).toBe(0);
  expect(s.raspberry_pi.undervoltage_active).toBe(false);
});
it("renders exact decimal integers without Number or global prototype changes", () => {
  for (const n of [
    0n,
    4n,
    9007199254740991n,
    9007199254740992n,
    9007199254740993n,
    18446744073709551615n,
  ])
    expect(decimalUInt64(n)).toBe(n.toString());
  expect(() => decimalUInt64(-1n)).toThrow(RangeError);
  expect(() => decimalUInt64(18446744073709551616n)).toThrow(RangeError);
  expect(Object.hasOwn(BigInt.prototype, "toJSON")).toBe(false);
});
it("retains the sign of a floating negative zero while rejecting uint64 -0", () => {
  const c = corpus.cases.find((c) => c.id === "float-cpu--0")!;
  expect(
    Object.is(
      parseSnapshot(Buffer.from(c.input_base64, "base64")).cpu
        .utilization_percent,
      -0,
    ),
  ).toBe(true);
});
it("rejects ill-formed UTF-16 input before TextEncoder could replace it", () => {
  for (const s of ["\ud800", "\udfff", "x\ud800y"])
    expect(() => parseDocument(`{"x":"${s}"}`)).toThrow(ContractError);
});
it("never echoes payload fragments in errors", () => {
  try {
    parseSnapshot('{"schema_version":"2.0","secret":"sensitive-fragment"');
  } catch (e) {
    expect(String(e)).not.toContain("sensitive-fragment");
  }
});
it("does not allow prototype members to satisfy known-field presence", () => {
  const c = corpus.cases[0]!;
  const bytes = Buffer.from(c.input_base64, "base64").toString("utf8");
  const hostname = parseSnapshot(bytes).host.hostname;
  const modified = bytes.replace(
    `"hostname":${JSON.stringify(hostname)}`,
    '"__proto__":{"hostname":"inherited"}',
  );
  expect(() => parseSnapshot(modified)).toThrow(ContractError);
  expect(Object.hasOwn(Object.prototype, "hostname")).toBe(false);
});
