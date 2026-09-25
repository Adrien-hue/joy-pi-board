import type { Overview } from "./contract";
import { requestOverview, RequestFailure, type FailureKind } from "./request";

export interface Clock {
  mono: number;
  wall: number;
}
export interface Runtime {
  now(): Clock;
  later(fn: () => void, ms: number): ReturnType<typeof setTimeout>;
  cancel(id: ReturnType<typeof setTimeout>): void;
  microtask(fn: () => void): void;
  request(signal: AbortSignal, received: () => void): Promise<Overview>;
}
export const browserRuntime: Runtime = {
  now: () => ({ mono: performance.now(), wall: Date.now() }),
  later: (fn, ms) => setTimeout(fn, ms),
  cancel: (id) => clearTimeout(id),
  microtask: (fn) => queueMicrotask(fn),
  request: (signal, received) => requestOverview(signal, fetch, received),
};
export interface View {
  board: "loading" | "usable" | "paused" | FailureKind;
  httpStatus: number | undefined;
  envelope: Overview | null;
  metrics: "current" | "partial" | "stale" | "none";
  age: number | null;
}
function elapsed(a: Clock, b: Clock): number | null {
  if (
    ![a.mono, a.wall, b.mono, b.wall].every(Number.isFinite) ||
    b.mono < a.mono ||
    b.wall < a.wall
  )
    return null;
  const n = Math.max(b.mono - a.mono, b.wall - a.wall);
  return Number.isFinite(n) ? n : null;
}
const cap = (a: number, b: number) => Math.min(30001, a + Math.min(30001, b));

