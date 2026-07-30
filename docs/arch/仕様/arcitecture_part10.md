# 10. Address Model

## 10.1 概要

本章では Device が保持する Address のデータモデルを定義する。

Address Model は以下を提供する。

- IPv4 Address 管理
- IPv6 Address 管理
- Address Lifecycle
- Address Expiration
- Address History の基礎情報

Address は Device に対して管理される。

---

## 10.2 Address Family

### 10.2.1 対応する Family

Address Family は以下のいずれかでなければならない。

```text
ipv4
ipv6
```

---

### 10.2.2 不正値

以下は許可してはならない。

```text
v4
```

```text
v6
```

```text
IPv4
```

```text
IPv6
```

内部的な保存形式は実装依存であるが、API 上の表現は統一しなければならない。

---

## 10.3 DeviceAddress

### 10.3.1 責務

DeviceAddress は Device の現在の IP Address を保持する。

---

### 10.3.2 必須属性

DeviceAddress は少なくとも以下を持たなければならない。

| Field | Description |
|---------|---------|
| device_id | Device |
| family | Address Family |
| address | IP Address |
| expires_at | Expiration |
| last_seen_at | Last Update |
| created_at | Creation Time |
| updated_at | Update Time |

---

## 10.4 Device と Address

### 10.4.1 関係

```text
Device
 ├─ IPv4 Address
 └─ IPv6 Address
```

---

### 10.4.2 保持数

同一 Device は最大 2 件の Address を保持できる。

```text
IPv4: 0..1
IPv6: 0..1
```

---

## 10.5 一意性制約

### 10.5.1 Family 一意性

以下を保証しなければならない。

```text
UNIQUE(device_id, family)
```

---

### 10.5.2 例

許可:

```text
Device A
 ├─ IPv4
 └─ IPv6
```

禁止:

```text
Device A
 ├─ IPv4
 └─ IPv4
```

---

## 10.6 Address Value

### 10.6.1 IPv4

有効な IPv4 Address を保持できなければならない。

例:

```text
192.0.2.1
```

---

### 10.6.2 IPv6

有効な IPv6 Address を保持できなければならない。

例:

```text
2001:db8::1
```

---

### 10.6.3 Validation

Address 登録時はアドレス形式を検証しなければならない。

---

### 10.6.4 Reachability

実装は Address の到達性を検証してはならない。

---

### 10.6.5 Ownership

実装は Address の所有性を検証してはならない。

---

## 10.7 submitted_ip

### 10.7.1 定義

Device が Request 内で明示的に通知した Address。

---

### 10.7.2 用途

以下に利用する。

```text
Address Update
```

---

### 10.7.3 Validation

Address Format は検証しなければならない。

---

### 10.7.4 Ownership Verification

以下を実施してはならない。

```text
Ownership Verification
```

---

### 10.7.5 Reachability Verification

以下を実施してはならない。

```text
Ping
```

```text
Traceroute
```

```text
Port Check
```

---

## 10.8 observed_ip

### 10.8.1 定義

HTTP Connection の Source Address。

---

### 10.8.2 利用用途

以下に利用してよい。

- AuditLog
- Rate Limit
- Abuse Analysis

---

### 10.8.3 Address Registration

observed_ip を Device Address として登録してはならない。

---

## 10.9 Address Read API

### 10.9.1 基本原則

Address Read API は active Address のみ返さなければならない。

---

### 10.9.2 非表示

以下は返してはならない。

```text
expired
```

```text
deleted
```

---

### 10.9.3 対象

対象となるのは以下である。

```text
active
```

---

## 10.10 Expiration

### 10.10.1 目的

古い Address を自動的に無効化する。

---

### 10.10.2 expires_at

Address は有効期限を持たなければならない。

```text
expires_at
```

---

### 10.10.3 判定

以下を満たす場合。

```text
expires_at > current_time
```

状態:

```text
active
```

---

以下を満たす場合。

```text
expires_at <= current_time
```

状態:

```text
expired
```

---

## 10.11 last_seen_at

### 10.11.1 目的

Device が最後に更新を行った時刻を保持する。

---

### 10.11.2 更新タイミング

以下で更新してよい。

- Update
- Replace
- Nop

---

### 10.11.3 Retry

同一 Request の再送で更新されてもよい。

---

## 10.12 Address Lifecycle

Address は以下の状態を持つ。

```text
active
expired
deleted
```

---

### active

利用可能状態。

---

### expired

期限切れ状態。

---

### deleted

削除済み状態。

---

## 10.13 Address Deletion

### 10.13.1 削除用途

Address は以下により削除される。

- CleanupJob
- Device Delete
- Replace Operation

---

### 10.13.2 履歴

削除後も UpdateHistory または AuditLog から確認できるべきである。

---

## 10.14 Monitoring

以下の Metrics を提供することが推奨される。

```text
active_addresses
```

```text
expired_addresses
```

```text
address_updates_total
```

```text
address_update_failures_total
```

---

## 10.15 不変条件

実装は以下を保証しなければならない。

```text
Address Family
∈ {ipv4, ipv6}
```

```text
UNIQUE(device_id, family)
```

```text
Read API
→ active only
```

```text
submitted_ip
→ format validated
```

```text
submitted_ip
→ ownership not verified
```

```text
submitted_ip
→ reachability not verified
```

```text
observed_ip
≠ Device Address
```

```text
Address Count per Device
<= 2
```

これらを満たさない実装は本仕様に適合しない。
