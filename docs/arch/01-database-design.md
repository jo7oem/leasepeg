# データベース設計: 判断理由と制約

## 概要

PostgreSQL + sqlc（型安全なSQL生成）。スキーマは `sql/postgres/schema/public/` + `query.sql`、マイグレーションは `sql/postgres/migrations/` で管理。

## テーブル設計の判断理由

### users
| カラム | 型 | 判断理由 |
|--------|-----|----------|
| `id` | `uuid` | 標準的な UUID v4（デフォルト）。ユーザー作成頻度が低いため v4 で十分。v7 不要。 |
| `user_name` | `varchar(32)` UK | 32 文字：可読性とインデックス効率のバランス。ログイン識別子として一意。 |
| `display_name` | `text` | ログイン名とは分離。Unicode 対応、長い表示名も許容。 |
| `created_at`/`updated_at` | `timestamptz` default `now()` | 監査証跡。`updated_at` はトリガーではなく UPDATE クエリで明示的に更新（制御を明確化）。 |

**なぜ uuidv7 じゃない？** ユーザー作成は低頻度。v4 のランダム性で十分。v7 は `uuidv7()` 拡張への依存が増えるだけ。

### devices
| カラム | 型 | 判断理由 |
|--------|-----|----------|
| `id` | `uuid` DEFAULT `uuidv7()` | **時系列順 UUID**：(1) 作成時刻で自然ソート、(2) B-tree インデックスの局所性向上、(3) デバッグ時の可読性。`uuid-ossp` または pgcrypto が必要。 |
| `user_id` | `uuid` FK→users | 所有権。CASCADE DELETE は使わない——孤児デバイスは監査のため残す。 |
| `name`/`memo` | `text` | 自由形式。長さ制約なし（PostgreSQL `text` = 無制限 varchar）。 |

**なぜここだけ uuidv7？** デバイスは頻繁に作成される（ユーザーあたり複数の可能性）。時系列順のメリットが拡張依存を上回る。

### records
| カラム | 型 | 判断理由 |
|--------|-----|----------|
| `id` | `uuid` DEFAULT `uuidv7()` | devices と同理由。 |
| `device_id` | `uuid` FK→devices | 各レコードは 1 デバイスに属する。 |
| `type` | `text` | 識別子（"ipv4", "ipv6", "hostname" 等）。スキーマ変更なしで拡張可能。 |
| `ip_address` | `inet` | **PostgreSQL ネイティブ `inet` 型**：CIDR 操作、インデックス、バリデーションが組み込み。 |
| `reported_at` | `timestamptz` DEFAULT `now()` | デバイスが報告した時刻。DB 挿入時刻（`created_at`）とは異なる可能性あり。 |

**ユニーク制約**: `(device_id, type)` — デバイスごとにタイプ別で 1 レコードのみ。「タイプごとの最新 IP」を強制。更新は `UpdateRecord` で行う。

### credentials
| カラム | 型 | 判断理由 |
|--------|-----|----------|
| `user_id` + `provider` | PK（複合） | ユーザーあたりプロバイダー 1 つ。 |
| `provider` | `text` | "local", "github", "google" 等。拡張可能。 |
| `secret_hash` | `text` NULLABLE | Bcrypt/argon2 ハッシュ。OAuth 時は NULL（秘密なし）。 |
| `subject` | `text` | プロバイダー側のユーザー ID（OAuth `sub`）。 |
| `created_at`/`updated_at` | `timestamptz` | ローテーション監査。 |

**制約**:
- `UK(provider, subject)` — 重複 OAuth アカウント防止。
- `CHECK (provider <> 'local' OR secret_hash IS NOT NULL)` — ローカル認証時はパスワード必須を DB レベルで強制。

**なぜ別テーブル？** ユーザーあたり複数認証方式（ローカル + OAuth）に対応。users テーブルに nullable カラムを増やさない。将来の passkeys/WebAuthn にも対応可能。

## クエリ層（sqlc）

`internal/generated/rdb/` に生成。各クエリが `*Queries` のメソッドに 1:1 対応。

**命名規則**: `GetXxx`（単体）、`ListXxx`（複数）、`CreateXxx`、`UpdateXxx`、`DeleteXxx`。sqlc アノテーション（`:one`, `:many`, `:exec`）と一致。

**なぜ ORM じゃない？**
- `query.sql` に SQL が露出 —— レビュー可能、最適化可能。
- ランタイム リフレクション オーバーヘッド ゼロ。
- SQL から Go 型を生成 —— コンパイル時にスキーマドリフト検出。

## マイグレーション戦略

- ツール: `golang-migrate/migrate`（CLI は `dev_tools/bin/migrate`）。
- ファイル: `sql/postgres/migrations/` に `NNNNNN_name.up.sql` / `.down.sql`。
- `create-migration.sh` が自動連番（6 桁ゼロ埋め）。
- **Down マイグレーション必須** —— 「破壊のみ」の変更は不可。

**なぜ sqlc のマイグレーションじゃない？** sqlc はマイグレーション管理しない。golang-migrate は成熟しており、トランザクション対応、CI/CD 向き。

## 未解決 / 未決事項

1. **`uuidv7()` 可用性**: PostgreSQL 16+ または `uuid-ossp`/`pgcrypto` 拡張が必要。ターゲット PG バージョンを確認。
2. **論理削除？**: 現状 `deleted_at` なし。必要なら全テーブルに追加 + 部分インデックス。
3. **`records` のパーティショニング？** 書き込み量が多い場合、`reported_at` で時間ベースパーティションが必要に。
4. **`records.type` を enum 化？** 現在はフリーテキスト。固定値なら `CHECK (type IN (...))` または `record_types` 参照テーブルを検討。
5. **認証情報ローテーション ポリシー**: `updated_at` で追跡するが、期限/強制ロジックは未実装。