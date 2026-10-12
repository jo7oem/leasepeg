# 運用: コマンド, 設定, テスト, CI/CD

## 設定

**ファイル**: `config.yaml`（`config.yaml.example` からコピー）

```yaml
server:
  host: "0.0.0.0"
  port: 8080

logging:
  level: "info"      # debug, info, warn, error
  format: "json"     # json, text

database:
  type: "memory"     # memory, postgres
  host: "localhost"
  port: 5432
  user: ""
  password: ""
  name: ""
```

**なぜ YAML で env-only じゃない？**
- 構造化設定（ネストセクション）が Go 構造体にクリーンマップ。
- `viper` や `koanf` で環境変数オーバーレイ可能（`SERVER_PORT`, `DATABASE_HOST` 等）。
- 例ファイルが全オプションを文書化。

## 開発ワークフロー

### ローカル開発開始
```bash
# ターミナル 1: PostgreSQL（postgres モード時）
podman compose up postgres

# ターミナル 2: ホットリロード付きアプリ（air）
make setup  # 初回のみ
air         # .go ファイル監視, リビルド, 再起動
```

### テスト実行
```bash
make test
# 同等: go test --race -parallel 8 ./...
```

**なぜ `--race -parallel 8`？**
- レース検出器がデータ競合を検知（並行ヘルスステートで重要）。
- Parallel 8: 一般的開発マシンで CPU 使用率と速度のバランス。

### Lint/Format
```bash
make fmt      # 修正: golangci-lint --fix + yamlfmt
make go-lint  # チェックのみ
make lint     # 全部 (go + redoc + yaml)
```

**なぜ golangci-lint v2？**
- 高速、Go 1.22+ ネイティブ対応。
- 設定は `.golangci.yml`（ここでは未表示）。

## データベース操作

### SQL 層生成
```bash
make sqlc-generate
# 入力: sql/postgres/query.sql + schema/
# 出力: internal/generated/rdb/*_gen.go
```

### マイグレーション作成
```bash
make migrate-new MIGRATION_NAME=add_device_index
# 作成: sql/postgres/migrations/000002_add_device_index.up/down.sql
```

### マイグレーション適用（ローカル）
```bash
# 要: 稼働中 postgres（podman compose up postgres）
make migrate-up
```

### ロールバック
```bash
make migrate-down  # 1 つ戻す
```

**マイグレーション命名規則**: `NNNNNN_動詞_名詞`（例: `000002_add_device_idx`, `000003_alter_records_add_col`）。

**なぜ起動時自動マイグレーションじゃない？**
- 本番での明示的制御。
- 適用前 SQL レビュー。
- ロールバック機能必須。

## API 生成

```bash
cd internal/generated/oapi/rest && go generate
# 入力: rest.yaml
# 出力: rest.gen.go
```

**Makefile 追加予定**（TODO）:
```makefile
.PHONY: oapi-generate
oapi-generate:
	cd internal/generated/oapi/rest && go generate
```

## CI/CD パイプライン（GitHub Actions 前提）

現 `Makefile` ターゲットが CI ステージに対応:

| ステージ | コマンド | 目的 |
|---------|---------|------|
| `lint` | `make lint` | 静的解析, フォーマット |
| `test` | `make test` | 単体テスト + レース検出 |
| `build` | `make build` | バイナリビルド（バージョン ldflags 付き） |

**Redoc lint** (`make redoc-lint`): `redoc-cli` でコンテナ内 OpenAPI スペック検証。

**バージョン注入**:
```makefile
VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
```
バイナリは `main.version` 変数でバージョン報告。

## 本番デプロイ考慮事項

### バイナリ
```bash
make build
# 出力: ./leasepeg（static 的, pgx で cgo あり）
```

### コンテナ（未定義）
```dockerfile
# Multi-stage: build → scratch/distroless
# config.yaml, migrate バイナリコピー, マイグレーション実行 → サーバー起動
```

### 本番での DB マイグレーション
- `migrate-up` を init container または別ジョブで実行。
- アプリは DB 就緒（ヘルスチェック）待機後サーブ開始。

### ヘルスチェック
- Liveness: `GET /healthcheck/live`（プロセス生存）
- Readiness: `GET /healthcheck/ready`（依存関係 OK）

**K8s プローブ**:
```yaml
livenessProbe:
  httpGet:
    path: /healthcheck/live
    port: 8080
  initialDelaySeconds: 5
readinessProbe:
  httpGet:
    path: /healthcheck/ready
    port: 8080
  initialDelaySeconds: 3
```

## 未解決 / 未決事項

1. **シークレット管理**: `config.yaml` に `password: ""`。環境変数, Vault, SealedSecrets どれを使う？
2. **ログ集約**: JSON 形式準備済み。Vector/Fluent Bit 設定必要。
3. **メトリクスエンドポイント**: 未公開。`server/rest` に `/metrics`（Prometheus）追加検討？
4. **グレースフルシャットダウンタイムアウト**: 設定可能？ 現状 `cmd/server.go` でハードコード。
5. **バックアップ/リストア**: PostgreSQL バックアップ戦略（pg_dump, PITR）未文書化。
6. **負荷テスト**: ベースラインなし。`hey` または `k6` スクリプトを `/ping` と将来エンドポイント用に追加。