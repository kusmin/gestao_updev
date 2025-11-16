API_SPEC := docs/api.yaml
BACKEND_DIR := backend
FRONTEND_DIR := frontend
BACKOFFICE_DIR := backoffice
TEST_DATABASE_URL ?= postgres://testuser:testpassword@localhost:5433/testdb?sslmode=disable

.PHONY: api-lint api-preview api-types backend-run api-contract-test \
	backend-contract-run backend-migrate backend-test backend-lint backend-build backend-tidy \
	frontend-install frontend-dev frontend-build frontend-preview \
	frontend-lint frontend-test compose-up compose-down compose-logs \
	compose-restart pre-commit-install pre-commit-run pre-commit-update \
	coverage coverage-backend coverage-frontend coverage-backoffice \
	update-deps update-backend-deps update-frontend-deps update-backoffice-deps update-workflow-deps update-tests-deps \
	backend-deps-check backoffice-install backoffice-deps-check backoffice-deps-update backoffice-validate backoffice-lint backoffice-build \
	frontend-deps-check frontend-deps-update frontend-validate \
	deps-check-all deps-update-all validate-all \
	security backend-security frontend-security backoffice-security

.PHONY: swagger
swagger:
	$(MAKE) -C $(BACKEND_DIR) swagger

api-lint:
	pnpm dlx @stoplight/spectral-cli lint $(API_SPEC)

api-preview:
	pnpm dlx @redocly/cli preview-docs $(API_SPEC)

api-types:
	mkdir -p $(FRONTEND_DIR)/src/types
	pnpm dlx openapi-typescript $(API_SPEC) -o $(FRONTEND_DIR)/src/types/api.d.ts
	mkdir -p $(BACKOFFICE_DIR)/src/types
	pnpm dlx openapi-typescript $(API_SPEC) -o $(BACKOFFICE_DIR)/src/types/api.d.ts

backend-run:
	$(MAKE) -C $(BACKEND_DIR) run

backend-contract-run:
	$(MAKE) -C $(BACKEND_DIR) run

backend-migrate:
	$(MAKE) -C $(BACKEND_DIR) migrate

backend-test:
	$(MAKE) -C $(BACKEND_DIR) test

.PHONY: backend-test-docker
backend-test-docker:
	$(MAKE) -C $(BACKEND_DIR) test-docker

backend-lint:
	$(MAKE) -C $(BACKEND_DIR) lint

backend-build:
	$(MAKE) -C $(BACKEND_DIR) build

backend-tidy:
	$(MAKE) -C $(BACKEND_DIR) tidy

api-contract-test:
	( \
	  if [ -f .env.test ]; then \
	    set -a; \
	    . ./.env.test; \
	    set +a; \
	  else \
	    echo ".env.test not found – continuing without injecting extra env vars"; \
	  fi; \
	  DB_URL="$${DATABASE_URL:-$(TEST_DATABASE_URL)}"; \
	  DATABASE_URL=$$DB_URL $(MAKE) backend-migrate; \
	  DATABASE_URL=$$DB_URL pnpm dlx dredd@14 docs/api.yaml http://127.0.0.1:8080 \
	    --hookfiles ./tests/dredd/hooks/basic-flow.js \
	    --server ./scripts/run_dredd_server.sh \
	    --server-wait 5 \
	)

frontend-install:
	pnpm --dir $(FRONTEND_DIR) install

frontend-dev:
	pnpm --dir $(FRONTEND_DIR) run dev

frontend-build:
	pnpm --dir $(FRONTEND_DIR) run build

frontend-preview:
	pnpm --dir $(FRONTEND_DIR) run preview

frontend-lint:
	pnpm --dir $(FRONTEND_DIR) run lint

frontend-test:
	pnpm --dir $(FRONTEND_DIR) run test

coverage: coverage-backend coverage-frontend coverage-backoffice

coverage-backend:
	DATABASE_URL=$(TEST_DATABASE_URL) $(MAKE) -C $(BACKEND_DIR) migrate
	SKIP_AUTO_MIGRATE=1 DATABASE_URL=$(TEST_DATABASE_URL) $(MAKE) -C $(BACKEND_DIR) coverage

coverage-frontend:
	pnpm --dir $(FRONTEND_DIR) run test -- --coverage.enabled true --coverage.reporter=text-summary --coverage.reporter=lcov --coverage.include='src/**/*.{ts,tsx}' --passWithNoTests

coverage-backoffice:
	pnpm --dir $(BACKOFFICE_DIR) run test -- --coverage.enabled true --coverage.reporter=text-summary --coverage.reporter=lcov --coverage.include='src/**/*.{ts,tsx}' --passWithNoTests

