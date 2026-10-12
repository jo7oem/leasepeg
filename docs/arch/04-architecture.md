# アーキテクチャ: レイヤー構造とデータフロー

## 全体レイヤー

```
┌─────────────────────────────────────────────────────────────┐
│  cmd/                    (エントリーポイント: server, client) │
├─────────────────────────────────────────────────────────────┤
│  server/rest/            (HTTP ルーティング, ミドルウェア)    │
├─────────────────────────────────────────────────────────────┤
│  logic/rest/             (ハンドラ: ビジネスロジック)         │
├─────────────────────────────────────────────────────────────┤
│  internal/generated/                                              │
│  ├── oapi/rest/          (OpenAPI: 型, インターフェース, サーバー)│
│  └── rdb/                (sqlc: Queries, モデル)             │
├─────────────────────────────────────────────────────────────┤
│  store/                  (状態: health, 将来: cache/DB)      │
│  └── memory/             (インメモリ実装)                    │
├─────────────────────────────────────────────────────────────┤
│  logging/                (構造化 slog ラッパー)              │
└─────────────────────────────────────────────────────────────┘
```

**依存ルール**: 外側が内側に依存。内側は外側をインポートしない。
- `cmd` → `server` → `logic` → `internal/generated` + `store`
- `store/memory` は依存なし（ポータブル）

## コンポーネント責務

### cmd/server.go
- 設定読み込み（`config.yaml` → 構造体）
- DB 接続（pgxpool）またはメモリストア選択
- ロガー初期化
- HTTP サーバー ライフサイクル（グレースフルシャットダウン）
- 配線: `store` → `logic` → `server` → `http.Server`

### server/rest/register.go
- `http.ServeMux` 合成
- `logic/rest` ハンドラを `/` にマウント
- リクエストログミドルウェア追加（`logging.Logger` 経由）

### logic/rest/rest.go
- `oapi.StrictServerInterface`（生成）を実装
- **DB 直接呼ばない** —— リポジトリ（将来）または store に委譲
- 現状: health/ping のみ → `store/memory` + context IP 抽出使用
- ミドルウェア: `strictIPRecordHandler` がクライアント IP を context に注入

### store/memory/
- ミューテックス保護された状態機械 `HealthCheckState`
- `GetHealthCheckState()` / `SetHealthCheckState()`
- **インターフェース**: まだなし（暗黙）。PostgreSQL 実装追加時に `HealthStore` インターフェース定義予定。

### internal/generated/rdb/
- `Queries` 構造体: 全 DB 操作（テーブルごと CRUD）
- モデル: `User`, `Device`, `Record`, `Credential`（pgx 型）
- **ビジネスロジックなし** —— 純粋データアクセス

### internal/generated/oapi/rest/
- `StrictServerInterface`: 全操作をメソッド化
- エンドポイントごとのリクエスト/レスポンス型
- `Handler()`: `http.Handler` コンストラクタ

## データフロー例

### GET /healthcheck/ready
```
HTTP Request
    → server/rest (ServeMux)
    → logic/rest.handler.GetHealthcheckReady()
    → store/memory.GetHealthCheckState()
    → 状態返却 → 200/503/500 にマッピング
    → JSON レスポンス
```

### 将来: POST /devices
```
HTTP Request
    → server/rest
    → logic/rest.handler.CreateDevice()
    → (バリデーション) → rdb.Queries.CreateDevice(ctx, params)
    → pgxpool → PostgreSQL
    → Device 返却 → API レスポンスにマッピング
    → JSON レスポンス
```

## なぜこのレイヤー分け？

| レイヤー | 責務 | テスタビリティ |
|---------|------|---------------|
| `cmd` | 合成, ライフサイクル | 統合テストのみ |
| `server` | HTTP 懸念（ルーティング, ミドルウェア） | `httptest` |
| `logic` | ビジネスルール, オーケストレーション | モック store/DB で単体テスト |
| `generated` | トランスポート/データ契約 | コンパイル時検証 |
| `store` | 状態永続化の抽象化 | 実装交換可能（memory ↔ postgres） |

**KISS の判断理由**:
- DI フレームワーク不使用 —— `cmd/server.go` で明示的配線。
- インターフェース爆発させない —— 2 つ目の実装ができてから定義。
- 生成コードを「契約」とみなす —— logic は実装するだけ、ラップしない。

## 現状のギャップ / 技術的負債

1. **リポジトリ層なし**: `logic` が `rdb.Queries` を直接呼ぶ（未実装）。追加時は `logic/repository/` か `store/rdb/` に配置。
2. **Config → DB 接続**: `config.yaml` に `database.type: memory|postgres` あるが `cmd/server.go` 未完全配線。`store/memory` のみ。
3. **Context 伝播**: `logging.Logger` を手動で渡している。リクエストスコープロガーには `context.WithValue` 検討。
4. **エラーハンドリング**: 統一エラー型なし。`logic` が `error` 返す → `oapi` が 500 にマップ。HTTP ステータス適切マッピング用 `AppError` (コード付き) 必要。

## 未解決 / 未決事項

1. **リポジトリパターンの配置**: `logic/repository/`（ドメイン）vs `store/rdb/`（インフラ）？ 推奨: `store/rdb/` —— `logic` を純粋に保つ。
2. **トランザクション管理**: `pgx.Tx` の所有者は？ リポジトリメソッドは `*Queries` 受取（プール使用）。複数テーブル操作時は `WithTransaction(ctx, fn)` ヘルパー必要。
3. **オブザーバビリティ**: メトリクス（Prometheus）、トレーシング（OTel） —— どこに注入？ `server/rest` ミドルウェアか `cmd`？
4. **グレースフルシャットダウン**: `store/memory` に `StateShuttingDown` あり。PostgreSQL store にはコネクションドレイン必要。
5. **マルチテナンシー**: 現スキーマは devices に `user_id`。API にはユーザーコンテキスト必要（認証から）。どこで強制？ ミドルウェア → context → リポジトリ。