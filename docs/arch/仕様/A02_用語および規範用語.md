---
Title: 2. 用語および規範用語
Description: 本章は「用語および規範用語」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「2. 用語および規範用語」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# 2. 用語および規範用語

## 2.1 規範用語

本仕様における以下の大文字表記は、BCP 14、RFC 2119 および RFC 8174 に記載された意味で解釈する。

- MUST
- MUST NOT
- REQUIRED
- SHALL
- SHALL NOT
- SHOULD
- SHOULD NOT
- RECOMMENDED
- NOT RECOMMENDED
- MAY
- OPTIONAL

日本語本文では以下の対応を用いる。

| 規範用語 | 日本語表現 |
|----------|-----------|
| MUST | しなければならない |
| MUST NOT | してはならない |
| SHOULD | すべきである |
| SHOULD NOT | すべきではない |
| MAY | してもよい |

大文字で表記されない同一単語については、通常の自然言語として解釈する。

---

## 2.2 用語定義

### User

人間の利用者。

Device、ManagementGroup、および各種 Credential を管理する主体。

---

### UserSession

User が管理操作を行うための短寿命セッション。

UserSession は人間向け Credential であり、機械利用を想定しない。

---

### Device

IP アドレスを通知する対象。

以下を含む。

- ルーター
- サーバー
- Linux ホスト
- エージェント
- 仮想マシン
- コンテナホスト

---

### DeviceCredential

Device が自身の情報を更新するために使用する認証情報。

DeviceCredential は自身が対応する Device のみ更新可能である。

---

### ManagementGroup

単一 User に属する Device の集合。

ManagementGroup はアドレス参照の論理的単位として使用する。

---

### GroupReadToken

ManagementGroup に属する Device の情報を参照するための読み取り専用 Credential。

---

### Credential

認証情報の総称。

本仕様では以下を含む。

- UserSession
- DeviceCredential
- GroupReadToken

---

### Address

Device に紐づく IP アドレス。

Address は IPv4 または IPv6 のいずれかである。

---

### Address Family

IP アドレス種別。

以下のいずれかを指す。

- IPv4
- IPv6

---

### submitted_ip

Device が API Request の中で明示的に申告した IP アドレス。

---

### observed_ip

HTTP Connection の Peer Address から観測された Source IP アドレス。

---

### connection address family

HTTP Connection の Peer Address から判定される Address Family。

---

### target_family

Address Update 操作で対象となる Address Family。

---

### Metadata

Device に紐付く UTF-8 Key-Value データ。

---

### AuditLog

監査目的で保存されるイベントログ。

以下を含む。

- 認証
- 認可
- Address 更新
- Metadata 更新
- Credential 発行
- Credential 失効
- 削除操作
- CleanupJob 実行
- External Integration 障害

---

### UpdateHistory

対象リソースの状態変化を記録する履歴。

AuditLog とは目的が異なる。

---

### CleanupJob

期限切れデータや不要データを削除する内部処理。

---

### External Integration

外部システムとの連携機能。

Webhook や外部 API 呼び出しを含む。

---

### active address

Address Record が存在し、

```text
expires_at > current_time
```

を満たす状態。

---

### expired address

Address Record が存在するが、

```text
expires_at <= current_time
```

を満たす状態。

---

### deleted address

Address Record が削除された状態。

通常 API から参照できない。

---

### Revoked Credential

失効済み Credential。

失効済み Credential は認証成功として扱ってはならない。

---

### Expired Credential

有効期限切れ Credential。

有効期限切れ Credential は認証成功として扱ってはならない。

---

### Token Prefix

Credential を識別するための非秘密情報。

監査や Rate Limit に使用できる。

---

### Authorization

認証済み主体が対象 Resource に対して操作可能か判定する処理。

---

### Authentication

Credential が正当であることを確認する処理。

---

### Least Privilege

最小権限原則。

主体には必要最低限の権限のみを与えるべきである。

---

### Ownership

Resource がどの User に属するかを表す概念。

---

### Invariant

常に成立しなければならない制約。

Invariant を破る操作は拒否しなければならない。

---

### Idempotency

同一または同等の Request が複数回実行されても、最終状態が不正にならない性質。

例:

```text
同じ IP アドレスを 3 回送信しても
最終状態は 1 回送信時と同じである。
```

---

### Eventual Consistency

非同期処理において、最終的に整合した状態へ収束する性質。

---

### Retry

通信失敗またはタイムアウト後の再実行。

---

### Rate Limit

一定期間内の Request 数を制限する仕組み。

---

### Endpoint Type

Address Family に応じた API Endpoint 種別。

本仕様では以下を含む。

- v4only
- v6only
- dualstack

---

### Internal API

システム内部利用を目的とする API。

外部公開を前提としない。

---

### Monitoring API

システム監視用 API。

ヘルスチェックや Metrics 提供を行う。

---

## 2.3 略語

| 略語 | 意味 |
|--------|--------|
| API | Application Programming Interface |
| DB | Database |
| IP | Internet Protocol |
| UUID | Universally Unique Identifier |
| TTL | Time To Live |
| CSRF | Cross-Site Request Forgery |
| OOF | Out Of Office |
| SRE | Site Reliability Engineering |
| UTF-8 | Unicode Transformation Format 8-bit |
| RFC | Request For Comments |

---

## 2.4 日時表現

特別な記載がない限り、日時は UTC を基準として扱う。

実装は内部保存時に UTC を使用することが推奨される。

---

## 2.5 UUID

外部公開 Identifier は UUIDv7 を使用する。

UUIDv7 の利用は推測困難性を向上させるが、認可制御の代替にはならない。
