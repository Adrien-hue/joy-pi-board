import { invalid, NumericToken, parseEnvelopeDocument } from "../health/json";
import { parseSnapshot } from "../health/snapshot";
import type { HealthSnapshot } from "../health/types";

export type Reason =
  | "timeout"
  | "connection_failed"
  | "upstream_error"
  | "invalid_response"
  | "unsupported_schema_version";
export interface Overview {
  generated_at: string;
  health: {
    availability: "available" | "unavailable";
    snapshot_state: "current" | "stale" | "none";
    last_success_at: string | null;
    reason: Reason | null;
    snapshot: HealthSnapshot | null;
  };
  age: number | null;
}
const rootKeys = ["generated_at", "health"];
const healthKeys = [
  "availability",
  "snapshot_state",
  "last_success_at",
  "reason",
  "snapshot",
];
function object(v: unknown, keys: string[]): Record<string, unknown> {
  if (!v || typeof v !== "object" || Array.isArray(v) || NumericToken.is(v))
    invalid("/", "object required");
  if (
    Object.keys(v).length !== keys.length ||
    keys.some((k) => !Object.hasOwn(v, k))
  )
    invalid("/", "wrapper members");
  return v as Record<string, unknown>;
}
function timestamp(v: unknown): string {
  if (typeof v !== "string") invalid("/", "timestamp required");
  const m = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})(?:\.\d+)?Z$/.exec(
    v,
  );
  if (!m) invalid("/", "UTC timestamp required");
  const y = Number(m[1]),
    month = Number(m[2]),
    day = Number(m[3]);
  const days = [
    31,
    y % 4 === 0 && (y % 100 !== 0 || y % 400 === 0) ? 29 : 28,
    31,
    30,
    31,
    30,
    31,
    31,
    30,
    31,
    30,
    31,
  ];
  if (
    month < 1 ||
    month > 12 ||
    day < 1 ||
    day > days[month - 1]! ||
    Number(m[4]) > 23 ||
    Number(m[5]) > 59 ||
    Number(m[6]) > 59
  )
    invalid("/", "calendar timestamp");
  return v;
}

// Only locate one source object AFTER guarded library parsing accepted grammar.
// Strings are skipped atomically; no numbers, values or JSON Pointer are parsed.
function snapshotSource(text: string): string | null {
  const stack: {
    health: boolean;
    key: string;
    start: number;
    snapshot: boolean;
  }[] = [];
  let result: string | null = null;
  for (let i = 0; i < text.length; i++) {
    const c = text[i];
    if (c === '"') {
      const start = i++;
      for (; text[i] !== '"'; i++) if (text[i] === "\\") i++;
      let next = i + 1;
      while (/[\t\r\n ]/.test(text[next] ?? "!")) next++;
      if (text[next] === ":" && stack.length) {
        const frame = stack[stack.length - 1]!;
        frame.key = JSON.parse(text.slice(start, i + 1)) as string;
        // The library may assign __proto__ instead of retaining an own key.
        // Enforce the two strict wrapper key sets on decoded source names too.
        // Health snapshot members retain their separate, tolerant policy.
        if (
          (stack.length === 1 && !rootKeys.includes(frame.key)) ||
          (stack.length === 2 &&
            frame.health &&
            !healthKeys.includes(frame.key))
        )
          invalid("/", "wrapper members");
      }
    } else if (c === "{" || c === "[") {
      const parent = stack[stack.length - 1];
      stack.push({
        health: c === "{" && stack.length === 1 && parent?.key === "health",
        key: "",
        start: i,
        snapshot:
          c === "{" &&
          stack.length === 2 &&
          parent?.health === true &&
          parent.key === "snapshot",
      });
    } else if (c === "}" || c === "]") {
      const frame = stack.pop();
      if (frame?.snapshot) result = text.slice(frame.start, i + 1);
    } else if (c === "," && stack.length) stack[stack.length - 1]!.key = "";
  }
  return result;
}
function freeze<T>(v: T): T {
  if (v && typeof v === "object") {
    Object.values(v).forEach(freeze);
    Object.freeze(v);
  }
  return v;
}
export function parseOverview(
  input: Uint8Array | string,
  ageHeader: string | null,
): Overview {
  const { value, text } = parseEnvelopeDocument(input);
  const root = object(value, rootKeys);
  const h = object(root.health, healthKeys);
  const generated_at = timestamp(root.generated_at);
  const last_success_at =
    h.last_success_at === null ? null : timestamp(h.last_success_at);
  const reasons: unknown[] = [
    "timeout",
    "connection_failed",
    "upstream_error",
    "invalid_response",
    "unsupported_schema_version",
  ];
  if (h.availability !== "available" && h.availability !== "unavailable")
    invalid("/health", "availability");
  if (
    h.snapshot_state !== "current" &&
    h.snapshot_state !== "stale" &&
    h.snapshot_state !== "none"
  )
    invalid("/health", "snapshot state");
  if (h.reason !== null && !reasons.includes(h.reason))
    invalid("/health", "reason");
  let age: number | null = null;
  if (ageHeader !== null) {
    if (
      !/^(0|[1-9][0-9]*)$/.test(ageHeader) ||
      ageHeader.length > 16 ||
      BigInt(ageHeader) > 9007199254740991n
    )
      invalid("/", "age header");
    age = Number(ageHeader);
  }
  if ((age === null) !== (last_success_at === null))
    invalid("/", "age presence");
  const source = snapshotSource(text);
  if (
    new TextEncoder().encode(text).length -
      (source === null ? 0 : new TextEncoder().encode(source).length) >
    1024
  )
    invalid("/", "wrapper overhead");
  const snapshot =
    h.snapshot === null
      ? null
      : source === null
        ? invalid("/", "snapshot source missing")
        : parseSnapshot(source);
  if (h.availability === "available") {
    if (
      h.snapshot_state !== "current" ||
      h.reason !== null ||
      !snapshot ||
      age === null ||
      age > 30000
    )
      invalid("/health", "current invariant");
  } else {
    if (h.reason === null || h.snapshot_state === "current")
      invalid("/health", "failure invariant");
    if (
      h.snapshot_state === "stale" &&
      (!snapshot || age === null || age > 30000)
    )
      invalid("/health", "stale invariant");
    if (
      h.snapshot_state === "none" &&
      (snapshot !== null || (age !== null && age <= 30000))
    )
      invalid("/health", "none invariant");
  }
  return freeze({
    generated_at,
    health: {
      availability: h.availability,
      snapshot_state: h.snapshot_state,
      last_success_at,
      reason: h.reason as Reason | null,
      snapshot,
    },
    age,
  });
}
