# CLAUDE.md

This file provides guidance to Claude Code when working with
code in this repository.

OwnCord is a self-hosted Windows chat platform with two
components: a Go server (`chatserver.exe`) and a Tauri v2
desktop client (Rust + TypeScript).

## Project Brain

This project uses an Obsidian vault at `docs/brain/`
as the single source of truth for project state.
Read and write to it during every session.

### On Session Start

1. Read `docs/brain/Dashboard.md` to get oriented
2. Read `docs/brain/02-Tasks/In Progress.md` to see what's active
3. Read `docs/brain/05-Bugs/Open Bugs.md` to see outstanding bugs
4. Create a session log at
   `docs/brain/03-Sessions/YYYY-MM-DD-summary.md`
   using `docs/brain/Templates/Session Log.md`

### On Session End

1. Update the session log with everything that was done
2. Move completed tasks from `In Progress.md` to `Done.md`
3. Update `In Progress.md` with any newly started but unfinished work
4. If any architectural decisions were made, log them in `docs/brain/04-Decisions/`

### Task Management

- Tasks live in `docs/brain/02-Tasks/` across files:
  `Backlog.md`, `In Progress.md`, `Done.md`
- Format: `- [ ] **T-XXX:** Description` (use incrementing IDs)
- When starting a task, move it from Backlog → In Progress
- When finishing, check the box and move it from
  In Progress → Done with a completion date
- New tasks discovered during work go into Backlog under the appropriate priority

### Decision Logging

- Any significant technical choice (library, arch,
  protocol, trade-off) gets a decision record
- Use template at `docs/brain/Templates/Decision.md`
- Save as `docs/brain/04-Decisions/DEC-XXX-short-title.md` (incrementing IDs)
- Statuses: `proposed` → `accepted` | `rejected` | `superseded`

### Bug Tracking

- Use template at `docs/brain/Templates/Bug Report.md`
- Save as `docs/brain/05-Bugs/BUG-XXX-short-title.md` (incrementing IDs)
- Update `docs/brain/05-Bugs/Open Bugs.md` — add to
  Active, move to Resolved when fixed
- Statuses: `open` → `investigating` → `fixed` | `wontfix`

### Requirements & Architecture

- When requirements change or are discovered, update `docs/brain/00-Overview/Requirements.md`
- When architecture evolves, update `docs/brain/01-Architecture/Design.md`
- When dependencies change, update `docs/brain/01-Architecture/Tech Stack.md`
- Always log the *reason* for changes via a decision record

### Conventions

- Use `[[wiki-links]]` for cross-references between vault files
- Use ISO dates: `YYYY-MM-DD`
- Replace `{{date}}` in templates with the actual date
- Keep files concise — prefer bullet points over prose
- Do NOT delete old session logs or decisions — they are the project history

## Codex CLI - Code REVIEW

After builds, run Codex for a second opinion:

codex exec --sandbox read-only \
"Review for bugs and logic errors"

## Reference Files (read before implementing)

All specs live in `docs/brain/06-Specs/`:

- **CHATSERVER.md** -- Master spec: phases, tasks, security
  priorities, Windows-specific details.
- **PROTOCOL.md** -- WebSocket message format. Every message
  type, payload shape, and rate limit. Server and client
  must agree on this exactly.
- **SCHEMA.md** -- SQLite table definitions, indexes, FTS5
  setup, permission bitfield definitions.
- **API.md** -- REST endpoints, request/response shapes,
  error codes. All paths start with `/api/v1/`.
- **SETUP.md** -- Tooling requirements for both server and
  client development.
- **CLIENT-ARCHITECTURE.md** -- Tauri v2 client project
  structure, component map, store design, and conventions.
- **TESTING-STRATEGY.md** -- Test infrastructure, coverage
  targets, and patterns for every test type.
- **protocol-schema.json** -- Machine-readable schema for all
  36 WebSocket message types with field definitions. Located
  at `docs/protocol-schema.json`.

## Project Structure

```text
OwnCord/
├── Server/                  # Go server (implemented)
│   ├── config/
│   ├── db/                  # + errors.go (sentinel errors)
│   ├── auth/
│   ├── api/                 # + metrics_handler.go
│   ├── ws/                  # Split: voice_join.go, voice_leave.go,
│   │                        #   voice_controls.go, voice_broadcast.go,
│   │                        #   errors.go, ringbuffer.go
│   ├── admin/static/
│   ├── migrations/
│   └── scripts/             # voice-test.sh
├── Client/
│   ├── tauri-client/        # Tauri v2 client
│   │   ├── src-tauri/       #   Rust backend
│   │   │   └── src/
│   │   ├── src/             #   TypeScript frontend
│   │   │   ├── lib/         #     Core services (incl. livekitSession.ts,
│   │   │   │                #       disposable.ts)
│   │   │   ├── stores/      #     Reactive state
│   │   │   ├── components/  #     UI components
│   │   │   ├── pages/       #     Page layouts
│   │   │   │   ├── ConnectPage.ts
│   │   │   │   ├── MainPage.ts
│   │   │   │   └── main-page/
│   │   │   │       ├── ChannelController.ts
│   │   │   │       ├── ChatArea.ts
│   │   │   │       ├── SidebarArea.ts
│   │   │   │       └── ...
│   │   │   └── styles/      #     CSS (from mockups)
│   │   └── tests/
│   │       ├── unit/
│   │       ├── integration/
│   │       └── e2e/
│   └── ui-mockup.html      # Design source of truth
└── docs/
```

