---
Title: 14. Metadata Model
Description: 本章は「Metadata Model」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「14. Metadata Model」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 14. Metadata Model

## 14.1 概要

本章では Device に付随する Metadata のデータモデルを定義する。

Metadata は Device に対する補足情報を格納するために使用する。

Metadata は以下の特徴を持つ。

- Key-Value 形式
- UTF-8 対応
- 小容量データ向け
- Device 単位
- 履歴記録対応

---

## 14.2 Metadata の目的

Metadata は Device の補助的な管理情報を保持する。

例:

```text
device_name
```

```text
location
```

```text
description
```

```text
owner_note
```

---

## 14.3 DeviceMetadata

### 14.3.1 責務

DeviceMetadata は Device に紐付く任意情報を保持する。

---

### 14.3.2 必須属性

少なくとも以下を持つ。

| Field | Description |
|---------|---------|
| device_id | Device |
| key | Metadata Key |
| value | Metadata Value |
| created_at | Creation Time |
| updated_at | Update Time |

---

## 14.4 Key

### 14.4.1 形式

Key は UTF-8 String とする。

---

### 14.4.2 最大長

以下を満たさなければならない。

```text
<= 256 bytes
```

---

### 14.4.3 大小文字

実装は以下を定義しなければならない。

```text
Case Sensitive
```

または

```text
Case Insensitive
```

同一実装内で一貫していなければならない。

---

### 14.4.4 空文字列

空文字列 Key を許可してはならない。

---

## 14.5 Value

### 14.5.1 形式

Value は UTF-8 String とする。

---

### 14.5.2 最大長

以下を満たさなければならない。

```text
<= 128 KiB
```

---

### 14.5.3 空文字列

以下は有効な値である。

```text
""
```

---

### 14.5.4 Null

以下は削除を意味する。

```json
{
  "key": null
}
```

---

## 14.6 Metadata 数制限

### 14.6.1 Device 単位

1 Device あたり最大数:

```text
3 Key-Value Pairs
```

---

### 14.6.2 超過時

制限超過時は Request を拒否しなければならない。

推奨応答:

```http
400 Bad Request
```

---

## 14.7 一意性制約

### 14.7.1 Device 内一意性

以下を保証しなければならない。

```text
UNIQUE(device_id, key)
```

---

### 14.7.2 例

許可:

```text
name
location
note
```

---

禁止:

```text
name
name
```

---

## 14.8 Metadata 作成

### 14.8.1 新規追加

存在しない Key の場合は作成する。

---

### 14.8.2 Audit

作成操作は AuditLog に記録すべきである。

---

### 14.8.3 UpdateHistory

作成操作は UpdateHistory に記録すべきである。

---

## 14.9 Metadata 更新

### 14.9.1 更新条件

既存 Key が存在する場合は更新する。

---

### 14.9.2 上書き

新しい Value に置き換える。

---

### 14.9.3 空文字列

以下は有効な更新である。

```json
{
  "note": ""
}
```

---

## 14.10 Metadata 削除

### 14.10.1 削除方法

以下を削除として扱う。

```json
{
  "note": null
}
```

---

### 14.10.2 結果

Metadata Record を削除してよい。

---

### 14.10.3 Audit

削除操作は AuditLog に記録すべきである。

---

## 14.11 更新セマンティクス

### 14.11.1 例

入力:

```json
{
  "name": "router-a",
  "note": "",
  "old_key": null
}
```

---

結果:

```text
name
→ "router-a"
```

```text
note
→ ""
```

```text
old_key
→ deleted
```

---

## 14.12 Metadata History

### 14.12.1 基本原則

Metadata Value は履歴へ保存してはならない。

---

### 14.12.2 理由

以下のため。

- Privacy
- Security
- Storage Cost

---

### 14.12.3 保存対象

保存してよい情報:

```text
device_id
```

```text
key
```

```text
operation
```

```text
timestamp
```

```text
actor
```

---

### 14.12.4 保存禁止

保存してはならない。

```text
value
```

---

## 14.13 UpdateHistory

### 14.13.1 Operation Type

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

### 14.13.2 例

```json
{
  "operation": "create",
  "key": "name"
}
```

---

```json
{
  "operation": "update",
  "key": "name"
}
```

---

```json
{
  "operation": "delete",
  "key": "note"
}
```

---

## 14.14 Authorization

### 14.14.1 UserSession

自身が所有する Device の Metadata のみ操作できる。

---

### 14.14.2 DeviceCredential

自身の Device の Metadata のみ更新できる。

---

### 14.14.3 GroupReadToken

Metadata の更新を許可してはならない。

---

## 14.15 Monitoring

実装は以下を提供してよい。

```text
metadata_entries_total
```

```text
metadata_updates_total
```

```text
metadata_update_failures_total
```

---

## 14.16 Error Handling

### Invalid Key

以下の場合は拒否する。

```text
> 256 bytes
```

---

### Invalid Value

以下の場合は拒否する。

```text
> 128 KiB
```

---

### Pair Limit

以下の場合は拒否する。

```text
> 3 Metadata Entries
```

---

## 14.17 不変条件

実装は以下を保証しなければならない。

```text
Key
→ UTF-8
```

```text
Key Length
<= 256 bytes
```

```text
Value Length
<= 128 KiB
```

```text
Metadata Count Per Device
<= 3
```

```text
UNIQUE(device_id, key)
```

```text
Null
→ Delete
```

```text
Empty String
→ Valid Value
```

```text
Metadata History
→ MUST NOT Store Value
```

```text
UpdateHistory.operation
∈ {create, update, delete}
```

これらを満たさない実装は本仕様に適合しない。
