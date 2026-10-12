# API 設計: 判断理由と進化パス

## 現状のエンドポイント（OpenAPI 3.0）

定義: `internal/generated/oapi/rest/rest.yaml`、生成サーバー: `internal/generated/oapi/rest/rest.gen.go`。

| メソッド | パス | ハンドラ | 認証 | 目的 |
|---------|------|---------|------|------|
| `GET` | `/healthcheck/live` | `GetHealthcheckLive` | ❌ | プロセス生存確認（k8s livenessProbe） |
| `GET` | `/healthcheck/ready` | `GetHealthcheckReady` | ❌ | 就緒確認（k8s readinessProbe） |
| `GET` | `/ping` | `GetPing` | ❌ | 送信元 IP 返却（疎通テスト） |

**レスポンススキーマ**: `healthcheckLiveResp`, `healthcheckReadyResp`, `Pong` —— すべて `rest.yaml` の components に定義。

## 設計判断の理由

### 1. Strict Server モード + 生成ハンドラ

```go
oapi.NewStrictHandler(handler{}, middlewares...)
```

**なぜ？**
- コンパイル時保証: ハンドラが全操作を実装していることを保証。
- `interface{}` なし —— 型付きリクエスト/レスポンスオブジェクト。
- ミドルウェアチェーン（`StrictMiddlewareFunc`）で横断的関心事を処理。

### 2. BasicAuth セキュリティスキーム（定義のみ、未適用）

```yaml
security:
  - basicAuth: []
components:
  securitySchemes:
    basicAuth:
      type: http
      scheme: basic
```

**なぜ定義だけで operation ごとに空？**
- 将来の認証へのプレースホルダー。
- OpenAPI 仕様に意図を文書化。
- 今適用すると `/healthcheck/*` と `/ping`（インフラプローブ）が壊れる。
- **判断**: 実エンドポイント追加時に operation 単位で適用。

### 3. IP 抽出ミドルウェア（`strictIPRecordHandler`）

`r.RemoteAddr` から送信元 IP を抽出 → context → `GetPing` ハンドラへ。

**なぜミドルウェアでハンドラロジックじゃない？**
- 関心の分離: トランスポート懸念（IP）vs ビジネスロジック。
- 将来、クライアント IP が必要なエンドポイントで再利用可能。
- Context ベース —— HTTP リクエストなしでテスト可能。

### 4. ヘルスチェック状態機械

`store/memory` が状態定義: `Starting → Ready ↔ Degraded → ShuttingDown`、`Error` は終端。

**なぜ自作状態機械でライブラリじゃない？**
- 依存最小化。
- 遷移がコード上で明示的。
- `Degraded` と `Error` を区別 —— 部分機能（例: DB 遅いがキャッシュ動作）を表現。

**レスポンスコード**:
- `200`: `Ready` または `Degraded`
- `503`: `Starting` または `ShuttingDown`
- `500`: `Error`

k8s セマンティクスと一致: 503 = トラフィック送らない、500 = 何か壊れている。

## 進化パス

### フェーズ 1: デバイス管理 API（予定）
```
POST   /devices              CreateDevice
GET    /devices              ListDevices (ユーザースコープ)
GET    /devices/{id}         GetDevice
PATCH  /devices/{id}         UpdateDevice
DELETE /devices/{id}         DeleteDevice
GET    /devices/{id}/records ListRecords
POST   /devices/{id}/records CreateRecord
PATCH  /records/{id}         UpdateRecord
DELETE /records/{id}         DeleteRecord
```

**認証**: Bearer token（JWT）または API key。BasicAuth は廃止。

### フェーズ 2: ユーザー/認証 API
```
POST   /auth/register        ローカル登録
POST   /auth/login           ローカルログイン → JWT
POST   /auth/refresh         トークンリフレッシュ
GET    /auth/providers       連携 OAuth 一覧
POST   /auth/link/{provider} OAuth 連携開始
DELETE /auth/unlink/{provider} OAuth 連携解除
```

### バージョニング戦略

**URL パスバージョニング**: フェーズ 1 から `/v1/...`。
- OpenAPI `servers[0].url` をバージョンごとに更新。
- 生成パッケージ: `oapi/v1/rest`（別 Go パッケージ）。
- 旧バージョンは非推奨まで維持。

**なぜヘッダーバージョニングじゃない？**
- URL に明示的 —— デバッグ、キャッシュ、ルーティングが容易。
- 標準的慣習（Stripe, GitHub, AWS）。

### エラーレスポンス形式（標準化予定）

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "human-readable",
    "details": [{"field": "name", "issue": "required"}]
  }
}
```

**なぜ RFC 7807 (Problem Details) じゃない？**
- 採用検討中。相互運用性高い。
- 判断保留 —— 現エンドポイントはエラー返さないため。

## 未解決 / 未決事項

1. **認証方式**: JWT（ステートレス）vs セッション（ステートフル）？ JWT は水平スケーリング向きだが失効戦略必要（ブロックリスト or 短期有効 + リフレッシュ）。
2. **レート制限**: どこで？ ミドルウェア（IP 単位）か認証連携（ユーザー単位）か？ まず IP 単位でミドルウェア実装。
3. **リクエストバリデーション**: oapi-codegen は path/query パラメータ検証。Body は？ 生成型に `go-playground/validator` タグ、またはカスタムミドルウェア。
4. **ページネーション**: カーソルベース（安定）vs オフセット？ `ListDevices`/`ListRecords` はカーソル推奨。
5. **Webhook/非同期**: デバイス IP 変更 —— push vs poll？ 現 `/ping` は poll。Webhook 購読 API が必要になる可能性。
6. **デバイス向け API key vs OAuth**: デバイスはプログラマティックアクセス必要。API key（プレフィックス + ハッシュ）の方が OAuth client credentials より簡潔。