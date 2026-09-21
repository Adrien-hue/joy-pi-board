export type UInt64 = bigint;
export type IssueCode =
  "unsupported" | "permission_denied" | "temporarily_unavailable";
export interface Issue {
  path: string;
  code: IssueCode;
  message: string;
}
export interface ByteGroup {
  total_bytes: UInt64 | null;
  available_bytes: UInt64 | null;
  used_bytes: UInt64 | null;
}
export interface NetworkInterface {
  name: string;
  state: string | null;
  rx_bytes: UInt64 | null;
  tx_bytes: UInt64 | null;
}
export interface HealthSnapshot {
  schema_version: "1.0";
  observed_at: string;
  host: { hostname: string };
  cpu: { utilization_percent: number | null; logical_cpu_count: UInt64 | null };
  load: {
    one_minute: number | null;
    five_minutes: number | null;
    fifteen_minutes: number | null;
  };
  memory: ByteGroup;
  root_filesystem: ByteGroup;
  uptime_seconds: UInt64 | null;
  network: { interfaces: NetworkInterface[] } | null;
  raspberry_pi: {
    soc_temperature_celsius: number | null;
    thermal_throttling_active: boolean | null;
    thermal_throttling_occurred_since_boot: boolean | null;
    undervoltage_active: boolean | null;
    undervoltage_occurred_since_boot: boolean | null;
  };
  issues: Issue[];
}
