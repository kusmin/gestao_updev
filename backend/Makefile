PHONY: test
test:
	DATABASE_URL="${DATABASE_URL:-postgres://testuser:testpassword@localhost:5433/testdb?sslmode=disable}" make migrate
	DATABASE_URL="${DATABASE_URL:-postgres://testuser:testpassword@localhost:5433/testdb?sslmode=disable}" go test -p 1 -coverprofile=coverage.out -coverpkg ./... ./...

.PHONY: test-docker
test-docker:
	docker compose -f ../docker-compose.test.yml run --rm backend-test

.PHONY: coverage
coverage:
	docker compose -f ../docker-compose.coverage.yml run --rm backend-coverage

.PHONY: lint
lint:
	golangci-lint run

.PHONY: govulncheck
govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

.PHONY: gosec
gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build:
	go build -o gestao_updev_api cmd/api/main.go

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: swagger
swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs

.PHONY: migrate
migrate:
	env DATABASE_URL="${DATABASE_URL}" go run ./cmd/migrate

.PHONY: seed
seed:
	go run ./cmd/seed
