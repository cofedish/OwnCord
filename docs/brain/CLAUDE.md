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
- **E2E-BEST-PRACTICES.md** -- E2E test patterns: persistent
  fixtures, login-once, selector/assertion best practices.
- **DM-SYSTEM.md** -- Direct message system architecture,
  server/client flows, authorization model.
- **THEME-SYSTEM.md** -- Theming system: built-in themes,
  custom JSON import/export, accent color, CSS injection
  prevention.
- **RECONNECTION.md** -- Reconnection protocol: seq numbers,
  ring buffer replay, state recovery flow.
- **VOICE-COMPARISON-MATRIX.md** -- Competitive comparison
  of 25 voice/video behaviors across Discord, TeamSpeak,
  Guilded. Gap priorities and OwnCord target behaviors.
- **protocol-schema.json** -- Machine-readable schema for all
  36 WebSocket message types with field definitions. Located
  at `docs/protocol-schema.json`.

## Repo Copilot Assets

Project-shared Copilot customizations live in `.github/instructions/`
and `.claude/skills/`. Use them as the first stop before inventing new
workflow guidance.

### Auto-Applied Instructions

- `admin-panel.instructions.md` -- Admin handlers, SSE log streaming,
  backups, updates, and admin-only workflows.
- `database-and-schema.instructions.md` -- SQLite schema, migrations,
  tracked `schema_versions`, and DB-layer sentinel-error rules.
- `design-system.instructions.md` -- Client design-system, theming,
  accessibility, and token usage rules.
- `dm-authorization.instructions.md` -- DM-specific authorization branch
  rule for channel-aware handlers and dispatcher flows.
- `go-server.instructions.md` -- Go server structure, safety rules,
  permissions, and testing expectations.
- `protocol-and-api.instructions.md` -- REST/WS synchronization,
  reconnection, rate limits, and protocol coupling.
- `tauri-frontend.instructions.md` -- TypeScript client/store patterns
  and existing networking abstractions.
- `tauri-rust.instructions.md` -- Tauri/Rust backend layout, IPC,
  TOFU, proxying, updater, and Windows boundary rules.
- `testing.instructions.md` -- TDD, coverage, regression tests, and
  validation command expectations.
- `vault-workflow.instructions.md` -- Project-brain workflow, session
  logs, task/bug/decision tracking, and required doc updates.

### On-Demand Repo Skills

- `owncord-patterns` -- Broad repo architecture and workflow patterns.
- `protocol-guard` -- Protocol/API changes, new message types, and
  contract drift prevention.
- `code-review` -- Repo-specific code review workflow for server-side
  and safety-sensitive changes.
- `security-audit` -- Auth, tokens, permissions, uploads, and general
  security review.
- `test-writer` -- Repo-specific Go test generation and regression
  coverage guidance.
- `owncord-reconnection` -- `last_seq`, replay buffer, state recovery,
  and reconnect bug work.
- `livekit-integration` -- LiveKit token flow, media lifecycle,
  device-switching, cleanup, and voice recovery work.
- `e2e-tier-selector` -- Choose unit vs integration vs mocked/prod/native
  E2E for client and Tauri changes.
- `config-and-secrets` -- Config loading, env overrides, TLS/LiveKit
  setup, and secret-safe configuration changes.
- `observability-debugging` -- Metrics, diagnostics, admin log stream,
  and persisted client log workflows.
- `windows-integration` -- Windows Credential Manager, TOFU/WebView2,
  PTT, tray, hotkeys, updater, and other desktop-bound paths.

## Project Structure

```text
OwnCord/
├── Server/                  # Go server (implemented)
│   ├── config/
│   ├── db/                  # + errors.go (sentinel errors)
│   ├── auth/
│   ├── api/                 # + totp_handler.go, metrics_handler.go, diagnostics_handler.go
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
│   │   │   │                #       audioPipeline.ts, audioElements.ts,
│   │   │   │                #       deviceManager.ts, connectionStats.ts,
│   │   │   │                #       disposable.ts, logPersistence.ts)
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

# Mutation testing
go-gremlins mutate ./...             # run gremlins mutation testing

# Load testing (WebSocket)
k6 run scripts/k6/ws-load.js         # WebSocket load test

# Chaos/network fault testing
bash scripts/toxiproxy/chaos-test.sh # toxiproxy fault injection
```

### Client (Tauri v2)

```bash
cd Client/tauri-client

# Install dependencies (first time only)
npm install

# Development (hot reload)
npm run tauri dev

# Build release
npm run tauri build

# Run tests
npm test                             # all tests (vitest)
npm run test:unit                    # unit tests only
npm run test:integration             # integration tests
npm run test:e2e                     # Playwright E2E (mocked Tauri)
npm run test:e2e:native              # Playwright E2E (real Tauri exe + CDP)
npm run test:e2e:prod                # Playwright E2E (prod build)
npm run test:e2e:ui                  # Playwright UI mode
npm run test:watch                   # watch mode for tests
npm run test:coverage                # coverage report
npm run test:mutate                  # Stryker mutation testing
npm run test:mutate:dry              # Stryker dry-run only

# Type checking
npm run typecheck                    # full typecheck (all sources)
npm run typecheck:build              # typecheck build-only tsconfig

# Linting
npm run lint                         # ESLint check (src/)
npm run lint:fix                     # ESLint auto-fix
```

