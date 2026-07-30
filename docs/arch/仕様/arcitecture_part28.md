---
Title: 28. Compatibility
Description: 本章は「Compatibility」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「28. Compatibility」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 28. Compatibility

## 28.1 概要

本章では本システムが考慮する互換性要件を定義する。

本システムは以下との互換性を重視する。

- Router Devices
- Legacy Clients
- IPv4 Networks
- IPv6 Networks
- Retry-capable Clients
- Self-hosted Environments

---

## 28.2 基本方針

### 28.2.1 後方互換性

可能な限り既存クライアントとの互換性を維持すべきである。

---

### 28.2.2 セキュリティ優先

互換性のためにセキュリティ要件を緩和してはならない。

---

### 28.2.3 標準準拠

可能な限り標準的な HTTP および認証方式を利用すべきである。

---

## 28.3 Router Compatibility

### 28.3.1 概要

一部 Router は柔軟な HTTP Client を持たない。

そのため追加互換 API を提供してよい。

---

### 28.3.2 GET Update

実装は GET による Address Update を提供してよい。

例:

```http
GET /v1/device/update?ip=203.0.113.1
```

---

### 28.3.3 Security

GET Update を提供する場合でも認証を要求しなければならない。

---

### 28.3.4 Logging

Credential 情報が URL に含まれないよう注意しなければならない。

---

## 28.4 Authentication Compatibility

### 28.4.1 Bearer Token

推奨方式。

---

### 28.4.2 Basic Authentication

以下のために提供してよい。

- Legacy Router
- Embedded Device

---

### 28.4.3 Digest Authentication

以下のために提供してよい。

- Router Compatibility
- Legacy Device Compatibility

---

### 28.4.4 Authorization

認証方式に関係なく権限は同一でなければならない。

---

## 28.5 IPv4 Compatibility

### 28.5.1 Support

IPv4 環境をサポートしなければならない。

---

### 28.5.2 Endpoint

以下を提供してよい。

```text
v4only
```

---

### 28.5.3 Address Update

IPv4 Address を登録できなければならない。

---

## 28.6 IPv6 Compatibility

### 28.6.1 Support

IPv6 環境をサポートしなければならない。

---

### 28.6.2 Endpoint

以下を提供してよい。

```text
v6only
```

---

### 28.6.3 Address Update

IPv6 Address を登録できなければならない。

---

## 28.7 Dual Stack Compatibility

### 28.7.1 Support

IPv4 / IPv6 の両方を同時にサポートしなければならない。

---

### 28.7.2 Endpoint

以下を提供してよい。

```text
dualstack
```

---

### 28.7.3 Independent Management

IPv4 と IPv6 は独立して更新できなければならない。

---

## 28.8 Client Retry Compatibility

### 28.8.1 Timeout

クライアントは Timeout 後に再送してよい。

---

### 28.8.2 Retry Safety

実装は Retry Safe でなければならない。

---

### 28.8.3 Idempotency

同一 Request の再送で不正状態になってはならない。

---

## 28.9 Deployment Compatibility

### 28.9.1 Single Node

単一ノード構成をサポートすべきである。

---

### 28.9.2 Multi Node

複数 API Server 構成をサポートしてよい。

---

### 28.9.3 Container Environment

以下をサポートしてよい。

- Docker
- Podman
- Kubernetes

---

## 28.10 Database Compatibility

### 28.10.1 Primary Target

対象:

```text
PostgreSQL
```

---

### 28.10.2 Extension Dependency

特殊 Extension への依存は最小限にすべきである。

---

## 28.11 API Compatibility

### 28.11.1 Version Prefix

API Version を明示すべきである。

例:

```text
/v1
```

---

### 28.11.2 Breaking Change

Breaking Change は新 Version として提供すべきである。

---

### 28.11.3 Error Format

エラー形式は Version 内で安定しているべきである。

---

## 28.12 Monitoring Compatibility

### 28.12.1 Prometheus

Prometheus Export を推奨する。

---

### 28.12.2 OpenTelemetry

OpenTelemetry をサポートしてよい。

---

## 28.13 External Integration Compatibility

### 28.13.1 Retry

外部 API Retry に耐えるべきである。

---

### 28.13.2 Duplicate Event

重複処理に耐えるべきである。

---

### 28.13.3 Eventual Consistency

結果整合性を前提とする。

---

## 28.14 Upgrade Compatibility

### 28.14.1 Migration

Schema Migration をサポートしなければならない。

---

### 28.14.2 Data Preservation

Upgrade 時にユーザーデータを消失させてはならない。

---

### 28.14.3 Rollback

可能であれば Rollback Strategy を提供すべきである。

---

## 28.15 不変条件

実装は以下を保証しなければならない。

```text
IPv4 Supported
```

```text
IPv6 Supported
```

```text
Dual Stack Supported
```

```text
Retry Safe
```

```text
GET Update
→ Optional
```

```text
Basic Authentication
→ Optional
```

```text
Digest Authentication
→ Optional
```

```text
Bearer Authentication
→ Supported
```

```text
API Versioning
→ Supported
```

```text
PostgreSQL
→ Supported
```

```text
Single Node Deployment
→ Supported
```

これらを満たさない実装は本仕様に適合しない。
