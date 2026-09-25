import { parse } from "lossless-json";

export type Classification = "invalid_response" | "unsupported_schema_version";
export class ContractError extends Error {
  constructor(
    public readonly kind: Classification,
    public readonly path: string,
    rule: string,
  ) {
    super(`${kind} at ${path}: ${rule}`);
  }
}
export function invalid(path: string, rule: string): never {
  throw new ContractError("invalid_response", path, rule);
}

export const MAXIMUM_BYTES = 65536;
export const MAXIMUM_DEPTH = 32;

// A private brand cannot be forged by an unknown JSON member or __proto__.
export class NumericToken {
  #text: string;
  constructor(text: string) {
    this.#text = text;
  }
  static is(value: unknown): value is NumericToken {
    return typeof value === "object" && value !== null && #text in value;
  }
  get text(): string {
    return this.#text;
  }
}

function unicode(text: string): void {
  for (let i = 0; i < text.length; i++) {
    const c = text.charCodeAt(i);
    if (c >= 0xdc00 && c <= 0xdfff) invalid("/", "unpaired low surrogate");
    if (c < 0xd800 || c > 0xdbff) continue;
    const low = text.charCodeAt(++i);
    if (!(low >= 0xdc00 && low <= 0xdfff))
      invalid("/", "unpaired high surrogate");
  }
}

// Defensive preflight, NOT a JSON parser. It produces no values and does not
// recognize numbers/keywords or validate grammar. The library owns parsing.
// Native JSON.parse is used ONLY on individual quoted string tokens, never numbers.
function preflight(text: string, maximumDepth: number): void {
  const containers: (Set<string> | null)[] = [];
  for (let i = 0; i < text.length; i++) {
    const c = text[i];
    if (c === '"') {
      const start = i++;
      for (; i < text.length && text[i] !== '"'; i++) {
        if (text[i] === "\\") i++;
      }
      let decoded: unknown;
      try {
        decoded = JSON.parse(text.slice(start, i + 1));
      } catch {
        invalid("/", "malformed JSON string");
      }
      if (typeof decoded !== "string") invalid("/", "string expected");
      unicode(decoded);
      let next = i + 1;
      while (next < text.length && /[ \t\r\n]/.test(text[next]!)) next++;
      if (text[next] === ":") {
        const keys = containers[containers.length - 1];
        if (keys) {
          if (keys.has(decoded)) invalid("/", "duplicate object key");
          keys.add(decoded);
        }
      }
    } else if (c === "{" || c === "[") {
      containers.push(c === "{" ? new Set<string>() : null);
      if (containers.length > maximumDepth)
        invalid("/", "container depth exceeded");
    } else if (c === "}" || c === "]") {
      containers.pop();
    }
  }
}

// Exported for bounded, in-memory JSON composition tests and future consumers.
// Schema validation is separate; no parsed value is trusted by this function.
function guarded(
  input: Uint8Array | string,
  maximumBytes: number,
  maximumDepth: number,
): { value: unknown; text: string } {
  let text: string;
  if (typeof input === "string") {
    if (input.length > maximumBytes)
      invalid("/", "document byte limit exceeded");
    unicode(input); // TextEncoder must not replace ill-formed UTF-16.
    if (new TextEncoder().encode(input).length > maximumBytes)
      invalid("/", "document byte limit exceeded");
    text = input;
  } else {
    if (input.byteLength > maximumBytes)
      invalid("/", "document byte limit exceeded");
    try {
      text = new TextDecoder("utf-8", { fatal: true, ignoreBOM: true }).decode(
        input,
      );
    } catch {
      invalid("/", "malformed UTF-8");
    }
  }
  preflight(text, maximumDepth);
  try {
    const value: unknown = parse(text, undefined, {
      parseNumber: (token) => {
        // Validate raw grammar as well as retaining tokens: custom callbacks
        // must not inadvertently accept a permissive numeric library path.
        if (!/^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?$/.test(token))
          invalid("/", "malformed number");
        return new NumericToken(token);
      },
      onDuplicateKey: () => invalid("/", "duplicate object key"),
    });
    return { value, text };
  } catch {
    // Library messages can echo input. Never expose them as our diagnostics.
    invalid("/", "malformed JSON");
  }
}

export function parseDocument(input: Uint8Array | string): unknown {
  return guarded(input, MAXIMUM_BYTES, MAXIMUM_DEPTH).value;
}
// Two fixed profiles, not caller-controlled relaxation of Health limits.
export function parseEnvelopeDocument(input: Uint8Array | string): {
  value: unknown;
  text: string;
} {
  return guarded(input, 81920, 34);
}
