---
Title: 6. エンティティモデル
Description: 本章は「エンティティモデル」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「6. エンティティモデル」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 6. エンティティモデル

## 6.1 概要

本システムは以下のエンティティによって構成される。

- User
- UserSession
- Device
- DeviceCredential
- ManagementGroup
- GroupReadToken
- DeviceAddress
- DeviceMetadata
- AuditLog
- UpdateHistory
- ExternalIntegration
- CleanupJob

各エンティティは明確な責務を持ち、責務の重複を避けなければならない。

---

## 6.2 エンティティ一覧

### 6.2.1 User

システム利用者を表す。

User は以下を所有する。

- Device
- ManagementGroup
- UserSession

---

### 6.2.2 UserSession

User が管理操作を行うためのセッション。

UserSession は短寿命でなければならない。

---

### 6.2.3 Device

Address を通知する対象。

例:

- Router
- Linux Host
- VM
- Container Host

---

### 6.2.4 DeviceCredential

Device の自己更新専用 Credential。

---

### 6.2.5 ManagementGroup

複数 Device をグループ化するための論理単位。

---

### 6.2.6 GroupReadToken

ManagementGroup 用の読み取り専用 Credential。

---

### 6.2.7 DeviceAddress

Device に割り当てられた IP Address。

---

### 6.2.8 DeviceMetadata

Device に付与される任意 Key-Value 情報。

---

### 6.2.9 AuditLog

監査イベント。

---

### 6.2.10 UpdateHistory

状態変更履歴。

---

### 6.2.11 ExternalIntegration

外部システム連携設定。

---

### 6.2.12 CleanupJob

定期クリーンアップ処理。

---

# 6.3 所有関係

## 6.3.1 User

User は最上位所有者である。

```text
User
 ├─ Device
 ├─ ManagementGroup
 └─ UserSession
```

---

## 6.3.2 Device Ownership

Device は単一 User に所属する。

```text
Device
 └─ Owner User
```

---

## 6.3.3 ManagementGroup Ownership

ManagementGroup は単一 User に所属する。

```text
ManagementGroup
 └─ Owner User
```

---

## 6.3.4 Credential Ownership

Credential の所属関係は以下とする。

```text
UserSession
 └─ User

DeviceCredential
 └─ Device

GroupReadToken
 └─ ManagementGroup
```

---

# 6.4 所有権制約

以下の制約を満たさなければならない。

---

## 6.4.1 Device Ownership Constraint

```text
Device.owner_user_id
!= NULL
```

---

## 6.4.2 Group Ownership Constraint

```text
ManagementGroup.owner_user_id
!= NULL
```

---

## 6.4.3 Membership Ownership Constraint

ManagementGroup に所属できる Device は同一 User が所有するもののみとする。

以下を満たさなければならない。

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

---

## 6.4.4 Rejection Rule

以下を拒否しなければならない。

```text
User A の Device

↓

User B の Group に追加
```

---

# 6.5 カーディナリティ

## 6.5.1 User

```text
User
 ├─ UserSession      0..N
 ├─ Device           0..N
 └─ ManagementGroup  0..N
```

---

## 6.5.2 Device

```text
Device
 ├─ DeviceCredential 0..N
 ├─ DeviceAddress    0..2
 └─ DeviceMetadata   0..3
```

---

### DeviceAddress

最大 2 件保持できる。

```text
IPv4 : 0..1
IPv6 : 0..1
```

---

## 6.5.3 ManagementGroup

```text
ManagementGroup
 ├─ Device          0..N
 └─ GroupReadToken  0..N
```

---

# 6.6 Device モデル

## 6.6.1 責務

Device は Address 更新の主体である。

---

## 6.6.2 Device Identifier

Device は公開 Identifier を持つ。

```text
public_id
```

型:

```text
UUIDv7
```

---

## 6.6.3 Device Name

Device は任意名称を持つ。

例:

```text
router-a
vpn-gateway
edge-node-1
```

---

# 6.7 DeviceCredential モデル

## 6.7.1 責務

DeviceCredential は Device を認証する。

---

## 6.7.2 権限

以下のみ許可する。

```text
Address Update
Metadata Update
```

---

## 6.7.3 禁止事項

以下を許可してはならない。

```text
Address Read
Device Management
User Management
AuditLog Read
```

---

# 6.8 ManagementGroup モデル

## 6.8.1 責務

Address Read を行う対象 Device 集合を定義する。

---

## 6.8.2 所属

ManagementGroup は単一 User に所属する。

---

## 6.8.3 Device Membership

同一 Device は複数 Group に所属してよい。

例:

```text
Device A
 ├─ Group X
 └─ Group Y
```

---

# 6.9 GroupReadToken モデル

## 6.9.1 責務

ManagementGroup の参照専用 Credential。

---

## 6.9.2 権限

許可:

```text
Address Read
```

---

禁止:

```text
Update
Delete
Management
```

---

# 6.10 DeviceAddress モデル

## 6.10.1 責務

Device の実 Address を保持する。

---

## 6.10.2 Address Family

以下のいずれか。

```text
ipv4
ipv6
```

---

## 6.10.3 一意性

以下を満たさなければならない。

```text
UNIQUE(device_id, family)
```

---

## 6.10.4 状態

Address は以下の状態を持つ。

```text
active
expired
deleted
```

---

# 6.11 DeviceMetadata モデル

## 6.11.1 責務

任意の Device 情報を保持する。

---

## 6.11.2 制限

1 Device あたり最大:

```text
3 Key-Value Pairs
```

---

## 6.11.3 Key Constraint

```text
UTF-8
<= 256 bytes
```

---

## 6.11.4 Value Constraint

```text
UTF-8
<= 128 KiB
```

---

# 6.12 AuditLog モデル

## 6.12.1 責務

監査情報を保持する。

---

## 6.12.2 特徴

AuditLog は削除されにくい。

以下によって削除してはならない。

```text
Device Delete
Group Delete
```

---

# 6.13 UpdateHistory モデル

## 6.13.1 責務

状態遷移履歴を保持する。

---

## 6.13.2 AuditLog との違い

AuditLog:

```text
誰が何をしたか
```

UpdateHistory:

```text
何がどう変わったか
```

---

# 6.14 ExternalIntegration モデル

## 6.14.1 責務

外部システムとの連携設定を保持する。

---

## 6.14.2 状態

以下を持つ。

```text
enabled
disabled
suspended
```

---

# 6.15 CleanupJob モデル

## 6.15.1 責務

不要データ削除。

---

## 6.15.2 対象

以下を処理する。

- expired Address
- old UpdateHistory
- expired Credential
- retained AuditLog

---

# 6.16 エンティティ不変条件

実装は以下を常に保証しなければならない。

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

```text
DeviceCredential
→ single Device
```

```text
GroupReadToken
→ single ManagementGroup
```

```text
UserSession
→ single User
```

```text
DeviceAddress
→ single Device
```

```text
Address Family
∈ { ipv4, ipv6 }
```

これらの条件を破る操作は拒否しなければならない。
