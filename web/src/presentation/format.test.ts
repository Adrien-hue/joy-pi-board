import { expect, it } from "vitest";
import { bytes, percent, floating, uptime, flag } from "./format";
it("S03-T08 preserves uint64 extrema, binary carries and ratios", () => {
  expect(bytes(0n)).toBe("0 B");
  expect(bytes(1024n)).toBe("1.0 KiB");
  expect(bytes(1048575n)).toBe("1.0 MiB");
  expect(bytes(18446744073709551615n)).toBe("16.0 EiB");
  expect(percent(1n, 16n)).toBe("6.3%");
  expect(percent(18446744073709551614n, 18446744073709551615n)).toBe("100.0%");
  expect(percent(0n, 0n)).toBe("Not applicable");
  expect(percent(null, 1n)).toBe("Unavailable");
  expect(uptime(18446744073709551615n)).toBe("213503982334601 d 7 h 0 m 15 s");
  expect(uptime(0n)).toBe("0 d 0 h 0 m 0 s");
});
it("S03-T08 keeps percentage semantics, finite float text, false and null", () => {
  expect(floating(0.5037, 1, "%")).toBe("0.5%");
  expect(floating(-0, 1)).toBe("0.0");
  expect(floating(Number.MAX_VALUE, 2)).not.toMatch(/Infinity|NaN/);
  expect(flag(false)).toBe("No");
  expect(flag(true)).toBe("Yes");
  expect(flag(null)).toBe("Unavailable");
  expect(bytes(null)).toBe("Unavailable");
});
