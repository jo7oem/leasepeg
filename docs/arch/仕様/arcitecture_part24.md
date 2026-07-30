---
Title: 24. Data Model
Description: 本章は「Data Model」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「24. Data Model」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 24. Data Model

## 24.1 概要

本章では本システムの論理データモデルを定義する。

本章は以下を目的とする。

- Entity 定義
- 永続化対象の明確化
- SQL 制約の基礎定義
- 実装間の互換性確保

物理スキーマは実装依存であるが、本章の論理モデルを満たさなければならない。

---

## 24.2 共通要件

### 24.2.1 Internal ID

各テーブルは内部識別子を持ってよい。

例:

```text
id BIGSERIAL
```

---

### 24.2.2 Public ID

外部公開対象 Resource は UUIDv7 を使用しなければならない。

例:

```text
public_id UUIDv7
```

---

### 24.2.3 Timestamp

以下を保持することを推奨する。

```text
created_at
```

```text
updated_at
```

---

## 24.3 users

### 目的

利用者情報。

### モデル

```text
users

id
public_id UUIDv7

name

created_at
updated_at
```

---

### 制約

```text
public_id UNIQUE
```

---

## 24.4 user_sessions

### 目的

UserSession 管理。

### モデル

```text
user_sessions

id
public_id UUIDv7

user_id

secret_hash

expires_at
revoked_at

created_at
last_used_at
```

---

### 制約

```text
user_id
→ users.id
```

---

### Security

```text
secret_hash
```

は平文保存してはならない。

---

## 24.5 devices

### 目的

Device 管理。

### モデル

```text
devices

id
public_id UUIDv7

owner_user_id

name

created_at
updated_at
```

---

### 備考

以下は使用しない。

```text
devices.enabled
```

利用目的が定義されていないため削除する。

---

### 制約

```text
owner_user_id
→ users.id
```

---

## 24.6 management_groups

### 目的

Device グループ管理。

### モデル

```text
management_groups

id
public_id UUIDv7

owner_user_id

name

created_at
updated_at
```

---

### 制約

```text
owner_user_id
→ users.id
```

---

## 24.7 management_group_devices

### 目的

ManagementGroup Membership。

### モデル

```text
management_group_devices

management_group_id
device_id

created_at
```

---

### 制約

```text
UNIQUE(
  management_group_id,
  device_id
)
```

---

### Ownership Rule

以下を満たさなければならない。

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

---

## 24.8 device_credentials

### 目的

Device 認証情報。

### モデル

```text
device_credentials

id
public_id UUIDv7

device_id

type
permission

secret_hash

token_prefix

expires_at
revoked_at

last_used_at

created_at
```

---

### type

以下を想定する。

```text
bearer
basic
digest
```

---

### permission

```text
write
```

のみ。

---

### 制約

```text
device_id
→ devices.id
```

---

### Security

```text
secret_hash
```

は平文を保存してはならない。

---

## 24.9 group_read_tokens

### 目的

ManagementGroup Read Token。

### モデル

```text
group_read_tokens

id
public_id UUIDv7

management_group_id

secret_hash

token_prefix

expires_at
revoked_at

last_used_at

created_at
```

---

### permission

```text
read
```

固定。

---

### 制約

```text
management_group_id
→ management_groups.id
```

---

## 24.10 device_addresses

### 目的

Device Address 保存。

### モデル

```text
device_addresses

device_id

family
address

expires_at

last_seen_at

created_at
updated_at
```

---

### Family

以下のみ。

```text
ipv4
```

```text
ipv6
```

---

### 制約

```text
UNIQUE(device_id, family)
```

---

## 24.11 device_metadata

### 目的

Metadata 保存。

### モデル

```text
device_metadata

device_id

key
value

created_at
updated_at
```

---

### 制約

```text
UNIQUE(device_id, key)
```

---

### 制限

```text
max metadata count = 3
```

---

### Value

最大:

```text
128 KiB
```

---

## 24.12 audit_logs

### 目的

監査ログ保存。

### モデル

```text
audit_logs

id
public_id UUIDv7

request_id

actor_type
actor_id

target_type
target_id

device_id
credential_id
management_group_id

observed_ip
submitted_ip

endpoint_type

mode

user_agent

result
reason

created_at
```

---

### Security

以下を保存してはならない。

```text
Bearer Token
```

```text
Password
```

```text
Digest Secret
```

```text
Session Secret
```

---

## 24.13 update_histories

### 目的

変更履歴保存。

### モデル

```text
update_histories

id
public_id UUIDv7

actor_type
actor_id

event_type

target_type
target_id

device_id

family

old_address
new_address

old_expires_at
new_expires_at

metadata_key
metadata_operation

result
reason

created_at
```

---

### Metadata 制約

保存してよい:

```text
metadata_key
```

```text
metadata_operation
```

---

保存禁止:

```text
metadata_value
```

---

## 24.14 integration_events

### 目的

External Integration Queue。

### モデル

```text
integration_events

id
public_id UUIDv7

event_type

target_type
target_id

payload

status

retry_count

last_error

created_at
updated_at
```

---

### Status

例:

```text
pending
```

```text
processing
```

```text
success
```

```text
failed
```

```text
dead_letter
```

---

## 24.15 cleanup_job_runs

### 目的

CleanupJob 実行履歴。

### モデル

```text
cleanup_job_runs

id
public_id UUIDv7

job_name

started_at
finished_at

result

deleted_record_count

failure_reason

triggered_by
```

---

## 24.16 Entity Relationship

### User

```text
User
 ├─ UserSession
 ├─ Device
 └─ ManagementGroup
```

---

### Device

```text
Device
 ├─ DeviceCredential
 ├─ DeviceAddress
 └─ DeviceMetadata
```

---

### ManagementGroup

```text
ManagementGroup
 ├─ Device Membership
 └─ GroupReadToken
```

---

## 24.17 Ownership Invariants

実装は以下を保証しなければならない。

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

---

```text
DeviceCredential
→ Exactly One Device
```

---

```text
GroupReadToken
→ Exactly One Group
```

---

```text
UserSession
→ Exactly One User
```

---

## 24.18 Data Model Invariants

実装は以下を保証しなければならない。

```text
public_id
→ UUIDv7
```

```text
DeviceAddress
→ <= 1 row per family
```

```text
Metadata Count
<= 3
```

```text
Metadata Value
<= 128 KiB
```

```text
Credential Secret
→ Hashed
```

```text
AuditLog
→ No Secret Storage
```

```text
UpdateHistory
→ No Metadata Value Storage
```

```text
Ownership Constraint
→ Enforced
```

これらを満たさない実装は本仕様に適合しない。
