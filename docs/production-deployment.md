# Production Deployment Notes

## Safe baseline

- Use the Docker stack or another Linux deployment with a real reverse proxy.
- Terminate public TLS at Caddy/Nginx/Traefik with a valid certificate.
- Keep OwnCord itself on a private listener only.
- Keep `/admin` off the public internet.
- Do not use self-signed public TLS for real users.

## Required changes before exposure

- Set a real public DNS name in `.env`.
- Review `deploy/config/owncord.container.yaml`.
- Keep `tls.mode: "off"` inside the container when a reverse proxy handles TLS.
- Keep `trusted_proxies` limited to the reverse-proxy network only.
- Leave `admin_allowed_cidrs` on localhost unless you have a dedicated VPN/admin network.

## Backups

At minimum back up:

- `deploy/runtime/data/chatserver.db`
- `deploy/runtime/data/uploads/`
- `deploy/runtime/data/backups/`
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

- SQLite is still a single-writer database. This is fine for hobby/small-home use, but it is not a strong scaling story.
- The desktop client now expects valid TLS for REST/media paths. For production this is the intended behavior.
- `/admin` now uses an `HttpOnly` cookie-backed session instead of `localStorage`.
- Linux/container deployments should update via image rebuilds, not the built-in server updater.

