import { test, expect } from "../web/node_modules/@playwright/test/index.mjs";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { createHash } from "node:crypto";
import {
  startBoard,
  fixture,
  markup,
  corpus,
  root,
} from "./browser-harness.mjs";

let server;
test.beforeEach(async ({ browser }, info) => {
  server = await startBoard();
  info.annotations.push(
    { type: "engine", description: browser.version() },
    {
      type: "platform",
      description: `${process.platform}/${process.arch}; Node ${process.versions.node}`,
    },
    {
      type: "binary-sha256",
      description: createHash("sha256")
        .update(readFileSync(server.binary))
        .digest("hex"),
    },
  );
});
test.afterEach(async () => {
  await server?.stop();
});
async function current(page) {
  await expect(
    page.getByText("Current observations", { exact: true }),
  ).toBeVisible({ timeout: 7000 });
}
async function capture(page, info, name) {
  await page.screenshot({
    path: info.outputPath(name + ".png"),
    fullPage: true,
  });
}
test("S03-T09/T10/T12 real HTTP, exact values, responsive keyboard and local CSP", async ({
  page,
}, info) => {
  const external = [],
    errors = [],
    requests = [];
  page.on("request", (r) => {
    requests.push(r.url());
    if (!r.url().startsWith(server.origin)) external.push(r.url());
  });
  page.on("pageerror", (e) => errors.push(e.message));
  await page.route("**/*", (route) =>
    route.request().url().startsWith(server.origin)
      ? route.continue()
      : route.abort(),
  );
  await page.addInitScript(() => {
    window.cspViolations = [];
    document.addEventListener("securitypolicyviolation", (e) =>
      window.cspViolations.push(e.violatedDirective),
    );
  });
  const response = await page.goto(server.origin);
  expect(response.headers()["content-security-policy"]).toContain(
    "connect-src 'self'",
  );
  await current(page);
  await expect(
    page.getByText("9007199254740993", { exact: true }),
  ).toBeVisible();
  await expect(
    page.getByText("18446744073709551615", { exact: true }),
  ).toBeVisible();
  await expect(page.getByText("0.5%", { exact: true })).toBeVisible();
  await page.getByText("Exact utilization", { exact: true }).focus();
  await page.keyboard.press("Enter");
  await expect(page.getByText("0.5037%", { exact: true })).toBeVisible();
  const before = await page.locator("details[open]").count();
  await expect
    .poll(() => server.state.calls, { timeout: 7000 })
    .toBeGreaterThan(1);
  expect(await page.locator("details[open]").count()).toBe(before);
  expect(
    requests.filter((p) => p.endsWith("/api/v1/overview")).length,
  ).toBeGreaterThanOrEqual(server.state.calls);
  expect(await page.evaluate(() => window.cspViolations)).toEqual([]);
  expect(external).toEqual([]);
  expect(errors).toEqual([]);
  // Playwright WebKit screenshot synchronization injects a blocked style tag.
  // Check production CSP before capture instrumentation, without relaxing CSP.
  for (const width of [360, 390, 768, 1280]) {
    await page.setViewportSize({ width, height: 900 });
    expect(
      await page.evaluate(
        () => document.documentElement.scrollWidth <= innerWidth,
      ),
    ).toBe(true);
    await capture(page, info, "nominal-" + width);
  }
  // Text enlargement is supplemental, not a claim of native browser zoom.
  await page.evaluate(() => {
    document.documentElement.style.fontSize = "32px";
  });
  await page.setViewportSize({ width: 360, height: 900 });
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= innerWidth,
    ),
  ).toBe(true);
  await capture(page, info, "text-200-percent");
});
test("S03-T10 real partial, Health failure, recovery and maximum original markup", async ({
  page,
}, info) => {
  await page.goto(server.origin);
  await current(page);
  const partial =
    corpus.find(
      (c) =>
        c.classification === "valid" &&
        c.id.includes("memory") &&
        c.id.includes("permission"),
    ) ||
    corpus.find(
      (c) => c.classification === "valid" && c.id.startsWith("partial"),
    );
  expect(partial).toBeTruthy();
  const memoryExact = page
    .getByRole("region", { name: "Memory", exact: true })
    .getByText("Exact used", { exact: true });
  await memoryExact.focus();
  await page.keyboard.press("Enter");
  server.state.body = fixture(partial.id);
  await expect(
    page.getByText("Current observations · Some measures unavailable", {
      exact: true,
    }),
  ).toBeVisible({ timeout: 7000 });
  await expect(memoryExact).toBeFocused();
  await expect(
    page
      .getByRole("region", { name: "Memory", exact: true })
      .locator("details[open]"),
  ).toHaveCount(1);
  await capture(page, info, "partial");
  server.state.status = 503;
  await expect(
    page.getByText("Stale observations", { exact: true }),
  ).toBeVisible({ timeout: 7000 });
  await capture(page, info, "stale");
  server.state.status = 200;
  server.state.body = markup();
  await current(page);
  await expect(
    page.getByText("18446744073709551615", { exact: true }),
  ).toBeVisible();
  await capture(page, info, "recovered-maximum");
});
test("S03-T06/T10 independent real expiry after Board stops", async ({
  page,
}, info) => {
  await page.goto(server.origin);
  await current(page);
  await server.kill();
  await expect(
    page.getByText("Board could not be reached").first(),
  ).toBeVisible({ timeout: 8000 });
  await expect(
    page.getByText("Stale observations", { exact: true }),
  ).toBeVisible();
  await capture(page, info, "board-unreachable-stale");
  await expect(
    page.getByText("9007199254740993", { exact: true }),
  ).not.toBeVisible({ timeout: 35000 });
  await expect(page.locator("h1")).toHaveText("Joy Pi Board");
  await capture(page, info, "expired-without-server");
});
test("S03-T05 synthetic lifecycle rejects old freshness and waits for a future slot", async ({
  page,
}) => {
  await page.goto(server.origin);
  await current(page);
  await page.evaluate(() =>
    window.dispatchEvent(
      new PageTransitionEvent("pagehide", { persisted: true }),
    ),
  );
  await expect(
    page.getByText("9007199254740993", { exact: true }),
  ).not.toBeVisible();
  const count = server.state.calls;
  await page.evaluate(() => {
    window.dispatchEvent(
      new PageTransitionEvent("pageshow", { persisted: true }),
    );
    window.dispatchEvent(
      new PageTransitionEvent("pageshow", { persisted: true }),
    );
  });
  expect(server.state.calls).toBe(count);
  await current(page);
});
test("S03-T02/T07 intercepted invalid Board envelope masks immediately; HTTP differs from network", async ({
  page,
}, info) => {
  await page.goto(server.origin);
  await current(page);
  await page.route("**/api/v1/overview", (route) =>
    route.fulfill({
      status: 200,
      headers: {
        "content-type": "application/json",
        "cache-control": "no-store",
      },
      body: "{}",
    }),
  );
  await expect(
    page.getByText("Board returned invalid data").first(),
  ).toBeVisible({ timeout: 7000 });
  await expect(
    page.getByText("9007199254740993", { exact: true }),
  ).not.toBeVisible();
  await capture(page, info, "invalid-board");
  await page.unroute("**/api/v1/overview");
  await page.route("**/api/v1/overview", (route) =>
    route.fulfill({ status: 503, body: "private error" }),
  );
  await expect(page.getByText("Board returned HTTP 503").first()).toBeVisible({
    timeout: 7000,
  });
  await expect(page.getByText("private error")).toHaveCount(0);
  await capture(page, info, "board-http");
});

