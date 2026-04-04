# Docker Deployment

This is the recommended Linux self-host path for OwnCord.

## What this stack does

- Runs the Go server in a non-root Linux container.
- Keeps OwnCord itself on plain HTTP inside the private Docker network.
- Terminates public TLS at Caddy.
- Blocks `/admin` at the public reverse proxy by default.
- Exposes the raw server only on `127.0.0.1:${OWNCORD_ADMIN_PORT}` for SSH tunneling or local admin access.

## Files

- `Dockerfile`
- `compose.yaml`
- `.env.example`
- `deploy/config/owncord.container.yaml`
- `deploy/Caddyfile`
- `deploy/systemd/owncord-compose.service`

## Quick start

1. Copy `.env.example` to `.env`.
2. Set `OWNCORD_DOMAIN` to the public DNS name you control.
3. Review `deploy/config/owncord.container.yaml`.
4. Start the stack:

```bash
docker compose --profile proxy up -d --build
```

5. Check status:

```bash
docker compose ps
docker compose logs -f --tail=200
```

## First admin access

`/admin` is intentionally blocked on the public reverse proxy. Use one of these:

1. SSH tunnel:

```bash
ssh -L 8080:127.0.0.1:8080 your-server
```

Then open `http://127.0.0.1:8080/admin`.

2. Local browser on the server host:

`http://127.0.0.1:${OWNCORD_ADMIN_PORT}/admin`

## Data locations

- Database: `deploy/runtime/data/chatserver.db`
- Uploads: `deploy/runtime/data/uploads/`
- Backups: `deploy/runtime/data/backups/`
- ACME cache: `deploy/runtime/caddy-data/`

## Updating

Rebuild and restart:

```bash
docker compose --profile proxy up -d --build
```

OwnCord's built-in server updater is intentionally unavailable on non-Windows hosts.

## Voice / LiveKit

The compose stack does not auto-deploy LiveKit.

Use an external LiveKit service and set:

- `voice.livekit_url`
- `voice.livekit_api_key`
- `voice.livekit_api_secret`

Do not enable `voice.livekit_binary` inside the container unless you fully trust that deployment path.

