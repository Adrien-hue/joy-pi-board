import { decimalUInt64 } from "../health/snapshot";
export const unavailable = "Unavailable";
const units = ["B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"];
export function bytes(value: bigint | null): string {
  if (value === null) return unavailable;
  decimalUInt64(value);
  let k = 0,
    d = 1n;
  while (k < 6 && value >= d * 1024n) {
    k++;
    d *= 1024n;
  }
  if (k === 0) return `${value} B`;
  let q = (2n * value * 10n + d) / (2n * d);
  if (q >= 10240n && k < 6) {
    k++;
    d *= 1024n;
    q = (2n * value * 10n + d) / (2n * d);
  }
  return `${q / 10n}.${q % 10n} ${units[k]}`;
}
export function percent(used: bigint | null, total: bigint | null): string {
  if (used === null || total === null) return unavailable;
  decimalUInt64(used);
  decimalUInt64(total);
  if (total === 0n) return "Not applicable";
  const q = (2n * used * 1000n + total) / (2n * total);
  return `${q / 10n}.${q % 10n}%`;
}
export function floating(
  value: number | null,
  digits: number,
  suffix = "",
): string {
  if (value === null) return unavailable;
  if (!Number.isFinite(value)) throw new RangeError("finite metric required");
  const n = Object.is(value, -0) ? 0 : value;
  return (
    new Intl.NumberFormat("en-GB", {
      minimumFractionDigits: digits,
      maximumFractionDigits: digits,
      notation: Math.abs(n) >= 1e12 ? "scientific" : "standard",
      useGrouping: false,
    }).format(n) + suffix
  );
}
export function uptime(value: bigint | null): string {
  if (value === null) return unavailable;
  decimalUInt64(value);
  return `${value / 86400n} d ${(value % 86400n) / 3600n} h ${(value % 3600n) / 60n} m ${value % 60n} s`;
}
export function flag(value: boolean | null): string {
  return value === null ? unavailable : value ? "Yes" : "No";
}
