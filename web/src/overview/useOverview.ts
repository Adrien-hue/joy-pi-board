import { useEffect, useLayoutEffect, useRef, useState } from "react";
import { browserRuntime, OverviewController, type View } from "./controller";

export function useOverview() {
  const gate = useRef<HTMLDivElement>(null);
  const [, update] = useState(0);
  const [owned] = useState(
    () => new OverviewController(browserRuntime, () => update((n) => n + 1)),
  );
  useEffect(() => {
    owned.setMask(() => {
      if (gate.current) {
        if (
          gate.current.contains(document.activeElement) &&
          document.visibilityState === "visible"
        ) {
          gate.current
            .closest("main")
            ?.querySelector<HTMLElement>("#services")
            ?.focus({ preventScroll: true });
        }
        gate.current.hidden = true;
      }
    });
    const hide = () => owned.pause();
    const show = () => {
      if (document.visibilityState === "visible") owned.resume();
    };
    const visibility = () =>
      document.visibilityState === "hidden" ? hide() : show();
    const pageshow = (e: PageTransitionEvent) => {
      if (e.persisted) show();
    };
    document.addEventListener("visibilitychange", visibility);
    document.addEventListener("freeze", hide);
    document.addEventListener("resume", show);
    window.addEventListener("pagehide", hide);
    window.addEventListener("pageshow", pageshow);
    owned.start(document.visibilityState === "visible");
    return () => {
      owned.stop();
      document.removeEventListener("visibilitychange", visibility);
      document.removeEventListener("freeze", hide);
      document.removeEventListener("resume", show);
      window.removeEventListener("pagehide", hide);
      window.removeEventListener("pageshow", pageshow);
    };
  }, [owned]);
  const view: View = owned.view();
  useLayoutEffect(() => {
    if (gate.current) gate.current.hidden = view.metrics === "none";
  });
  // Events which expose exact values recheck eligibility, even with delayed timers.
  const guard = () => {
    if (owned.view().metrics === "none" && gate.current)
      gate.current.hidden = true;
  };
  return { view, gate, guard };
}
