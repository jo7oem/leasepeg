# 20. CleanupJob

## 20.1 概要

本章では CleanupJob の要件を定義する。

CleanupJob は不要データおよび保持期間を超過したデータを削除するための内部処理である。

CleanupJob の目的は以下である。

- Storage Optimization
- Performance Maintenance
- Retention Policy Enforcement
- Operational Automation

---

## 20.2 基本方針

### 20.2.1 自動化

CleanupJob は自動実行できなければならない。

---

### 20.2.2 安全性

CleanupJob は有効データを削除してはならない。

---

### 20.2.3 冪等性

同一 CleanupJob を複数回実行しても不正状態になってはならない。

---

### 20.2.4 監査可能性

CleanupJob の実行結果は監査可能でなければならない。

---

## 20.3 Cleanup 対象

### 20.3.1 Expired Address

CleanupJob は Expired Address を削除してよい。

条件:

```text
expires_at <= now
```

---

### 20.3.2 Expired Credential

失効済みかつ保持不要となった Credential 関連データを削除してよい。

---

### 20.3.3 UpdateHistory

保持期間を超過した UpdateHistory を削除してよい。

---

### 20.3.4 AuditLog

保持期間を超過した AuditLog を削除してよい。

---

### 20.3.5 Temporary Data

実装依存の一時データを削除してよい。

---

## 20.4 Cleanup 対象外

### 20.4.1 Active Address

削除してはならない。

---

### 20.4.2 Valid Credential

削除してはならない。

---

### 20.4.3 Active User Resource

削除してはならない。

---

### 20.4.4 保持期間内 AuditLog

削除してはならない。

---

## 20.5 起動方法

### 20.5.1 Scheduler

定期実行してよい。

例:

```text
Every Hour
```

```text
Every Day
```

---

### 20.5.2 Internal API

内部 API から起動してよい。

---

### 20.5.3 Manual Trigger

管理者による明示的起動を許可してよい。

---

### 20.5.4 Worker Startup

Worker 起動時に実行してもよい。

---

## 20.6 Internal Cleanup API

### 20.6.1 利用目的

CleanupJob の管理用途。

---

### 20.6.2 公開制御

既定では外部公開してはならない。

---

### 20.6.3 強制公開時

公開する場合は強い認証を要求しなければならない。

---

### 20.6.4 分離

Application API と分離すべきである。

---

## 20.7 Cleanup 実行手順

### 20.7.1 推奨フロー

```text
Start
↓
Acquire Lock
↓
Scan Targets
↓
Delete Eligible Records
↓
Write AuditLog
↓
Release Lock
↓
Finish
```

---

### 20.7.2 排他制御

複数 CleanupJob の同時実行を防止することを推奨する。

---

### 20.7.3 Retry

失敗した場合は再実行してよい。

---

## 20.8 CleanupJob Run Record

### 20.8.1 目的

実行履歴を保存する。

---

### 20.8.2 保存項目

少なくとも以下を保持すべきである。

| Field | Description |
|---------|---------|
| job_name | Job Name |
| started_at | Start Time |
| finished_at | Finish Time |
| result | Result |
| deleted_record_count | Deleted Count |
| failure_reason | Failure Reason |
| triggered_by | Trigger Source |

---

## 20.9 Trigger Source

以下をサポートすべきである。

```text
scheduler
```

```text
manual
```

```text
system
```

```text
api
```

---

## 20.10 Audit Requirements

### 20.10.1 実行記録

CleanupJob 実行を AuditLog に記録しなければならない。

---

### 20.10.2 記録項目

少なくとも以下を記録すべきである。

```text
job_name
```

```text
started_at
```

```text
result
```

```text
deleted_record_count
```

---

### 20.10.3 Failure

失敗時も監査記録を残さなければならない。

---

## 20.11 Failure Handling

### 20.11.1 Partial Success

部分成功を許可してよい。

---

### 20.11.2 Rollback

大量削除を単一 Transaction とする必要はない。

---

### 20.11.3 Retry

次回実行時に継続処理できるべきである。

---

## 20.12 Monitoring

### 20.12.1 Metrics

以下を提供することを推奨する。

```text
cleanup_job_runs_total
```

```text
cleanup_job_failures_total
```

```text
cleanup_job_deleted_records_total
```

```text
cleanup_job_duration_seconds
```

---

### 20.12.2 Alert

以下について Alert を設定してよい。

- Repeated Failure
- Long Execution Time
- Excessive Deletion Count

---

## 20.13 Retention Policy

### AuditLog

既定:

```text
13 months
```

---

### UpdateHistory

実装依存。

---

### Address

Expired 状態になった後、CleanupJob により削除してよい。

---

### Credential

失効後、保持期間経過時に削除してよい。

---

## 20.14 不変条件

実装は以下を保証しなければならない。

```text
Active Address
→ MUST NOT Be Deleted
```

```text
Expired Address
→ May Be Deleted
```

```text
CleanupJob
→ Audit Logged
```

```text
CleanupJob
→ Retry Safe
```

```text
CleanupJob
→ Idempotent
```

```text
Internal Cleanup API
→ Not Public By Default
```

```text
Cleanup Failure
→ Logged
```

```text
Run Result
→ Recorded
```

```text
Deleted Count
→ Recorded
```

これらを満たさない実装は本仕様に適合しない。
