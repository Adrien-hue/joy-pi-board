import { ContractError, invalid, NumericToken, parseDocument } from "./json";
import type { ByteGroup, HealthSnapshot, IssueCode } from "./types";

const maximum = 18446744073709551615n;
function object(value: unknown, path: string): Record<string, unknown> {
  if (
    typeof value !== "object" ||
    value === null ||
    Array.isArray(value) ||
    NumericToken.is(value)
  )
    invalid(path, "required object");
  return value as Record<string, unknown>;
}
function field(o: Record<string, unknown>, key: string, path: string): unknown {
  if (!Object.hasOwn(o, key)) invalid(path, "required field absent");
  return o[key];
}
function text(value: unknown, path: string): string {
  if (typeof value !== "string" || value.length === 0)
    invalid(path, "required non-empty string");
  return value;
}
function uint(value: unknown, path: string): bigint | null {
  if (value === null) return null;
  if (!NumericToken.is(value) || !/^(0|[1-9][0-9]*)$/.test(value.text))
    invalid(path, "uint64 decimal token required");
  // Length bound avoids allocating a huge bigint for an already invalid integer.
  if (value.text.length > 20) invalid(path, "uint64 out of range");
  const n = BigInt(value.text);
  if (n > maximum) invalid(path, "uint64 out of range");
  return n;
}
function floating(
  value: unknown,
  path: string,
  min = -Infinity,
  max = Infinity,
): number | null {
  if (value === null) return null;
  if (!NumericToken.is(value)) invalid(path, "number required");
  const n = Number(value.text);
  if (!Number.isFinite(n) || n < min || n > max)
    invalid(path, "finite number outside domain");
  return n;
}
function boolean(value: unknown, path: string): boolean | null {
  if (value !== null && typeof value !== "boolean")
    invalid(path, "boolean required");
  return value;
}
function array(value: unknown, path: string): unknown[] {
  if (!Array.isArray(value)) invalid(path, "required array");
  return value;
}
function timestamp(value: unknown): string {
  const s = text(value, "/observed_at");
  const m =
    /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})(?:\.(\d+))?Z$/.exec(s);
  if (!m) invalid("/observed_at", "RFC3339 UTC ending Z required");
  const year = Number(m[1]),
    month = Number(m[2]),
    day = Number(m[3]);
  const leap = year % 4 === 0 && (year % 100 !== 0 || year % 400 === 0);
  const days = [31, leap ? 29 : 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31];
  if (
    month < 1 ||
    month > 12 ||
    day < 1 ||
    day > days[month - 1]! ||
    Number(m[4]) > 23 ||
    Number(m[5]) > 59 ||
    Number(m[6]) > 59
  )
    invalid("/observed_at", "invalid calendar timestamp");
  // Health cannot encode Go's zero time. Nanosecond interpretation matches Go.
  if (
    s.startsWith("0001-01-01T00:00:00") &&
    !/[1-9]/.test((m[7] ?? "").slice(0, 9))
  )
    invalid("/observed_at", "zero timestamp");
  return s;
}
function allOrNone(path: string, ...values: unknown[]): void {
  if (!values.every((v) => (v === null) === (values[0] === null)))
    invalid(path, "all leaves available or null");
}
function group(root: Record<string, unknown>, key: string): ByteGroup {
  const p = `/${key}`,
    o = object(field(root, key, p), p);
  const total_bytes = uint(
    field(o, "total_bytes", `${p}/total_bytes`),
    `${p}/total_bytes`,
  );
  const available_bytes = uint(
    field(o, "available_bytes", `${p}/available_bytes`),
    `${p}/available_bytes`,
  );
  const used_bytes = uint(
    field(o, "used_bytes", `${p}/used_bytes`),
    `${p}/used_bytes`,
  );
  allOrNone(p, total_bytes, available_bytes, used_bytes);
  if (
    total_bytes !== null &&
    available_bytes !== null &&
    used_bytes !== null &&
    (available_bytes > total_bytes ||
      used_bytes !== total_bytes - available_bytes)
  )
    invalid(p, "inconsistent byte group");
  return { total_bytes, available_bytes, used_bytes };
}
function issueCode(value: unknown, path: string): IssueCode {
  if (
    value !== "unsupported" &&
    value !== "permission_denied" &&
    value !== "temporarily_unavailable"
  )
    invalid(path, "unsupported issue code");
  return value;
}
// UTF-8 lexicographic ordering equals Unicode scalar ordering, NOT UTF-16 order.
function compareNames(a: string, b: string): number {
  const aa = Array.from(a, (s) => s.codePointAt(0)!);
  const bb = Array.from(b, (s) => s.codePointAt(0)!);
  for (let i = 0; i < Math.min(aa.length, bb.length); i++) {
    if (aa[i] !== bb[i]) return aa[i]! - bb[i]!;
  }
  return aa.length - bb.length;
}

