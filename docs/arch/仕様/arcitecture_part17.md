---
Title: 17. AuditLog
Description: 本章は「AuditLog」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「17. AuditLog」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 17. AuditLog

## 17.1 概要

本章では監査ログ（AuditLog）の要件を定義する。

AuditLog の目的は以下である。

- セキュリティ監査
- インシデント調査
- 操作追跡
- 認証失敗分析
- 外部連携障害分析
- コンプライアンス対応

AuditLog は UpdateHistory と異なる責務を持つ。

---

## 17.2 基本方針

### 17.2.1 改ざん耐性

AuditLog は改ざん困難でなければならない。

---

### 17.2.2 削除耐性

通常操作によって AuditLog を削除してはならない。

---

### 17.2.3 追跡可能性

重要なイベントは必ず AuditLog に記録しなければならない。

---

### 17.2.4 非公開性

AuditLog は権限を持つ User のみ参照できる。

---

## 17.3 記録対象イベント

### 17.3.1 Authentication Success

以下を記録しなければならない。

```text
Authentication Success
```

---

### 17.3.2 Authentication Failure

以下を記録しなければならない。

```text
Authentication Failure
```

---

### 17.3.3 Authorization Failure

以下を記録しなければならない。

```text
Authorization Failure
```

---

### 17.3.4 DeviceCredential Issue

以下を記録しなければならない。

```text
DeviceCredential Issue
```

---

### 17.3.5 DeviceCredential Revoke

以下を記録しなければならない。

```text
DeviceCredential Revoke
```

---

### 17.3.6 GroupReadToken Issue

以下を記録しなければならない。

```text
GroupReadToken Issue
```

---

### 17.3.7 GroupReadToken Revoke

以下を記録しなければならない。

```text
GroupReadToken Revoke
```

---

### 17.3.8 Address Update

以下を記録しなければならない。

```text
Address Update
```

---

### 17.3.9 Metadata Update

以下を記録しなければならない。

```text
Metadata Update
```

---

### 17.3.10 Device Delete

以下を記録しなければならない。

```text
Device Delete
```

---

### 17.3.11 ManagementGroup Delete

以下を記録しなければならない。

```text
ManagementGroup Delete
```

---

### 17.3.12 User Delete

以下を記録しなければならない。

```text
User Delete
```

---

### 17.3.13 Membership Change

以下を記録しなければならない。

```text
ManagementGroup Membership Change
```

---

### 17.3.14 Rate Limit

以下を記録しなければならない。

```text
Rate Limit Activated
```

---

### 17.3.15 External Integration Failure

以下を記録しなければならない。

```text
External Integration Failure
```

---

### 17.3.16 External Integration Suspension

以下を記録しなければならない。

```text
External Integration Suspended
```

---

### 17.3.17 External Integration Resume

以下を記録しなければならない。

```text
External Integration Resumed
```

---

### 17.3.18 CleanupJob

以下を記録しなければならない。

```text
CleanupJob Run
```

---

## 17.4 AuditLog Data Model

### 17.4.1 必須項目

少なくとも以下を保持すべきである。

| Field | Description |
|---------|---------|
| id | Internal ID |
| public_id | UUIDv7 |
| request_id | Request Identifier |
| actor_type | User / Device / System |
| actor_id | Actor Identifier |
| target_type | Resource Type |
| target_id | Resource Identifier |
| result | Success / Failure |
| created_at | Timestamp |

---

### 17.4.2 推奨項目

以下を保存することを推奨する。

| Field | Description |
|---------|---------|
| device_id | Device |
| credential_id | Credential |
| management_group_id | Group |
| observed_ip | Source Address |
| submitted_ip | Submitted Address |
| endpoint_type | v4only/v6only/dualstack |
| mode | Update Mode |
| user_agent | Client Info |
| reason | Failure Reason |

---

## 17.5 Actor Type

次の Actor をサポートすべきである。

```text
user
```

```text
device
```

```text
system
```

```text
worker
```

---

## 17.6 Result

次の Result をサポートすべきである。

```text
success
```

```text
failure
```

---

## 17.7 機微情報

### 17.7.1 保存禁止

以下を保存してはならない。

```text
Bearer Token
```

```text
Basic Password
```

```text
Digest Secret
```

```text
Session Secret
```

```text
Plaintext Credential
```

---

### 17.7.2 Metadata Value

Metadata の Full Value は保存してはならない。

---

### 17.7.3 Token Prefix

Token Prefix は保存してよい。

ただし認証情報として利用してはならない。

---

## 17.8 AuditLog Read

### 17.8.1 UserSession

UserSession による参照を許可してよい。

---

### 17.8.2 所有権

User は自身の Resource に関する AuditLog のみ参照できなければならない。

---

### 17.8.3 GroupReadToken

参照を許可してはならない。

---

### 17.8.4 DeviceCredential

参照を許可してはならない。

---

## 17.9 保持期間

### 17.9.1 既定値

AuditLog の既定保持期間。

```text
13 months
```

---

### 17.9.2 Cleanup

保持期間超過後に CleanupJob が削除してよい。

---

### 17.9.3 User Delete

User 削除によって AuditLog を削除してはならない。

---

### 17.9.4 Device Delete

Device 削除によって AuditLog を削除してはならない。

---

### 17.9.5 ManagementGroup Delete

ManagementGroup 削除によって AuditLog を削除してはならない。

---

## 17.10 AuditLog Integrity

### 17.10.1 Immutable Policy

記録済み AuditLog は変更しないことを推奨する。

---

### 17.10.2 Append Only

Append Only 方式を推奨する。

```text
Insert
✓

Update
×
```

---

## 17.11 Monitoring

実装は以下の Metrics を提供してよい。

```text
audit_events_total
```

```text
audit_failures_total
```

```text
authentication_failure_total
```

```text
authorization_failure_total
```

---

## 17.12 不変条件

実装は以下を保証しなければならない。

```text
Authentication Failure
→ Audit Logged
```

```text
Authorization Failure
→ Audit Logged
```

```text
Credential Operation
→ Audit Logged
```

```text
Delete Operation
→ Audit Logged
```

```text
AuditLog
→ Retained For 13 Months
```

```text
User Delete
→ MUST NOT Delete AuditLog
```

```text
Device Delete
→ MUST NOT Delete AuditLog
```

```text
ManagementGroup Delete
→ MUST NOT Delete AuditLog
```

```text
Secret
→ MUST NOT Be Logged
```

これらを満たさない実装は本仕様に適合しない。
