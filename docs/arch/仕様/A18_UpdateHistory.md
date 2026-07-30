---
Title: 18. UpdateHistory
Description: 本章は「UpdateHistory」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「18. UpdateHistory」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 18. UpdateHistory

## 18.1 概要

本章では UpdateHistory の要件を定義する。

UpdateHistory は Resource の状態変化を追跡するための履歴である。

AuditLog が「誰が何をしたか」を記録するのに対し、UpdateHistory は「何がどのように変化したか」を記録する。

---

## 18.2 目的

UpdateHistory は以下の目的で保持する。

- Address 変更追跡
- Metadata 変更追跡
- 障害分析
- 運用分析
- 状態遷移の可視化
- データ監査

---

## 18.3 AuditLog との違い

### AuditLog

主目的:

```text
Who did what
```

例:

```text
User A updated Device X
```

---

### UpdateHistory

主目的:

```text
What changed
```

例:

```text
Address changed
192.0.2.1
↓
198.51.100.5
```

---

## 18.4 基本モデル

UpdateHistory は少なくとも以下を保持しなければならない。

| Field | Description |
|---------|---------|
| actor | 操作主体 |
| timestamp | 実行時刻 |
| event_type | イベント種別 |
| target_type | 対象種別 |
| target_id | 対象識別子 |
| result | 実行結果 |

---

## 18.5 Actor

以下をサポートすべきである。

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

## 18.6 Event Type

少なくとも以下をサポートすべきである。

```text
address_update
```

```text
metadata_create
```

```text
metadata_update
```

```text
metadata_delete
```

```text
address_delete
```

```text
credential_revoke
```

---

## 18.7 Address Update History

### 18.7.1 目的

Address の変更履歴を保持する。

---

### 18.7.2 保存項目

少なくとも以下を保持すべきである。

| Field | Description |
|---------|---------|
| actor | 操作者 |
| timestamp | 更新時刻 |
| device_id | Device |
| family | Address Family |
| old_address | 更新前 |
| new_address | 更新後 |
| old_expires_at | 更新前 TTL |
| new_expires_at | 更新後 TTL |
| mode | Update Mode |
| result | 結果 |
| reason | 理由 |

---

### 18.7.3 Address Value

Address Update History では Address Value を保存してよい。

---

### 18.7.4 Update Mode

以下を記録すべきである。

```text
update
```

```text
replace
```

```text
nop
```

---

## 18.8 Metadata Update History

### 18.8.1 目的

Metadata の変更追跡。

---

### 18.8.2 保存項目

少なくとも以下を保持すべきである。

| Field | Description |
|---------|---------|
| actor | 操作者 |
| timestamp | 更新時刻 |
| device_id | Device |
| key | Metadata Key |
| operation | 操作種別 |
| result | 結果 |
| reason | 理由 |

---

### 18.8.3 Operation

以下を使用する。

```text
create
```

```text
update
```

```text
delete
```

---

### 18.8.4 Value 保存禁止

Metadata Value は保存してはならない。

禁止:

```text
metadata_value
```

---

### 18.8.5 理由

以下のため。

- Privacy
- Storage Cost
- Security

---

## 18.9 Delete History

### 18.9.1 Device Delete

記録することを推奨する。

---

### 18.9.2 Address Delete

記録することを推奨する。

---

### 18.9.3 Metadata Delete

記録しなければならない。

---

### 18.9.4 Group Delete

記録することを推奨する。

---

## 18.10 Result

### Success

```text
success
```

---

### Failure

```text
failure
```

---

### Partial Success

実装は以下をサポートしてよい。

```text
partial_success
```

---

## 18.11 Timestamp

### 18.11.1 必須

全履歴は Timestamp を持たなければならない。

---

### 18.11.2 推奨

UTC を推奨する。

---

## 18.12 Retention

### 18.12.1 保持期間

UpdateHistory の保持期間は実装依存とする。

---

### 18.12.2 Cleanup

CleanupJob は古い UpdateHistory を削除してよい。

---

### 18.12.3 AuditLog との差異

AuditLog より短い保持期間としてもよい。

---

## 18.13 Read Access

### UserSession

参照を許可してよい。

---

### DeviceCredential

参照を許可してはならない。

---

### GroupReadToken

参照を許可してはならない。

---

## 18.14 Monitoring

実装は以下を提供してよい。

```text
update_history_events_total
```

```text
update_history_retained_records
```

```text
update_history_cleanup_total
```

---

## 18.15 不変条件

実装は以下を保証しなければならない。

```text
Address Update
→ History Recorded
```

```text
Metadata Create
→ History Recorded
```

```text
Metadata Update
→ History Recorded
```

```text
Metadata Delete
→ History Recorded
```

```text
Metadata Value
→ MUST NOT Be Stored
```

```text
Operation
∈ {create, update, delete}
```

```text
Timestamp
→ Required
```

```text
Result
→ Recorded
```

```text
UpdateHistory
≠ AuditLog
```

これらを満たさない実装は本仕様に適合しない。
