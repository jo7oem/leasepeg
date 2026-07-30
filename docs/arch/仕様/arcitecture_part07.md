---
Title: 7. 識別子および所有権制約
Description: 本章は「識別子および所有権制約」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「7. 識別子および所有権制約」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 7. 識別子および所有権制約

## 7.1 概要

本章では、本システムにおける識別子の取り扱いおよび所有権制約を定義する。

実装は以下を保証しなければならない。

- Resource Identifier の一貫性
- 所有権の明確化
- 認可チェックの強制
- UUIDv7 の利用
- 所有者不一致の防止

---

## 7.2 識別子ポリシー

### 7.2.1 基本原則

外部公開されるすべての Resource は公開用 Identifier を持たなければならない。

実装は内部 ID を利用者へ公開してはならない。

---

### 7.2.2 公開識別子

公開識別子には UUIDv7 を使用する。

例:

```text
018f6f39-3f5c-7d2c-85c2-cabb7de57992
```

---

### 7.2.3 UUIDv7 採用理由

UUIDv7 を採用する理由は以下である。

- 推測が困難
- 時系列ソート可能
- 分散生成可能
- 高い一意性

---

## 7.3 内部識別子

### 7.3.1 Internal ID

Database 内部では整数 ID を使用してもよい。

例:

```text
id BIGSERIAL
```

---

### 7.3.2 非公開原則

Internal ID は API を通じて公開してはならない。

誤:

```json
{
  "id": 42
}
```

正:

```json
{
  "public_id": "018f6f39-3f5c-7d2c-85c2-cabb7de57992"
}
```

---

## 7.4 Resource Identifier

### 7.4.1 User

```text
user_public_id
```

型:

```text
UUIDv7
```

---

### 7.4.2 Device

```text
device_public_id
```

型:

```text
UUIDv7
```

---

### 7.4.3 ManagementGroup

```text
management_group_public_id
```

型:

```text
UUIDv7
```

---

### 7.4.4 Credential

```text
credential_public_id
```

型:

```text
UUIDv7
```

---

### 7.4.5 AuditLog

```text
audit_log_public_id
```

型:

```text
UUIDv7
```

---

### 7.4.6 UpdateHistory

```text
update_history_public_id
```

型:

```text
UUIDv7
```

---

## 7.5 所有権モデル

すべての Resource は単一 User が所有する。

---

### 7.5.1 Device

```text
Device
 └─ Owner User
```

---

### 7.5.2 ManagementGroup

```text
ManagementGroup
 └─ Owner User
```

---

### 7.5.3 DeviceCredential

```text
DeviceCredential
 └─ Device
      └─ User
```

所有権は Device を経由して判定する。

---

### 7.5.4 GroupReadToken

```text
GroupReadToken
 └─ ManagementGroup
      └─ User
```

所有権は ManagementGroup を経由して判定する。

---

## 7.6 所有権判定

### 7.6.1 基本原則

認証成功後、必ず所有権判定を実施しなければならない。

以下を実行しなければならない。

```text
Authentication
↓
Authorization
↓
Operation
```

---

### 7.6.2 禁止事項

以下は許可してはならない。

```text
Authentication Only
```

---

### 7.6.3 UUID の誤用禁止

UUIDv7 を利用していても認可を省略してはならない。

誤った考え方:

```text
UUID が推測できない
↓
認可不要
```

正しい設計:

```text
UUIDv7
+
Authorization Check
```

---

## 7.7 Device 所有権制約

以下を満たさなければならない。

```text
Device.owner_user_id
!= NULL
```

---

### 7.7.1 Device 作成

Device 作成時には必ず Owner User を設定する。

---

### 7.7.2 Device 更新

Device 更新時は Owner User を変更してはならない。

---

## 7.8 ManagementGroup 所有権制約

### 7.8.1 Owner 必須

以下を満たさなければならない。

```text
ManagementGroup.owner_user_id
!= NULL
```

---

### 7.8.2 更新制約

Group Owner の変更は明示的操作によってのみ許可してよい。

---

## 7.9 Membership 制約

### 7.9.1 基本条件

ManagementGroup に Device を追加する場合、以下を満たさなければならない。

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

---

### 7.9.2 拒否条件

以下は拒否しなければならない。

```text
User A
 └─ Device X

User B
 └─ ManagementGroup Y
```

```text
Device X
↓
Group Y
```

---

### 7.9.3 Authorization Failure

所有権が一致しない場合は認可エラーとしなければならない。

---

## 7.10 Credential 所有権制約

### 7.10.1 DeviceCredential

以下を満たさなければならない。

```text
DeviceCredential
→ exactly one Device
```

---

### 7.10.2 GroupReadToken

以下を満たさなければならない。

```text
GroupReadToken
→ exactly one ManagementGroup
```

---

### 7.10.3 UserSession

以下を満たさなければならない。

```text
UserSession
→ exactly one User
```

---

## 7.11 API における識別子

### 7.11.1 Request

API Request は公開識別子を使用する。

例:

```http
GET /v1/devices/{device_public_id}
```

---

### 7.11.2 Response

API Response は公開識別子のみを返す。

例:

```json
{
  "public_id": "018f6f39-3f5c-7d2c-85c2-cabb7de57992"
}
```

---

### 7.11.3 Internal ID 非公開

以下を返してはならない。

```json
{
  "id": 42
}
```

---

## 7.12 Audit 対象

以下のイベントは監査対象とする。

- 所有権違反
- Authorization Failure
- Cross User Access Attempt
- Invalid Membership Update

---

## 7.13 不変条件

実装は以下の条件を常に保証しなければならない。

### User Ownership

```text
Device.owner_user_id != NULL
```

```text
ManagementGroup.owner_user_id != NULL
```

---

### Membership Ownership

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

---

### Credential Binding

```text
DeviceCredential
→ single Device
```

```text
GroupReadToken
→ single ManagementGroup
```

```text
UserSession
→ single User
```

---

### Public Identifier

```text
public_id
∈ UUIDv7
```

---

### Authorization

```text
Authentication
≠ Authorization
```

認証成功は認可成功を意味しない。

認可チェックを省略してはならない。
