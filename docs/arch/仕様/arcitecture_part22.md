# 22. Monitoring and Observability

## 22.1 概要

本章では本システムの監視および可観測性要件を定義する。

本システムは運用可能性を重視し、少なくとも以下を提供しなければならない。

- Health Check
- Readiness Check
- Metrics
- Operational Logs

---

## 22.2 目的

Monitoring の目的は以下である。

- 障害検知
- 可用性監視
- 性能監視
- セキュリティ監視
- Capacity Planning
- 障害解析

---

## 22.3 Monitoring API

### 22.3.1 役割

Monitoring API はシステム監視専用 API とする。

---

### 22.3.2 Port 分離

Monitoring API は Application API と別 Port で Listen しなければならない。

例:

```text
Application API
:8443

Monitoring API
:9090
```

---

### 22.3.3 非公開原則

Monitoring API は既定で外部公開してはならない。

---

### 22.3.4 例外

外部公開する場合は別認証を要求しなければならない。

---

## 22.4 Health Check

### 22.4.1 Liveness

プロセス生存確認。

例:

```http
GET /healthcheck/live
```

---

### 22.4.2 成功例

```http
200 OK
```

---

### 22.4.3 Failure

致命的エラー時は非成功ステータスを返してよい。

---

## 22.5 Readiness Check

### 22.5.1 目的

トラフィック受入可能状態を判定する。

---

### 22.5.2 Endpoint

例:

```http
GET /healthcheck/ready
```

---

### 22.5.3 判定対象

少なくとも以下を確認すべきである。

- Database
- Queue
- Internal Dependencies

---

### 22.5.4 Failure

依存システムが利用不可の場合は Ready として扱ってはならない。

---

## 22.6 Metrics

### 22.6.1 形式

以下を推奨する。

```text
Prometheus
```

または

```text
OpenTelemetry
```

---

### 22.6.2 Endpoint

例:

```http
GET /metrics
```

---

## 22.7 Authentication Metrics

### 22.7.1 Success Count

```text
authentication_success_total
```

---

### 22.7.2 Failure Count

```text
authentication_failure_total
```

---

### 22.7.3 Authorization Failure

```text
authorization_failure_total
```

---

## 22.8 Rate Limit Metrics

### 22.8.1 Activation Count

```text
rate_limit_activations_total
```

---

### 22.8.2 Source IP

実装は以下も提供してよい。

```text
rate_limit_source_ip_total
```

---

### 22.8.3 Credential

```text
rate_limit_credential_total
```

---

## 22.9 Address Metrics

### 22.9.1 Update Count

```text
address_updates_total
```

---

### 22.9.2 Update Failure

```text
address_update_failures_total
```

---

### 22.9.3 Active Address

```text
active_addresses
```

---

### 22.9.4 Expired Address

```text
expired_addresses
```

---

## 22.10 Credential Metrics

### 22.10.1 DeviceCredential Count

```text
device_credentials_total
```

---

### 22.10.2 GroupReadToken Count

```text
group_read_tokens_total
```

---

### 22.10.3 Revoked Credential

```text
revoked_credentials_total
```

---

### 22.10.4 Expired Credential

```text
expired_credentials_total
```

---

## 22.11 External Integration Metrics

### 22.11.1 Success

```text
external_integration_success_total
```

---

### 22.11.2 Failure

```text
external_integration_failure_total
```

---

### 22.11.3 Suspension

```text
external_integration_suspended_total
```

---

### 22.11.4 Retry

```text
external_integration_retry_total
```

---

## 22.12 Cleanup Metrics

### 22.12.1 Run Count

```text
cleanup_job_runs_total
```

---

### 22.12.2 Failure Count

```text
cleanup_job_failures_total
```

---

### 22.12.3 Deleted Records

```text
cleanup_job_deleted_records_total
```

---

### 22.12.4 Duration

```text
cleanup_job_duration_seconds
```

---

## 22.13 API Metrics

### 22.13.1 Requests

```text
api_requests_total
```

---

### 22.13.2 Latency

```text
api_request_duration_seconds
```

---

### 22.13.3 Errors

```text
api_errors_total
```

---

### 22.13.4 HTTP Status

実装は以下を集計してよい。

```text
2xx
```

```text
4xx
```

```text
5xx
```

---

## 22.14 Database Metrics

### 22.14.1 DB Error

```text
db_errors_total
```

---

### 22.14.2 Query Latency

```text
db_query_duration_seconds
```

---

### 22.14.3 Connection Usage

```text
db_connections_active
```

---

## 22.15 Worker Metrics

### 22.15.1 Worker Error

```text
worker_errors_total
```

---

### 22.15.2 Queue Depth

```text
queue_depth
```

---

### 22.15.3 Event Retry

```text
worker_retries_total
```

---

### 22.15.4 Dead Letter

```text
dead_letter_events_total
```

---

## 22.16 Operational Logs

### 22.16.1 目的

運用障害解析を支援する。

---

### 22.16.2 記録対象

少なくとも以下を記録すべきである。

- Worker Failure
- CleanupJob Result
- Database Error
- Queue Failure
- Integration Failure

---

### 22.16.3 Level

推奨:

```text
INFO
WARN
ERROR
```

---

## 22.17 Alerting

### 22.17.1 推奨対象

以下に Alert を設定することを推奨する。

- Authentication Failure Spike
- Database Error Spike
- Queue Backlog
- Worker Failure
- Cleanup Failure
- External Integration Failure

---

### 22.17.2 Critical Alert

以下は Critical として扱ってよい。

```text
Database Unavailable
```

```text
API Unavailable
```

---

## 22.18 不変条件

実装は以下を保証しなければならない。

```text
Monitoring API
→ Separate Port
```

```text
Monitoring API
→ Not Public By Default
```

```text
Liveness Endpoint
→ Available
```

```text
Readiness Endpoint
→ Available
```

```text
Authentication Metrics
→ Collected
```

```text
Address Metrics
→ Collected
```

```text
External Integration Metrics
→ Collected
```

```text
Cleanup Metrics
→ Collected
```

```text
Operational Failure
→ Logged
```

これらを満たさない実装は本仕様に適合しない。
