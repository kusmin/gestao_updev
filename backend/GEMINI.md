# GEMINI Project Context: gestao_updev (backend)

## Project Overview

This project is the Go backend for `gestao_updev`, a SaaS platform for local businesses. It uses the Gin web framework and is designed with a multi-tenant architecture. The API is documented using Swagger. The logging is handled by the `zap` library, and it uses a standardized JSON response format. The project is containerized using Docker and uses `golangci-lint` for linting. A request ID middleware is used to trace requests.

## Building and Running

### Using Make
The project includes a `Makefile` that simplifies the execution of common tasks:

*   `make test`: Run all tests.
*   `make lint`: Run the linter.
*   `make tidy`: Tidy the `go.mod` file.
*   `make build`: Build the application.
*   `make run`: Run the application.
*   `make swagger`: Generate the Swagger documentation.

### Manual Commands
As the project is in its initial phase, there are no explicit build or run scripts. However, based on the code and Dockerfile, the following commands can be used:

#### Running the application
```bash
go run cmd/api/main.go
```

#### Building the application
```bash
go build -o gestao_updev_api cmd/api/main.go
```

#### Building the Docker image
```bash
docker build -t gestao_updev_api .
```

#### Running the Docker container
```bash
docker run -p 8080:8080 gestao_updev_api
```

#### Running the linter
```bash
golangci-lint run
```

#### Generating Swagger documentation
```bash
swag init
```

## Development Conventions

### General

*   **Language:** All code, comments, and documentation should be in Portuguese.
*   **Formatting:** Use `gofmt` for formatting the code.
*   **Linting:** Use `golangci-lint` to check the code for issues.

### Backend (Go)

*   **Framework:** Gin
*   **Architecture:** Layered architecture (Handler, Service, Repository). The current structure suggests this, but it's not fully implemented.
*   **Multi-tenancy:** The `X-Tenant-ID` header is required for most requests to identify the tenant. Some public routes like authentication and health checks are exempt.
*   **Configuration:** Configuration is managed through environment variables.
*   **Logging:** The `zap` library is used for logging.
*   **API Responses:** The API uses a standardized JSON response format for successes and errors.
*   **Request Tracing:** A middleware injects a unique `X-Request-ID` header into each request.
*   **API Documentation:** The API is documented using Swagger. The documentation is generated using `swaggo/swag`.
*   **Dependencies:** The project uses Go modules to manage dependencies. The `go.sum` file contains the checksums of the dependencies.

### Contribution

*   **Branching:** Create a new branch for each feature (`feature/nome-da-feature`).
*   **Commits:** Use Conventional Commits (`feat:`, `fix:`, `docs:`, etc.).
*   **Pull Requests:** Submit a pull request with a clear description of the changes.
