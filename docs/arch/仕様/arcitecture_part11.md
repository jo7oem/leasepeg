# 11. Address Lifecycle

## 11.1 概要

本章では DeviceAddress のライフサイクルを定義する。

Address は論理的に以下の状態を持つ。

```text
active
expired
deleted
```

Address は必ずいずれかの状態に属さなければならない。

---

## 11.2 ライフサイクル状態

### 11.2.1 State Diagram

```text
            +--------+
            | active |
            +----+---+
                 |
                 | expires_at reached
                 v
            +--------+
            |expired |
            +----+---+
                 |
                 | cleanup
                 v
            +--------+
            |deleted |
            +--------+
```

または以下によって active から deleted へ直接遷移してもよい。

```text
Device Delete
Replace Delete
Administrative Delete
```

---

## 11.3 Active

### 11.3.1 定義

以下を満たす状態。

```text
Address Record Exists
AND
expires_at > now
```

---

### 11.3.2 特徴

Active Address は利用可能である。

Address Read API は Active Address を返却しなければならない。

---

### 11.3.3 Read Visibility

以下から参照可能。

- UserSession
- GroupReadToken

権限要件を満たす場合に限る。

---

### 11.3.4 更新

Active Address は以下によって更新できる。

- update
- replace
- nop

---

## 11.4 Expired

### 11.4.1 定義

以下を満たす状態。

```text
Address Record Exists
AND
expires_at <= now
```

---

### 11.4.2 特徴

Address Record は存在する。

ただし有効な Address として扱ってはならない。

---

### 11.4.3 Read Visibility

Expired Address は通常の Address Read API から返してはならない。

---

### 11.4.4 再利用

同じ Device に対する新しい Update により、再び Active Address が作成されてもよい。

---

### 11.4.5 保持

Expired Address は一定期間 Database に残してもよい。

即時削除は必須ではない。

---

## 11.5 Deleted

### 11.5.1 定義

Address Record が削除された状態。

---

### 11.5.2 特徴

Deleted Address は通常 API から参照できない。

---

### 11.5.3 履歴

Deleted 状態であっても以下から確認できるべきである。

- AuditLog
- UpdateHistory

---

### 11.5.4 再作成

Deleted 状態の後、新たな Address が作成されてもよい。

これは新しい Active Address として扱う。

---

## 11.6 Expiration

### 11.6.1 目的

古い Address を自動的に無効化する。

---

### 11.6.2 判定条件

Active:

```text
expires_at > now
```

Expired:

```text
expires_at <= now
```

---

### 11.6.3 Time Source

判定に使用する時刻は System Clock を使用する。

実装は UTC 使用を推奨する。

---

## 11.7 Address 作成

### 11.7.1 初回登録

Address が登録された場合、

```text
active
```

状態として生成しなければならない。

---

### 11.7.2 必須項目

以下を持たなければならない。

- address
- family
- expires_at
- created_at
- updated_at

---

## 11.8 Address 更新

### 11.8.1 Update

対象 Family の Address を更新する。

状態は Active のままである。

---

### 11.8.2 Replace

対象 Family を残し、他 Family を削除する。

削除された Address は Deleted とみなす。

---

### 11.8.3 Nop

Address 値は変更しない。

以下のみ更新してよい。

```text
expires_at
last_seen_at
updated_at
```

---

## 11.9 Replace による削除

### 11.9.1 IPv4 Replace

例:

```text
Before:
IPv4 active
IPv6 active

Replace IPv4

After:
IPv4 active
IPv6 deleted
```

---

### 11.9.2 IPv6 Replace

例:

```text
Before:
IPv4 active
IPv6 active

Replace IPv6

After:
IPv4 deleted
IPv6 active
```

---

## 11.10 Device Delete

### 11.10.1 動作

Device 削除時は関連 Address を削除しなければならない。

---

### 11.10.2 結果

Address は Deleted 状態となる。

---

### 11.10.3 Audit

削除イベントは AuditLog に記録しなければならない。

---

## 11.11 CleanupJob

### 11.11.1 目的

Expired Address を物理削除する。

---

### 11.11.2 対象

以下を対象としてよい。

```text
expired
```

---

### 11.11.3 実行方式

以下のいずれでもよい。

- Scheduled
- Manual Trigger
- Internal API

---

## 11.12 Read API 要件

### 11.12.1 Active Only

Address Read API は Active Address のみ返さなければならない。

---

### 11.12.2 Expired Exclusion

Expired Address を返してはならない。

---

### 11.12.3 Deleted Exclusion

Deleted Address を返してはならない。

---

## 11.13 Audit Requirements

以下の状態遷移は AuditLog に記録しなければならない。

### 作成

```text
none
→ active
```

---

### 期限切れ

```text
active
→ expired
```

---

### 削除

```text
active
→ deleted
```

または

```text
expired
→ deleted
```

---

## 11.14 UpdateHistory Requirements

UpdateHistory は少なくとも以下を保持すべきである。

- old state
- new state
- timestamp
- actor
- reason

---

## 11.15 Monitoring Metrics

以下の Metrics を提供することが推奨される。

```text
active_addresses
```

```text
expired_addresses
```

```text
address_lifecycle_transitions_total
```

```text
address_deleted_total
```

---

## 11.16 不変条件

実装は以下を保証しなければならない。

```text
active
→ Readable
```

```text
expired
→ Not Returned By Read API
```

```text
deleted
→ Not Returned By Read API
```

```text
expires_at <= now
→ expired
```

```text
expired
→ may be removed by CleanupJob
```

```text
Device Delete
→ Address Deleted
```

```text
Address Lifecycle
∈ {active, expired, deleted}
```

これらの条件を満たさない実装は本仕様に適合しない。