### Dev Tools

```bash
# Server
go install github.com/air-verse/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/go-gremlins/gremlins/cmd/go-gremlins@latest

# Client (installed via npm via package.json)
# vitest, playwright, typescript, vite, stryker — all in package.json

# Optional global tools
# k6 — load testing (https://k6.io/docs/get-started/installation/)
# toxiproxy — chaos testing (https://github.com/shopify/toxiproxy)
```

## Branch Strategy

- `main` -- stable releases
- `dev` -- active development

## Key Features

- **Voice & video chat**: LiveKit-powered voice and video.
  LiveKit server runs as a companion process alongside
  `chatserver.exe`. Client connects via `livekitSession.ts`
  (facade pattern delegating to `audioPipeline.ts`,
  `audioElements.ts`, `deviceManager.ts`). VideoGrid
  component replaces chat area when cameras are active
  with Discord-style fixed 16:9 aspect ratio tiles.
  Sidebar stream preview: hover voice channel to see
  live video thumbnail of active screenshares.
  VAD runs on AudioWorklet (`public/vad-worklet.js`) with
  setTimeout fallback. Token TTL: 24h (refresh at 23h).
  Ghost state cleanup: retry with exponential backoff.
  Speaker indicators: pulsing green glow animation.
  Device hot-swap: auto-fallback to default + toast.
  Listen-only mode: "Grant Microphone" recovery button.
  Connection quality: auto-expand stats on degradation.
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
- **Connection quality indicator**: Signal-bars icon +
  ping text in VoiceWidget header. Clicking it expands a
  transport statistics pane (outgoing/incoming rates with
  Mbps, packets, RTT, session totals). Polls WebRTC stats
  every 2s via `lib/connectionStats.ts`. Color-coded: green
  (<100ms), yellow (100-200ms), red (>200ms). Auto-expands
  stats pane on poor/bad quality (3s debounce). Quality
  change callback: `onQualityChanged` with debounce.
- **Compact mode**: CSS class `.compact-mode` on body
  reduces spacing, avatar sizes, and font sizes throughout.
- **Admin IP restriction**: `/admin` routes restricted to
  `admin_allowed_cidrs` in server config (default: private
  networks only). Middleware in `api/middleware.go`.
- **Metrics endpoint**: `GET /api/v1/metrics` (admin IP
  restricted) returns uptime, goroutines, heap, connected
  users, voice sessions, LiveKit health.
- **Reconnection with state recovery**: Client tracks `seq`
  numbers on all server broadcasts. On reconnect, sends
  `last_seq` in auth; server replays missed events from a
  1000-event ring buffer. Falls back to full `ready` if too
  far behind.
- **Heartbeat monitoring**: Server sweeps for stale
  connections every 30s, kicks clients with no activity for
  90s.
- **Direct Messages (1-on-1)**: Full-stack DM implementation.
  Server: `db/dm_queries.go` with GetOrCreateDMChannel,
  GetUserDMChannels, OpenDM, CloseDM. REST: `POST /api/v1/dms`,
  `GET /api/v1/dms`, `DELETE /api/v1/dms/{channelId}`.
  WebSocket DM events: `dm_channel_open`, `dm_channel_close`.
  Client: `dm.store.ts` for state, dispatcher handles DM events,
  "+" button opens member picker to start DMs. DM header shows
  `@ username` with live status. Auto-reopen on message.
- **Unified Sidebar Layout**: Single 240px sidebar. Sections:
  Server header → DM Preview (top 3, bubbles on new message) →
  Text Channels → Voice Channels → Members (collapsible,
  persisted state) → Voice Widget → User Bar.
  DM header shows total unread badge. "View all messages"
  link switches to full DM mode. Members section collapses
  to header-only bar. Architecture in
  `pages/main-page/SidebarArea.ts`.
- **Quick-Switch Server Overlay**: 🚪 button in UserBar opens
  overlay showing favorited servers. Click to disconnect and
  switch via sessionStorage handoff to ConnectPage.
  Component: `QuickSwitchOverlay.ts`.
- **Voice call duration timer**: Elapsed time counter in
  VoiceWidget (MM:SS / HH:MM:SS). Stored as `joinedAt` in
  `voice.store.ts`, rendered by 1-second interval in
  `VoiceWidget.ts`. Local-only (each user sees their own
  timer). Resets on leave/disconnect.