## Build Commands

### Server (Go)

```bash
cd Server
go build -o chatserver.exe -ldflags "-s -w -X main.version=1.0.0" .
go test ./...                        # all tests
go test ./... -cover                 # with coverage
```

### Client (Tauri v2)

```bash
cd Client/tauri-client

# Development (hot reload)
npm run tauri dev

# Build release
npm run tauri build

# Run tests
npm test                             # all tests (vitest)
npm run test:unit                    # unit tests only
npm run test:integration             # integration tests
npm run test:e2e                     # Playwright E2E tests
npm run test:coverage                # with coverage report
```

### Dev Tools

```bash
# Server
go install github.com/air-verse/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Client (installed via npm)
# vitest, playwright, typescript, vite — all in package.json
```

## Branch Strategy

- `main` -- stable releases
- `dev` -- active development

## Key Features

- **Voice & video chat**: LiveKit-powered voice and video.
  LiveKit server runs as a companion process alongside
  `chatserver.exe`. Client connects via `livekitSession.ts`.
  VideoGrid component replaces chat area when cameras are
  active.
- **GIF picker**: Tenor API v2 integration via
  `lib/tenor.ts`. Uses Google's public anonymous API key.
  Picker in MessageInput sends GIF URL as message content;
  inline image rendering in `renderers.ts`.
- **Push-to-talk**: Rust-side `GetAsyncKeyState` polling
  (`src-tauri/src/ptt.rs`) — non-consuming, works globally.
  Key capture UI in KeybindsTab with 10s timeout.
  Client-side wiring in `lib/ptt.ts`.
- **Desktop notifications**: `lib/notifications.ts` — Tauri
  plugin-notification with Web Notification API fallback.
  Taskbar flash, notification sound, @everyone suppression.
- **Compact mode**: CSS class `.compact-mode` on body
  reduces spacing, avatar sizes, and font sizes throughout.
- **Admin IP restriction**: `/admin` routes restricted to
  `admin_allowed_cidrs` in server config (default: private
  networks only). Middleware in `api/middleware.go`.
- **Metrics endpoint**: `GET /api/v1/metrics` (admin IP
  restricted) returns uptime, goroutines, heap, connected
  users.
- **Reconnection with state recovery**: Client tracks `seq`
  numbers on all server broadcasts. On reconnect, sends
  `last_seq` in auth; server replays missed events from a
  1000-event ring buffer. Falls back to full `ready` if too
  far behind.
- **Heartbeat monitoring**: Server sweeps for stale
  connections every 30s, kicks clients with no activity for
  90s.

## Critical Rules (always apply)

- **API paths**: Always `/api/v1/*` (matches server router)
- **Roles**: Always use role NAME strings ("admin",
  "member"), never numeric role\_id in UI-facing code
- **Rate limiting**: Client must respect PROTOCOL.md
  limits (typing 1/3s, presence 1/10s, voice 20/s)
- **Status values**: Only `online`, `idle`, `dnd`,
  `offline`. Never `invisible`.
- **Tenor API key**: The key in `lib/tenor.ts` is Google's
  public anonymous key — not a secret. Do not move to env.

## Conventions & Details (see canonical files in docs/brain/)

- **Client architecture & conventions**:
  06-Specs/CLIENT-ARCHITECTURE.md
- **Server spec & conventions**: 06-Specs/CHATSERVER.md
- **Security rules**: 06-Specs/CHATSERVER.md (Security section)
- **Testing requirements**: 06-Specs/TESTING-STRATEGY.md
- **Coverage target**: 80%+ (TDD: RED → GREEN → IMPROVE)

## gstack Skills

gstack is installed at `~/.claude/skills/gstack`.

- **Web browsing**: Always use `/browse` from gstack for
  all web browsing. Never use `mcp__claude-in-chrome__*`
  tools.

Available skills:

- `/plan-ceo-review` — CEO-level plan review
- `/plan-eng-review` — Engineering plan review
- `/review` — Code review
- `/ship` — Ship checklist
- `/browse` — Headless browser for QA and browsing
- `/qa` — QA testing
- `/qa-only` — QA testing (no fixes)
- `/setup-browser-cookies` — Configure browser cookies
- `/retro` — Retrospective
- `/document-release` — Document a release
