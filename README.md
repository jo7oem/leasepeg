# leasepeg

このリポジトリは **leasepeg** の実装です。

## データベースマイグレーション

### 前提条件
- Docker と Docker Compose がインストールされていること
- `make setup` で `golang-migrate` バイナリが `$(DEV_BIN)/migrate` にインストールされます

### マイグレーションディレクトリ
```
sql/postgres/migrations/
```
このディレクトリに `000001_<説明>.up.sql` と `000001_<説明>.down.sql` を作成します。サンプルは `sql/postgres/migrations/000001_init.up.sql` と `000001_init.down.sql` にあります。

### マイグレーションの実行
- ローカルで実行する場合（Makefile のターゲット）
  ```bash
  make migrate-up   # 最新のマイグレーションを適用
  make migrate-down # 1 ステップだけロールバック
  ```
- Docker Compose の `migrate` プロファイルを使う場合
  ```bash
  docker compose --profile migrate up migrate
  ```
  これによりコンテナ内で `migrate` が実行され、`up` が適用されます。

### 新しいマイグレーションの作成手順
1. `sql/postgres/migrations` に `NNNNN_<名前>.up.sql` と `NNNNN_<名前>.down.sql` を作成
2. `make migrate-up` で適用確認
3. 必要に応じて `docker compose --profile migrate up migrate` でも同様に適用できます

### 注意点
- データベースの接続文字列は `postgres://postgres:postgres@localhost:5432/mydatabase?sslmode=disable`（ローカル） または `postgres://postgres:postgres@postgres:5432/mydatabase?sslmode=disable`（Docker 内） です。
- 変更を加えるたびに `make migrate-up` で動作確認を行ってください。