export function parseSnapshot(input: Uint8Array | string): HealthSnapshot {
  return validateSnapshot(parseDocument(input));
}

// Private: only receives a tree after the global structural checks above.
function validateSnapshot(value: unknown): HealthSnapshot {
  const root = object(value, "/");
  const version = field(root, "schema_version", "/schema_version");
  if (typeof version !== "string")
    invalid("/schema_version", "required string");
  if (version !== "1.0")
    throw new ContractError(
      "unsupported_schema_version",
      "/schema_version",
      "unsupported version",
    );
  const host = object(field(root, "host", "/host"), "/host");
  const cpu = object(field(root, "cpu", "/cpu"), "/cpu");
  const load = object(field(root, "load", "/load"), "/load");
  const pi = object(
    field(root, "raspberry_pi", "/raspberry_pi"),
    "/raspberry_pi",
  );
  const s: HealthSnapshot = {
    schema_version: version,
    observed_at: timestamp(field(root, "observed_at", "/observed_at")),
    host: {
      hostname: text(
        field(host, "hostname", "/host/hostname"),
        "/host/hostname",
      ),
    },
    cpu: {
      utilization_percent: floating(
        field(cpu, "utilization_percent", "/cpu/utilization_percent"),
        "/cpu/utilization_percent",
        0,
        100,
      ),
      logical_cpu_count: uint(
        field(cpu, "logical_cpu_count", "/cpu/logical_cpu_count"),
        "/cpu/logical_cpu_count",
      ),
    },
    load: {
      one_minute: floating(
        field(load, "one_minute", "/load/one_minute"),
        "/load/one_minute",
        0,
      ),
      five_minutes: floating(
        field(load, "five_minutes", "/load/five_minutes"),
        "/load/five_minutes",
        0,
      ),
      fifteen_minutes: floating(
        field(load, "fifteen_minutes", "/load/fifteen_minutes"),
        "/load/fifteen_minutes",
        0,
      ),
    },
    memory: group(root, "memory"),
    root_filesystem: group(root, "root_filesystem"),
    uptime_seconds: uint(
      field(root, "uptime_seconds", "/uptime_seconds"),
      "/uptime_seconds",
    ),
    network: null,
    raspberry_pi: {
      soc_temperature_celsius: floating(
        field(
          pi,
          "soc_temperature_celsius",
          "/raspberry_pi/soc_temperature_celsius",
        ),
        "/raspberry_pi/soc_temperature_celsius",
      ),
      thermal_throttling_active: boolean(
        field(
          pi,
          "thermal_throttling_active",
          "/raspberry_pi/thermal_throttling_active",
        ),
        "/raspberry_pi/thermal_throttling_active",
      ),
      thermal_throttling_occurred_since_boot: boolean(
        field(
          pi,
          "thermal_throttling_occurred_since_boot",
          "/raspberry_pi/thermal_throttling_occurred_since_boot",
        ),
        "/raspberry_pi/thermal_throttling_occurred_since_boot",
      ),
      undervoltage_active: boolean(
        field(pi, "undervoltage_active", "/raspberry_pi/undervoltage_active"),
        "/raspberry_pi/undervoltage_active",
      ),
      undervoltage_occurred_since_boot: boolean(
        field(
          pi,
          "undervoltage_occurred_since_boot",
          "/raspberry_pi/undervoltage_occurred_since_boot",
        ),
        "/raspberry_pi/undervoltage_occurred_since_boot",
      ),
    },
    issues: [],
  };
  if (s.cpu.logical_cpu_count === 0n)
    invalid("/cpu/logical_cpu_count", "must be positive");
  allOrNone("/load", ...Object.values(s.load));
  const firmware = [
    s.raspberry_pi.thermal_throttling_active,
    s.raspberry_pi.thermal_throttling_occurred_since_boot,
    s.raspberry_pi.undervoltage_active,
    s.raspberry_pi.undervoltage_occurred_since_boot,
  ];
  allOrNone("/raspberry_pi", ...firmware);
  const network = field(root, "network", "/network");
  if (network !== null) {
    const n = object(network, "/network");
    const interfaces = array(
      field(n, "interfaces", "/network/interfaces"),
      "/network/interfaces",
    );
    if (interfaces.length > 64)
      invalid("/network/interfaces", "at most 64 interfaces");
    s.network = {
      interfaces: interfaces.map((entry, i) => {
        const p = `/network/interfaces/${i}`,
          o = object(entry, p),
          state = field(o, "state", `${p}/state`);
        return {
          name: text(field(o, "name", `${p}/name`), `${p}/name`),
          state: state === null ? null : text(state, `${p}/state`),
          rx_bytes: uint(
            field(o, "rx_bytes", `${p}/rx_bytes`),
            `${p}/rx_bytes`,
          ),
          tx_bytes: uint(
            field(o, "tx_bytes", `${p}/tx_bytes`),
            `${p}/tx_bytes`,
          ),
        };
      }),
    };
    for (let i = 1; i < s.network.interfaces.length; i++)
      if (
        compareNames(
          s.network.interfaces[i - 1]!.name,
          s.network.interfaces[i]!.name,
        ) >= 0
      )
        invalid(
          "/network/interfaces",
          "names must be unique and sorted by UTF-8 bytes",
        );
  }
  s.issues = array(field(root, "issues", "/issues"), "/issues").map(
    (entry, i) => {
      const p = `/issues/${i}`,
        o = object(entry, p);
      return {
        path: text(field(o, "path", `${p}/path`), `${p}/path`),
        code: issueCode(field(o, "code", `${p}/code`), `${p}/code`),
        message: text(field(o, "message", `${p}/message`), `${p}/message`),
      };
    },
  );
  const expected: string[] = [];
  let useful = false;
  const add = (v: unknown, path: string) => {
    if (v === null) expected.push(path);
    else useful = true;
  };
  add(s.uptime_seconds, "/uptime_seconds");
  add(s.cpu.logical_cpu_count, "/cpu/logical_cpu_count");
  add(s.load.one_minute, "/load");
  add(s.memory.total_bytes, "/memory");
  add(s.root_filesystem.total_bytes, "/root_filesystem");
  for (const [key, v] of Object.entries(s.raspberry_pi))
    add(v, `/raspberry_pi/${key}`);
  if (s.network === null) expected.push("/network");
  else
    s.network.interfaces.forEach((n, i) => {
      useful = true;
      add(n.state, `/network/interfaces/${i}/state`);
      add(n.rx_bytes, `/network/interfaces/${i}/rx_bytes`);
      add(n.tx_bytes, `/network/interfaces/${i}/tx_bytes`);
    });
  add(s.cpu.utilization_percent, "/cpu/utilization_percent");
  if (s.issues.length !== expected.length)
    invalid("/issues", "issue count does not match unavailable observations");
  s.issues.forEach((issue, i) => {
    if (issue.path !== expected[i])
      invalid(`/issues/${i}/path`, "unexpected issue path or order");
  });
  const firmwareIssues = s.issues.filter(
    (issue) =>
      issue.path.startsWith("/raspberry_pi/") &&
      issue.path !== "/raspberry_pi/soc_temperature_celsius",
  );
  if (
    firmwareIssues.some(
      (issue) =>
        issue.code !== firmwareIssues[0]!.code ||
        issue.message !== firmwareIssues[0]!.message,
    )
  )
    invalid("/issues", "firmware issues must share code and message");
  if (!useful) invalid("/", "no useful observation");
  return s;
}

export function decimalUInt64(value: bigint): string {
  if (typeof value !== "bigint" || value < 0n || value > maximum)
    throw new RangeError("uint64 required");
  return value.toString(10);
}
