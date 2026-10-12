# Operations: Commands, Config, Testing, CI/CD

## Configuration

**File**: `config.yaml` (copy from `config.yaml.example`)

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

**Why YAML not env-only?**
- Structured config (nested sections) maps cleanly to Go struct.
- `viper` or `koanf` can overlay env vars (`SERVER_PORT`, `DATABASE_HOST`).
- Example file documents all options.

## Development Workflow

### Start Local Dev
```bash
# Terminal 1: PostgreSQL (if using postgres mode)
podman compose up postgres

# Terminal 2: App with hot reload (air)
make setup  # first time only
air         # watches .go files, rebuilds, restarts
```

### Run Tests
```bash
make test
# Equivalent: go test --race -parallel 8 ./...
```

**Why `--race -parallel 8`?**
- Race detector catches data races (critical for concurrent health state).
- Parallel 8: balances CPU usage on typical dev machines.

### Lint/Format
```bash
make fmt      # fix: golangci-lint --fix + yamlfmt
make go-lint  # check only
make lint     # all (go + redoc + yaml)
```

**Why golangci-lint v2?**
- Faster, native Go 1.22+ support.
- Config in `.golangci.yml` (not shown here).

## Database Operations

### Generate SQL Layer
```bash
make sqlc-generate
# Reads: sql/postgres/query.sql + schema/
# Writes: internal/generated/rdb/*_gen.go
```

### Create Migration
```bash
make migrate-new MIGRATION_NAME=add_device_index
# Creates: sql/postgres/migrations/000002_add_device_index.up/down.sql
```

### Apply Migrations (Local)
```bash
# Requires running postgres (podman compose up postgres)
make migrate-up
```

### Rollback
```bash
make migrate-down  # rolls back 1 migration
```

**Migration naming convention**: `NNNNNN_verb_noun` (e.g., `000002_add_device_idx`, `000003_alter_records_add_col`).

**Why not auto-migrate on startup?**
- Explicit control for production.
- Review SQL before apply.
- Rollback capability required.

## API Generation

```bash
cd internal/generated/oapi/rest && go generate
# Reads: rest.yaml
# Writes: rest.gen.go
```

**Add to Makefile** (TODO):
```makefile
.PHONY: oapi-generate
oapi-generate:
	cd internal/generated/oapi/rest && go generate
```

## CI/CD Pipeline (GitHub Actions implied)

Current `Makefile` targets map to CI stages:

| Stage | Command | Purpose |
|-------|---------|---------|
| `lint` | `make lint` | Static analysis, formatting |
| `test` | `make test` | Unit + race detection |
| `build` | `make build` | Compile binary (with version ldflags) |

**Redoc lint** (`make redoc-lint`): Validates OpenAPI spec via `redoc-cli` in container.

**Version injection**: 
```makefile
VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
```
Binary reports version via `main.version` variable.

## Production Deployment Considerations

### Binary
```bash
make build
# Output: ./leasepeg (static-ish, cgo for pgx)
```

### Container (Not Yet Defined)
```dockerfile
# Multi-stage: build → scratch/distroless
# Copy config.yaml, migrate binary, run migrations → start server
```

### Database Migration in Prod
- Run `migrate-up` as init container or separate job.
- App waits for DB ready (healthcheck) before serving.

### Health Checks
- Liveness: `GET /healthcheck/live` (process alive)
- Readiness: `GET /healthcheck/ready` (dependencies OK)

**K8s probes**:
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

## Unresolved / Open Questions

1. **Secrets management**: `config.yaml` has `password: ""`. Use env vars, Vault, or SealedSecrets?
2. **Log aggregation**: JSON format ready. Vector/Fluent Bit config needed.
3. **Metrics endpoint**: Not exposed. Add `/metrics` (Prometheus) in `server/rest`?
4. **Graceful shutdown timeout**: Configurable? Currently hardcoded in `cmd/server.go`.
5. **Backup/restore**: PostgreSQL backup strategy (pg_dump, PITR) not documented.
6. **Load testing**: No baseline. Add `hey` or `k6` script for `/ping` and future endpoints.