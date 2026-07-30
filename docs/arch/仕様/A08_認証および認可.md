---
Title: 8. 認証および認可
Description: 本章は「認証および認可」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「8. 認証および認可」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 8. 認証および認可

## 8.1 概要

本章では、本システムにおける認証および認可モデルを定義する。

本システムは以下の認証主体を持つ。

- UserSession
- DeviceCredential
- GroupReadToken

実装は認証と認可を明確に分離しなければならない。

```text
Authentication
↓
Authorization
↓
Operation
```

認証成功は認可成功を意味しない。

---

## 8.2 UserSession

### 8.2.1 目的

UserSession は人間利用者による管理操作のための Credential である。

---

### 8.2.2 権限

UserSession は以下を実行できる。

- Device Create
- Device Update
- Device Delete
- ManagementGroup Create
- ManagementGroup Update
- ManagementGroup Delete
- DeviceCredential Issue
- DeviceCredential Revoke
- GroupReadToken Issue
- GroupReadToken Revoke
- Address Read
- Metadata Read
- AuditLog Read
- Cleanup API Call

---

### 8.2.3 制限

UserSession は以下を満たさなければならない。

```text
maximum lifetime = 30 minutes
```

---

### 8.2.4 長寿命化禁止

User 権限を持つ永続 Credential を発行してはならない。

以下は禁止とする。

```text
Permanent User Token
```

```text
API Key with User Privilege
```

---

### 8.2.5 Cookie Security

Cookie によって UserSession を管理する場合、以下を満たさなければならない。

```text
Secure
HttpOnly
SameSite=Lax or stricter
```

---

### 8.2.6 CSRF Protection

状態変更を伴う API は CSRF 対策を実施しなければならない。

対象:

- POST
- PUT
- PATCH
- DELETE

---

## 8.3 DeviceCredential

### 8.3.1 目的

DeviceCredential は Device が自身の情報を更新するための Credential である。

---

### 8.3.2 対応方式

実装は以下をサポートしてよい。

```text
Bearer Token
```

```text
Basic Authentication
```

```text
Digest Authentication
```

---

### 8.3.3 権限

DeviceCredential は以下のみ実行できる。

```text
Address Update
Metadata Update
```

---

### 8.3.4 禁止事項

以下を許可してはならない。

```text
Address Read
```

```text
Group Read
```

```text
Device Management
```

```text
Credential Management
```

```text
User Management
```

---

### 8.3.5 Device Binding

DeviceCredential は単一 Device に紐付かなければならない。

```text
DeviceCredential
→ One Device
```

---

### 8.3.6 Device 解決

Device Update API では Request に Device ID を含めてはならない。

誤:

```json
{
  "device_id": "..."
}
```

正:

```text
Credential
↓
Device Resolve
↓
Update
```

---

### 8.3.7 Revocation

実装は Credential の即時失効をサポートしなければならない。

---

### 8.3.8 Expiration

以下を満たさなければならない。

```text
expires_at != NULL
```

---

### 8.3.9 Last Used

以下を記録しなければならない。

```text
last_used_at
```

---

### 8.3.10 Secret Storage

Secret を平文保存してはならない。

---

## 8.4 GroupReadToken

### 8.4.1 目的

ManagementGroup 内 Device の Address 参照を行う。

---

### 8.4.2 権限

許可:

```text
Address Read
```

---

禁止:

```text
Address Update
```

```text
Metadata Update
```

```text
Credential Management
```

```text
User Management
```

---

### 8.4.3 Group Binding

以下を満たさなければならない。

```text
GroupReadToken
→ One ManagementGroup
```

---

### 8.4.4 Revocation

失効は即時反映されなければならない。

---

### 8.4.5 Expiration

以下を満たさなければならない。

```text
expires_at != NULL
```

---

### 8.4.6 Last Used

以下を保持しなければならない。

```text
last_used_at
```

---

## 8.5 Credential Lifetime Policy

### 8.5.1 基本要件

以下の Credential は有効期限を持たなければならない。

- DeviceCredential
- GroupReadToken

---

### 8.5.2 Default Cap

既定値:

```text
365 days
```

---

### 8.5.3 Maximum Cap

実装は以下を上限としてよい。

```text
99 years
```

---

### 8.5.4 Expired Credential

期限切れ Credential は認証成功として扱ってはならない。

---

### 8.5.5 Revoked Credential

失効 Credential は認証成功として扱ってはならない。

---

## 8.6 Token Prefix

### 8.6.1 目的

Credential 識別補助として利用する。

---

### 8.6.2 要件

以下を満たさなければならない。

```text
>= 16 random bytes
```

または同等以上の衝突耐性を持つ値。

---

### 8.6.3 一意性

以下を満たさなければならない。

```text
UNIQUE(type, token_prefix)
```

---

### 8.6.4 利用用途

利用してよい用途:

- Audit
- Rate Limit
- Credential Lookup

---

### 8.6.5 禁止用途

以下は禁止する。

```text
Token Prefix Only Authentication
```

---

### 8.6.6 Error Response

認証失敗時に以下を返してはならない。

```text
token_prefix exists
```

```text
token_prefix invalid
```

---

## 8.7 Secret Storage Policy

### 8.7.1 基本方針

認証 Secret は平文保存してはならない。

---

### 8.7.2 Bearer Token

保存方式:

```text
Hash
```

---

### 8.7.3 Basic Password

保存方式:

```text
Hash
```

---

### 8.7.4 UserSession Secret

保存方式:

```text
Hash
```

---

### 8.7.5 Token Display

Token は発行時のみ表示してよい。

発行後に Secret を再表示してはならない。

---

## 8.8 Authentication Failure

### 8.8.1 情報秘匿

認証失敗時に以下を漏洩してはならない。

- Credential Exists
- Device Exists
- Token Exists
- Expired
- Revoked

---

### 8.8.2 共通応答

実装はなるべく統一的な応答を返すべきである。

例:

```http
401 Unauthorized
```

---

## 8.9 Authorization Requirements

### 8.9.1 Device

Device Credential は自身の Device のみ更新できなければならない。

---

### 8.9.2 Group

GroupReadToken は自身の Group のみ参照できなければならない。

---

### 8.9.3 User

UserSession は自身が所有する Resource のみ操作できなければならない。

---

## 8.10 不変条件

実装は以下を保証しなければならない。

```text
UserSession
→ User Privilege
```

```text
DeviceCredential
→ Write Only
```

```text
GroupReadToken
→ Read Only
```

```text
User Privilege Credential
→ Not Long-Lived
```

```text
Credential
→ Revocable
```

```text
Credential
→ Expirable
```

```text
Authentication
≠ Authorization
```

これらの条件を満たさない実装は本仕様に適合しない。
