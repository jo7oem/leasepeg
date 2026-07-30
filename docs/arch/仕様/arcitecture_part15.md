---
Title: 15. API
Description: 本章は「API」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「15. API」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 15. API

## 15.1 概要

本章では本システムが提供する API の共通要件および各 API カテゴリの要件を定義する。

API は以下のカテゴリで構成される。

- Device Update API
- ManagementGroup Read API
- User Management API
- Monitoring API
- Internal API

---

## 15.2 API 共通要件

### 15.2.1 Transport

Production 環境では HTTPS を使用しなければならない。

```text
HTTP
→ MUST NOT
```

```text
HTTPS
→ MUST
```

---

### 15.2.2 Content-Type

実装は Content-Type を検証しなければならない。

例:

```http
Content-Type: application/json
```

---

### 15.2.3 Request Size

実装は Request Body Size を制限しなければならない。

制限値は実装依存とする。

---

### 15.2.4 Authentication

Non-Public API は認証を要求しなければならない。

---

### 15.2.5 Authorization

認証成功後は認可を実施しなければならない。

```text
Authentication
≠
Authorization
```

---

### 15.2.6 Request ID

全 API は Request ID を生成または受理すべきである。

例:

```http
X-Request-ID
```

---

### 15.2.7 Retry Safety

同一 Request の再送により不正状態になってはならない。

---

## 15.3 API Versioning

### 15.3.1 Version Prefix

API は Version Prefix を持つべきである。

例:

```text
/v1
```

---

### 15.3.2 Breaking Change

Breaking Change は新 Version として提供すべきである。

例:

```text
/v1
```

↓

```text
/v2
```

---

## 15.4 Device Update API

### 15.4.1 目的

Device が Address および Metadata を更新するための API。

---

### 15.4.2 Authentication

以下を使用する。

```text
DeviceCredential
```

---

### 15.4.3 Bearer Example

```http
POST /v1/device/address

Authorization: Bearer <device_token>
```

---

### 15.4.4 Basic Example

```http
POST /v1/device/address

Authorization: Basic xxxxx
```

---

### 15.4.5 Device Identification

Device は Credential から解決しなければならない。

---

### 15.4.6 禁止事項

Request 内で Device ID を指定させてはならない。

禁止例:

```json
{
  "device_id": "..."
}
```

---

### 15.4.7 許可操作

以下のみ。

```text
Address Update
```

```text
Metadata Update
```

---

### 15.4.8 禁止操作

以下を許可してはならない。

```text
Address Read
```

```text
User Management
```

```text
Credential Management
```

---

## 15.5 ManagementGroup Address Read API

### 15.5.1 目的

ManagementGroup 内 Device の Address を参照する。

---

### 15.5.2 Authentication

以下を使用する。

```text
GroupReadToken
```

---

### 15.5.3 Example

```http
GET /v1/management-groups/{group_id}/addresses

Authorization: Bearer <group_read_token>
```

---

### 15.5.4 参照範囲

Token に紐付いた Group のみ参照できなければならない。

---

### 15.5.5 Return Data

返却対象:

```text
active Address only
```

---

### 15.5.6 禁止操作

```text
Update
```

```text
Delete
```

```text
Management
```

---

## 15.6 User Management API

### 15.6.1 目的

管理者向け API。

---

### 15.6.2 Authentication

以下を使用する。

```text
UserSession
```

---

### 15.6.3 Device Management

許可:

```text
Create
Update
Delete
```

---

### 15.6.4 Group Management

許可:

```text
Create
Update
Delete
```

---

### 15.6.5 Credential Management

許可:

```text
Issue
Revoke
```

---

### 15.6.6 Audit Access

許可:

```text
AuditLog Read
```

---

### 15.6.7 Address Access

許可:

```text
Address Read
```

---

### 15.6.8 Cleanup Access

提供する場合に限り許可:

```text
Cleanup API
```

---

## 15.7 Monitoring API

### 15.7.1 目的

監視およびヘルスチェック。

---

### 15.7.2 推奨 Endpoint

```http
GET /healthcheck/live
```

```http
GET /healthcheck/ready
```

```http
GET /metrics
```

---

### 15.7.3 Port Separation

Application API と別 Port で Listen しなければならない。

---

### 15.7.4 公開制御

既定で外部公開してはならない。

---

## 15.8 Internal API

### 15.8.1 目的

運用管理用。

---

### 15.8.2 用途

例:

```text
Cleanup Trigger
```

```text
Worker Control
```

```text
Diagnostics
```

---

### 15.8.3 公開制御

外部公開してはならない。

---

### 15.8.4 認証

公開する場合は Application API とは別の強い認証を要求しなければならない。

---

## 15.9 HTTP Method Policy

### Read

以下を使用する。

```http
GET
```

---

### Create

以下を使用する。

```http
POST
```

---

### Update

以下を使用する。

```http
PUT
```

または

```http
PATCH
```

---

### Delete

以下を使用する。

```http
DELETE
```

---

## 15.10 Response Format

### 15.10.1 JSON

API Response は JSON を使用すべきである。

---

### 15.10.2 Success Example

```json
{
  "result": "success"
}
```

---

### 15.10.3 Error Example

```json
{
  "error": "invalid_request",
  "message": "request is invalid",
  "request_id": "..."
}
```

---

## 15.11 Audit Requirements

以下を AuditLog に記録しなければならない。

- Authentication Success
- Authentication Failure
- Authorization Failure
- Update
- Delete
- Credential Operation

---

## 15.12 Authorization Requirements

実装は以下を保証しなければならない。

```text
UserSession
→ Owned Resources Only
```

```text
DeviceCredential
→ Bound Device Only
```

```text
GroupReadToken
→ Bound Group Only
```

---

## 15.13 Error Handling

### Validation Failure

推奨:

```http
400 Bad Request
```

---

### Authentication Failure

推奨:

```http
401 Unauthorized
```

---

### Authorization Failure

推奨:

```http
403 Forbidden
```

---

### Resource Not Found

推奨:

```http
404 Not Found
```

---

### Rate Limit

推奨:

```http
429 Too Many Requests
```

---

### Server Error

推奨:

```http
500 Internal Server Error
```

---

## 15.14 不変条件

実装は以下を保証しなければならない。

```text
Non-Public API
→ Authentication Required
```

```text
Authentication
→ Authorization Required
```

```text
Device Update API
→ DeviceCredential Only
```

```text
Group Read API
→ GroupReadToken Only
```

```text
User Management API
→ UserSession Only
```

```text
Monitoring API
→ Separate Port
```

```text
Internal API
→ Not Public By Default
```

```text
Request Retry
→ Safe
```

これらを満たさない実装は本仕様に適合しない。
