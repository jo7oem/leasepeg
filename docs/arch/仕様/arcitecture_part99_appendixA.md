---
Title: Appendix A. 実装ガイドライン（非規範）
Description: 本章は「実装ガイドライン（非規範）」を対象とし、適用範囲と満たすべき制約を規定する。
Summary: |
  「Appendix A. 実装ガイドライン（非規範）」の要件、判断基準、関連する前提を定義する。
  本章で規定する制約に違反する実装または運用は許容しない。
---

# Appendix A. 実装ガイドライン（非規範）

> 本付録は参考情報であり、規範仕様ではない。
>
> 本付録の内容が本文と矛盾する場合は本文を優先する。

---

# A.1 目的

本付録は、本仕様に準拠した実装を行う際の推奨事項を示す。

以下を目的とする。

- 実装者向け補足情報の提供
- 推奨アーキテクチャの提示
- 運用上のベストプラクティス共有

---

# A.2 推奨技術スタック

## A.2.1 API Server

推奨例:

```text
Go 1.24+
```

Web Framework 例:

```text
net/http
chi
gin
echo
```

---

## A.2.2 Database

推奨:

```text
PostgreSQL 16+
```

---

## A.2.3 Queue

実装例:

```text
PostgreSQL Queue
```

```text
Redis
```

```text
RabbitMQ
```

---

## A.2.4 Monitoring

推奨:

```text
Prometheus
```

```text
Grafana
```

---

# A.3 推奨ディレクトリ構成

```text
cmd/
  api/
  worker/

internal/
  auth/
  device/
  group/
  metadata/
  audit/
  integration/

pkg/
  api/
  logging/

migrations/

configs/

docs/
```

---

# A.4 API Server 推奨構成

```text
HTTP Server
 ├─ Middleware
 ├─ Authentication
 ├─ Authorization
 ├─ Validation
 ├─ Service Layer
 └─ Repository Layer
```

---

# A.5 Transaction Guidelines

以下は単一 Transaction を推奨する。

- Device Create
- Device Delete
- Membership Change
- Credential Revoke
- Address Update

例:

```text
BEGIN

Authorization Check

Update

AuditLog Insert

COMMIT
```

---

# A.6 Logging Guidelines

## 推奨 Log Level

```text
DEBUG
INFO
WARN
ERROR
```

---

## Log 出力例

```json
{
  "request_id": "...",
  "actor": "device",
  "operation": "address_update",
  "result": "success"
}
```

---

## 出力禁止

以下を出力しない。

```text
Token
Password
Session Secret
Digest Secret
```

---

# A.7 Migration Guidelines

推奨:

```text
Forward Only Migration
```

---

推奨ツール例:

```text
golang-migrate
```

```text
goose
```

---

# A.8 Testing Guidelines

## Unit Test

対象:

- Validation
- Authorization
- Lifecycle Logic
- Metadata Logic

---

## Integration Test

対象:

- Database
- API
- Worker

---

## Security Test

対象:

- Authentication
- Authorization
- Rate Limit
- CSRF

---

# A.9 Backup Considerations

本仕様は Backup を規定しないが、以下を推奨する。

対象:

- PostgreSQL
- AuditLog
- UpdateHistory

---

推奨:

```text
Daily Backup
```

```text
Encrypted Storage
```

---

# A.10 Deployment Examples

## Single Node

```text
Nginx
 └─ API
 └─ Worker
 └─ PostgreSQL
```

---

## High Scale

```text
Load Balancer
 ├─ API
 ├─ API
 ├─ API

Queue

Worker
Worker
Worker

PostgreSQL
```

---

# A.11 Operational Checklist

- [ ] HTTPS 有効
- [ ] Monitoring 有効
- [ ] AuditLog 保持確認
- [ ] CleanupJob 自動実行確認
- [ ] Backup 実施確認
- [ ] Credential Rotation 手順確認
- [ ] Alert 動作確認

---

# A.12 Appendix Completion

本付録は実装者向け補足情報を定義する。

本文の要件を満たす限り、実装方式は自由である。

---
End of Appendix A
