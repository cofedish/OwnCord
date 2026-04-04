# Production Architecture

This fork targets a single production architecture. There is no separate
"legacy", "dev-mode", or SQLite-first deployment model in the long-term plan.

## Target stack

- Application server: Go
- Desktop client: Tauri v2 (`Rust + TypeScript`)
- Reverse proxy and TLS termination: Caddy
- Primary database: PostgreSQL
- Cache / coordination / rate limits / ephemeral state: Redis
- Binary object storage: S3-compatible backend (`MinIO` for self-host)
- Realtime transport: WebSocket
- Voice / video: LiveKit

## Production topology

```text
Desktop Client
    |
    | HTTPS / WSS
    v
Caddy
    |
    v
OwnCord API / WS nodes
    |             |             |
    |             |             +--> Redis
    |             |
    |             +-----------------> PostgreSQL
    |
    +--------------------------------> S3 / MinIO
```

## Data responsibilities

### PostgreSQL

Use PostgreSQL for:

- users
- sessions
- invites
- channels
- messages
- reactions
- roles and permissions
- audit log
- attachment metadata
- durable settings

### Redis

Use Redis for:

- rate limiting
- presence fanout
- WebSocket coordination across nodes
- short-lived tickets
- ephemeral caches
- job coordination

Redis is not the primary source of truth.

### S3 / MinIO

Use object storage for:

- attachments
- avatars
- exported archives
- media derivatives

Application code should persist metadata in PostgreSQL and store binary
payloads by object key in the blob store.

## Non-goals

The production target is explicitly not:

- SQLite as the primary database
- local filesystem as the primary object store
- MongoDB as an additional persistence layer
- sticky-session-only horizontal scaling

## Implementation order

1. Lock in production-only config and deployment topology.
2. Introduce database / blob-store / cache abstraction boundaries.
3. Port durable persistence to PostgreSQL.
4. Port attachment storage to S3-compatible object storage.
5. Move rate limits, presence, and coordination to Redis.
6. Rework runtime and CI for multi-service production validation.

## Release gate

The project should not be called production-ready until:

- PostgreSQL is the default and only primary database
- binary objects are served from an object-store-backed design
- Redis handles distributed ephemeral state
- CI covers integration against PostgreSQL, Redis, and MinIO
- deployment and backup/restore procedures are validated
