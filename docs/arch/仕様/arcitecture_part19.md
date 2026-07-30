# 19. 冪等性および結果整合性

## 19.1 概要

本章では本システムにおける冪等性（Idempotency）および結果整合性（Eventual Consistency）の要件を定義する。

本システムは以下を保証しなければならない。

- Retry Safe
- Last Accepted Operation Wins
- Eventual Consistency
- Async Processing Compatibility
- External Integration Consistency

---

## 19.2 基本方針

### 19.2.1 冪等性

同一 Request が複数回実行されても最終結果が不正になってはならない。

---

### 19.2.2 結果整合性

非同期処理を含む場合でも、最終的に正しい状態へ収束しなければならない。

---

### 19.2.3 再送前提

実装は Client Retry を前提として設計しなければならない。

---

## 19.3 Address Update の冪等性

### 19.3.1 update

同一 Address への update は複数回実行されてもよい。

例:

```text
Initial:
192.0.2.1

Update:
192.0.2.10

Retry:
192.0.2.10
```

結果:

```text
192.0.2.10
```

---

### 19.3.2 replace

replace は複数回実行されてもよい。

例:

```text
IPv4 active
IPv6 active

replace ipv4

IPv4 active
IPv6 deleted

replace ipv4 again

IPv4 active
IPv6 deleted
```

---

### 19.3.3 nop

nop は何度実行されてもよい。

---

## 19.4 Metadata Update の冪等性

### 19.4.1 同一値更新

同じ値への更新は成功として扱ってよい。

例:

```json
{
  "name": "router-a"
}
```

を複数回送信してもよい。

---

### 19.4.2 空文字列

空文字列への更新も同様とする。

---

### 19.4.3 Delete

Metadata Delete は複数回実行されてもよい。

---

## 19.5 Last Accepted Operation Wins

### 19.5.1 原則

同一 Resource に複数 Update が発生した場合、

最後に受理された操作が最終状態を決定する。

---

### 19.5.2 例

時系列:

```text
T1
Address=192.0.2.1
```

```text
T2
Address=192.0.2.2
```

```text
T3
Address=192.0.2.3
```

結果:

```text
192.0.2.3
```

---

### 19.5.3 Worker Retry

Worker Retry によって古い状態へ戻してはならない。

---

## 19.6 Request Retry

### 19.6.1 Timeout

Client は Timeout 後に再送してよい。

---

### 19.6.2 必須要件

実装は Retry に耐えなければならない。

---

### 19.6.3 失敗例

以下は許可されない。

```text
Retry
↓
Duplicate Address
```

---

### 19.6.4 成功例

```text
Retry
↓
Same State
```

---

## 19.7 非同期処理

### 19.7.1 原則

非同期処理は結果整合性モデルを採用する。

---

### 19.7.2 API Server

API Server は Database 更新を完了させる。

---

### 19.7.3 Worker

Worker は後続処理を担当する。

---

### 19.7.4 外部連携

External Integration は非同期で処理する。

---

## 19.8 Eventual Consistency

### 19.8.1 定義

一時的な不一致が存在してもよい。

---

### 19.8.2 最終状態

最終的には API Server が保持する状態と一致しなければならない。

---

### 19.8.3 例

```text
API Update Success
```

↓

```text
Worker Queue Pending
```

↓

```text
External API Updated
```

---

## 19.9 External Integration Consistency

### 19.9.1 原則

External Integration は最終的に API 状態へ収束しなければならない。

---

### 19.9.2 同期実行禁止

以下は禁止する。

```text
Request
↓
External API Call
↓
Response
```

---

### 19.9.3 推奨構成

```text
Request
↓
DB Update
↓
Event Queue
↓
Response
↓
Worker
↓
External API
```

---

## 19.10 Event Processing

### 19.10.1 Retry

Event は再実行可能でなければならない。

---

### 19.10.2 Duplicate Event

重複 Event が発生しても不正状態になってはならない。

---

### 19.10.3 Ordering

可能な限り順序を維持することを推奨する。

---

## 19.11 Failure Recovery

### 19.11.1 Worker Failure

Worker Failure 時も API Server の整合性は維持しなければならない。

---

### 19.11.2 Queue Failure

Queue 障害後に処理を再開できるべきである。

---

### 19.11.3 External API Failure

External API Failure により内部状態を巻き戻してはならない。

---

## 19.12 Audit Requirements

以下を AuditLog に記録すべきである。

- Retry
- Duplicate Detection
- Worker Failure
- Queue Failure
- External Synchronization Failure

---

## 19.13 Monitoring

実装は以下を提供してよい。

```text
request_retries_total
```

```text
worker_retries_total
```

```text
duplicate_events_total
```

```text
external_sync_failures_total
```

```text
eventual_consistency_repairs_total
```

---

## 19.14 不変条件

実装は以下を保証しなければならない。

```text
Same Request
→ Same Final State
```

```text
Retry
→ Safe
```

```text
Duplicate Event
→ Safe
```

```text
Last Accepted Operation
→ Wins
```

```text
External Integration
→ Eventually Consistent
```

```text
Worker Retry
→ MUST NOT Corrupt State
```

```text
External Failure
→ MUST NOT Rollback Internal State
```

```text
Async Processing
→ Eventually Consistent
```

これらを満たさない実装は本仕様に適合しない。