compose-up:
	docker compose up --build

compose-down:
	docker compose down --remove-orphans

compose-logs:
	docker compose logs -f

compose-restart:
	docker compose down --remove-orphans
	docker compose up --build

pre-commit-install:
	pre-commit install

pre-commit-run:
	pre-commit run --all-files

pre-commit-update:
	pre-commit autoupdate

update-deps: update-backend-deps update-frontend-deps update-backoffice-deps update-workflow-deps update-tests-deps

update-backend-deps:
	$(MAKE) -C $(BACKEND_DIR) tidy
	cd $(BACKEND_DIR) && go get -u ./... && go mod tidy

update-frontend-deps:
	pnpm --dir $(FRONTEND_DIR) update
	pnpm --dir $(FRONTEND_DIR) install

update-backoffice-deps:
	pnpm --dir $(BACKOFFICE_DIR) update
	pnpm --dir $(BACKOFFICE_DIR) install

update-workflow-deps:
	pnpm update

update-tests-deps:
	pnpm --dir tests/e2e update
	pnpm --dir tests/e2e install
	pnpm --dir tests/postman update
	pnpm --dir tests/postman install

.PHONY: security backend-security frontend-security backoffice-security
security: backend-security frontend-security backoffice-security ## Run all security scanners

backend-security:
	@echo ">>> Running govulncheck..."
	cd $(BACKEND_DIR) && go install golang.org/x/vuln/cmd/govulncheck@latest
	cd $(BACKEND_DIR) && govulncheck ./...
	@echo ">>> Running gosec..."
	cd $(BACKEND_DIR) && go install github.com/securego/gosec/v2/cmd/gosec@latest
	cd $(BACKEND_DIR) && gosec ./...

frontend-security:
	@echo ">>> Running Trivy (frontend)..."
	pnpm dlx --package=trivy trivy fs --exit-code 1 --severity HIGH,CRITICAL $(FRONTEND_DIR)

backoffice-security:
	@echo ">>> Running Trivy (backoffice)..."
	pnpm dlx --package=trivy trivy fs --exit-code 1 --severity HIGH,CRITICAL $(BACKOFFICE_DIR)

##@ Dependencies & Validation
.PHONY: backoffice-install
backoffice-install: ## Install backoffice dependencies
	@echo ">>> Installing backoffice dependencies..."
	pnpm --dir $(BACKOFFICE_DIR) install

.PHONY: backend-deps-check
backend-deps-check: ## Check for backend dependency updates
	@echo ">>> Checking for backend dependency updates..."
	cd $(BACKEND_DIR) && go list -u -m all

.PHONY: backoffice-deps-check
backoffice-deps-check: ## Interactively check for backoffice dependency updates
	@echo ">>> Checking for backoffice dependency updates..."
	pnpm dlx --prefix $(BACKOFFICE_DIR) ncu

.PHONY: backoffice-deps-update
backoffice-deps-update: ## Update backoffice dependencies
	@echo ">>> Updating backoffice dependencies..."
	pnpm dlx --prefix $(BACKOFFICE_DIR) ncu -u
	$(MAKE) backoffice-install

.PHONY: backoffice-validate
backoffice-validate: backoffice-lint backoffice-build ## Validate backoffice project

.PHONY: backoffice-lint
backoffice-lint: ## Lint backoffice project
	@echo ">>> Linting backoffice project..."
	pnpm --dir $(BACKOFFICE_DIR) run lint

.PHONY: backoffice-build
backoffice-build: ## Build backoffice project
	@echo ">>> Building backoffice project..."
	pnpm --dir $(BACKOFFICE_DIR) run build

.PHONY: frontend-deps-check
frontend-deps-check: ## Interactively check for frontend dependency updates
	@echo ">>> Checking for frontend dependency updates..."
	pnpm dlx --prefix $(FRONTEND_DIR) ncu

.PHONY: frontend-deps-update
frontend-deps-update: ## Update frontend dependencies
	@echo ">>> Updating frontend dependencies..."
	pnpm dlx --prefix $(FRONTEND_DIR) ncu -u
	$(MAKE) frontend-install

.PHONY: frontend-validate
frontend-validate: frontend-lint frontend-build ## Validate frontend project

.PHONY: deps-check-all
deps-check-all: backend-deps-check backoffice-deps-check frontend-deps-check ## Check for all dependency updates

.PHONY: deps-update-all
deps-update-all: backoffice-deps-update frontend-deps-update ## Update all dependencies

.PHONY: validate-all
validate-all: backoffice-validate frontend-validate ## Validate all projects
