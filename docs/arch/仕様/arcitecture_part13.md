# 13. Endpoint Address Family

## 13.1 概要

本章では Address Update API における Endpoint Address Family の取り扱いを定義する。

本システムは以下の Endpoint 種別を提供する。

- v4only
- v6only
- dualstack

Endpoint Type は Address Family 判定および監査目的で利用される。

---

## 13.2 Endpoint Type

### 13.2.1 v4only

IPv4 通信専用 Endpoint。

例:

```text
v4only.example.com
```

---

### 13.2.2 v6only

IPv6 通信専用 Endpoint。

例:

```text
v6only.example.com
```

---

### 13.2.3 dualstack

IPv4 および IPv6 の両方を受け付ける Endpoint。

例:

```text
dualstack.example.com
```

---

## 13.3 v4only Endpoint

### 13.3.1 目的

IPv4 専用クライアントをサポートする。

---

### 13.3.2 接続条件

以下のみ受け付ける。

```text
IPv4 Connection
```

---

### 13.3.3 IPv6 接続

以下は受け付けてはならない。

```text
IPv6 Connection
```

---

### 13.3.4 target_family

submitted_ip が存在しない場合:

```text
target_family = ipv4
```

---

## 13.4 v6only Endpoint

### 13.4.1 目的

IPv6 専用クライアントをサポートする。

---

### 13.4.2 接続条件

以下のみ受け付ける。

```text
IPv6 Connection
```

---

### 13.4.3 IPv4 接続

以下は受け付けてはならない。

```text
IPv4 Connection
```

---

### 13.4.4 target_family

submitted_ip が存在しない場合:

```text
target_family = ipv6
```

---

## 13.5 dualstack Endpoint

### 13.5.1 目的

IPv4 および IPv6 の双方をサポートする。

---

### 13.5.2 接続条件

以下を受け付ける。

```text
IPv4 Connection
```

```text
IPv6 Connection
```

---

### 13.5.3 target_family 判定

submitted_ip が存在しない場合:

```text
target_family
=
connection address family
```

---

### 13.5.4 例

IPv4 Connection:

```text
198.51.100.1
```

結果:

```text
target_family = ipv4
```

---

IPv6 Connection:

```text
2001:db8::1
```

結果:

```text
target_family = ipv6
```

---

## 13.6 submitted_ip 優先ルール

### 13.6.1 原則

submitted_ip が存在する場合、Endpoint Type より submitted_ip の Family を優先する。

---

### 13.6.2 判定式

```text
target_family
=
family(submitted_ip)
```

---

### 13.6.3 IPv4 Connection + IPv6 submitted_ip

入力:

```text
Connection:
198.51.100.5

submitted_ip:
2001:db8::1
```

結果:

```text
target_family = ipv6
```

---

### 13.6.4 IPv6 Connection + IPv4 submitted_ip

入力:

```text
Connection:
2001:db8::100

submitted_ip:
203.0.113.20
```

結果:

```text
target_family = ipv4
```

---

### 13.6.5 制限

実装は Connection Family と submitted_ip Family の不一致のみを理由に Request を拒否してはならない。

---

## 13.7 Address Registration

### 13.7.1 Connection Family

Connection Family は Registration Target を強制しない。

---

### 13.7.2 submitted_ip

submitted_ip が存在する場合、登録対象は submitted_ip によって決定する。

---

### 13.7.3 observed_ip

observed_ip を Device Address として扱ってはならない。

---

## 13.8 Endpoint Type Recording

### 13.8.1 目的

監査および運用分析のため。

---

### 13.8.2 AuditLog

実装は Endpoint Type を AuditLog に記録すべきである。

例:

```text
v4only
```

```text
v6only
```

```text
dualstack
```

---

### 13.8.3 利用目的

以下に利用してよい。

- 運用分析
- トラブルシュート
- 統計情報
- 移行分析

---

## 13.9 Metrics

実装は以下の Metrics を提供してよい。

```text
endpoint_requests_total
```

```text
endpoint_v4_requests_total
```

```text
endpoint_v6_requests_total
```

```text
endpoint_dualstack_requests_total
```

---

## 13.10 Endpoint Security

### 13.10.1 Endpoint Type

Endpoint Type は認証方式を変更してはならない。

---

### 13.10.2 Authorization

Endpoint Type に関係なく Authorization を実施しなければならない。

---

### 13.10.3 Family Validation

Address Family の検証は実施しなければならない。

---

## 13.11 Operational Considerations

### 13.11.1 Dual Stack 推奨

新規導入では dualstack Endpoint の利用を推奨する。

---

### 13.11.2 Legacy Support

v4only Endpoint は IPv4 専用環境向けに提供してよい。

---

### 13.11.3 IPv6 Migration

v6only Endpoint は IPv6 専用環境への移行支援として利用してよい。

---

## 13.12 不変条件

実装は以下を保証しなければならない。

```text
Endpoint Type
∈ {v4only, v6only, dualstack}
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
IPv4 Connection + IPv6 submitted_ip
→ allowed
```

```text
IPv6 Connection + IPv4 submitted_ip
→ allowed
```

```text
observed_ip
≠ Device Address
```

```text
Endpoint Type
→ AuditLog Recordable
```

これらを満たさない実装は本仕様に適合しない。
