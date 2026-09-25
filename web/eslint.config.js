import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import reactHooks from "eslint-plugin-react-hooks";

export default [
  { ignores: ["web/dist/**", "web/node_modules/**"] },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    files: ["web/src/**/*.{ts,tsx}"],
    languageOptions: { globals: globals.browser },
    plugins: { "react-hooks": reactHooks },
    rules: reactHooks.configs.recommended.rules,
  },
  {
    files: ["scripts/browser.spec.mjs"],
    languageOptions: { globals: globals.browser },
  },
  {
    files: ["web/*.js", "web/*.ts", "scripts/*.mjs"],
    languageOptions: { globals: globals.node },
  },
];
