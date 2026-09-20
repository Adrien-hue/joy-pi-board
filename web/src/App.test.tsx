import { renderToStaticMarkup } from "react-dom/server";
import { describe, expect, it } from "vitest";
import { App } from "./App";

describe("S00 shell", () => {
  it("renders an accessible, honest foundation without provider controls or data", () => {
    const html = renderToStaticMarkup(<App />);
    expect(html).toContain("<h1>Joy Pi Board</h1>");
    expect(html).toContain('aria-labelledby="foundation-title"');
    expect(html).toContain("Frontend ready");
    expect(html).toContain("No health data is");
    expect(html).not.toMatch(/<button|<input|<table|https?:\/\//);
  });
});
