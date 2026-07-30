---
Title: 30. References
Description: 本章は「References」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「30. References」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 30. References

## 30.1 概要

本章では本仕様が参照する標準仕様および参考資料を定義する。

本章の参照文書は、本仕様における規範用語、HTTP セマンティクス、認証方式、および UUID 利用方針の根拠となる。

---

## 30.2 Normative References

### RFC 2119

Key words for use in RFCs to Indicate Requirement Levels

URL:

```text
https://www.rfc-editor.org/info/rfc2119
```

利用箇所:

- MUST
- MUST NOT
- SHOULD
- SHOULD NOT
- MAY

---

### RFC 8174

Ambiguity of Uppercase vs Lowercase in RFC 2119 Key Words

URL:

```text
https://www.rfc-editor.org/info/rfc8174
```

利用箇所:

- 規範用語の解釈

---

### RFC 3986

Uniform Resource Identifier (URI): Generic Syntax

URL:

```text
https://www.rfc-editor.org/info/rfc3986
```

利用箇所:

- API URI 設計

---

### RFC 7235

HTTP/1.1 Authentication

URL:

```text
https://www.rfc-editor.org/info/rfc7235
```

利用箇所:

- Authentication Model
- Basic Authentication
- Digest Authentication

---

### RFC 7617

The "Basic" HTTP Authentication Scheme

URL:

```text
https://www.rfc-editor.org/info/rfc7617
```

利用箇所:

- Basic Authentication

---

### RFC 7616

HTTP Digest Access Authentication

URL:

```text
https://www.rfc-editor.org/info/rfc7616
```

利用箇所:

- Digest Authentication

---

### RFC 7519

JSON Web Token (JWT)

URL:

```text
https://www.rfc-editor.org/info/rfc7519
```

利用箇所:

- UserSession 実装時の参考資料

備考:

本仕様は JWT を要求しない。

---

## 30.3 UUID References

### RFC 9562

Universally Unique IDentifiers (UUIDs)

URL:

```text
https://www.rfc-editor.org/info/rfc9562
```

利用箇所:

- UUIDv7

---

### UUID Version 7

利用箇所:

```text
public_id
```

対象:

- User
- Device
- ManagementGroup
- Credential
- AuditLog
- UpdateHistory
- IntegrationEvent
- CleanupJobRun

---

## 30.4 HTTP References

### HTTP Semantics

RFC 9110

URL:

```text
https://www.rfc-editor.org/info/rfc9110
```

利用箇所:

- HTTP Method
- Status Code
- Request Semantics

---

### HTTP Caching

RFC 9111

URL:

```text
https://www.rfc-editor.org/info/rfc9111
```

参考資料。

---

## 30.5 Security References

### OWASP Top 10

URL:

```text
https://owasp.org/www-project-top-ten/
```

利用箇所:

- Input Validation
- Authentication
- Authorization
- Secret Management

---

### OWASP ASVS

Application Security Verification Standard

URL:

```text
https://owasp.org/www-project-application-security-verification-standard/
```

利用箇所:

- Security Requirements
- Session Management
- Credential Storage

---

### OWASP Cheat Sheet Series

URL:

```text
https://cheatsheetseries.owasp.org/
```

利用箇所:

- CSRF Protection
- Session Security
- Logging

---

## 30.6 Database References

### PostgreSQL Documentation

URL:

```text
https://www.postgresql.org/docs/
```

利用箇所:

- Transaction
- Constraint
- Index
- Foreign Key

---

## 30.7 Monitoring References

### Prometheus

URL:

```text
https://prometheus.io/docs/
```

利用箇所:

- Metrics Export

---

### OpenTelemetry

URL:

```text
https://opentelemetry.io/docs/
```

利用箇所:

- Metrics
- Tracing
- Observability

---

## 30.8 API Description References

### OpenAPI Specification

URL:

```text
https://spec.openapis.org/oas/latest.html
```

利用箇所:

- API Definition
- Security Schemes

---

## 30.9 Informative References

以下は参考資料であり、本仕様の規範文書ではない。

### Twelve-Factor App

URL:

```text
https://12factor.net/
```

---

### Google SRE Book

URL:

```text
https://sre.google/sre-book/table-of-contents/
```

---

### Site Reliability Workbook

URL:

```text
https://sre.google/workbook/
```

---

## 30.10 Reference Usage Policy

本仕様中で RFC、OWASP、OpenAPI 等への参照があった場合は以下の順序で解釈する。

```text
本仕様
↓
Normative References
↓
Informative References
```

参照文書と本仕様が矛盾する場合、本仕様を優先する。

---

## 30.11 Specification Completion Statement

本仕様は以下を定義する。

- Architecture
- Entity Model
- Authentication
- Authorization
- Address Model
- Address Lifecycle
- Metadata Model
- Update Semantics
- AuditLog
- UpdateHistory
- CleanupJob
- External Integration
- Monitoring
- Security Requirements
- Data Model
- SQL Constraints
- OpenAPI Security Schemes
- Error Responses
- Compatibility
- Conformance Checklist

本仕様を満たす実装は、セルフホスト型 IP アドレス管理プラットフォーム Draft v2 に適合する。

---
End of Document
