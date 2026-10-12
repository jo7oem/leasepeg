# API Design: Decisions and Evolution Path

## Current Endpoints (OpenAPI 3.0)

Defined in `internal/generated/oapi/rest/rest.yaml`, generated server in `internal/generated/oapi/rest/rest.gen.go`.

| Method | Path | Handler | Auth | Purpose |
|--------|------|---------|------|---------|
| `GET` | `/healthcheck/live` | `GetHealthcheckLive` | ❌ | Process alive (k8s livenessProbe) |
| `GET` | `/healthcheck/ready` | `GetHealthcheckReady` | ❌ | Ready to serve (k8s readinessProbe) |
| `GET` | `/ping` | `GetPing` | ❌ | Echo client IP (connectivity test) |

**Response schemas**: `healthcheckLiveResp`, `healthcheckReadyResp`, `Pong` — all in `rest.yaml` components.

## Design Decisions

### 1. Strict Server Mode + Generated Handlers

```go
oapi.NewStrictHandler(handler{}, middlewares...)
```

**Why?** 
- Compile-time guarantee: handler implements all operations.
- No `interface{}` — typed request/response objects.
- Middleware chain (`StrictMiddlewareFunc`) for cross-cutting concerns.

### 2. BasicAuth Security Scheme (Defined, Not Enforced)

```yaml
security:
  - basicAuth: []
components:
  securitySchemes:
    basicAuth:
      type: http
      scheme: basic
```

**Why defined but empty per-operation?**
- Placeholder for future auth.
- OpenAPI spec documents intent.
- Enforcing now would break `/healthcheck/*` and `/ping` (infra probes).
- **Decision**: Apply per-operation when real endpoints added.

### 3. IP Extraction Middleware (`strictIPRecordHandler`)

Extracts source IP from `r.RemoteAddr` → context → `GetPing` handler.

**Why middleware not handler logic?**
- Separation: transport concern (IP) vs business logic.
- Reusable for future endpoints needing client IP.
- Context-based — testable without HTTP request.

### 4. Health Check State Machine

`store/memory` defines states: `Starting → Ready ↔ Degraded → ShuttingDown`, `Error` terminal.

**Why custom state machine not library?**
- Minimal dependency.
- Explicit transitions visible in code.
- `Degraded` distinct from `Error` — partial functionality (e.g., DB slow but cache works).

**Response codes**:
- `200`: `Ready` or `Degraded`
- `503`: `Starting` or `ShuttingDown`
- `500`: `Error`

Matches k8s semantics: 503 = don't send traffic, 500 = something broken.

## Evolution Path

### Phase 1: Device Management API (Planned)
```
POST   /devices              CreateDevice
GET    /devices              ListDevices (user-scoped)
GET    /devices/{id}         GetDevice
PATCH  /devices/{id}         UpdateDevice
DELETE /devices/{id}         DeleteDevice
GET    /devices/{id}/records ListRecords
POST   /devices/{id}/records CreateRecord
PATCH  /records/{id}         UpdateRecord
DELETE /records/{id}         DeleteRecord
```

**Auth**: Bearer token (JWT) or API key. BasicAuth removed.

### Phase 2: User/Auth API
```
POST   /auth/register        Local registration
POST   /auth/login           Local login → JWT
POST   /auth/refresh         Token refresh
GET    /auth/providers       List linked OAuth
POST   /auth/link/{provider} Initiate OAuth link
DELETE /auth/unlink/{provider} Unlink OAuth
```

### Versioning Strategy

**URL path versioning**: `/v1/...` from Phase 1.
- OpenAPI `servers[0].url` updated per version.
- Generated package: `oapi/v1/rest` (separate go package).
- Old version maintained until deprecation.

**Why not header versioning?**
- Explicit in URL — easier debugging, caching, routing.
- Standard practice (Stripe, GitHub, AWS).

### Error Response Format (To Be Standardized)

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "human-readable",
    "details": [{"field": "name", "issue": "required"}]
  }
}
```

**Why not RFC 7807 (Problem Details)?** 
- Consider adopting. More interoperable.
- Decision pending — current endpoints don't return errors yet.

## Unresolved / Open Questions

1. **Auth mechanism**: JWT (stateless) vs session (stateful)? JWT fits horizontal scaling but needs revocation strategy (blocklist or short expiry + refresh).
2. **Rate limiting**: Where? Middleware (per-IP) or auth-aware (per-user)? Start with per-IP in middleware.
3. **Request validation**: oapi-codegen validates path/query params. Body validation? Use `go-playground/validator` tags on generated types or custom middleware.
4. **Pagination**: Cursor-based (stable) vs offset? Cursor preferred for `ListDevices`/`ListRecords`.
5. **Webhooks/async**: Device IP changes — push vs poll? Current `/ping` is poll. Webhook subscription API may be needed.
6. **API key vs OAuth for devices**: Devices may need programmatic access. API key (prefix + hash) simpler than OAuth client credentials.