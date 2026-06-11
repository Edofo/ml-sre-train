# SLOs — recommender service

Status: active · Window: 30-day rolling · Owner: @Edofo · Review cadence: quarterly

## 1. Service and scope

These SLOs cover the user-facing endpoint of the recommender service: `POST /recommend`.

Excluded: `/health`, `/ready` (Kubernetes probes) and `/metrics` (Prometheus scrapes). They carry no user traffic and are filtered out of the RED metrics by the instrumentation middleware.

## 2. SLI definitions

### Availability

Proportion of `POST /recommend` requests answered without a server error.

- **Good event**: response with status `!~ 5..`. 4xx responses count as good: a malformed request is a client fault and must not burn the service's budget.
- **Source**: `http_requests_total{path="/recommend"}` (Prometheus, scraped every 15s).

```promql
1 - (
  sum(rate(http_requests_total{path="/recommend",status=~"5.."}[30d]))
  /
  sum(rate(http_requests_total{path="/recommend"}[30d]))
)
```

### Latency

Proportion of `POST /recommend` requests served in under **500 ms**.

- **Good event**: request observed in the `le="0.5"` histogram bucket. 500 ms is the product tolerance for an inline recommendation call; beyond that the caller falls back to a static list.
- **Source**: `http_request_duration_seconds` histogram. The threshold must be an existing bucket boundary (0.5 is part of the default buckets).

```promql
sum(rate(http_request_duration_seconds_bucket{path="/recommend",le="0.5"}[30d]))
/
sum(rate(http_request_duration_seconds_count{path="/recommend"}[30d]))
```

## 3. Objectives

| SLI | Objective | Window |
|---|---|---|
| Availability | **99.5 %** | 30 days rolling |
| Latency (< 500 ms) | **99 %** | 30 days rolling |

**Why not 99.9 %?** Each additional nine multiplies the engineering cost. A recommendation is a best-effort feature: callers degrade gracefully to a popular-items fallback, so a small error rate is invisible to end users. 99.5 % keeps a meaningful but achievable budget for a service deployed several times a day.

## 4. Error budget

| SLO | Budget over 30 days | Equivalent full outage |
|---|---|---|
| Availability 99.5 % | 0.5 % of requests | ~3.6 hours |
| Latency 99 % | 1 % of requests may be slow | — |

Policy (inspired by the [SRE Workbook example](https://sre.google/workbook/error-budget-policy/)):
- **Budget healthy** → ship features normally.
- **Budget exhausted** → feature freeze on the service; only reliability work until the 30-day SLI is back above target.
- Every budget-exhaustion event requires a postmortem (see `reliability/postmortems/`).

## 5. Alerting — multi-window multi-burn-rate

Alert rules are generated from the spec in this directory (see §7) following [Alerting on SLOs](https://sre.google/workbook/alerting-on-slos/).

| Burn rate | Long / short window | Severity | Expected response |
|---|---|---|---|
| 14.4× | 1 h / 5 m | `critical` (page) | Act immediately: budget gone in ~2 days |
| 6× | 6 h / 30 m | `critical` (page) | Act immediately |
| 3× | 24 h / 2 h | `warning` (ticket) | Handle within the day |
| 1× | 72 h / 6 h | `warning` (ticket) | Investigate this week |

Both windows must burn at once: the short window confirms the issue is still ongoing, which keeps false positives near zero. Runbooks: `reliability/runbooks/` (placeholder, Phase 3).

## 6. Dashboards

- **Recommend Traffic Dashboard** (Grafana, provisioned from `observability/grafana/`): rate, error ratio, p50/p95/p99.
- SLI/burn-rate series are available under `slo:sli_error:ratio_rate*` for ad-hoc queries.

## 7. Operations — regenerating the rules

The **source of truth** is the Sloth spec [`recommender-slos.yaml`](recommender-slos.yaml). The [`rules.yaml`](rules.yaml) PrometheusRule is **generated — never edit it by hand**.

```bash
make slo-generate   # sloth validate + sloth generate
kubectl apply -f observability/slo/rules.yaml
```

Workflow: edit the spec → `make slo-generate` → commit **both files** → apply (manual today, ArgoCD later). Changing an objective is an engineering decision: update §3 of this document in the same PR.

## 8. References

- [Implementing SLOs](https://sre.google/workbook/implementing-slos/) · [Alerting on SLOs](https://sre.google/workbook/alerting-on-slos/) · [Example SLO Document](https://sre.google/workbook/slo-document/)
- [Sloth](https://sloth.dev/) — SLO spec format and rule generation
