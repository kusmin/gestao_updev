# Repository Guidelines

## Project Structure & Module Organization

- React + TypeScript sources live in `src/`. `main.tsx` wires Vite’s root, `App.tsx` hosts page layout, and `lib/` holds shared hooks/utilities.
- Styles reside in `App.css` and `index.css`; update them alongside component changes.
- API contracts land in `src/types/` via `npm run generate:api-types`, which shells into `make -C .. api-types` to stay aligned with the backend schema.
- Build artifacts are emitted to `dist/`. Project-level tooling sits beside `vite.config.ts`, `eslint.config.js`, and `tsconfig*.json`.

## Build, Test, and Development Commands

- `npm run dev` ― launches Vite with hot reload (runs `generate:api-types` first).
- `npm run build` ― type-checks with `tsc -b` and produces an optimized Vite build in `dist/`.
- `npm run preview` ― serves the latest build for smoke testing.
- `npm run lint` ― runs ESLint with the repo config, failing on unused disable directives or warnings.
- `npm run generate:api-types` ― regenerates OpenAPI-derived TypeScript models; run this whenever backend contracts move.

## Coding Style & Naming Conventions

- Use modern React function components with hooks; prefer module-scoped helpers over default exports.
- Keep two-space indentation, single quotes in JSX/TSX, and trailing commas where ESLint expects them.
- Components: `PascalCase` (`UserTable.tsx`); utilities: `camelCase`; constants/enums: `SCREAMING_SNAKE_CASE`.
- Never edit files under `src/types/` manually—treat them as generated artifacts.

## Testing Guidelines

- No automated suite ships yet; when adding one, favor Vitest + React Testing Library (Vite-native) and colocate specs as `<module>.test.tsx`.
- Aim for meaningful branch coverage on data loaders and shared hooks; new features should carry smoke tests at minimum.
- Until CI is wired, run `npx vitest run --coverage` locally (add a `test` script mirroring that invocation) before requesting review.

## Commit & Pull Request Guidelines

- History currently only contains `Initial commit`; keep the imperative tone and adopt Conventional Commits (e.g., `feat: add billing dashboard filters`).
- Every PR should explain _what_ changed, _why_, and how reviewers can verify it; link Linear/Jira issues when applicable.
- Attach screenshots or terminal output for UI-visible updates, note any schema regenerations (`generate:api-types`), and ensure CI/lint commands pass locally.

## Security & Configuration Tips

- Store environment-specific values in `.env.local` (ignored by default); never commit secrets.
- Regenerate API types immediately after backend contract bumps to avoid stale type assumptions.
- Audit new dependencies for browser bundle impact; prefer splitting large utilities into `lib/` to keep lazy-loading simple.