- **Auto-login**: One server profile can be set as auto-login
  (lightning bolt toggle on server cards). On startup, loads
  saved credentials from Windows Credential Manager and
  attempts automatic login. Falls back to login form on
  failure or 2FA. Only one profile can be auto-login at a
  time. Enforces `rememberPassword: true`.
- **Server health with online users**: `GET /api/v1/health`
  returns `online_users` count. Server cards on ConnectPage
  show latency + online user count. Health checks repeat
  every 15 seconds so offline servers update automatically.
- **OC Neon Glow Theme + Theming System**: New default theme
  with cyan (#00c8ff) → purple (#7b2fff) gradient. Theme manager
  in `lib/themes.ts` supports built-in + custom themes via
  JSON import/export. Accent color picker overrides theme accent.
  Theme file: `theme-neon-glow.css`. Restored on app startup
  (both theme class and accent color override).
- **Discord-style Settings Panel**: Centered floating panel with
  blurred backdrop (8px blur), rounded corners (12px), scale
  animation on open, and click-outside-to-close behavior.
  Component: `components/SettingsOverlay.ts`. CSS: `.settings-overlay`
  (full-screen backdrop) and `.settings-panel` (900px wide card with
  sidebar navigation, 80vh height). Tabs: Account, Appearance,
  Notifications, Text & Images, Accessibility, Voice & Audio,
  Keybinds, Advanced, Logs.
- **Two-factor authentication (TOTP)**: 2FA protection with TOTP
  enrollment/disable in Settings > Account. Server: 10-minute auth
  challenges with 5-attempt rate limiting per challenge, 10 req/min
  per IP. Client: QR code display, backup code backup, password
  confirmation. Server-wide `require_2fa` policy enforcement. Login
  flow returns `requires_2fa: true` + `partial_token` (10min TTL) on
  second-factor challenge; client shows TOTP overlay. Partial tokens
  expire on successful 2FA completion and on logout.
- **Account deletion**: `DELETE /api/v1/auth/account` with
  password confirmation. Anonymizes user (username →
  `[deleted-{id}]`, clears password/avatar/TOTP, bans row).
  Soft-deletes messages, removes sessions/DM participation/
  reactions/read states. Blocks last-admin deletion. Client:
  "Danger Zone" section in AccountTab with inline confirmation.
- **Observability & debugging**: Structured logging across all
  layers. Server: enriched HTTP request logs (client_ip, req_id,
  bytes), WS disconnect stats (duration, msgs sent/received/dropped,
  voice channel), LiveKit webhook event logging. Client: JSONL log
  persistence to disk with 5-day rotation (`lib/logPersistence.ts`),
  LiveKit ICE candidate logging (host/srflx/relay types), room
  lifecycle events, WS reconnection tracking. Rust proxies:
  structured `log` crate logging for TLS handshakes, TOFU checks,
  connection lifecycle. Cache clear buttons in Settings > Advanced.
  Diagnostics endpoint: `GET /api/v1/diagnostics/connectivity`.

## Design System
Always read DESIGN.md before making any visual or UI decisions.
All font choices, colors, spacing, and aesthetic direction are defined there.
Do not deviate without explicit user approval.
In QA mode, flag any code that doesn't match DESIGN.md.

## Critical Rules (always apply)

- **API paths**: Always `/api/v1/*` (matches server router)
- **Roles**: Always use role NAME strings ("admin",
  "member"), never numeric role\_id in UI-facing code
- **Rate limiting**: Client must respect PROTOCOL.md
  limits (typing 1/3s, presence 1/10s, voice 20/s)
- **Status values**: Only `online`, `idle`, `dnd`,
  `offline`. Never `invisible`. Server validates against
  this allowlist in `handlers_presence.go`.
- **Username rules**: 2-32 runes, printable only, no
  control chars. Enforced by `auth.ValidateUsername()` in
  both registration and admin setup.
- **Backup path security**: All backup handlers use
  `filepath.Abs` + prefix check (not just string-contains)
  to prevent path traversal on Windows.
- **SSE auth**: Admin log stream uses single-use tickets
  (30s TTL) via `POST /admin/api/logs/ticket`, NOT tokens
  in URL query parameters.
- **Third-party fetches**: `acceptInvalidCerts: true` is
  ONLY for server URLs. Third-party fetches (OG previews,
  external images) must use `acceptInvalidCerts: false`.
- **Content-Type validation**: Data URIs from remote
  fetches must validate MIME type against `SAFE_MIME_TYPES`
  allowlist before construction.
- **Tenor API key**: The key in `lib/tenor.ts` is Google's
  public anonymous key — not a secret. Do not move to env.
- **DM authorization**: DM channels use `IsDMParticipant`
  checks instead of role-based permissions. Every handler
  that touches a channel must branch on `ch.Type == "dm"`
  and verify participant membership. This applies to WS
  handlers (channel_focus, typing, chat_send, edit, delete,
  reaction) and REST handlers (GET messages, pins).

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
