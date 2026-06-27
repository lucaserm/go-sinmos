# SINMOS — Frontend

SPA for the SINMOS school monitoring/security system, talking to the Go API in
[`../backend`](../backend).

**Stack:** Bun · Vite · React 19 · React Router · TanStack Query · React Hook Form + Zod · Tailwind v4.

## Development

```bash
bun install
bun run dev      # http://localhost:5173
```

Vite proxies `/api` → `http://localhost:8080`, so the dev app is same-origin
with the API (the Go server has no CORS). Start the backend first:

```bash
cd ../backend
docker compose up -d        # postgres
goose up                    # run migrations (reads ../backend/.env)
go run ./cmd                # API on :8080
```

Register the first user (no seed exists):

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"name":"Coord","code":"coord1","password":"secret6","role":"ADMIN"}'
```

Then log in at `/login` with `coord1` / `secret6`.

## Scripts

- `bun run dev` — dev server
- `bun run build` — type-check + production build
- `bun run lint` — oxlint
- `bun run preview` — serve the production build

## Configuration

`.env` → `VITE_API_BASE` (default `/api/v1`). For a non-localhost deploy, either
serve `dist/` from the Go server or add CORS to the API.

## Structure

```
src/
  auth/         AuthProvider, useAuth, route guards
  components/   UI kit, Layout, Table, Modal, Toast, QueryState
  lib/          axios client + token refresh, resource factory, formatters
  pages/        Landing, Login, Dashboard
    students/   list + detail
    occurrences/ list + detail + register form (disciplinary workflow)
    admin/      registration hub (config-driven CRUD via registry.ts)
  types/        API DTO types
```

## Roles

- **ADMIN** (Coordenação) — everything, incl. cadastros and warning approval.
- **SUPPORT** (Assistência) — students + register/track occurrences.
- **RECEPTION** (Portaria) — student lookup only.

> UI role-gating is cosmetic. Most API endpoints are currently unauthenticated
> server-side — enforce auth on the backend before relying on it.

## Known backend follow-ups (not frontend bugs)

- No CORS on the API (dev relies on the Vite proxy).
- Only `GET /auth/me` and `POST /occurrences` enforce auth.
- No per-student read for the guardian/subject join tables, so the student
  detail page can't list horários/responsáveis yet.
- No image upload — students take a `photo_url` string only.
- No server-side search; list pages fetch a wide page and filter client-side.
- Schedule times are parsed as `HH:MM` (Go layout `15:04`), not `HH:MM:SS`.
