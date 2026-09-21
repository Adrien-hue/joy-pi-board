import { defineConfig } from "vitest/config";

export default defineConfig({
  base: "/",
  build: {
    outDir: "dist",
    emptyOutDir: true,
    assetsDir: "assets",
    target: ["chrome111", "edge111", "firefox115", "safari16.4"],
    sourcemap: false,
  },
  test: {
    environment: "node",
    include: ["src/**/*.test.{ts,tsx}"],
  },
});
