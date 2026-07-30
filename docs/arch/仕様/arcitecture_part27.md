# 27. Error Response

## 27.1 概要

本章では API のエラーレスポンス要件を定義する。

本システムは以下を目的とする。

- 一貫性のあるエラー形式
- Client 実装容易性
- 情報漏洩防止
- 監査容易性

---

## 27.2 基本原則

### 27.2.1 共通形式

API は共通エラー形式を提供すべきである。

---

### 27.2.2 情報漏洩防止

エラーレスポンスから以下を推測できてはならない。

- Credential の存在
- Device の存在
- User の存在
- Token Prefix の有効性
- Resource の所有者

---

### 27.2.3 Request ID

エラーレスポンスは Request ID を含むべきである。

---

## 27.3 共通エラー形式

### 27.3.1 JSON Format

推奨形式:

```json
{
  "error": "invalid_request",
  "message": "request is invalid",
  "request_id": "018f70f2-83d0-7be2-8d63-f95f559a381b"
}
```

---

### 27.3.2 error

機械処理向け識別子。

例:

```text
invalid_request
```

```text
authentication_failed
```

```text
authorization_failed
```

---

### 27.3.3 message

人間向け説明。

---

### 27.3.4 request_id

運用調査用。

---

## 27.4 Validation Error

### 27.4.1 用途

入力値が不正な場合。

---

### 27.4.2 推奨 Status

```http
400 Bad Request
```

---

### 27.4.3 例

```json
{
  "error": "invalid_request",
  "message": "invalid address format"
}
```

---

### 27.4.4 許可事項

修正可能な入力エラー理由は返してよい。

例:

```text
invalid ipv4 format
```

---

## 27.5 Authentication Error

### 27.5.1 用途

認証失敗。

---

### 27.5.2 推奨 Status

```http
401 Unauthorized
```

---

### 27.5.3 応答例

```json
{
  "error": "authentication_failed"
}
```

---

### 27.5.4 情報漏洩禁止

以下を返してはならない。

```text
credential_not_found
```

```text
credential_expired
```

```text
credential_revoked
```

```text
device_not_found
```

---

## 27.6 Authorization Error

### 27.6.1 用途

認可失敗。

---

### 27.6.2 推奨 Status

```http
403 Forbidden
```

---

### 27.6.3 応答例

```json
{
  "error": "authorization_failed"
}
```

---

### 27.6.4 情報漏洩禁止

他 User の Resource 存在を推測できる情報を返してはならない。

---

## 27.7 Resource Not Found

### 27.7.1 用途

Resource が存在しない。

---

### 27.7.2 推奨 Status

```http
404 Not Found
```

---

### 27.7.3 応答例

```json
{
  "error": "resource_not_found"
}
```

---

### 27.7.4 認可との関係

存在確認攻撃を防ぐため、

実装は以下を返してもよい。

```http
404 Not Found
```

認可失敗時であっても。

---

## 27.8 Conflict Error

### 27.8.1 用途

Resource State Conflict。

---

### 27.8.2 推奨 Status

```http
409 Conflict
```

---

### 27.8.3 例

```text
duplicate membership
```

```text
duplicate metadata key
```

---

## 27.9 Rate Limit Error

### 27.9.1 用途

Rate Limit 適用時。

---

### 27.9.2 推奨 Status

```http
429 Too Many Requests
```

---

### 27.9.3 Retry-After

実装は以下を返してよい。

```http
Retry-After
```

---

### 27.9.4 情報漏洩禁止

Rate Limit 応答から Resource 存在を推測できてはならない。

---

## 27.10 Unsupported Media Type

### 27.10.1 用途

不正 Content-Type。

---

### 27.10.2 推奨 Status

```http
415 Unsupported Media Type
```

---

### 27.10.3 例

```text
text/plain
```

```text
application/xml
```

サポートしない場合。

---

## 27.11 Payload Too Large

### 27.11.1 用途

Request Size 超過。

---

### 27.11.2 推奨 Status

```http
413 Payload Too Large
```

---

## 27.12 Internal Error

### 27.12.1 用途

予期しないエラー。

---

### 27.12.2 推奨 Status

```http
500 Internal Server Error
```

---

### 27.12.3 応答例

```json
{
  "error": "internal_error",
  "request_id": "..."
}
```

---

### 27.12.4 詳細情報

Stack Trace を返してはならない。

---

## 27.13 Service Unavailable

### 27.13.1 用途

一時的障害。

---

### 27.13.2 推奨 Status

```http
503 Service Unavailable
```

---

### 27.13.3 例

- Database Unavailable
- Queue Unavailable

---

## 27.14 Error Codes

推奨 Error Code 一覧。

| Error | Meaning |
|---------|---------|
| invalid_request | Validation Failed |
| authentication_failed | Authentication Failed |
| authorization_failed | Authorization Failed |
| resource_not_found | Not Found |
| conflict | State Conflict |
| rate_limited | Rate Limited |
| unsupported_media_type | Invalid Content-Type |
| payload_too_large | Request Too Large |
| internal_error | Internal Failure |
| service_unavailable | Temporary Failure |

---

## 27.15 Audit Requirements

以下を AuditLog に記録すべきである。

- Authentication Failure
- Authorization Failure
- Rate Limit
- Internal Error
- Validation Failure

---

## 27.16 Security Requirements

実装は以下を保証しなければならない。

```text
Authentication Error
→ No Credential Disclosure
```

```text
Authorization Error
→ No Ownership Disclosure
```

```text
Internal Error
→ No Stack Trace Disclosure
```

```text
Rate Limit Error
→ No Resource Disclosure
```

---

## 27.17 不変条件

実装は以下を保証しなければならない。

```text
All Errors
→ JSON Compatible
```

```text
Authentication Failure
→ 401
```

```text
Authorization Failure
→ 403
```

```text
Rate Limited
→ 429
```

```text
Unsupported Media Type
→ 415
```

```text
Payload Too Large
→ 413
```

```text
Unexpected Failure
→ 500
```

```text
No Credential Disclosure
```

```text
No Resource Disclosure
```

```text
No Secret Disclosure
```

これらを満たさない実装は本仕様に適合しない。
