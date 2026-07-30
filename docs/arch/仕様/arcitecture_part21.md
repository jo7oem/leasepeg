# 21. External API Integration

## 21.1 概要

本章では External API Integration の要件を定義する。

External API Integration は、本システムの状態変化を外部システムへ通知するための機能である。

本機能は以下を目的とする。

- Webhook 連携
- 外部 CMDB 更新
- 外部 DNS 管理
- 通知システム連携
- 運用自動化

---

## 21.2 基本原則

### 21.2.1 非同期実行

External API Integration は非同期で実行しなければならない。

---

### 21.2.2 API 応答分離

External API の成否は API Request の成否に影響してはならない。

---

### 21.2.3 結果整合性

External API は最終的に内部状態と一致しなければならない。

---

### 21.2.4 冪等性

同一 Event の重複送信に耐えなければならない。

---

## 21.3 イベント処理モデル

### 21.3.1 基本フロー

```text
Client Request
↓
API Server
↓
Database Update
↓
Integration Event Create
↓
Response Return
↓
Worker
↓
External API
```

---

### 21.3.2 禁止構成

以下は禁止する。

```text
Request
↓
External API
↓
Database Update
↓
Response
```

---

### 21.3.3 理由

以下を回避するため。

- Timeout
- Latency Increase
- External Dependency Failure
- Availability Reduction

---

## 21.4 Integration Event

### 21.4.1 目的

External API 呼び出し対象を管理する。

---

### 21.4.2 必須項目

少なくとも以下を保持すべきである。

| Field | Description |
|---------|---------|
| event_type | Event Type |
| target_type | Resource Type |
| target_id | Resource Identifier |
| payload | Event Payload |
| status | Processing Status |
| retry_count | Retry Count |
| last_error | Last Error |
| created_at | Creation Time |
| updated_at | Update Time |

---

## 21.5 Event Types

少なくとも以下をサポートすべきである。

```text
address_updated
```

```text
address_deleted
```

```text
device_created
```

```text
device_deleted
```

```text
metadata_updated
```

```text
credential_revoked
```

---

## 21.6 Worker Processing

### 21.6.1 責務

Worker は以下を実行する。

- Event Fetch
- External Call
- Retry
- Failure Recording

---

### 21.6.2 並列実行

複数 Worker による処理を許可してよい。

---

### 21.6.3 冪等性

重複実行によって不正状態を発生させてはならない。

---

## 21.7 Retry Policy

### 21.7.1 必要性

一時的障害を吸収するため Retry を実施する。

---

### 21.7.2 管理項目

少なくとも以下を定義すべきである。

```text
retry_count
```

```text
retry_interval
```

```text
timeout
```

---

### 21.7.3 Backoff

Exponential Backoff を推奨する。

---

### 21.7.4 Retry Exhaustion

上限到達後は Failure として扱う。

---

## 21.8 Dead Letter Policy

### 21.8.1 目的

再試行不能な Event を隔離する。

---

### 21.8.2 条件

以下を満たした場合。

```text
Retry Count > Max Retry
```

---

### 21.8.3 処理

Dead Letter Queue または同等機能へ移動してよい。

---

### 21.8.4 Audit

Dead Letter 化を AuditLog に記録すべきである。

---

## 21.9 Integration State

### 21.9.1 状態一覧

以下をサポートする。

```text
enabled
```

```text
disabled
```

```text
suspended
```

---

### 21.9.2 enabled

正常動作状態。

---

### 21.9.3 disabled

管理者停止状態。

---

### 21.9.4 suspended

障害による自動停止状態。

---

## 21.10 Suspension Policy

### 21.10.1 Failure Tracking

連続 Failure 数を記録しなければならない。

---

### 21.10.2 自動停止

一定回数以上失敗した場合、自動停止してよい。

---

### 21.10.3 状態保持

停止状態を永続化しなければならない。

---

### 21.10.4 Audit

停止イベントを AuditLog に記録しなければならない。

---

## 21.11 Recovery Policy

### 21.11.1 状態確認

User は停止状態を確認できるべきである。

---

### 21.11.2 Failure Reason

Failure Reason を保持しなければならない。

---

### 21.11.3 Resume

User は管理画面から Resume を実行できるべきである。

---

### 21.11.4 Audit

Resume 操作は AuditLog に記録しなければならない。

---

## 21.12 User Notification

### 21.12.1 障害通知

Integration が停止した場合、User に通知してよい。

---

### 21.12.2 回復通知

Resume 後に通知してよい。

---

## 21.13 Security

### 21.13.1 Secret

External API Secret は平文保存してはならない。

---

### 21.13.2 Timeout

External API Call は Timeout を持たなければならない。

---

### 21.13.3 Isolation

External API 障害によって内部状態を破壊してはならない。

---

## 21.14 Audit Requirements

以下を AuditLog に記録しなければならない。

- Integration Enabled
- Integration Disabled
- Integration Suspended
- Integration Resumed
- Integration Failure
- Dead Letter

---

## 21.15 Monitoring

以下の Metrics を提供することを推奨する。

```text
external_integration_success_total
```

```text
external_integration_failure_total
```

```text
external_integration_suspended_total
```

```text
external_integration_retry_total
```

```text
external_integration_dead_letter_total
```

---

## 21.16 不変条件

実装は以下を保証しなければならない。

```text
External Integration
→ Asynchronous
```

```text
API Response
→ Independent Of External API
```

```text
Retry
→ Safe
```

```text
Dead Letter
→ Supported
```

```text
Repeated Failure
→ May Suspend Integration
```

```text
Suspended State
→ Persisted
```

```text
Resume
→ User Controllable
```

```text
Integration Failure
→ Audit Logged
```

```text
External Failure
→ MUST NOT Corrupt Internal State
```

これらを満たさない実装は本仕様に適合しない。
