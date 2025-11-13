.PHONY: test
test:
	go test -coverprofile=coverage.out -coverpkg ./... ./...
	go tool cover -func=coverage.out | tail -n 1

.PHONY: lint
lint:
	golangci-lint run

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: build
build:
	go build -o gestao_updev_api cmd/api/main.go

.PHONY: run
run:
	export APP_ENV=development; \
	export HTTP_PORT=8080; \
	export LOG_LEVEL=info; \
	go run cmd/api/main.go

.PHONY: swagger
swagger:
	go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/api/main.go -o docs

.PHONY: migrate
migrate:
	go run ./cmd/migrate

.PHONY: seed
seed:
	go run ./cmd/seed
