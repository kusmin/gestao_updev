API_SPEC := docs/api.yaml
BACKEND_DIR := backend
FRONTEND_DIR := frontend
BACKOFFICE_DIR := backoffice

.PHONY: api-lint api-preview api-types backend-run api-contract-test \
	backend-contract-run backend-migrate backend-test backend-lint backend-build backend-tidy \
	frontend-install frontend-dev frontend-build frontend-preview \
	frontend-lint frontend-test compose-up compose-down compose-logs \
	compose-restart pre-commit-install pre-commit-run pre-commit-update

.PHONY: swagger
swagger:
	$(MAKE) -C $(BACKEND_DIR) swagger

api-lint:
	npx @stoplight/spectral-cli lint $(API_SPEC)

api-preview:
	npx @redocly/cli preview-docs $(API_SPEC)

api-types:
	mkdir -p $(FRONTEND_DIR)/src/types
	npx openapi-typescript $(API_SPEC) -o $(FRONTEND_DIR)/src/types/api.d.ts
	mkdir -p $(BACKOFFICE_DIR)/src/types
	npx openapi-typescript $(API_SPEC) -o $(BACKOFFICE_DIR)/src/types/api.d.ts

backend-run:
	$(MAKE) -C $(BACKEND_DIR) run

backend-contract-run:
	$(MAKE) -C $(BACKEND_DIR) run

backend-migrate:
	$(MAKE) -C $(BACKEND_DIR) migrate

backend-test:
	$(MAKE) -C $(BACKEND_DIR) test

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
	  $(MAKE) backend-migrate; \
	  npx dredd@14 --config tests/dredd/dredd.yml \
	)

frontend-install:
	npm --prefix $(FRONTEND_DIR) install

frontend-dev:
	npm --prefix $(FRONTEND_DIR) run dev

frontend-build:
	npm --prefix $(FRONTEND_DIR) run build

frontend-preview:
	npm --prefix $(FRONTEND_DIR) run preview

frontend-lint:
	npm --prefix $(FRONTEND_DIR) run lint

frontend-test:
	npm --prefix $(FRONTEND_DIR) run test

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
