.PHONY: test
test:
	go test ./...

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
	go run cmd/api/main.go

.PHONY: swagger
swagger:
	swag init
