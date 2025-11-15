# Repository Guidelines

## Project Structure & Module Organization
- `backend/` (Go under `cmd/`, `internal/<domínio>`, shared `pkg/`), `frontend/` and `backoffice/` (React + Vite + TS) form the core apps.
- Knowledge base lives in `docs/`; infra assets stay in `docker/` plus Compose files at the root; helper automation belongs in `scripts/`.
- `tests/` nests Dredd (`tests/dredd/dredd.yml`), E2E and Postman suites; commit generated API types per app in `src/types/` and keep UI assets scoped to each `public/`.

## Build, Test, and Development Commands
- `make api-lint | api-preview | api-types` lint, preview and emit TypeScript definitions from `docs/api.yaml`.
- `make backend-run/test/lint/build` wraps `go run`, `go test ./...`, `golangci-lint` and binary production.
- `npm --prefix frontend run dev` (same for `backoffice`) launches Vite with automatic type generation; `npm --prefix <app> run test -- --coverage` executes Vitest suites.
- `make api-contract-test` migrates and runs Dredd, while `docker compose up --build` spins up PostgreSQL + backend + UIs for manual QA.

## Coding Style & Naming Conventions
- Format Go code with `gofmt`/`goimports`, enforce `golangci-lint`, use `camelCase` locals, `PascalCase` exports, and suffixes `Service`, `Repository`, `Request`, `Response`.
- Frontend code must satisfy ESLint + Prettier + strict TypeScript; components use PascalCase, hooks `useSomething`, shared logic in `src/lib` or `src/hooks`.
- Keep all docs in Portuguese, respect `.editorconfig`, prefer CSS Modules or Tailwind colocated with components, and never mutate React Query caches directly.

## Testing Guidelines
- Backend suites live in `*_test.go` tables using `stretchr/testify`; target ≥70 % coverage and execute `make backend-test` before any PR.
- React tests use Vitest + Testing Library with names like `shouldDisplayAgendaWhen...`; aim for ≥50 % coverage confirmed via `make coverage-frontend` / `coverage-backoffice`.
- Re-run `make api-contract-test` for every change touching handlers or `docs/api.yaml`, and keep UI regression checks inside `tests/e2e/`.

## Commit & Pull Request Guidelines
- Follow Conventional Commits (`feat: agenda sync`, `fix: auth middleware`) and branches `feature/<slug>`, `fix/<slug>`, `chore/<slug>`.
- PR descriptions should link issues, cite executed commands (backend/frontend tests, contract, build) and include screenshots or curl snippets for UX/API tweaks.
- Secure two approvals (or one senior reviewer for critical work), wait for green CI (lint, build, coverage, API Spec Quality), and squash noisy WIP before merging.
