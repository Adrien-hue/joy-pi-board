import { build } from "../web/node_modules/vite/dist/node/index.js";
await build({
  configFile: false,
  root: "web",
  build: {
    outDir: "../out/measurement",
    emptyOutDir: true,
    target: ["chrome111", "edge111", "firefox115", "safari16.4"],
    lib: {
      entry: "measurement-entry.ts",
      formats: ["es"],
      fileName: () => "consumer.mjs",
    },
  },
});
