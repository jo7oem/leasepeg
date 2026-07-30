# 3. 設計原則

## 3.1 基本方針

本システムは以下の設計原則に基づいて構築する。

- 権限分離
- 最小権限
- 認証情報の長寿命化抑制
- 監査可能性
- 結果整合性
- 運用容易性
- セキュア・バイ・デフォルト
- IPv4 / IPv6 対称性

実装は、利便性よりも安全性および運用可能性を優先しなければならない。

---

## 3.2 権限分離

本システムは認証主体ごとに権限を分離する。

以下の権限モデルを採用する。

| 主体 | Credential | 権限 |
|--------|--------|--------|
| Human User | UserSession | 管理操作 |
| Device | DeviceCredential | 自 Device 更新 |
| External Consumer | GroupReadToken | 読み取り専用 |

---

### 3.2.1 UserSession

UserSession は人間による管理操作専用である。

UserSession は以下を実行できる。

- Device 作成
- Device 更新
- Device 削除
- ManagementGroup 作成
- ManagementGroup 更新
- ManagementGroup 削除
- Credential 発行
- Credential 失効
- Address 参照
- AuditLog 参照

---

### 3.2.2 DeviceCredential

DeviceCredential は Device の自己更新専用である。

DeviceCredential は以下のみ実行できる。

- Address 更新
- Metadata 更新

DeviceCredential は以下を実行してはならない。

- Address 参照
- AuditLog 参照
- Device 管理
- Credential 管理
- User 管理

---

### 3.2.3 GroupReadToken

GroupReadToken は読み取り専用 Credential である。

GroupReadToken は以下のみ実行できる。

- Group Address Read

GroupReadToken は以下を実行してはならない。

- Update
- Delete
- Credential 管理
- Group 管理
- User 管理

---

### 3.2.4 Read / Write 分離

実装は Read-Write 両権限を持つ長寿命 Credential を発行してはならない。

以下を満たさなければならない。

```text
Read XOR Write
```

すなわち、

- Read 専用
- Write 専用

のいずれかでなければならない。

---

## 3.3 最小権限

すべての主体には必要最小限の権限のみを付与しなければならない。

---

### 3.3.1 DeviceCredential

DeviceCredential は対応する Device のみ更新できなければならない。

他 Device へのアクセスを許可してはならない。

---

### 3.3.2 GroupReadToken

GroupReadToken は自身に対応する ManagementGroup のみ参照できなければならない。

---

### 3.3.3 UserSession

UserSession は自身が所有する Resource のみ操作できなければならない。

---

### 3.3.4 UUID と認可

UUIDv7 を使用している場合でも認可チェックを省略してはならない。

以下は誤りである。

```text
UUID が推測できない
↓
Authorization 不要
```

正しい設計は以下である。

```text
UUIDv7
+
Authorization Check
```

---

## 3.4 所有権モデル

すべての Resource は単一 User に帰属する。

---

### 3.4.1 Device

Device は単一 User に所属する。

```text
Device
 └── Owner User
```

---

### 3.4.2 ManagementGroup

ManagementGroup は単一 User に所属する。

```text
ManagementGroup
 └── Owner User
```

---

### 3.4.3 Group Membership

ManagementGroup に追加できる Device は同一 User が所有する Device のみとする。

実装は以下を保証しなければならない。

```text
Device.owner_user_id
==
ManagementGroup.owner_user_id
```

この条件を満たさない Membership 追加は拒否しなければならない。

---

## 3.5 認証情報の長寿命化抑制

認証情報は必要以上に長寿命化してはならない。

---

### 3.5.1 UserSession

UserSession は短寿命でなければならない。

以下を満たさなければならない。

```text
maximum lifetime = 30 minutes
```

---

### 3.5.2 DeviceCredential

DeviceCredential は長寿命であってもよい。

ただし以下を満たさなければならない。

- expires_at を持つ
- revoke 可能
- last_used_at を持つ

---

### 3.5.3 GroupReadToken

GroupReadToken は長寿命であってもよい。

ただし以下を満たさなければならない。

- expires_at を持つ
- revoke 可能
- last_used_at を持つ

---

### 3.5.4 永続 Credential の禁止事項

実装は以下を禁止しなければならない。

```text
無期限 User Credential
```

すべての Credential は失効可能でなければならない。

---

## 3.6 セキュア・バイ・デフォルト

実装は安全な設定を既定値として採用しなければならない。

---

### 3.6.1 非公開 API

管理 API は認証必須でなければならない。

公開アクセスを許可してはならない。

---

### 3.6.2 Monitoring API

Monitoring API は既定で外部公開してはならない。

---

### 3.6.3 Internal API

Internal API は既定で外部公開してはならない。

---

### 3.6.4 Secret

認証情報を平文保存してはならない。

---

## 3.7 監査可能性

重要な操作は追跡可能でなければならない。

---

### 3.7.1 AuditLog

以下は AuditLog に記録しなければならない。

- Authentication
- Authorization Failure
- Credential Issue
- Credential Revoke
- Address Update
- Metadata Update
- Delete
- CleanupJob

---

### 3.7.2 UpdateHistory

変更履歴は UpdateHistory に保存する。

AuditLog と UpdateHistory は目的を分離しなければならない。

---

## 3.8 結果整合性

本システムは結果整合性を重視する。

---

### 3.8.1 Last Accepted Operation Wins

複数更新が発生した場合、最後に受理された操作結果を採用しなければならない。

---

### 3.8.2 Retry Safety

同一 Request の再送が発生しても不正状態にしてはならない。

---

### 3.8.3 External Integration

外部連携は非同期に実行しなければならない。

最終的に API Server の状態と整合することを保証しなければならない。

---

## 3.9 IPv4 / IPv6 対称性

IPv4 と IPv6 は可能な限り対称に扱う。

実装は以下を満たすべきである。

- IPv4 のみを特別扱いしない
- IPv6 のみを特別扱いしない
- API は両 Address Family を一貫した形で扱う

---

## 3.10 運用容易性

本システムは小規模環境でも運用できることを重視する。

運用担当者は以下を容易に実施できるべきである。

- Credential 失効
- Address 状態確認
- AuditLog 確認
- UpdateHistory 確認
- CleanupJob 実行
- External Integration 状態確認

以上を満たすことを設計原則とする。
