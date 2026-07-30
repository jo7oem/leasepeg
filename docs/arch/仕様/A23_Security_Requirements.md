---
Title: 23. Security Requirements
Description: 本章は「Security Requirements」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「23. Security Requirements」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 23. Security Requirements

## 23.1 概要

本章では本システムのセキュリティ要件を定義する。

本システムは以下を基本原則とする。

- Least Privilege
- Defense in Depth
- Secure by Default
- Auditability
- Credential Protection

---

## 23.2 Transport Security

### 23.2.1 HTTPS 必須

Production 環境では HTTPS を使用しなければならない。

```text
HTTPS
→ MUST
```

---

### 23.2.2 HTTP 禁止

認証情報を平文 HTTP で送信してはならない。

```text
HTTP
→ MUST NOT
```

---

### 23.2.3 TLS Termination

TLS Termination は本システム外部で実施してもよい。

例:

- Nginx
- Envoy
- HAProxy
- Load Balancer

---

## 23.3 Secret Storage

### 23.3.1 平文保存禁止

認証情報を平文保存してはならない。

---

### 23.3.2 Bearer Token

保存形式:

```text
Hash
```

---

### 23.3.3 Basic Authentication

Password は Hash 保存しなければならない。

---

### 23.3.4 UserSession

Session Secret は Hash 保存しなければならない。

---

### 23.3.5 Digest Authentication

保存方式は実装ごとに明確化しなければならない。

---

### 23.3.6 発行後表示

Credential Secret は発行時のみ表示してよい。

以後再表示してはならない。

---

## 23.4 Credential Security

### 23.4.1 Expiration

Credential は有効期限を持たなければならない。

対象:

- DeviceCredential
- GroupReadToken

---

### 23.4.2 Revocation

Credential は即時失効できなければならない。

---

### 23.4.3 Last Used

以下を保持しなければならない。

```text
last_used_at
```

---

### 23.4.4 Token Length

Bearer Token は十分な長さを持たなければならない。

推奨:

```text
>= 32 Characters
```

---

## 23.5 Token Prefix

### 23.5.1 要件

Token Prefix は衝突耐性を持たなければならない。

推奨:

```text
>= 16 Random Bytes
```

---

### 23.5.2 認証禁止

Token Prefix のみで認証してはならない。

---

### 23.5.3 Error Leakage

認証失敗時に Token Prefix の有効性を返してはならない。

---

## 23.6 Authentication

### 23.6.1 DeviceCredential

DeviceCredential は自身の Device のみ操作できなければならない。

---

### 23.6.2 GroupReadToken

GroupReadToken は読み取り専用でなければならない。

---

### 23.6.3 UserSession

UserSession のみ User 権限を持てる。

---

### 23.6.4 Long-Lived User Credential

発行してはならない。

---

## 23.7 Authorization

### 23.7.1 必須要件

すべての Non-Public API は Authorization を実施しなければならない。

---

### 23.7.2 UUIDv7

UUIDv7 を利用していても Authorization を省略してはならない。

---

### 23.7.3 Ownership Verification

以下を検証しなければならない。

```text
device_id
```

```text
group_id
```

```text
token_id
```

---

### 23.7.4 Ownership Invariant

以下を保証しなければならない。

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

---

## 23.8 Input Validation

### 23.8.1 IP Address

IP Address Format を検証しなければならない。

---

### 23.8.2 Address Family

検証対象:

```text
ipv4
```

```text
ipv6
```

---

### 23.8.3 Update Mode

検証対象:

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

### 23.8.4 Metadata Key

以下を検証しなければならない。

- UTF-8 Validity
- Length

---

### 23.8.5 Metadata Value

以下を検証しなければならない。

- UTF-8 Validity
- Length

---

### 23.8.6 Request Size

Request Body Size を制限しなければならない。

---

### 23.8.7 Content-Type

Content-Type を検証しなければならない。

---

## 23.9 Session Security

### 23.9.1 Secure

Cookie は Secure 属性を持たなければならない。

---

### 23.9.2 HttpOnly

Cookie は HttpOnly 属性を持たなければならない。

---

### 23.9.3 SameSite

Cookie は以下を満たさなければならない。

```text
SameSite=Lax
```

またはそれ以上。

---

### 23.9.4 Expiration

Session は短寿命でなければならない。

```text
<= 30 Minutes
```

---

## 23.10 CSRF Protection

### 23.10.1 対象

状態変更 API。

例:

```http
POST
PUT
PATCH
DELETE
```

---

### 23.10.2 要件

CSRF 対策を実施しなければならない。

---

## 23.11 Error Handling

### 23.11.1 情報漏洩防止

以下を漏洩してはならない。

- Credential Exists
- Credential Revoked
- Credential Expired
- User Exists
- Device Exists

---

### 23.11.2 共通応答

認証失敗時は統一的な応答を推奨する。

---

## 23.12 Audit Security

### 23.12.1 Logging Required

以下を記録しなければならない。

- Authentication Failure
- Authorization Failure
- Credential Revoke
- Delete Operation

---

### 23.12.2 Logging Prohibited

以下を記録してはならない。

```text
Password
```

```text
Token Secret
```

```text
Session Secret
```

```text
Digest Secret
```

---

## 23.13 External Integration Security

### 23.13.1 Timeout

External API Call は Timeout を持たなければならない。

---

### 23.13.2 Secret Protection

External Integration Secret は平文保存してはならない。

---

### 23.13.3 Isolation

外部システム障害によって内部状態を破壊してはならない。

---

## 23.14 Monitoring Security

### 23.14.1 Monitoring API

既定で外部公開してはならない。

---

### 23.14.2 Internal API

既定で外部公開してはならない。

---

### 23.14.3 Separate Port

Monitoring API は別 Port を使用しなければならない。

---

## 23.15 Security Invariants

実装は以下を保証しなければならない。

```text
Secrets
→ Never Stored In Plaintext
```

```text
Authentication
→ Required
```

```text
Authorization
→ Required
```

```text
UserSession
→ Short Lived
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
Password
→ Hashed
```

```text
Token
→ Hashed
```

```text
CSRF Protection
→ Enabled
```

```text
Monitoring API
→ Not Public By Default
```

```text
Internal API
→ Not Public By Default
```

```text
Authentication Failure
→ No Information Disclosure
```

```text
Authorization Failure
→ No Information Disclosure
```

これらを満たさない実装は本仕様に適合しない。
