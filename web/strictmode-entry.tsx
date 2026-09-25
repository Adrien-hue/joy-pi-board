// Test-only development React entry, emitted outside web/dist.
import { StrictMode, useEffect } from "react";
import { createRoot } from "react-dom/client";
import { App } from "./src/App";
let setups = 0,
  cleanups = 0;
function ReplayProbe() {
  useEffect(() => {
    document.documentElement.dataset.effectSetups = String(++setups);
    return () => {
      document.documentElement.dataset.effectCleanups = String(++cleanups);
    };
  }, []);
  return null;
}
createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ReplayProbe />
    <App />
  </StrictMode>,
);
