import { useId, type ReactNode } from "react";
import { useOverview } from "./overview/useOverview";
import type { View } from "./overview/controller";
import type { ByteGroup, HealthSnapshot } from "./health/types";
import {
  bytes,
  flag,
  floating,
  percent,
  unavailable,
  uptime,
} from "./presentation/format";

function Metric({
  label,
  value,
  exact,
  path,
  snapshot,
}: {
  label: string;
  value: string;
  exact?: string | undefined;
  path?: string;
  snapshot?: HealthSnapshot;
}) {
  const id = useId();
  const issues = snapshot?.issues.filter((issue) => issue.path === path) ?? [];
  return (
    <div className="metric">
      <dt>{label}</dt>
      <dd aria-describedby={issues.length ? id : undefined}>{value}</dd>
      {exact !== undefined && (
        <dd className="exact">
          <details>
            <summary>Exact {label.toLowerCase()}</summary>
            <span>{exact}</span>
          </details>
        </dd>
      )}
      {issues.length > 0 && (
        <dd id={id} className="issue">
          {issues.map((issue) => (
            <span key={issue.path}>
              {issue.message} <small>({issue.code})</small>
            </span>
          ))}
        </dd>
      )}
    </div>
  );
}
function Group({
  title,
  group,
  path,
  snapshot,
}: {
  title: string;
  group: ByteGroup;
  path: string;
  snapshot: HealthSnapshot;
}) {
  const id = useId();
  const issue = snapshot.issues.find((i) => i.path === path);
  return (
    <section className="card" aria-labelledby={id}>
      <h2 id={id}>{title}</h2>
      <dl aria-describedby={issue ? id + "-issue" : undefined}>
        <Metric
          label="Used"
          value={bytes(group.used_bytes)}
          exact={
            group.used_bytes === null
              ? unavailable
              : `${group.used_bytes} bytes`
          }
        />
        <Metric
          label="Usage"
          value={percent(group.used_bytes, group.total_bytes)}
        />
        <Metric
          label="Available"
          value={bytes(group.available_bytes)}
          exact={
            group.available_bytes === null
              ? unavailable
              : `${group.available_bytes} bytes`
          }
        />
        <Metric
          label="Total"
          value={bytes(group.total_bytes)}
          exact={
            group.total_bytes === null
              ? unavailable
              : `${group.total_bytes} bytes`
          }
        />
      </dl>
      {issue && (
        <p className="issue" id={id + "-issue"}>
          {issue.message} <small>({issue.code})</small>
        </p>
      )}
    </section>
  );
}
const reasons = {
  timeout: "Health did not respond in time",
  connection_failed: "Health could not be reached",
  upstream_error: "Health returned an error",
  invalid_response: "Health returned invalid data",
  unsupported_schema_version: "Health uses an unsupported schema version",
};
export function boardText(view: View): string {
  switch (view.board) {
    case "loading":
      return "Loading observations";
    case "usable":
      return "Board reachable";
    case "paused":
      return "Updates paused — waiting for fresh data";
    case "http":
      return `Board returned HTTP ${view.httpStatus}`;
    case "network":
      return "Board could not be reached";
    case "timeout":
      return "Board did not respond in time";
    case "contract":
      return "Board returned invalid data";
  }
}
function Section({ title, children }: { title: string; children: ReactNode }) {
  const id = useId();
  return (
    <section className="card" aria-labelledby={id}>
      <h2 id={id}>{title}</h2>
      {children}
    </section>
  );
}
export function Dashboard({
  view,
  gate,
  guard,
}: {
  view: View;
  gate?: React.Ref<HTMLDivElement>;
  guard?: () => void;
}) {
  const shown = view.metrics !== "none";
  const snapshot = shown ? view.envelope?.health.snapshot : null;
  const health = view.envelope?.health;
  const freshReport =
    view.board === "usable" && (shown || health?.snapshot_state === "none");
  const healthText = freshReport
    ? `Health ${health?.availability ?? "unknown"}`
    : "Health availability unknown";
  const freshness =
    view.metrics === "current"
      ? "Current observations"
      : view.metrics === "partial"
        ? "Current observations · Some measures unavailable"
        : view.metrics === "stale"
          ? "Stale observations"
          : "No presentable observations";
  const announcement = `${boardText(view)}. ${healthText}. ${freshness}.`;
  const pi = snapshot?.raspberry_pi;
  return (
    <main>
      <a className="skip" href="#services">
        Skip to service information
      </a>
      <header className="page-header">
        <p className="eyebrow">
          JOY PI HOME <span>/ BOARD</span>
        </p>
        <h1>{snapshot?.host.hostname ?? "Joy Pi Board"}</h1>
        <p className="intro">Your host, at a glance.</p>
        <div className="badges">
          <span>{boardText(view)}</span>
          <span>{healthText}</span>
        </div>
      </header>
      <p
        className="sr-only"
        role="status"
        aria-live="polite"
        aria-atomic="true"
      >
        {announcement}
      </p>
      <div className={`freshness ${view.metrics === "stale" ? "stale" : ""}`}>
        <strong>{freshness}</strong>
        {view.age !== null && shown && (
          <span>
            Data age: at most approximately {Math.ceil(view.age / 1000)} s
          </span>
        )}
        {freshReport && health?.reason && <p>{reasons[health.reason]}</p>}
      </div>
      {!shown && (
        <section className="empty">
          <h2>
            {view.board === "loading"
              ? "Connecting to Board"
              : "Waiting for observations"}
          </h2>
          <p>
            {view.board === "paused"
              ? "Metrics stay hidden until a new response is validated."
              : "Available observations will appear automatically. No action is needed."}
          </p>
        </section>
      )}
      <div
        ref={gate}
        hidden={!shown}
        onClickCapture={guard}
        onKeyDownCapture={guard}
        onFocusCapture={guard}
      >
        {snapshot && (
          <>
            <div className="primary-grid">
              <Section title="CPU">
                <dl>
                  <Metric
                    label="Utilization"
                    value={floating(snapshot.cpu.utilization_percent, 1, "%")}
                    exact={
                      snapshot.cpu.utilization_percent === null
                        ? unavailable
                        : `${snapshot.cpu.utilization_percent}%`
                    }
                    path="/cpu/utilization_percent"
                    snapshot={snapshot}
                  />
                  <Metric
                    label="Logical CPUs"
                    value={
                      snapshot.cpu.logical_cpu_count?.toString() ?? unavailable
                    }
                    path="/cpu/logical_cpu_count"
                    snapshot={snapshot}
                  />
                </dl>
              </Section>
              <Group
                title="Memory"
                group={snapshot.memory}
                path="/memory"
                snapshot={snapshot}
              />
              <Group
                title="Root filesystem"
                group={snapshot.root_filesystem}
                path="/root_filesystem"
                snapshot={snapshot}
              />
              <Section title="SoC temperature">
                <dl>
                  <Metric
                    label="Temperature"
                    value={floating(
                      snapshot.raspberry_pi.soc_temperature_celsius,
                      1,
                      " °C",
                    )}
                    exact={
                      snapshot.raspberry_pi.soc_temperature_celsius === null
                        ? unavailable
                        : `${snapshot.raspberry_pi.soc_temperature_celsius} °C`
                    }
                    path="/raspberry_pi/soc_temperature_celsius"
                    snapshot={snapshot}
                  />
                </dl>
              </Section>
            </div>
            <div className="secondary-grid">
              <Section title="Load averages">
                <dl
                  className="load-grid"
                  aria-describedby={
                    snapshot.load.one_minute === null ? "load-issue" : undefined
                  }
                >
                  {Object.entries(snapshot.load).map(([key, value], i) => (
                    <Metric
                      key={key}
                      label={`${[1, 5, 15][i]} minute${i ? "s" : ""}`}
                      value={floating(value, 2)}
                      exact={value === null ? unavailable : String(value)}
                    />
                  ))}
                </dl>
                {snapshot.issues
                  .filter((i) => i.path === "/load")
                  .map((i) => (
                    <p className="issue" id="load-issue" key={i.path}>
                      {i.message} <small>({i.code})</small>
                    </p>
                  ))}
              </Section>
              <Section title="Uptime">
                <dl>
                  <Metric
                    label="At observation"
                    value={uptime(snapshot.uptime_seconds)}
                    exact={
                      snapshot.uptime_seconds === null
                        ? unavailable
                        : `${snapshot.uptime_seconds} seconds`
                    }
                    path="/uptime_seconds"
                    snapshot={snapshot}
                  />
                </dl>
              </Section>
            </div>
            <section className="network" aria-labelledby="network-title">
              <div className="section-heading">
                <h2 id="network-title" tabIndex={-1}>
                  Network interfaces
                </h2>
                <p>RX / TX are cumulative counters, not rates.</p>
              </div>
              {snapshot.network === null ? (
                <p className="card">
                  {unavailable}.{" "}
                  {snapshot.issues.find((i) => i.path === "/network")?.message}
                </p>
              ) : snapshot.network.interfaces.length === 0 ? (
                <p className="card">No interfaces reported.</p>
              ) : (
                <div className="network-grid">
                  {snapshot.network.interfaces.map((n, i) => (
                    <section className="card" key={n.name}>
                      <h3>{n.name}</h3>
                      <dl>
                        <Metric
                          label="Operational state"
                          value={n.state ?? unavailable}
                          path={`/network/interfaces/${i}/state`}
                          snapshot={snapshot}
                        />
                        <Metric
                          label="RX · received bytes"
                          value={n.rx_bytes?.toString() ?? unavailable}
                          path={`/network/interfaces/${i}/rx_bytes`}
                          snapshot={snapshot}
                        />
                        <Metric
                          label="TX · transmitted bytes"
                          value={n.tx_bytes?.toString() ?? unavailable}
                          path={`/network/interfaces/${i}/tx_bytes`}
                          snapshot={snapshot}
                        />
                      </dl>
                    </section>
                  ))}
                </div>
              )}
            </section>
            <section aria-labelledby="pi-title">
              <div className="section-heading">
                <h2 id="pi-title">Raspberry Pi indicators</h2>
                <p>Observed flags, without a health score.</p>
              </div>
              <div className="secondary-grid">
                {(["active", "occurred_since_boot"] as const).map(
                  (suffix, index) => (
                    <Section
                      title={
                        index ? "Occurred since boot" : "Active at observation"
                      }
                      key={suffix}
                    >
                      <dl>
                        {(["thermal_throttling", "undervoltage"] as const).map(
                          (name, i) => {
                            const key = `${name}_${suffix}` as const;
                            return (
                              <Metric
                                key={key}
                                label={
                                  i ? "Undervoltage" : "Thermal throttling"
                                }
                                value={flag(pi![key])}
                                path={`/raspberry_pi/${key}`}
                                snapshot={snapshot}
                              />
                            );
                          },
                        )}
                      </dl>
                    </Section>
                  ),
                )}
              </div>
            </section>
          </>
        )}
      </div>
      <footer id="services" className="services" tabIndex={-1}>
        <h2>Services &amp; freshness</h2>
        <dl>
          <Metric label="Board" value={boardText(view)} />
          <Metric label="Health" value={healthText} />
          <Metric
            label="Health observation (UTC)"
            value={snapshot?.observed_at ?? "Not currently presented"}
          />
          <Metric
            label="Last validated receipt by Board (UTC)"
            value={health?.last_success_at ?? "No success reported"}
          />
          <Metric
            label="Board response generated (UTC)"
            value={view.envelope?.generated_at ?? "No response yet"}
          />
        </dl>
        <p className="principle">Board presents. Health observes.</p>
      </footer>
    </main>
  );
}
export function App() {
  const state = useOverview();
  return <Dashboard {...state} />;
}
