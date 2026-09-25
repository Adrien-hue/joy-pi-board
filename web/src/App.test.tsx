import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import { App, Dashboard } from "./App";
import { readFileSync } from "node:fs";
import { parseOverview } from "./overview/contract";
describe("S03-T09 dashboard structure", () => {
  it("starts honestly without metrics, presumed service availability or controls", () => {
    const html = renderToStaticMarkup(<App />);
    expect(html).toContain("<h1>Joy Pi Board</h1>");
    expect(html).toContain('role="status"');
    expect(html).toContain("Loading observations");
    expect(html).not.toMatch(/<button|<input|https?:\/\//);
    expect(html).not.toContain("Health available");
  });
  it("renders exact counters and all metric families without fabricated rates", () => {
    const corpus = JSON.parse(
      readFileSync(
        new URL("../../testdata/s01/cases.json", import.meta.url),
        "utf8",
      ),
    ) as { cases: { input_base64: string }[] };
    const s = Buffer.from(corpus.cases[0]!.input_base64, "base64").toString();
    const envelope = parseOverview(
      `{"generated_at":"2026-09-23T00:00:00Z","health":{"availability":"available","snapshot_state":"current","last_success_at":"2026-09-23T00:00:00Z","reason":null,"snapshot":${s}}}`,
      "0",
    );
    const html = renderToStaticMarkup(
      <Dashboard
        view={{
          board: "usable",
          httpStatus: undefined,
          envelope,
          metrics: "current",
          age: 0,
        }}
      />,
    );
    for (const text of [
      "0.5%",
      "9007199254740993",
      "18446744073709551615",
      "Memory",
      "Root filesystem",
      "Temperature",
      "Load averages",
      "Uptime",
      "Active at observation",
      "Occurred since boot",
      "Yes",
      "No",
      "Services &amp; freshness",
    ])
      expect(html).toContain(text);
    expect(html).not.toMatch(/<script|<button|<input|bytes\/s/);
  });
});
