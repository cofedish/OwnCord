# Production Deployment Notes

## Safe baseline

- Use the Docker stack or another Linux deployment with a real reverse proxy.
- Terminate public TLS at Caddy/Nginx/Traefik with a valid certificate.
- Keep OwnCord itself on a private listener only.
- Use PostgreSQL, Redis, and S3-compatible object storage as first-class infrastructure.
- Keep `/admin` off the public internet.
- Do not use self-signed public TLS for real users.

## Required changes before exposure

- Set a real public DNS name in `.env`.
- Review `deploy/config/owncord.container.yaml`.
- Keep `tls.mode: "off"` inside the container when a reverse proxy handles TLS.
- Keep `trusted_proxies` limited to the reverse-proxy network only.
- In the provided Docker profile, `admin_allowed_cidrs` intentionally also includes the private Docker bridge so host-to-container localhost access works through NAT.

## Backups

At minimum back up:

- PostgreSQL dumps or physical volume backups from `deploy/runtime/postgres/`
- Redis only if you intentionally rely on persisted Redis state
- object data from `deploy/runtime/minio-data/`
- app-local data from `deploy/runtime/data/` while the persistence migration is still in flight
- `deploy/runtime/caddy-data/`

Example:

```bash
tar -C deploy/runtime -czf owncord-backup-$(date +%F).tar.gz data caddy-data
```

## Restoring

1. Stop the stack.
2. Restore `deploy/runtime/data`.
3. Restore `deploy/runtime/caddy-data` if you want to keep existing ACME state.
4. Start the stack again.

## Operational caveats

- The target architecture for this fork is PostgreSQL + Redis + object storage. The application migration is being moved onto that stack and should not stop at SQLite/local-files.
- The desktop client now expects valid TLS for REST/media paths. For production this is the intended behavior.
- `/admin` now uses an `HttpOnly` cookie-backed session instead of `localStorage`.
- Linux/container deployments should update via image rebuilds, not the built-in server updater.
