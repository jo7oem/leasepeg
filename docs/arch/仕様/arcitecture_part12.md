# 12. Address Update Semantics

## 12.1 概要

本章では Address Update API の更新規則を定義する。

本システムは以下の更新モードを提供する。

- update
- replace
- nop

実装は本章で定義するセマンティクスに従わなければならない。

---

## 12.2 用語

### submitted_ip

Device が明示的に通知した Address。

---

### observed_ip

HTTP Connection の Source Address。

---

### target_family

更新対象となる Address Family。

---

### connection address family

Connection の Source Address から判定された Family。

---

## 12.3 target_family の決定

### 12.3.1 submitted_ip が存在する場合

以下を使用する。

```text
target_family
=
family(submitted_ip)
```

---

### 12.3.2 submitted_ip が存在しない場合

以下を使用する。

```text
target_family
=
connection address family
```

---

### 12.3.3 例

#### IPv4 Connection + submitted_ip IPv6

```text
Connection:
198.51.100.1

submitted_ip:
2001:db8::1

target_family:
ipv6
```

許可される。

---

#### IPv6 Connection + submitted_ip IPv4

```text
Connection:
2001:db8::10

submitted_ip:
203.0.113.5

target_family:
ipv4
```

許可される。

---

## 12.4 submitted_ip

### 12.4.1 Address Validation

実装はフォーマット検証を行わなければならない。

---

### 12.4.2 Ownership Verification

以下を行ってはならない。

```text
Address Ownership Verification
```

---

### 12.4.3 Reachability Verification

以下を行ってはならない。

```text
Ping
```

```text
Traceroute
```

```text
TCP Connect
```

---

### 12.4.4 理由

Address の正当性は Device または User の責任とする。

---

## 12.5 更新モード

更新モードは以下のいずれかとする。

```text
update
replace
nop
```

---

## 12.6 update

### 12.6.1 目的

対象 Family のみ更新する。

---

### 12.6.2 IPv4 Update

初期状態:

```text
IPv4 = A
IPv6 = B
```

更新:

```text
target_family = ipv4
new_address = C
```

結果:

```text
IPv4 = C
IPv6 = B
```

---

### 12.6.3 IPv6 Update

初期状態:

```text
IPv4 = A
IPv6 = B
```

更新:

```text
target_family = ipv6
new_address = D
```

結果:

```text
IPv4 = A
IPv6 = D
```

---

### 12.6.4 他 Family

対象外 Family を変更してはならない。

---

## 12.7 replace

### 12.7.1 目的

対象 Family のみ残す。

---

### 12.7.2 IPv4 Replace

更新前:

```text
IPv4 = A
IPv6 = B
```

更新:

```text
target_family = ipv4
```

更新後:

```text
IPv4 = A
IPv6 = deleted
```

---

### 12.7.3 IPv6 Replace

更新前:

```text
IPv4 = A
IPv6 = B
```

更新:

```text
target_family = ipv6
```

更新後:

```text
IPv4 = deleted
IPv6 = B
```

---

### 12.7.4 削除履歴

Replace による削除は UpdateHistory および AuditLog に記録すべきである。

---

## 12.8 nop

### 12.8.1 目的

Address を変更せず TTL を延長する。

---

### 12.8.2 変更可能項目

以下のみ更新してよい。

```text
expires_at
```

```text
last_seen_at
```

```text
updated_at
```

---

### 12.8.3 Address Value

Address Value を更新してはならない。

---

### 12.8.4 例

更新前:

```text
address = 192.0.2.1
expires_at = T1
```

nop 実行後:

```text
address = 192.0.2.1
expires_at = T2
```

---

## 12.9 nop + submitted_ip

### 12.9.1 禁止

nop に submitted_ip を含めてはならない。

---

### 12.9.2 応答

以下を返さなければならない。

```http
400 Bad Request
```

---

### 12.9.3 理由

nop は Address 更新を伴わないため。

---

## 12.10 同一 Address 更新

### 12.10.1 許可

同一 Address への Update は成功として扱ってよい。

---

### 12.10.2 更新項目

以下は更新してよい。

```text
expires_at
```

```text
last_seen_at
```

```text
updated_at
```

---

### 12.10.3 例

更新前:

```text
192.0.2.1
```

更新後:

```text
192.0.2.1
```

成功として扱ってよい。

---

## 12.11 Retry Safety

### 12.11.1 要件

同一 Request の再送によって不正状態になってはならない。

---

### 12.11.2 update

複数回実行されてもよい。

---

### 12.11.3 replace

複数回実行されてもよい。

---

### 12.11.4 nop

複数回実行されてもよい。

---

## 12.12 Audit Requirements

以下は AuditLog 記録対象とする。

- update
- replace
- nop
- validation failure

---

### 記録推奨項目

```text
actor
device_id
family
mode
old_address
new_address
result
timestamp
```

---

## 12.13 UpdateHistory Requirements

UpdateHistory は少なくとも以下を保持すべきである。

```text
device_id
```

```text
family
```

```text
old_address
```

```text
new_address
```

```text
mode
```

```text
timestamp
```

---

## 12.14 不変条件

実装は以下を保証しなければならない。

```text
submitted_ip
→ target_family
```

```text
submitted_ip present
→ family(submitted_ip)
```

```text
submitted_ip absent
→ connection family
```

```text
update
→ target_family only
```

```text
replace
→ target_family remains
```

```text
nop
→ no address change
```

```text
nop + submitted_ip
→ 400 Bad Request
```

```text
Address Ownership Verification
→ MUST NOT
```

```text
Address Reachability Verification
→ MUST NOT
```

これらを満たさない実装は本仕様に適合しない。
