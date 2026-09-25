import { parseOverview, type Overview } from "./contract";
export type FailureKind = "network" | "http" | "timeout" | "contract";
export class RequestFailure extends Error {
  constructor(
    public readonly kind: FailureKind,
    public readonly status?: number,
  ) {
    super(kind);
  }
}
export async function requestOverview(
  signal: AbortSignal,
  fetcher: typeof fetch = fetch,
  received: () => void = () => {},
): Promise<Overview> {
  let response: Response;
  try {
    response = await fetcher("/api/v1/overview", {
      method: "GET",
      mode: "same-origin",
      credentials: "same-origin",
      cache: "no-store",
      redirect: "error",
      signal,
    });
  } catch {
    throw new RequestFailure("network");
  }
  const reader = response.body?.getReader();
  let complete = false;
  const abort = () => {
    void reader?.cancel().catch(() => undefined);
  };
  signal.addEventListener("abort", abort, { once: true });
  try {
    if (response.status !== 200)
      throw new RequestFailure("http", response.status);
    const mime = response.headers.get("content-type") ?? "";
    if (
      !/^application\/json(?:\s*;\s*charset\s*=\s*(?:utf-8|"utf-8"))?$/i.test(
        mime,
      ) ||
      !response.headers
        .get("cache-control")
        ?.split(",")
        .some((v) => v.trim().toLowerCase() === "no-store")
    )
      throw new RequestFailure("contract");
    const encoding = response.headers.get("content-encoding");
    if (encoding !== null && encoding.toLowerCase() !== "identity")
      throw new RequestFailure("contract");
    const length = response.headers.get("content-length");
    if (
      length !== null &&
      (!/^(0|[1-9][0-9]*)$/.test(length) ||
        length.length > 5 ||
        Number(length) > 81920)
    )
      throw new RequestFailure("contract");
    if (!reader) throw new RequestFailure("contract");
    const chunks: Uint8Array[] = [];
    let size = 0;
    while (true) {
      if (signal.aborted) throw new RequestFailure("network");
      const chunk = await reader.read();
      if (chunk.done) break;
      size += chunk.value.byteLength;
      if (size > 81920) throw new RequestFailure("contract");
      chunks.push(chunk.value);
    }
    received();
    if (signal.aborted) throw new RequestFailure("network");
    if (length !== null && size !== Number(length))
      throw new RequestFailure("contract");
    const bytes = new Uint8Array(size);
    let offset = 0;
    for (const chunk of chunks) {
      bytes.set(chunk, offset);
      offset += chunk.length;
    }
    const result = parseOverview(
      bytes,
      response.headers.get("x-joy-pi-snapshot-age-ms"),
    );
    complete = true;
    return result;
  } catch (e) {
    if (e instanceof RequestFailure) throw e;
    throw new RequestFailure("contract");
  } finally {
    signal.removeEventListener("abort", abort);
    if (!complete) {
      try {
        await reader?.cancel();
      } catch {
        /* safe diagnostic only */
      }
    }
    reader?.releaseLock();
  }
}
