# 25. SQL Constraints

## 25.1 概要

本章では Database レベルで維持すべき制約を定義する。

本章の目的は以下である。

- データ整合性維持
- 重複防止
- 所有権保証
- Application Bug の影響緩和
- 安全なデータ管理

実装は可能な限り Database Constraint を利用して整合性を保証すべきである。

---

## 25.2 Constraint 方針

### 25.2.1 原則

整合性は以下の複数層で維持する。

```text
Application Validation
+
Transaction
+
Database Constraint
```

---

### 25.2.2 推奨

可能な限り Database 側で保証することを推奨する。

---

### 25.2.3 Authorization

Authorization 自体を SQL Constraint で代替してはならない。

---

## 25.3 Public ID

### 25.3.1 一意性

公開 Identifier は一意でなければならない。

---

### 25.3.2 制約

```sql
UNIQUE(public_id)
```

---

### 25.3.3 対象

以下のテーブルに適用する。

- users
- devices
- management_groups
- credentials
- audit_logs
- update_histories
- integration_events
- cleanup_job_runs

---

## 25.4 Device Address Constraint

### 25.4.1 Family 一意性

同一 Device に対して同一 Family を複数保持してはならない。

---

### 25.4.2 制約

```sql
UNIQUE(device_id, family)
```

---

### 25.4.3 許可例

```text
Device A
 ├─ IPv4
 └─ IPv6
```

---

### 25.4.4 禁止例

```text
Device A
 ├─ IPv4
 └─ IPv4
```

---

## 25.5 Metadata Constraint

### 25.5.1 Key 一意性

同一 Device 内で Metadata Key は一意でなければならない。

---

### 25.5.2 制約

```sql
UNIQUE(device_id, key)
```

---

### 25.5.3 許可例

```text
name
location
note
```

---

### 25.5.4 禁止例

```text
name
name
```

---

## 25.6 Membership Constraint

### 25.6.1 重複禁止

同一 Device を同一 Group へ複数回追加してはならない。

---

### 25.6.2 制約

```sql
UNIQUE(management_group_id, device_id)
```

---

### 25.6.3 理由

以下を防ぐ。

- 重複 Membership
- 集計不整合
- Read API 重複

---

## 25.7 Token Prefix Constraint

### 25.7.1 要件

Token Prefix は Credential Type ごとに一意でなければならない。

---

### 25.7.2 制約

```sql
UNIQUE(type, token_prefix)
```

---

### 25.7.3 理由

以下を防止する。

- 衝突
- Audit 混乱
- Rate Limit 誤判定

---

## 25.8 Foreign Keys

### 25.8.1 users

関連先:

```text
root entity
```

---

### 25.8.2 devices

```sql
FOREIGN KEY(owner_user_id)
REFERENCES users(id)
```

---

### 25.8.3 management_groups

```sql
FOREIGN KEY(owner_user_id)
REFERENCES users(id)
```

---

### 25.8.4 user_sessions

```sql
FOREIGN KEY(user_id)
REFERENCES users(id)
```

---

### 25.8.5 device_credentials

```sql
FOREIGN KEY(device_id)
REFERENCES devices(id)
```

---

### 25.8.6 group_read_tokens

```sql
FOREIGN KEY(management_group_id)
REFERENCES management_groups(id)
```

---

### 25.8.7 device_addresses

```sql
FOREIGN KEY(device_id)
REFERENCES devices(id)
```

---

### 25.8.8 device_metadata

```sql
FOREIGN KEY(device_id)
REFERENCES devices(id)
```

---

### 25.8.9 management_group_devices

```sql
FOREIGN KEY(management_group_id)
REFERENCES management_groups(id)
```

```sql
FOREIGN KEY(device_id)
REFERENCES devices(id)
```

---

## 25.9 Nullable Constraint

### 25.9.1 Owner

以下は NULL を許可してはならない。

```text
owner_user_id
```

---

### 25.9.2 Public ID

以下は NULL を許可してはならない。

```text
public_id
```

---

### 25.9.3 Credential Expiration

以下は NULL を許可してはならない。

```text
expires_at
```

対象:

- DeviceCredential
- GroupReadToken

---

## 25.10 Check Constraints

### 25.10.1 Address Family

許可値:

```sql
CHECK (
 family IN ('ipv4', 'ipv6')
)
```

---

### 25.10.2 Credential Permission

DeviceCredential:

```sql
CHECK (
 permission = 'write'
)
```

---

GroupReadToken:

```sql
CHECK (
 permission = 'read'
)
```

---

### 25.10.3 Metadata Length

実装は Metadata 制限を保証しなければならない。

例:

```sql
CHECK (
 octet_length(key) <= 256
)
```

---

## 25.11 Ownership Constraint

### 25.11.1 概要

以下の Invariant を保証しなければならない。

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

---

### 25.11.2 実装方法

以下のいずれでもよい。

- Transaction
- Trigger
- Constraint
- Application Logic

---

### 25.11.3 違反時

Membership 作成を拒否しなければならない。

---

## 25.12 Deletion Behavior

### 25.12.1 方針

削除時の挙動を明確にしなければならない。

---

### 25.12.2 Cascade

例:

```text
Device
→ DeviceMetadata
```

---

### 25.12.3 Non-Cascade

例:

```text
Device
↛ AuditLog
```

---

### 25.12.4 AuditLog

AuditLog は Cascade Delete してはならない。

---

### 25.12.5 UpdateHistory

UpdateHistory も Cascade Delete しないことを推奨する。

---

## 25.13 Transaction Requirements

### 25.13.1 Membership Create

以下は単一 Transaction を推奨する。

```text
Ownership Check
+
Membership Insert
```

---

### 25.13.2 Credential Revoke

失効処理は Atomic であるべきである。

---

### 25.13.3 Device Delete

関連削除は Transaction で処理することを推奨する。

---

## 25.14 Index Recommendations

### Device

```sql
INDEX(owner_user_id)
```

---

### Group

```sql
INDEX(owner_user_id)
```

---

### Credential

```sql
INDEX(token_prefix)
```

---

### Address

```sql
INDEX(expires_at)
```

---

### AuditLog

```sql
INDEX(created_at)
```

---

### UpdateHistory

```sql
INDEX(created_at)
```

---

## 25.15 不変条件

実装は以下を保証しなければならない。

```text
UNIQUE(device_id, family)
```

```text
UNIQUE(device_id, key)
```

```text
UNIQUE(management_group_id, device_id)
```

```text
UNIQUE(public_id)
```

```text
UNIQUE(type, token_prefix)
```

```text
owner_user_id
→ NOT NULL
```

```text
expires_at
→ NOT NULL
```

```text
family
∈ {ipv4, ipv6}
```

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

```text
AuditLog
→ MUST NOT Cascade Delete
```

これらを満たさない実装は本仕様に適合しない。
