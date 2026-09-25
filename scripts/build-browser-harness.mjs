import { build } from "../web/node_modules/vite/dist/node/index.js";
import { fileURLToPath } from "node:url";
const web = fileURLToPath(new URL("../web/", import.meta.url));
await build({
  configFile: false,
  root: web,
  define: { "process.env.NODE_ENV": JSON.stringify("development") },
  build: {
    outDir: "../out/s03-strictmode",
    emptyOutDir: true,
    target: ["chrome111", "edge111", "firefox115", "safari16.4"],
    lib: {
      entry: "strictmode-entry.tsx",
      formats: ["es"],
      fileName: () => "strictmode.mjs",
    },
    minify: false,
  },
});
