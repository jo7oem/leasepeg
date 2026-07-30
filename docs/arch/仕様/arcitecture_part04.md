あｓｄ# 4. システム範囲

## 4.1 概要

本システムは、セルフホスト型 IP アドレス管理プラットフォームとして、Device のアドレス管理および監査機能を提供する。

本システムは以下を提供する。

- Device Address Update API
- ManagementGroup Address Read API
- User Management API
- Device Metadata 管理
- Credential 発行および失効
- AuditLog
- UpdateHistory
- CleanupJob
- External API Integration
- Monitoring API

---

## 4.2 提供機能

### 4.2.1 Device Address Management

Device の IPv4 および IPv6 Address を管理する。

実装は以下をサポートしなければならない。

- Address 登録
- Address 更新
- Address TTL 更新
- Address 削除
- Address Lifecycle 管理

---

### 4.2.2 Device Metadata Management

Device に任意 Metadata を付与できなければならない。

Metadata は以下の用途を想定する。

- デバイス名
- 設置場所
- 説明文
- 管理タグ

---

### 4.2.3 ManagementGroup

複数 Device を論理的にグループ化できなければならない。

ManagementGroup は以下の用途を想定する。

- システム単位
- 拠点単位
- サービス単位
- 運用単位

---

### 4.2.4 Credential Management

以下の Credential を管理する。

- UserSession
- DeviceCredential
- GroupReadToken

管理対象は以下を含む。

- 発行
- 失効
- 有効期限管理
- 利用履歴管理

---

### 4.2.5 Audit Function

重要な操作を監査可能にしなければならない。

システムは AuditLog を提供する。

---

### 4.2.6 History Function

変更履歴を保持しなければならない。

システムは UpdateHistory を提供する。

---

### 4.2.7 Cleanup Function

不要データの削除を行わなければならない。

システムは CleanupJob を提供する。

---

### 4.2.8 External Integration

外部システムと連携できなければならない。

実装方式は任意である。

例:

- Webhook
- REST API
- メッセージキュー

---

### 4.2.9 Monitoring

監視機能を提供しなければならない。

Monitoring API は以下を提供する。

- Liveness Check
- Readiness Check
- Metrics

---

## 4.3 管理対象

### 4.3.1 User

システム利用者。

---

### 4.3.2 Device

アドレス通知対象。

---

### 4.3.3 ManagementGroup

Device の集合。

---

### 4.3.4 Credential

認証情報。

---

### 4.3.5 Address

IPv4 または IPv6 Address。

---

### 4.3.6 Metadata

Key-Value Information。

---

### 4.3.7 AuditLog

監査ログ。

---

### 4.3.8 UpdateHistory

変更履歴。

---

## 4.4 サポート対象

### 4.4.1 IPv4

システムは IPv4 Address をサポートしなければならない。

---

### 4.4.2 IPv6

システムは IPv6 Address をサポートしなければならない。

---

### 4.4.3 Dual Stack

システムは Dual Stack をサポートしなければならない。

同一 Device は以下を同時保持できる。

```text
IPv4 Address
IPv6 Address
```

---

## 4.5 API の範囲

本システムは以下の API を提供する。

### Device API

DeviceCredential 向け API。

用途:

- Address Update
- Metadata Update

---

### Group Read API

GroupReadToken 向け API。

用途:

- Address Read

---

### User Management API

UserSession 向け API。

用途:

- Device 管理
- Group 管理
- Credential 管理
- AuditLog 参照

---

### Monitoring API

運用監視向け API。

用途:

- Health Check
- Metrics

---

### Internal API

内部運用向け API。

用途:

- CleanupJob
- Worker Management

---

## 4.6 非対象範囲

本システムは以下を提供しない。

---

### 4.6.1 DNS Server

本システムは DNS Server を提供しない。

---

### 4.6.2 Dynamic DNS Service

本システムは Dynamic DNS Service そのものを提供しない。

---

### 4.6.3 IP Ownership Verification

以下を検証しない。

- Address 所有権
- Address 到達性

---

### 4.6.4 Network Reachability

以下を実施しない。

- Ping
- Traceroute
- Port Scan

---

### 4.6.5 Certificate Management

以下は対象外とする。

- TLS Certificate 発行
- Certificate Renewal

---

### 4.6.6 Reverse Proxy

以下は対象外とする。

- Nginx 設定
- Apache 設定
- Envoy 設定

---

### 4.6.7 WAF

以下は対象外とする。

- WAF
- IDS
- IPS

---

### 4.6.8 Backup / Restore

Backup および Restore は対象外とする。

実装は以下を保証するのみとする。

- Backup を妨げない
- Restore を妨げない

---

### 4.6.9 High Availability

以下は本仕様の対象外とする。

- Multi Region
- Geo Replication
- Active Active 構成

---

## 4.7 デプロイメントスコープ

本システムは以下を対象とする。

### 小規模環境

- 個人利用
- 自宅環境

---

### 中規模環境

- 小規模企業
- 拠点管理

---

### 自己運用環境

- VPS
- オンプレミス
- Private Cloud

---

## 4.8 責務分離

システム責務を以下のように分離する。

| 領域 | 責務 |
|--------|--------|
| Application | Address 管理 |
| Database | 永続化 |
| Worker | 非同期処理 |
| Reverse Proxy | TLS 終端 |
| Monitoring | 監視 |
| Infrastructure | Backup |

---

## 4.9 スコープ外変更

以下の追加機能は将来拡張として扱う。

- Multi User Collaboration
- Shared ManagementGroup
- RBAC
- SSO
- Federation
- Notification Center
- DNS Integration
- DHCP Integration

本仕様では定義しない。