test("S03-T03/T10 shared valid corpus through Health HTTP, Board and production browser", async ({
  page,
}) => {
  test.setTimeout(120000);
  for (const c of corpus.filter((c) => c.classification === "valid")) {
    server.state.body = fixture(c.id);
    await page.goto(server.origin);
    await expect(
      page.getByText("Health available", { exact: true }).first(),
    ).toBeVisible();
    const expected = c.expected;
    await expect(page.locator("h1")).toHaveText(expected.host.hostname);
    for (const n of expected.network?.interfaces ?? []) {
      for (const value of [n.rx_bytes, n.tx_bytes])
        if (value !== null)
          await expect(
            page
              .locator(".network-grid")
              .getByText(value.replace("u64:", ""), { exact: true })
              .first(),
          ).toBeVisible();
    }
  }
});
test("S03-T09/T12 long labels, 64 interfaces and untrusted issue text", async ({
  page,
}, info) => {
  const long = "host-" + "long-label-".repeat(40);
  server.state.body = Buffer.from(
    fixture("interfaces-64")
      .toString()
      .replace('"joy-pi"', JSON.stringify(long)),
  );
  await page.setViewportSize({ width: 360, height: 900 });
  await page.goto(server.origin);
  await current(page);
  await expect(page.locator(".network-grid section")).toHaveCount(64);
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= innerWidth,
    ),
  ).toBe(true);
  await capture(page, info, "long-64-interfaces");
  const partial = corpus.find(
    (c) => c.id === "partial-permission_denied-memory",
  );
  const message = partial.expected.issues[0].message;
  const hostile =
    "<script>window.compromised=true</script> " +
    "Long unavailable message ".repeat(80);
  server.state.body = Buffer.from(
    fixture(partial.id)
      .toString()
      .replace(JSON.stringify(message), JSON.stringify(hostile)),
  );
  await page.goto(server.origin);
  await expect(page.getByText(hostile, { exact: false })).toBeVisible();
  expect(await page.evaluate(() => window.compromised)).toBeUndefined();
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= innerWidth,
    ),
  ).toBe(true);
  await capture(page, info, "long-text-issue");
});

test("S03-T04 actual React development StrictMode effect replay", async ({
  page,
}) => {
  await page.route("**/__test/strictmode", (route) =>
    route.fulfill({
      contentType: "text/html",
      body: '<!doctype html><html lang="en"><body><div id="root"></div><script type="module" src="/__test/strictmode.mjs"></script></body></html>',
    }),
  );
  await page.route("**/__test/strictmode.mjs", (route) =>
    route.fulfill({
      contentType: "text/javascript",
      body: readFileSync(resolve(root, "out/s03-strictmode/strictmode.mjs")),
    }),
  );
  await page.goto(server.origin + "/__test/strictmode");
  await current(page);
  await expect(page.locator("html")).toHaveAttribute("data-effect-setups", "2");
  await expect(page.locator("html")).toHaveAttribute(
    "data-effect-cleanups",
    "1",
  );
  expect(server.state.calls).toBe(1);
  await expect.poll(() => server.state.calls, { timeout: 7000 }).toBe(2);
  await page.goto("about:blank");
});

test("S03-T07/T12 Health absent before first success remains a usable Board response", async ({
  page,
}, info) => {
  server.state.status = 503;
  await page.goto(server.origin);
  await expect(
    page.getByText("Health unavailable", { exact: true }).first(),
  ).toBeVisible();
  await expect(
    page.getByText("Board reachable", { exact: true }).first(),
  ).toBeVisible();
  await expect(
    page.getByText("No success reported", { exact: true }),
  ).toBeVisible();
  await expect(page.getByText("9007199254740993", { exact: true })).toHaveCount(
    0,
  );
  await capture(page, info, "none-before-first-success");
});
