---
Title: 26. OpenAPI Security Schemes
Description: 本章は「OpenAPI Security Schemes」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「26. OpenAPI Security Schemes」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 26. OpenAPI Security Schemes

## 26.1 概要

本章では OpenAPI 定義における Security Scheme を定義する。

本システムは認証主体ごとに異なる Security Scheme を使用する。

対象は以下とする。

- DeviceCredential
- GroupReadToken
- Basic Authentication
- Digest Authentication
- UserSession

---

## 26.2 設計原則

### 26.2.1 権限分離

認証方式ごとに権限を分離しなければならない。

---

### 26.2.2 OpenAPI 明示

各 API は利用可能な Security Scheme を明示しなければならない。

---

### 26.2.3 最小権限

Security Scheme は必要最小限の権限のみを持つべきである。

---

## 26.3 Device Bearer Authentication

### 26.3.1 用途

DeviceCredential 用。

---

### 26.3.2 権限

許可:

```text
Address Update
Metadata Update
```

---

### 26.3.3 禁止

```text
Address Read
```

```text
Management Operation
```

---

### 26.3.4 OpenAPI 定義例

```yaml
deviceBearerAuth:
  type: http
  scheme: bearer
  bearerFormat: opaque-device-write-token
```

---

## 26.4 Group Bearer Authentication

### 26.4.1 用途

GroupReadToken 用。

---

### 26.4.2 権限

許可:

```text
Address Read
```

---

### 26.4.3 禁止

```text
Address Update
```

```text
Management Operation
```

---

### 26.4.4 OpenAPI 定義例

```yaml
groupBearerAuth:
  type: http
  scheme: bearer
  bearerFormat: opaque-group-read-token
```

---

## 26.5 Basic Authentication

### 26.5.1 用途

レガシー Router 互換性のために提供してよい。

---

### 26.5.2 権限

DeviceCredential と同等。

---

### 26.5.3 OpenAPI 定義例

```yaml
basicAuth:
  type: http
  scheme: basic
```

---

### 26.5.4 Security Requirement

Password は Hash 保存しなければならない。

---

## 26.6 Digest Authentication

### 26.6.1 用途

レガシー Device 互換性のために提供してよい。

---

### 26.6.2 権限

DeviceCredential と同等。

---

### 26.6.3 OpenAPI 定義例

```yaml
digestAuth:
  type: http
  scheme: digest
```

---

### 26.6.4 Storage Requirement

Digest Credential の保存方式は実装ごとに明確化しなければならない。

---

## 26.7 User Session Authentication

### 26.7.1 用途

User Management API。

---

### 26.7.2 権限

許可:

```text
Device Management
```

```text
Group Management
```

```text
Credential Management
```

```text
AuditLog Read
```

---

### 26.7.3 OpenAPI 定義例

```yaml
userSessionAuth:
  type: apiKey
  in: cookie
  name: session
```

---

## 26.8 Cookie Security Requirements

### 26.8.1 Secure

```text
Secure
```

必須。

---

### 26.8.2 HttpOnly

```text
HttpOnly
```

必須。

---

### 26.8.3 SameSite

以下を必須とする。

```text
SameSite=Lax
```

以上。

---

## 26.9 Security Requirement Mapping

### Device Update API

```yaml
security:
  - deviceBearerAuth: []
```

---

### Group Read API

```yaml
security:
  - groupBearerAuth: []
```

---

### User Management API

```yaml
security:
  - userSessionAuth: []
```

---

### Legacy Device API

```yaml
security:
  - basicAuth: []
```

または

```yaml
security:
  - digestAuth: []
```

---

## 26.10 Multiple Scheme Support

### 26.10.1 許容

同一 Endpoint が複数認証方式をサポートしてもよい。

---

### 26.10.2 Device API

例:

```yaml
security:
  - deviceBearerAuth: []
  - basicAuth: []
  - digestAuth: []
```

---

### 26.10.3 権限

認証方式が異なっても権限は同一でなければならない。

---

## 26.11 Authorization Rules

認証成功後は必ず Authorization を実施しなければならない。

---

### DeviceCredential

```text
Bound Device Only
```

---

### GroupReadToken

```text
Bound Group Only
```

---

### UserSession

```text
Owned Resources Only
```

---

## 26.12 OpenAPI Documentation Requirements

OpenAPI 定義は少なくとも以下を説明すべきである。

- Authentication Method
- Required Headers
- Cookie Requirements
- Authorization Requirements
- Error Responses

---

## 26.13 Error Handling

認証失敗時の推奨応答。

```http
401 Unauthorized
```

---

認可失敗時の推奨応答。

```http
403 Forbidden
```

---

認証失敗時に以下を返してはならない。

```text
Credential Exists
```

```text
Credential Revoked
```

```text
Credential Expired
```

---

## 26.14 Security Invariants

実装は以下を保証しなければならない。

```text
DeviceCredential
→ Write Only
```

```text
GroupReadToken
→ Read Only
```

```text
UserSession
→ User Privilege
```

```text
Authentication
≠ Authorization
```

```text
UserSession Cookie
→ Secure
```

```text
UserSession Cookie
→ HttpOnly
```

```text
UserSession Cookie
→ SameSite=Lax+
```

```text
Basic Password
→ Hashed
```

```text
Bearer Token
→ Hashed
```

```text
Digest Credential
→ Protected
```

これらを満たさない実装は本仕様に適合しない。