// One mounted page owns this controller. No global cache or persistence.
export class OverviewController {
  private mounted = false;
  private visible = false;
  private generation = 0;
  private epoch = 0;
  private cadence: ReturnType<typeof setTimeout> | undefined;
  private ageTimer: ReturnType<typeof setTimeout> | undefined;
  private expiry: ReturnType<typeof setTimeout> | undefined;
  private active: {
    abort: AbortController;
    timer: ReturnType<typeof setTimeout> | undefined;
  } | null = null;
  private anchor: { last: Clock; age: number } | null = null;
  private state: View = {
    board: "loading",
    httpStatus: undefined,
    envelope: null,
    metrics: "none",
    age: null,
  };
  constructor(
    private readonly runtime: Runtime,
    private readonly changed: () => void,
    private mask: () => void = () => {},
  ) {}
  setMask(mask: () => void): void {
    this.mask = mask;
  }
  start(visible: boolean): void {
    this.mounted = true;
    this.visible = visible;
    this.epoch = this.runtime.now().mono;
    const generation = ++this.generation;
    if (visible)
      this.runtime.microtask(() => {
        if (this.mounted && this.visible && generation === this.generation) {
          this.launch();
          this.schedule();
        }
      });
    else this.state = { ...this.state, board: "paused" };
  }
  stop(): void {
    this.pause();
    this.mounted = false;
  }
  pause(): void {
    if (!this.mounted) return;
    this.visible = false;
    this.generation++;
    this.anchor = null;
    this.mask();
    this.clearTimers();
    this.active?.abort.abort();
    if (this.active?.timer !== undefined)
      this.runtime.cancel(this.active.timer);
    this.state = { ...this.state, board: "paused", metrics: "none", age: null };
    this.changed();
  }
  resume(): void {
    if (!this.mounted || this.visible) return;
    this.visible = true;
    this.schedule(); // no immediate foreground retry
  }
  private clearTimers(): void {
    for (const id of [this.cadence, this.ageTimer, this.expiry])
      if (id !== undefined) this.runtime.cancel(id);
    this.cadence = this.ageTimer = this.expiry = undefined;
  }
  private schedule(): void {
    if (!this.mounted || !this.visible) return;
    if (this.cadence !== undefined) this.runtime.cancel(this.cadence);
    const now = this.runtime.now().mono;
    if (!Number.isFinite(now)) {
      this.invalidate();
      return;
    }
    if (!Number.isFinite(this.epoch) || now < this.epoch) this.epoch = now;
    const delay = 5000 - ((now - this.epoch) % 5000);
    const due = now + delay;
    this.cadence = this.runtime.later(() => {
      this.cadence = undefined;
      // Discard slots missed through suspension/event-loop delay.
      if (this.runtime.now().mono < due + 5000) this.launch();
      this.schedule();
    }, delay);
  }
  private invalidate(): void {
    this.anchor = null;
    this.mask();
    this.state = { ...this.state, metrics: "none", age: null };
  }
  private launch(): void {
    if (!this.mounted || !this.visible || this.active) return;
    const generation = this.generation,
      start = this.runtime.now();
    const attempt = {
      abort: new AbortController(),
      timer: undefined as ReturnType<typeof setTimeout> | undefined,
    };
    this.active = attempt;
    let terminal = false;
    let receipt: Clock | null = null;
    const durationAt = (now: Clock): number | null => {
      if (receipt === null) return elapsed(start, now);
      const before = elapsed(start, receipt),
        after = elapsed(receipt, now);
      if (before === null || after === null || !Number.isFinite(before + after))
        return null;
      return before + after;
    };
    const eligible = () =>
      this.mounted &&
      this.visible &&
      generation === this.generation &&
      this.active === attempt &&
      !terminal;
    const fail = (kind: FailureKind, status?: number) => {
      if (!eligible()) return;
      terminal = true;
      this.state = { ...this.state, board: kind, httpStatus: status };
      if (kind === "contract") this.invalidate();
      this.changed();
      this.armAge();
    };
    attempt.timer = this.runtime.later(() => {
      fail("timeout");
      attempt.abort.abort();
    }, 2000);
    void this.runtime
      .request(attempt.abort.signal, () => {
        receipt = this.runtime.now();
      })
      .then(
        (envelope) => {
          if (!eligible()) return;
          const now = this.runtime.now(),
            duration = durationAt(now);
          if (duration === null) {
            fail("contract");
            return;
          }
          if (duration >= 2000) {
            fail("timeout");
            return;
          }
          terminal = true;
          // Full request duration includes streaming, parsing and publication work.
          this.anchor =
            envelope.age === null
              ? null
              : { last: now, age: cap(envelope.age, duration) };
          this.state = {
            board: "usable",
            httpStatus: undefined,
            envelope,
            metrics: "none",
            age: null,
          };
          this.changed();
          this.armAge();
        },
        (error) => {
          const duration = durationAt(this.runtime.now());
          if (duration === null) fail("contract");
          else if (duration >= 2000) fail("timeout");
          else
            fail(
              error instanceof RequestFailure ? error.kind : "network",
              error instanceof RequestFailure ? error.status : undefined,
            );
        },
      )
      .finally(() => {
        if (attempt.timer !== undefined) this.runtime.cancel(attempt.timer);
        // Hold admission until asynchronous reader cleanup actually settled.
        if (this.active === attempt) this.active = null;
      });
  }
  view(): View {
    if (!this.visible || !this.anchor || !this.state.envelope?.health.snapshot)
      return { ...this.state, metrics: "none", age: null };
    const now = this.runtime.now(),
      delta = elapsed(this.anchor.last, now);
    if (delta === null) {
      this.invalidate();
      return { ...this.state, metrics: "none", age: null };
    }
    this.anchor = { last: now, age: cap(this.anchor.age, delta) };
    const age = this.anchor.age;
    const metrics =
      age > 30000
        ? "none"
        : this.state.board !== "usable" ||
            this.state.envelope.health.snapshot_state === "stale"
          ? "stale"
          : this.state.envelope.health.snapshot.issues.length
            ? "partial"
            : "current";
    if (metrics === "none") this.mask();
    return { ...this.state, metrics, age };
  }
  private armAge(): void {
    for (const id of [this.expiry, this.ageTimer])
      if (id !== undefined) this.runtime.cancel(id);
    this.expiry = this.ageTimer = undefined;
    const view = this.view();
    if (
      !this.mounted ||
      !this.visible ||
      view.metrics === "none" ||
      view.age === null
    )
      return;
    this.expiry = this.runtime.later(
      () => {
        this.view();
        this.changed();
        this.armAge();
      },
      Math.max(1, Math.floor(30000 - view.age) + 1),
    );
    this.ageTimer = this.runtime.later(() => {
      this.view();
      this.changed();
      this.armAge();
    }, 1000);
  }
}
