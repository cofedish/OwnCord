# Contributing to OwnCord

This repository is a Linux-first self-host fork of OwnCord. Contributions
that help harden the server, improve the Docker/Compose story, or fix
bugs discovered during local audits are very welcome.

Because the project is in early alpha, please check the issue tracker
before starting work on anything non-trivial — the direction of the fork
changes often and it's easy to waste effort on something that has
already been rewritten on a branch.

## Before you open an issue

1. Search existing issues:
   https://github.com/cofedish/OwnCord/issues
2. If you're reporting a bug, include:
   - OwnCord version / commit you're running;
   - deployment mode (bare metal / Docker / Compose);
   - host OS and relevant reverse-proxy setup;
   - reproduction steps;
   - server logs around the failure (scrub tokens and user data first).
3. If you're reporting a security issue, **do not** open a public
   issue. Use GitHub's private security advisory form instead:
   https://github.com/cofedish/OwnCord/security/advisories/new

## Pull requests

- Target branch: `main`.
- Keep each PR focused on a single topic.
- Match the existing formatting (see `.editorconfig`). In particular:
  - Go sources: tabs, `gofmt`.
  - TypeScript / JavaScript: 2 spaces, LF line endings.
- Run `make lint` and `make test` locally before submitting.
- Write commit messages in the style of the existing history
  (imperative mood, prefix with area — `server:`, `client:`, `infra:`,
  `docs:`, `fix:`, `chore:`, `refactor:`).

## Security hardening

This fork exists specifically because the upstream project had several
security issues that blocked our self-host use case. PRs that continue
that line of work — tightening defaults, shrinking attack surface,
killing unnecessary capabilities in the Docker image — are prioritized.

## License

By submitting a contribution you agree that it will be released under
the same license as the project (see [`LICENSE`](LICENSE)).
