# 16. 削除セマンティクス

## 16.1 概要

本章では Resource の削除時に適用される動作を定義する。

削除対象には以下を含む。

- Device
- ManagementGroup
- User
- DeviceAddress
- Metadata
- Credential

削除時においても監査可能性を維持しなければならない。

---

## 16.2 削除原則

### 16.2.1 監査保持

削除後であっても AuditLog を失ってはならない。

---

### 16.2.2 履歴保持

削除後であっても UpdateHistory を保持してよい。

---

### 16.2.3 整合性維持

削除により Ownership Invariant を破ってはならない。

---

### 16.2.4 冪等性

同一削除操作が複数回実行されても不正状態になってはならない。

---

## 16.3 Device Delete

### 16.3.1 概要

Device 削除時は関連データを適切に削除または失効しなければならない。

---

### 16.3.2 DeviceCredential

関連する DeviceCredential をすべて失効しなければならない。

```text
Device
↓
Delete

All DeviceCredentials
↓
Revoked
```

---

### 16.3.3 Address

関連する Address を削除しなければならない。

対象:

- IPv4
- IPv6

---

### 16.3.4 Metadata

関連 Metadata を削除しなければならない。

---

### 16.3.5 Group Membership

ManagementGroup から Device を除去しなければならない。

---

### 16.3.6 AuditLog

AuditLog は削除してはならない。

---

### 16.3.7 UpdateHistory

UpdateHistory は削除してはならないことを推奨する。

---

### 16.3.8 Address Visibility

削除後は以下から返してはならない。

```text
Address Read API
```

---

### 16.3.9 Metadata Visibility

削除後は以下から返してはならない。

```text
Metadata API
```

---

## 16.4 ManagementGroup Delete

### 16.4.1 概要

ManagementGroup 削除時は Group 関連 Resource のみを削除する。

---

### 16.4.2 GroupReadToken

関連する GroupReadToken をすべて失効しなければならない。

---

### 16.4.3 Membership

Membership を削除しなければならない。

---

### 16.4.4 Device

Device 自体を削除してはならない。

例:

```text
Delete Group
≠
Delete Device
```

---

### 16.4.5 AuditLog

AuditLog は削除してはならない。

---

### 16.4.6 UpdateHistory

UpdateHistory は削除してはならないことを推奨する。

---

## 16.5 User Delete

### 16.5.1 概要

User Delete は最上位削除操作である。

---

### 16.5.2 UserSession

関連 UserSession をすべて失効しなければならない。

---

### 16.5.3 Devices

所有する Device へ Device Delete Semantics を適用しなければならない。

---

### 16.5.4 Groups

所有する ManagementGroup へ Group Delete Semantics を適用しなければならない。

---

### 16.5.5 Credentials

関連 Credential をすべて失効しなければならない。

---

### 16.5.6 Audit Retention

AuditLog は保持期間が終了するまで保持しなければならない。

---

### 16.5.7 Retention Period

既定保持期間:

```text
13 months
```

---

### 16.5.8 Visibility

削除済み User の AuditLog は管理者向け画面からのみ参照できるべきである。

---

## 16.6 Address Delete

### 16.6.1 発生条件

以下によって削除される。

- Replace
- Device Delete
- CleanupJob
- Administrative Delete

---

### 16.6.2 Lifecycle

```text
active
↓
deleted
```

または

```text
active
↓
expired
↓
deleted
```

---

### 16.6.3 Read Visibility

削除後は Read API から返してはならない。

---

## 16.7 Metadata Delete

### 16.7.1 発生条件

以下を削除として扱う。

```json
{
  "key": null
}
```

---

### 16.7.2 UpdateHistory

削除イベントは履歴へ記録しなければならない。

---

### 16.7.3 Value Logging

削除時も Value を履歴保存してはならない。

---

## 16.8 Credential Delete

### 16.8.1 DeviceCredential

削除または失効できなければならない。

---

### 16.8.2 GroupReadToken

削除または失効できなければならない。

---

### 16.8.3 UserSession

削除または失効できなければならない。

---

### 16.8.4 即時反映

失効は即時反映されなければならない。

---

## 16.9 Audit Requirements

以下の削除イベントは AuditLog に記録しなければならない。

- Device Delete
- Group Delete
- User Delete
- Address Delete
- Metadata Delete
- Credential Revoke

---

### 記録推奨項目

```text
actor
resource_type
resource_id
result
timestamp
reason
```

---

## 16.10 UpdateHistory Requirements

削除系履歴には少なくとも以下を保持すべきである。

```text
resource_type
```

```text
resource_id
```

```text
operation = delete
```

```text
timestamp
```

---

## 16.11 Cleanup Interaction

### 16.11.1 Expired Address

CleanupJob によって削除してよい。

---

### 16.11.2 AuditLog

保持期間超過後に削除してよい。

---

### 16.11.3 UpdateHistory

保持期間経過後に削除してよい。

---

## 16.12 不変条件

実装は以下を保証しなければならない。

```text
Device Delete
→ Revoke DeviceCredentials
```

```text
Device Delete
→ Delete Addresses
```

```text
Device Delete
→ Delete Metadata
```

```text
Device Delete
→ Remove Memberships
```

```text
ManagementGroup Delete
→ Revoke GroupReadTokens
```

```text
ManagementGroup Delete
→ MUST NOT Delete Device
```

```text
User Delete
→ Revoke UserSessions
```

```text
Delete Operation
→ Audit Logged
```

```text
AuditLog
→ Retained
```

```text
Deleted Resource
→ Not Returned By API
```

これらを満たさない実装は本仕様に適合しない。
