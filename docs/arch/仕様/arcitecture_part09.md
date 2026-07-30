# 9. 認証失敗処理および Rate Limit

## 9.1 概要

本章では、認証失敗時の振る舞い、および DoS 攻撃や Credential 推測攻撃を緩和するための Rate Limit 要件を定義する。

実装は以下を保証しなければならない。

- 情報漏洩防止
- Credential Enumeration 防止
- Brute Force 攻撃耐性
- Rate Limit による保護
- Audit 記録

---

## 9.2 認証失敗時の応答原則

### 9.2.1 基本要件

認証失敗時の応答は情報を漏洩してはならない。

以下を推測できてはならない。

- Credential の存在
- Device の存在
- Group の存在
- User の存在
- Token Prefix の存在

---

### 9.2.2 応答統一

以下の状況で応答差異を作ってはならない。

```text
Credential Not Found
```

```text
Credential Revoked
```

```text
Credential Expired
```

```text
Credential Invalid
```

---

### 9.2.3 許可される応答例

```http
401 Unauthorized
```

```json
{
  "error": "authentication_failed"
}
```

---

### 9.2.4 禁止例

以下のような応答を返してはならない。

```json
{
  "error": "device_not_found"
}
```

```json
{
  "error": "token_expired"
}
```

```json
{
  "error": "credential_revoked"
}
```

---

## 9.3 Credential Enumeration 防止

### 9.3.1 要件

攻撃者が Resource 存在有無を推測できてはならない。

---

### 9.3.2 非公開情報

以下は非公開情報として扱う。

- Device Exists
- Group Exists
- Credential Exists
- User Exists

---

### 9.3.3 実装要件

認証失敗時は同一カテゴリのエラーとして扱うべきである。

---

## 9.4 Authentication Failure Logging

### 9.4.1 記録対象

認証失敗は AuditLog に記録しなければならない。

---

### 9.4.2 記録項目

少なくとも以下を保存すべきである。

```text
timestamp
source_ip
credential_type
token_prefix
result
```

---

### 9.4.3 Secret Logging 禁止

以下を記録してはならない。

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

## 9.5 Rate Limit の目的

Rate Limit は以下を防止するために存在する。

### Brute Force

```text
Credential Guessing
```

---

### Enumeration

```text
Resource Discovery
```

---

### Abuse

```text
API Abuse
```

---

### DoS

```text
Request Flood
```

---

## 9.6 認証前 Rate Limit

### 9.6.1 概要

認証前に適用する Rate Limit。

---

### 9.6.2 Source IP

以下で制限してよい。

```text
Source IP
```

---

### 9.6.3 Token Prefix

以下で制限してよい。

```text
Token Prefix
```

---

### 9.6.4 Endpoint

以下で制限してよい。

```text
Endpoint
```

---

### 9.6.5 Request Path

以下で制限してよい。

```text
Request Path
```

---

## 9.7 認証後 Rate Limit

### 9.7.1 概要

認証成功後に適用する制限。

---

### 9.7.2 Credential

制限単位:

```text
Credential ID
```

---

### 9.7.3 Device

制限単位:

```text
Device ID
```

---

### 9.7.4 User

制限単位:

```text
User ID
```

---

### 9.7.5 Group

制限単位として利用してもよい。

```text
ManagementGroup ID
```

---

## 9.8 Rate Limit 超過時の処理

### 9.8.1 応答コード

以下を返してよい。

```http
429 Too Many Requests
```

---

### 9.8.2 Retry-After

以下を返してよい。

```http
Retry-After
```

---

### 9.8.3 情報秘匿

Rate Limit 応答から Resource 情報を漏洩してはならない。

---

## 9.9 認証失敗遅延

### 9.9.1 目的

高速総当たり攻撃を遅延させる。

---

### 9.9.2 固定遅延

実装は固定遅延を導入してよい。

例:

```text
100ms
```

---

### 9.9.3 ランダム遅延

実装はランダム遅延を導入してよい。

例:

```text
100ms - 500ms
```

---

### 9.9.4 強制失敗

閾値超過時は常に認証失敗として扱ってよい。

---

## 9.10 Token Prefix Protection

### 9.10.1 要件

Token Prefix は Credential 本体ではない。

---

### 9.10.2 利用可能用途

許可:

```text
Audit
```

```text
Rate Limit
```

```text
Monitoring
```

---

### 9.10.3 禁止用途

禁止:

```text
Authentication
```

---

### 9.10.4 Error Leakage

以下を返してはならない。

```text
token_prefix exists
```

```text
token_prefix unknown
```

---

## 9.11 Abuse Detection

### 9.11.1 目的

異常利用を検知する。

---

### 9.11.2 監視例

以下を監視してよい。

```text
Authentication Failure Rate
```

```text
Rate Limit Trigger Rate
```

```text
Credential Guess Attempts
```

---

### 9.11.3 Audit

異常検知イベントは AuditLog に記録すべきである。

---

## 9.12 Monitoring Metrics

以下の Metrics を提供することが推奨される。

```text
authentication_success_total
authentication_failure_total

authorization_failure_total

rate_limit_activations_total

rate_limit_source_ip_total

rate_limit_credential_total

rate_limit_user_total
```

---

## 9.13 Security Requirements

実装は以下を満たさなければならない。

```text
Authentication Failure
≠ Resource Disclosure
```

```text
Credential Not Found
=
Credential Expired
=
Credential Revoked
```

利用者から識別できてはならない。

---

## 9.14 不変条件

実装は以下を保証しなければならない。

```text
Authentication Failure
→ No Resource Disclosure
```

```text
Rate Limit
→ Independent Of Resource Visibility
```

```text
Token Prefix
→ Not Authentication Secret
```

```text
Authentication Failure
→ Audit Logged
```

```text
Rate Limit Trigger
→ Audit Logged
```

これらを満たさない実装は本仕様に適合しない。
