# Database Design: Rationale and Constraints

## Overview

PostgreSQL + sqlc (type-safe SQL generation). Schema defined in `sql/postgres/schema/public/` + `query.sql`, migrations in `sql/postgres/migrations/`.

## Table Design Decisions

### users
| Column | Type | Rationale |
|--------|------|-----------|
| `id` | `uuid` | Standard UUID v4 (default), not v7. Users created less frequently; v4 sufficient. |
| `user_name` | `varchar(32)` UK | 32 chars: balances readability with index efficiency. Unique for login identity. |
| `display_name` | `text` | Separate from login name; allows Unicode, longer display names. |
| `created_at`/`updated_at` | `timestamptz` default `now()` | Audit trail. `updated_at` manually maintained in UPDATE queries (not trigger) for explicit control. |

**Why not uuidv7?** User creation is low-frequency; v4 randomness acceptable. v7 adds dependency on `uuidv7()` extension.

### devices
| Column | Type | Rationale |
|--------|------|-----------|
| `id` | `uuid` DEFAULT `uuidv7()` | **Time-ordered UUID** for: (1) natural sort by creation time, (2) better B-tree index locality, (3) debugging readability. Requires `uuid-ossp` or pgcrypto. |
| `user_id` | `uuid` FK→users | Ownership. Cascade delete not used—orphaned devices retained for audit. |
| `name`/`memo` | `text` | Free-form; no length constraint (PostgreSQL `text` = unbounded varchar). |

**Why uuidv7 only here?** Devices created frequently (per-user, potentially many). Time-ordering benefits outweigh extension dependency.

### records
| Column | Type | Rationale |
|--------|------|-----------|
| `id` | `uuid` DEFAULT `uuidv7()` | Same rationale as devices. |
| `device_id` | `uuid` FK→devices | Each record belongs to one device. |
| `type` | `text` | Discriminator (e.g., "ipv4", "ipv6", "hostname"). Extensible without schema change. |
| `ip_address` | `inet` | **Native PostgreSQL `inet` type**: CIDR operations, indexing, validation built-in. |
| `reported_at` | `timestamptz` DEFAULT `now()` | When device reported. Not `created_at`—may differ from DB insert time. |

**Unique Constraint**: `(device_id, type)` — one active record per type per device. Enforces "latest IP per type" semantics. UPDATE used for refresh (see `UpdateRecord`).

### credentials
| Column | Type | Rationale |
|--------|------|-----------|
| `user_id` + `provider` | PK (composite) | One credential per provider per user. |
| `provider` | `text` | "local", "github", "google", etc. Extensible. |
| `secret_hash` | `text` NULLABLE | Bcrypt/argon2 hash. NULL for OAuth (no secret). |
| `subject` | `text` | Provider's user ID (OAuth `sub`). |
| `created_at`/`updated_at` | `timestamptz` | Rotation audit. |

**Constraints**:
- `UK(provider, subject)` — prevents duplicate OAuth accounts.
- `CHECK (provider <> 'local' OR secret_hash IS NOT NULL)` — enforces password for local auth at DB level.

**Why separate table?** Multiple auth methods per user (local + OAuth). Avoids nullable columns on users table. Future-proof for passkeys/WebAuthn.

## Query Layer (sqlc)

Generated in `internal/generated/rdb/`. Each query maps 1:1 to a method on `*Queries`.

**Naming convention**: `GetXxx` (single), `ListXxx` (multiple), `CreateXxx`, `UpdateXxx`, `DeleteXxx`. Matches sqlc annotations (`:one`, `:many`, `:exec`).

**Why not ORM?** 
- Explicit SQL visible in `query.sql` — reviewable, optimizable.
- Zero runtime reflection overhead.
- Go types generated from SQL — compile-time schema drift detection.

## Migration Strategy

- Tool: `golang-migrate/migrate` (CLI in `dev_tools/bin/migrate`).
- Files: `NNNNNN_name.up.sql` / `.down.sql` in `sql/postgres/migrations/`.
- `create-migration.sh` auto-numbers (6 digits, zero-padded).
- **Down migrations required** — no "destructive only" changes.

**Why not sqlc's migration?** sqlc doesn't manage migrations. golang-migrate is mature, supports transactions, CI/CD friendly.

## Unresolved / Open Questions

1. **`uuidv7()` availability**: Requires PostgreSQL 16+ or `uuid-ossp`/`pgcrypto` extension. Confirm target PG version.
2. **Soft delete?**: No `deleted_at` currently. If needed, add to all tables + partial indexes.
3. **Partitioning for `records`?** High write volume? Time-based partitioning on `reported_at` may be needed later.
4. **`records.type` enum?** Currently free text. Consider `CHECK (type IN (...))` or separate `record_types` table if fixed set.
5. **Credential rotation policy**: `updated_at` tracks rotation but no expiry/enforcement logic yet.