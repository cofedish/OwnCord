# OwnCord Fork

Linux-first, self-host-oriented fork of OwnCord with Docker/Compose deployment, safer internet-facing defaults, and a first pass of security hardening.

## What this fork is

This repository keeps the OwnCord codebase, but changes the operational focus:

- the server is intended to run on Linux
- Docker/Compose is the recommended deployment path
- public exposure is expected to happen behind a reverse proxy such as Caddy or Nginx
- several security issues found during local audit were fixed
- the desktop client is still a separate Tauri application and remains Windows-oriented for now

## Current status

This project is still early-stage software. The fork is safer and easier to self-host than the upstream state I audited, but it is still not something I would classify as mature secure messaging infrastructure.

Use it for:

- lab environments
- hobby self-hosting
- small trusted groups after proper TLS and reverse-proxy setup

Do not use it for:

- sensitive communications
- hostile multi-tenant deployments
- high-scale production workloads

## What changed in this fork

- Dockerfile and Compose deployment for the server
- Linux-friendly deployment docs and production notes
- localhost/VPN-oriented admin exposure model
- admin UI moved away from `localStorage` token persistence to `HttpOnly` cookie-backed session flow
- private attachment responses no longer advertise public caching
- desktop client no longer bypasses TLS validation for REST/media requests
- production devtools removed from the default client build
- safer Tauri defaults and reduced risky production surface
- Linux/container deployments are expected to update via image rebuild/redeploy instead of the built-in server updater

## Quick start

### Server

1. Clone the repository.
2. Copy `.env.example` to `.env`.
3. Set `OWNCORD_DOMAIN`.
4. Review `deploy/config/owncord.container.yaml`.
5. Start the stack:

```bash
docker compose --profile proxy up -d --build
```

6. Reach the admin panel through a local tunnel:

```bash
ssh -L 8080:127.0.0.1:8080 your-server
```

7. Open `http://127.0.0.1:8080/admin`, create the owner account, then generate invite codes for users.

### Client

The desktop client is installed separately from the server.

To build it from source:

```bash
cd Client/tauri-client
npm install
npm run tauri build
```

The Windows installer is produced in:

`Client/tauri-client/src-tauri/target/release/bundle/nsis/`

Important:

- this fork no longer allows insecure REST/media TLS bypass in the client
- the client should connect to a real HTTPS endpoint behind your reverse proxy
- the admin-local `127.0.0.1:8080` listener is not the normal user-facing endpoint

## Security posture

The fork includes practical fixes for the issues found during audit, but “improved” does not mean “fully trusted”.

Main hardening areas already addressed:

- desktop TLS trust model
- admin session handling
- attachment caching/privacy
- Linux deployment path
- reverse-proxy/admin exposure defaults

Main limitations that still remain:

- SQLite is still the database and still a scaling bottleneck
- the project is still alpha-grade software
- the desktop client still needs more runtime hardening and broader end-to-end validation
- voice/LiveKit remains an operationally sensitive component

## Deployment model

Recommended production-ish shape:

- OwnCord server in Docker
- Caddy or Nginx handling public TLS
- `/admin` blocked publicly and reached only via localhost/VPN/SSH tunnel
- regular file-level backups of `deploy/runtime/data`
- image rebuild/redeploy for updates

See also:

- [README.md](README.md)
- [README.ru.md](README.ru.md)
- [docs/docker-deployment.md](docs/docker-deployment.md)
- [docs/production-deployment.md](docs/production-deployment.md)

## Repository layout

- `Server/` — Go server
- `Client/tauri-client/` — Tauri desktop client
- `deploy/` — Compose, Caddy, config, systemd examples
- `docs/` — deployment and project docs

## License

AGPL-3.0
