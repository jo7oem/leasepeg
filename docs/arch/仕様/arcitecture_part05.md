---
Title: 5. アーキテクチャ
Description: 本章は「アーキテクチャ」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「5. アーキテクチャ」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 5. アーキテクチャ

## 5.1 アーキテクチャ概要

本システムは、シンプルなセルフホスト運用からスケールアウト構成まで対応できるアーキテクチャを採用する。

システムは以下のコンポーネントで構成される。

- Reverse Proxy
- API Server
- Worker
- PostgreSQL
- Monitoring Endpoint

---

## 5.2 基本アーキテクチャ

### 5.2.1 論理構成

```text
+----------------+
| Client Device  |
+--------+-------+
         |
         v
+----------------+
| Reverse Proxy  |
+--------+-------+
         |
         v
+----------------+
|   API Server   |
+---+--------+---+
    |        |
    |        |
    v        v
+------+  +------+
| DB   |  |Queue |
+------+  +------+
              |
              v
         +---------+
         | Worker  |
         +---------+
```

---

### 5.2.2 コンポーネント責務

| Component | Responsibility |
|------------|------------|
| Reverse Proxy | TLS Termination |
| API Server | API Processing |
| PostgreSQL | Persistence |
| Worker | Async Processing |
| Monitoring | Health Check / Metrics |

---

## 5.3 最小構成

### 5.3.1 シングルノード構成

小規模環境では以下の構成を推奨する。

```text
Host
 ├─ Reverse Proxy
 ├─ API Server
 ├─ Worker
 └─ PostgreSQL
```

---

### 5.3.2 適用対象

以下の用途を想定する。

- 個人利用
- ホームラボ
- 小規模組織
- VPN管理

---

### 5.3.3 要求事項

最小構成であっても以下を満たさなければならない。

- AuditLog
- UpdateHistory
- Credential Revocation
- Monitoring API

---

## 5.4 スケールアウト構成

### 5.4.1 高負荷構成

大規模環境では以下の構成を採用してよい。

```text
                +-------------+
                | LoadBalancer|
                +------+------+
                       |
         +-------------+-------------+
         |                           |
         v                           v

  +-------------+             +-------------+
  | API Server  |             | API Server  |
  +------+------+             +------+------+
         |                           |
         +-------------+-------------+
                       |
                       v

                +-------------+
                | PostgreSQL  |
                +------+------+

                       |
                       v

                +-------------+
                | Queue       |
                +------+------+
                       |
          +------------+------------+
          |                         |
          v                         v

    +-----------+             +-----------+
    | Worker 1  |             | Worker 2  |
    +-----------+             +-----------+
```

---

### 5.4.2 水平スケーリング

API Server は Stateless とすることが推奨される。

以下を前提とする。

```text
Session State
↓
Database
```

または

```text
Session Store
```

---

### 5.4.3 Worker スケーリング

Worker は複数起動できなければならない。

以下を満たすべきである。

- Event Queue 消費
- Retry 処理
- Dead Letter 処理

---

## 5.5 API Server

### 5.5.1 責務

API Server は以下を担当する。

- Authentication
- Authorization
- Validation
- Database Operation
- AuditLog Recording
- Event Generation

---

### 5.5.2 非同期原則

API Server は以下を同期実行してはならない。

- External API Call
- Retry Processing
- Cleanup Processing

これらは Worker が実行する。

---

### 5.5.3 Stateless

API Server は可能な限り Stateless にすべきである。

以下をメモリに保持してはならない。

- Device State
- Address State
- Credential State

永続状態は Database を唯一の正とする。

---

## 5.6 PostgreSQL

### 5.6.1 採用理由

以下の要件を満たすため、PostgreSQL を採用する。

- Transaction
- Foreign Key
- Constraint
- JSON Support
- Reliability

---

### 5.6.2 永続化対象

以下を保存する。

- User
- Device
- Credential
- Address
- Metadata
- AuditLog
- UpdateHistory

---

### 5.6.3 整合性

Database は整合性維持の最終責任を持つ。

可能な限り以下を利用する。

- Foreign Key
- Unique Constraint
- Transaction

---

## 5.7 Worker

### 5.7.1 責務

Worker は以下を担当する。

- External Integration
- Retry
- CleanupJob
- Scheduled Task

---

### 5.7.2 API 分離

Worker は外部 Client Request を直接受けてはならない。

---

### 5.7.3 冪等性

Worker は同一イベントの再処理に耐えなければならない。

以下を保証すべきである。

```text
Retry Safe
```

---

## 5.8 Queue

### 5.8.1 利用目的

Queue は以下のために使用する。

- Event Delivery
- Retry
- Backpressure

---

### 5.8.2 任意実装

実装は以下のいずれでもよい。

- PostgreSQL Queue
- Redis
- RabbitMQ
- Kafka

本仕様では特定実装を要求しない。

---

## 5.9 Monitoring API

### 5.9.1 目的

Monitoring API は以下を提供する。

- Health Check
- Readiness Check
- Metrics

---

### 5.9.2 エンドポイント例

```http
GET /healthcheck/live
GET /healthcheck/ready
GET /metrics
```

---

### 5.9.3 公開制御

Monitoring API は既定で外部公開してはならない。

---

## 5.10 Internal API

### 5.10.1 用途

以下を提供してよい。

- CleanupJob Trigger
- Worker Control
- Internal Diagnostics

---

### 5.10.2 セキュリティ

Internal API は外部公開してはならない。

外部公開が必要な場合は強い認証を必須とする。

---

## 5.11 ネットワーク分離

本システムは以下の通信経路を分離することが推奨される。

### Public Path

```text
Internet
 ↓
Application API
```

---

### Internal Path

```text
Worker
Database
Monitoring
```

---

### Management Path

```text
Admin Access
```

---

## 5.12 障害耐性

### 5.12.1 API Server 障害

API Server 障害時でも Database の整合性は維持されなければならない。

---

### 5.12.2 Worker 障害

Worker 障害時でも Address Update API は継続利用できるべきである。

External Integration は遅延してもよい。

---

### 5.12.3 Database 障害

Database が利用不能な場合、書き込み処理は失敗してよい。

ただし不整合を発生させてはならない。

---

## 5.13 アーキテクチャ原則

実装は以下を満たさなければならない。

- API Server は Stateless
- Database は Source of Truth
- Worker は Async Processing
- Monitoring は独立 Port
- Internal API は非公開
- External Integration は非同期
- Retry Safe
- Last Accepted Operation Wins

以上をアーキテクチャ要件とする。
