# GEMINI Project Context: gestao_updev

## Project Overview

This project, `gestao_updev`, is a SaaS platform designed for local businesses like barbershops and clothing stores. The goal is to provide a centralized dashboard for managing clients, appointments, inventory, and sales. The architecture is based on a Go backend and a React/Next.js frontend, with a PostgreSQL database. The system is designed to be multi-tenant.

The project is currently in the planning phase, with the `backend` and `frontend` directories still empty. All the information is based on the documentation in the `docs` directory.

## Building and Running

As the project is in its initial phase, there are no scripts or commands to build or run the application yet. However, based on the documentation, the following commands are expected to be used:

### Backend (Go)

```bash
# TODO: Add build and run commands for the backend.
# Example:
# go build ./cmd/api
# ./api
```

### Frontend (React/Next.js)

```bash
# TODO: Add build and run commands for the frontend.
# Example:
# npm install
# npm run dev
```

## Development Conventions

### General

*   **Language:** All code, comments, and documentation should be in Portuguese.
*   **Formatting:** Use `gofmt` for the backend and `prettier` for the frontend.

### Backend (Go)

*   **Framework:** Gin or Fiber.
*   **Architecture:** Layered architecture (Handler, Service, Repository).
*   **Database:** PostgreSQL with SQLC.
*   **Multi-tenancy:** `tenant_id` column on shared tables.

### Frontend (React/Next.js)

*   **Framework:** Next.js.
*   **UI Kit:** To be defined (Material-UI, Chakra UI, or Tailwind CSS).
*   **State Management:** Zustand or Redux.

### Contribution

*   **Branching:** Create a new branch for each feature (`feature/nome-da-feature`).
*   **Commits:** Use Conventional Commits (`feat:`, `fix:`, `docs:`, etc.).
*   **Pull Requests:** Submit a pull request with a clear description of the changes.
