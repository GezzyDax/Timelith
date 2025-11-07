# Repository Guidelines

## Project Structure & Module Organization
- `go-backend/cmd/server` hosts the API entrypoint, while `internal/{api,auth,scheduler,telegram,...}` contain focused packages for request handling, scheduling, and Telegram orchestration. Keep business logic inside `internal` to preserve boundaries.
- `web-ui/src/app` provides the Next.js routing layer; shared UI, hooks, and schema utilities live under `src/{components,lib,types}`. Favor colocated feature folders when extending dashboard screens.
- `scripts/` holds helper automation, `docker-compose.yml` wires Postgres/Redis/backend/frontend, and `uploads/` is mounted for media payloads.

## Build, Test, and Development Commands
```bash
# Backend
cd go-backend && make install        # sync Go modules
make run                             # start API locally
make build                           # compile binary to bin/server

# Frontend
cd web-ui && npm install             # install dependencies
npm run dev                          # launch Next.js dev server
docker compose up -d                 # full stack: db, cache, API, UI
```
Use `docker compose logs go-backend -f` or `npm run start` for production-like verification.

## Coding Style & Naming Conventions
- Backend code is formatted with `gofmt` (tabs for indentation, camelCase identifiers). Keep handlers thin and push orchestration into `internal/*` services. Run `go fmt ./...` before committing.
- Frontend follows Next.js + TypeScript defaults: 2-space indentation, PascalCase for components, camelCase for hooks/utilities. `npm run lint` enforces ESLint + Next rules. Tailwind utility classes should be ordered logically (layout → spacing → color).

## Testing Guidelines
- Go tests live beside source files as `_test.go`; cover scheduler logic, Telegram adapters, and database services with table-driven tests. Execute `cd go-backend && make test`.
- The web UI currently lacks automated tests; when adding them prefer React Testing Library under `web-ui/src/__tests__` and mirror file names (e.g., `schedule-form.test.tsx`). Document any manual verification steps in pull requests until coverage is added.

## Commit & Pull Request Guidelines
- Follow the existing short, present-tense summaries (`git log` shows examples like “Initial implementation…”). Reference issue IDs when applicable (`feat: add cron preview (#42)`).
- Pull requests should describe user-facing impact, mention affected services (`go-backend`, `web-ui`, infra), list test commands executed, and include screenshots or curl snippets for dashboard/UI changes.

## Security & Configuration Tips
- Never commit `.env` secrets; rely on `.env.example` for placeholders. Ensure `ENCRYPTION_KEY`, `JWT_SECRET`, and Telegram credentials are injected via environment variables or Docker secrets.
- When debugging locally, avoid logging full access tokens or verification codes—scrub sensitive fields before printing and prefer the shared `logger` package for structured output.
