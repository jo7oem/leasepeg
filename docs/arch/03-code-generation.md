# コード生成ワークフロー

## ツール

| ツール | 目的 | 設定 | 出力先 |
|--------|------|------|--------|
| `sqlc` | SQL → Go（型安全クエリ） | `sqlc.yaml` | `internal/generated/rdb/` |
| `oapi-codegen` | OpenAPI → Go（サーバー/クライアント） | `internal/generated/oapi/rest/generate.go` | `internal/generated/oapi/rest/rest.gen.go` |

両方とも `make setup` で `dev_tools/bin/` にインストール。

## sqlc 設定（`sqlc.yaml`）

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "sql/postgres/query.sql"
    schema: "sql/postgres/schema/public"
    gen:
      go:
        package: "rdb"
        out: "internal/generated/rdb"
        sql_package: "pgx/v5"
        output_db_file_name: "db_gen.go"
        output_models_file_name: "models_gen.go"
        output_files_suffix: "_gen"
```

**主な選択理由**:
- `pgx/v5` —— ネイティブ PostgreSQL ドライバー、`database/sql` のオーバーヘッドなし。
- `_gen` サフィックス —— 生成コードと手書きコードを明確に区別。
- 単一パッケージ `rdb` —— クエリとモデルを同梱（循環依存なし）。

## oapi-codegen 設定（`generate.go`）

```go
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest \
  -generate types,strict-server,chi-server \
  -package rest \
  -o rest.gen.go \
  rest.yaml
```

**主な選択理由**:
- `strict-server` —— コンパイル時インターフェース準拠保証。
- `chi-server` —— `chi.Router` 互換ハンドラ生成（実態は `http.ServeMux` ラッパーで使用）。
- `types` —— リクエスト/レスポンス構造体生成。
- **生成しない**: `client`（`oapi.yaml` に `client: true` あるが generate 指示に含めず —— 意図的）。

## 再生成コマンド

```bash
# SQL 層
make sqlc-generate

# API 層
cd internal/generated/oapi/rest && go generate

# 両方
make sqlc-generate && cd internal/generated/oapi/rest && go generate
```

**CI**: 両方 `make ci` → `make lint` で実行（専用の generate チェックなし；古ければコンパイルエラーになる）。

## なぜこの組み合わせ？

| 層 | 代替案 | 採用理由 |
|-----|--------|----------|
| DB | GORM, sqlx, ent | sqlc: SQL ファースト、ランタイムリフレクションゼロ、コンパイル時安全性 |
| API | 手書き, gRPC, GraphQL | oapi-codegen: OpenAPI 標準、HTTP/JSON ネイティブ、strict server |

**gRPC ではない理由**: 外部コンシューマーは HTTP/JSON 必要。内部だけなら gRPC もありだが複雑化（protobuf, gateway）。

**GraphQL ではない理由**: 現スコープでは過剰。REST + OpenAPI で十分。

## 生成ファイル —— 編集禁止

| ファイル | 再生成コマンド | 編集可？ |
|----------|---------------|---------|
| `internal/generated/rdb/models_gen.go` | `make sqlc-generate` | ❌ |
| `internal/generated/rdb/query.sql_gen.go` | `make sqlc-generate` | ❌ |
| `internal/generated/rdb/db_gen.go` | `make sqlc-generate` | ❌ |
| `internal/generated/oapi/rest/rest.gen.go` | `go generate` | ❌ |

**拡張パターン**:
- DB: `logic/` に手書きリポジトリを作り `*rdb.Queries` をラップ。
- API: `logic/rest/` に手書き `handler` 構造体を作り生成インターフェースを実装。

## バージョン固定

```makefile
$(DEV_BIN)/sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

$(DEV_BIN)/oapi-codegen:  # Makefile に未追加
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
```

**なぜ `@latest` で固定しない？**
- `dev_tools/bin/` はローカル、コミットしない。
- `make setup` で最新取得 —— 開発者がアップグレード制御。
- CI で固定したいなら `go.toolchain` ディレクティブや `go.work` で対応可能。

## 未解決 / 未決事項

1. **oapi-codegen を Makefile 化**: 現状サブディレクトリで `go generate`。整合性のため `make oapi-generate` ターゲット追加。
2. **クライアント生成**: `oapi.yaml` に `client: true` あるが generate 指示に含めず。外部 SDK 必要時に有効化。
3. **sqlc `emit_json_tags`**: 無効。モデルの JSON シリアライズ必要なら `emit_json_tags: true` 追加。
4. **複数 SQL ファイル**: 現状単一 `query.sql`。肥大化したらドメイン別（user, device, record）に分割？ sqlc は複数クエリファイル対応。
5. **生成テストヘルパー**: sqlc はモック用 `Querier` インターフェース生成可能。`emit_interface: true` を有効化検討。